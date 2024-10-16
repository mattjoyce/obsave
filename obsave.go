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

	"github.com/gookit/ini/v2"
)

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

// Function to expand the home directory if "~" is used in the path
func expandPath(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		usr, err := user.Current()
		if err != nil {
			return "", err
		}
		return filepath.Join(usr.HomeDir, strings.TrimPrefix(path, "~")), nil
	}
	return path, nil
}

func getConfigPath() string {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("Error getting current user:", err)
		os.Exit(1)
	}
	return filepath.Join(usr.HomeDir, configDir, configFile)
}

// Function to load config file for default vault path using gookit/ini
func loadConfig() string {
	configPath := getConfigPath()

	err := ini.LoadExists(configPath)
	if err != nil {
		fmt.Printf("Warning: No config file found at %s. Please set the vault manually.\n", configPath)
		return ""
	}

	vaultPath := ini.String("VaultPath")
	debugLog("Loaded vault path from config: " + vaultPath)
	return vaultPath
}

// Function to create frontmatter
func createFrontmatter(name, tags string) string {
	frontmatter := fmt.Sprintf("---\ntitle: %s\n", name)
	if tags != "" {
		tagList := strings.Split(tags, ",")
		for i := range tagList {
			tagList[i] = strings.TrimSpace(tagList[i])
		}
		frontmatter += fmt.Sprintf("tags: [%s]\n", strings.Join(tagList, ", "))
	}
	frontmatter += fmt.Sprintf("date: %s\n---\n", time.Now().Format("2006-01-02 15:04"))
	debugLog("Frontmatter created: " + frontmatter)
	return frontmatter
}

// Function to save content to Obsidian vault
func saveToObsidian(name, content, tags, vaultPath, overwriteMode string) error {
	fileName := name + ".md"
	filePath := filepath.Join(vaultPath, fileName)

	// Handle file existence cases
	if _, err := os.Stat(filePath); err == nil {
		if overwriteMode == "serialize" {
			baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			counter := 1
			for {
				serializedFileName := fmt.Sprintf("%s_%d.md", baseName, counter)
				serializedFilePath := filepath.Join(vaultPath, serializedFileName)
				if _, err := os.Stat(serializedFilePath); os.IsNotExist(err) {
					filePath = serializedFilePath
					break
				}
				counter++
			}
			debugLog("Serialized file path: " + filePath)
		} else if overwriteMode != "overwrite" {
			return fmt.Errorf("file '%s' already exists. Use --overwrite-mode=overwrite or --overwrite-mode=serialize", fileName)
		} else {
			debugLog("Overwriting existing file: " + filePath)
		}
	} else {
		debugLog("New file created: " + filePath)
	}

	// Create frontmatter and save to the file
	frontmatter := createFrontmatter(name, tags)
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
	fmt.Printf("Note saved to %s\n", filePath)
	return nil
}

func main() {
	// Command-line flags
	name := flag.String("n", "", "Name of the note (short: -n, long: --name)")
	tags := flag.String("t", "", "Comma-separated list of tags (short: -t, long: --tag)")
	vaultPath := flag.String("v", "", "Path to Obsidian vault folder (short: -v, long: --vault)")
	overwriteMode := flag.String("overwrite-mode", "fail", "Overwrite mode: 'overwrite' or 'serialize'")
	debugFlag := flag.Bool("debug", false, "Enable debug mode (long: --debug)")
	flag.Parse()

	// Enable debug mode based on flag
	debugMode = *debugFlag
	if debugMode {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		debugLog("Debug mode enabled")
	}

	// Load config vault path if not provided via argument
	if *vaultPath == "" {
		*vaultPath = loadConfig()
		if *vaultPath == "" {
			fmt.Println("Error: Vault path is required either through config or --vault argument.")
			os.Exit(1)
		}
	}

	// Expand the "~" if used in the vault path
	expandedVaultPath, err := expandPath(*vaultPath)
	if err != nil {
		fmt.Println("Error expanding vault path:", err)
		os.Exit(1)
	}
	debugLog("Vault path expanded: " + expandedVaultPath)

	// Check for required arguments
	if *name == "" {
		fmt.Println("Error: Name is required.")
		flag.Usage()
		os.Exit(1)
	}

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

	// Save the content to the Obsidian vault
	err = saveToObsidian(*name, content, *tags, expandedVaultPath, *overwriteMode)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving note: %v\n", err)
		os.Exit(1)
	}
}
