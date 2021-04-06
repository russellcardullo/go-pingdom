package acceptance

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/nordcloud/go-pingdom/pingdomext"
	"github.com/stretchr/testify/assert"
)

var client_ext *pingdomext.Client

var runExtAcceptance bool

func init() {
	if os.Getenv("PINGDOM_EXT_ACCEPTANCE") == "1" {
		runExtAcceptance = true

		config := pingdomext.ClientConfig{

			HTTPClient: &http.Client{
				Timeout: time.Second * 10,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			},
		}
		client_ext, _ = pingdomext.NewClientWithConfig(config)

	}
}

func TestIntegrations(t *testing.T) {
	if !runExtAcceptance {
		t.Skip()
	}

	integration := pingdomext.WebHookIntegration{
		Active:     false,
		ProviderID: 2,
		UserData: &pingdomext.WebHookData{
			Name: "wlwu-tets-1",
			URL:  "http://www.example.com",
		},
	}

	createMsg, err := client_ext.Integrations.Create(&integration)
	assert.NoError(t, err)
	assert.NotNil(t, createMsg)
	assert.NotEmpty(t, createMsg)
	assert.True(t, createMsg.Status)

	integrationID := createMsg.ID

	listMsg, err := client_ext.Integrations.List()
	assert.NoError(t, err)
	assert.NotNil(t, listMsg)
	assert.NotEmpty(t, listMsg)

	getMsg, err := client_ext.Integrations.Read(integrationID)
	assert.NoError(t, err)
	assert.NotNil(t, getMsg)
	assert.NotEmpty(t, getMsg)
	assert.NotEmpty(t, getMsg.CreatedAt)
	assert.NotEmpty(t, getMsg.Name)

	integration.Active = true
	integration.UserData.Name = "wlwu-tets-update"
	integration.UserData.URL = "http://www.example1.com"

	updateMsg, err := client_ext.Integrations.Update(integrationID, &integration)
	assert.NoError(t, err)
	assert.NotNil(t, updateMsg)
	assert.True(t, updateMsg.Status)

	delMsg, err := client_ext.Integrations.Delete(integrationID)
	assert.NoError(t, err)
	assert.NotNil(t, delMsg)
	assert.True(t, delMsg.Status)

	listProviderMsg, err := client_ext.Integrations.ListProviders()
	assert.NoError(t, err)
	assert.NotNil(t, listProviderMsg)
	assert.Equal(t, len(listProviderMsg), 2)

}
