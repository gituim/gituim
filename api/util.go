package api

import (
	"com/gitlab/gituim/repository"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

func getVar(w http.ResponseWriter, r *http.Request, varName string) (string, bool) {
	vars := mux.Vars(r)
	value, ok := vars[varName]
	if !ok {
		http.Error(w, "invalid repository name", http.StatusBadRequest)
		return "", ok
	}
	return value, ok
}

func handleError(err error, w http.ResponseWriter) {
	if errors.Is(err, repository.NotFoundError) {
		http.Error(w, repository.NotFoundError.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
