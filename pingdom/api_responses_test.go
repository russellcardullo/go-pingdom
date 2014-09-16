package pingdom

import "testing"

func TestPingdomError(t *testing.T) {
	pe := PingdomError{StatusCode: 400, StatusDesc: "Bad Request", Message: "Missing param foo"}
	want := "400 Bad Request: Missing param foo"
	if e := pe.Error(); e != want {
		t.Errorf("Error() returned '%+v', want '%+v'", e, want)
	}
}
