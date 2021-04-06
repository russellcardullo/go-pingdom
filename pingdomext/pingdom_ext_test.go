package pingdomext

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/nordcloud/go-pingdom/pingdom"
	"github.com/stretchr/testify/assert"
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
	client = &Client{
		JWTToken: "my_jwt_token",
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		Integrations: nil,
	}
	client.Integrations = &IntegrationService{client: client}

	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func teardown() {
	server.Close()
}

func TestNewClientWithConfig(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		w.Header().Add("Set-Cookie", "pingdom_login_session_id=qw4us4Ed7aLSGugMRDHkqM9G6mwuKdn9Hz90r6IHhRc%3D; Path=/; HttpOnly; Secure")
		w.Header().Add("Location", "https://my.solarwinds.cloud/login?response_type=code&scope=openid%20swicus&client_id=pingdom&state=htILEppzoMPtb6UjOdM98XPS3Mcwkr3Y&redirect_uri=https%3A%2F%2Fmy.pingdom.com%2Fauth%2Fswicus%2Fcallback")
		_, _ = fmt.Fprintf(w, "{}")
	})

	mux.HandleFunc("/v1/login", func(w http.ResponseWriter, r *http.Request) {
		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		_, _ = fmt.Fprintf(w,
			`{"RedirectUrl": "https://my.pingdom.com/auth/swicus/callback?code=70kRkkAB7OIv5YYTPR6LpHH-2jMbtaDEHScLDw1amfw.baMoW3w-HkNXOj_I8pv580mRwBjIRVdFLW3cXFGRX9o&scope=openid+swicus&state=htILEppzoMPtb6UjOdM98XPS3Mcwkr3Y"}`,
		)
	})

	mux.HandleFunc("/auth/swicus/callback", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		w.Header().Add("Set-Cookie", "jwt=my_test_token")
		_, _ = fmt.Fprintf(w, "{}")
	})

	url, err := url.Parse(server.URL)
	assert.NotEmpty(t, url)
	assert.NoError(t, err)

	c, err := NewClientWithConfig(ClientConfig{
		Username: "test_user",
		Password: "test_pwd",
		BaseURL:  url.String(),
		AuthURL:  url.String() + "/v1/login",
		HTTPClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, c.JWTToken, "my_test_token")
	assert.NotNil(t, c.Integrations)
}

func TestNewClientWithConfig2(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		w.Header().Add("Set-Cookie", "pingdom_login_session_id=qw4us4Ed7aLSGugMRDHkqM9G6mwuKdn9Hz90r6IHhRc%3D; Path=/; HttpOnly; Secure")
		w.Header().Add("Location", "https://my.solarwinds.cloud/login?response_type=code&scope=openid%20swicus&client_id=pingdom&state=htILEppzoMPtb6UjOdM98XPS3Mcwkr3Y&redirect_uri=https%3A%2F%2Fmy.pingdom.com%2Fauth%2Fswicus%2Fcallback")
		_, _ = fmt.Fprintf(w, "{}")
	})

	url, err := url.Parse(server.URL)
	assert.NotEmpty(t, url)
	assert.NoError(t, err)

	c, err := NewClientWithConfig(ClientConfig{
		Username: "test_user",
		Password: "test_pwd",
		BaseURL:  url.String(),
		AuthURL:  url.String() + "/v1/login",
		HTTPClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
	})
	assert.Error(t, err)
	assert.Nil(t, c)

}

func TestClient_NewRequest(t *testing.T) {

	setup()
	defer teardown()

	req, err := client.NewRequest("GET", "/data/v3/integration", nil)

	assert.NoError(t, err)
	assert.Equal(t, "GET", req.Method)
	assert.Equal(t, client.BaseURL.String()+"/data/v3/integration", req.URL.String())
}

func TestClient_Do(t *testing.T) {
	setup()
	defer teardown()
	type foo struct {
		A string
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `{"A":"a"}`)
	})

	req, _ := client.NewRequest("GET", "/", nil)
	body := new(foo)
	want := &foo{"a"}

	_, err := client.Do(req, body)
	assert.NoError(t, err)
	assert.Equal(t, want, body)
}

func Test_decodeResponse(t *testing.T) {
	type args struct {
		r *http.Response
		v interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				r: &http.Response{
					Body: ioutil.NopCloser(strings.NewReader(`
					{
						"integration": {
							"status": true,
							"id": 112396
						}
					}`)),
				},
				v: &integrationJSONResponse{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := decodeResponse(tt.args.r, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("decodeResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
func TestValidateResponse(t *testing.T) {
	valid := &http.Response{
		Request:    &http.Request{},
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(strings.NewReader("OK")),
	}

	assert.NoError(t, validateResponse(valid))

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

	want := &pingdom.PingdomError{StatusCode: 400, StatusDesc: "Bad Request", Message: "This is an error"}
	assert.Equal(t, want, validateResponse(invalid))
}

func Test_getCookie(t *testing.T) {
	type args struct {
		resp *http.Response
		name string
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Cookie
		wantErr bool
	}{
		{
			name: "response with session cookie",
			args: args{
				name: "pingdom_login_session_id",
				resp: &http.Response{
					Header: http.Header{
						"Set-Cookie": {"pingdom_login_session_id=xxxxxxx", "Path=/", "HttpOnly", "Secure"},
					},
				},
			},
			want: &http.Cookie{
				Name:  "pingdom_login_session_id",
				Value: "xxxxxxx",
				Raw:   "pingdom_login_session=xxxxxxx",
			},
			wantErr: false,
		},
		{
			name: "response without session cookie",
			args: args{
				name: "pingdom_login_session_id",
				resp: &http.Response{
					Header: http.Header{
						"Set-Cookie": {"pingdom_login_session=xxxxxxx", "Path=/", "HttpOnly", "Secure"},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCookie(tt.args.resp, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("getCookie() = %v, want %v", got, tt.want)
			}
		})
	}
}
