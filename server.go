package main

import (
	"bytes"
	"database/sql"
	"fmt"
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

//TableBody stores the HTML formatted table rows with the corresposding data
//for each customer record from the SQL database
type TableBody struct {
	TableFields []byte
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
func createRecord(name string, age string, id string, address string) Customer {
	customer := Customer{
		name:    name,
		age:     age,
		id:      id,
		address: address,
	}
	return customer
}

//createHTMLTable converts a given Customer slice into a byte stream
//TRIED CONVERTING IT INTO string but the HTMLEscaper functions that are a part of
//convert all the tags into escape sequences. Byte Streams get past those filters
func createHTMLTable(customers *[]Customer) []byte {
	var buffer bytes.Buffer
	for _, value := range *customers {
		buffer.WriteString("<tr> ")
		buffer.WriteString("<td>" + value.id + "</td> ")
		buffer.WriteString("<td>" + value.name + "</td> ")
		buffer.WriteString("<td>" + value.age + "</td> ")
		buffer.WriteString("<td>" + value.address + "</td> ")
		buffer.WriteString("</tr>\n")
	}
	return buffer.Bytes()
}

//addRecord function accepts input from the form and then adds the input to
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
	customer := createRecord(cname, cage, cid, caddress)
	log.Println("Record Created Successfully!")
	log.Println(customer)
	dbAdd(customer)
	log.Println("Customer Added Successfully!")
	customers := dbSpan()
	log.Println(customers)
	byteObject := TableBody{TableFields: createHTMLTable(&customers)}
	t, _ := template.ParseFiles("index")
	t.Execute(writer, byteObject)

}

//renderIndex function renders the index.html page on request from the clients
func renderIndex(writer http.ResponseWriter, request *http.Request) {
	log.Println("RENDERING index.html")
	customers := dbSpan()
	log.Println(customers)
	byteObject := TableBody{TableFields: createHTMLTable(&customers)}
	t, _ := template.ParseFiles("index")
	t.Execute(writer, byteObject)
}

//dbAdd function receives an object of type customer and adds it to the
//MySQL Database
func dbAdd(customer Customer) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	var buffer bytes.Buffer
	buffer.WriteString("INSERT INTO `CustomerData` (CustomerId, CustomerName, CustomerAge, CustomerAddress) VALUES (")
	customerStringSlice := []string{customer.id, customer.name, customer.age, customer.address}
	buffer.WriteString("'" + customerStringSlice[0] + "', '")
	buffer.WriteString(customerStringSlice[1] + "', '")
	buffer.WriteString(customerStringSlice[2] + "', '")
	buffer.WriteString(customerStringSlice[3] + "');")
	log.Println(buffer.String())
	stmt, err := db.Prepare(buffer.String())
	if err != nil {
		log.Println(err)
	}
	res, err := stmt.Exec()
	if err != nil {
		log.Println(err)
		log.Println(res)
	}
}

//dbSpan function sends a query to MySQL Database and then collects the response
//which is used to display the results onto HTML page
func dbSpan() []Customer {
	err := db.Ping()
	if err != nil {
		log.Println("MySQL.Sphere Live check failed")
	}
	rows, err := db.Query("SELECT * FROM CustomerData ")
	if err != nil {
		log.Println("MySQL.Sphere Query Failed")
	} else {
		if rows == nil {
			log.Println("MySQL.Sphere No Records Exist")
		}
	}
	var name, age, id, address string
	var customers []Customer
	for rows.Next() {
		err := rows.Scan(&id, &name, &age, &address)
		if err != nil {
			log.Println("MySQL.sphere Record parsing failed!")
			log.Fatal(err)
		}
		customer := createRecord(name, age, id, address)
		customers = append(customers, customer)
	}
	err = rows.Err()
	if err != nil {
		//log.Fatal(err)
		fmt.Println("err")
	}
	return customers
}

func main() {
	var err error
	db, err = sql.Open("mysql", "sphere:SphereWorkAssessment234()@tcp(127.0.0.1:3306)/sphere")
	if err != nil {
		log.Printf("MySQL.Sphere not responding%e\n", err)
	}
	err = db.Ping()
	if err != nil {
		log.Println("MySQL.Sphere Offline")
	}
	if err != nil {
		log.Println("MySQL.Sphere DB Init Failed")
	} else {
		log.Println("MySQL.Sphere DB Init Success")
	}
	defer db.Close()
	http.HandleFunc("/addRecord", addRecord)
	http.HandleFunc("/index.html", renderIndex)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
