#!/bin/bash

port=${PENTADB_PORT:-'4567'}
path=${PENTADB_PATH:-'/var/lib/pentadbs'}

# exec
server_dir="/pentadb/server"
go run ${server_dir}/server.go ${server_dir}/server_node.go -p ${port} -a ${path}