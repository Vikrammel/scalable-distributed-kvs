#!/bin/bash
#to run, use 'chmod +x single-site_test.sh' then './single-site_test.sh'

echo ""
echo "========== Single-Site Key-Value Store Test =========="
echo ""
echo "Test script checks putting/getting/deleting entries from the KVS"
echo ""
read -p "Make sure to build image from project root directory as 'mycontainer' before continuing. Continue? (y/n)" CONT
echo ""
if [ "$CONT" = "y" ]; then
    (
    echo ""
    echo ""
    echo "========== Single-Site KVS Test Output =========="
    echo ""
    echo ""
    echo ""
    echo "Stopping all containers ..."
    echo ""
    docker kill $(docker ps -q)
    echo ""
    echo ""
    echo "Creating docker network 'mynet' ..."
    echo ""
    sudo docker network create --subnet 10.0.0.0/16 mynet
    echo ""
    echo ""
    echo "Starting node with IP:Port 10.0.0.24:8080 ..."
    echo ""
    docker run -p 8084:8080 --ip=10.0.0.24 --net=mynet -e K=3 -e VIEW="10.0.0.21:8080,10.0.0.22:8080,10.0.0.23:8080,10.0.0.24:8080" -e IPPORT="10.0.0.24:8080" mycontainer &
    echo ""
    echo ""
    sleep 15s
    echo "======== PUT 5 {'<fruit>':'<color>'} pairs in KVS ========"
    curl -i http://10.0.0.24:8080/kv-store/kiwi -d "val=green" -X PUT
    echo ""
    echo ""
    curl -i http://10.0.0.24:8080/kv-store/apple -d "val=red" -X PUT
    echo ""
    echo ""
    curl -i http://10.0.0.24:8080/kv-store/banana -d "val=yellow" -X PUT
    echo ""
    echo ""
    curl -i http://10.0.0.24:8080/kv-store/grape -d "val=purple" -X PUT
    echo ""
    echo ""
    curl -i http://10.0.0.24:8080/kv-store/papaya -d "val=orange" -X PUT
    echo ""
    echo ""
    echo "======== GET key-value pair for key 'banana' from KVS ========"
    curl -i http://10.0.0.24:8080/kv-store/banana -X GET
    echo ""
    echo ""
    echo "======== GET all key-value pairs from KVS ========"
    curl -i http://10.0.0.24:8080/kv-store -X GET
    echo ""
    echo ""
    echo "======== DEL key-value pair for key 'grape' from KVS ========"
    curl -i http://10.0.0.24:8080/kv-store/grape -X DEL
    echo ""
    echo ""
    echo "======== GET key-value pair for key 'grape' from KVS ========"
    curl -i http://10.0.0.24:8080/kv-store/grape -X GET
    echo ""
    echo ""
    echo "======== GET all key-value pairs from KVS ========"
    curl -i http://10.0.0.24:8080/kv-store -X GET
    echo ""
    echo ""
    echo ""
    echo "Stopping all containers ..."
    echo ""
    docker kill $(docker ps -q)
    echo ""
    echo ""
    echo ""
    echo "========== End of Single-Site KVS Tests =========="
    )>single-site_test_output.txt
fi