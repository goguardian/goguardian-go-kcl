#!/bin/bash

counter=0

while true; do
    aws kinesis put-record --region us-east-1 --stream-name sample_kinesis_stream --cli-binary-format raw-in-base64-out --data $RANDOM --partition-key $RANDOM > /dev/null
    ((counter++))
    echo "Records put: $counter"
    sleep 0.5
done