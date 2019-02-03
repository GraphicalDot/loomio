package controllers

import (
	"go-api/config"
	"net/http"

	//"github.com/davecgh/go-spew/spew"
	"io"
)

func HomeHandler(appContext *config.AppContext, w http.ResponseWriter, r *http.Request) (int, error) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
	return http.StatusOK, nil
}
