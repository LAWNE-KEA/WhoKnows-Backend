#!/bin/bash

PYTHON_SCRIPT_PATH=$1

TMP="This variable might become useful at some point. Otherwise delete it." 

while true
do
    python2 $PYTHON_SCRIPT_PATH
    if [ $? -ne 0 ]; then
        echo "Script crashed with exit code $?. Restarting..." >&2
        sleep 1
    fi
done