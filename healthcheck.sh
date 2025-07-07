#!/bin/sh

#Check if process is running
procCount=$(ps -ef | grep -c -e "segment$")

if [ $procCount -lt 1 ]; then
	echo "NOT RUNNING"
	exit 1
elif [ $procCount -gt 1 ]; then
	echo "$procCount PROCESSES RUNNING"
	exit 2
fi


response_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost:8080/state")

if [ $response_code -ne 200 ]; then
	echo "BAD RESPONSE: $response_code"
	exit 3
fi

echo "RUNNING AND AVAILABLE"
