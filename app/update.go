//functions relating to updating variables or view components of each node

package main

import (
	"ipsorting"
)

func updateCDict(){
    //sort view and update cDict off of view
    view = ipsorting.SortIPs(view)
    for index,node := range view{
        clusterID = index/nodesPerCluster
		cDict[node] = clusterID
	}
    for node,_ := range cDict{
        if stringInSlice(node, view){
			delete(cDict, node)
		}
	}
}