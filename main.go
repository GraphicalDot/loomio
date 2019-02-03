package main

import (
	"fmt"
	"log"
	"loomio/config"
	"loomio/controllers"
	"net/http"
	"time"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/syndtr/goleveldb/leveldb"
)

func loggingMiddleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		t1 := time.Now()
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}
	return http.HandlerFunc(fn)
}

func main() {
	fmt.Println("Here is the loom.io test")
	configuration := config.Readconfig("./config.json")

	r := mux.NewRouter()

	db, err := leveldb.OpenFile(configuration.LevelDB.PATH, nil)
	if err != nil {
		log.Printf("The database couldn be found", err)
	}

	database := config.Database{db}
	context := config.AppContext{Database: database}
	defer db.Close()

	AddUserContextHandler := &config.ContextHandler{&context, controllers.AddUserHandler}
	r.Methods("POST").Path("/user").Name("AddUser").Handler(AddUserContextHandler)

	GetUserContextHandler := &config.ContextHandler{&context, controllers.GetUserHandler}
	r.Methods("GET").Path("/user/{id:[0-9a-zA-Z-]+}").Name("GetUser").Handler(GetUserContextHandler)
	r.Use(loggingMiddleware)

	fmt.Println("Here you go")
	allowedOrigins := handlers.AllowedOrigins([]string{"*"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	srv := &http.Server{
		Addr: "0.0.0.0:8001",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}
	/*

	    // Run our server in a goroutine so that it doesn't block.
	    go func() {
	        if err := srv.ListenAndServe(); err != nil {
	          log.Println(err)
	      }
	  }()
	*/
	log.Fatal(srv.ListenAndServe(), handlers.CORS(allowedMethods, allowedOrigins)(r))
}
