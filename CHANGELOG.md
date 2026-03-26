# Changelog

All notable changes to this project will be documented in this file.

This project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html). The changelog is generated and managed by [sley](https://github.com/indaco/sley).

## v0.2.4 - 2026-03-26

### 🩹 Fixes

- **sitenavigationelementlist:** set Name and Description in FromSitemapFile ([11c10af](https://github.com/indaco/teseo/commit/11c10af))

### 📖 Documentation

- **README:** add more badges ([170e63c](https://github.com/indaco/teseo/commit/170e63c)) ([#19](https://github.com/indaco/teseo/pull/19))

### 🏡 Chores

- add prek.toml configuration ([a751bc7](https://github.com/indaco/teseo/commit/a751bc7))
- migrate task runner from Taskfile/Make to just ([f2f607d](https://github.com/indaco/teseo/commit/f2f607d)) ([#16](https://github.com/indaco/teseo/pull/16))
- setup github.com/indaco/sley ([483a43d](https://github.com/indaco/teseo/commit/483a43d)) ([#15](https://github.com/indaco/teseo/pull/15))
- regenerate \*\_templ.go files with templ v0.3.924 ([1dc7ca5](https://github.com/indaco/teseo/commit/1dc7ca5))
- **demos:** fix type for AreaServed in seo.templ ([36a50ef](https://github.com/indaco/teseo/commit/36a50ef))

### 🤖 CI

- harden GitHub Actions workflows with zizmor recommendations ([736fdc2](https://github.com/indaco/teseo/commit/736fdc2)) ([#17](https://github.com/indaco/teseo/pull/17))

### 📦 Build

- upgrade templ to v0.3.1001 ([8da3e96](https://github.com/indaco/teseo/commit/8da3e96)) ([#18](https://github.com/indaco/teseo/pull/18))
- update templ to v0.3.924 ([e7944e6](https://github.com/indaco/teseo/commit/e7944e6))

**Full Changelog:** [v0.2.3...v0.2.4](https://github.com/indaco/teseo/compare/v0.2.3...v0.2.4)

### ❤️ Contributors

- [@indaco](https://github.com/indaco)

## v0.2.3 - 2025-07-21

[compare changes](https://github.com/indaco/teseo/compare/v0.2.2...v0.2.3)

### 🩹 Fixes

- Merge pull request [#14](https://github.com/indaco/teseo/pull/14) from [@chrisbward](https://github.com/chrisbward)

### 💅 Refactors

- **schemaorg/types:** Rename to StringList + add full test coverage ([50e65ae](https://github.com/indaco/teseo/commit/50e65ae))

### 📦 Build

- Update templ to v0.3.920 ([d68bd77](https://github.com/indaco/teseo/commit/d68bd77))

### ❤️ Contributors

- Indaco ([@indaco](https://github.com/indaco))
- chrisbward ([@chrisbward](https://github.com/chrisbward)

## v0.2.2 - 2025-06-25

[compare changes](https://github.com/indaco/teseo/compare/v0.2.1...v0.2.2)

### 📦 Build

- Update templ to v0.3.906 ([bc3ba51](https://github.com/indaco/teseo/commit/bc3ba51))

### ❤️ Contributors

- Indaco ([@indaco](https://github.com/indaco))

## v0.2.1 - 2025-06-25

[compare changes](https://github.com/indaco/teseo/compare/v0.2.0...v0.2.1)

### 🩹 Fixes

- Correct `SiteNavigationElement` JSON-LD type usage ([2cc6815](https://github.com/indaco/teseo/commit/2cc6815))

### 📖 Documentation

- **README:** Update `SiteNavigationElement` section to reflect new ItemList-based API ([825fa45](https://github.com/indaco/teseo/commit/825fa45))

### 📦 Build

- Update templ to v0.3.898 ([bcd575d](https://github.com/indaco/teseo/commit/bcd575d))

### 🏡 Chore

- **go.mod:** Set go 1.23.0 and add toolchain directive ([6138602](https://github.com/indaco/teseo/commit/6138602))
- Regenerate \*\_templ.go files with templ v0.3.898 ([dd81e83](https://github.com/indaco/teseo/commit/dd81e83))
- Update devbox.lock ([32359f6](https://github.com/indaco/teseo/commit/32359f6))

### 🎨 Styles

- Resolve linter warnings ([8435be3](https://github.com/indaco/teseo/commit/8435be3))

### 🤖 CI

- Run ci workflow on all branches ([2e927fe](https://github.com/indaco/teseo/commit/2e927fe))
- Update golangci-lint-action to v8 ([01a0cc3](https://github.com/indaco/teseo/commit/01a0cc3))

### ❤️ Contributors

- Indaco ([@indaco](https://github.com/indaco))

## v0.2.0 - 2025-04-07

[compare changes](https://github.com/indaco/teseo/compare/v0.1.0...v0.2.0)

### 🚀 Enhancements

- **schemaorg:** Add validation methods for entities ([fc43c71](https://github.com/indaco/teseo/commit/fc43c71))

### 🩹 Fixes

- Test command in pre-push and pre-commit hooks ([be39f98](https://github.com/indaco/teseo/commit/be39f98))
- **twittercard:** Ensure default values in NewCard function ([a78d486](https://github.com/indaco/teseo/commit/a78d486))
- **schemaorg/website:** Prop in ensureDefaults ([f9d1744](https://github.com/indaco/teseo/commit/f9d1744))

### 💅 Refactors

- **opengraph:** Unify metadata rendering and improve optional field handling ([a6f3bd2](https://github.com/indaco/teseo/commit/a6f3bd2))
- **schemaorg:** Simplify `ToGoHTMLJsonLd` method ([c6232d9](https://github.com/indaco/teseo/commit/c6232d9))
- Improve unique key generation and URL handling ([4996ac2](https://github.com/indaco/teseo/commit/4996ac2))
- **opengraph:** Replace `templ.ToGoHTML` with `teseo.RenderToHTML` ([d8313c5](https://github.com/indaco/teseo/commit/d8313c5))
- HTML rendering in _TwitterCard_ with `RenderToHTML` ([587ea2a](https://github.com/indaco/teseo/commit/587ea2a))
- **schemaorg:** Move `Organization` and `Person` structs to relative files ([7dbc023](https://github.com/indaco/teseo/commit/7dbc023))

### 📖 Documentation

- **README:** Add contributing guide and update sections ([3200ba5](https://github.com/indaco/teseo/commit/3200ba5))
- **README:** Update badges and reorganize the content ([d6c63b3](https://github.com/indaco/teseo/commit/d6c63b3))
- **README:** Remove slash from meta output ([0ee3696](https://github.com/indaco/teseo/commit/0ee3696))

### 📦 Build

- Bump Go version and update dependencies ([6483d38](https://github.com/indaco/teseo/commit/6483d38))

### 🏡 Chore

- **demos:** Regenerate with templ v0.3.857 ([7174d92](https://github.com/indaco/teseo/commit/7174d92))

### ✅ Tests

- **opengraph:** Add tests ([86f3d24](https://github.com/indaco/teseo/commit/86f3d24))
- Replace manual slice search with slices.Contains ([d1de916](https://github.com/indaco/teseo/commit/d1de916))

### 🤖 CI

- Add Git hooks for commit message validation and testing ([7dfb6d8](https://github.com/indaco/teseo/commit/7dfb6d8))
- Add lint, test and coverage workflows ([2a09be7](https://github.com/indaco/teseo/commit/2a09be7))
- Fix indentation in coverage.yml ([9be8d48](https://github.com/indaco/teseo/commit/9be8d48))
- Simplify Makefile and Taskfile for clarity ([acb267d](https://github.com/indaco/teseo/commit/acb267d))
- Allow pull requests from all branches in CI workflow ([268bdde](https://github.com/indaco/teseo/commit/268bdde))
- Add GitHub Actions workflow for release notes generation ([be5e3f9](https://github.com/indaco/teseo/commit/be5e3f9))

### ❤️ Contributors

- Indaco ([@indaco](https://github.com/indaco))

## v0.1.0 - 2024-10-07

### 🏡 Chore

- initial release

### ❤️ Contributors

- Indaco <github@mircoveltri.me>
