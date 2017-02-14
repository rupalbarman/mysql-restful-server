package main 

import (
	"database/sql"
	"net/http"
	"fmt"
	"html/template"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func homePage(w http.ResponseWriter, r *http.Request) {
	if r.Method== "GET"{
		t, _:=template.ParseFiles("static/index.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	
	if r.Method== "GET"{

		t, _:=template.ParseFiles("static/login.html")
		t.Execute(w, nil)

	} else{

		r.ParseForm()
		username:= r.FormValue("username")
		password:= r.FormValue("password")
		fmt.Println("Checking user:", username, " ", password)
		var dbUser string
		var dbPass string

		err:= db.QueryRow("SELECT username, password FROM users WHERE username=? ",username).Scan(&dbUser, &dbPass)
		if err== nil {
			fmt.Println("User authorized")
		}
		if dbUser== username && dbPass== password {
			fmt.Println("Logged in", " ", username)
		}

		http.Redirect(w, r, "/", http.StatusFound)
		
	}
}

func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method== "GET"{

		t, _:=template.ParseFiles("static/signup.html")
		t.Execute(w, nil)

	} else{

		r.ParseForm()
		username:= r.FormValue("username")
		password:= r.FormValue("password")
		fmt.Println("New User:", username, " ", password)

		var dbUser string
		var dbPass string

		err:= db.QueryRow("SELECT username, password FROM users WHERE username=? ",username).Scan(&dbUser, &dbPass)
		if err== nil {
			fmt.Println("User authorized")
		}

		if dbUser == username {
			fmt.Println("Already there in DB", " ", username)
			http.Redirect(w, r, "/users", http.StatusFound)
			return
		}

		_,err= db.Exec("INSERT into users(username, password) VALUES (?,?)", username, password)
		if err== nil {
			fmt.Println("Successfully inserted account details")
		}

		http.Redirect(w, r, "/login", http.StatusFound)
		
	}
}

func users(w http.ResponseWriter, r *http.Request) {
	type User struct {
		DbUser string
		DbPass string
	}

	users:= make([]User, 0)

	rows,err:= db.Query("SELECT username, password FROM users")
	defer rows.Close()
		if err== nil {
			fmt.Println("All Users")

			for rows.Next(){
				u:= User{}
				_= rows.Scan(&u.DbUser, &u.DbPass)
				users= append(users, u)
				fmt.Println(u.DbUser, " ", u.DbPass)
				//fmt.Fprintf(w, "%s\t%s\n", u.DbUser, u.DbPass)
			}
			fmt.Println(users)
			t, _:=template.ParseFiles("static/users.html")
			t.Execute(w, users)

		} else {
			fmt.Println(err.Error())
		}

}

func manageUsers(w http.ResponseWriter, r *http.Request) {
	
	if r.Method== "GET"{

		t, _:=template.ParseFiles("static/manage.html")
		t.Execute(w, nil)

	} else{

		r.ParseForm()
		username:= r.FormValue("username")
		operation:= r.FormValue("operation")
		var dbUser string

		err:= db.QueryRow("SELECT username FROM users WHERE username=? ",username).Scan(&dbUser)
		if err== nil {
			fmt.Println("Look up success")
		}
		if dbUser== username {
			fmt.Println("User found", " ", username)

			switch {
			case operation == "delete" :
				err:= db.QueryRow("DELETE FROM users WHERE username=? ",username)
				if err== nil {
					fmt.Println(operation, "success")
					fmt.Println(err)
				}
			}
		}
		http.Redirect(w, r, "/manage", http.StatusFound)
		
	}
}

func main() {
	db,err= sql.Open("mysql", "root:rupal@/books")
	if err!= nil {
		panic(err.Error())
	}

	defer db.Close()

	err= db.Ping()
	if err!= nil {
		panic(err.Error())
	}

	http.HandleFunc("/", homePage)
	http.HandleFunc("/login", login)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/users", users)
	http.HandleFunc("/manage", manageUsers)
	http.ListenAndServe(":8000", nil)
}
