package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

func main() {
	fmt.Println("Hello World")
	var port string
	flag.StringVar(&port, "p", ":8000", "port on which the server will listen")

	flag.Parse()

	cache := map[string][]byte{}

	handlerFunc := func(w http.ResponseWriter, r *http.Request) {
		log.Println("requesting URL ", r.URL.String())
		w.Header().Set("Access-Control-Allow-Origin", "*")

		fname := r.URL.Path
		if len(fname) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Invalid URL %v", fname)
			return
		}

		if data, ok := cache[fname]; ok {
			log.Printf("found file %v in cache", fname)
			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
			return
		}

		b, err := ioutil.ReadFile(path.Clean(fmt.Sprintf("./wwwroot/%v", fname)))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "File not found %v", fname)
			return
		}

		data := struct {
			Type string `json:"type,omitempty"`
			Data string `json:"data,omitempty"`
		}{
			Type: http.DetectContentType(b),
			Data: base64.StdEncoding.EncodeToString(b),
		}

		res, _ := json.Marshal(data)
		cache[fname] = res
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}

	http.HandleFunc("/", handlerFunc)

	log.Println("listening on port :", port)
	http.ListenAndServe(port, nil)

}
