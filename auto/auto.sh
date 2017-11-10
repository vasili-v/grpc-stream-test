#!/bin/bash

if [ -z "$1" -o -z "$2" ]; then
	echo "Expected server and client addresses got \"$1\" and \"$2\". Exiting..."
	exit 1
fi

echo "Copying script to \"$1\"..."
scp prepare.sh $1:

echo "Copying script to \"$2\"..."
scp prepare.sh $2:

echo "Preparing environment on \"$1\"..."
ssh $1 ./prepare.sh

echo "Preparing environment on \"$2\"..."
ssh $2 ./prepare.sh

echo "Installing and starting server on \"$1\"..."
ssh -tt $1 src/github.com/vasili-v/grpc-stream-test/auto/server.sh &
SPID=$!
echo "Running server shell with PID=$SPID..."
trap "echo \"Killing server shell with PID=$SPID...\" && kill $SPID" EXIT

echo "Installing client on \"$2\"..."
ssh $2 src/github.com/vasili-v/grpc-stream-test/auto/client.sh

echo "Cooldown for 5s..."
sleep 5

for i in 32 64 96 128 192 256 384 512 1024 ; do
	echo "Running $i streams test on \"$2\" (using \"$1\" as a server)..."
	ssh $2 src/github.com/vasili-v/grpc-stream-test/auto/test.sh $1 $i
done

echo "Making summary for all tests..."
ssh $2 'src/github.com/vasili-v/grpc-stream-test/summary.py -csv $HOME'
