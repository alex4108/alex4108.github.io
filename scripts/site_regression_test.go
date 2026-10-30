package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"html"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

var sourceRoot = flag.String("source", ".", "Jekyll source directory")
var baseline = flag.String("baseline", "origin/master", "Git baseline for preserved content")
var siteRoot = flag.String("site", "", "Production Jekyll output directory (required for rendered checks)")
var siteTitle = flag.String("title", "Alex Schittko", "Expected rendered site.title")

const origin = "https://schittko.me"

var pages = []struct{ path, title string }{
	{"index.html", ""},
	{"about/index.html", "About"},
	{"archives/index.html", "Writing"},
	{"work/index.html", "Work"},
	{"2021/02/04/Welcome/index.html", "Welcome!"},
	{"2021/02/04/How-we-used-xcp-ng-to-run-Rosetta-@-home/index.html", "How we used XCP-ng to run Rosetta@Home"},
	{"2021/02/04/SQL-Server-4600-DBs/index.html", "Alex, how do we save money on our SQL servers?"},
	{"2021/03/12/Azure-SQL-Iaas-Public/index.html", "Public access to an Azure SQL IaaS Failover Cluster Instance"},
	{"2021/12/04/My-First-Hackathon/index.html", "My First Hackathon"},
	{"preview/index.html", ""},
	{"preview/about/index.html", "About"},
	{"preview/writing/index.html", "Writing"},
	{"preview/projects/index.html", ""},
	{"preview/case-studies/index.html", ""},
}

var attributePattern = regexp.MustCompile(`([A-Za-z_:][-A-Za-z0-9_:.]*)\s*=\s*(?:"([^"]*)"|'([^']*)')`)
var tagPattern = regexp.MustCompile(`(?s)<[^>]*>`)

func readFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(content)
}

func gitFile(t *testing.T, path string) string {
	t.Helper()
	output, err := exec.Command("git", "-C", *sourceRoot, "show", *baseline+":"+path).CombinedOutput()
	if err != nil {
		t.Fatalf("read %s:%s: %v\n%s", *baseline, path, err, output)
	}
	return string(output)
}

func markdownBody(t *testing.T, content string) string {
	t.Helper()
	header := regexp.MustCompile(`(?s)^---\n.*?\n---\n`).FindStringIndex(content)
	if header == nil {
		t.Fatal("missing YAML front matter")
	}
	return strings.TrimSpace(content[header[1]:])
}

func renderedFile(t *testing.T, path string) string {
	t.Helper()
	if *siteRoot == "" {
		t.Fatal("rendered checks require -site pointing to an actual production Jekyll build")
	}
	return readFile(t, filepath.Join(*siteRoot, filepath.FromSlash(path)))
}

func attributes(tag string) map[string]string {
	values := make(map[string]string)
	for _, match := range attributePattern.FindAllStringSubmatch(tag, -1) {
		values[strings.ToLower(match[1])] = html.UnescapeString(match[2] + match[3])
	}
	return values
}

func visibleText(content string) string {
	return strings.Join(strings.Fields(html.UnescapeString(tagPattern.ReplaceAllString(content, ""))), " ")
}

func pageURL(path string) string {
	return origin + "/" + strings.TrimSuffix(path, "index.html")
}

func decodeXML(t *testing.T, path string, value any) {
	t.Helper()
	decoder := xml.NewDecoder(strings.NewReader(renderedFile(t, path)))
	if err := decoder.Decode(value); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	for {
		token, err := decoder.Token()
		if err == io.EOF {
			return
		}
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		if text, ok := token.(xml.CharData); ok && strings.TrimSpace(string(text)) == "" {
			continue
		}
		t.Fatalf("%s: unexpected content after XML document", path)
	}
}

func TestExplicitSiteTitle(t *testing.T) {
	for _, name := range []string{"default", "preview"} {
		t.Run(name, func(t *testing.T) {
			layout := readFile(t, filepath.Join(*sourceRoot, "_layouts", name+".html"))
			titles := regexp.MustCompile(`(?s)<title>.*?</title>`).FindAllString(layout, -1)
			if len(titles) != 1 {
				t.Fatalf("want one explicit site.title-aware title, got %d", len(titles))
			}
			if !strings.Contains(titles[0], "{{ site.title | escape }}") {
				t.Error("explicit title must HTML-escape site.title")
			}
			seo := regexp.MustCompile(`{%\s*seo\b.*?%}`).FindAllString(layout, -1)
			if len(seo) != 1 || !strings.Contains(seo[0], "title=false") {
				t.Errorf("want one SEO block with its title disabled, got %q", seo)
			}
		})
	}
}

func TestOriginalBlogPresentation(t *testing.T) {
	index := readFile(t, filepath.Join(*sourceRoot, "index.html"))
	for _, want := range []string{"for post in paginator.posts", `<article class="post">`, "post.description", "Read more here"} {
		if !strings.Contains(index, want) {
			t.Errorf("homepage lost original blog-list feature %q", want)
		}
	}

	for _, unwanted := range []string{`class="hero"`, `class="proof"`, `class="card-grid"`, "Operating principles"} {
		if strings.Contains(index, unwanted) {
			t.Errorf("homepage still contains redesign feature %q", unwanted)
		}
	}

	css := readFile(t, filepath.Join(*sourceRoot, "assets", "css", "style.css"))
	for _, selector := range []string{"#page", "#header", "#content", "#footer"} {
		rule := regexp.MustCompile(regexp.QuoteMeta(selector) + `\s*\{([^}]+)}`).FindStringSubmatch(css)
		if len(rule) != 2 || !regexp.MustCompile(`width:\s*500px;`).MatchString(rule[1]) {
			t.Errorf("%s must retain the original 500px desktop layout", selector)
		}
	}
	for _, want := range []string{"font-family: Verdana, serif;", "font-size: 13px;", "left: -240px;", "(max-width: 870px)"} {
		if !strings.Contains(css, want) {
			t.Errorf("stylesheet lost original theme rule %q", want)
		}
	}
	for _, unwanted := range []string{".hero", ".proof-strip", ".card-grid", ".work-card", ".project-list"} {
		if strings.Contains(css, unwanted) {
			t.Errorf("stylesheet still contains redesign selector %q", unwanted)
		}
	}
}

func TestHistoricalListMarkers(t *testing.T) {
	for _, article := range []struct {
		path         string
		lists, steps int
	}{
		{"2021/02/04/How-we-used-xcp-ng-to-run-Rosetta-@-home/index.html", 1, 6},
		{"2021/02/04/SQL-Server-4600-DBs/index.html", 5, 11},
	} {
		t.Run(article.path, func(t *testing.T) {
			ordered := regexp.MustCompile(`(?s)<ol>.*</ol>`).FindString(renderedFile(t, article.path))
			if strings.Count(ordered, "<ol>") != article.lists || strings.Count(ordered, "</ol>") != article.lists ||
				strings.Count(ordered, "<li>") != article.steps || strings.Count(ordered, "</li>") != article.steps {
				t.Fatalf("historical procedure must retain %d ordered lists and %d steps", article.lists, article.steps)
			}
			t.Logf("historical ordered markup: %d lists, %d steps", article.lists, article.steps)
		})
	}

	for _, stylesheet := range []struct{ name, content string }{
		{"source", readFile(t, filepath.Join(*sourceRoot, "assets", "css", "style.css"))},
		{"rendered", renderedFile(t, "assets/css/style.css")},
	} {
		for _, marker := range []struct{ list, style string }{{"ol", "decimal"}, {"ul", "disc"}} {
			t.Run(stylesheet.name+"/"+marker.list, func(t *testing.T) {
				rule := `(?m)^` + marker.list + `\s*\{\s*list-style:\s*` + marker.style + `;\s*\}`
				if !regexp.MustCompile(rule).MatchString(stylesheet.content) {
					t.Errorf("%s stylesheet must explicitly style %s as %s in its own rule", stylesheet.name, marker.list, marker.style)
				}
			})
		}
	}
}

func TestWritingLabels(t *testing.T) {
	oldLabel := regexp.MustCompile(`(?i)Writing\s*/\s*Archive|>\s*Archives?\s*<|title:\s*Archives?\b|writing archive`)
	for _, path := range []string{"_layouts/default.html", "_layouts/preview.html", "index.html", "archives/index.html", "work/index.html", "preview/index.html", "preview/writing/index.html"} {
		t.Run(path, func(t *testing.T) {
			content := readFile(t, filepath.Join(*sourceRoot, filepath.FromSlash(path)))
			if oldLabel.MatchString(content) {
				t.Error("public section labels must say Writing, not Archive or Writing / Archive")
			}
			if !strings.Contains(content, "Writing") {
				t.Error("missing visible Writing label")
			}
		})
	}
	writing := readFile(t, filepath.Join(*sourceRoot, "archives", "index.html"))
	if !strings.Contains(writing, "\ntitle: Writing\n") || !strings.Contains(writing, "\npermalink: /archives/\n") {
		t.Error("Writing must retain its title and /archives/ URL")
	}
}

func TestConservativeAbout(t *testing.T) {
	about := readFile(t, filepath.Join(*sourceRoot, "about", "index.markdown"))
	want := markdownBody(t, gitFile(t, "preview/about/index.markdown"))
	got := markdownBody(t, about)
	if got != want {
		t.Errorf("About must preserve the complete conservative preview body: got %d bytes, want %d", len(got), len(want))
	}
	if !strings.Contains(about, "\nlayout: page\n") || !strings.Contains(about, "\npermalink: /about/\n") || strings.Contains(about, "\nsitemap: false\n") {
		t.Error("public About must use its public layout, URL and sitemap visibility")
	}
}

func TestWorkProvenance(t *testing.T) {
	work := readFile(t, filepath.Join(*sourceRoot, "work", "index.html"))
	for _, project := range []string{"ha-kbeacon-ring", "AstroInstaller", "meshfleet", "discord-photo-reaper", "bitbucket-azuredevops-migrate", "jellycleaner", "Approova"} {
		if !strings.Contains(work, "https://github.com/alex4108/"+project) {
			t.Errorf("lost public project link %s", project)
		}
	}
	for _, metric := range []string{"In 2021", "4,600", "1.25 TB", "90-minute", "30%"} {
		if !strings.Contains(work, metric) {
			t.Errorf("missing historical migration context/result %q", metric)
		}
	}
	for _, unsupported := range []string{"challenge-response", "Bluetooth proxies", "dry runs", "rate-limit backoff", "unattended installation", "guarded cleanup", "A Go application", "assign access"} {
		if strings.Contains(work, unsupported) {
			t.Errorf("project claim lacks original published provenance: %q", unsupported)
		}
	}
	for _, redesign := range []string{"proof-strip", "project-list", "case-study", "eyebrow", "section-heading"} {
		if strings.Contains(work, `class="`+redesign) {
			t.Errorf("Work must use the original theme, not redesign class %q", redesign)
		}
	}
}

func TestPreviewSafety(t *testing.T) {
	layout := readFile(t, filepath.Join(*sourceRoot, "_layouts", "preview.html"))
	if !strings.Contains(layout, `<meta name="robots" content="noindex, follow">`) {
		t.Error("retained preview must remain noindex")
	}
	if regexp.MustCompile(`(?i)<script\b`).MatchString(layout) {
		t.Error("retained preview must not reintroduce legacy executable JavaScript")
	}
	for _, path := range []string{"preview/index.html", "preview/about/index.markdown", "preview/writing/index.html", "preview/projects/index.html", "preview/case-studies/index.html"} {
		if !strings.Contains(readFile(t, filepath.Join(*sourceRoot, filepath.FromSlash(path))), "\nsitemap: false\n") {
			t.Errorf("%s must remain excluded from sitemap", path)
		}
	}
}

func TestRenderedPages(t *testing.T) {
	titleTags := regexp.MustCompile(`(?is)<title\b[^>]*>(.*?)</title\s*>`)
	entities := regexp.MustCompile(`&(?:amp|lt|gt|quot|apos|#[0-9]+|#x[0-9A-Fa-f]+);`)
	headTags := regexp.MustCompile(`(?is)<(?:meta|link)\b[^>]*>`)
	scripts := regexp.MustCompile(`(?is)<script\b([^>]*)>(.*?)</script\s*>`)
	header := regexp.MustCompile(`(?is)<header\b[^>]*id="header"[^>]*>(.*?)</header>`)
	navigation := regexp.MustCompile(`(?is)<nav\b[^>]*aria-label="Primary navigation"[^>]*>(.*?)</nav>`)
	anchors := regexp.MustCompile(`(?is)<a\b[^>]*>(.*?)</a>`)

	for _, page := range pages {
		t.Run(page.path, func(t *testing.T) {
			content := renderedFile(t, page.path)
			wantTitle := *siteTitle
			if page.title != "" && page.title != *siteTitle {
				wantTitle = page.title + " | " + wantTitle
			}
			titles := titleTags.FindAllStringSubmatch(content, -1)
			if len(titles) != 1 {
				t.Fatalf("want exactly one rendered title, got %d", len(titles))
			}
			if html.UnescapeString(titles[0][1]) != wantTitle || strings.ContainsAny(entities.ReplaceAllString(titles[0][1], ""), "&<>\"'") {
				t.Errorf("want escaped title %q, got %q", wantTitle, titles[0][1])
			}
			t.Logf("title=%q raw=%q", html.UnescapeString(titles[0][1]), titles[0][1])
			if strings.Count(content, "<!-- Begin Jekyll SEO tag ") != 1 || strings.Count(content, "<!-- End Jekyll SEO tag -->") != 1 {
				t.Error("want exactly one complete SEO metadata block")
			}

			metadata := make(map[string][]string)
			canonical, rss := 0, 0
			styles := make(map[string]int)
			for _, tag := range headTags.FindAllString(content, -1) {
				attr := attributes(tag)
				if strings.HasPrefix(strings.ToLower(tag), "<meta") {
					key := attr["name"]
					if key == "" {
						key = attr["property"]
					}
					metadata[key] = append(metadata[key], attr["content"])
				}
				switch attr["rel"] {
				case "canonical":
					canonical++
					if attr["href"] != pageURL(page.path) {
						t.Errorf("wrong canonical URL %q", attr["href"])
					}
				case "alternate":
					if attr["type"] == "application/rss+xml" {
						rss++
						if attr["href"] != "/feed.xml" || attr["title"] != *siteTitle {
							t.Errorf("RSS discovery lost dynamic title or URL: %v", attr)
						}
					}
				case "stylesheet":
					styles[attr["href"]]++
				}
			}
			if canonical != 1 || rss != 1 || len(styles) != 2 || styles["/assets/css/style.css"] != 1 || styles["/assets/css/colors-dark.css"] != 1 {
				t.Errorf("invalid canonical/RSS/original-theme links: canonical=%d rss=%d styles=%v", canonical, rss, styles)
			}
			for _, key := range []string{"generator", "author", "description", "og:title", "og:locale", "og:description", "og:url", "og:site_name", "og:type", "twitter:card", "twitter:title"} {
				if len(metadata[key]) != 1 || metadata[key][0] == "" {
					t.Errorf("missing or duplicate SEO metadata %q: %v", key, metadata[key])
				}
			}
			if len(metadata["generator"]) == 1 && metadata["generator"][0] != "Jekyll v3.10.0" {
				t.Errorf("not the Pages 232 Jekyll build: %v", metadata["generator"])
			}
			if len(metadata["og:url"]) == 1 && metadata["og:url"][0] != pageURL(page.path) {
				t.Errorf("wrong Open Graph URL: %v", metadata["og:url"])
			}
			robots := strings.ToLower(strings.Join(metadata["robots"], ","))
			if len(metadata["robots"]) > 1 || strings.HasPrefix(page.path, "preview/") != strings.Contains(robots, "noindex") {
				t.Errorf("incorrect production/preview indexing: %v", metadata["robots"])
			}

			jsonScripts := scripts.FindAllStringSubmatch(content, -1)
			if len(jsonScripts) != 1 || attributes(jsonScripts[0][1])["type"] != "application/ld+json" {
				t.Fatalf("want one JSON-LD script and no executable JavaScript, got %d scripts", len(jsonScripts))
			}
			var schema struct {
				Context string `json:"@context"`
				Type    string `json:"@type"`
				URL     string `json:"url"`
			}
			if err := json.Unmarshal([]byte(jsonScripts[0][2]), &schema); err != nil {
				t.Fatal(err)
			}
			if schema.Context != "https://schema.org" || schema.Type == "" || schema.URL != pageURL(page.path) {
				t.Errorf("invalid SEO JSON-LD: %+v", schema)
			}

			headerMatch := header.FindStringSubmatch(content)
			if len(headerMatch) != 2 {
				t.Fatal("missing original site header")
			}
			headerLink := anchors.FindStringSubmatch(headerMatch[1])
			if len(headerLink) != 2 || visibleText(headerLink[1]) != *siteTitle {
				t.Errorf("header lost escaped dynamic site.title %q", *siteTitle)
			}
			navMatch := navigation.FindStringSubmatch(content)
			if len(navMatch) != 2 {
				t.Fatal("missing accessible primary navigation")
			}
			writingLinks := 0
			for _, anchor := range anchors.FindAllStringSubmatch(navMatch[1], -1) {
				if attributes(anchor[0])["href"] == "/archives/" && visibleText(anchor[1]) == "Writing" {
					writingLinks++
				}
			}
			if writingLinks != 1 || strings.Contains(visibleText(navMatch[1]), "Archive") {
				t.Error("primary navigation must label /archives/ as Writing exactly once")
			}
		})
	}
}

func TestRenderedBlogAndAssets(t *testing.T) {
	for _, path := range []string{"index.html", "archives/index.html"} {
		content := renderedFile(t, path)
		if strings.Count(content, `<article class="post">`) != 5 {
			t.Errorf("%s must retain all five original-style post summaries", path)
		}
		for _, page := range pages {
			if strings.HasPrefix(page.path, "2021/") && !strings.Contains(content, `href="/`+strings.TrimSuffix(page.path, "index.html")+`"`) {
				t.Errorf("%s lost published article %s", path, page.path)
			}
		}
	}
	about := renderedFile(t, "about/index.html")
	headings := regexp.MustCompile(`(?is)<h[1-6]\b[^>]*>(.*?)</h[1-6]>`)
	renderedHeadings := make(map[string]bool)
	for _, heading := range headings.FindAllStringSubmatch(about, -1) {
		renderedHeadings[visibleText(heading[1])] = true
	}
	for _, heading := range regexp.MustCompile(`(?m)^#{1,6}\s+(.+)$`).FindAllStringSubmatch(gitFile(t, "preview/about/index.markdown"), -1) {
		if !renderedHeadings[heading[1]] {
			t.Errorf("About lost rendered technology/personal section %q", heading[1])
		}
	}
	for _, path := range []string{"assets/css/style.css", "assets/css/colors-dark.css"} {
		if renderedFile(t, path) != readFile(t, filepath.Join(*sourceRoot, filepath.FromSlash(path))) {
			t.Errorf("rendered theme asset differs from source: %s", path)
		}
	}
	css := renderedFile(t, "assets/css/style.css")
	for _, match := range regexp.MustCompile(`url\(['"]?([^'")\s]+)['"]?\)`).FindAllStringSubmatch(css, -1) {
		reference, err := url.Parse(match[1])
		if err != nil {
			t.Fatal(err)
		}
		target := filepath.Clean(filepath.Join("assets/css", filepath.FromSlash(reference.Path)))
		if reference.IsAbs() || reference.Host != "" || filepath.IsAbs(reference.Path) || target == ".." || strings.HasPrefix(target, ".."+string(filepath.Separator)) {
			t.Fatalf("theme asset is not a local in-site reference: %q", match[1])
		}
		if renderedFile(t, target) == "" {
			t.Errorf("empty theme asset %s", target)
		}
	}
	palette := renderedFile(t, "assets/css/colors-dark.css")
	for _, color := range []string{"#002a35", "#839496", "#2aa198", "#073642", "#cb4b16"} {
		if !strings.Contains(palette, color) {
			t.Errorf("lost original dark-palette color %s", color)
		}
	}
}

func TestFeedAndSitemapXML(t *testing.T) {
	postURLs, siteURLs := make(map[string]bool), make(map[string]bool)
	for _, page := range pages {
		if strings.HasPrefix(page.path, "2021/") {
			postURLs[pageURL(page.path)] = true
		}
		if !strings.HasPrefix(page.path, "preview/") {
			siteURLs[pageURL(page.path)] = true
		}
	}
	var feed struct {
		XMLName xml.Name `xml:"rss"`
		Version string   `xml:"version,attr"`
		Channel struct {
			Title string `xml:"title"`
			Links []struct {
				XMLName xml.Name
				Text    string `xml:",chardata"`
				Href    string `xml:"href,attr"`
				Rel     string `xml:"rel,attr"`
				Type    string `xml:"type,attr"`
			} `xml:"link"`
			Items []struct {
				Title       string `xml:"title"`
				Description string `xml:"description"`
				Link        string `xml:"link"`
				PubDate     string `xml:"pubDate"`
				GUID        struct {
					Value     string `xml:",chardata"`
					Permalink string `xml:"isPermaLink,attr"`
				} `xml:"guid"`
			} `xml:"item"`
		} `xml:"channel"`
	}
	decodeXML(t, "feed.xml", &feed)
	if feed.Version != "2.0" || feed.Channel.Title != *siteTitle || len(feed.Channel.Items) != 5 || len(feed.Channel.Links) != 2 {
		t.Fatalf("invalid RSS version/title/item/link counts: version=%s title=%q items=%d links=%d", feed.Version, feed.Channel.Title, len(feed.Channel.Items), len(feed.Channel.Links))
	}
	homepageLinks, selfLinks := 0, 0
	for _, link := range feed.Channel.Links {
		switch link.XMLName.Space {
		case "":
			homepageLinks++
			if link.Text != origin+"/" {
				t.Errorf("wrong RSS homepage %q", link.Text)
			}
		case "http://www.w3.org/2005/Atom":
			selfLinks++
			if link.Href != origin+"/feed.xml" || link.Rel != "self" || link.Type != "application/rss+xml" {
				t.Errorf("invalid RSS self link %+v", link)
			}
		default:
			t.Errorf("unexpected RSS link namespace %q", link.XMLName.Space)
		}
	}
	if homepageLinks != 1 || selfLinks != 1 {
		t.Errorf("want one RSS homepage and one Atom self link, got %d and %d", homepageLinks, selfLinks)
	}
	for _, item := range feed.Channel.Items {
		if !postURLs[item.Link] || item.GUID.Value != item.Link || item.GUID.Permalink != "true" || item.Title == "" || item.Description == "" {
			t.Errorf("invalid, missing or duplicate RSS article %q", item.Link)
		}
		delete(postURLs, item.Link)
		if _, err := time.Parse(time.RFC1123Z, item.PubDate); err != nil {
			t.Errorf("invalid RSS article date: %v", err)
		}
	}
	var sitemap struct {
		XMLName xml.Name `xml:"http://www.sitemaps.org/schemas/sitemap/0.9 urlset"`
		URLs    []struct {
			Loc     string `xml:"loc"`
			Lastmod string `xml:"lastmod"`
		} `xml:"url"`
	}
	decodeXML(t, "sitemap.xml", &sitemap)
	if len(sitemap.URLs) != len(siteURLs) {
		t.Errorf("want %d public sitemap URLs, got %d", len(siteURLs), len(sitemap.URLs))
	}
	for _, entry := range sitemap.URLs {
		if !siteURLs[entry.Loc] {
			t.Errorf("invalid, duplicate or preview sitemap URL %q", entry.Loc)
		}
		delete(siteURLs, entry.Loc)
		if entry.Lastmod != "" {
			if _, err := time.Parse(time.RFC3339, entry.Lastmod); err != nil {
				t.Errorf("invalid sitemap lastmod: %v", err)
			}
		}
	}
	if len(postURLs) != 0 || len(siteURLs) != 0 {
		t.Errorf("missing RSS/sitemap URLs: posts=%v site=%v", postURLs, siteURLs)
	}
	if !strings.Contains(renderedFile(t, "robots.txt"), "Sitemap: "+origin+"/sitemap.xml") {
		t.Error("robots.txt lost the canonical sitemap URL")
	}
}

func TestProtectedSourceBytes(t *testing.T) {
	output, err := exec.Command("git", "-C", *sourceRoot, "ls-tree", "-r", "--name-only", *baseline, "--", "_posts").CombinedOutput()
	if err != nil {
		t.Fatalf("list baseline posts: %v\n%s", err, output)
	}
	posts := strings.Fields(string(output))
	entries, err := os.ReadDir(filepath.Join(*sourceRoot, "_posts"))
	if err != nil {
		t.Fatal(err)
	}
	if len(posts) != 5 || len(entries) != 5 {
		t.Fatalf("must retain exactly five historical source posts: baseline=%d source=%d", len(posts), len(entries))
	}
	for _, path := range append(posts, "CNAME", "preview/about/index.markdown") {
		if readFile(t, filepath.Join(*sourceRoot, filepath.FromSlash(path))) != gitFile(t, path) {
			t.Errorf("%s must remain byte-identical to %s", path, *baseline)
		}
	}
	if renderedFile(t, "CNAME") != gitFile(t, "CNAME") {
		t.Error("rendered CNAME differs from the original")
	}
	for _, path := range []string{"scripts", "README.md", "Gemfile", "Gemfile.lock"} {
		_, err := os.Stat(filepath.Join(*siteRoot, path))
		if err == nil {
			t.Errorf("private build input leaked into public output: %s", path)
		} else if !os.IsNotExist(err) {
			t.Fatal(err)
		}
	}
}
