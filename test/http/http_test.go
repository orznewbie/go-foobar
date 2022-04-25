package http

import (
	"github.com/orznewbie/gotmpl/pkg/log"
	"net/http"
	"testing"
)

func sayHello(w http.ResponseWriter, r *http.Request) {
	log.Info(r)
	w.Write([]byte("Hello World!"))
}

func sayByeBye(w http.ResponseWriter, r *http.Request) {
	log.Info(r)
	w.Write([]byte("Bye Bye!"))
}

func TestHttp(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello/", sayHello)
	mux.HandleFunc("/bye", sayByeBye)

	if err := http.ListenAndServe(":1234", mux); err != nil {
		t.Fatal(err)
	}
}
