package pingdom

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
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

func TestNewRequest(t *testing.T) {
	setup()
	defer teardown()

	req, err := client.NewRequest("GET", "/checks", nil)
	if err != nil {
		t.Errorf("NewRequest returned error: %v", err)
	}

	if req.Method != "GET" {
		t.Errorf("NewRequest Method returned %+v, want %+v", req.Method, "GET")
	}

	if req.URL.String() != client.BaseURL.String()+"/checks" {
		t.Errorf("NewRequest URL returned %+v, want %+v", req.URL.String(), client.BaseURL.String()+"/checks")
	}
}

func TestValidateResponse(t *testing.T) {
	valid := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("OK")),
	}

	if err := ValidateResponse(valid); err != nil {
		t.Errorf("ValidateResponse with valid response returned error %+v", err)
	}

	invalid := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusBadRequest,
		Body: ioutil.NopCloser(strings.NewReader(`{
			"error" : {
				"statuscode": 400,
				"statusdesc": "Bad Request",
				"errormessage": "This is an error"
			}
		}`)),
	}

	want := &PingdomError{400, "Bad Request", "This is an error"}
	if err := ValidateResponse(invalid); !reflect.DeepEqual(err, want) {
		t.Errorf("ValidateResponse with invalid response returned %+v, want %+v", err, want)
	}

}

func TestParams(t *testing.T) {
	check := Check{Name: "fake check", Hostname: "example.com"}
	params := check.Params()
	want := map[string]string{
		"name":                     "fake check",
		"host":                     "example.com",
		"paused":                   "false",
		"resolution":               "0",
		"sendtoemail":              "false",
		"sendtosms":                "false",
		"sendtotwitter":            "false",
		"sendtoiphone":             "false",
		"sendtoandroid":            "false",
		"sendnotificationwhendown": "0",
		"notifyagainevery":         "0",
		"notifywhenbackup":         "false",
		"type":                     "http",
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Check.Params() returned %+v, want %+v", params, want)
	}
}

func TestValidate(t *testing.T) {
	check := Check{Name: "fake check", Hostname: "example.com", Resolution: 15}
	if err := Validate(&check); err != nil {
		t.Errorf("Validate with valid check returned error %+v", err)
	}

	check = Check{Name: "fake check", Hostname: "example.com"}
	if err := Validate(&check); err == nil {
		t.Errorf("Validate with invalid check expected error, returned nil")
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

	want := []Check{
		Check{
			ID:               85975,
			Name:             "My check 1",
			LastErrorTime:    1297446423,
			LastResponseTime: 355,
			LastTestTime:     1300977363,
			Hostname:         "example.com",
			Resolution:       1,
			Status:           "up",
		},
		Check{
			ID:               161748,
			Name:             "My check 2",
			LastErrorTime:    1299194968,
			LastResponseTime: 1141,
			LastTestTime:     1300977268,
			Hostname:         "mydomain.com",
			Resolution:       5,
			Status:           "up",
		},
		Check{
			ID:               208655,
			Name:             "My check 3",
			LastErrorTime:    1300527997,
			LastResponseTime: 800,
			LastTestTime:     1300977337,
			Hostname:         "example.net",
			Resolution:       1,
			Status:           "down",
		},
	}

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

	newCheck := Check{Name: "My new HTTP check", Hostname: "example.com", Resolution: 5}
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

	check, err := client.ReadCheck(85975)
	if err != nil {
		t.Errorf("ReadCheck returned error: %v", err)
	}

	want := &Check{
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
	}
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

	updateCheck := Check{Name: "Updated Check", Hostname: "example2.com", Resolution: 5}
	msg, err := client.UpdateCheck(12345, updateCheck)
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

	msg, err := client.DeleteCheck(12345)
	if err != nil {
		t.Errorf("DeleteCheck returned error: %v", err)
	}

	want := &PingdomResponse{Message: "Deletion of check was successful!"}
	if !reflect.DeepEqual(msg, want) {
		t.Errorf("DeleteCheck returned %+v, want %+v", msg, want)
	}
}
