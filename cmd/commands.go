// cmd/commands.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"your_project/api"
	"your_project/models"
	"your_project/utils"
)

func main() {
	// Initialize logger
	utils.InitLogger(false)

	// Load configuration
	config, err := utils.LoadConfig("")
	if err != nil {
		utils.ErrorLogger.Println("Failed to load config:", err)
		os.Exit(1)
	}

	// Create API client
	client, err := api.NewClient(config.BaseURL)
	if err != nil {
		utils.ErrorLogger.Println("Failed to create API client:", err)
		os.Exit(1)
	}

	// Parse command-line arguments
	if len(os.Args) < 2 {
		fmt.Println("Expected 'create', 'read', 'update', 'delete' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "create":
		createCommand(client, os.Args[2:])
	case "read":
		readCommand(client, os.Args[2:])
	case "update":
		updateCommand(client, os.Args[2:])
	case "delete":
		deleteCommand(client, os.Args[2:])
	default:
		fmt.Println("Expected 'create', 'read', 'update', 'delete' subcommands")
		os.Exit(1)
	}
}

func createCommand(client *api.Client, args []string) {
	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
	resourceType := createCmd.String("type", "", "Resource type")
	attributes := createCmd.String("attributes", "", "Resource attributes in JSON format")
	createCmd.Parse(args)

	if *resourceType == "" || *attributes == "" {
		createCmd.Usage()
		os.Exit(1)
	}

	var attrs map[string]interface{}
	err := json.Unmarshal([]byte(*attributes), &attrs)
	if err != nil {
		utils.ErrorLogger.Println("Invalid attributes JSON:", err)
		os.Exit(1)
	}

	resource := &models.Resource{
		Type:       *resourceType,
		Attributes: attrs,
	}

	createdResource, err := client.Create(resource)
	if err != nil {
		utils.ErrorLogger.Println("Failed to create resource:", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(createdResource, "", "  ")
	if err != nil {
		utils.ErrorLogger.Println("Failed to marshal response:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func readCommand(client *api.Client, args []string) {
	readCmd := flag.NewFlagSet("read", flag.ExitOnError)
	resourceType := readCmd.String("type", "", "Resource type")
	id := readCmd.String("id", "", "Resource ID")
	readCmd.Parse(args)

	if *resourceType == "" || *id == "" {
		readCmd.Usage()
		os.Exit(1)
	}

	resource, err := client.Read(*resourceType, *id)
	if err != nil {
		utils.ErrorLogger.Println("Failed to read resource:", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(resource, "", "  ")
	if err != nil {
		utils.ErrorLogger.Println("Failed to marshal response:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func updateCommand(client *api.Client, args []string) {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	resourceType := updateCmd.String("type", "", "Resource type")
	id := updateCmd.String("id", "", "Resource ID")
	attributes := updateCmd.String("attributes", "", "Resource attributes in JSON format")
	updateCmd.Parse(args)

	if *resourceType == "" || *id == "" || *attributes == "" {
		updateCmd.Usage()
		os.Exit(1)
	}

	var attrs map[string]interface{}
	err := json.Unmarshal([]byte(*attributes), &attrs)
	if err != nil {
		utils.ErrorLogger.Println("Invalid attributes JSON:", err)
		os.Exit(1)
	}

	resource := &models.Resource{
		Type:       *resourceType,
		ID:         *id,
		Attributes: attrs,
	}

	updatedResource, err := client.Update(resource)
	if err != nil {
		utils.ErrorLogger.Println("Failed to update resource:", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(updatedResource, "", "  ")
	if err != nil {
		utils.ErrorLogger.Println("Failed to marshal response:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func deleteCommand(client *api.Client, args []string) {
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	resourceType := deleteCmd.String("type", "", "Resource type")
	id := deleteCmd.String("id", "", "Resource ID")
	deleteCmd.Parse(args)

	if *resourceType == "" || *id == "" {
		deleteCmd.Usage()
		os.Exit(1)
	}

	err := client.Delete(*resourceType, *id)
	if err != nil {
		utils.ErrorLogger.Println("Failed to delete resource:", err)
		os.Exit(1)
	}

	fmt.Println("Resource deleted successfully.")
}
