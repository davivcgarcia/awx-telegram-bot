/*
Copyright 2018 Davi Garcia (davivcgarcia@gmail.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/davivcgarcia/awx-telegram-bot/config"
)

// AWXAuthToken stores authentication token retrived from AWX Auth API.
type AWXAuthToken struct {
	Token   string
	Expires string
}

// awxAuthenticate does Authentication using AWX API, returning the Token value.
func awxAuthenticate() (string, error) {
	jsonData := map[string]string{"username": config.EnvVars["AWX_USER"], "password": config.EnvVars["AWX_PASS"]}
	jsonValue, _ := json.Marshal(jsonData)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	response, err := client.Post(config.EnvVars["AWX_URL"]+"/authtoken/", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("[ERROR] Failed executing HTTP request with error %s\n", err)
		return "", err
	}
	data, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Printf("[ERROR] Failed reading the response buffer with error %s\n", err)
		return "", err
	}
	awxToken := &AWXAuthToken{}
	err = json.Unmarshal(data, &awxToken)
	if err != nil {
		log.Printf("[ERROR] Failed parsing response with error %s\n", err)
		return "", err
	}
	log.Printf("[INFO] Ansible AWX API authentication succeeded!")
	return awxToken.Token, nil
}

// awxLaunchJob function makes RESTful API call to AWX and launch a Job Template.
func awxLaunchJob(token string, jobID string) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	req, _ := http.NewRequest("POST", config.EnvVars["AWX_URL"]+"/job_templates/"+jobID+"/launch/", nil)
	req.Header.Add("Authorization", "Token "+token)
	_, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] Failed executing HTTP request with error %s\n", err)
		return err
	}
	log.Printf("[INFO] Ansible AWX API request for Job Template ID launch %s succeeded!", jobID)
	return nil
}

// healthHandler does a dummy health status http endpoint
func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok!")
}

// StartAWXController does AWX Controller initialization
func StartAWXController(exit chan bool) {
	log.Println("[INFO] Starting AWX Controller")
	http.HandleFunc("/health", healthHandler)
	http.ListenAndServe(":8080", nil)
	exit <- true
}
