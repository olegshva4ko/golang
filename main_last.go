package main
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	_"github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	)
const(CONN_HOST = "localhost"
	CONN_PORT = "8080"
	DRIVER_NAME = "mysql"
	DATA_SOURCE_NAME = "root:pswd@/book_shelve")

var db *sql.DB
var connectionError error

func init() {
	db, connectionError = sql.Open(DRIVER_NAME, DATA_SOURCE_NAME)
	if connectionError != nil {
		log.Fatal("error connecting to database :: ", connectionError)
	}
}

func createRecord(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("going to insert record in database for name : ", name[0])
		stmt, err := db.Prepare("INSERT books SET name = ?")

		if err != nil {
			log.Print("error preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0])
		if err != nil {
			log.Print("error executing query :: ", err)
			return
		}

		id, err := result.LastInsertId()
		_, _ = fmt.Fprintf(w, "Last Inserted Record Id is :: %s", strconv.FormatInt(id, 10))
	} else {
		_, _ = fmt.Fprintf(w, "Error occurred while creating record in database for name :: %s", name[0])
	}
}

func getCurrentDb(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT DATABASE() as db")
	if err != nil {
		log.Print("error executing query :: ", err)
		return
	}

	var db string
	for rows.Next() {
		_ = rows.Scan(&db)
	}

	_, _ = fmt.Fprintf(w, "Current Database is :: %s", db)
}

type Books struct {
	Id int `json:"uid"`
	Name string `json:"name"`
}

func readRecords(w http.ResponseWriter, r *http.Request) {
	log.Print("reading records from database")
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		log.Print("error occurred while executing select query :: ",err)
		return
	}
	bookArray := []Books{}
	for rows.Next() {

		var uid int
		var name string
		err = rows.Scan(&uid, &name)
		book := Books{Id: uid, Name: name}
		bookArray = append(bookArray, book)
	}
	json.NewEncoder(w).Encode(bookArray)
}

func updateRecord(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("going to update record in database for id :: ", id)
		stmt, err := db.Prepare("UPDATE books SET name=? where uid = ?")
		if err != nil {
			log.Print("error occurred while preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0], id)
		if err != nil {
			log.Print("error occurred while executing query :: ", err)
			return
		}
		rowsAffected, err := result.RowsAffected()
		fmt.Fprintf(w, "Number of rows updated in database are :: %d",rowsAffected)
	} else {
		fmt.Fprintf(w, "Error occurred while updating record in database for id :: %s", id)
		}
}

func deleteRecord(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	name, ok := vals["name"]
	if ok {
		log.Print("going to delete record in database for name :: ", name[0])
		stmt, err := db.Prepare("DELETE from books where name = ?")
		if err != nil {
			log.Print("error occurred while preparing query :: ", err)
			return
		}
		result, err := stmt.Exec(name[0])
		if err != nil {
			log.Print("error occurred while executing query :: ", err)
			return
		}
		rowsAffected, err := result.RowsAffected()
		fmt.Fprintf(w, "Number of rows deleted in database are :: %d", rowsAffected)
	} else{
		fmt.Fprintf(w, "Error occurred while deleting record in database for name %s", name[0])
	}
}

func main() {
//	http.HandleFunc("/", getCurrentDb)
	router := mux.NewRouter()
	router.HandleFunc("/", getCurrentDb)
	router.HandleFunc("/books/create", createRecord).Methods("POST")
	router.HandleFunc("/books", readRecords).Methods("GET")
	router.HandleFunc("/books/update/{id}",
		updateRecord).Methods("PUT")
	router.HandleFunc("/books/delete",
		deleteRecord).Methods("DELETE")
	http.Handle("/", router)
	defer db.Close()
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)

	if err != nil {
		log.Fatal("error starting http server :: ", err)
		return
	}
}