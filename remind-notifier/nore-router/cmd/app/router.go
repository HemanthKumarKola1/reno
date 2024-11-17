package main

import (
	"net/http"

	"router.com/pkg/user"
	"router.com/repo"
)

func route(rc *repo.RedisClient) *http.ServeMux {

	mux := http.NewServeMux()
	u := user.NewUser(rc)
	mux.HandleFunc("/server1", u.ValidateUser)

	return mux
}
