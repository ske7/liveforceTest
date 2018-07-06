package main

import (
	"database/sql"
	"errors"
	"time"
)

// User struct
type User struct {
	Name     string
	Email    string
	Password string
	Role     string
}

// UserToPrint struct
type UserToPrint struct {
	ID          int
	Name        string
	Email       string
	DateCreated *time.Time
	Role        string
}

// GetPasswordByEmail return password hash from DB
func GetPasswordByEmail(email string) (int, string, error) {
	var password string
	var id int
	if err := GetDB().QueryRow(
		"SELECT id, password FROM public.users WHERE email = $1", email).Scan(&id, &password); err == sql.ErrNoRows {
		return 0, "", errors.New("login: non-existent user")
	} else if err != nil {
		return 0, "", err
	}
	return id, password, nil
}

// GetRole return role of the user from DB
func GetRole(uid int) (string, error) {
	var role string
	if err := GetDB().QueryRow(
		"SELECT role FROM public.users WHERE id = $1", uid).Scan(&role); err == sql.ErrNoRows {
		return "", errors.New("login: non-existent user")
	} else if err != nil {
		return "", err
	}
	return role, nil
}

// GetUsers return all users from DB
func GetUsers(role string) (users []UserToPrint, err error) {
	var sql string
	if role == "user" {
		sql = `SELECT id, name, email, date_created, role FROM public.users where role <> 'admin'`
	} else {
		sql = `SELECT id, name, email, date_created, role FROM public.users`
	}
	rows, err := GetDB().Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id          int
			name        string
			email       string
			dateCreated *time.Time
			role        string
		)
		if err = rows.Scan(&id, &name, &email, &dateCreated, &role); err != nil {
			return nil, err
		}
		users = append(users, UserToPrint{id, name, email, dateCreated, role})

	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// Add Saves the User to DB
func (a *User) Add() error {
	sql := `INSERT INTO public.users (name, email, password, date_created, role)
          VALUES ($1, $2, $3, CURRENT_TIMESTAMP, $4)`

	_, err := GetDB().Exec(sql, a.Name, a.Email, HashPassword([]byte(a.Password)), a.Role)
	return err
}

// Del delete the User from DB
func (a *User) Del() error {
	sql := `delete from public.users where email = $1`

	_, err := GetDB().Exec(sql, a.Email)
	return err
}
