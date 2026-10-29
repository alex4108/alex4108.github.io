# schittko.me

Alex Schittko's Jekyll site, retaining the original dark theme, blog list, and complete technology notes on About. Writing stays at `/archives/`; published article URLs and `CNAME` are unchanged. Pages contain no analytics, remote fonts, or executable client-side JavaScript.

## Local build and checks

Ruby and Bundler are required. The bundle pins GitHub Pages 232 / Jekyll 3.10.0. Python 3 runs the existing link checker; Go 1.18+ runs the additional regression checks. Fetch `origin/master` before checking preserved source content against it.

Keep dependencies and generated output outside the checkout, including the tracked legacy `_site` files:

```sh
export BUNDLE_PATH="$(mktemp -d)"
export BUNDLE_APP_CONFIG="$BUNDLE_PATH/config"
export SITE_OUTPUT="$(mktemp -d)"
bundle install
TZ=UTC JEKYLL_ENV=production bundle exec jekyll build --safe --strict_front_matter --destination "$SITE_OUTPUT"
python3 scripts/check_site.py "$SITE_OUTPUT"
GOMAXPROCS=2 go test scripts/site_regression_test.go -v -args -source "$PWD" -site "$SITE_OUTPUT" -baseline origin/master
GOMAXPROCS=2 go vet scripts/site_regression_test.go
```

The unchanged Python checker validates internal anchor/source links, fragments, and all five published article URLs, and checks that feed/robots/sitemap files exist. The Go checks additionally parse RSS and sitemap XML, verify their exact public URL sets, protect original post/CNAME bytes and the complete conservative About text, and check the rendered theme, Writing labels, titles, canonical links, SEO metadata, and preview indexing.

Both layouts explicitly escape `site.title` in the HTML title and header. A page title is prepended when present; `seo title=false` retains one SEO metadata block without generating a second title. To check escaping, build to another external destination with an additional Jekyll configuration (`--config _config.yml,/absolute/path/to/title.yml`) containing:

```yaml
title: 'R&D <Engineering> "Alex"'
```

Run the same Go checks against that destination with `-title 'R&D <Engineering> "Alex"'`. They also verify the RSS channel/discovery title and the header use the override.

The retained `/preview/` pages remain `noindex` and are excluded from the sitemap. To preview locally:

```sh
bundle exec jekyll serve --destination "$SITE_OUTPUT"
```
