# obsave

`obsave` is a command-line utility that allows you to pipe text content into an Obsidian vault, adding YAML front matter in the process. It offers features like front matter customization, debug mode, and options to control how existing files are handled.

## Changelog

### Latest Changes
- Enhanced debug logging with detailed configuration and operation reporting
- Added extended help system:
  - `-h` shows concise usage information
  - `--help` displays detailed help with examples
- Added short form options:
  - `-n` for `--name`
  - `-t` for `--tags`
  - `-p` for `--properties`
  - `-v` for `--verbose`
  - `-ob` for `--vault`
  - `-c` for `--config`
- Improved configuration handling:
  - Changed config from positional argument to flag option (`-c` or `--config`)
  - Implemented clear configuration precedence
  - Better handling of default and specified config files

## Features

- **Front Matter Generation**: Automatically adds YAML front matter to your notes, including fields such as `title`, `tags`, and custom `properties`.
- **Custom Properties**: Add custom key-value pairs to the front matter.
- **Multiple Configuration Files**: Support for multiple YAML configuration files, allowing for different setups for various projects or use cases.
- **Flexible Tag and Property Handling**: Options to replace, add, or merge tags and properties.
- **Overwrite Modes**: Control how to handle existing files.
- **Dry Run Mode**: Simulate operations without writing files.
- **Debug Mode**: Enable detailed logging.
- **Verbose Mode**: Print the full path of saved files.

> [!IMPORTANT]
> ## Changelog
> - **19-Oct-2024**: Added `--verbose` option to output final filename and path
> - **17-Oct-2024**: 
>   - Switched to YAML for configuration
>   - Implemented separate configuration files
> - **16-Oct-2024**:
>   - Implemented `--debug` option for verbose output
>   - Added `--dry-run` option and support for custom properties
> - **15-Oct-2024**: Initial commit of Obsave project

## Installation

To install `obsave`, you need to have Go installed on your machine. You can install the utility with:

```bash
go install github.com/mattjoyce/obsave@latest
```

This will install the utility to your `$GOPATH/bin` directory.

## Command Line Options

### Core Options
- `-n, --name`: Note name/title (required if not in config)
- `-ob, --vault`: Path to Obsidian vault (required if not in config)
- `-t, --tags`: Comma-separated list of tags
- `-p, --properties`: Custom frontmatter properties (format: key=value;key2=value2)

### Configuration
- `-c, --config`: Specify a config file to use (default: ~/.config/obsave/config)

### Control Options
- `--overwrite-mode`: How to handle existing files ("overwrite", "serialize", or "fail")
- `--tags-handling`: How to handle tags ("replace", "add", or "merge")
- `--properties-handling`: How to handle properties ("replace", "add", or "merge")

### Output Control
- `-v, --verbose`: Print the full path of the saved file
- `--debug`: Enable detailed debug logging
- `--dry-run`: Simulate the operation without writing files

### Help and Documentation
- `-h`: Display concise usage information
- `--help`: Display detailed help with examples

## Usage

### Basic Example

You can pipe text content into `obsave` and save it to your Obsidian vault:

```bash
echo "This is my note content." | obsave -n "MyNote" -t "project,example" -p "author=John Doe;status=Draft" -ob "~/Documents/ObsidianVault"
```

### Additional Options

- **Debug Mode**:
  Enable detailed logging to troubleshoot or inspect the internal operations of `obsave`:
  ```bash
  echo "Test content" | obsave -n "DebugTest" -ob "~/vault" --debug
  ```

- **Dry Run**:
  Use the `--dry-run` option to simulate the operation without writing the file:
  ```bash
  echo "Test content" | obsave -n "TestNote" -ob "~/vault" --dry-run
  ```

- **Overwrite Mode**:
  Specify how to handle existing files:
  ```bash
  echo "New content" | obsave -n "ExistingNote" -ob "~/vault" --overwrite-mode "overwrite"
  ```

## Configuration

Configuration files use YAML format and are stored in `~/.config/obsave/`. The configuration system follows a clear precedence order:

1. Start with empty configuration
2. Load default config file (~/.config/obsave/config) if it exists
3. Load specified config file (if -c/--config provided)
4. Apply command line options (these always override config file settings)

> [!IMPORTANT]
> Windows users should use `%USERPROFILE%` instead of `~/` for the config directory.  Which might be `C:\Users\<username>\`

### Required Settings
The following settings must be provided either through a config file or command line options:
- `name` (via -n/--name)
- `vault_path` (via -ob/--vault)

### Example Configuration File
```yaml
vault_path: "~/Documents/ObsidianVault"
overwrite_mode: "fail"
tags:
  - default
  - obsave
properties:
  tool: obsave
  version: "1.0"
debug: false
dry_run: false
tags_handling: "merge"
properties_handling: "merge"
name: "Default Note Name"  # optional
verbose: false
```

### Configuration Scenarios

1. **No Config File**:
   - Must provide --name and --vault options
   - Other options use program defaults
   ```bash
   echo "Content" | obsave -n "Note" -ob "~/vault"
   ```

2. **Default Config Exists**:
   - Settings from ~/.config/obsave/config are used
   - Can override any setting via command line
   ```bash
   # Use config but override name
   echo "Content" | obsave -n "Custom Name"
   ```

3. **Specified Config**:
   - Loads specified config instead of default
   - Can still override via command line
   ```bash
   echo "Content" | obsave -c custom-config -n "Override Name"
   ```

## Tag and Property Handling Options

obsave provides flexible options for managing tags and properties. These options control how the command-line inputs interact with existing configuration or content.

### Tags Handling

Use the `--tags-handling` option to specify how tags should be managed. Available modes are:

- **Replace**: Completely replaces existing tags with the ones provided in the CLI.
  ```bash
  echo "Content" | obsave -n "TaggedNote" -t "example,notes,project" --tags-handling "replace"
  ```
  * If config had: `["existing", "tags"]`
  * CLI tags: `"example,notes,project"`
  * Result: `["example", "notes", "project"]`

- **Add**: Adds new tags from the CLI, avoiding duplicates.
  ```bash
  echo "Content" | obsave -n "TaggedNote" -t "example,notes,project" --tags-handling "add"
  ```
  * If config had: `["project", "important"]`
  * CLI tags: `"example,notes,project"`
  * Result: `["project", "important", "example", "notes"]`

- **Merge**: Adds all tags from the CLI, allowing duplicates.
  ```bash
  echo "Content" | obsave -n "TaggedNote" -t "example,notes,project" --tags-handling "merge"
  ```
  * If config had: `["project", "important"]`
  * CLI tags: `"example,notes,project"`
  * Result: `["project", "important", "example", "notes", "project"]`

### Properties Handling

Use the `--properties-handling` option to specify how properties should be managed. Available modes are:

- **Replace**: Completely replaces existing properties with the ones provided in the CLI.
  ```bash
  echo "Content" | obsave -n "PropertyNote" -p "status=Review;author=John" --properties-handling "replace"
  ```
  * If config had: `{"category": "Work", "priority": "High"}`
  * CLI properties: `"status=Review;author=John"`
  * Result: `{"status": "Review", "author": "John"}`

- **Add**: Adds new properties from the CLI, without overwriting existing ones.
  ```bash
  echo "Content" | obsave -n "PropertyNote" -p "status=Review;author=John" --properties-handling "add"
  ```
  * If config had: `{"category": "Work", "priority": "High"}`
  * CLI properties: `"status=Review;author=John"`
  * Result: `{"category": "Work", "priority": "High", "status": "Review", "author": "John"}`

- **Merge**: Adds all properties from the CLI, overwriting existing ones with the same key.
  ```bash
  echo "Content" | obsave -n "PropertyNote" -p "status=Review;priority=Medium" --properties-handling "merge"
  ```
  * If config had: `{"category": "Work", "priority": "High"}`
  * CLI properties: `"status=Review;priority=Medium"`
  * Result: `{"category": "Work", "priority": "Medium", "status": "Review"}`

## Cloning and Building from Source

### Prerequisites

Make sure you have Go installed on your system. You can install Go from [here](https://golang.org/dl/).

### Steps to Clone and Build:

1. **Clone the Repository**:
   ```bash
   git clone https://github.com/mattjoyce/obsave.git
   cd obsave
   ```

2. **Build the Project**:
   To build the project into an executable binary, run:
   ```bash
   go build -o obsave
   ```

   This will generate the `obsave` executable in the current directory.

3. **Run the Utility**:
   You can now run `obsave` directly:
   ```bash
   ./obsave -n "TestNote" -t "example" -ob "~/Documents/ObsidianVault"
   ```

4. **Install Locally** (Optional):
   If you want to install `obsave` to your system's `$GOPATH/bin` directory for global usage, run:
   ```bash
   go install
   ```

   After installation, you can use `obsave` from anywhere in your terminal:
   ```bash
   obsave -n "NewNote" -ob "~/ObsidianVault"
   ```
