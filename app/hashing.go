//functions pertaining to consistent hashing
package main

import (
	"encoding/hex"
	"crypto/md5"
	"strconv"
)

//update hash values for clusters
func updateHashRing() {
	updateCDict()
	
	//make a set of unique vals (which are clusterIDs) from cDict to 
	//generate buckets for
	var clusterIDs []int
    for _, clusterInd := range cDict{
        if clusterInd == int(len(view)/nodesPerCluster){
			continue
		}
		clusterIDs = append(clusterIDs, clusterInd)
	}

	//traverse through set to generate 250 buckets for each clusters

	hasher := md5.New()
	//clear hashRing and hashClusterMap
	hashRing = hashRing[:0]
	hashClusterMap = make(map[string]int)

    for cIndex := range clusterIDs {
        for j := 0; j < 250; j++ { //generate 250 buckets per cluster
			hasher.Write([]byte("hash string" + strconv.Itoa(cIndex) + strconv.Itoa(j) + newline))
			hashStr := hex.EncodeToString(hasher.Sum(nil))
			hashRing = insortStringIntoSlice(hashRing, hashStr)
			hashClusterMap[hashStr] = cIndex
		}
	}
}

//checks which cluster a key belongs to
func checkKeyHash(key string) int {
	hasher := md5.New()

	hasher.Write([]byte(key + newline))
	ringIndex := -1 //where the hash of the key would go in the hash ring
	
	for ind, hash := range hashRing {
		if key > hash {
			continue
		}
		ringIndex = ind
	}
	if ringIndex < 0 {
		ringIndex = 0
	}
	ringHash := hashRing[ringIndex] //get corresponding hash
	return hashClusterMap[ringHash] //return clusterID corresponsing to hash
}