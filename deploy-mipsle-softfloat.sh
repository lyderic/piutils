#!/bin/bash

REMOTE_MIPS_MACHINE=gl-inet

main() {
	if [ -z ${1} ] ; then
		usage
	fi
	gofile="${1}"
  # utility is gofile without the '.go' extension
	utility="${gofile%.*}"
	GOARCH=mipsle GOMIPS=softfloat go build -v "${gofile}" && scp "${utility}" root@${REMOTE_MIPS_MACHINE}:/mnt/mmcblk0p1/bin
	if [ $? -eq 0 ] ; then
		echo "${utility} deployed"
	fi
}

usage() {
	echo "Usage: $(basename ${0}) <utility.go>"
	exit 23
}

main $@
