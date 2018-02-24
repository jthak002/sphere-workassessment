package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//Customer struct to store the details of the Customers.
type Customer struct {
	name    string
	age     string
	id      string
	address string
}

//unexported global db variable to store the database open variable
var db *sql.DB

//ParseRecord function uses the Customer data structure and returns a string containing
//the customer data fprmatted into a table component
func (customer *Customer) ParseRecord() string {
	tableData := "<tr><td>" + string(customer.id) + "</td><td>" + customer.name + "</td><td>" + string(customer.age) + "</td><td>" + customer.address + "</td></tr>"
	return tableData
}

//CreateRecord function makes an instance of the customer struct for SQL storage purposes
func CreateRecord(name string, age string, id string, address string) *Customer {
	customer := Customer{
		name:    name,
		age:     age,
		id:      id,
		address: address,
	}
	return &customer
}

//addRecord function accepts input from the form and thenadds the input to
//the MySQL Database. The page refreshes the page to display the new values of the table.
func addRecord(writer http.ResponseWriter, request *http.Request) {
	log.Println("RENDERING addRecord.html")
	request.ParseForm()
	formMap := request.Form
	var cname, caddress, cid, cage string
	for key, value := range formMap {
		if key == "frm_name" {
			cname = strings.Join(value, "")
		} else if key == "frm_age" {
			cage = strings.Join(value, "")
		} else if key == "frm_id" {
			cid = strings.Join(value, "")
		} else if key == "frm_address" {
			caddress = strings.Join(value, "")
		}
	}
	customer := CreateRecord(cname, cage, cid, caddress)
	log.Println(*customer)

}

//renderIndex function renders the index.html page on request from the clients
func renderIndex(writer http.ResponseWriter, request *http.Request) {
	log.Println("RENDERING index.html")
	t, _ := template.ParseFiles("index.html")
	lol := "Hello"
	t.Execute(writer, lol)
}

//dbinit function intiilizes the MySQL Database to accept inputs from the application form
//If the SQL database exists, and is ready to accept input the process will continue normally
func dbInit() error {
	db, err := sql.Open("mysql", "sphere:SphereWorkAssessment234()@tcp(127.0.0.1:3306)/sphere")
	if err != nil {
		log.Printf("MySQL.Sphere not responding%e\n", err)
	}
	err = db.Ping()
	if err != nil {
		log.Println("MySQL.Sphere Offline")
	}
	return err
}
func main() {
	err := dbInit() //Initialize the MySQL Database - Store in db global variable
	if err != nil {
		log.Println("MySQL.Sphere DB Init Failed")
	} else {
		log.Println("MySQL.Sphere DB Init Success")
	}
	defer db.Close()
	http.HandleFunc("/addRecord", addRecord)
	http.HandleFunc("/", renderIndex)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
