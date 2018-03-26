//main file for running node in kvs
//written by Vikram Melkote
//note: remember, when file is saved in vs code, code is auto formatted and unused
// code is removed. Make sure to comment out unused code or it won't persist
//'cluster' and 'partition' refer to the same thing

package main

import (
	"time"
	// "encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"ipsorting"

	"github.com/gorilla/mux" //import router
)

//node's vars
var keyVals map[string]string      //map (dictionary) of string:string
var ipPort string                  //node's own "<IP:Port>"
var view []string                  //node's view, initially passed in through env
var notInView []string             //nodes not in view, pinged in heartbeat in case they're back
var nodesPerCluster int            //number of nodes in each cluster ('K' from env)
var clusterID int                  //cluster index node belongs to
var isReplica bool                 //node type of this node
var vClock map[string]int          //maps keys to logical clock values
var storedTimeStamp map[string]int //maps keys to latest timestamp
var hashRing []string              //sorted ring of hashes
var hashClusterMap map[string]int  //maps buckets to cluster indexes
var cDict map[string]int           //dictionary mapping node IPs to cluster indicies
var localCluster []string          //list of IPs in the local cluster (could be proxy cluster)
var proxies []string               //list of proxy IPs in the network

//set up request properties
var requestTimeout = time.Duration(175 * time.Millisecond)
var httpClient = http.Client{
	Timeout: requestTimeout,
}

//Strings to prepend onto URL.
var httpStr = "http://"
var kvStr = "/kv-store/"
var newline = "\n"

//init
func main() {
	router := mux.NewRouter()              //init router
	keyVals = make(map[string]string, 100) //initialize with 100 keys

	//get environmental info
	ipPort = os.Getenv("IPPORT") //get node's ipPort from env
	var err error
	nodesPerCluster, err = strconv.Atoi(os.Getenv("K"))
	if err != nil {
		log.Println(err)
	}

	//initialize views, all in 'view' env go to notInView because we don't know if they are up,
	//we'll check for them in the first heartbeat
	notInView = strings.Split(os.Getenv("VIEW"), ",") //get node's initial view from env
	notInView = ipsorting.SortIPs(notInView)
	//initially we're only sure that the node itself is online
	view = append(view, ipPort)

	//set initial vars
	clusterID = 0
	isReplica = false //initialize as proxy

	sortedView := ipsorting.SortIPs(view)
	log.Println(sortedView)

	heartBeat() //start heartBeat()

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

	log.Fatal(http.ListenAndServe(ipPort, router))
}
