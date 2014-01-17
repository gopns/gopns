package model

type User struct {
	id string
}

func NewUser(id string) (user *User) {
	return &User{id}
}

func (user User) Id() string {
	return user.id
}

func (user *User) SetId(id string) {
	user.id = id
}
