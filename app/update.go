//functions relating to updating variables or view components of each node

package main

import (
	"ipsorting"
)

//sort view and update cDict off of view
func updateCDict() {
	view = ipsorting.SortIPs(view)
	for index, node := range view {
		clusterID = index / nodesPerCluster
		cDict[node] = clusterID
	}
	for node := range cDict {
		if stringInSlice(node, view) {
			delete(cDict, node)
		}
	}
}

//removes replica IP from localCluster and view of current node
func removeReplica(ip string) {
	if len(localCluster) == nodesPerCluster { //make sure localCluster is not proxy cluster

		if stringInSlice(ip, localCluster) {
			remove(ip, localCluster)
		}

		if stringInSlice(ip, view) {
			remove(ip, view)
		}
	}
}

//removes proxy IP from localCluster and view of current node
func removeProxy(ip string) {
	if len(localCluster) != nodesPerCluster { //make sure localCluster is proxy cluster

		if stringInSlice(ip, localCluster) {
			remove(ip, localCluster)
		}

		if stringInSlice(ip, proxies) {
			remove(ip, proxies)
		}

		if stringInSlice(ip, view) {
			remove(ip, view)
		}
	}
}
