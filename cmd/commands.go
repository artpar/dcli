// cmd/commands.go

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"jsonapi-cli-llm/api"
	"jsonapi-cli-llm/models"
	"jsonapi-cli-llm/utils"
	"os"
	"strings"
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
		fmt.Println("Expected 'create', 'read', 'update', 'delete', 'list', 'relation' subcommands")
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
	case "list":
		listCommand(client, os.Args[2:])
	case "relation":
		relationCommand(client, os.Args[2:])
	default:
		fmt.Println("Expected 'create', 'read', 'update', 'delete', 'list', 'relation' subcommands")
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

func listCommand(client *api.Client, args []string) {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	resourceType := listCmd.String("type", "", "Resource type")
	pageNumber := listCmd.String("page[number]", "", "Page number")
	pageSize := listCmd.String("page[size]", "", "Page size")
	filters := listCmd.String("filter", "", "Filters in key1:value1,key2:value2 format")
	sort := listCmd.String("sort", "", "Sort fields, e.g., 'name,-created_at'")
	include := listCmd.String("include", "", "Related resources to include")
	fields := listCmd.String("fields", "", "Fields to return in key1:field1,field2;key2:field3 format")
	listCmd.Parse(args)

	if *resourceType == "" {
		listCmd.Usage()
		os.Exit(1)
	}

	options := &api.ListOptions{
		Page:   make(map[string]string),
		Filter: make(map[string]string),
		Fields: make(map[string]string),
	}

	if *pageNumber != "" {
		options.Page["number"] = *pageNumber
	}
	if *pageSize != "" {
		options.Page["size"] = *pageSize
	}
	if *filters != "" {
		filterPairs := strings.Split(*filters, ",")
		for _, pair := range filterPairs {
			kv := strings.SplitN(pair, ":", 2)
			if len(kv) == 2 {
				options.Filter[kv[0]] = kv[1]
			}
		}
	}
	if *sort != "" {
		options.Sort = *sort
	}
	if *include != "" {
		options.Include = *include
	}
	if *fields != "" {
		fieldGroups := strings.Split(*fields, ";")
		for _, group := range fieldGroups {
			kv := strings.SplitN(group, ":", 2)
			if len(kv) == 2 {
				options.Fields[kv[0]] = kv[1]
			}
		}
	}

	doc, err := client.List(*resourceType, options)
	if err != nil {
		utils.ErrorLogger.Println("Failed to list resources:", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		utils.ErrorLogger.Println("Failed to marshal response:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func relationCommand(client *api.Client, args []string) {
	if len(args) < 1 {
		fmt.Println("Expected 'get', 'update', 'add', 'remove' subcommands")
		os.Exit(1)
	}

	switch args[0] {
	case "get":
		getRelationCommand(client, args[1:])
	case "update":
		updateRelationCommand(client, args[1:])
	case "add":
		addRelationCommand(client, args[1:])
	case "remove":
		removeRelationCommand(client, args[1:])
	default:
		fmt.Println("Expected 'get', 'update', 'add', 'remove' subcommands")
		os.Exit(1)
	}
}

func getRelationCommand(client *api.Client, args []string) {
	getRelCmd := flag.NewFlagSet("relation get", flag.ExitOnError)
	resourceType := getRelCmd.String("type", "", "Resource type")
	id := getRelCmd.String("id", "", "Resource ID")
	relation := getRelCmd.String("relation", "", "Relation name")
	getRelCmd.Parse(args)

	if *resourceType == "" || *id == "" || *relation == "" {
		getRelCmd.Usage()
		os.Exit(1)
	}

	doc, err := client.FetchRelations(*resourceType, *id, *relation)
	if err != nil {
		utils.ErrorLogger.Println("Failed to fetch relation:", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		utils.ErrorLogger.Println("Failed to marshal response:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func updateRelationCommand(client *api.Client, args []string) {
	updateRelCmd := flag.NewFlagSet("relation update", flag.ExitOnError)
	resourceType := updateRelCmd.String("type", "", "Resource type")
	id := updateRelCmd.String("id", "", "Resource ID")
	relation := updateRelCmd.String("relation", "", "Relation name")
	data := updateRelCmd.String("data", "", "Relation data in JSON format")
	updateRelCmd.Parse(args)

	if *resourceType == "" || *id == "" || *relation == "" || *data == "" {
		updateRelCmd.Usage()
		os.Exit(1)
	}

	var relationData interface{}
	err := json.Unmarshal([]byte(*data), &relationData)
	if err != nil {
		utils.ErrorLogger.Println("Invalid data JSON:", err)
		os.Exit(1)
	}

	doc, err := client.UpdateRelationship(*resourceType, *id, *relation, relationData)
	if err != nil {
		utils.ErrorLogger.Println("Failed to update relation:", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		utils.ErrorLogger.Println("Failed to marshal response:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func addRelationCommand(client *api.Client, args []string) {
	addRelCmd := flag.NewFlagSet("relation add", flag.ExitOnError)
	resourceType := addRelCmd.String("type", "", "Resource type")
	id := addRelCmd.String("id", "", "Resource ID")
	relation := addRelCmd.String("relation", "", "Relation name")
	data := addRelCmd.String("data", "", "Relation data in JSON format")
	addRelCmd.Parse(args)

	if *resourceType == "" || *id == "" || *relation == "" || *data == "" {
		addRelCmd.Usage()
		os.Exit(1)
	}

	var relationData interface{}
	err := json.Unmarshal([]byte(*data), &relationData)
	if err != nil {
		utils.ErrorLogger.Println("Invalid data JSON:", err)
		os.Exit(1)
	}

	doc, err := client.AddToRelationship(*resourceType, *id, *relation, relationData)
	if err != nil {
		utils.ErrorLogger.Println("Failed to add to relation:", err)
		os.Exit(1)
	}

	output, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		utils.ErrorLogger.Println("Failed to marshal response:", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}

func removeRelationCommand(client *api.Client, args []string) {
	removeRelCmd := flag.NewFlagSet("relation remove", flag.ExitOnError)
	resourceType := removeRelCmd.String("type", "", "Resource type")
	id := removeRelCmd.String("id", "", "Resource ID")
	relation := removeRelCmd.String("relation", "", "Relation name")
	data := removeRelCmd.String("data", "", "Relation data in JSON format")
	removeRelCmd.Parse(args)

	if *resourceType == "" || *id == "" || *relation == "" || *data == "" {
		removeRelCmd.Usage()
		os.Exit(1)
	}

	var relationData interface{}
	err := json.Unmarshal([]byte(*data), &relationData)
	if err != nil {
		utils.ErrorLogger.Println("Invalid data JSON:", err)
		os.Exit(1)
	}

	err = client.DeleteFromRelationship(*resourceType, *id, *relation, relationData)
	if err != nil {
		utils.ErrorLogger.Println("Failed to remove from relation:", err)
		os.Exit(1)
	}

	fmt.Println("Relation updated successfully.")
}
