package pingdom

import (
	"reflect"
	"testing"
)

func TestHttpCheckToParamsMinimal(t *testing.T) {
	check := HttpCheck{
		BaseCheck: BaseCheck{
			Name: "fake check",
			Host: "example.com",
		},
	}
	params := check.toParams()
	want := map[string]string{
		"type": "http",
		"name": "fake check",
		"host": "example.com",
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Check.toParams() returned %+v, want %+v", params, want)
	}
}

func TestHttpCheckToParamsAll(t *testing.T) {
	check := HttpCheck{
		BaseCheck: BaseCheck{
			Name:                     "fake check",
			Host:                     "example.com",
			Paused:                   OptBool(true),
			Resolution:               OptInt(5),
			SendNotificationWhenDown: OptInt(2),
			NotifyAgainEvery:         OptInt(5),
			NotifyWhenBackup:         OptBool(false),
			Tags:                     OptStr("tag1"),
			ProbeFilters:             OptStr("probefilter"),
			IPv6:                     OptBool(false),
			ResponseTimeThreshold: OptInt(30),
			UserIds:               &[]int{11111111, 22222222},
			IntegrationIds:        &[]int{33333333, 44444444},
			TeamIds:               &[]int{55555555, 66666666},
		},
		Url:        OptStr("/foo"),
		Encryption: OptBool(true),
		RequestHeaders: map[string]string{
			"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
			"Pragma":     "no-cache",
		},
		Username: OptStr("user"),
		Password: OptStr("pass"),
	}
	params := check.toParams()
	want := map[string]string{
		"type":                     "http",
		"name":                     "fake check",
		"host":                     "example.com",
		"url":                      "/foo",
		"resolution":               "5",
		"paused":                   "true",
		"sendnotificationwhendown": "2",
		"notifyagainevery":         "5",
		"notifywhenbackup":         "false",
		"ipv6":                     "false",
		"responsetime_threshold": "30",
		"requestheader0":         "Pragma:no-cache",
		"requestheader1":         "User-Agent:Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
		"auth":                   "user:pass",
		"encryption":             "true",
		"tags":                   "tag1",
		"probe_filters":          "probefilter",
		"userids":                "11111111,22222222",
		"integrationids":         "33333333,44444444",
		"teamids":                "55555555,66666666",
	}

	for k, v := range params {
		if want[k] != v {
			t.Errorf("the fuck: %v: %v", k, want[k])
		}
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Check.toParams() returned %+v, want %+v", params, want)
	}
}

func TestHttpCheckValid(t *testing.T) {
	check := HttpCheck{BaseCheck: BaseCheck{Name: "fake check", Host: "example.com", Resolution: OptInt(15)}}
	if err := check.valid(); err != nil {
		t.Errorf("Valid with valid check returned error %+v", err)
	}

	check = HttpCheck{BaseCheck: BaseCheck{Name: "fake check", Host: "example.com", Resolution: OptInt(0)}}
	if err := check.valid(); err == nil {
		t.Errorf("Valid with invalid check (`Resolution` == 0) expected error, returned nil")
	}

	check = HttpCheck{
		BaseCheck: BaseCheck{
			Name:       "fake check",
			Host:       "example.com",
			Resolution: OptInt(15),
		},
		ShouldContain:    OptStr("foo"),
		ShouldNotContain: OptStr("bar"),
	}
	if err := check.valid(); err == nil {
		t.Errorf("Valid with invalid check (`ShouldContain` and `ShouldNotContain` defined) expected error, returned nil")
	}

}

func TestPingCheckToParams(t *testing.T) {
	check := PingCheck{BaseCheck: BaseCheck{Name: "fake check", Host: "example.com", UserIds: &[]int{11111111, 22222222}, IntegrationIds: &[]int{33333333, 44444444}}}
	params := check.toParams()
	want := map[string]string{
		"type":           "ping",
		"name":           "fake check",
		"host":           "example.com",
		"userids":        "11111111,22222222",
		"integrationids": "33333333,44444444",
	}

	if !reflect.DeepEqual(params, want) {
		t.Errorf("Check.PostParams() returned %+v, want %+v", params, want)
	}
}

func TestPingCheckValid(t *testing.T) {
	check := PingCheck{BaseCheck{Name: "fake check", Host: "example.com", Resolution: OptInt(15)}}
	if err := check.valid(); err != nil {
		t.Errorf("Valid with valid check returned error %+v", err)
	}

	check = PingCheck{BaseCheck{Name: "fake check", Host: "example.com", Resolution: OptInt(23)}}
	if err := check.valid(); err == nil {
		t.Errorf("Valid with invalid check expected error, returned nil")
	}
}
