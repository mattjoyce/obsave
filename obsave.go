package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	VaultPath          string            `yaml:"vault_path"`
	OverwriteMode      string            `yaml:"overwrite_mode"`
	Tags               []string          `yaml:"tags"`
	Properties         map[string]string `yaml:"properties"`
	Debug              bool              `yaml:"debug"`
	DryRun             bool              `yaml:"dry_run"`
	TagsHandling       string            `yaml:"tags_handling"`
	PropertiesHandling string            `yaml:"properties_handling"`
	Name               string            `yaml:"name"`
	Passthrough        bool              `yaml:"passthrough"`
	Verbose            bool              `yaml:"verbose"`
}

// configs should be in ~/.config/obsave/

const configDir = ".config/obsave"
const configFile = "config"

// Global variable to track if debug mode is enabled
var debugMode bool

// Function to log debug messages when debug mode is enabled
func debugLog(message string) {
	if debugMode {
		log.Println("[DEBUG]", message)
	}
}

// Function to expand home directories and clean up paths
func expandAndCleanPath(path string) (string, error) {
	// Handle home directory expansion (for Unix-like systems)
	if path[:2] == "~/" {
		usr, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = filepath.Join(usr, path[2:])
	}

	// Clean up the path to remove unnecessary parts like "./", "../", etc.
	cleanPath := filepath.Clean(path)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func setDefaultsAndOverrides(config *Config, overwriteModeFlag string) {
	// Default values if neither config nor flag is set
	if config.OverwriteMode == "" {
		config.OverwriteMode = "fail" // Base default
	}

	// Command line flag overrides everything if specified
	if overwriteModeFlag != "" {
		debugLog(fmt.Sprintf("Overwrite mode set from: %s to %s", config.OverwriteMode, overwriteModeFlag))
		config.OverwriteMode = overwriteModeFlag
	}

	// Validate the final value
	switch config.OverwriteMode {
	case "fail", "overwrite", "serialize":
		// Valid values
	default:
		log.Printf("Invalid overwrite mode '%s', using default 'fail'", config.OverwriteMode)
		config.OverwriteMode = "fail"
	}

	debugLog(fmt.Sprintf("Final overwrite mode: %s", config.OverwriteMode))
}

func loadConfig(configFile string) (*Config, *string, error) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	// Build the path to the config file using the configName
	configPath := filepath.Join(usr.HomeDir, configDir, configFile)
	debugLog("Config Path : " + configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, &configPath, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}
	debugLog("Config loaded : " + configFile)

	return &config, &configPath, nil
}

func createFrontmatter(config *Config) string {
	frontmatter := fmt.Sprintf("---\ntitle: %s\n", config.Name)

	// Handle tags if they exist in the config
	if len(config.Tags) > 0 {
		frontmatter += fmt.Sprintf("tags: [%s]\n", strings.Join(config.Tags, ", "))
	}

	// Add the date
	frontmatter += fmt.Sprintf("date: %s\n", time.Now().Format("2006-01-02"))

	// Append custom frontmatter properties (key-value pairs from config.Properties)
	for key, value := range config.Properties {
		frontmatter += fmt.Sprintf("%s: %s\n", key, value)
	}

	frontmatter += "---\n"
	return frontmatter
}

func saveToObsidian(content string, config *Config) (string, error) {
	fileName := config.Name + ".md"
	filePath := filepath.Join(config.VaultPath, fileName)

	// Handle file existence cases
	if _, err := os.Stat(filePath); err == nil {
		if config.OverwriteMode == "serialize" {
			baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			counter := 1
			for {
				serializedFileName := fmt.Sprintf("%s_%d.md", baseName, counter)
				serializedFilePath := filepath.Join(config.VaultPath, serializedFileName)
				if _, err := os.Stat(serializedFilePath); os.IsNotExist(err) {
					filePath = serializedFilePath
					break
				}
				counter++
			}
			debugLog("Serialized file path: " + filePath)
		} else if config.OverwriteMode != "overwrite" {
			return "", fmt.Errorf("file '%s' already exists. Use --overwrite-mode=overwrite or --overwrite-mode=serialize", fileName)
		} else {
			debugLog("Overwriting existing file: " + filePath)
		}
	} else {
		debugLog("New file created: " + filePath)
	}

	// Create frontmatter using the updated createFrontmatter function
	frontmatter := createFrontmatter(config)

	if config.DryRun {
		fmt.Println("Dry-run: The following content would be saved:")
		fmt.Println(frontmatter + "\n" + content)
		return filePath, nil
	}

	// Write the content to the file
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Prepend the frontmatter to the content
	_, err = file.WriteString(frontmatter + "\n" + content)
	if err != nil {
		return "", err
	}

	debugLog("Note saved successfully at: " + filePath)
	return filePath, nil
}

func replaceProperties(config *Config, cliProperties string) {
	// Initialize or clear the map
	config.Properties = make(map[string]string)

	// Split the CLI properties into key-value pairs
	pairs := strings.Split(cliProperties, ";")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			config.Properties[key] = value
		}
	}
}

func addProperties(config *Config, cliProperties string) {
	// Split the CLI properties into key-value pairs
	pairs := strings.Split(cliProperties, ";")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			if _, exists := config.Properties[key]; !exists {
				config.Properties[key] = value
			}
		}
	}
}

func mergeProperties(config *Config, cliProperties string) {
	// Split the CLI properties into key-value pairs
	pairs := strings.Split(cliProperties, ";")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			config.Properties[key] = value
		}
	}
}

func replaceTags(config *Config, cliTags string) {
	// Split the CLI tags into a slice of strings
	config.Tags = strings.Split(cliTags, ",")

	// Trim spaces around each tag
	for i, tag := range config.Tags {
		config.Tags[i] = strings.TrimSpace(tag)
	}
}

func addTags(config *Config, cliTags string) {
	// Split the CLI tags into a slice of strings
	newTags := strings.Split(cliTags, ",")

	// Trim spaces and add only new tags
	for _, tag := range newTags {
		tag = strings.TrimSpace(tag)
		if !contains(config.Tags, tag) { // Ensure no duplicates
			config.Tags = append(config.Tags, tag)
		}
	}
}

// Helper function to check if a tag already exists
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func mergeTags(config *Config, cliTags string) {
	// Split the CLI tags into a slice of strings
	newTags := strings.Split(cliTags, ",")

	// Trim spaces and add all new tags (duplicates allowed)
	for _, tag := range newTags {
		tag = strings.TrimSpace(tag)
		config.Tags = append(config.Tags, tag)
	}
}

func printExtendedHelp() {
	fmt.Println(`
Obsave - Obsidian Note Creation Utility

DESCRIPTION:
    A utility for creating and managing notes in an Obsidian vault with flexible
    frontmatter handling, tag management, and property customization.

BASIC USAGE:
    echo "Your note content" | obsave -n "Note Title" -ob ~/vault/path

CONFIGURATION:
    Default config location: ~/.config/obsave/config
    Config precedence:
    1. Default config (if exists)
    2. Specified config file (-c/--config)
    3. Command line options

EXAMPLES:
    # Create a simple note
    echo "Meeting minutes" | obsave -n "Team Meeting" -ob ~/Notes

    # Use tags and properties
    echo "Project specs" | obsave -n "Project Alpha" -ob ~/Notes \
        -t "project,specs" -p "status=draft;priority=high"

    # Use a specific config file
    echo "Custom note" | obsave -c my-config -n "Custom Note"

    # Merge new tags with config defaults
    echo "Tagged content" | obsave -n "Tagged Note" -ob ~/Notes \
        -t "new-tag" --tags-handling merge

FRONTMATTER:
    The generated note includes YAML frontmatter with:
    - title: from name option
    - date: auto-generated (YYYY-MM-DD)
    - tags: from config and/or command line
    - custom properties: from config and/or command line

For more information and examples, visit:
https://github.com/mattjoyce/obsave`)
}

func init() {
	// Customize the usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nFor detailed help and examples, use -h or --help\n")
	}
}

func main() {
	// Command-line flags with both long and short forms

	// Add explicit help flag (in addition to automatic -h)
	help := flag.Bool("help", false, "Display detailed help information")

	var name string
	flag.StringVar(&name, "name", "", "Name of the note")
	flag.StringVar(&name, "n", "", "Name of the note (shorthand)")

	var tags string
	flag.StringVar(&tags, "tags", "", "Comma-separated list of tags")
	flag.StringVar(&tags, "t", "", "Comma-separated list of tags (shorthand)")

	var properties string
	flag.StringVar(&properties, "properties", "", "Custom frontmatter properties key:value pairs (e.g., author=John;status=Draft)")
	flag.StringVar(&properties, "p", "", "Custom frontmatter properties key:value pairs (shorthand)")

	var vaultPath string
	flag.StringVar(&vaultPath, "vault", "", "Path to Obsidian vault folder")
	flag.StringVar(&vaultPath, "ob", "", "Path to Obsidian vault folder (shorthand)")

	var verbose bool
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose mode")
	flag.BoolVar(&verbose, "v", false, "Enable verbose mode (shorthand)")

	var configFileFlag string
	flag.StringVar(&configFileFlag, "config", "", "Name of the config file to use")
	flag.StringVar(&configFileFlag, "c", "", "Name of the config file to use (shorthand)")

	// Existing flags without short forms - removed defaults where appropriate
	var overwriteMode string
	flag.StringVar(&overwriteMode, "overwrite-mode", "", "Overwrite mode: 'overwrite' or 'serialize'")
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	dryRun := flag.Bool("dry-run", false, "Simulate the run without writing files")
	var tagsHandling string
	flag.StringVar(&tagsHandling, "tags-handling", "", "Tags handling mode: 'replace', 'add', or 'merge'")
	var propertiesHandling string
	flag.StringVar(&propertiesHandling, "properties-handling", "", "Properties handling mode: 'replace', 'add', or 'merge'")

	var passthrough bool
	flag.BoolVar(&passthrough, "passthrough", false, "Pass input through to stdout while saving to vault")

	// Parse command-line flags
	flag.Parse()

	// Check for help flag first
	if *help {
		printExtendedHelp()
		os.Exit(0)
	}

	// Initialize config with safe defaults
	config := &Config{
		OverwriteMode:      "fail",  // Safe default
		TagsHandling:       "merge", // Safe default
		PropertiesHandling: "merge", // Safe default
	}

	// Enable debug mode if set
	if *debugFlag {
		debugMode = true
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		debugLog("Debug mode enabled")
	}

	// 1. Try to load default config if it exists
	defaultConfig, configConfigPath, err := loadConfig(configFile)
	if err == nil {
		// Only use default config if it was successfully loaded
		config = defaultConfig
		debugLog("Default config loaded from: " + *configConfigPath)
	}

	// 2. If a specific config file was provided, load and use it instead
	if configFileFlag != "" {
		specifiedConfig, specifiedConfigPath, err := loadConfig(configFileFlag)
		if err != nil {
			log.Fatalf("Error loading specified config file: %v, from path: %s", err, *specifiedConfigPath)
		}
		config = specifiedConfig
		debugLog("Specified config loaded from: " + *specifiedConfigPath)
	}

	// 3. Override with command-line options if provided
	if vaultPath != "" {
		debugLog(fmt.Sprintf("Vault path set from: %s to %s", config.VaultPath, vaultPath))
		config.VaultPath = vaultPath
	}
	if overwriteMode != "" {
		debugLog(fmt.Sprintf("Overwrite mode set from: %s to %s", config.OverwriteMode, overwriteMode))
		config.OverwriteMode = overwriteMode
	}
	if *debugFlag {
		config.Debug = true
	}
	if *dryRun {
		debugLog("Dry run enabled")
		config.DryRun = true
	}
	if tagsHandling != "" {
		debugLog(fmt.Sprintf("Tags handling set from: %s to %s", config.TagsHandling, tagsHandling))
		config.TagsHandling = tagsHandling
	}
	if propertiesHandling != "" {
		debugLog(fmt.Sprintf("Properties handling set from: %s to %s", config.PropertiesHandling, propertiesHandling))
		config.PropertiesHandling = propertiesHandling
	}
	if verbose {
		config.Verbose = true
	}
	if name != "" {
		config.Name = name
	}
	if tags != "" {
		debugLog(fmt.Sprintf("Old tags: [%s]", strings.Join(config.Tags, ", ")))
		switch config.TagsHandling {
		case "replace":
			debugLog("Tags handling mode: replace")
			replaceTags(config, tags)
		case "add":
			debugLog("Tags handling mode: add")
			addTags(config, tags)
		default: // "merge" is default
			debugLog("Tags handling mode: merge")
			mergeTags(config, tags)
		}
		debugLog("New tags: " + tags)
	}
	if properties != "" {
		debugLog("Old properties:")
		for k, v := range config.Properties {
			debugLog(fmt.Sprintf("  %s: %s", k, v))
		}

		switch config.PropertiesHandling {
		case "replace":
			debugLog("Properties handling mode: replace")
			replaceProperties(config, properties)
		case "add":
			debugLog("Properties handling mode: add")
			addProperties(config, properties)
		default: // "merge" is default
			debugLog("Properties handling mode: merge")
			mergeProperties(config, properties)
		}
		debugLog("Final properties:")
		for k, v := range config.Properties {
			debugLog(fmt.Sprintf("  %s: %s", k, v))
		}
	}

	if passthrough {
		config.Passthrough = true
	}

	// 4. Check mandatory options
	mandatoryError := false
	if config.Name == "" {
		fmt.Println("Error: Note name is required (use --name or -n)")
		mandatoryError = true
	}

	if config.VaultPath == "" {
		fmt.Println("Error: Vault path is required (use --vault or -ob)")
		mandatoryError = true
	}

	if mandatoryError {
		fmt.Println("\nUsage information:")
		flag.Usage()
		os.Exit(1)
	}

	if config.Debug {
		debugLog("Final configuration:")
		debugLog(fmt.Sprintf("  Name: %s", config.Name))
		debugLog(fmt.Sprintf("  Vault Path: %s", config.VaultPath))
		debugLog(fmt.Sprintf("  Overwrite Mode: %s", config.OverwriteMode))
		debugLog(fmt.Sprintf("  Tags Handling: %s", config.TagsHandling))
		debugLog(fmt.Sprintf("  Properties Handling: %s", config.PropertiesHandling))
		debugLog(fmt.Sprintf("  Debug: %v", config.Debug))
		debugLog(fmt.Sprintf("  Dry Run: %v", config.DryRun))
		debugLog(fmt.Sprintf("  Verbose: %v", config.Verbose))
	}

	// Expand the "~" if used in the vault path
	expandedVaultPath, err := expandAndCleanPath(config.VaultPath)
	if err != nil {
		fmt.Println("Error expanding vault path:", err)
		os.Exit(1)
	}
	config.VaultPath = expandedVaultPath
	debugLog("Vault path expanded: " + expandedVaultPath)

	// Read piped input (from stdin)
	scanner := bufio.NewScanner(os.Stdin)
	var contentBuilder strings.Builder
	for scanner.Scan() {
		line := scanner.Text()
		contentBuilder.WriteString(scanner.Text() + "\n")
		if config.Passthrough {
			fmt.Println(line)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	content := contentBuilder.String()
	debugLog("Content read from stdin: " + content)

	// Save the content to the Obsidian vault, or simulate if dry-run is enabled
	fullFilename, err := saveToObsidian(content, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving note: %v\n", err)
		os.Exit(1)
	}

	if config.Verbose {
		fmt.Println(fullFilename)
	}

}
