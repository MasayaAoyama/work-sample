package main

import (
	"database/sql"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"os"
)

var version = "v1.0"

var (
	dbuser   string
	dbpass   string
	dbhost   string
	dbport   string
	dbname   string
	hostname string
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "/ is requested")
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
	id, name := getUser()
	fmt.Fprintf(w, fmt.Sprintf("[%s] %s: name is [%d %s]\n", version, hostname, id, name))
	fmt.Printf(fmt.Sprintf("[%s] %s: name is [%d %s]\n", version, hostname, id, name))

}

func addUserHandler(w http.ResponseWriter, r *http.Request) {
	name := addUser()
	fmt.Fprintf(w, fmt.Sprintf("[%s] %s: added user [%s]\n", version, hostname, name))
	fmt.Printf(fmt.Sprintf("[%s] %s: added user [%s]\n", version, hostname, name))

}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	check := checkDBConnection()
	if !check {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, fmt.Sprintf("Cannot connect to DB at %s\n", hostname))
		fmt.Printf(fmt.Sprintf("Cannot connect to DB at %s\n", hostname))
	}
	fmt.Printf(fmt.Sprintf("Can connect to DB at %s\n", hostname))
}

func main() {
	getFromEnv()
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/getuser", getUserHandler)
	mux.HandleFunc("/adduser", addUserHandler)
	mux.HandleFunc("/healthz", healthzHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Printf("fin ListenAndServe(): %s", err)
	}

}

func getFromEnv() {
	dbuser = os.Getenv("DBUSER")
	dbpass = os.Getenv("DBPASS")
	dbhost = os.Getenv("DBHOST")
	dbport = os.Getenv("DBPORT")
	dbname = os.Getenv("DBNAME")
	hostname, _ = os.Hostname()
}

func getUser() (int, string) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpass, dbhost, dbport, dbname))
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	stmtOut, err := db.Prepare("SELECT id,name FROM users ORDER BY RAND() LIMIT 1;")
	if err != nil {
		panic(err.Error())
	}
	defer stmtOut.Close()

	var name string
	var id int

	err = stmtOut.QueryRow().Scan(&id, &name)
	if err != nil {
		panic(err.Error())
	}

	return id, name
}

func addUser() string {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpass, dbhost, dbport, dbname))
	if err != nil {
		return err.Error()
	}
	defer db.Close()

	stmtIns, err := db.Prepare("INSERT INTO users(name) VALUES(?)")
	if err != nil {
		return err.Error()
	}
	defer stmtIns.Close()

	name := randomdata.SillyName()
	_, err = stmtIns.Exec(name)
	if err != nil {
		return err.Error()
	}

	return name
}

func checkDBConnection() bool {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbuser, dbpass, dbhost, dbport, dbname))
	if err != nil {
		return false
	}
	defer db.Close()
	return true
}
