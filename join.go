package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

// login if user info is correct
func login(c echo.Context) error {
	femail := c.FormValue("email")
	fpass := c.FormValue("password")
	userid, username, email, pass := getUsername(femail)

	if pass == fpass && femail == email {
		//userSession[email] = username
		setSession(c, username, userid)
		return c.Redirect(http.StatusSeeOther, "/") // 303 code
		// TODO redirect to latest page
	}

	data := make(map[string]interface{}, 2)
	data["userid"] = nil
	data["error"] = "user information is not correct"
	return c.Render(200, "login.html", data)
}

// get an username by email
func getUsername(femail string) (int, string, string, string) {
	var name, email, password string
	var userid int
	err := db.QueryRow(
		"SELECT userid, username, email, password FROM comments.users WHERE email = ?",
		femail).Scan(&userid, &name, &email, &password)
	if err != nil {
		fmt.Println("Error with db.QueryRow", err.Error())
	}
	return userid, name, email, password
}
func insertUser(user, pass, email string) error {
	_, err := db.Exec(
		"INSERT INTO comments.users(username, password, email) VALUES ( ?, ?, ?)",
		user, pass, email)

	// if there is an error inserting, handle it
	if err != nil {
		fmt.Println(err)
		return err
	}
	// be careful deferring Queries if you are using transactions
	return nil
}

func setSession(c echo.Context, username string, userid int) {
	sess, _ := session.Get("session", c)
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60, // = 1h,
		HttpOnly: true,    // no websocket or any thing else
	}
	sess.Values["username"] = username
	sess.Values["userid"] = userid
	sess.Save(c.Request(), c.Response())
}

func signup(c echo.Context) error {
	username := c.FormValue("username")
	pass := c.FormValue("password")
	email := c.FormValue("email")
	err := insertUser(username, pass, email)
	if err != nil {
		//fmt.Println(err)
		return c.Render(200, "sign.html", "wrrone")
	}
	return c.Redirect(http.StatusSeeOther, "/login") // 303 code
}

func signPage(c echo.Context) error {
	data := make(map[string]interface{}, 1)
	sess, _ := session.Get("session", c)
	data["userid"] = sess.Values["userid"]
	data["username"] = sess.Values["username"]
	return c.Render(200, "sign.html", data)
	//fmt.Println( c.Render(200, "sign.html", sess.Values["userid"].(int))); return nil
}

func loginPage(c echo.Context) error {
	data := make(map[string]interface{}, 1)
	sess, _ := session.Get("session", c)
	data["userid"] = sess.Values["userid"]
	data["username"] = sess.Values["username"]
	return c.Render(200, "login.html", data)
	//fmt.Println( c.Render(200, "login.html", data)); return nil
}
