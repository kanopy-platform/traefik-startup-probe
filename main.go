package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type Providers struct {
	Kubernetes Kubernetes `json:"kubernetes"`
}

type Kubernetes struct {
	Backends  map[string]interface{} `json:"backends"`
	Frontends Frontends              `json:"frontends"`
}

type Frontends map[string]Frontend

type Frontend map[string]interface{}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	Client HTTPClient
)

func init() {
	Client = &http.Client{}
}

func countFrontends() (int, error) {

	var providers Providers

	req, err := http.NewRequest(http.MethodGet, "http://localhost:8080/api/providers", nil)
	if err != nil {
		return 0, err
	}

	resp, err := Client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&providers)

	if err != nil {
		return 0, err
	}

	return len(providers.Kubernetes.Frontends), nil
}

func main() {
	// start with Service Unavailable
	code := 503

	curFrontends := 0
	prevFrontends := 0
	serverInitialized := false

	logLevelString, ok := os.LookupEnv("LOG_LEVEL")
	if !ok {
		logLevelString = "info"
	}

	logLevel, err := log.ParseLevel(logLevelString)
	if err != nil {
		logLevel = log.DebugLevel
	}
	log.SetLevel(logLevel)

	log.Info("startup probe server started")

	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if !serverInitialized {
			if curFrontends, err = countFrontends(); err != nil {
				log.Error(err)
			}
			log.Debug(fmt.Sprintf("frontends: %v (previous: %v)", curFrontends, prevFrontends))
			if curFrontends != 0 && curFrontends == prevFrontends {
				code = 200
				serverInitialized = true
			}
		}

		prevFrontends = curFrontends

		log.Debug(fmt.Sprintf("serverInitialized: %v", serverInitialized))

		w.WriteHeader(code)
		_, err = w.Write([]byte(http.StatusText(code)))
		if err != nil {
			log.Error(err)
		}

		log.Info(fmt.Sprintf("%v %v %v", r.URL, code, r.UserAgent()))

	})

	log.Fatal(http.ListenAndServe(":8083", nil))

}
