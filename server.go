package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Customer struct { //Instance of each Customer
	name    string
	age     byte
	id      int
	address string
}

//ParseRecord function uses the Customer data structure and returns a string containing
//the customer data fprmatted into a table component
func (customer *Customer) ParseRecord() string {
	tableData := "<tr><td>" + string(customer.id) + "</td><td>" + customer.name + "</td><td>" + string(customer.age) + "</td><td>" + customer.address + "</td></tr>"
	return tableData
}
func addRecord(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("RENDERING addRecord.html")
}

func renderIndex(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("RENDERING index.html")
	t, _ := template.ParseFiles("index.html")
	t.Execute(writer)
}

func main() {
	http.HandleFunc("/addRecord", addRecord)
	http.HandleFunc("/", renderIndex)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
