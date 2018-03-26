//functions relating to updating variables or view components of each node

package main

import (
	"log"
	"strings"
	"time"

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

//removes IP from localCluster, view, notInView, cDict, proxies of current node
//also removes from notInView if node did not crash (ie deliberately removed) so it
//isn't added back by heartbeat
func removeNode(ip string, crash bool) {
	if stringInSlice(ip, localCluster) {
		remove(ip, localCluster)
		localCluster = ipsorting.SortIPs(localCluster)
	}

	if stringInSlice(ip, proxies) {
		remove(ip, proxies)
		proxies = ipsorting.SortIPs(proxies)
	}

	if stringInSlice(ip, view) {
		remove(ip, view)
		view = ipsorting.SortIPs(view)
	}

	_, nodeInCDict := cDict[ip]
	if nodeInCDict {
		delete(cDict, ip)
	}

	if crash && stringInSlice(ip, notInView) == false {
		notInView = append(notInView, ip)
		notInView = ipsorting.SortIPs(notInView)
	}

	if crash == false && stringInSlice(ip, notInView) {
		remove(ip, notInView)
		notInView = ipsorting.SortIPs(notInView)
	}
}

//heartbeat function to ping nodes and adjust the view
func heartBeat() {
	for {
		// heart = threading.Timer(3.0, heartBeat)
		// heart.daemon = True
		// heart.start()

		// if firstHeartBeat:
		//     firstHeartBeat = False
		//     return

		time.Sleep(2 * time.Second) //run heartbeat every 2s

		log.Println("My IP: " + ipPort + newline +
			"View: " + strings.Join(view, ", ") + newline +
			"Cluster Map: " + clusterMapToString(cDict) + newline +
			"Proxies: " + strings.Join(proxies, ", "))

		//loop through notInView to see if any nodes in there have come online
		for _, node := range notInView {
			alive := pingNode(node)
			if alive {
				if stringInSlice(node, view) == false {
					view = append(view, node)
				}
				remove(node, notInView)
			}
		}

		//loop through view to see if any nodes in there have now gone offline
		for _, node := range view {
			if node == ipPort { //don't ping self
				continue
			}

			alive := pingNode(node)
			if alive == false {
				removeNode(node, true)
			}
		}

		/* old python code to translate/improve
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
}
