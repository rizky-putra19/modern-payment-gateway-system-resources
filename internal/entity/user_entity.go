package entity

import "time"

type RolesEntity struct {
	Id     int    `db:"id" json:"id"`
	Roles  string `db:"role_name" json:"roles"`
	Access string `db:"permission_desc" json:"access"`
}

type ListUsersEntity struct {
	Id        int       `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Status    string    `db:"status" json:"status"`
	Roles     string    `db:"role_name" json:"roles"`
}
