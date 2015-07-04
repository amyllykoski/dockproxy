package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Agent configuration, from configuration file.
type Agent struct {
	Name      string `json:"name"`
	IpAddress string `json:"ipAddress"`
	Port      string `json:"port"`
}

var agents []Agent

func readConfiguration(configFile string) {
	file, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&agents)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(agents)
}

// Command line flags.
type Flags struct {
	listenIP   string
	listenPort string
	configFile string
}

type BuildMessage struct {
	BuilderName string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status"`
}

// Stores the latest JSONs from the builders.
var buildMessages map[string]BuildMessage

func getCmdLineArgs() Flags {
	listenIP := flag.String("lip", "0.0.0.0", "listen IP address")
	listenPort := flag.String("lp", "8007", "listen port")
	configFile := flag.String("c", "config.json", "Proxy configuration file")
	flag.Parse()

	return Flags{*listenIP, *listenPort, *configFile}
}

func serveStaticFiles() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	log.Println("Serving static files in port 8004.")
	http.ListenAndServe("0.0.0.0:8004", nil)
}

func handleRequest(w http.ResponseWriter, url string) {
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
		//os.Exit(1)
		return
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			//os.Exit(1)
			return
		}
		//fmt.Printf("%s\n", string(contents))
		io.WriteString(w, string(contents))
	}
}

func main() {

	flags := getCmdLineArgs()
	fmt.Println("Command Line Flags:")
	fmt.Println("[listenUrl: " + flags.listenIP + "]")
	fmt.Println("[listenPort: " + flags.listenPort + "]")
	fmt.Println("[configFile: " + flags.configFile + "]")

	server := http.Server{
		Addr:    flags.listenIP + ":" + flags.listenPort,
		Handler: &myHandler{},
	}

	readConfiguration(flags.configFile)
	buildMessages = make(map[string]BuildMessage)
	go serveStaticFiles()
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

	var buildMessage BuildMessage
	json.Unmarshal(body, &buildMessage)

	switch r.Method {
	case "POST":
		buildMessages[buildMessage.BuilderName] = buildMessage
		log.Println("BuildMessages: %s ", buildMessages)
	case "GET":
		w.Header().Set("Access-Control-Allow-Credentials", "false")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		arr := make([]BuildMessage, 0, len(buildMessages))
		for _, msg := range buildMessages {
			log.Println("Value: ", msg)
			arr = append(arr, msg)
		}
		ret, err := json.Marshal(arr)
		if err != nil {
			panic(err)
		}
		io.WriteString(w, string(ret))
	default:
		log.Println("Unsupported method: ", r.Method)
	}

	//fmt.Println("Body: " + string(body))
}

func handleGetAgentConfiguration(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "false")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	arr := make([]Agent, 0, len(agents))
	for _, msg := range agents {
		arr = append(arr, msg)
	}

	ret, err := json.Marshal(arr)
	if err != nil {
		panic(err)
	}
	io.WriteString(w, string(ret))
	//fmt.Println("Body: " + string(body))
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if strings.Contains(r.URL.String(), "/build") {
		handleBuildMessages(w, r)
		return
	}

	if strings.Contains(r.URL.String(), "/agents") {
		handleGetAgentConfiguration(w, r)
		return
	}

	if strings.Contains(r.URL.String(), "/favicon") {
		return
	}

	redirectURL := strings.Replace(r.URL.String(), "/_", "http://", 1)
	fmt.Printf("Orig: %s, redirect: %s\n", r.URL.String(), redirectURL)
	handleRequest(w, redirectURL)
}
