package main

import (
	"fmt"
	"strconv"

	"github.com/russellcardullo/pingdom"
)

func main() {
	client := pingdom.NewClient("username", "password", "api_key")

	// List all checks
	checks, _ := client.ListChecks()
	fmt.Println("All checks:", checks)

	// Create a new check
	check, _ := client.CreateCheck(pingdom.HttpCheck{"test_check", "example.com"})
	fmt.Println("Created check:", check)

	// Get details for a check
	details, _ := client.ReadCheck(strconv.Itoa(check.ID))
	fmt.Println("Details:", details)

	// Update a check
	upMsg, _ := client.UpdateCheck(strconv.Itoa(check.ID), "modified_check", "foo.com")
	fmt.Println("Modified check, message:", upMsg)

	// Delete a check
	delMsg, _ := client.DeleteCheck(strconv.Itoa(check.ID))
	fmt.Println("Deleted check, message:", delMsg)
}
