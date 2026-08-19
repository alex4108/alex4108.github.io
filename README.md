# schittko.me

Alex Schittko's Jekyll site. The site keeps the original theme and builds to static HTML without analytics, remote fonts, or client-side JavaScript.

## Local build and checks

Ruby and Bundler are required.

```sh
bundle install
JEKYLL_ENV=production bundle exec jekyll build --strict_front_matter
python3 scripts/check_site.py _site
```

The checker validates generated internal links and fragments, the RSS/robots/sitemap outputs, and every previously published article URL. To preview locally:

```sh
bundle exec jekyll serve
```
