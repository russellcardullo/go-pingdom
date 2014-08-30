package pingdom

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	server *httptest.Server
)

func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// test client
	client = NewClient("fake_email@example.com", "12345", "my_api_key")
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if want != r.Method {
		t.Errorf("Request method = %v, want %v", r.Method, want)
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient("user", "password", "key")
	if c.client != http.DefaultClient {
		t.Errorf("NewClient client = %v, want http.DefaultClient", c.client)
	}

	if c.BaseURL.String() != defaultBaseURL {
		t.Errorf("NewClient BaseURL = %v, want %v", c.BaseURL.String(), defaultBaseURL)
	}
}

func TestListChecks(t *testing.T) {
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
					"type": "http"
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
					"type": "ping"
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
					"type": "http"
				}
			]
		}`)
	})

	checks, err := client.ListChecks()
	if err != nil {
		t.Errorf("ListChecks returned error: %v", err)
	}

	want := &Check{ID: 85975, Name: "My check 1"}
	if !reflect.DeepEqual(checks, want) {
		t.Errorf("ListChecks returned %+v, want %+v", checks, want)
	}
}

func TestCreateCheck(t *testing.T) {
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

	newCheck := HttpCheck{"My new HTTP check", "example.com"}
	check, err := client.CreateCheck(newCheck)
	if err != nil {
		t.Errorf("CreateCheck returned error: %v", err)
	}

	want := &Check{ID: 138631, Name: "My new HTTP check"}
	if !reflect.DeepEqual(check, want) {
		t.Errorf("CreateCheck returned %+v, want %+v", check, want)
	}
}

func TestReadCheck(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/checks/85975", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{
			"check" : {
				"id" : 85975,
				"name" : "My check 7",
				"resolution" : 1,
				"sendtoemail" : false,
				"sendtosms" : false,
				"sendtotwitter" : false,
				"sendtoiphone" : false,
				"sendnotificationwhendown" : 0,
				"notifyagainevery" : 0,
				"notifywhenbackup" : false,
				"created" : 1240394682,
				"type" : {
				  "http" : {
					"url" : "/",
					"port" : 80,
					"requestheaders" : {
					  "User-Agent" : "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)"
					}
				  }
				},
				"hostname" : "s7.mydomain.com",
				"status" : "up",
				"lasterrortime" : 1293143467,
				"lasttesttime" : 1294064823
			}
		}`)
	})

	check, err := client.ReadCheck("85975")
	if err != nil {
		t.Errorf("ReadCheck returned error: %v", err)
	}

	want := &Check{ID: 85975, Name: "My check 7"}
	if !reflect.DeepEqual(check, want) {
		t.Errorf("ReadCheck returned %+v, want %+v", check, want)
	}

}

func TestUpdateCheck(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/checks/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		fmt.Fprint(w, `{"message":"Modification of check was successful!"}`)
	})

	msg, err := client.UpdateCheck("12345", "updated_check", "example2.com")
	if err != nil {
		t.Errorf("UpdateCheck returned error: %v", err)
	}

	want := &PingdomResponse{Message: "Modification of check was successful!"}
	if !reflect.DeepEqual(msg, want) {
		t.Errorf("UpdateCheck returned %+v, want %+v", msg, want)
	}
}

func TestDeleteCheck(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/api/2.0/checks/12345", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		fmt.Fprint(w, `{"message":"Deletion of check was successful!"}`)
	})

	msg, err := client.DeleteCheck("12345")
	if err != nil {
		t.Errorf("DeleteCheck returned error: %v", err)
	}

	want := &PingdomResponse{Message: "Deletion of check was successful!"}
	if !reflect.DeepEqual(msg, want) {
		t.Errorf("DeleteCheck returned %+v, want %+v", msg, want)
	}
}
