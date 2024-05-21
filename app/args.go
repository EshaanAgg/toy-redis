package main

import (
	"flag"
)

type Args struct {
	port       int
	replicaof  string
	dir        string
	dbfilename string
}

func GetArgs() Args {
	port := flag.Int("port", 6379, "The port on which the Redis server listens")
	replicaof := flag.String("replicaof", "", "The host and port of the master server")
	dir := flag.String("dir", "/tmp/redis", "The directory in which to store the database files")
	dbfilename := flag.String("dbfilename", "dump.rdb", "The name of the database file")

	flag.Parse()

	return Args{
		port:       *port,
		replicaof:  *replicaof,
		dir:        *dir,
		dbfilename: *dbfilename,
	}
}
