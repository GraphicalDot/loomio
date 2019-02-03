package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
)

type AppResponse struct {
	Message string                 `json:"message"`
	Error   bool                   `json:"error"`
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
}

//If more keys are present in the config.json, then
//their types must be defined here , so that json
//can be unmarshled in those new fields
type LevelDB struct {
	PATH string
}

type settings struct {
	LevelDB LevelDB
}

/*
type MongoType struct {
	Host     string
	DB       string
	Username string
	Password string
}


*/
func Readconfig(path string) *settings {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("File error: %v\n", err)
		os.Exit(1)
	}
	//fmt.Printf("%s\n", string(file))

	//m := new(Dispatch)
	//var m interface{}
	var jsontype settings
	err = json.Unmarshal(file, &jsontype)
	if err != nil {
		fmt.Printf("File error: %v\n", err)

	}
	return &jsontype
}

//If more collections needs to be accessed they must
//be defined here
type Database struct {
	DBSession *leveldb.DB //this needs to be changed and must be type leveldb Object
}

type AppContext struct {
	Database Database
	//RethinkSession *r.Session
}

type ContextHandler struct {
	*AppContext
	Handler func(*AppContext, http.ResponseWriter, *http.Request) (int, error)
}

func (ahandler ContextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	//The mux will run this struct method automatically, This method on the basis
	// of the return from Hanndler will decide its response
	status, err := ahandler.Handler(ahandler.AppContext, w, r)
	if err != nil {
		log.Println("Method not found")
		switch status {
		case http.StatusNotFound:
			http.NotFound(w, r)
		case http.StatusInternalServerError:
			http.Error(w, http.StatusText(status), status)
		default:
			fmt.Println("Not found")
			http.Error(w, http.StatusText(405), 405)
		}
	}
}
