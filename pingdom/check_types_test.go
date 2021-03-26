package pingdom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpCheckPutParams(t *testing.T) {
	verifyCertificate := true
	sslDownDaysBefore := 10

	tests := []struct {
		name       string
		giveCheck  HttpCheck
		wantParams map[string]string
	}{
		{
			name: "parametrizes http check",
			giveCheck: HttpCheck{
				Name:     "fake check",
				Hostname: "example.com",
				Url:      "/foo",
				RequestHeaders: map[string]string{
					"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
					"Pragma":     "no-cache",
				},
				Username:                 "user",
				Password:                 "pass",
				IntegrationIds:           []int{33333333, 44444444},
				UserIds:                  []int{123, 456},
				TeamIds:                  []int{789},
				ResponseTimeThreshold:    2300,
				Resolution:               5,
				VerifyCertificate:        &verifyCertificate,
				SSLDownDaysBefore:        &sslDownDaysBefore,
				SendNotificationWhenDown: 3,
			},
			wantParams: map[string]string{
				"name":                     "fake check",
				"host":                     "example.com",
				"paused":                   "false",
				"resolution":               "5",
				"notifyagainevery":         "0",
				"notifywhenbackup":         "false",
				"url":                      "/foo",
				"requestheader0":           "Pragma:no-cache",
				"requestheader1":           "User-Agent:Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
				"auth":                     "user:pass",
				"encryption":               "false",
				"shouldnotcontain":         "",
				"postdata":                 "",
				"integrationids":           "33333333,44444444",
				"tags":                     "",
				"probe_filters":            "",
				"userids":                  "123,456",
				"teamids":                  "789",
				"responsetime_threshold":   "2300",
				"verify_certificate":       "true",
				"ssl_down_days_before":     "10",
				"sendnotificationwhendown": "3",
			},
		},
		{
			name: "parametrizes http check without optional fields",
			giveCheck: HttpCheck{
				Name:     "fake check",
				Hostname: "example.com",
				Url:      "/foo",
				RequestHeaders: map[string]string{
					"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
					"Pragma":     "no-cache",
				},
				Username:              "user",
				Password:              "pass",
				IntegrationIds:        []int{33333333, 44444444},
				UserIds:               []int{123, 456},
				TeamIds:               []int{789},
				ResponseTimeThreshold: 2300,
			},
			wantParams: map[string]string{
				"name":                   "fake check",
				"host":                   "example.com",
				"paused":                 "false",
				"notifyagainevery":       "0",
				"notifywhenbackup":       "false",
				"url":                    "/foo",
				"requestheader0":         "Pragma:no-cache",
				"requestheader1":         "User-Agent:Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
				"auth":                   "user:pass",
				"encryption":             "false",
				"shouldnotcontain":       "",
				"postdata":               "",
				"integrationids":         "33333333,44444444",
				"tags":                   "",
				"probe_filters":          "",
				"userids":                "123,456",
				"teamids":                "789",
				"responsetime_threshold": "2300",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			params := tt.giveCheck.PutParams()
			assert.Equal(tst, tt.wantParams, params)
		})
	}
}

func TestHttpCheckPostParams(t *testing.T) {
	verifyCertificate := true
	sslDownDaysBefore := 10

	check := HttpCheck{
		Name:     "fake check",
		Hostname: "example.com",
		Url:      "/foo",
		RequestHeaders: map[string]string{
			"User-Agent": "Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
			"Pragma":     "no-cache",
		},
		Username:              "user",
		Password:              "pass",
		IntegrationIds:        []int{33333333, 44444444},
		UserIds:               []int{123, 456},
		TeamIds:               []int{789},
		ResponseTimeThreshold: 2300,
		VerifyCertificate:     &verifyCertificate,
		SSLDownDaysBefore:     &sslDownDaysBefore,
	}
	want := map[string]string{
		"name":                   "fake check",
		"host":                   "example.com",
		"paused":                 "false",
		"notifyagainevery":       "0",
		"notifywhenbackup":       "false",
		"type":                   "http",
		"url":                    "/foo",
		"requestheader0":         "Pragma:no-cache",
		"requestheader1":         "User-Agent:Pingdom.com_bot_version_1.4_(http://www.pingdom.com/)",
		"auth":                   "user:pass",
		"encryption":             "false",
		"integrationids":         "33333333,44444444",
		"userids":                "123,456",
		"teamids":                "789",
		"responsetime_threshold": "2300",
		"verify_certificate":     "true",
		"ssl_down_days_before":   "10",
	}

	params := check.PostParams()
	assert.Equal(t, want, params)
}

func TestHttpCheckValid(t *testing.T) {
	check := HttpCheck{Name: "fake check", Hostname: "example.com"}
	assert.NoError(t, check.Valid())

	badCheck := HttpCheck{Name: "fake check", Hostname: "example.com", Resolution: 6}
	assert.Error(t, badCheck.Valid())

	badContainsCheck := HttpCheck{
		Name:             "fake check",
		Hostname:         "example.com",
		Resolution:       15,
		ShouldContain:    "foo",
		ShouldNotContain: "bar",
	}
	assert.Error(t, badContainsCheck.Valid())
}

func TestPingCheckPostParams(t *testing.T) {
	check := PingCheck{
		Name:                  "fake check",
		Hostname:              "example.com",
		IntegrationIds:        []int{33333333, 44444444},
		UserIds:               []int{123, 456},
		TeamIds:               []int{789},
		ResponseTimeThreshold: 2300,
	}
	want := map[string]string{
		"name":                   "fake check",
		"host":                   "example.com",
		"paused":                 "false",
		"notifyagainevery":       "0",
		"notifywhenbackup":       "false",
		"type":                   "ping",
		"integrationids":         "33333333,44444444",
		"userids":                "123,456",
		"teamids":                "789",
		"responsetime_threshold": "2300",
	}

	params := check.PostParams()
	assert.Equal(t, want, params)
}

func TestPingCheckPutParams(t *testing.T) {
	check := PingCheck{
		Name:           "fake check",
		Hostname:       "example.com",
		IntegrationIds: []int{33333333, 44444444},
		UserIds:        []int{123, 456},
		TeamIds:        []int{789},
		Resolution:     5,
	}
	want := map[string]string{
		"name":             "fake check",
		"host":             "example.com",
		"resolution":       "5",
		"paused":           "false",
		"notifyagainevery": "0",
		"notifywhenbackup": "false",
		"integrationids":   "33333333,44444444",
		"probe_filters":    "",
		"userids":          "123,456",
		"teamids":          "789",
	}

	params := check.PutParams()
	assert.Equal(t, want, params)
}

func TestPingCheckValid(t *testing.T) {
	check := PingCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	assert.NoError(t, check.Valid())

	badCheck := PingCheck{Name: "fake check", Resolution: 10}
	assert.Error(t, badCheck.Valid())
}

func TestTCPCheckPostParams(t *testing.T) {
	check := TCPCheck{
		Name:           "fake check",
		Hostname:       "example.com",
		IntegrationIds: []int{33333333, 44444444},
		UserIds:        []int{123, 456},
		TeamIds:        []int{789},
		Port:           8080,
		StringToSend:   "Hello World",
		StringToExpect: "Hi there",
	}
	want := map[string]string{
		"name":             "fake check",
		"host":             "example.com",
		"paused":           "false",
		"notifyagainevery": "0",
		"notifywhenbackup": "false",
		"type":             "tcp",
		"integrationids":   "33333333,44444444",
		"userids":          "123,456",
		"teamids":          "789",
		"port":             "8080",
		"stringtosend":     "Hello World",
		"stringtoexpect":   "Hi there",
	}

	params := check.PostParams()
	assert.Equal(t, want, params)
}

func TestTCPCheckValid(t *testing.T) {
	check := TCPCheck{Name: "fake check", Hostname: "example.com", Resolution: 15, Port: 8080}
	assert.NoError(t, check.Valid())

	badCheck := TCPCheck{Name: "fake check", Hostname: "example.com", Resolution: 15}
	assert.Error(t, badCheck.Valid())

	badPortCheck := TCPCheck{Name: "fake check", Hostname: "example.com", Port: 66666}
	assert.Error(t, badPortCheck.Valid())
}

func TestDNSCheckPutParams(t *testing.T) {
	tests := []struct {
		name       string
		giveCheck  DNSCheck
		wantParams map[string]string
	}{
		{
			name: "parametrizes DNS check with all fields",
			giveCheck: DNSCheck{
				Name:                     "fake check",
				Hostname:                 "example.com",
				ExpectedIP:               "192.168.1.1",
				NameServer:               "8.8.8.8",
				IntegrationIds:           []int{33333333, 66666666},
				Resolution:               10,
				Paused:                   false,
				SendNotificationWhenDown: 3,
				NotifyAgainEvery:         5,
				NotifyWhenBackup:         false,
				Tags:                     "abc,efg,xyz",
				ProbeFilters:             "region: NA",
				UserIds:                  []int{123, 456},
				TeamIds:                  []int{789},
			},
			wantParams: map[string]string{
				"name":                     "fake check",
				"host":                     "example.com",
				"expectedip":               "192.168.1.1",
				"nameserver":               "8.8.8.8",
				"paused":                   "false",
				"resolution":               "10",
				"notifyagainevery":         "5",
				"notifywhenbackup":         "false",
				"integrationids":           "33333333,66666666",
				"tags":                     "abc,efg,xyz",
				"probe_filters":            "region: NA",
				"userids":                  "123,456",
				"teamids":                  "789",
				"sendnotificationwhendown": "3",
			},
		},
		{
			name: "parametrizes http check without optional fields",
			giveCheck: DNSCheck{
				Name:       "fake check",
				Hostname:   "example.com",
				ExpectedIP: "192.168.1.1",
				NameServer: "8.8.8.8",
			},
			wantParams: map[string]string{
				"name":             "fake check",
				"host":             "example.com",
				"expectedip":       "192.168.1.1",
				"nameserver":       "8.8.8.8",
				"paused":           "false",
				"notifyagainevery": "0",
				"notifywhenbackup": "false",
				"integrationids":   "",
				"tags":             "",
				"probe_filters":    "",
				"userids":          "",
				"teamids":          "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(tst *testing.T) {
			params := tt.giveCheck.PutParams()
			assert.Equal(tst, tt.wantParams, params)
		})
	}
}

func TestDNSCheckPostParams(t *testing.T) {
	check := DNSCheck{
		Name:           "fake check",
		Hostname:       "example.com",
		ExpectedIP:     "192.168.1.1",
		NameServer:     "8.8.8.8",
		IntegrationIds: []int{33333333, 44444444},
		UserIds:        []int{123, 456},
		TeamIds:        []int{789},
	}
	want := map[string]string{
		"name":             "fake check",
		"host":             "example.com",
		"expectedip":       "192.168.1.1",
		"nameserver":       "8.8.8.8",
		"paused":           "false",
		"notifyagainevery": "0",
		"notifywhenbackup": "false",
		"type":             "dns",
		"integrationids":   "33333333,44444444",
		"userids":          "123,456",
		"teamids":          "789",
	}

	params := check.PostParams()
	assert.Equal(t, want, params)
}

func TestDNSCheckValid(t *testing.T) {
	check := DNSCheck{Name: "fake check", Hostname: "example.com", ExpectedIP: "192.168.0.1", NameServer: "8.8.8.8", Resolution: 15}
	assert.NoError(t, check.Valid())

	badCheck := DNSCheck{Name: "fake check", Hostname: "example.com"}
	assert.Error(t, badCheck.Valid())

	badNameServerCheck := DNSCheck{Name: "fake check", Hostname: "example.com", ExpectedIP: "192.168.0.1"}
	assert.Error(t, badNameServerCheck.Valid())
}

func TestValidCommonParameters(t *testing.T) {
	assert.Error(t, validCommonParameters("", "example.com", 5))
	assert.Error(t, validCommonParameters("Test Name", "", 5))
	assert.Error(t, validCommonParameters("Test Name", "example.com", 7))
	assert.NoError(t, validCommonParameters("Test Name", "example.com", 0))
}

func TestSummaryPerformanceRequestValid(t *testing.T) {
	t.Run("missing field 'id'", func(t *testing.T) {
		assert.Equal(t, ErrMissingId, SummaryPerformanceRequest{}.Valid())
	})

	t.Run("resolution", func(t *testing.T) {
		assert.Nil(t, SummaryPerformanceRequest{
			Id:         123,
			Resolution: "hour",
		}.Valid())
		assert.Nil(t, SummaryPerformanceRequest{
			Id:         123,
			Resolution: "day",
		}.Valid())
		assert.Nil(t, SummaryPerformanceRequest{
			Id:         123,
			Resolution: "week",
		}.Valid())
		assert.Equal(t, ErrBadResolution, SummaryPerformanceRequest{
			Id:         123,
			Resolution: "month",
		}.Valid())

	})
}

func TestSummaryPerformanceRequestGetParams(t *testing.T) {
	id := 1337
	t.Run("empty request", func(t *testing.T) {
		want := map[string]string{}

		params := SummaryPerformanceRequest{
			Id: id,
		}.GetParams()

		assert.Equal(t, want, params)
	})

	t.Run("with some params", func(t *testing.T) {
		want := map[string]string{
			"resolution":    "week",
			"includeuptime": "true",
		}

		params := SummaryPerformanceRequest{
			Id:            id,
			IncludeUptime: true,
			Resolution:    "week",
		}.GetParams()

		assert.Equal(t, want, params)
	})
}
