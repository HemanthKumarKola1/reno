package main

import (
	"net/http"

	"github.com/gocql/gocql"
	"notify.com/repo"
)

func main() {

	// Connect to Cassandra
	cluster := gocql.NewCluster("cassandra_host1", "cassandra_host2")
	cluster.Keyspace = "your_keyspace"
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()
	repo.NewRepo(session)

	// TODO: gRPC
	http.ListenAndServe(":8080", route())
}

func route() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/server1", Handler)

	return mux
}
