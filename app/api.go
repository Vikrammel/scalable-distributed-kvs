//functions that are mapped to the different routes of nodes' API

package main

import (
	"encoding/json"
	// "log"
	// "fmt"
	"net/http"
	// "io/ioutil"
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

func _getDB(w http.ResponseWriter, r *http.Request){
	//return {"result": "success", "dict": json.dumps(d), "causal_payload": json.dumps(vClock),
	//	"timestamp": json.dumps(storedTimeStamp)}, 200
	json.NewEncoder(w).Encode(&map[string]string{"result": "success", "kvs": stringMapToString(keyVals), 
		"causal_payload": stringIntMapToString(vClock), "timestamp": stringIntMapToString(storedTimeStamp)}) //200
}

/*
func forwardPut(clusterID int, key string, value string, causalPayload int, timestamp int){
	//try requesting random replicas in cluster
	noResp := true
	dataCluster := getPartition(clusterID)
	for noResp {
		if len(dataCluster) <= 0{
			json.NewEncoder(w).Encode(&map[string]string{"result": "success", "kvs": kvStr, 
				"causal_payload": stringIntMapToString(vClock), "timestamp": stringIntMapToString(storedTimeStamp)}) //200
		}
	}

}

def forwardPut(cluster, key, value, causalPayload, timestamp):
    #Try requesting random replicas
    noResp = True
    dataCluster = getPartition(cluster)
    while noResp:
        if dataCluster is None:
            return {'result': 'error', 'msg': 'Server unavailable'}, 500
        repIp = random.choice(dataCluster)
        try:
            response = requests.put((http_str + repIp + kv_str + key), 
                data = {'val': value, 'causal_payload': causalPayload, 'timestamp': timestamp})
        except requests.exceptions.RequestException as exc: #Handle replica failure
            dataCluster.remove(repIp)
            removeReplica(repIp)
            notInView.append(repIp)
            notInView = sortIPs(notInView)
            continue
        noResp = False
    #_print(response.json(), 'Fp')
return response.json()
*/