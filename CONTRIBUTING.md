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

- [golangci-lint](https://golangci-lint.run/): For linting Go code.
- [go-task](https://taskfile.dev/) or `make`: For running project tasks.
- [modernize](https://pkg.go.dev/golang.org/x/tools/gopls/internal/analysis/modernize): Run the modernizer analyzer to simplify code by using modern constructs.

## Setting Up Git Hooks

Git hooks are used to enforce code quality and streamline the workflow. Follow these steps to set them up:

### Using Devbox

If using `devbox`, Git hooks are automatically installed when you run `devbox shell`. No further action is required.

### Manual Setup

For users not using `devbox`, follow the steps below to manually install the Git hooks:

1. Clone the repository:

   ```bash
   git clone https://github.com/indaco/teseo.git
   cd teseo
   ```

2. Install the Git Hooks

   **Unix-based systems (Linux, macOS):**

   ```bash
   sh .scripts/setup-hooks.sh
   ```

   **Windows systems:**

   ```cmd
   .scripts\setup-hooks.bat
   ```

## Running Tasks

This project provides both a `Makefile` and a `Taskfile` for running various tasks. You can use either `make` or `task` to execute the tasks, depending on your preference.

### View all available tasks

- _Makefile_: `make help`
- _Taskfile_: `task --list-all`

#### Common tasks

```bash
build:               # Build for production
clean:               # Clean the build directory and Go cache
dev:                 # Run the demos server
lint:                # Run golangci-lint
modernize:           # Run go-modernize
templ:               # Run templ fmt and templ generate commands on the demos
test:                # Run all tests and generate coverage report
test/coverage:       # Run go tests and use go tool cover
test/force:          # Clean go tests cache and run all tests
```
