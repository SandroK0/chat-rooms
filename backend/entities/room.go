package entities

import (
	"fmt"
)

type Room struct {
	Name  string
	Users []*User
}

func NewRoom(roomName string) *Room {
	return &Room{Name: roomName}
}

func (r *Room) AddUser(user *User) error {
	for _, u := range r.Users {
		if u.Name == user.Name {
			return fmt.Errorf("user %q already exists in room %q", user.Name, r.Name)
		}
	}
	r.Users = append(r.Users, user)
	return nil
}

func (r *Room) GetUserByToken(token string) *User {
	for _, u := range r.Users {
		if u.Token == token {
			return u
		}
	}
	return nil
}

func (r *Room) RemoveUser(token string) {
	for i, u := range r.Users {
		if u.Token == token {
			r.Users = append(r.Users[:i], r.Users[i+1:]...)
			return
		}
	}
}

func (r *Room) IsEmpty() bool {
	return len(r.Users) == 0
}
