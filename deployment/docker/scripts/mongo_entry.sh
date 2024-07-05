#!/bin/bash


echo "rs.initiate()" > /docker-entrypoint-initdb.d/1-init-replicaset.js
echo "db = db.getSiblingDB(process.env[$0]);" > /docker-entrypoint-initdb.d/2-init-db-collection.js
echo "db.createCollection($1, { capped: false });" >> /docker-entrypoint-initdb.d/2-init-db-collection.js
echo "db.init.insert([{ message: $2 }]);" >> /docker-entrypoint-initdb.d/2-init-db-collection.js

/usr/local/bin/docker-entrypoint.sh mongod --replSet rs0 --bind_ip_all --noauth "'MONGO_APP_DATABASE'" "'init'" "'db initialized successfully'"
