//functions that provide basic utilities and tools
package main

import (
	"bytes"
	"fmt"
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

func clusterMapToString(m map[string]int) string {
	mapString := new(bytes.Buffer)
	for node, clusterIndex := range m {
		fmt.Fprintf(mapString, "%s=\"%d\"\n", node, clusterIndex)
	}
	return mapString.String()
}
