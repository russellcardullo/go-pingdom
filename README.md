# pingdom-go #

pingdom-go is a Go client library for the Pingdom API.

This currently only supports working with basic HTTP checks.

## Usage ##

Construct a new Pingdom client:
  
```go
client := pingdom.NewClient("pingdom_username", "pingdom_password", "pingdom_api_key")
```

Get a list of all checks:

```go
checks, err := client.ListChecks()
```

Create a new HTTP check:

```go
newCheck := pingdom.HttpCheck{"test_check", "example.com"}
check, err := client.CreateCheck(newCheck)
fmt.Println("Created check:", check) // Check{ID, Name, Host}
```

Get details for a specific check

```go
check, err := client.GetCheck(12345)
```

Update a check

```go
msg, err := client.UpdateCheck(12345, "my check name", "example2.com")
```

Delete a check

```go
msg, err := client.DeleteCheck(12345)
```

