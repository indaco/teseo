<h1 align="center">
  teseo
</h1>

<h2 align="center" style="font-size: 1.5em;">
  Go templ components for SEO.
</h2>

<p align="center">
  <a href="https://github.com/indaco/teseo/actions/workflows/ci.yml" target="_blank">
    <img src="https://github.com/indaco/teseo/actions/workflows/ci.yml/badge.svg" alt="CI" />
  </a>
  <a href="https://codecov.io/gh/indaco/teseo" target="_blank">
    <img src="https://codecov.io/gh/indaco/teseo/branch/main/graph/badge.svg" alt="Code coverage" />
  </a>
  <a href="https://goreportcard.com/report/github.com/indaco/teseo" target="_blank">
    <img src="https://goreportcard.com/badge/github.com/indaco/teseo" alt="Go Report Card" />
  </a>
  <a href="https://github.com/indaco/teseo/actions/workflows/security.yml" target="_blank">
    <img src="https://github.com/indaco/teseo/actions/workflows/security.yml/badge.svg" alt="Security Scan" />
  </a>
  <a href="https://github.com/indaco/teseo/releases" target="_blank">
    <img src="https://img.shields.io/github/v/tag/indaco/teseo?label=version&sort=semver&color=4c1" alt="version">
  </a>
  <a href="https://pkg.go.dev/github.com/indaco/teseo" target="_blank">
    <img src="https://pkg.go.dev/badge/github.com/indaco/teseo.svg" alt="Go Reference" />
  </a>
  <a href="LICENSE" target="_blank">
    <img src="https://img.shields.io/badge/license-mit-blue?style=flat-square" alt="License" />
  </a>
  <a href="https://www.jetify.com/devbox" target="_blank">
    <img src="https://www.jetify.com/img/devbox/shield_moon.svg" alt="Built with Devbox" />
  </a>
</p>

**teseo** provides a rich set of SEO-focused Go structs that follow **Schema.org** and **OpenGraph** specifications, with helpers to render structured data using either [templ](https://github.com/a-h/templ) components or Go’s built-in `html/template`.

Whether you're working with _Schema.org JSON-LD_, _OpenGraph_, or _Twitter Cards_, `teseo` simplifies the process of adding SEO metadata to your web applications.

## Features

- Complete support for **Schema.org JSON-LD** types
- Built-in types for **OpenGraph** meta tags
- Easy-to-generate **Twitter Card** metadata
- Dual rendering support: `templ` components or Go `template/html`
- Developer-friendly API with helpers and factory methods

## Installation

Add this package to your project:

```bash
go get github.com/indaco/teseo@latest
```

## Supported Data Types

### Schema.org JSON-LD

- Article
- BreadcrumbList
- Event
- FAQPage
- LocalBusiness
- Organization
- Person
- Product
- SiteNavigationElement
- WebPage
- WebSite

### OpenGraph

- Article
- Audio
- Book
- Business
- Event
- MusicAlbum
- MusicPlaylist
- MusicSong
- MusicRadioStation
- Place
- Profile
- Product
- ProductGroup
- Restaurant
- Video
- VideoEpisode
- VideoMovie
- Website

### Twitter Cards

- Summary Card
- Summary with Large Image
- App Card
- Player Card

## Usage

### Schema.org JSON-LD

Each entity includes:

- `ToJsonLd()` → renders a `templ.Component`
- `ToGoHTMLJsonLd()` → returns a `template.HTML string`

You can define data using plain structs or with provided factory functions.

#### Example: WebPage

```templ
package pages

import "github.com/indaco/teseo/schemaorg"

templ HomePage() {
 {{
    webpage := &schemaorg.WebPage{
        URL:         "https://www.example.com",
        Name:        "Example WebPage",
        Headline:    "Welcome to Example WebPage",
        Description: "This is an example webpage.",
        About:       "Something related to the home page",
        Keywords:    "example, webpage, demo",
        InLanguage:  "en",
    }
 }}
 <!DOCTYPE html>
 <html lang="en">
   <head>
      <meta charset="UTF-8"/>
      <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
      <title>teseo - homepage</title>
      <!-- render JSON-LD here -->
      @webpage.ToJsonLd()
    </head>
    <body>
       <!-- your content -->
    </body>
 </html>
}
```

The expected output:

```html
<script type="application/ld+json">
  {
    "@context": "https://schema.org",
    "@type": "WebPage",
    "url": "https://www.example.com",
    "name": "Example WebPage",
    "headline": "Welcome to Example WebPage",
    "description": "This is an example webpage",
    "about": "Something related to the home page",
    "keywords": "example, webpage, demo",
    "inLanguage": "en"
  }
</script>
```

#### Example: BreadcrumbList using NewBreadcrumbListFromUrl

`teseo` also provides utility methods such as `NewBreadcrumbListFromUrl`, which helps you automatically generate a breadcrumb list based on the full page URL. This method is invaluable during development, as it helps quickly generate and structure breadcrumb navigation for dynamic or complex URLs, making debugging faster and more efficient.

```templ
package main

import (
    "github.com/indaco/teseo/schemaorg"
    "github.com/indaco/teseo"
    "net/http"
)

func HandleAbout(w http.ResponseWriter, r *http.Request) {
    pageURL := teseo.GetFullURL(r) // Helper function to get the full URL from the request
    breadcrumbList, err := schemaorg.NewBreadcrumbListFromUrl(pageURL)
    if err != nil {
        fmt.Println("Error generating breadcrumb list:", err)
        return
    }

    err = pages.AboutPage(breadcrumbList).Render(r.Context(), w)
    if err != nil {
        return
    }
}

templ AboutPage(breadcrumbList *schemaorg.BreadcrumbList) {
    <!DOCTYPE html>
    <html lang="en">
        <head>
            <meta charset="UTF-8"/>
            <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
            <title>teseo - about</title>
            <!-- render JSON-LD here -->
            @breadcrumbList.ToJsonLd()
        </head>
        <body>
            <!-- your content -->
        </body>
    </html>
}
```

The expected output for a URL like `https://www.example.com/about`:

```html
<script type="application/ld+json">
  {
    "@context": "https://schema.org",
    "@type": "BreadcrumbList",
    "itemListElement": [
      {
        "@type": "ListItem",
        "position": 1,
        "name": "Home",
        "item": "https://www.example.com"
      },
      {
        "@type": "ListItem",
        "position": 2,
        "name": "About",
        "item": "https://www.example.com/about"
      }
    ]
  }
</script>
```

#### SiteNavigationElementList: JSON-LD and Sitemap Generation

The **SiteNavigationElementList** represents a Schema.org `ItemList` composed of `SiteNavigationElement` entries. It can be used to structure navigation menus as JSON-LD and optionally generate a sitemap XML file.

**Example usage in Templ page:**

```go
package pages

import "github.com/indaco/teseo/schemaorg"

templ HomePage() {
 {{
    sne := schemaorg.NewSiteNavigationElementList(
      "main",
      []schemaorg.SiteNavigationElement{
        schemaorg.NewSimpleSiteNavigationElement(1, "Home", "https://www.example.com"),
        schemaorg.NewSimpleSiteNavigationElement(2, "About", "https://www.example.com/about"),
      },
    )
 }}
 <!DOCTYPE html>
 <html lang="en">
   <head>
      <meta charset="UTF-8"/>
      <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
      <title>teseo - homepage</title>
      <!-- Render JSON-LD -->
      @sne.ToJsonLd()
    </head>
    <body>
       <!-- your content -->
    </body>
 </html>
}
```

The expected output:

```html
<script type="application/ld+json">
  {
    "@context": "https://schema.org",
    "@type": "ItemList",
    "identifier": "main",
    "itemListElement": [
      {
        "@type": "SiteNavigationElement",
        "position": 1,
        "name": "Home",
        "url": "https://www.example.com"
      },
      {
        "@type": "SiteNavigationElement",
        "position": 2,
        "name": "About",
        "url": "https://www.example.com/about"
      }
    ]
  }
</script>
```

**Sitemap XML Generation:**

```go
package handlers

import (
  "log"
  "net/http"

  "github.com/indaco/teseo/schemaorg"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
  sne := schemaorg.NewSiteNavigationElementList(
    "main",
    []schemaorg.SiteNavigationElement{
      schemaorg.NewSimpleSiteNavigationElement(1, "Home", "https://www.example.com"),
      schemaorg.NewSimpleSiteNavigationElement(2, "About", "https://www.example.com/about"),
    },
  )

  err := sne.ToSitemapFile("./_demos/statics/sitemap.xml")
  if err != nil {
    log.Fatalf("Failed to generate sitemap: %v", err)
  }

  err = pages.HomePage(sne).Render(r.Context(), w)
  if err != nil {
    log.Printf("render error: %v", err)
  }
}
```

Similarly, the `FromSitemapFile` method allows you to parse a sitemap XML file and populate the `SiteNavigationElementList` struct. This is especially useful for debugging or importing existing sitemaps into your application logic.

### OpenGraph Meta Tags

For **OpenGraph**, entities come with `ToMetaTags` and `ToGoHTMLMetaTags` methods that generates the necessary meta tags for OpenGraph data. Similar to Schema.org, you can either create the entity via a **pure struct** or a **factory method**. Here’s an example for generating meta tags for an _Article_:

```templ
package pages

import "github.com/indaco/teseo/opengraph"

templ FirstArticle() {
 {{
    article := &opengraph.Article{
        Title:       "Example Article",
        URL:         "https://www.example.com/article/example-article",
        Description: "This is an example article description.",
        Image:       "https://www.example.com/images/article.jpg",
    }
 }}
 <!DOCTYPE html>
 <html lang="en">
   <head>
      <meta charset="UTF-8"/>
      <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
      <title>teseo - first article</title>
      <!-- render opengraph metadata here -->
      @article.ToMetaTags()
    </head>
    <body>
       <!-- your content -->
    </body>
 </html>
}
```

The expected output:

```html
<meta property="og:type" content="article" />
<meta property="og:title" content="Example Article" />
<meta
  property="og:url"
  content="https://www.example.com/article/example-article"
/>
<meta
  property="og:description"
  content="This is an example article description."
/>
<meta
  property="og:image"
  content="https://www.example.com/images/article.jpg"
/>
```

### Twitter Cards

For **Twitter Cards**, you can also use either the **pure struct** or **factory methods** to generate Twitter Card meta tags via the `ToMetaTags` and `ToGoHTMLMetaTags` methods. Here’s how to generate a _Twitter Summary Card_.

```templ
package pages

import "github.com/indaco/teseo/twittercard"

templ AboutMe() {
 {{
    twCard := &twittercard.TwitterCard{
        Card:        twittercard.CardSummary,
        Title:       "Example Summary",
        Description: "This is an example summary card.",
        Image:       "https://www.example.com/summary.jpg",
        Site:        "@example_site",
    }
 }}
 <!DOCTYPE html>
 <html lang="en">
   <head>
      <meta charset="UTF-8"/>
      <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
      <title>teseo - first article</title>
      <!-- render twittercard metadata here -->
      @twCard.ToMetaTags()
    </head>
    <body>
      <!-- your content -->
    </body>
 </html>
}
```

The expected output:

```html
<meta name="twitter:card" content="summary" />
<meta name="twitter:title" content="Example Summary" />
<meta name="twitter:description" content="This is an example summary card." />
<meta name="twitter:image" content="https://www.example.com/summary.jpg" />
<meta name="twitter:site" content="@example_site" />
```

This works for all supported Twitter Cards (e.g., App Card, Player Card, etc.).

## Demo

Check out the [\_demos](_demos/) folder for real-world usage of:

- JSON-LD structured data
- OpenGraph meta tags
- Twitter Card metadata

### Run the demo

```bash
just dev # http://localhost:3300
```

## Contributing

Contributions are welcome!

See the [Contributing Guide](/CONTRIBUTING.md) for setup instructions.

## License

This project is licensed under the MIT License – see the [LICENSE](./LICENSE) file for details.
