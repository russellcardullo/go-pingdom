package solarwinds

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

type GraphQLRequest struct {
	OperationName string      `json:"operationName"`
	Variables     interface{} `json:"variables"`
	Query         string      `json:"query"`
	ResponseType  string      `json:"-"` // Not required by GraphQL schema
}

type GraphQLResponse map[string]interface{}

func NewGraphQLResponse(body io.Reader, key string) (*GraphQLResponse, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(b)
	log.Printf("The response body is: %s", bodyStr)
	root := map[string]interface{}{}
	if err := json.NewDecoder(strings.NewReader(bodyStr)).Decode(&root); err != nil {
		return nil, err
	}
	data, ok := root["data"].(map[string]interface{})
	if !ok {
		body, _ := json.Marshal(root)
		return nil, fmt.Errorf("request failed with response: %v", string(body))
	}
	graphQLResp := GraphQLResponse{}
	for k, v := range data[key].(map[string]interface{}) {
		graphQLResp[k] = v
	}
	return &graphQLResp, nil
}

func (r GraphQLResponse) isSuccess() bool {
	if success, ok := r["success"]; ok {
		return success.(bool)
	} else {
		return true
	}
}

func (r GraphQLResponse) message() string {
	if msg, ok := r["message"]; ok {
		return msg.(string)
	} else {
		return ""
	}
}
