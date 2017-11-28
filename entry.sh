#!/bin/bash

port=${PENTADB_PORT:-'4567'}
path=${PENTADB_PATH:-'/tmp/pentadb'}

base_dir="/pentadb/src/"

packages=(
    github.com/golang/snappy
    github.com/satori/go.uuid
    github.com/syndtr/goleveldb
)

# delete invalid packages
# then download the newer packages
for package in ${packages[@]}
do
    rm -r ${base_dir}${package}
    go get ${package}
done

# exec
server_dir="/pentadb/src/github.com/shenaishiren/pentadb/commands"
go run ${server_dir}/server.go -p ${port} -a ${path}