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

package model

import (
	"log"

	"github.com/davivcgarcia/awx-telegram-bot/config"
	"gopkg.in/mgo.v2"
)

// global MGO Session instance
var mainSession *mgo.Session

// init initialize the MGO Session object
func init() {
	session, err := mgo.Dial(config.EnvVars["MONGO_URL"])
	if err != nil {
		log.Fatalf("[FATAL] Database problem: %s", err)
	}
	mainSession = session
}

// generateConnection creates the objects for database operations
func generateConnection() (*mgo.Session, *mgo.Collection) {
	collection := mainSession.DB(config.EnvVars["MONGO_DB"]).C(config.EnvVars["MONGO_COLLECTION"])
	return mainSession.Copy(), collection
}
