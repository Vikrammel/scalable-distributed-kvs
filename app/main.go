//main file for running node in kvs
//written by Vikram and Mike
//note: remember, when file is saved in vs code, code is auto formatted and unused
// code is removed. Make sure to comment out unused code or it won't persist

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux" //import router
)

var keyVals map[string]string //map (dictionary) of string:string

// GetAllKeys displays all from the keyVals var
func GetAllKeys(w http.ResponseWriter, r *http.Request) {
	keyVals["apple"] = "red" //hard coded entry to test
	json.NewEncoder(w).Encode(keyVals)
}

// GetKey displays a single data
func GetKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for key, value := range keyVals {
		if key == params["key"] {
			json.NewEncoder(w).Encode(value)
			return
		}
	}
	type noKeyError struct {
		errStr string `json:"errStr"`
		msg    string `json:"msg"`
	}
	json.NewEncoder(w).Encode(&noKeyError{errStr: "error", msg: "key not found"})
}

// // PutKey creates/updates a key
// func PutKey(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
//
// }

// // DeleteKey deletes a key
// func DeleteKey(w http.ResponseWriter, r *http.Request) {
// 	params := mux.Vars(r)
//
// }

func main() {
	router := mux.NewRouter() //init router

	//funcs for routes
	router.HandleFunc("/kv-store", GetAllKeys).Methods("GET")
	router.HandleFunc("/kv-store/{key}", GetKey).Methods("GET")
	// 	router.HandleFunc("/kv-store/{key}", PutKey).Methods("POST")
	// 	router.HandleFunc("/kv-store/{key}", DeleteKey).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
