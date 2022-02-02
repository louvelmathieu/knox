package knox

import (
	"encoding/base64"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

/**
 * 401 Handler for all REST API
 */
func error401(w http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.RequestURI, "Method": r.Method}).Info("Error 401 - Unauthorized")
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusUnauthorized,
		Message: "Unauthorized",
	})
	// Dont use http.Error, it force "Content-type" header
}

/**
 * 404 Handler for all REST API
 */
func error404(w http.ResponseWriter, r *http.Request) {

	log.WithFields(log.Fields{"URL": r.RequestURI, "Method": r.Method}).Info("Error 404 - Page not found")
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusNotFound,
		Message: "Page not found",
	})
	// Dont use http.Error, it force "Content-type" header
}

/**
 * 405 Handler for all REST API
 */
func error405(w http.ResponseWriter, r *http.Request) {

	if r.Method == "OPTIONS" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "OPTIONS, HEAD, GET, POST, PUT, DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Content-type", "application/json; charset=utf-8")
		return
	}

	log.WithFields(log.Fields{"URL": r.RequestURI, "Method": r.Method}).Warning("Error 405 - Method not allowed")
	w.Header().Add("Content-type", "application/json; charset=utf-8")
	w.Header().Add("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(Response{
		Status:  http.StatusMethodNotAllowed,
		Message: "Method not allowed",
	})
}

/**
 * robots.txt
 */
func robot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/plain")
	w.Header().Del("Access-Control-Allow-Origin")
	w.Header().Del("Access-Control-Allow-Headers")
	w.Write([]byte("User-agent: *\nDisallow: /"))
}

/**
 * favicon.ico
 * Return blank square favicon
 */
func favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "image/x-icon")
	w.Header().Del("Access-Control-Allow-Origin")
	w.Header().Del("Access-Control-Allow-Headers")
	data := "data:image/x-icon;base64,AAABAAEAEBAAAAAAAABoBQAAFgAAACgAAAAQAAAAIAAAAAEACAAAAAAAAAEAAAAAAAAAAAAAAAEAAAAAAAD///8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
	i := strings.Index(data, ",") // The actual image starts after the ","
	dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(data[i+1:]))
	io.Copy(w, dec)
}

/**
 * Prob URL for K8s
 */
func probe(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(`{"ping":"pong"}`))
}
