//functions that are mapped to the different routes of nodes' API

package main

import (
	"encoding/json"
	// "log"
	"net/http"
	// "os"
	// "strings"
	// "strconv"

	// "ipsorting"

	"github.com/gorilla/mux" //import router
)

// GetAllKeys displays all from the keyVals var
func GetAllKeys(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(keyVals) //200
}

// GetKey returns a single key-value pair
func GetKey(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	if params["key"] == "get_partition_id" {
		json.NewEncoder(w).Encode(clusterID) //200
		return
	}

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

	//error checking inputs
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

// DeleteAll deletes all key:value pairs in the kvs map
func DeleteAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	for key := range keyVals {
		delete(keyVals, key)
	}

	json.NewEncoder(w).Encode(&map[string]string{"Success": "Key-Value Store cleared"}) //200
}
