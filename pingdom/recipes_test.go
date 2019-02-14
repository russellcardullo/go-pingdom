package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipiesList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms.recipes", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"recipes": {
				"76795": {
					"name": "Login/Logout",
					"status": "SUCCESSFUL",
					"kitchen": "us-east",
					"active": "YES",
					"created_at": 1550066917,
					"interval": 10,
					"use_legacy_notifications": false
				},
				"76791": {
					"name": "Login/Logout",
					"status": "SUCCESSFUL",
					"kitchen": "eu",
					"active": "YES",
					"created_at": 1550058211,
					"interval": 10,
					"use_legacy_notifications": false
				}
			},
			"counts": {
				"total": 2,
				"limited": 2,
				"filtered": 2
			}
		}`)
	})

	want := &listRecipesJsonResponse{
		Recipes: map[string]RecipeResponse{
			"76795": RecipeResponse{
				Name:                   "Login/Logout",
				Status:                 "SUCCESSFUL",
				Kitchen:                "us-east",
				Active:                 "YES",
				CreatedAt:              1550066917,
				Interval:               10,
				UseLegacyNotifications: false,
			},
			"76791": RecipeResponse{
				Name:                   "Login/Logout",
				Status:                 "SUCCESSFUL",
				Kitchen:                "eu",
				Active:                 "YES",
				CreatedAt:              1550058211,
				Interval:               10,
				UseLegacyNotifications: false,
			},
		},
	}

	checks, err := client.Recipe.List()
	assert.NoError(t, err)
	assert.Equal(t, want, checks)
}

func TestRecipeSummary(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/tms.summary.performance/76795", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"summary": {
				"hours": [
					{
						"#": [
							{
								"avgresponse": 1239,
								"command": "Go to URL https://app.signal-ai.com"
							},
							{
								"avgresponse": 159,
								"command": "Fill in field username with ops@signal-ai.com"
							},
							{
								"avgresponse": 174,
								"command": "Fill in field password with 4Hb6Ix^TJLcI4dELT$"
							},
							{
								"avgresponse": 1173,
								"command": "Click #kc-login"
							},
							{
								"avgresponse": 5799,
								"command": "Wait for element #sgStoryCard > div to exist"
							},
							{
								"avgresponse": 691,
								"command": "Click #account-actions"
							},
							{
								"avgresponse": 287,
								"command": "Click Logout"
							},
							{
								"avgresponse": 762,
								"command": "Click #yesButton"
							}
						],
						"avgresponse": 10289,
						"starttime": 1550124000,
						"revision": 696516
					},
					{
						"#": [
							{
								"avgresponse": 1232,
								"command": "Go to URL https://app.signal-ai.com"
							},
							{
								"avgresponse": 168,
								"command": "Fill in field username with ops@signal-ai.com"
							},
							{
								"avgresponse": 157,
								"command": "Fill in field password with 4Hb6Ix^TJLcI4dELT$"
							},
							{
								"avgresponse": 1160,
								"command": "Click #kc-login"
							},
							{
								"avgresponse": 5784,
								"command": "Wait for element #sgStoryCard > div to exist"
							},
							{
								"avgresponse": 672,
								"command": "Click #account-actions"
							},
							{
								"avgresponse": 262,
								"command": "Click Logout"
							},
							{
								"avgresponse": 772,
								"command": "Click #yesButton"
							}
						],
						"avgresponse": 10211,
						"starttime": 1550127600,
						"revision": 696516
					}
				]
			}
		}`)
	})

	want := &RecipeSummaryResponse{Summary: struct {
		Hours []struct {
			Data []struct {
				Avgresponse int    "json:\"avgresponse\""
				Command     string "json:\"command\""
			} "json:\"#\""
			Avgresponse int "json:\"avgresponse\""
			Starttime   int "json:\"starttime\""
			Revision    int "json:\"revision\""
		} "json:\"hours\""
	}{Hours: []struct {
		Data []struct {
			Avgresponse int    "json:\"avgresponse\""
			Command     string "json:\"command\""
		} "json:\"#\""
		Avgresponse int "json:\"avgresponse\""
		Starttime   int "json:\"starttime\""
		Revision    int "json:\"revision\""
	}{struct {
		Data []struct {
			Avgresponse int    "json:\"avgresponse\""
			Command     string "json:\"command\""
		} "json:\"#\""
		Avgresponse int "json:\"avgresponse\""
		Starttime   int "json:\"starttime\""
		Revision    int "json:\"revision\""
	}{Data: []struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 1239, Command: "Go to URL https://app.signal-ai.com"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 159, Command: "Fill in field username with ops@signal-ai.com"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 174, Command: "Fill in field password with 4Hb6Ix^TJLcI4dELT$"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 1173, Command: "Click #kc-login"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 5799, Command: "Wait for element #sgStoryCard > div to exist"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 691, Command: "Click #account-actions"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 287, Command: "Click Logout"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 762, Command: "Click #yesButton"}}, Avgresponse: 10289, Starttime: 1550124000, Revision: 696516}, struct {
		Data []struct {
			Avgresponse int    "json:\"avgresponse\""
			Command     string "json:\"command\""
		} "json:\"#\""
		Avgresponse int "json:\"avgresponse\""
		Starttime   int "json:\"starttime\""
		Revision    int "json:\"revision\""
	}{Data: []struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 1232, Command: "Go to URL https://app.signal-ai.com"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 168, Command: "Fill in field username with ops@signal-ai.com"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 157, Command: "Fill in field password with 4Hb6Ix^TJLcI4dELT$"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 1160, Command: "Click #kc-login"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 5784, Command: "Wait for element #sgStoryCard > div to exist"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 672, Command: "Click #account-actions"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 262, Command: "Click Logout"}, struct {
		Avgresponse int    "json:\"avgresponse\""
		Command     string "json:\"command\""
	}{Avgresponse: 772, Command: "Click #yesButton"}}, Avgresponse: 10211, Starttime: 1550127600, Revision: 696516}}}}

	checks, err := client.Recipe.Summary(76795)
	assert.NoError(t, err)
	assert.Equal(t, want, checks)
}
