<h1 align="center" style="font-size: 2.5rem;">teseo</h1>
<h2 align="center">Go templ components for SEO</h2>
<p align="center">
    <a href="https://github.com/indaco/teseo/blob/main/LICENSE" target="_blank">
        <img src="https://img.shields.io/badge/license-mit-blue?style=flat-square&logo=none" alt="license" />
    </a>
     &nbsp;
     <a href="https://goreportcard.com/report/github.com/indaco/teseo" target="_blank">
        <img src="https://goreportcard.com/badge/indaco/teseo" alt="go report card" />
    </a>
    &nbsp;
    <a href="https://pkg.go.dev/github.com/indaco/teseo/" target="_blank">
        <img src="https://pkg.go.dev/badge/github.com/indaco/teseo/.svg" alt="go reference" />
    </a>
    &nbsp;
    <a href="https://www.jetify.com/devbox/docs/contributor-quickstart/">
      <img
          src="https://www.jetify.com/img/devbox/shield_moon.svg"
          alt="Built with Devbox"
      />
  </a>
</p>

`teseo` provides a comprehensive list of SEO-related data types (go structs) that adhere to **Schema.org** and **OpenGraph** specifications, with methods to easily generate [templ](https://github.com/a-h/templ) components or standard `template/html` output from them.

Whether you are looking to implement _Schema.org JSON-LD_, _OpenGraph_, or _Twitter Cards_, `teseo` helps you generate SEO-friendly meta information effortlessly.

## Features

- A comprehensive list of useful **Schema.org JSON-LD** types.
- A comprehensive list of useful **OpenGraph** meta tags.
- Support for **Twitter Cards** meta tags.
- Easy-to-use methods to generate JSON-LD and meta tags.
- Render data types as **templ components** or using **template/html**.

## Supported Data Types

### Schema.org JSON-LD Entities

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

### OpenGraph Data Types

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

## Installation

Add this package to your project:

```bash
go get github.com/indaco/teseo@latest
```

## Usage

### Schema.org JSON-LD

For **Schema.org JSON-LD**, each entity provides `ToJsonLd` and `ToGoHTMLJsonLd` methods. You can render the structured data as a templ component or as an HTML string, suitable for Go's `template/html`. Entities can be created using **pure structs** or **factory methods**.

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

#### SiteNavigationElement: JSON-LD and Sitemap Generation

The **SiteNavigationElement** represents a Schema.org object that can be used to structure site navigation data. This entity supports both JSON-LD generation and the creation of a sitemap XML file.

**Factory method usage:**

```go
package pages

import "github.com/indaco/teseo/schemaorg"

templ HomePage() {
 {{
    sne := schemaorg.NewSiteNavigationElementWithItemList(
      "Main Navigation",
      "https://www.example.com",
      []schemaorg.ItemListElement{
        {Name: "Home", URL: "https://www.example.com", Position: 1},
        {Name: "About", URL: "https://www.example.com/about", Position: 2},
      },
    )
 }}
 <!DOCTYPE html>
 <html lang="en">
   <head>
      <meta charset="UTF-8"/>
      <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
      <title>teseo - homepage</title>
      <!-- render JSON-LD here -->
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
  "@type": "SiteNavigationElement",
  "name": "Main Navigation",
  "url": "https://www.example.com",
  "position": 1,
  "itemListElement": [
    {"@type": "ListItem", "position": 1, "name": "Home", "url": "https://www.example.com"},
    {"@type": "ListItem", "position": 2, "name": "About", "url": "https://www.example.com/about"}
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
  sne := schemaorg.NewSiteNavigationElementWithItemList(
    "Main Navigation",
    "https://www.example.com",
    []schemaorg.ItemListElement{
      {Name: "Home", URL: "https://www.example.com", Position: 1},
      {Name: "About", URL: "https://www.example.com/about", Position: 2},
    },
  )


  err := sne.ToSitemapFile("./_demos/statics/sitemap.xml")
  if err != nil {
    log.Fatalf("Failed to generate sitemap: %v", err)
  }

  err = pages.HomePage(sne).Render(r.Context(), w)
  if err != nil {
    return
  }
}
```

Similarly, the `FromSitemapFile` method allows you to parse a sitemap XML file and populate the `SiteNavigationElement` struct. This can speed up the debugging process and is particularly useful when working with dynamically generated sitemaps.

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
<meta property="og:type" content="article"/>
<meta property="og:title" content="Example Article"/>
<meta property="og:url" content="https://www.example.com/article/example-article"/>
<meta property="og:description" content="This is an example article description."/>
<meta property="og:image" content="https://www.example.com/images/article.jpg"/>
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
<meta name="twitter:card" content="summary"/>
<meta name="twitter:title" content="Example Summary"/>
<meta name="twitter:description" content="This is an example summary card."/>
<meta name="twitter:image" content="https://www.example.com/summary.jpg"/>
<meta name="twitter:site" content="@example_site"/>
```

This works for all supported Twitter Cards (e.g., App Card, Player Card, etc.).

## Demo

Check out the _demos folder for real-world usage of:

- JSON-LD structured data
- OpenGraph meta tags
- Twitter Card metadata

### Run the demo

```bash
# Taskfile
task dev # http://localhost:7332

# Makefile
make dev # http://localhost:7332
```

## Contributing

Contributions are welcome!

See the [Contributing Guide](/CONTRIBUTING.md) for setup instructions.

## License

This project is licensed under the MIT License – see the [LICENSE](./LICENSE) file for details.
