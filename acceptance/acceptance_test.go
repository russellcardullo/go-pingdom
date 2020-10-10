package acceptance

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/russellcardullo/go-pingdom/pingdom"
	"github.com/stretchr/testify/assert"
)

var client *pingdom.Client

var runAcceptance bool

func init() {
	if os.Getenv("PINGDOM_ACCEPTANCE") == "1" {
		runAcceptance = true

		config := pingdom.ClientConfig{
			APIToken: os.Getenv("PINGDOM_API_TOKEN"),
			HTTPClient: &http.Client{
				Timeout: time.Second * 10,
			},
		}
		client, _ = pingdom.NewClientWithConfig(config)
	}
}

func TestListChecks(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	checks, err := client.Checks.List()
	assert.NoError(t, err)
	assert.NotNil(t, checks)
}

func TestHTTPCheck(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	newCheck := pingdom.HttpCheck{
		Name:                     "Test Check",
		Hostname:                 "example.com",
		Resolution:               5,
		SendNotificationWhenDown: 100,
		Tags:                     "tag",
	}
	check, err := client.Checks.Create(&newCheck)
	assert.NoError(t, err)
	assert.NotNil(t, check)

	newCheck.Name = "Test Check 2"
	up, err := client.Checks.Update(check.ID, &newCheck)
	assert.NoError(t, err)
	assert.NotNil(t, up)

	resp, err := client.Checks.Read(check.ID)
	assert.NoError(t, err)
	assert.Equal(t, newCheck.Name, resp.Name)

	delMsg, err := client.Checks.Delete(check.ID)
	assert.NoError(t, err)
	assert.NotNil(t, delMsg)
}

func TestTagSupport(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	newCheck := pingdom.HttpCheck{
		Name:                     "Test Check",
		Hostname:                 "example.com",
		Resolution:               5,
		SendNotificationWhenDown: 100,
		Tags:                     "tag",
	}
	check, err := client.Checks.Create(&newCheck)
	assert.NoError(t, err)
	assert.NotNil(t, check)

	params := make(map[string]string)
	params["include_tags"] = "true"
	params["tags"] = "tag"

	checks, err := client.Checks.List(params)
	assert.NoError(t, err)
	assert.NotNil(t, checks)
	assert.Equal(t, 1, len(checks))

	delMsg, err := client.Checks.Delete(check.ID)
	assert.NoError(t, err)
	assert.NotNil(t, delMsg)
}

func TestProbes(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	params := make(map[string]string)

	probes, err := client.Probes.List(params)
	assert.NoError(t, err)
	assert.NotNil(t, probes)
	assert.NotEmpty(t, probes)
}

func TestContacts(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}

	contact := pingdom.Contact{
		Name:   "Test User",
		Paused: false,
		NotificationTargets: pingdom.NotificationTargets{
			SMS: []pingdom.SMSNotification{
				{
					CountryCode: "00",
					Number:      "5555555555",
					Provider:    "nexmo",
					Severity:    "LOW",
				},
				{
					CountryCode: "00",
					Number:      "5555555555",
					Provider:    "nexmo",
					Severity:    "HIGH",
				},
			},
		},
	}

	createMsg, err := client.Contacts.Create(&contact)
	assert.NoError(t, err)
	assert.NotNil(t, createMsg)
	assert.NotEmpty(t, createMsg)

	contact.ID = createMsg.ID

	listMsg, err := client.Contacts.List()
	assert.NoError(t, err)
	assert.NotNil(t, listMsg)
	assert.NotEmpty(t, listMsg)

	contact.NotificationTargets.SMS[0].Number = "2222222222"
	updateMsg, err := client.Contacts.Update(contact.ID, &contact)
	assert.NoError(t, err)
	assert.NotNil(t, updateMsg)

	delMsg, err := client.Contacts.Delete(contact.ID)
	assert.NoError(t, err)
	assert.NotNil(t, delMsg)
}

func TestTeams(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}

	team := pingdom.Team{
		Name:      "Test team",
		MemberIDs: []int{},
	}

	createMsg, err := client.Teams.Create(&team)
	assert.NoError(t, err)
	assert.NotNil(t, createMsg)
	assert.NotEmpty(t, createMsg)

	team.ID = createMsg.ID

	listMsg, err := client.Teams.List()
	assert.NoError(t, err)
	assert.NotNil(t, listMsg)
	assert.NotEmpty(t, listMsg)

	team.Name = "Test team renamed"
	updateMsg, err := client.Teams.Update(team.ID, &team)
	assert.NoError(t, err)
	assert.NotNil(t, updateMsg)
	assert.NotEmpty(t, updateMsg)

	delMsg, err := client.Teams.Delete(team.ID)
	assert.NoError(t, err)
	assert.NotNil(t, delMsg)
}

func TestTeamAndContactConnections(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}

	contact := pingdom.Contact{
		Name:   "Test User",
		Paused: false,
		NotificationTargets: pingdom.NotificationTargets{
			SMS: []pingdom.SMSNotification{
				{
					CountryCode: "00",
					Number:      "5555555555",
					Provider:    "nexmo",
					Severity:    "LOW",
				},
				{
					CountryCode: "00",
					Number:      "5555555555",
					Provider:    "nexmo",
					Severity:    "HIGH",
				},
			},
		},
	}

	team := pingdom.Team{
		Name:      "Test team",
		MemberIDs: []int{},
	}

	createTeamMsg, err := client.Teams.Create(&team)
	assert.NoError(t, err)
	assert.NotNil(t, createTeamMsg)
	assert.NotEmpty(t, createTeamMsg)

	team.ID = createTeamMsg.ID

	createContactMsg, err := client.Contacts.Create(&contact)
	assert.NoError(t, err)
	assert.NotNil(t, createContactMsg)
	assert.NotEmpty(t, createContactMsg)

	contact.ID = createContactMsg.ID

	team.MemberIDs = append(team.MemberIDs, contact.ID)

	// Verify we can add contacts
	updateMsg, err := client.Teams.Update(team.ID, &team)
	assert.NoError(t, err)
	assert.NotNil(t, updateMsg)
	assert.NotEmpty(t, updateMsg)
	assert.NotEmpty(t, updateMsg.Members)

	team.MemberIDs = []int{}
	// Verify we can remove contacts
	updateMsg, err = client.Teams.Update(team.ID, &team)
	assert.NoError(t, err)
	assert.NotNil(t, updateMsg)
	assert.NotEmpty(t, updateMsg)
	assert.Empty(t, updateMsg.Members)

	delMsg, err := client.Teams.Delete(team.ID)
	assert.NoError(t, err)
	assert.NotNil(t, delMsg)
}
