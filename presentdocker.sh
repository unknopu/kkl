#!/bin/bash

if (! docker stats --no-stream >/dev/null 2>&1); then
    echo "[!] docker desktop does not open yet"
    exit 1
fi

MY_NETWORK="mynet"
MY_MONGO="mongodb"
MY_REDIS="redis"
MY_API="testapi"

if docker network ls | grep -q $MY_NETWORK; then
    echo "[!] found network: $MY_NETWORK"
else
    echo "[*] initiate network: mynet"
    docker network create $MY_NETWORK
fi

if [ "$(docker ps -aq -f name=$MY_MONGO)" ]; then
    echo "[!] found mongodb://$MY_MONGO:27017"
else
    echo "[*] initiate mongodb: 'mongodb://mongodb:27017'"
    docker run --name $MY_MONGO -it -p 27017:27017 --network $MY_NETWORK -d mongo
fi

if [ "$(docker ps -aq -f name=$MY_REDIS)" ]; then
    echo "[!] found redis"
else
    echo "[*] initiate redis"
    docker run --name $MY_REDIS -p 6379:6379 --network $MY_NETWORK -v redisdata:/data -d redis
fi

sleep 1
if [ "$(docker ps -aq -f name=$MY_API)" ]; then
    echo "[!] found API"
else
    echo "[*] initiate $MY_API..."
    docker build -t $MY_API .
    docker run --name $MY_API -p 4000:4000 --network $MY_NETWORK -d $MY_API
fi


