#!/bin/bash
response=$(curl 16.171.225.38:32063)

if [[ $(echo "$response" | wc -w) -gt 20 ]]
then
    echo "true"
    exit 0
else
    exit 1
fi