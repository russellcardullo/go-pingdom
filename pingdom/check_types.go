package pingdom

import (
	"fmt"
	"sort"
	"strconv"
)


// BaseCheck represents the basic parameters of a Pingdom check. Optional
// parameters are pointers.
type BaseCheck struct {
	Name                     string
	Host                     string
	Type                     string
	Paused                   *bool
	Resolution               *int
	UserIds                  *[]int
	SendNotificationWhenDown *int
	NotifyAgainEvery         *int
	NotifyWhenBackup         *bool
	Tags                     *string
	ProbeFilters             *string
	IPv6                     *bool
	ResponseTimeThreshold    *int
	IntegrationIds           *[]int
	TeamIds                  *[]int
}

func (bc *BaseCheck) toParams() map[string]interface{} {
	return map[string]interface{}{
		"name":                     bc.Name,
		"host":                     bc.Host,
		"type":                     bc.Type,
		"paused":                   bc.Paused,
		"resolution":               bc.Resolution,
		"userids":                  bc.UserIds,
		"sendnotificationwhendown": bc.SendNotificationWhenDown,
		"notifyagainevery":         bc.NotifyAgainEvery,
		"notifywhenbackup":         bc.NotifyWhenBackup,
		"tags":                     bc.Tags,
		"probe_filters":            bc.ProbeFilters,
		"ipv6":                     bc.IPv6,
		"responsetime_threshold": bc.ResponseTimeThreshold,
		"integrationids":         bc.IntegrationIds,
		"teamids":                bc.TeamIds,
	}
}

func (bc *BaseCheck) valid() error {
	if bc.Name == "" {
		return fmt.Errorf("Invalid value for `Name`.  Must contain non-empty string")
	}

	if bc.Host == "" {
		return fmt.Errorf("Invalid value for `Hostname`.  Must contain non-empty string")
	}

	if bc.Resolution != nil && *bc.Resolution != 1 && *bc.Resolution != 5 && *bc.Resolution != 15 &&
		*bc.Resolution != 30 && *bc.Resolution != 60 {
		return fmt.Errorf("Invalid value %v for `Resolution`.  Allowed values are [1,5,15,30,60].", bc.Resolution)
	}

	return nil
}

// HttpCheck represents a Pingdom http check.
type HttpCheck struct {
	BaseCheck
	Url              *string
	Encryption       *bool
	Port             *int
	Username         *string
	Password         *string
	ShouldContain    *string
	ShouldNotContain *string
	PostData         *string
	RequestHeaders   map[string]string
}

func (hc *HttpCheck) toParams() map[string]string {
	m := hc.BaseCheck.toParams()
	m["type"] = "http"
	m["url"] = hc.Url
	m["encryption"] = hc.Encryption
	m["port"] = hc.Port

	// Convert auth
	if hc.Username != nil && hc.Password != nil {
		m["auth"] = fmt.Sprintf("%s:%s", *hc.Username, *hc.Password)
	}

	// ShouldContain and ShouldNotContain are mutually exclusive.
	// But we must define one so they can be emptied if required.
	if hc.ShouldContain != nil {
		m["shouldcontain"] = hc.ShouldContain
	} else {
		m["shouldnotcontain"] = hc.ShouldNotContain
	}

	m["postdata"] = hc.PostData

	// Convert headers
	var headers []string
	for k := range hc.RequestHeaders {
		headers = append(headers, k)
	}
	sort.Strings(headers)
	for i, k := range headers {
		m[fmt.Sprintf("requestheader%d", i)] = fmt.Sprintf("%s:%s", k, hc.RequestHeaders[k])
	}

	return valuesToStrings(m)
}

// Determine whether the HttpCheck contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (ck *HttpCheck) valid() error {
	if err := ck.BaseCheck.valid(); err != nil {
		return err
	}

	if ck.ShouldContain != nil && ck.ShouldNotContain != nil {
		return fmt.Errorf("`ShouldContain` and `ShouldNotContain` must not be declared at the same time")
	}

	return nil
}

// PingCheck represents a Pingdom ping check
type PingCheck struct {
	BaseCheck
}

func (pc *PingCheck) toParams() map[string]string {
	m := pc.BaseCheck.toParams()
	m["type"] = "ping"
	return valuesToStrings(m)
}

// Determine whether the PingCheck contains valid fields.  This can be
// used to guard against sending illegal values to the Pingdom API
func (pc *PingCheck) valid() error {
	if err := pc.BaseCheck.valid(); err != nil {
		return err
	}

	return nil
}

func valuesToStrings(orig map[string]interface{}) map[string]string {
	m := make(map[string]string)
	for k, val := range orig {
		if val == nil {
			// we drop all unset optionial fields
			continue
		}
		switch v := val.(type) {
		case string:
			m[k] = v
		case *string:
			if v != nil {
				m[k] = *v
			}
		case bool:
			m[k] = strconv.FormatBool(v)
		case *bool:
			if v != nil {
				m[k] = strconv.FormatBool(*v)
			}
		case int:
			m[k] = strconv.Itoa(v)
		case *int:
			if v != nil {
				m[k] = strconv.Itoa(*v)
			}
		case *[]int:
			if v != nil {
				m[k] = intListToCDString(*v)
			}
		default:
			panic("unexpected type")
		}
	}
	return m
}

func intListToCDString(integers []int) string {
	var CDString string
	for i, item := range integers {
		if i == 0 {
			CDString = strconv.Itoa(item)
		} else {
			CDString = fmt.Sprintf("%v,%d", CDString, item)
		}
	}
	return CDString
}
