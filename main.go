package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func handleRequest(w http.ResponseWriter, url string) {
	log.Println("Handling request: " + url)
	if len(url) < 5 {
		io.WriteString(w, "Cannot redirect to "+url)
		return
	}
	
	w.Header().Set("Access-Control-Allow-Credentials", "false")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		io.WriteString(w, "Error in HTTP Proxy")
		os.Exit(1)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			os.Exit(1)
		}
		fmt.Printf("%s\n", string(contents))
		io.WriteString(w, string(contents))
	}
}

func main() {
	server := http.Server{
		Addr:    ":8007",
		Handler: &myHandler{},
	}
	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	redirectURL := strings.Replace(r.URL.String(), "/_", "http://", 1)
	fmt.Printf("Orig: %s, redirect: %s\n", r.URL.String(), redirectURL)
	handleRequest(w, redirectURL)
}
