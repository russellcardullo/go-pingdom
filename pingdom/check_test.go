package pingdom

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckServiceList(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/checks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"checks": [
				{
					"hostname": "example.com",
					"id": 85975,
					"lasterrortime": 1297446423,
					"lastresponsetime": 355,
					"lasttesttime": 1300977363,
					"name": "My check 1",
					"resolution": 1,
					"status": "up",
					"type": "http",
					"tags": [
						{
							"name": "apache",
							"type": "a",
							"count": 2
						}
					]
				},
				{
					"hostname": "mydomain.com",
					"id": 161748,
					"lasterrortime": 1299194968,
					"lastresponsetime": 1141,
					"lasttesttime": 1300977268,
					"name": "My check 2",
					"resolution": 5,
					"status": "up",
					"type": "ping",
					"tags": [
						{
							"name": "nginx",
							"type": "u",
							"count": 1
						}
					]
				},
				{
					"hostname": "example.net",
					"id": 208655,
					"lasterrortime": 1300527997,
					"lastresponsetime": 800,
					"lasttesttime": 1300977337,
					"name": "My check 3",
					"resolution": 1,
					"status": "down",
					"type": "http",
					"tags": [
						{
							"name": "apache",
							"type": "a",
							"count": 2
						}
					]
				}
			]
		}`)
	})

	var countA, countB float64 = 1, 2

	want := []CheckResponse{
		{
			ID:               85975,
			Name:             "My check 1",
			LastErrorTime:    1297446423,
			LastResponseTime: 355,
			LastTestTime:     1300977363,
			Hostname:         "example.com",
			Resolution:       1,
			Status:           "up",
			Type: CheckResponseType{
				Name: "http",
			},
			Tags: []CheckResponseTag{
				{
					Name:  "apache",
					Type:  "a",
					Count: countB,
				},
			},
		},
		{
			ID:               161748,
			Name:             "My check 2",
			LastErrorTime:    1299194968,
			LastResponseTime: 1141,
			LastTestTime:     1300977268,
			Hostname:         "mydomain.com",
			Resolution:       5,
			Status:           "up",
			Type: CheckResponseType{
				Name: "ping",
			},
			Tags: []CheckResponseTag{
				{
					Name:  "nginx",
					Type:  "u",
					Count: countA,
				},
			},
		},
		{
			ID:               208655,
			Name:             "My check 3",
			LastErrorTime:    1300527997,
			LastResponseTime: 800,
			LastTestTime:     1300977337,
			Hostname:         "example.net",
			Resolution:       1,
			Status:           "down",
			Type: CheckResponseType{
				Name: "http",
			},
			Tags: []CheckResponseTag{
				{
					Name:  "apache",
					Type:  "a",
					Count: countB,
				},
			},
		},
	}

	checks, err := client.Checks.List()
	if assert.Nil(t, err) {

		// remove tags to check separately
		assert.Equal(t, want, checks, "Checks.List() should return correct result")
	}
}

func TestCheckServiceCreate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/checks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprint(w, `{
			"check":{
				"id":138631,
				"name":"My new HTTP check"
			}
		}`)
	})

	newCheck := HttpCheck{
		Name:           "My new HTTP check",
		Hostname:       "example.com",
		Resolution:     5,
		ContactIds:     []int{11111111, 22222222},
		IntegrationIds: []int{33333333, 44444444},
	}
	check, err := client.Checks.Create(&newCheck)
	if err != nil {
		t.Errorf("CreateCheck returned error: %v", err)
	}

	want := &CheckResponse{ID: 138631, Name: "My new HTTP check"}
	if !reflect.DeepEqual(check, want) {
		t.Errorf("CreateCheck returned %+v, want %+v", check, want)
	}
}

func TestCheckServiceRead(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/checks/85975", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"check" : {
        "contactids": [
            11111111,
            22222222
        ],
        "created" : 1240394682,
        "hostname" : "s7.mydomain.com",
        "id" : 85975,
        "integrationids": [
            33333333,
            44444444
        ],
        "ipv6": false,
        "lasterrortime" : 1293143467,
        "lasttesttime" : 1294064823,
        "name" : "My check 7",
        "notifyagainevery" : 0,
        "notifywhenbackup" : false,
        "probe_filters": [],
        "resolution" : 1,
        "sendnotificationwhendown" : 0,
        "sendtoandroid" : false,
        "sendtoemail" : false,
        "sendtoiphone" : false,
        "sendtosms" : false,
        "sendtotwitter" : false,
        "status" : "up",
        "tags": [],
        "type" : {
          "http" : {
            "encryption": false,
            "port" : 80,
            "requestheaders" : {
              "User-Agent" : "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)"
            },
            "url" : "/"
          }
        }
			}
		}`)
	})

	check, err := client.Checks.Read(85975)
	if err != nil {
		t.Errorf("ReadCheck returned error: %v", err)
	}

	want := &CheckResponse{
		ID:                       85975,
		Name:                     "My check 7",
		Resolution:               1,
		SendToEmail:              false,
		SendToTwitter:            false,
		SendToIPhone:             false,
		SendNotificationWhenDown: 0,
		NotifyAgainEvery:         0,
		NotifyWhenBackup:         false,
		Created:                  1240394682,
		Hostname:                 "s7.mydomain.com",
		Status:                   "up",
		LastErrorTime:            1293143467,
		LastTestTime:             1294064823,
		Type: CheckResponseType{
			Name: "http",
			HTTP: &CheckResponseHTTPDetails{
				Url:              "/",
				Encryption:       false,
				Port:             80,
				Username:         "",
				Password:         "",
				ShouldContain:    "",
				ShouldNotContain: "",
				PostData:         "",
				RequestHeaders: map[string]string{
					"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
				},
			},
		},
		ContactIds:     []int{11111111, 22222222},
		IntegrationIds: []int{33333333, 44444444},
		Tags:           []CheckResponseTag{},
	}

	if !reflect.DeepEqual(check, want) {
		t.Errorf("ReadCheck returned %+v, want %+v", check, want)
	}

}

func TestCheckServiceUpdate(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/checks/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"message":"Modification of check was successful!"}`)
	})

	updateCheck := HttpCheck{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
	msg, err := client.Checks.Update(12345, &updateCheck)
	if err != nil {
		t.Errorf("UpdateCheck returned error: %v", err)
	}

	want := &PingdomResponse{Message: "Modification of check was successful!"}
	if !reflect.DeepEqual(msg, want) {
		t.Errorf("UpdateCheck returned %+v, want %+v", msg, want)
	}
}

func TestCheckServiceDelete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/checks/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"message":"Deletion of check was successful!"}`)
	})

	msg, err := client.Checks.Delete(12345)
	if err != nil {
		t.Errorf("DeleteCheck returned error: %v", err)
	}

	want := &PingdomResponse{Message: "Deletion of check was successful!"}
	if !reflect.DeepEqual(msg, want) {
		t.Errorf("DeleteCheck returned %+v, want %+v", msg, want)
	}
}
