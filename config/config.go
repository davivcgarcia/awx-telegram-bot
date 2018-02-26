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

package config

import (
	"fmt"
	"log"
	"os"
)

// EnvVars stores configurations loaded from environment variables.
var EnvVars map[string]string

// This init function will collect environment variables and load them into the EnvVars map.
// At the end, it will verify if any environment variables are empty, and panic if so.
func init() {
	log.Println("[INFO] Loading configuration")
	EnvVars = map[string]string{
		"TELEGRAM_API_TOKEN":     os.Getenv("TELEGRAM_API_TOKEN"),
		"TELEGRAM_CHAT_ID":       os.Getenv("TELEGRAM_CHAT_ID"),
		"AWX_URL":                os.Getenv("AWX_URL"),
		"AWX_USER":               os.Getenv("AWX_USER"),
		"AWX_PASS":               os.Getenv("AWX_PASS"),
		"AWX_JOB_TEMPLATE_START": os.Getenv("AWX_JOB_TEMPLATE_START"),
		"AWX_JOB_TEMPLATE_STOP":  os.Getenv("AWX_JOB_TEMPLATE_STOP"),
		"MONGO_URL":              os.Getenv("MONGO_URL"),
		"MONGO_DB":               os.Getenv("MONGO_DB"),
		"MONGO_COLLECTION":       os.Getenv("MONGO_COLLECTION"),
	}
	for key, value := range EnvVars {
		if value == "" {
			log.Panicln(fmt.Sprintf("Missing required parameters: %s", key))
		}
	}
}
