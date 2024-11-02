## Shell Completion

### ZSH Completion

Obsave includes a completion script for zsh that provides:
- Automatic completion of command flags and options
- Dynamic completion of config files from your obsave config directory
- Directory completion for vault paths
- Predefined completions for modes (overwrite, tags, properties handling)

#### Installation

1. Create the completions directory if it doesn't exist:
```bash
mkdir -p ~/.zsh/completions
```

2. Copy the completion script:
```bash
cp completions/_obsave ~/.zsh/completions/
```

3. Add the completions directory to your fpath by adding this line to your `~/.zshrc`:
```bash
fpath=(~/.zsh/completions $fpath)
```

4. Reload completion scripts:
```bash
autoload -U compinit && compinit
```

#### Usage Examples

After installation, you can use TAB completion with obsave:

```bash
obsave -[TAB]                 # Show all available flags
obsave --config [TAB]         # Show available config files
obsave --vault [TAB]          # Browse directories for vault path
obsave --overwrite-mode [TAB] # Show available overwrite modes
```

The completion script will also show brief descriptions of each option as you tab through them.
