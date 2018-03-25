//functions relating to updating variables or view components of each node

package main

import (
	"log"
	"strings"

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

//heartbeat function to ping nodes and adjust the view
func heartBeat() {
	// heart = threading.Timer(3.0, heartBeat)
	// heart.daemon = True
	// heart.start()

	// if firstHeartBeat:
	//     firstHeartBeat = False
	//     return

	log.Println("My IP: " + ipPort + newline +
		"View: " + strings.Join(view, ", ") + newline +
		"Cluster Map: " + clusterMapToString(cDict) + newline +
		"Proxies: " + strings.Join(proxies, ", "))

	/*
		    for ip in notInView: //check if any nodes not currently in view came back online
		        try:
		            response = (requests.get((http_str + ip + kv_str + "get_node_details"), timeout=2)).json()
		            if response['result'] == 'success':
		                notInView.remove(ip)
		                view.append(ip)
		                view = sortIPs(view)
		                #updateHashRing()
		        except: #Handle no response from i
		            pass
		    for ip in view:
		        if ip != IpPort:
		            try:
		                response = (requests.get((http_str + ip + kv_str + "get_node_details"), timeout=2)).json()
		            except requests.exceptions.RequestException as exc: #Handle no response from ip
		                if ip in replicas:
		                    removeReplica(ip)
		                elif ip in proxies:
		                    removeProxie(ip)
		                notInView.append(ip)
		                notInView = sortIPs(notInView)
		                #updateHashRing()
		    updateRatio()
		    print("reps " + str(replicas))
		    print("prox " + str(proxies))
			sys.stdout.flush()
	*/
}
