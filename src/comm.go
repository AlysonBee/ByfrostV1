package main

import (
	"byfrostV1/errors"
	"byfrostV1/mapping"
	"byfrostV1/server"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var Collapse = false
var CollapseList = []string{}
var StrangePatterns = []string{}

type responseData struct {
	Token string
}

type responsePayload struct {
	Collapse     bool     `json:"collapse"`
	CollapseList []string `json:"collapselist"`
	Tokens       []*server.TokenJSON
	Body         *server.BodyJSON
}

func isPathToOpen(response string) bool {
	segments := strings.Split(response, ":")
	if segments[0] == "location" {
		return true
	}
	return false
}

func parseDeref(derefFunction string, derefToken string) string {
	derefHappening := strings.Split(derefFunction, derefToken)
	if len(derefHappening) > 1 {
		return derefHappening[1]
	}
	return derefFunction
}

func RequestCode(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	var resp responseData
	if err := json.Unmarshal(body, &resp); err != nil {
		panic(err)
	}

	var ret responsePayload
	if isPathToOpen(resp.Token) {
		filePath := strings.Split(resp.Token, ":")
		exec.Command("open", filePath[1]).Run()
		return
	} else {
		segments := strings.Split(resp.Token, "-")
		// Check for a dereference ID.

		GlobalNamespace = *toScan
		SearchImport(*toScan, segments[1], segments[0])
		body, status := getDisplay(segments[1], segments[0])

		if status == errors.ErrorCodeAddDisplayNode {
			return
		}
		if body == nil {
			return
		}
		ret.Tokens = body.Tokens
		ret.Body = body
		ret.Collapse = Collapse
		ret.CollapseList = CollapseList
	}

	jsonBytes, err := structToJSON(ret)
	if err != nil {
		panic(err)
	}

	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)

	Collapse = false
	CollapseList = []string{}
}

func StringToBytes(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func HelloIndexer(w http.ResponseWriter, r *http.Request) {
	initDisplay()
	mapping.UsedPositions = [][]int{}
	allToState(RootPath, true)
	Collapse = false
	CollapseList = []string{}

	displayBody, status := getDisplay("main", "")
	if status == errors.ErrorCodeAddDisplayNode {
		return
	}
	jsonBytes, err := structToJSON(displayBody)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)

	return
}

func CodeServer() {

	r := mux.NewRouter()

	r.HandleFunc("/hello-indexer", HelloIndexer)
	r.HandleFunc("/hello-indexer/request", RequestCode).Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})

	handler := c.Handler(r)

	srv := &http.Server{
		Handler: handler,
		Addr:    ":4201",
	}

	srv.ListenAndServe()
}
