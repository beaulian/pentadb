#!/bin/sh

port=${PENTADB_PORT:-'4567'}
path=${PENTADB_PATH:-'/tmp/pentadb'}

base_dir="/pentadb/src/"

packages=(
    "github.com/golang/snappy",
    "github.com/satori/go.uuid",
    "github.com/syndtr/goleveldb/leveldb"
)

# delete invalid packages
# then download the newer packages
for package in packages
do
    rm -r ${base_dir}${package}
    go get ${package}
done

# exec
cd /pentadb/src/github.com/shenaishiren/pentadb/commands
go run server.go -p ${port} -a ${path}