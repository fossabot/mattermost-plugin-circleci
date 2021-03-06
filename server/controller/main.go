package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"

	"github.com/chetanyakan/mattermost-plugin-circleci/server/config"
	"github.com/chetanyakan/mattermost-plugin-circleci/server/util"
)

type Endpoint struct {
	Path         string
	Method       string
	Execute      func(w http.ResponseWriter, r *http.Request)
	RequiresAuth bool
}

// Endpoints is a map of endpoint key to endpoint object
// Usage: getEndpointKey(GetMetadata): GetMetadata
var Endpoints = map[string]*Endpoint{
	getEndpointKey(circleCIBuildFinished): circleCIBuildFinished,
}

// Uniquely identifies an endpoint using path and method
func getEndpointKey(endpoint *Endpoint) string {
	return util.GetKeyHash(endpoint.Path + "-" + endpoint.Method)
}

// InitAPI initializes the REST API
func InitAPI() *mux.Router {
	r := mux.NewRouter()
	handleStaticFiles(r)

	s := r.PathPrefix("/api/v1").Subrouter()
	for _, endpoint := range Endpoints {
		handler := endpoint.Execute
		if endpoint.RequiresAuth {
			handler = handleAuthRequired(endpoint)
		}
		s.HandleFunc(endpoint.Path, handler).Methods(endpoint.Method)
	}

	return r
}

// handleStaticFiles handles the static files under the assets directory.
func handleStaticFiles(r *mux.Router) {
	bundlePath, err := config.Mattermost.GetBundlePath()
	if err != nil {
		config.Mattermost.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	// This will serve static files from the 'assets' directory under '/static/<filename>'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}

func handleAuthRequired(endpoint *Endpoint) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if Authenticated(w, r) {
			endpoint.Execute(w, r)
		}
	}
}

// Authenticated verifies if provided request is performed by a logged-in Mattermost user.
func Authenticated(w http.ResponseWriter, r *http.Request) bool {
	userID := r.Header.Get(config.HeaderMattermostUserID)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return false
	}

	return true
}

func returnStatusOK(w http.ResponseWriter) {
	w.Header().Set("Content-Alias", "application/json")
	m := make(map[string]string)
	m[model.STATUS] = model.STATUS_OK
	_, _ = w.Write([]byte(model.MapToJson(m)))
}
