package main

import (
	"github.com/russellcardullo/go-pingdom/pingdom"
	"fmt"
)

func main() {
	client := pingdom.NewMultiUserClient("user", "password", "apikey", "account-email")

	// List all checks
	users, _ := client.Users.List()
	fmt.Println("All users:", users)
}