<h1 align="center" style="font-size: 2.5rem;">teseo</h1>
<h2 align="center">Go templ components for SEO</h2>
<p align="center">
    <a href="https://github.com/indaco/teseo/blob/main/LICENSE" target="_blank">
        <img src="https://img.shields.io/badge/license-mit-blue?style=flat-square&logo=none" alt="license" />
    </a>
     &nbsp;
     <a href="https://goreportcard.com/report/github.com/indaco/teseo/" target="_blank">
        <img src="https://goreportcard.com/badge/github.com/indaco/teseo" alt="go report card" />
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
- Easy-to-use functions to generate JSON-LD and meta tags.
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

```sh
go get github.com/indaco/teseo@latest
```

## Usage

### Schema.org JSON-LD

For **Schema.org JSON-LD**, each entity provides `ToTemplJsonLd` and `ToGoHTMLJsonLd` functions. You can render the structured data as a templ component or as an HTML string, suitable for Go's `template/html`. Entities can be created using **pure structs** or **factory methods**.

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

`teseo` also provides utility functions such as `NewBreadcrumbListFromUrl`, which helps you automatically generate a breadcrumb list based on the full page URL.

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

Then render it in your templ component:

```templ
package pages

import "github.com/indaco/teseo/schemaorg"

templ HomePage(sne *schemaorg.SiteNavigationElement) {
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

**Parsing Sitemap XML:**

You can also parse an existing sitemap XML file and populate the `SiteNavigationElement` struct:

```go
package handlers

import (
  "log"
  "net/http"

  "github.com/indaco/teseo"
  "github.com/indaco/teseo/_demos/pages"
  "github.com/indaco/teseo/schemaorg"
)

func HandleAbout(w http.ResponseWriter, r *http.Request) {
  sne := &schemaorg.SiteNavigationElement{}
  err = sne.FromSitemapFile("./_demos/statics/sitemap.xml")
  if err != nil {
    log.Fatalf("Failed to read sitemap: %v", err)
  }

  err = pages.AboutPage(sne).Render(r.Context(), w)
  if err != nil {
    return
  }
}
```

Then render it in your templ component as the example above.

### OpenGraph Meta Tags

For **OpenGraph**, entities come with a `ToMetaTags` function that generates the necessary meta tags for OpenGraph data. Similar to Schema.org, you can either create the entity via a **pure struct** or a **factory method**. Here’s an example for generating meta tags for an _Article_:

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

For **Twitter Cards**, you can also use either the **pure struct** or **factory methods** to generate Twitter Card meta tags via the `ToMetaTags` function.. Here’s how to generate a _Twitter Summary Card_.

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

A sample website is available in the **_demos** folder, which demonstrates how to integrate teseo for generating structured data and metadata. This demo serves as a reference for implementing Schema.org JSON-LD, OpenGraph, and Twitter Cards in your own web applications.

Feel free to explore the demo to see real-world usage of the library and how easily you can add SEO-friendly metadata to your Go web projects.

### Run the demo

```bash
# Taskfile
task live # http://localhost:7332
# Makefile
make live # http://localhost:7332
```

## Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

### Development Environment Setup

To set up a development environment for this repository, you can use [devbox](https://www.jetify.com/devbox) along with the provided `devbox.json` configuration file.

1. Install devbox by following the instructions in the [devbox documentation](https://www.jetify.com/devbox/docs/installing_devbox/).
2. Clone this repository to your local machine.
3. Navigate to the root directory of the cloned repository.
4. Run `devbox install` to install all packages mentioned in the `devbox.json` file.
5. Run `devbox shell --pure` to start a new shell with access to the environment.
6. Once the devbox environment is set up, you can start developing, testing, and contributing to the repository.

### Running Tasks

This project provides both a `Makefile` and a `Taskfile` for running various tasks. You can use either `make` or `task` to execute the tasks, depending on your preference.

To view all available tasks, run:

- **Makefile**: `make help`
- **Taskfile**: `task --list-all`

Available tasks:

```bash
live      # Run the demos live server with templ watch mode.
templ     # Run templ fmt and templ generate commands.
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
