#!/bin/bash

port=${PENTADB_PORT:-'4567'}
path=${PENTADB_PATH:-'/tmp/pentadb'}

# exec
server_dir="/pentadb/src/github.com/shenaishiren/pentadb/commands"
go run ${server_dir}/server.go -p ${port} -a ${path}