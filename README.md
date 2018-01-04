#Distributed, Fault-Tolerant Key-Value Store written in Golang#
##by Vikram Melkote and Mike Hamilton###

###Set-Up###

* to make a docker subnet:

	`sudo docker network create --subnet 10.0.0.0/16 mynet`

* build the docker container:

	`sudo docker build -t mycontainer .`

###Usage###

* to start a full node example:

	`sudo docker run -p 8084:8080 --ip=10.0.0.24 --net=mynet -e K=3 -e VIEW="10.0.0.21:8080,10.0.0.22:8080,10.0.0.23:8080,10.0.0.24:8080" -e IPPORT="10.0.0.24:8080" mycontainer`

* to start a half node example:

	`sudo docker run -p 8084:8080 --ip=10.0.0.24 --net=mynet -e IPPORT="10.0.0.24:8080" mycontainer`

* add `-d` to docker run commands to run in node in detached node (console output of node not visible)

* to send a request to a node:

	**GET**

	`curl -i http://10.0.0.22:8080/kv-store -X GET`

	or

	`curl -i http://10.0.0.22:8080/kv-store/{key} -X GET`

	or

	`curl -i http://10.0.0.22:8080/kv-store/get_node_details -X GET` (not yet implemented)

	**PUT**

	`curl -i  http://localhost:1337/kv-store/kiwi -d "val=green" -X PUT`

	or

	`curl -i http://10.0.0.21:8080/kv-store/update_view?type=add -d "ip_port=10.0.0.22:8080" -X PUT` (not yet implemented)

	or

	`curl -i http://10.0.0.21:8080/kv-store/update_view?type=add -d "ip_port=10.0.0.22:8080" -X PUT` (not yet implemented)
	
	**DEL**

	`curl -i http://10.0.0.22:8080/kv-store/{key} -X DEL`

	or

	`curl -i http://10.0.0.22:8080/kv-store -X DEL`

to view kvs of a node:

* find out node's docker bridge ip with `ip addr show` (look for 'docker')

* open browser (preferrably firefox or with JSON viewer extension for other browsers)

* navigate to `http://<nodebridgeip:bridgeport>/kv-store` optionally add `/<key>`

###TO-DO (differently than old flask build):###

* Abstract helper functions into separate files for readability

* fix print functions to be one function with different tags like Tom had it

* implement proper chord where a node only stores r buckets before and after itself
	instead of all buckets

* make read repair work so it doesn't try to send a whole dict if it's bugger than max request size.
	Look into setting up a fixed sized buffer to send huge dicts in chunks

* make sorting of arrays more efficient

###Problems to debug from old flask implementation (if encountered after rebuild):###

* dataCluster is empty in getPartition()

* nodes that get added in test 6 get added to view too many times.
	maybe a problem with update_view add? or _update! or updateRatio?

* _getAllKeys! causes broken pipe (I think in test 8)

###Resources###

Golang

[Golang programming in one video](https://www.youtube.com/watch?v=CF9S4QZuV30)

[Building a RESTful API with Golang](https://www.codementor.io/codehakase/building-a-restful-api-with-golang-a6yivzqdo)

[Golang official specification](https://golang.org/ref/spec)

[Encoding/decoding JSON in Golang](https://kev.inburke.com/kevin/golang-json-http/)

[Golang HTTP status presets](https://golang.org/src/net/http/status.go)

Golang + Docker

[Using godep vendor with docker](https://stackoverflow.com/questions/40340860/godep-vendor-with-docker)

[Golang and Docker for development and production](https://medium.com/statuscode/golang-docker-for-development-and-production-ce3ad4e69673)

[How to create the smallest possible Docker image for your goland application](http://blog.cloud66.com/how-to-create-the-smallest-possible-docker-image-for-your-golang-application/)

Golang + vscode

[issue with delv for golang<=1.7](https://github.com/derekparker/delve/issues/936)

[godoc command is not available fix](https://github.com/Microsoft/vscode-go/issues/446)

* note: beware of nothing things down in your code, vscode marks extraneous code as wrong or even 

	removes it entirely upon save

