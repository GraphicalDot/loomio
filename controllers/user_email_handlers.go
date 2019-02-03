package controllers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"loomio/config"
	"loomio/models"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUserEmailHandler(appContext *config.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "Invalid post paramteres", http.StatusBadRequest)
		}
	}()
	email := mux.Vars(r)["email"]
	user, err := getUserEmail(&appContext.Database, email)
	//getUser(&appContext.Database, uid, &user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var result config.AppResponse
	if err != nil {
		result = config.AppResponse{Message: "User is not present", Success: true, Error: false}

	} else {
		result = config.AppResponse{Message: fmt.Sprintf("%v", *user), Success: true, Error: false}

	}

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	json.NewEncoder(w).Encode(result)

	return http.StatusOK, nil
}

func getUserEmail(dbSession *config.Database, email string) (*models.User, error) {
	iter := dbSession.DBSession.NewIterator(nil, nil)
	var user models.User

	for iter.Next() {
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		//key := iter.Key()

		value := iter.Value()

		reader := bytes.NewReader(value)
		// make a decoder
		decoder := gob.NewDecoder(reader)
		// decode it int

		decoder.Decode(&user)
		if user.Email == email {
			log.Println("User matched", user)
			return &user, nil
		}
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		log.Println("Error in the ietration object", err)
	}
	return nil, errors.New("Iteration ended")
}
