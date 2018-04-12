//functions that provide basic utilities and tools
package main

import (
	"strings"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"io"
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

func intInSlice(a int, list []int) bool {
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

//converts a map of string:int to a string
func stringIntMapToString(m map[string]int) string {
	mapString := new(bytes.Buffer)
	for node, clusterIndex := range m {
		fmt.Fprintf(mapString, "%s : %d \n", node, clusterIndex)
	}
	return mapString.String()
}

//converts a map of string:string to a string
func stringMapToString(m map[string]string) string {
	mapString := new(bytes.Buffer)
	fmt.Fprintf(mapString, "{")
	index := 0
	for key, val := range m {
		fmt.Fprintf(mapString, "'%s':'%s'", key, val)
		if index < len(m) - 1 {
			fmt.Fprintf(mapString, ",")
		}
		index++
	}
	fmt.Fprintf(mapString, "}")
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

//insert a string into a slice in the right place
func insortStringIntoSlice(slice []string, str string) []string {
	var returnSlice []string
	inserted := false
	for ind, currStr := range slice {
		if strings.Compare(str, currStr) == -1 || strings.Compare(str, currStr) == 0 {
			returnSlice = append(returnSlice, slice[:ind]...)
			returnSlice = append(returnSlice, str)
			returnSlice = append(returnSlice, slice[ind:]...)
			inserted = true
		}
	}
	if inserted == false {
		copy(returnSlice, slice)
		returnSlice = append(returnSlice, str)
	}
	return returnSlice
}

//gets the ips of the nodes in a partition
func getPartition(clusterID int) []string{
	var returnSlice []string
	for node, clusterIndex := range cDict {
		if clusterIndex == clusterID {
			returnSlice = append(returnSlice, node)
		}
	}
	return returnSlice
}

//https://gist.github.com/maniankara/a10d19960293b34b608ac7ef068a3d63
func putRequest(url string, data io.Reader) int {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, url, data)
	if err != nil {
		// handle error
		return 0
	}
	_, err = client.Do(req)
	if err != nil {
		// handle error
		return 0
	}
	return 1
}

/*
def getPartition(num):
    partitionStart = num*K
    partitionEnd = partitionStart + K
    if partitionEnd > len(view):
        return None
    membersRange = range(partitionStart, partitionEnd)
    members = []
    for node in view:
        if view.index(node) in membersRange:
            members.append(node)
return members
*/