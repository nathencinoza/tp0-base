#!/bin/bash

SERVER_CONTAINER_NAME="server"
NETWORK_NAME="tp0_testing_net"
TEST_MESSAGE="Testing! :)"

RESPONSE=$(docker run --rm --network $NETWORK_NAME busybox sh -c "echo '$TEST_MESSAGE' | nc $SERVER_CONTAINER_NAME 12345")

if [ "$RESPONSE" = "$TEST_MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi
