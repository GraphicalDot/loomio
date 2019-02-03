package controllers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"loomio/config"
	"loomio/models"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

func AddUserHandler(appContext *config.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "Invalid post paramteres", http.StatusBadRequest)
		}
	}()
	data, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	//If there is an error reading the request params, The code will panic,
	// You can handle the panic by using function closures as is handled in
	//user login
	if err != nil {
		panic(err.Error())
	}

	//Creating an instance of User struct and Unmarshalling incoming
	// json into the User staruct instance
	var user models.User
	err = json.Unmarshal(data, &user) //address needs to be passed, If you wont pass a pointer,

	//Creating a new instance of response object so that response can be returned in JSON
	if err != nil {
		errString := fmt.Sprintf("%v", err)
		json.NewEncoder(w).Encode(config.AppResponse{errString,
			false, true, nil})
		return http.StatusUnauthorized, nil
	}

	fmt.Println(user)
	uid, err := AddUser(appContext.Database, &user)

	result := config.AppResponse{Message: fmt.Sprintf("%v", err), Success: false, Error: true}

	if err != nil {
		json.NewEncoder(w).Encode(result)
		return http.StatusUnauthorized, nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	result = config.AppResponse{Message: fmt.Sprintf("%v", uid), Success: true, Error: false}

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	json.NewEncoder(w).Encode(result)
	return http.StatusOK, nil
}

func AddUser(dbSession config.Database, user *models.User) (string, error) {

	user_id := uuid.NewV4()
	user.UserID = user_id.String()
	ok := RetreiveSecondaryIndex(dbSession, user.Email)
	if ok {
		log.Println("This is the user id found in secondary index", ok)
		return "", errors.New("The email is already registered")

	}

	var network bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&network)

	err := enc.Encode(user)
	if err != nil {
		log.Println("Error in encoding gob")
		return "", err
	}

	log.Printf("This is the user %v", user)
	err = dbSession.DBSession.Put([]byte(user.UserID), network.Bytes(), nil)
	//dberr := userCollection.Insert(user)
	fmt.Println(err)
	if err != nil {
		log.Println(err)
		return "", err
	}

	AddSecondaryIndex(dbSession, user.Username, user.Email)

	return user.UserID, nil
}
