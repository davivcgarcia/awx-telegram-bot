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
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Registry type model information for persistent lab registries
type Registry struct {
	Active   bool
	Username string
	CheckIn  time.Time
	CheckOut time.Time
}

// GetAllRegistries retrives all registries from database
func GetAllRegistries() []Registry {
	session, collection := generateConnection()
	defer session.Close()
	var result []Registry
	collection.Find(bson.M{}).All(&result)
	return result
}

// GetAllActiveRegistries retrives all active registries from database
func GetAllActiveRegistries() []Registry {
	session, collection := generateConnection()
	defer session.Close()
	var result []Registry
	collection.Find(bson.M{"active": true}).All(&result)
	return result
}

// GetUserActiveRegistry retrives all active registries for a specific user from database
func GetUserActiveRegistry(username string) Registry {
	session, collection := generateConnection()
	defer session.Close()
	var result Registry
	collection.Find(bson.M{"active": true, "username": username}).One(&result)
	return result
}

// userHasActiveRegistry determines if an user has an active registry persisted (checked-in)
func userHasActiveRegistry(username string) bool {
	session, collection := generateConnection()
	defer session.Close()
	result, _ := collection.Find(bson.M{"active": true, "username": username}).Count()
	switch result {
	case 0:
		return false
	case 1:
		return true
	default:
		panic(fmt.Errorf("[PANIC] Database inconsistent for user %s", username))
	}
}

// CheckInRegistry adds an active registry (check-in) for a user in the database
func CheckInRegistry(username string) error {
	if userHasActiveRegistry(username) {
		return fmt.Errorf("User %s has already an active check-in registry", username)
	}
	session, collection := generateConnection()
	defer session.Close()
	collection.Insert(Registry{Active: true, Username: username, CheckIn: time.Now()})
	return nil
}

// CheckOutRegistry inactivates the active registry (check-out) for a user in the database
func CheckOutRegistry(username string) error {
	if !userHasActiveRegistry(username) {
		return fmt.Errorf("User %s does not have an active check-in registry", username)
	}
	session, collection := generateConnection()
	defer session.Close()
	query := bson.M{"active": true, "username": username}
	change := bson.M{"$set": bson.M{"active": false, "checkout": time.Now()}}
	collection.Update(query, change)
	return nil
}

// CheckOutAllRegistries inactivates all registries (check-out) in the database
func CheckOutAllRegistries() error {
	session, collection := generateConnection()
	defer session.Close()
	query := bson.M{"active": true}
	change := bson.M{"$set": bson.M{"active": false, "checkout": time.Now()}}
	collection.Update(query, change)
	return nil
}
