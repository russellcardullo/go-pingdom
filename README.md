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
fmt.Println("Checks:", checks) // [{ID Name} ...]
```

Create a new HTTP check:

```go
newCheck := pingdom.Check{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.CreateCheck(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Get details for a specific check

```go
check, err := client.GetCheck(12345)
```

Update a check

```go
updatedCheck := pingdom.Check{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
msg, err := client.UpdateCheck(12345, &updatedCheck)
```

Delete a check

```go
msg, err := client.DeleteCheck(12345)
```

