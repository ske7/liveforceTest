package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
	"os"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("USER_SESSION_SECRET")))

type userCredentials struct {
	Email    string
	Password string
}

// LoginHandler check the user credentials
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var userData userCredentials

	if r.Body == nil {
		http.Error(w, "Please send a request JSON", http.StatusBadRequest)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&userData)
	if err != nil && err != io.EOF {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := userData.Email
	password := []byte(userData.Password)

	uid, originalPassword, err := GetPasswordByEmail(email)
	if err != nil {
		if err.Error() == "login: non-existent user" {
			http.Error(w, "User is not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	session, _ := store.Get(r, "session.id")
	if CheckPassword(originalPassword, password) {
		session.Values["authenticated"] = true
		session.Values["uid"] = uid
		session.Save(r, w)
	} else {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	w.Write([]byte("Logged In successfully"))
}

// LogoutHandler removes the user session
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session.id")
	session.Values["authenticated"] = false
	session.Save(r, w)
	w.Write([]byte("Log Out done"))
}

// Helper function for checking auth and getting logged user role
func checkUserAuthAndReturnRole(w http.ResponseWriter, r *http.Request) (role string, err error) {
	session, _ := store.Get(r, "session.id")
	if (session.Values["authenticated"] != nil) && session.Values["authenticated"] != false {
		uid := session.Values["uid"].(int)
		role, err = GetRole(uid)
		if err != nil {
			if err.Error() == "login: non-existent user" {
				return "", errors.New("User is not found")
			}
			return "", err
		}
	} else {
		return "", errors.New("User is not allowed the action. Please login")
	}
	return role, nil
}

// AddUserHandler add user to DB table
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	if role, err := checkUserAuthAndReturnRole(w, r); err == nil {
		var addUserData User

		if r.Body == nil {
			http.Error(w, "Please send a request JSON", http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&addUserData)
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if addUserData.Role == "admin" && role != "admin" {
			http.Error(w, "Add user with role 'admin' is forbidden for non admins", http.StatusBadRequest)
			return
		}

		err = addUserData.Add()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("User added successfully!"))

	} else {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}

// DeleteUserHandler delete user from DB table
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if role, err := checkUserAuthAndReturnRole(w, r); err == nil {
		if role == "user" {
			http.Error(w, "Delete action is forbidden for 'user' role", http.StatusForbidden)
			return
		}
		var delUserData User

		if r.Body == nil {
			http.Error(w, "Please send a request JSON", http.StatusBadRequest)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&delUserData)
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = delUserData.Del()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Write([]byte("User deleted successfully!"))

	} else {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}

// GetUsersHandler prints users table
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if role, err := checkUserAuthAndReturnRole(w, r); err == nil {
		users, err := GetUsers(role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i := range users {
			fmt.Fprintf(w, "%v\t%v\t%v\t%v\t%v\n", users[i].ID, users[i].Name, users[i].Email, users[i].DateCreated, users[i].Role)
		}

	} else {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
}
