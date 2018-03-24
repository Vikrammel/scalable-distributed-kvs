//main file for running node in kvs
//written by Vikram
//note: remember, when file is saved in vs code, code is auto formatted and unused
// code is removed. Make sure to comment out unused code or it won't persist

package main

import (
	// "encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"strconv"

	"ipsorting"

	"github.com/gorilla/mux" //import router
)

//node's vars
var keyVals map[string]string //map (dictionary) of string:string
var ipport string             //node's own "<IP:Port>"
var view []string             //node's initial view passed in through env
var nodesPerCluster int		  //number of nodes in each cluster ('K' from env)
var partitionID int			  //partition index node belongs to

//init
func main() {
	router := mux.NewRouter()                    //init router
	keyVals = make(map[string]string, 100)       //initialize with 100 keys

	//get environmental info
	ipport = os.Getenv("IPPORT")                 //get node's ipport from env
	var err error
	nodesPerCluster, err = strconv.Atoi(os.Getenv("K"))
	if err != nil {
		log.Println(err)
	}
	view = strings.Split(os.Getenv("VIEW"), ",") //get node's initial view from env

	//set vars
	partitionID = 0
	sortedView := ipsorting.SortIPs(view)
	log.Println(sortedView)

	//funcs for routes (with and without slashes at the end of URL)
	router.HandleFunc("/kv-store", GetAllKeys).Methods("GET")
	router.HandleFunc("/kv-store/", GetAllKeys).Methods("GET")
	router.HandleFunc("/kv-store/{key}", GetKey).Methods("GET")
	router.HandleFunc("/kv-store/{key}/", GetKey).Methods("GET")
	router.HandleFunc("/kv-store/{key}", PutKey).Methods("PUT")
	router.HandleFunc("/kv-store/{key}/", PutKey).Methods("PUT")
	router.HandleFunc("/kv-store/{key}", DeleteKey).Methods("DEL")
	router.HandleFunc("/kv-store/{key}/", DeleteKey).Methods("DEL")
	router.HandleFunc("/kv-store", DeleteAll).Methods("DEL")
	router.HandleFunc("/kv-store/", DeleteAll).Methods("DEL")

	log.Fatal(http.ListenAndServe(ipport, router))
}
