#!/bin/bash

RASPBERRY_PI=pi3

main() {
	if [ -z ${1} ] ; then
		usage
	fi
	gofile="${1}"
  # utility is gofile without the '.go' extension
	utility="${gofile%.*}"
	GOBIN=$HOME/go/bin GOOS=linux GOARCH=arm go build -v "${gofile}" && scp "${utility}" ${RASPBERRY_PI}:bin
	if [ $? -eq 0 ] ; then
		echo "${utility} deployed"
	fi
	#rm "${utility}"
}

usage() {
	echo "Usage: $(basename ${0}) <utility.go>"
	exit 23
}

main $@
