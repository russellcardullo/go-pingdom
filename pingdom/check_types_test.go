package pingdom

import (
	"reflect"
	"testing"
)

func TestHttpCheckParams(t *testing.T) {
	check := HttpCheck{Name: "fake check", Hostname: "example.com"}
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

func TestHttpCheckValid(t *testing.T) {
	check := HttpCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	if err := check.Valid(); err != nil {
		t.Errorf("Valid with valid check returned error %+v", err)
	}

	check = HttpCheck{Name: "fake check", Hostname: "example.com"}
	if err := check.Valid(); err == nil {
		t.Errorf("Valid with invalid check expected error, returned nil")
	}
}

func TestPingCheckParams(t *testing.T) {
	check := PingCheck{Name: "fake check", Hostname: "example.com"}
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
		"type":                     "ping",
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Check.Params() returned %+v, want %+v", params, want)
	}
}

func TestPingCheckValid(t *testing.T) {
	check := PingCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	if err := check.Valid(); err != nil {
		t.Errorf("Valid with valid check returned error %+v", err)
	}

	check = PingCheck{Name: "fake check", Hostname: "example.com"}
	if err := check.Valid(); err == nil {
		t.Errorf("Valid with invalid check expected error, returned nil")
	}
}
