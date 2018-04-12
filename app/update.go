//functions relating to updating variables or view components of each node

package main

import (
	"log"
	"strings"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"reflect"
	// "bytes"
	// "net/http"

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
		if stringInSlice(node, view) == false{
			delete(cDict, node)
		}
	}
}

//removes IP from localCluster, view, notInView, cDict, proxies of current node
//also removes from notInView if node did not crash (ie deliberately removed) so it
//isn't added back by heartbeat
func removeNode(ip string, crash bool) {
	if stringInSlice(ip, localCluster) {
		localCluster = remove(ip, localCluster)
		localCluster = ipsorting.SortIPs(localCluster)
	}

	if stringInSlice(ip, proxies) {
		proxies = remove(ip, proxies)
		proxies = ipsorting.SortIPs(proxies)
	}

	if stringInSlice(ip, view) {
		view = remove(ip, view)
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
		notInView = remove(ip, notInView)
		notInView = ipsorting.SortIPs(notInView)
	}
}

func updateDatabase(){
	for _,node := range localCluster {
		if node == ipPort {
			continue
		}
		rs, err := http.Get(httpStr + node + kvStr + "_getDB")
		if err != nil {
			continue
		}
		defer rs.Body.Close()

		bodyBytes, err := ioutil.ReadAll(rs.Body)
		if err != nil {
			continue
		}
		msgMap := make(map[string]interface{})
		err2 := json.Unmarshal(bodyBytes, &msgMap)
		if err2 != nil {
			continue
		}
		log.Println("In UpdateDatabase(): ")
		log.Println(msgMap["kvs"])
	}
}
/*
def updateDatabase():
    global replicas, notInView
    for ip in replicas:
        if ip == IpPort:
            continue
        try:
            #TODO: create _getAll! function, returning [d, vClock, storedTimeStamp]
            response = (requests.get((http_str + ip + kv_str + '_getAllKeys!'), timeout=5)).json()
            try:
                responseD = json.loads(response['dict'])
            except:
                _print("Can't get data from a halfNode")
                continue
            responseCausal = json.loads(response['causal_payload'])
            responseTime = json.loads(response['timestamp'])
            for key in json.loads(response['dict']):
                if (d.get(key) == None or responseCausal[key] > vClock[key] or
                   (responseCausal[key] == vClock[key] and responseTime[key] > storedTimeStamp[key])):
                    d[key] = responseD[key].encode('ascii', 'ignore')
                    vClock[key] = responseCausal[key]
                    storedTimeStamp[key] = responseTime[key]
        except requests.exceptions.RequestException: #Handle no response from ip
            _print("updateDatabase timeout occured.")
            removeReplica(ip)
            notInView.append(ip)
			notInView = sortIPs(notInView)
*/

//update ratio of proxies / replicas
func updateRatio() {

	newPartition := indexOf(ipPort, view) / nodesPerCluster
	proxyPartition := len(view) / nodesPerCluster

	if len(view) >= nodesPerCluster {
		// log.Println( "partition: " + strconv.Itoa(clusterID) + newline +
		// 		"newPartition: " + strconv.Itoa(newPartition) + newline +
		// 		"proxy partition: " + strconv.Itoa(proxyPartition) )

		// The partition this node belongs to had changed.
		if clusterID != newPartition {
			clusterID = newPartition
			localCluster = localCluster[:0]
			localCluster = getPartition(clusterID)
		}
		if clusterID >= proxyPartition { //this node is a proxy
			if stringInSlice(ipPort, proxies) == false {

				// if len(keyVals) != 0 {
				// 	// for key, val := range keyVals {
				// 		// 	forwardPut(0, key, val, vClock[key], storedTimeStamp[key])
				// 		// response = requests.put((http_str + repIp + kv_str + key),
				// 		// data = {'val': value, 'causal_payload': causalPayload, 'timestamp': timestamp})
				// 	// }
				// }

				keyVals = make(map[string]string)
				vClock = make(map[string]int)
				storedTimeStamp = make(map[string]int)

				isReplica = false
				//_print("Converted to " + getReplicaDetail() + " at 0.5")
				proxies = append(proxies, ipPort)
				proxies = ipsorting.SortIPs(proxies)
				localCluster = localCluster[:0]
			}
		} else {
			isReplica = true
			//_print("Converted to " + str(getReplicaDetail()) + " at 1")
		}
	}

	for index, node := range view {
		// This is a proxy.
		if index/nodesPerCluster >= proxyPartition {
			if stringInSlice(node, proxies) == false {
				if stringInSlice(node, localCluster) && isReplica {
					localCluster = remove(node, localCluster)
				}
				proxies = append(proxies, node)
				proxies = ipsorting.SortIPs(proxies)
				if node == ipPort {
					isReplica = false
					//_print("Converted to " + str(getReplicaDetail()) + " at 2")
				}
			}
		} else if indexOf(node, view)/nodesPerCluster == clusterID {
			// This is a replica within the same partition.
			if node == ipPort {
				isReplica = true
				//_print("Converted to " + str(getReplicaDetail()) + " at 3")

				if stringInSlice(node, localCluster) == false {
					localCluster = append(localCluster, node)
					localCluster = ipsorting.SortIPs(localCluster)
					if stringInSlice(node, proxies) {
						proxies = remove(node, proxies)
					}
				}
				updateDatabase()
			}
			if isReplica && stringInSlice(node, localCluster) == false {
				if stringInSlice(node, proxies) {
					proxies = remove(node, proxies)
				}
				localCluster = append(localCluster, node)
				localCluster = ipsorting.SortIPs(localCluster)
			}
		} else {
			// This is a replica in another partition.
			if stringInSlice(node, proxies) {
				proxies = remove(node, proxies)
			}
			if stringInSlice(node, localCluster) {
				localCluster = remove(node, localCluster)
			}
		}
	}
	// updateHashRing()
	updateCDict() //in old flask code, called from within updateHashRing
}

//heartbeat function to ping nodes and adjust the view
func heartBeat() {
	for {
		time.Sleep(2 * time.Second) //run heartbeat every 2s

		log.Println("My IP: " + ipPort + newline +
			"View: " + strings.Join(view, ", ") + newline +
			"notInView: " + strings.Join(notInView, ", ") + newline +
			"Cluster Map: " + stringIntMapToString(cDict) + newline +
			"Proxies: " + strings.Join(proxies, ", "))

		//loop through notInView to see if any nodes in there have come online
		for _, node := range notInView {
			alive := pingNode(node)
			nodeStatus := " dead"
			if alive {
				nodeStatus = " alive"
			}
			log.Println(node + nodeStatus)
			if alive {
				if stringInSlice(node, view) == false {
					view = append(view, node)
				}
				notInView = remove(node, notInView)
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

		//sort arrays since we might have removed from / added to them
		notInView = ipsorting.SortIPs(notInView)
		view = ipsorting.SortIPs(view)

		updateRatio()
	}
}
