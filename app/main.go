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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyVals) //200
	return
}

// GetKey displays a single data
func GetKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for key, value := range keyVals {
		if key == params["key"] {
			json.NewEncoder(w).Encode(value) //200
			return
		}
	}

	w.WriteHeader(http.StatusNotFound) //404
	json.NewEncoder(w).Encode(&map[string]string{"Error": "Key not found"})
}

// PutKey creates/updates a key
func PutKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	value := ""

	//erro checking inputs
	if r.Body == nil {
		//http.Error(w, "Please send a request body", 400)
		w.WriteHeader(http.StatusBadRequest) //400
		json.NewEncoder(w).Encode(&map[string]string{"Error": "No request body"})
		return
	}
	value = r.PostFormValue("val")
	if !(len(value) > 0) {
		value = ""
	} //make sure value is a valid empty str
	keyVals[params["key"]] = value //store/update user's value for key

	json.NewEncoder(w).Encode(&map[string]string{"Success": "key value pair {'" + params["key"] + "':'" + value + "'} updated"}) //200
}

// DeleteKey deletes a key
func DeleteKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for key, value := range keyVals {
		if key == params["key"] {
			delete(keyVals, params["key"])
			json.NewEncoder(w).Encode(&map[string]string{"Success": "key value pair {'" + params["key"] + "':'" + value + "'} deleted"}) //200
			return
		}
	}

	w.WriteHeader(http.StatusNotFound) //404
	json.NewEncoder(w).Encode(&map[string]string{"Error": "Key not found"})
}

func main() {
	router := mux.NewRouter()              //init router
	keyVals = make(map[string]string, 100) //initialize with 100 keys

	//hard coded entries to test
	keyVals["apple"] = "red"
	keyVals["banana"] = "yellow"
	keyVals["lime"] = "green"
	keyVals["blueberry"] = "blue"
	keyVals["orange"] = "orange"
	keyVals["grape"] = "purple"

	//funcs for routes (with and without slashes at the end of URL)
	router.HandleFunc("/kv-store", GetAllKeys).Methods("GET")
	router.HandleFunc("/kv-store/", GetAllKeys).Methods("GET")
	router.HandleFunc("/kv-store/{key}", GetKey).Methods("GET")
	router.HandleFunc("/kv-store/{key}/", GetKey).Methods("GET")
	router.HandleFunc("/kv-store/{key}", PutKey).Methods("PUT")
	router.HandleFunc("/kv-store/{key}/", PutKey).Methods("PUT")
	router.HandleFunc("/kv-store/{key}", DeleteKey).Methods("DEL")
	router.HandleFunc("/kv-store/{key}/", DeleteKey).Methods("DEL")

	log.Fatal(http.ListenAndServe("localhost:1337", router))
}
