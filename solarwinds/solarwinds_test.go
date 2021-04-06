package solarwinds

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
	"net/http"
	"net/http/httptest"
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
	client, _ = NewClient(ClientConfig{
		Username: "chszchen@nordcloud.com",
		Password: "abcdefg",
	})
	client.baseURL = server.URL
}

func teardown() {
	server.Close()
}

func TestNewClient(t *testing.T) {
	c, err := NewClient(ClientConfig{})
	assert.NoError(t, err)
	assert.Equal(t, c.baseURL, defaultBaseURL)
}

func TestExtractCSRFToken(t *testing.T) {
	doc, _ := html.Parse(strings.NewReader(`
<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<link rel="shortcut icon" href="https://cdn.solarwinds.cloud/common-settings/v713/favicon.ico" />
	<link href="https://fonts.googleapis.com/css?family=Open+Sans:400,600&display=swap" rel="stylesheet">
	<meta name="viewport" content="width=device-width,initial-scale=1" />
	<meta name="theme-color" content="#000000" />
	<link rel="manifest" href="https://cdn.solarwinds.cloud/common-settings/v713/manifest.json" />
	<title>Settings</title>
	<script src="//assets.adobedtm.com/764583179334/ccf7525f7c9f/launch-d8cc1fb68e59.min.js"></script>
	<link href="https://cdn.solarwinds.cloud/common-settings/v713/static/css/2.221d85f4.chunk.css" rel="stylesheet">
	<meta name="csrf-token" content="fbO8qrEt-qGJ3jtQctuzcbVfBD47Quy-RE_Q" />
	<script>
		window.serverEnv={"DEPLOYMENT_ENVIRONMENT":"production","APPOPTICS_URL":"https://my.appoptics.com","APP_VERSION":"v450","APP_URL":"https://my.solarwinds.cloud"}
	</script>
</head>

<body><noscript>You need to enable JavaScript to run this app.</noscript>
	<div id="root"></div>
	<script type="text/javascript">
		try{_satellite.pageBottom()}catch{console.log("Adobe Analytics failed to load")}
	</script>
	<script>
		!function(e,n,t,a,i){var c,o,s,d,p;for((i=e[a]=e[a]||{})._q=i._q||[],o=0,s=(c=["initialize","identify","updateOptions","pageLoad","track"]).length;o<s;++o)!function(e){i[e]=i[e]||function(){i._q[e===c[0]?"unshift":"push"]([e].concat([].slice.call(arguments,0)))}}(c[o]);(d=n.createElement(t)).async=!0,d.src="https://cdn.pendo.io/agent/static/56593e57-5efc-42b3-58a0-4daa189e8069/pendo.js",(p=n.getElementsByTagName(t)[0]).parentNode.insertBefore(d,p)}(window,document,"script","pendo")
	</script>
	<script>
		!function(e){function t(t){for(var n,l,i=t[0],c=t[1],s=t[2],f=0,p=[];f<i.length;f++)l=i[f],Object.prototype.hasOwnProperty.call(o,l)&&o[l]&&p.push(o[l][0]),o[l]=0;for(n in c)Object.prototype.hasOwnProperty.call(c,n)&&(e[n]=c[n]);for(a&&a(t);p.length;)p.shift()();return u.push.apply(u,s||[]),r()}function r(){for(var e,t=0;t<u.length;t++){for(var r=u[t],n=!0,i=1;i<r.length;i++){var c=r[i];0!==o[c]&&(n=!1)}n&&(u.splice(t--,1),e=l(l.s=r[0]))}return e}var n={},o={1:0},u=[];function l(t){if(n[t])return n[t].exports;var r=n[t]={i:t,l:!1,exports:{}};return e[t].call(r.exports,r,r.exports,l),r.l=!0,r.exports}l.m=e,l.c=n,l.d=function(e,t,r){l.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:r})},l.r=function(e){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},l.t=function(e,t){if(1&t&&(e=l(e)),8&t)return e;if(4&t&&"object"==typeof e&&e&&e.__esModule)return e;var r=Object.create(null);if(l.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var n in e)l.d(r,n,function(t){return e[t]}.bind(null,n));return r},l.n=function(e){var t=e&&e.__esModule?function(){return e.default}:function(){return e};return l.d(t,"a",t),t},l.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},l.p="https://cdn.solarwinds.cloud/common-settings/v713/";var i=this["webpackJsonp@solarwindscloud/settings-client"]=this["webpackJsonp@solarwindscloud/settings-client"]||[],c=i.push.bind(i);i.push=t,i=i.slice();for(var s=0;s<i.length;s++)t(i[s]);var a=c;r()}([])
	</script>
	<script src="https://cdn.solarwinds.cloud/common-settings/v713/static/js/2.acaf28fd.chunk.js"></script>
	<script src="https://cdn.solarwinds.cloud/common-settings/v713/static/js/main.82b65646.chunk.js"></script>
</body>
</html>
`))
	token, err := extractCSRFToken(doc)
	assert.NoError(t, err)
	assert.Equal(t, "fbO8qrEt-qGJ3jtQctuzcbVfBD47Quy-RE_Q", token)
}

func TestLogin(t *testing.T) {
	setup()
	defer teardown()

	swicus := RandString(10)
	mux.HandleFunc("/v1/login", func(w http.ResponseWriter, r *http.Request) {
		if m := "POST"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		state := r.URL.Query()["state"]
		w.Header().Add(headerNameSetCookie, "Swicus-auth=; Path=/; Expires=Thu, 01 Jan 1970 00:00:00 GMT; HttpOnly; Secure")
		w.Header().Add(headerNameSetCookie, fmt.Sprintf("%v=%v", cookieNameSwicus, swicus)+"; Path=/; Expires=Tue, 06 Apr 2021 09:18:16 GMT; Max-Age=1209600; HttpOnly; Secure; SameSite=None")
		body := fmt.Sprintf(
			`{"RedirectUrl": "https://my.solarwinds.cloud/common/auth/callback?code=txsXjr18udIjF4sdy6fG2fqqHlTK9qY7ePtDppVJiP0.Wb8vBKTZYhdo8GrQBC_-a5nLmDP2thYzsCvkeAfUhS8&scope=openid+Swicus&state=%s"}`,
			state)
		_, _ = fmt.Fprint(w, body)
	})
	result, err := client.login()
	assert.NoError(t, err)
	assert.Equal(t, swicus, result.Swicus)
}

func TestObtainSwiSettings(t *testing.T) {
	setup()
	defer teardown()
	swiSettings := RandString(10)
	mux.HandleFunc("/common/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add(headerNameSetCookie, fmt.Sprintf("%v=%v", cookieNameSwiSettings, swiSettings)+"; Path=/; Expires=Tue, 06 Apr 2021 11:14:34 GMT; HttpOnly; Secure; SameSite=None")
		http.Redirect(w, r, "/foo", http.StatusFound)
	})
	err := client.obtainSwiSettings()
	assert.NoError(t, err)
	assert.Equal(t, swiSettings, client.swiSettings)
}

func TestObtainToken(t *testing.T) {
	setup()
	defer teardown()

	tokenStr := "fbO8qrEt-qGJ3jtQctuzcbVfBD47Quy-RE_Q"
	mux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `
<!doctype html>
<html>
<head>
	<meta charset="utf-8" />
	<link rel="shortcut icon" href="https://cdn.solarwinds.cloud/common-settings/v713/favicon.ico" />
	<link href="https://fonts.googleapis.com/css?family=Open+Sans:400,600&display=swap" rel="stylesheet">
	<meta name="viewport" content="width=device-width,initial-scale=1" />
	<meta name="theme-color" content="#000000" />
	<link rel="manifest" href="https://cdn.solarwinds.cloud/common-settings/v713/manifest.json" />
	<title>Settings</title>
	<script src="//assets.adobedtm.com/764583179334/ccf7525f7c9f/launch-d8cc1fb68e59.min.js"></script>
	<link href="https://cdn.solarwinds.cloud/common-settings/v713/static/css/2.221d85f4.chunk.css" rel="stylesheet">
	<meta name="csrf-token" content="fbO8qrEt-qGJ3jtQctuzcbVfBD47Quy-RE_Q" />
	<script>
		window.serverEnv={"DEPLOYMENT_ENVIRONMENT":"production","APPOPTICS_URL":"https://my.appoptics.com","APP_VERSION":"v450","APP_URL":"https://my.solarwinds.cloud"}
	</script>
</head>

<body><noscript>You need to enable JavaScript to run this app.</noscript>
	<div id="root"></div>
	<script type="text/javascript">
		try{_satellite.pageBottom()}catch{console.log("Adobe Analytics failed to load")}
	</script>
	<script>
		!function(e,n,t,a,i){var c,o,s,d,p;for((i=e[a]=e[a]||{})._q=i._q||[],o=0,s=(c=["initialize","identify","updateOptions","pageLoad","track"]).length;o<s;++o)!function(e){i[e]=i[e]||function(){i._q[e===c[0]?"unshift":"push"]([e].concat([].slice.call(arguments,0)))}}(c[o]);(d=n.createElement(t)).async=!0,d.src="https://cdn.pendo.io/agent/static/56593e57-5efc-42b3-58a0-4daa189e8069/pendo.js",(p=n.getElementsByTagName(t)[0]).parentNode.insertBefore(d,p)}(window,document,"script","pendo")
	</script>
	<script>
		!function(e){function t(t){for(var n,l,i=t[0],c=t[1],s=t[2],f=0,p=[];f<i.length;f++)l=i[f],Object.prototype.hasOwnProperty.call(o,l)&&o[l]&&p.push(o[l][0]),o[l]=0;for(n in c)Object.prototype.hasOwnProperty.call(c,n)&&(e[n]=c[n]);for(a&&a(t);p.length;)p.shift()();return u.push.apply(u,s||[]),r()}function r(){for(var e,t=0;t<u.length;t++){for(var r=u[t],n=!0,i=1;i<r.length;i++){var c=r[i];0!==o[c]&&(n=!1)}n&&(u.splice(t--,1),e=l(l.s=r[0]))}return e}var n={},o={1:0},u=[];function l(t){if(n[t])return n[t].exports;var r=n[t]={i:t,l:!1,exports:{}};return e[t].call(r.exports,r,r.exports,l),r.l=!0,r.exports}l.m=e,l.c=n,l.d=function(e,t,r){l.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:r})},l.r=function(e){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},l.t=function(e,t){if(1&t&&(e=l(e)),8&t)return e;if(4&t&&"object"==typeof e&&e&&e.__esModule)return e;var r=Object.create(null);if(l.r(r),Object.defineProperty(r,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var n in e)l.d(r,n,function(t){return e[t]}.bind(null,n));return r},l.n=function(e){var t=e&&e.__esModule?function(){return e.default}:function(){return e};return l.d(t,"a",t),t},l.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},l.p="https://cdn.solarwinds.cloud/common-settings/v713/";var i=this["webpackJsonp@solarwindscloud/settings-client"]=this["webpackJsonp@solarwindscloud/settings-client"]||[],c=i.push.bind(i);i.push=t,i=i.slice();for(var s=0;s<i.length;s++)t(i[s]);var a=c;r()}([])
	</script>
	<script src="https://cdn.solarwinds.cloud/common-settings/v713/static/js/2.acaf28fd.chunk.js"></script>
	<script src="https://cdn.solarwinds.cloud/common-settings/v713/static/js/main.82b65646.chunk.js"></script>
</body>
</html>
`)
	})
	err := client.obtainToken(&loginResult{
		RedirectUrl: server.URL + "/settings",
	})
	assert.NoError(t, err)
	assert.Equal(t, tokenStr, client.csrfToken)
}
