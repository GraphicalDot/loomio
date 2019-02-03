package controllers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"loomio/config"
	"loomio/models"
	"net/http"

	"github.com/gorilla/mux"
)

func GetUserHandler(appContext *config.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "Invalid post paramteres", http.StatusBadRequest)
		}
	}()
	uid := mux.Vars(r)["id"]
	log.Printf("This is the id %v", uid)
	var user models.User

	getUser(&appContext.Database, uid, &user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.

	result := config.AppResponse{Message: fmt.Sprintf("%v", user), Success: true, Error: false}
	json.NewEncoder(w).Encode(result)
	return http.StatusOK, nil
}

func getUser(dbSession *config.Database, key string, user *models.User) error {
	data, err := dbSession.DBSession.Get([]byte(key), nil)
	if err != nil {
		fmt.Println("Failed to retrieve results:", err)
		return err
	}

	reader := bytes.NewReader(data)
	// make a decoder
	decoder := gob.NewDecoder(reader)
	// decode it int

	decoder.Decode(&user)
	return nil
}
