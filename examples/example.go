package main

import (
	"fmt"

	"github.com/russellcardullo/pingdom"
)

func main() {
	client := pingdom.NewClient("username", "password", "api_key")

	// List all checks
	checks, _ := client.ListChecks()
	fmt.Println("All checks:", checks)

	// Create a new check
	check, _ := client.CreateCheck(pingdom.Check{Name: "Test Check", Hostname: "example.com"})
	fmt.Println("Created check:", check)

	// Get details for a check
	details, _ := client.ReadCheck(check.ID)
	fmt.Println("Details:", details)

	// Update a check
	check, _ := client.UpdateCheck(check.ID, pingdom.Check{Name: "Modified Check", Hostname: "example2.com"})
	fmt.Println("Modified check, message:", upMsg)

	// Delete a check
	delMsg, _ := client.DeleteCheck(check.ID)
	fmt.Println("Deleted check, message:", delMsg)
}
