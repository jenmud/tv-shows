package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test__getIntFromEnvvar(t *testing.T) {
	os.Clearenv()

	err := os.Setenv("PORT", fmt.Sprintf("%d", 8001))
	assert.Nil(t, err)

	port, err := getIntFromEnvvar("PORT")
	assert.Nil(t, err)
	assert.Equal(t, 8001, port)
}

func Test__getIntFromEnvvar__MissingEnvKey(t *testing.T) {
	os.Clearenv()
	port, err := getIntFromEnvvar("NoSuchThing")
	assert.Error(t, err)
	assert.Equal(t, 0, port)
}

func Test__getIntFromEnvvar__ConvertToIntFailed(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("PORT", "800a")
	assert.Nil(t, err)

	port, err := getIntFromEnvvar("PORT")
	assert.Error(t, err)
	assert.Equal(t, 0, port)
}

func TestGetListenPort__FromValidEnvVar(t *testing.T) {
	os.Clearenv()

	err := os.Setenv("PORT", fmt.Sprintf("%d", 8001))
	assert.Nil(t, err)

	port := GetListenPort()
	assert.Equal(t, 8001, port)
}

func TestGetListenPort__MissingEnvVar(t *testing.T) {
	os.Clearenv()

	// Default port is 8000
	port := GetListenPort()
	assert.Equal(t, 8000, port)
}

func TestGetListenPort__FailedConvertToIntFromEnvVar(t *testing.T) {
	os.Clearenv()
	err := os.Setenv("PORT", "800a")
	assert.Nil(t, err)

	// Default port is 8000
	port := GetListenPort()
	assert.Equal(t, 8000, port)
}

func TestIndexHandler__PostWithNobody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Index))
	defer server.Close()

	r := bytes.NewReader([]byte(""))
	res, err := http.Post(server.URL, "application/json", r)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)

	httpErr := &HttpJSONParseError{}
	payload, _ := ioutil.ReadAll(res.Body)
	err = json.Unmarshal(payload, httpErr)
	assert.Nil(t, err)

	assert.Equal(
		t,
		"Could not decode request: JSON parsing failed",
		httpErr.Error,
	)
}

func TestIndexHandler__PostWithValidJSONButIncorrectData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Index))
	defer server.Close()

	r := bytes.NewReader([]byte("{}"))
	res, err := http.Post(server.URL, "application/json", r)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	payload, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "{}", string(payload))
}

func TestIndexHandler__PostValidRequestJSONData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(Index))
	defer server.Close()

	requestData, err := ioutil.ReadFile("testdata/request.json")
	assert.Nil(t, err)

	r := bytes.NewReader(requestData)
	res, err := http.Post(server.URL, "application/json", r)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)

	payload, err := ioutil.ReadAll(res.Body)
	assert.Nil(t, err)

	resp := &ResponseJson{}
	err = json.Unmarshal(payload, resp)
	assert.Nil(t, err)

	// Based on the sample request payload, we are expecting
	// 7 shows with that have a episode count > 0.
	assert.Equal(t, 7, len(resp.Payload))

	// Lets check the first item is what we are expecting.
	item := resp.Payload[0]
	assert.Equal(t, "http://catchup.ninemsn.com.au/img/jump-in/shows/16KidsandCounting1280.jpg", item.Image)
	assert.Equal(t, "show/16kidsandcounting", item.Slug)
	assert.Equal(t, "16 Kids and Counting", item.Title)
}

func TestModel_ParseRequestJSON(t *testing.T) {
	requestData, err := ioutil.ReadFile("testdata/request.json")
	assert.Nil(t, err)

	data, err := ParseRequestJSON(requestData)
	assert.Nil(t, err)

	// According to the sample request payload, we will
	// have 10 items in the payload.
	assert.Equal(t, 10, len(data))
}

func TestModel_ParseRequestJSON__incorrectJSON(t *testing.T) {
	data, err := ParseRequestJSON([]byte("{broken_json}"))
	assert.Error(t, err)
	assert.Equal(t, 0, len(data))
}

func TestModel_FilterTVShowsForDRM(t *testing.T) {
	requestData, err := ioutil.ReadFile("testdata/request.json")
	assert.Nil(t, err)

	data, err := ParseRequestJSON(requestData)
	assert.Nil(t, err)

	shows := FilterTVShowsForDRM(data)

	// Based on the sample request payload, we are expecting
	// 8 shows with DRM's enabled.
	assert.Equal(t, 8, len(shows))
}

func TestModel_FilterTVShowsWithEpisodes(t *testing.T) {
	requestData, err := ioutil.ReadFile("testdata/request.json")
	assert.Nil(t, err)

	data, err := ParseRequestJSON(requestData)
	assert.Nil(t, err)

	shows := FilterTVShowsWithEpisodes(data)

	// Based on the sample request payload, we are expecting
	// 7 shows with that have a episode count > 0.
	assert.Equal(t, 7, len(shows))
}

func TestModel_MakeResponseJson(t *testing.T) {
	mockRequestData := `
	{
		"payload": [
			{
			    "country": "UK",
			    "drm": true,
			    "episodeCount": 3,
			    "image": {
				"showImage": "http://mock.com/pic.jpg"
			    },
			    "slug": "show/mock",
			    "title": "Mock title"
			}
		]
	}
	`
	data, err := ParseRequestJSON([]byte(mockRequestData))
	assert.Nil(t, err)

	response, err := MakeResponseJson(data)
	assert.Nil(t, err)

	expectedResponseJSON := &ResponseJson{}
	err = json.Unmarshal(response, expectedResponseJSON)
	assert.Nil(t, err)

	item := expectedResponseJSON.Payload[0]
	assert.Equal(t, "http://mock.com/pic.jpg", item.Image)
	assert.Equal(t, "show/mock", item.Slug)
	assert.Equal(t, "Mock title", item.Title)
}
