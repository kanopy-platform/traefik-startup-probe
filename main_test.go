package main

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDoType func(req *http.Request) (*http.Response, error)

type MockClient struct {
	MockDo MockDoType
}

func (m *MockClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}
func TestCountFrontends(t *testing.T) {
	json := `
		{
			"kubernetes": {
				"backends": {},
				"frontends": {
					"frontend1.example.local/": {},
					"frontend2.example.local/": {}
				}
			}
		}
	`
	resp := io.NopCloser(bytes.NewReader([]byte(json)))
	client := &MockClient{
		MockDo: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       resp,
			}, nil
		},
	}

	result, err := countFrontends(client)
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, 2, result)
}
