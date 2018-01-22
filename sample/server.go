package main

import (
	"fmt"
	"github.com/cheikhshift/form"
	"log"
	"net/http"
)

const Hostsubstr = "%s://%s/%s"

func redirect(port int) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// remove/add not default ports from req.Host
		schema := "http"
		if port == 443 {
			schema = "https"
		}
		target := fmt.Sprintf(Hostsubstr, schema, req.URL.Host, req.URL.Path)
		if len(req.URL.RawQuery) > 0 {
			target += fmt.Sprintf("?%s", req.URL.RawQuery)
		}

		http.Redirect(w, req, target, 301)
	}
}

func main() {

	fmt.Println("Hello monde!")
	go http.ListenAndServe(":80", http.HandlerFunc(redirect(8090)))

	http.HandleFunc("/home/random/path", form.Handler)
	http.HandleFunc("/home/random/path2", form.Handler)

	errgos := http.ListenAndServe(":8090", nil)
	if errgos != nil {
		log.Fatal(errgos)
	}

}
