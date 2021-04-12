package acceptance

import (
	"github.com/nordcloud/go-pingdom/solarwinds"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/nordcloud/go-pingdom/pingdom"
	"github.com/stretchr/testify/assert"
)

var client *pingdom.Client

var runAcceptance bool

func init() {
	if os.Getenv("PINGDOM_ACCEPTANCE") == "1" {
		runAcceptance = true

		config := pingdom.ClientConfig{
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

func TestDNSCheck(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}
	newCheck := pingdom.DNSCheck{
		Name:                     "Test Check",
		Hostname:                 "example.com",
		ExpectedIP:               "2606:2800:220:1:248:1893:25c8:1946",
		NameServer:               "a.iana-servers.net",
		SendNotificationWhenDown: 100,
		Tags:                     "dns",
	}
	check, err := client.Checks.Create(&newCheck)
	assert.NoError(t, err)
	assert.NotNil(t, check)
	assert.Equal(t, check.Name, "Test Check")

	newCheck.Name = "Test Check 2"
	up, err := client.Checks.Update(check.ID, &newCheck)
	assert.NoError(t, err)
	assert.NotNil(t, up)

	resp, err := client.Checks.Read(check.ID)
	assert.NoError(t, err)
	assert.Equal(t, newCheck.Name, resp.Name)
	assert.Equal(t, resp.Resolution, 5)

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

func TestOccurrences(t *testing.T) {
	if !runAcceptance {
		t.Skip()
	}

	now := time.Now()
	from := now.Add(24 * time.Hour)
	to := from.Add(1 * time.Hour)

	// To create a new occurrence, we need to create a maintenance.
	maintenance := pingdom.MaintenanceWindow{
		Description:    "Acceptance test - " + solarwinds.RandString(10),
		From:           from.Unix(),
		To:             to.Unix(),
		RecurrenceType: "day",
		RepeatEvery:    1,
		EffectiveTo:    to.Add(3 * 24 * time.Hour).Unix(),
	}
	createMaintenanceMsg, err := client.Maintenances.Create(&maintenance)
	assert.NoError(t, err)
	assert.NotNil(t, createMaintenanceMsg)
	assert.NotEmpty(t, createMaintenanceMsg)

	occurrences, err := client.Occurrences.List(pingdom.ListOccurrenceQuery{
		MaintenanceId: int64(createMaintenanceMsg.ID),
	})
	assert.NoError(t, err)
	assert.True(t, len(occurrences) > 1)
	occurrence := occurrences[0]

	newTo := time.Unix(occurrence.To, 0).Add(1 * time.Hour).Unix()
	resp, err := client.Occurrences.Update(occurrence.Id, pingdom.Occurrence{
		From: occurrence.From,
		To:   newTo,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	afterUpdate, err := client.Occurrences.Read(occurrence.Id)
	assert.NoError(t, err)
	assert.Equal(t, newTo, afterUpdate.To)

	resp, err = client.Occurrences.Delete(occurrence.Id)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	resp, err = client.Maintenances.Delete(createMaintenanceMsg.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	idsToDelete := make([]int64, 0, len(occurrences)-1)
	for _, occ := range occurrences[1:] {
		idsToDelete = append(idsToDelete, occ.Id)
	}
	resp, err = client.Occurrences.MultiDelete(idsToDelete)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp)

	occurrences, err = client.Occurrences.List(pingdom.ListOccurrenceQuery{
		MaintenanceId: int64(createMaintenanceMsg.ID),
	})
	assert.NoError(t, err)
	assert.True(t, len(occurrences) == 0)
}
