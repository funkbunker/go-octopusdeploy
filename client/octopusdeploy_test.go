package client

import (
	"net/http"
	"os"
	"testing"

	"github.com/dghubble/sling"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := &http.Client{}
	url := os.Getenv("OCTOPUS_URL")
	apiKey := os.Getenv("OCTOPUS_APIKEY")
	spaceID := os.Getenv("OCTOPUS_SPACE_ID")

	testCases := []struct {
		name    string
		isValid bool
		client  *http.Client
		url     string
		apiKey  string
		spaceID string
	}{
		{"EmptyURL", false, client, emptyString, apiKey, spaceID},
		{"EmptyURLWithSpace", false, client, whitespaceString, apiKey, spaceID},
		{"EmptyAPIKey", false, client, url, emptyString, emptyString},
		{"EmptyAPIKeyWithSpace", false, client, url, whitespaceString, spaceID},
		{"InvalidAPIKey", false, client, url, "API-***************************", spaceID},
		{"ValidAPIKey", true, client, url, apiKey, spaceID},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := NewClient(tc.client, tc.url, tc.apiKey, tc.spaceID)

			if !tc.isValid {
				assert.Error(t, err)
				assert.Nil(t, client)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, client)
		})
	}
}

func TestGetWithEmptyParameters(t *testing.T) {
	resource, err := apiGet(nil, nil, emptyString)

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestGetWithEmptySling(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	resource, err := apiGet(nil, input, "fake-path")

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestGetWithEmptyPath(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	resource, err := apiGet(sling.New(), input, emptyString)

	assert.Error(t, err)
	assert.Nil(t, resource)

	resource, err = apiGet(sling.New(), input, whitespaceString)

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestAddWithEmptyParameters(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiAdd(nil, input, response, emptyString)

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestAddWithEmptySling(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiAdd(nil, input, response, "fake-path")

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestAddWithEmptyPath(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiAdd(sling.New(), input, response, emptyString)

	assert.Error(t, err)
	assert.Nil(t, resource)

	resource, err = apiAdd(sling.New(), input, response, whitespaceString)

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestPostWithEmptyParameters(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiPost(nil, input, response, emptyString)

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestPostWithEmptySling(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiPost(nil, input, response, "fake-path")

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestPostWithEmptyPath(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiPost(sling.New(), input, response, emptyString)

	assert.Error(t, err)
	assert.Nil(t, resource)

	resource, err = apiPost(sling.New(), input, response, whitespaceString)

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestUpdateWithEmptyParameters(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiUpdate(nil, input, response, emptyString)

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestUpdateWithEmptySling(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiUpdate(nil, input, response, "fake-path")

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestUpdateWithEmptyPath(t *testing.T) {
	input := &inputTestValueStruct{test: "fake-value"}
	response := &inputTestResponseStruct{test: "fake-value"}

	resource, err := apiUpdate(sling.New(), input, response, emptyString)

	assert.Error(t, err)
	assert.Nil(t, resource)

	resource, err = apiUpdate(sling.New(), input, response, whitespaceString)

	assert.Error(t, err)
	assert.Nil(t, resource)
}

func TestDeleteWithEmptyParameters(t *testing.T) {
	err := apiDelete(nil, emptyString)
	assert.Error(t, err)
}

func TestDeleteWithEmptySling(t *testing.T) {
	err := apiDelete(nil, "fake-path")
	assert.Error(t, err)
}

func TestDeleteWithEmptyPath(t *testing.T) {
	err := apiDelete(nil, emptyString)
	assert.Error(t, err)

	err = apiDelete(nil, whitespaceString)
	assert.Error(t, err)
}

type inputTestValueStruct struct {
	test string
}

func (i *inputTestValueStruct) GetID() string {
	return "fake-ID"
}

func (i *inputTestValueStruct) Validate() error {
	return nil
}

type inputTestResponseStruct struct {
	test string
}

func (i *inputTestResponseStruct) GetID() string {
	return "fake-ID"
}

func (i *inputTestResponseStruct) Validate() error {
	return nil
}
