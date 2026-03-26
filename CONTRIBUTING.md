# Contributing

Contributions are welcome! Feel free to open an issue or submit a pull request.

## Development Environment Setup

### Using Devbox (Recommended)

To set up a development environment for this repository, you can use [devbox](https://www.jetify.com/devbox) along with the provided `devbox.json` configuration file.

1. Install devbox by following [these instructions](https://www.jetify.com/devbox/docs/installing_devbox/).
2. Clone this repository to your local machine.

   ```bash
   git clone https://github.com/indaco/teseo.git
   cd teseo
   ```

3. Run `devbox install` to install all dependencies specified in `devbox.json`.
4. Enter the environment with `devbox shell --pure`.
5. Start developing, testing, and contributing!

### Manual Setup

If you prefer not to use Devbox, ensure you have the following tools installed:

- [Go](https://go.dev/)
- [just](https://github.com/casey/just): For running project tasks.
- [golangci-lint](https://golangci-lint.run/): For linting Go code.
- [modernize](https://pkg.go.dev/golang.org/x/tools/gopls/internal/analysis/modernize): Run the modernizer analyzer to simplify code by using modern constructs.
- [prek](https://github.com/indaco/prek): For managing git hooks.

## Setting Up Git Hooks

Git hooks are used to enforce code quality and streamline the workflow.

### Using Devbox

If using `devbox`, Git hooks are automatically installed via [prek](https://github.com/indaco/prek) when you run `devbox shell`. No further action is required.

### Manual Setup

For users not using `devbox`, install [prek](https://github.com/indaco/prek) and run:

```bash
prek install --hook-type commit-msg --hook-type pre-push
```

## Running Tasks

This project uses [just](https://github.com/casey/just) as a task runner.

### View all available tasks

```bash
just
```

### Common tasks

```bash
just clean            # Clean the build directory and Go cache
just fmt              # Format code
just modernize        # Run go-modernize with auto-fix
just lint             # Run golangci-lint
just check            # Run fmt, modernize, lint, and reportcard
just test             # Run all tests and print code coverage value
just test-force       # Clean go tests cache and run all tests
just test-coverage    # Run all tests and generate coverage report
just test-race        # Run all tests with race detector
just templ            # Run templ fmt and templ generate on the demos
just dev              # Run the demos server
```
