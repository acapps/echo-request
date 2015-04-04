package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/justinas/alice"
)

func main() {

	stdHandler := alice.New(logHandler, recoverHandler)

	http.Handle("/", stdHandler.ThenFunc(defaultHandler))

	log.Infof("Server is listening on port: %s", "9999")
	listener, _ := net.Listen("tcp4", ":9999")

	log.Fatal(http.Serve(listener, nil))

}

func defaultHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ouput, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%s", ouput)
}

// logHandler, is a wrapper for all incoming requests. This just provides
// a centralized approach to logging.
func logHandler(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
