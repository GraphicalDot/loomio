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

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	fmt.Println(err)
	fmt.Println(uid)

	if err != nil {
		json.NewEncoder(w).Encode(`{message: fmt.Sprintf("%v", err), error: true, success:false}`)
		return http.StatusUnauthorized, nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	result := config.AppResponse{Message: fmt.Sprintf("%v", uid), Success: true, Error: false}

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	json.NewEncoder(w).Encode(result)
	return http.StatusOK, nil
}

func AddUser(dbSession config.Database, user *models.User) (string, error) {

	user_id := uuid.NewV4()
	user.UserID = user_id.String()
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

func FindUser(userCollection *mgo.Collection, email string, password string) bool {

	var user models.User
	dberr := userCollection.Find(bson.M{"email": email, "password": password}).One(&user)

	if dberr != nil {
		log.Fatal(dberr)
		return false
	}
	return true
}

// UpdateAlbum updates an Album in the DB (not used for now)
func UpdateAlbum(userCollection *mgo.Collection, user models.User) bool {
	dberr := userCollection.UpdateId(user.UserID, user)

	if dberr != nil {
		log.Println(dberr)
		return false
	}
	return true
}
