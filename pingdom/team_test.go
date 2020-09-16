package pingdom

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMTeamServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alerting/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"teams": [
			  	{
					"id": 1,
					"name": "Team Rocket",
					"members": [
						{
							"id": 1,
							"name": "John Doe",
							"type": "user"
						}
					]
				},
				{
					"id": 2,
					"name": "The A-Team",
					"members": [
						{
							"id": 2,
							"name": "John 'Hannibal' Smith",
							"type": "user"
						},
						{
							"id": 3,
							"name": "Templeton 'Faceman' Peck",
							"type": "contact"
						}
					]
				}
			]
		  }`,
		)
	})
	want := []TeamResponse{
		{
			ID:   1,
			Name: "Team Rocket",
			Members: []TeamMemberResponse{
				{
					ID:   1,
					Name: "John Doe",
					Type: "user",
				},
			},
		},
		{
			ID:   2,
			Name: "The A-Team",
			Members: []TeamMemberResponse{
				{
					ID:   2,
					Name: "John 'Hannibal' Smith",
					Type: "user",
				},
				{
					ID:   3,
					Name: "Templeton 'Faceman' Peck",
					Type: "contact",
				},
			},
		},
	}

	teams, err := client.Teams.List()
	assert.NoError(t, err)
	assert.Equal(t, want, teams, "Teams.List() should return correct result")
}

func TestTeamServiceCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alerting/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"team": {
			  	"id": 12345678
			}
		  }`)
	})

	team := Team{
		Name:      "Operations",
		MemberIDs: []int{12345, 54321},
	}

	want := &TeamResponse{
		ID: 12345678,
	}

	teams, err := client.Teams.Create(&team)
	assert.NoError(t, err)
	assert.Equal(t, want, teams, "Teams.Create() should return correct result")
}

func TestTeamServiceRead(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alerting/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"team": {
				"id": 1,
				"name": "Team Rocket",
				"members": [
					{
						"id": 1,
						"name": "John Doe",
						"type": "user"
					},
					{
						"id": 4,
						"name": "Sidekick Jimmy",
						"type": "contact"
					}
				]
			}
		}`)
	})

	want := &TeamResponse{
		ID:   1,
		Name: "Team Rocket",
		Members: []TeamMemberResponse{
			{
				ID:   1,
				Name: "John Doe",
				Type: "user",
			},
			{
				ID:   4,
				Name: "Sidekick Jimmy",
				Type: "contact",
			},
		},
	}

	team, err := client.Teams.Read(1)
	assert.NoError(t, err)
	assert.Equal(t, want, team, "Teams.Read() should return correct result")
}

func TestTeamServiceUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alerting/teams/65", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{
			"team": {
				"id": 65,
				"name": "Operations",
				"members": [
				{
					"id": 10034512,
					"name": "John \"Hannibal\" Smith",
					"type": "contact"
				},
				{
					"id": 10043154,
					"name": "Templeton \"Face(man)\" Peck",
					"type": "contact"
				}
				]
			}
		}`)
	})

	updateTeam := Team{
		Name:      "Operations",
		MemberIDs: []int{10034512, 10043154},
	}

	want := &TeamResponse{
		ID:   65,
		Name: "Operations",
		Members: []TeamMemberResponse{
			{
				ID:   10034512,
				Name: "John \"Hannibal\" Smith",
				Type: "contact",
			},
			{
				ID:   10043154,
				Name: "Templeton \"Face(man)\" Peck",
				Type: "contact",
			},
		},
	}

	team, err := client.Teams.Update(65, &updateTeam)
	assert.NoError(t, err)
	assert.Equal(t, want, team, "Teams.Update() should return correct result")
}

func TestTeamServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/alerting/teams/1234", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{
			"message": "Deletion of team 1234 was successful"
	}`)
	})
	want := &TeamDeleteResponse{Message: "Deletion of team 1234 was successful"}

	team, err := client.Teams.Delete(1234)
	assert.NoError(t, err)
	assert.Equal(t, want, team, "Teams.Delete() should return correct result")
}
