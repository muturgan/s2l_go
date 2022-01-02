#!/bin/sh
SCRIPT_PATH=$(dirname $0)
export $(cat $SCRIPT_PATH/.env | xargs) && $SCRIPT_PATH/main
