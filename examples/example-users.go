package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"github.com/russellcardullo/go-pingdom/pingdom"
)

type Credentials struct {
	User 			string `json:"user"`
	Password 		string `json:"password"`
	ApiKey 			string `json:"apikey"`
	AccountEmail 	string `json:"accountEmail"`
}

func getConfig() Credentials {
	// Config Example
	// { "user" : "", "password" : "", "apikey" : "", "accountEmail" : "" }
	if len(os.Args) < 2 {
		fmt.Println("You must provide a configuration file.")
	}
	args := os.Args[1:] // get the config file from cli
	// Open our jsonFile
	jsonFile, err := os.Open(args[0])
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	fmt.Println("Successfully Opened: " + args[0])

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)


	// we initialize our Users array
	var config Credentials


	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &config)

	fmt.Println(config)

	return config
}

func main() {
	config := getConfig()
	client := pingdom.NewMultiUserClient(config.User, config.Password, config.ApiKey, config.AccountEmail)

	//List all checks
	users, _ := client.Users.List()
	fmt.Println("All users:", users)


	user := pingdom.User{
		Username : "example-user",
	}
	u, _ := client.Users.Create(&user)
	fmt.Println(u)
}