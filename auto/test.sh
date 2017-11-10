#!/bin/bash

if [ -z "$1" -o -z "$2" ]; then
	echo "Expected server address and number of streams got \"$1\" and \"$2\". Exiting..."
	exit 1
fi

echo "Testing server $1 with $2 streams..."
$HOME/bin/gst-client -s $1:5555 -r -n 4000000 -c $2 > $HOME/test-$2.json
