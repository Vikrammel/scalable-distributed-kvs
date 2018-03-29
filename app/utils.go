//functions that provide basic utilities and tools
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
)

//checks if string is in slice
//https://stackoverflow.com/a/15323988
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

//gets index of string in string slice
func indexOf(str string, strSlice []string) int {
	for ind, val := range strSlice {
		if val == str {
			return ind
		}
	}
	return -1 //not in slice
}

//removes string from slice
//https://stackoverflow.com/a/34070691
func remove(r string, s []string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

//converts the map of node IPs:clusterID to a string
func clusterMapToString(m map[string]int) string {
	mapString := new(bytes.Buffer)
	for node, clusterIndex := range m {
		fmt.Fprintf(mapString, "%s=\"%d\"\n", node, clusterIndex)
	}
	return mapString.String()
}

//pings input node to check if it's online, returns true if it is, else false
func pingNode(node string) bool {
	// Make request for node's clusterID
	log.Println("pinging " + node)
	rs, err := httpClient.Get(httpStr + node + kvStr + "get_partition_id")
	// Error handling
	if err == nil {

		defer rs.Body.Close()

		_, err := ioutil.ReadAll(rs.Body)
		if err == nil {
			log.Println(node + " response parsed successfully")
			return true //node is online, response is good
		}
		log.Println(node + " response parsing failed")
		return false //bad response
		// bodyString := string(bodyBytes)
	}
	log.Println(node + " request error: " + err.Error())
	return false //error pinging node, consider offline
}
