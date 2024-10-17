# obsave

`obsave` is a command-line utility that allows you to pipe text content into an Obsidian vault, adding YAML front matter in the process. It offers features like front matter customization, debug mode, and options to control how existing files are handled.

## Features

- **Front Matter Generation**: Automatically adds YAML front matter to your notes, including fields such as `title`, `tags`, and custom `properties`.
- **Custom Properties**: Use the `--properties` flag to add custom key-value pairs to the front matter (e.g., `author=John Doe;status=Draft`).
- **Multiple Configuration Files**: Support for multiple YAML configuration files, allowing for different setups for various projects or use cases.
- **Flexible Tag and Property Handling**: Options to replace, add, or merge tags and properties from the command line with those in the configuration.
- **Overwrite Modes**: Control how to handle existing files using the `--overwrite-mode` option. Available modes:
  - `overwrite`: Overwrites the existing file.
  - `serialize`: Creates a new version of the file with an incremented suffix (e.g., `note_1.md`).
  - `fail`: Aborts the operation if the file already exists.
- **Dry Run Mode**: Simulate the operation without writing the file by using the `--dry-run` flag.
- **Debug Mode**: Enable detailed logging with the `--debug` flag.

## Installation

To install `obsave`, you need to have Go installed on your machine. You can install the utility with:

```bash
go install github.com/mattjoyce/obsave@latest
```

This will install the utility to your `$GOPATH/bin` directory.

## Usage

### Basic Example

You can pipe text content into `obsave` and save it to your Obsidian vault:

```bash
echo "This is my note content." | obsave --name "MyNote" --tags "project,example" --properties "author=John Doe;status=Draft" --vault "~/Documents/ObsidianVault"
```

In this example:
- `--name "MyNote"` specifies the name of the note.
- `--tags "project,example"` adds tags to the note.
- `--properties "author=John Doe;status=Draft"` adds custom key-value pairs to the front matter.
- `--vault "~/Documents/ObsidianVault"` specifies the path to your Obsidian vault.

### Additional Options

- **Debug Mode**:
  Enable detailed logging to troubleshoot or inspect the internal operations of `obsave`:
  ```bash
  echo "Test content" | obsave --name "DebugTest" --vault "~/vault" --debug
  ```

- **Dry Run**:
  Use the `--dry-run` option to simulate the operation without writing the file:
  ```bash
  echo "Test content" | obsave --name "TestNote" --vault "~/vault" --dry-run
  ```

- **Overwrite Mode**:
  Specify how to handle existing files:
  ```bash
  echo "New content" | obsave --name "ExistingNote" --vault "~/vault" --overwrite-mode "overwrite"
  ```

- **Tags Handling**:
  Control how tags are managed:
  ```bash
  echo "Content" | obsave --name "TaggedNote" --tags "new,tags" --tags-handling "merge"
  ```

- **Properties Handling**:
  Control how properties are managed:
  ```bash
  echo "Content" | obsave --name "PropertyNote" --properties "status=Review" --properties-handling "add"
  ```

### Tag and Property Handling Options

obsave provides flexible options for managing tags and properties. These options control how the command-line inputs interact with existing configuration or content.

#### Tags Handling

Use the `--tags-handling` option to specify how tags should be managed. Available modes are:

- **Replace**: Completely replaces existing tags with the ones provided in the CLI.
  ```bash
  echo "Content" | obsave --name "TaggedNote" --tags "example,notes,project" --tags-handling "replace"
  ```
  * If config had: `["existing", "tags"]`
  * CLI tags: `"example,notes,project"`
  * Result: `["example", "notes", "project"]`

- **Add**: Adds new tags from the CLI, avoiding duplicates.
  ```bash
  echo "Content" | obsave --name "TaggedNote" --tags "example,notes,project" --tags-handling "add"
  ```
  * If config had: `["project", "important"]`
  * CLI tags: `"example,notes,project"`
  * Result: `["project", "important", "example", "notes"]`

- **Merge**: Adds all tags from the CLI, allowing duplicates.
  ```bash
  echo "Content" | obsave --name "TaggedNote" --tags "example,notes,project" --tags-handling "merge"
  ```
  * If config had: `["project", "important"]`
  * CLI tags: `"example,notes,project"`
  * Result: `["project", "important", "example", "notes", "project"]`

#### Properties Handling

Use the `--properties-handling` option to specify how properties should be managed. Available modes are:

- **Replace**: Completely replaces existing properties with the ones provided in the CLI.
  ```bash
  echo "Content" | obsave --name "PropertyNote" --properties "status=Review;author=John" --properties-handling "replace"
  ```
  * If config had: `{"category": "Work", "priority": "High"}`
  * CLI properties: `"status=Review;author=John"`
  * Result: `{"status": "Review", "author": "John"}`

- **Add**: Adds new properties from the CLI, without overwriting existing ones.
  ```bash
  echo "Content" | obsave --name "PropertyNote" --properties "status=Review;author=John" --properties-handling "add"
  ```
  * If config had: `{"category": "Work", "priority": "High"}`
  * CLI properties: `"status=Review;author=John"`
  * Result: `{"category": "Work", "priority": "High", "status": "Review", "author": "John"}`

- **Merge**: Adds all properties from the CLI, overwriting existing ones with the same key.
  ```bash
  echo "Content" | obsave --name "PropertyNote" --properties "status=Review;priority=Medium" --properties-handling "merge"
  ```
  * If config had: `{"category": "Work", "priority": "High"}`
  * CLI properties: `"status=Review;priority=Medium"`
  * Result: `{"category": "Work", "priority": "Medium", "status": "Review"}`

## Configuration

The utility now uses YAML configuration files stored in `~/.config/obsave/`. The default configuration file is named `config`, but you can specify different configuration files as needed.

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
```

### Using Multiple Configuration Files

You can use different configuration files by specifying them as the first argument:

```bash
echo "Custom config content" | obsave custom_config --name "CustomNote"
```

This will use the configuration from `~/.config/obsave/custom_config`.

[!note]
All option can be specified in the config yaml, except `name` which must be specified as a cli option.

## Cloning and Building from Source

If you'd like to clone the repository and build `obsave` locally, follow these steps:

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
   ./obsave --name "TestNote" --tags "example" --vault "~/Documents/ObsidianVault"
   ```

4. **Install Locally** (Optional):
   If you want to install `obsave` to your system's `$GOPATH/bin` directory for global usage, run:
   ```bash
   go install
   ```

   After installation, you can use `obsave` from anywhere in your terminal:
   ```bash
   obsave --name "NewNote" --vault "~/ObsidianVault"
   ```
