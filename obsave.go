package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	VaultPath     string            `yaml:"vault_path"`  
	OverwriteMode string            `yaml:"overwrite_mode"` 
	Tags          []string          `yaml:"tags"`        
	Properties    map[string]string `yaml:"properties"`  
	Debug         bool              `yaml:"debug"`
	DryRun        bool              `yaml:"dry_run"`
	TagsHandling  string            `yaml:"tags_handling"` 
	PropertiesHandling string `yaml:"properties_handling"`
	Name          string            `yaml:"name"`
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

// Function to parse the custom frontmatter string
func parseCustomClasses(classString string) (map[string]string, error) {
	frontmatter := make(map[string]string)

	if classString == "" {
		return frontmatter, nil
	}

	pairs := strings.Split(classString, ";")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid pair: %s", pair)
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		frontmatter[key] = value
	}

	return frontmatter, nil
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

func getConfigPath() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}
	return filepath.Join(usr.HomeDir, configDir, configFile)
}


func loadConfig(configFile string) (*Config, error) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}

	// Build the path to the config file using the configName
	configPath := filepath.Join(usr.HomeDir, configDir, configFile)
	debugLog("Config Path : " + configPath)

	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %v", err)
	}
	debugLog("Config loaded : " + configFile)

	return &config, nil
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


func saveToObsidian(content string, config *Config) error {
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
					return fmt.Errorf("file '%s' already exists. Use --overwrite-mode=overwrite or --overwrite-mode=serialize", fileName)
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
			return nil
	}

	// Write the content to the file
	file, err := os.Create(filePath)
	if err != nil {
			return err
	}
	defer file.Close()

	// Prepend the frontmatter to the content
	_, err = file.WriteString(frontmatter + "\n" + content)
	if err != nil {
			return err
	}

	debugLog("Note saved successfully at: " + filePath)
	return nil
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


func main() {
	// Command-line flags
	name := flag.String("name", "", "Name of the note")
	tags := flag.String("tags", "", "Comma-separated list of tags")
	properties := flag.String("properties", "", "Custom frontmatter properties key:value pairs (e.g., author=John;status=Draft)")
	vaultPath := flag.String("vault", "", "Path to Obsidian vault folder")
	overwriteMode := flag.String("overwrite-mode", "fail", "Overwrite mode: 'overwrite' or 'serialize'")
	debugFlag := flag.Bool("debug", false, "Enable debug mode")
	dryRun := flag.Bool("dry-run", false, "Simulate the run without writing files")
	tagsHandling := flag.String("tags-handling", "merge", "Tags handling mode: 'replace', 'add', or 'merge'")
	propertiesHandling := flag.String("properties-handling", "merge", "Properties handling mode: 'replace', 'add', or 'merge'")
	flag.Parse()


	// Check for a non-hyphen-prefixed positional argument for the config file
	configName := configFile // Default config file
	if len(os.Args) > 1 && !strings.HasPrefix(os.Args[1], "-") {
		configName = os.Args[1] // Use the positional argument as the config name
		// Re-parse the remaining arguments after the config name
		flag.CommandLine.Parse(os.Args[2:])
	} else {
		// Parse command-line flags normally
		flag.Parse()
	}

	// Load config values and merge them with command-line options
	config, err := loadConfig(configName)
	if err != nil {
			log.Fatalf("Error loading config: %v", err)
	}

	// Overwrite config values with command-line flags (if provided)
	if *vaultPath != "" {
		config.VaultPath = *vaultPath
	}
	if *overwriteMode != "" {
		config.OverwriteMode = *overwriteMode
	}
	if *debugFlag {
		config.Debug = true
	}
	if *dryRun {
		config.DryRun = true
	}
	if *tagsHandling != "" {
		config.TagsHandling = *tagsHandling
	}
	if *propertiesHandling != "" {
		config.PropertiesHandling = *propertiesHandling
	}

	switch config.PropertiesHandling {
	case "wipe": // or "replace"
			replaceProperties(config, *properties)
	case "add":
			addProperties(config, *properties)
	case "merge": // or "merge replace"
			mergeProperties(config, *properties)
	}

	switch config.TagsHandling {
	case "replace":
		replaceTags(config, *tags)
	case "add":
		addTags(config, *tags)
	case "merge":
		mergeTags(config, *tags)
	}

  // Enable debug mode based on flag
	if config.Debug {
		debugMode = true
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		debugLog("Debug mode enabled")
	}

	if config.VaultPath == "" {
		fmt.Println("Error: Vault path is required either through config or --vault argument.")
			os.Exit(1)
	}
	

	// Expand the "~" if used in the vault path
	expandedVaultPath, err := expandAndCleanPath(config.VaultPath)
	if err != nil {
		fmt.Println("Error expanding vault path:", err)
		os.Exit(1)
	}
	config.VaultPath = expandedVaultPath
	debugLog("Vault path expanded: " + expandedVaultPath)

	// Check for required arguments
	if *name == "" {
		fmt.Println("Error: Name is required.")
		flag.Usage()
		os.Exit(1)
	}

	config.Name = *name
	
	// Read piped input (from stdin)
	scanner := bufio.NewScanner(os.Stdin)
	var contentBuilder strings.Builder
	for scanner.Scan() {
		contentBuilder.WriteString(scanner.Text() + "\n")
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
	content := contentBuilder.String()
	debugLog("Content read from stdin: " + content)

	// Save the content to the Obsidian vault, or simulate if dry-run is enabled
	err = saveToObsidian(content, config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving note: %v\n", err)
		os.Exit(1)
	}
}
