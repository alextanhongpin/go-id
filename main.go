package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/alextanhongpin/go-id/ids"

	_ "github.com/lib/pq"
)

const (
	appID  = "app"
	userID = "user"
)

type AppID string

func (a AppID) Prefix() string {
	return appID
}

type UserID string

func (u UserID) Prefix() string {
	return userID
}

type User struct {
	ID   ids.ID[UserID]
	Name string
}

func main() {
	db, err := sql.Open("postgres", "postgres://john:123456@127.0.0.1:5432/test?sslmode=disable")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var u User
	if err := db.QueryRow(`select id, name from users`).Scan(&u.ID, &u.Name); err != nil {
		panic(err)
	}
	fmt.Println("got user", u)

	var u2 User
	if err := db.QueryRow(`select id, name from users where id = $1`, u.ID).Scan(&u2.ID, &u2.Name); err != nil {
		panic(err)
	}
	fmt.Println("got user by name", u2)

	u3 := User{
		ID:   ids.New[UserID](),
		Name: "alice",
	}
	fmt.Println("inserting", u3)
	_, err = db.Exec(`insert into users (id, name) values ($1, $2)`, u3.ID, u3.Name)
	if err != nil {
		panic(err)
	}

	var id = ids.New[AppID]()
	b, err := json.MarshalIndent(appRequest{
		ID: id,
	}, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
	var req userRequest
	if err := json.Unmarshal(b, &req); err != nil {
		panic(err)
	}

	fmt.Println(req)
	fmt.Println(req.ID.Prefix())
	fmt.Println(req.ID.UUID())
}

type appRequest struct {
	ID ids.ID[AppID] `json:"id"`
}

type userRequest struct {
	ID ids.ID[UserID] `json:"id"`
}
