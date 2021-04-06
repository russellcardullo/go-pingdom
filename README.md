# go-pingdom #

[![Build Status](https://travis-ci.org/nordcloud/go-pingdom.svg?branch=master)](https://travis-ci.org/nordcloud/go-pingdom) [![Go Report Card](https://goreportcard.com/badge/github.com/nordcloud/go-pingdom/pingdom)](https://goreportcard.com/report/github.com/nordcloud/go-pingdom/pingdom) [![GoDoc](https://godoc.org/github.com/nordcloud/go-pingdom/pingdom?status.svg)](https://godoc.org/github.com/nordcloud/go-pingdom/pingdom)

go-pingdom is a Go client library for the Pingdom API.

This currently supports working with HTTP, ping checks, and TCP checks.

**Important**: The current version of this library only supports the Pingdom 3.1 API.  If you are still using the deprecated Pingdom 2.1 API please pin your dependencies to tag v1.1.0 of this library.

## Usage ##

### Client ###

Construct a new Pingdom client:

```go
client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
    APIToken: "pingdom_api_token",
})
```

Using a Pingdom client, you can access supported services.

You can override the timeout or other parameters by passing a custom http client:
```go
client, err := pingdom.NewClientWithConfig(pingdom.ClientConfig{
    APIToken: "pingdom_api_token",
    HTTPClient: &http.Client{
        Timeout: time.Second * 10,
    },
})
```

The `APIToken` can also implicitly be provided by setting the environment variable `PINGDOM_API_TOKEN`:

```bash
export PINGDOM_API_TOKEN=pingdom_api_token
./your_application
```


### Pindom Extension Client ###

Construct a new Pingdom extension client:

```go
client_ext, err := pingdomext.NewClientWithConfig(pingdomext.ClientConfig{
    Username: "test_user",
    Password: "test_pwd",
    HTTPClient: &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    },
})
```

Using a Pingdom extention client, you can access supported services, like integration service.

You must override the CheckRedirect since there have multiple redirect while get the jwt token for access api. 

The `Username` and `Password` can also implicitly be provided by setting the environment variable `SOLARWINDS_USER` and `SOLARWINDS_PASSWD`:

```bash
export SOLARWINDS_USER=test_user
export SOLARWINDS_PASSWD=test_pwd
./your_application
```

### CheckService ###

This service manages pingdom Checks which are represented by the `Check` struct.
When creating or updating Checks you must specify at a minimum the `Name`, `Hostname`
and `Resolution`.  Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Checks from Pingdom: https://www.pingdom.com/features/api/documentation/#ResourceChecks

Get a list of all checks:

```go
checks, err := client.Checks.List()
fmt.Println("Checks:", checks) // [{ID Name} ...]
```

Create a new HTTP check:

```go
newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new Ping check:
```go
newCheck := pingdom.PingCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new TCP check:
```go
newCheck := pingdom.TCPCheck{Name: "Test Check", Hostname: "example.com", Port: 25, StringToSend: "HELO foo.com", StringToExpect: "250 mail.test.com", Resolution: 5}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Create a new DNS check:
```go
newCheck := pingdom.DNSCheck{
    Name: "fake check",
    Hostname: "example.com",
    ExpectedIP: "192.168.1.1",
    NameServer: "8.8.8.8",
}
check, err := client.Checks.Create(&newCheck)
fmt.Println("Created check:", check) // {ID, Name}
```

Get details for a specific check:

```go
checkDetails, err := client.Checks.Read(12345)
```

For checks with detailed information, check the specific details in
the field `Type` (e.g. `checkDetails.Type.HTTP`).

Update a check:

```go
updatedCheck := pingdom.HttpCheck{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
msg, err := client.Checks.Update(12345, &updatedCheck)
```

Delete a check:

```go
msg, err := client.Checks.Delete(12345)
```

Create a check with basic alert notification to a user.

```go
newCheck := pingdom.HttpCheck{Name: "Test Check", Hostname: "example.com", Resolution: 5, SendNotificationWhenDown: 2, UserIds []int{12345}}
checkResponse, err := client.Checks.Create(&newCheck)
```

### MaintenanceService ###

This service manages pingdom Maintenances which are represented by the `Maintenance` struct.
When creating or updating Maintenances you must specify at a minimum the `Description`, `From`
and `To`.  Other fields are optional but if not set will be given the zero
values for the underlying type.

More information on Maintenances from Pingdom: https://www.pingdom.com/resources/api/2.1#ResourceMaintenance

Get a list of all maintenances:

```go
maintenances, err := client.Maintenances.List()
fmt.Println("Maintenances:", maintenances) // [{ID Description} ...]
```

Create a new Maintenance Window:

```go
m := pingdom.MaintenanceWindow{
    Description: "My Maintenance",
    From:        1,
    To:          1234567899,
}
maintenance, err := client.Maintenances.Create(&m)
fmt.Println("Created MaintenanceWindow:", maintenance) // {ID Description}
```

Get details for a specific maintenance:

```go
maintenance, err := client.Maintenances.Read(12345)
```

Update a maintenance: (Please note, that based on experience, you are allowed to modify only `Description`, `EffectiveTo` and `To`)

```go
updatedMaintenance := pingdom.MaintenanceWindow{
    Description: "My Maintenance",
    To:          1234567999,
}
msg, err := client.Maintenances.Update(12345, &updatedMaintenance)
```

Delete a maintenance:

Note: that only future maintenance window can be deleted. This means that both `To` and `From` should be in future.

```go
msg, err := client.Maintenances.Delete(12345)
```

After contacting Pingdom, the better approach would be to use update function and setting `To` and `EffectiveTo` to current time

```go
maintenance, _ := client.Maintenances.Read(12345)

m := pingdom.MaintenanceWindow{
    Description: maintenance.Description,
    From:        maintenance.From,
    To:          1,
    EffectiveTo: 1,
}

maintenanceUpdate, err := client.Maintenances.Update(12345, &m)
```

### ProbeService ###

This service gets pingdom Probes which are represented by the `Probes` struct.

More information on Probes from Pingdom: https://www.pingdom.com/resources/api/2.1#ResourceProbes
Several parameters are supported for filtering output. Please see them in Pingdom API documentation.

**NOTE:** Official documentation does not specify that `region` is returned for every probe entry, but it does and you can use it.

Get a list of all probes:

```go
params := make(map[string]string)

probes, err := client.Probes.List(params)
fmt.Println("Probes:", probes) // [{ID Name} ...]

for _, probe := range probes {
    fmt.Println("Probe region:", probe.Region)  // Probe region: EU
}
```

### TeamService ###

This service manages pingdom Teams which are represented by the `Team` struct.
When creating or updating Teams you must specify the `Name` and `MemberIDs`,
though `MemberIDs` may be an empty slice.
More information on Teams from Pingdom: https://docs.pingdom.com/api/#tag/Teams

Get a list of all teams:

```go
teams, err := client.Teams.List()
fmt.Println("Teams:", teams) // [{ID Name MemberIDs} ...]
```

Create a new Team:

```go
t := pingdom.TeamData{
    Name: "Team",
    MemberIDs: []int{},
}
team, err := client.Teams.Create(&t)
fmt.Println("Created Team:", team) // {ID Name MemberIDs}
```

Get details for a specific team:

```go
team, err := client.Teams.Read(12345)
```

Update a team:

```go
modifyTeam := pingdom.TeamData{
    Name:    "New Name"
    MemberIDs: []int{123, 678},
}
team, err := client.Teams.Update(12345, &modifyTeam)
```

Delete a team:

```go
team, err := client.Teams.Delete(12345)
```

### ContactService ###

This service manages users and their contact information which is represented by the `Contact` struct.
More information from Pingdom: https://docs.pingdom.com/api/#tag/Contacts

Get all contact info:

```go
contacts, err := client.Contacts.List()
fmt.Println(contacts)
```

Create a new contact:

```go
contact := Contact{
    Name: "John Doe",
    Paused: false,
    NotificationTargets: NotificationTargets{
        SMS: []SMSNotificationTarget{
            {
                Number: "5555555555",
                CountryCode: "1",
                Provider: "Verizon",
            }
        }
    }
}
contactId, err := client.Contacts.Create(contact)
fmt.Println("New Contact ID: ", contactId.Id)
```

Update a contact

```go
contactId := 1234

contact := Contact{
    Name : "John Doe",
    Paused : false,
    NotificationTargets: NotificationTargets{
        SMS: []SMSNotificationTarget{
            {
                Number: "5555555555",
                CountryCode: "1",
                Provider: "T-Mobile",
            }
        }
    }
}
result, err := client.Contacts.Update(contactId, contact)
fmt.Println(result.Message)
```

Delete a contact

```go
contactId := 1234

result, err := client.Contacts.Delete(contactId)
fmt.Println(result.Message)
```


### IntegrationService ###

This service manages pingdom Integrations which are represented by the `Integration` struct. Now only support manages the WebHook Integrations.
When creating or updating Integrations you must specify the `Active`, `ProviderID` and `WebHookData`.  


Get a list of all integrations:

```go
integrations, err := client_ext.Integrations.List()
fmt.Println("Integrations:", integrations) 
```

Create a new WebHook Integration:

```go
newIntegration := pingdomext.WebHookIntegration{
	Active:     false,
	ProviderID: 2,
	UserData: &pingdomext.WebHookData{
		Name: "tets-1",
		URL:  "http://www.example.com",
	},
}
integrationStatus, err := client_ext.Integrations.Create(&newIntegration)
fmt.Println("Created integration:", integrationStatus) 
```

Get details for a specific integration:

```go
integrationDetail, err := client_ext.Integrations.Read(12345)
```


Update a integration:

```go
updatedIntegration := pingdomext.WebHookIntegration{
	Active:     true,
	ProviderID: 2,
	UserData: &pingdomext.WebHookData{
		Name: "tets-3",
		URL:  "http://www.example5.com",
	},
}
updateMsg, err := client_ext.Integrations.Update(12345, &updatedIntegration)
```

Delete a integration:

```go
delMsg, err := client_ext.Integrations.Delete(12345)
```

List all integration providers:

```go
listProviders, err := client_ext.Integrations.ListProviders()
```



## Development ##

### Acceptance Tests ###

You can run acceptance tests against the actual pingdom API to test any changes:
```
PINGDOM_API_TOKEN=[api token] make acceptance
```

In order to run acceptance tests against the pingdom extension API, the following environment variables must be set:
```
SOLARWINDS_USER=[username] SOLARWINDS_PASSWD=[password] make acceptance
```

Note that this will create actual resources in your Pingdom account.  The tests will make a best effort to clean up but these would
not be guaranteed on test failures depending on the nature of the failure.
