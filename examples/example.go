package main

import (
	"fmt"

	"github.com/heroku/go-pingdom/pingdom"
)

func main() {
	client := pingdom.NewClient("username", "password", "api_key")

	// List all checks
	checks, _ := client.Checks.List()
	fmt.Println("All checks:", checks)

	// Create a new http check
	newCheck := pingdom.HttpCheck{BaseCheck: pingdom.BaseCheck{Name: "Test Check", Host: "example.com", Resolution: pingdom.OptInt(5)}}
	check, _ := client.Checks.Create(&newCheck)
	fmt.Println("Created check:", check) // {ID, Name}

	// Create a new ping check
	newPingCheck := pingdom.PingCheck{BaseCheck: pingdom.BaseCheck{Name: "Test Ping", Host: "example.com", Resolution: pingdom.OptInt(1)}}
	pingcheck, _ := client.Checks.Create(&newPingCheck)
	fmt.Println("Created check:", pingcheck) // {ID, Name}

	// Get details for a check
	details, _ := client.Checks.Read(check.ID)
	fmt.Println("Details:", details)

	// Update a check
	updatedCheck := pingdom.HttpCheck{BaseCheck: pingdom.BaseCheck{Name: "Updated Check", Host: "example2.com", Resolution: pingdom.OptInt(5)}}
	upMsg, _ := client.Checks.Update(check.ID, &updatedCheck)
	fmt.Println("Modified check, message:", upMsg)

	// Delete a check
	delMsg, _ := client.Checks.Delete(check.ID)
	fmt.Println("Deleted check, message:", delMsg)

}
