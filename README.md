
# obsave

`obsave` is a command-line utility that allows you to pipe text content into an Obsidian vault, adding YAML front matter in the process. It offers features like front matter customization, debug mode, and options to control how existing files are handled.

## Features

- **Front Matter Generation**: Automatically adds YAML front matter to your notes, including fields such as `title`, `tags`, and custom properties.
- **Custom Properties**: Use the `--properties` flag to add custom key-value pairs to the front matter (e.g., `author=John Doe;status=Draft`).
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
  - `overwrite`: Overwrites the file if it already exists.
  - `serialize`: Saves a new version of the file if it exists.
  - `fail`: Fails if the file already exists.
  
  Example:
  ```bash
  echo "New content" | obsave --name "ExistingNote" --vault "~/vault" --overwrite-mode "overwrite"
  ```

## Configuration

The utility uses a configuration file stored in `~/.config/obsave/config` to store default settings, such as the vault path. You can edit this file to set your preferred default vault path, or override it using the `--vault` flag.

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
   If you want to install `obsave` to your system’s `$GOPATH/bin` directory for global usage, run:
   ```bash
   go install
   ```

   After installation, you can use `obsave` from anywhere in your terminal:
   ```bash
   obsave --name "NewNote" --vault "~/ObsidianVault"
   ```

