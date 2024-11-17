package user

import (
	"fmt"
	"net/http"

	"router.com/repo"
)

type User struct {
	*repo.RedisClient
}

func NewUser(rc *repo.RedisClient) *User {
	return &User{rc}
}

func (u User) ValidateUser(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	password := r.Header.Get("password")

	if u.Validate(username, password) != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Authentication successful!")

}
