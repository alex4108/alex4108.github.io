#!/usr/bin/env python3
"""Check generated internal links and preserve the site's published article URLs."""

from __future__ import annotations

import argparse
from html.parser import HTMLParser
from pathlib import Path
from urllib.parse import unquote, urljoin, urlparse

SITE_ORIGIN = "https://schittko.me"
PUBLISHED_ARTICLE_PATHS = (
    "/2021/02/04/How-we-used-xcp-ng-to-run-Rosetta-@-home/",
    "/2021/02/04/SQL-Server-4600-DBs/",
    "/2021/02/04/Welcome/",
    "/2021/03/12/Azure-SQL-Iaas-Public/",
    "/2021/12/04/My-First-Hackathon/",
)


class DocumentParser(HTMLParser):
    def __init__(self) -> None:
        super().__init__(convert_charrefs=True)
        self.ids: set[str] = set()
        self.links: list[tuple[str, str, set[str]]] = []

    def handle_starttag(self, tag: str, attrs: list[tuple[str, str | None]]) -> None:
        values = {key: value or "" for key, value in attrs}
        if values.get("id"):
            self.ids.add(values["id"])
        if tag == "a" and values.get("name"):
            self.ids.add(values["name"])
        if tag == "a" and values.get("href"):
            self.links.append(("href", values["href"], set(values.get("rel", "").split())))
        for attribute in ("src",):
            if values.get(attribute):
                self.links.append((attribute, values[attribute], set()))


def output_for_url(root: Path, path: str) -> Path | None:
    relative = unquote(path).lstrip("/")
    candidate = root / relative
    options = [candidate]
    if path.endswith("/") or not relative:
        options.insert(0, candidate / "index.html")
    elif not Path(relative).suffix:
        options.extend((root / f"{relative}.html", candidate / "index.html"))
    return next((item for item in options if item.is_file()), None)


def source_url(root: Path, source: Path) -> str:
    relative = source.relative_to(root).as_posix()
    if relative == "index.html":
        return "/"
    if relative.endswith("/index.html"):
        return f"/{relative[:-10]}"
    return f"/{relative}"


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("site", nargs="?", default="_site", type=Path)
    args = parser.parse_args()
    root = args.site.resolve()
    if not root.is_dir():
        raise SystemExit(f"Generated site not found: {root}")

    documents: dict[Path, DocumentParser] = {}
    for document in root.rglob("*.html"):
        parsed = DocumentParser()
        parsed.feed(document.read_text(encoding="utf-8"))
        documents[document.resolve()] = parsed

    errors: list[str] = []
    checked = 0
    for source, document in documents.items():
        base = f"{SITE_ORIGIN}{source_url(root, source)}"
        for attribute, raw_link, rel in document.links:
            if raw_link.startswith(("mailto:", "tel:", "data:", "javascript:")):
                continue
            resolved = urlparse(urljoin(base, raw_link))
            if resolved.scheme not in ("http", "https") or resolved.netloc != "schittko.me":
                continue
            target = output_for_url(root, resolved.path)
            checked += 1
            if target is None:
                errors.append(f"{source.relative_to(root)}: missing {attribute} target {raw_link}")
                continue
            if resolved.fragment and target.suffix == ".html":
                target_parser = documents.get(target.resolve())
                if target_parser and unquote(resolved.fragment) not in target_parser.ids:
                    errors.append(f"{source.relative_to(root)}: missing fragment in {raw_link}")

    for published_path in PUBLISHED_ARTICLE_PATHS:
        if output_for_url(root, published_path) is None:
            errors.append(f"published article URL was not generated: {published_path}")

    for required in ("robots.txt", "sitemap.xml", "feed.xml"):
        if not (root / required).is_file():
            errors.append(f"required output was not generated: /{required}")

    if errors:
        print("Site checks failed:")
        for error in errors:
            print(f"  - {error}")
        return 1

    print(
        f"Site checks passed: {checked} internal references across "
        f"{len(documents)} HTML files; {len(PUBLISHED_ARTICLE_PATHS)} published article URLs preserved."
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
