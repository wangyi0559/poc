#!/bin/bash


function installNodeModules() {
	echo
	if [ -d node_modules ]; then
		echo "============== node modules installed already ============="
	else
		echo "============== Installing node modules ============="
		npm install
	fi
	echo
}


installNodeModules

docker run -p 27017:27017 -v $PWD/db:/data/db -d mongo:3.2
sleep 5
npm start
