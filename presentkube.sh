#!/bin/bash

echo "[*] initiate deployment/service redis: mynet"
kubectl create deployment redis --image=redis
kubectl expose deployment/redis --port 6379 --name redis --type NodePort

echo "[*] initiate deployment/service mongodb"
kubectl create deployment mongo --image=mongo
kubectl expose deployment/mongodb --port 27017 --name mongodb --type NodePort
# kubectl port-forward deployment/mongo 27017:27017

echo "[*] initiate deployment/service testapi"
kubectl create deployment testapi --image=unknopu/testapi
kubectl expose deployment/testapi --port 4000 --name testapi --type LoadBalancer