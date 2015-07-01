package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Command line flags.
type flags struct {
	listenIP   string
	listenPort string
}

func getCmdLineArgs() flags {
	listenIP := flag.String("lip", "0.0.0.0", "listen IP address")
	listenPort := flag.String("lp", "8007", "listen port")
	flag.Parse()

	return flags{*listenIP, *listenPort}
}

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

	flags := getCmdLineArgs()
	fmt.Println("Command Line Flags:")
	fmt.Println("[listenUrl: " + flags.listenIP + "]")
	fmt.Println("[listenPort: " + flags.listenPort + "]")

	server := http.Server{
		Addr:    flags.listenIP + ":" + flags.listenPort,
		Handler: &myHandler{},
	}

	server.ListenAndServe()
}

// Build messages are sent to /build Rest endpoint. The POST method is used
// for updating build information. The GET method retrieves the latest.
// JSON is as follows:
//   { "name" : "name of the builder",
//     "image" : "name of the image being built",
//	   "status" : "build | tag | push | done | error"
//   }
func handleBuildMessages(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	switch r.Method {
	case "POST":
		fmt.Println("Got POST")
	case "GET":
		fmt.Println("Got GET")
	default:
		log.Println("Unsupported method: ", r.Method)
	}

	//fmt.Println("Body: " + string(body))
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.String(), "/build") {
		handleBuildMessages(w, r)
		return
	}

	redirectURL := strings.Replace(r.URL.String(), "/_", "http://", 1)
	fmt.Printf("Orig: %s, redirect: %s\n", r.URL.String(), redirectURL)
	handleRequest(w, redirectURL)
}
