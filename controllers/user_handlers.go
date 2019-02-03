package controllers

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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

func GetUserHandler(appContext *config.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	defer func() {
		if r := recover(); r != nil {
			http.Error(w, "Invalid post paramteres", http.StatusBadRequest)
		}
	}()
	uid := mux.Vars(r)["id"]
	log.Printf("This is the id %v", uid)

	getUser(&appContext.Database, uid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	json.NewEncoder(w).Encode(`{"success": true, error: false}`)
	return http.StatusOK, nil
}

func getUser(dbSession *config.Database, key string) error {
	data, err := dbSession.DBSession.Get([]byte(key), nil)
	if err != nil {
		fmt.Println("Failed to retrieve results:", err)
		return err
	}
	fmt.Printf("Failed to retrieve results:%v", data)

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
	err = AddUser(appContext.Database, &user)
	fmt.Println(err)

	if err != nil {
		json.NewEncoder(w).Encode(`{message: fmt.Sprintf("%v", err), error: true, success:false}`)
		return http.StatusUnauthorized, nil
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	json.NewEncoder(w).Encode(`{"success": true, error: false}`)
	return http.StatusOK, nil
}

func AddUser(dbSession config.Database, user *models.User) error {

	user_id := uuid.NewV4()
	user.UserID = user_id.String()
	var network bytes.Buffer // Stand-in for a network connection
	enc := gob.NewEncoder(&network)

	err := enc.Encode(user)
	if err != nil {
		log.Println("Error in encoding gob")
	}

	log.Printf("This is the user %v", user)
	err = dbSession.DBSession.Put([]byte(user.UserID), network.Bytes(), nil)
	//dberr := userCollection.Insert(user)
	fmt.Println(err)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
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
