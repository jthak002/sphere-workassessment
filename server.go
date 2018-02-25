package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
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

//PageBody stres the entirety of the page in a struct. The fixed parts are
//indexPreTable and indexPostTable. tableFields stores the HTML formatted table content
type PageBody struct {
	indexPreTable  []byte
	indexPostTable []byte
}

var pageBody PageBody

//unexported global db variable to store the database open variable
var db *sql.DB

//pageRenderInit reads the pre - and post-table parts of the HTML document
//since this doesn't change, the operation can be computed only once.
func pageRenderInit() error {
	indexPreTableLocal, err := ioutil.ReadFile("index-pre")
	if err != nil {
		log.Fatal("index.pre render failed")
	}
	indexPostTableLocal, err := ioutil.ReadFile("index-post")
	if err != nil {
		log.Fatal("index.post render failed")
	}
	pageBody = PageBody{indexPreTable: indexPreTableLocal,
		indexPostTable: indexPostTableLocal}
	return err
}

//createHTMLTable converts a given Customer slice into a byte stream
//od HTML formatted Table Data
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

//ParseRecord function uses the Customer data structure and returns a string containing
//the customer data formatted into a table component
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

//renderIndex function renders the index.html page on request from the clients
func renderIndex(writer http.ResponseWriter, request *http.Request) {
	log.Println("Server: render index.html")
	customers := dbSpan()
	tableFields := createHTMLTable(&customers)
	preAndTable := append(pageBody.indexPreTable, tableFields...)
	finalObject := append(preAndTable, pageBody.indexPostTable...)
	fmt.Fprintf(writer, string(finalObject))

}

//dbAdd function receives an object of type customer and adds it to the
//MySQL Database
func dbAdd(customer Customer) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	var buffer bytes.Buffer
	//CONSOLIDATING THE ENTIRE SQL COMMAND INTO A SINGLE STATEMENT USING A bytes.Buffer
	buffer.WriteString("INSERT INTO `CustomerData` (CustomerId, CustomerName, CustomerAge, CustomerAddress) VALUES (")
	//POPULATING THE ACTUAL VALUES FROM THE SQL DATABASE INTO THE buffer
	customerStringSlice := []string{customer.id, customer.name, customer.age, customer.address}
	buffer.WriteString("'" + customerStringSlice[0] + "', '")
	buffer.WriteString(customerStringSlice[1] + "', '")
	buffer.WriteString(customerStringSlice[2] + "', '")
	buffer.WriteString(customerStringSlice[3] + "');")
	//USING Prepare() TO CREATE THE SQL COMMAND INTO a Statement OBJECT
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

//addRecord function accepts input from the form and then adds the input to
//the MySQL Database. The page refreshes the page to display the new values of the table.
func addRecord(writer http.ResponseWriter, request *http.Request) {
	log.Println("Server: render /addRecord")
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
	log.Println("MySQL.sphere:Record Created Successfully!")
	dbAdd(customer)
	log.Println("MySQL.sphereCustomer Added Successfully!")
	customers := dbSpan()
	tableFields := createHTMLTable(&customers)
	html1 := append(pageBody.indexPreTable, tableFields...)
	finalObject := append(html1, pageBody.indexPostTable...)
	fmt.Fprintf(writer, string(finalObject))
}

func main() {
	var err error          //func main shared error variable
	err = pageRenderInit() //Read the Pre-Table and Post-Table part of the HTML Page
	if err != nil {
		log.Println("Page Render Failed!")
	} else {
		log.Println("Page Render Successful!")
	}
	//Initialize the connection to the SQL Database
	db, err = sql.Open("mysql", "sphere:SphereWorkAssessment234()@tcp(127.0.0.1:3306)/sphere")
	//SQL Error Connection Troubleshooting
	if err != nil {
		log.Printf("MySQL.Sphere not responding%e\n", err)
	}
	//SQL Connection Life Check - Performed before each SQL Transaction
	err = db.Ping()
	if err != nil {
		log.Println("MySQL.Sphere Offline")
	}
	//If SQL Connection established and Live display SUCCESS
	if err != nil {
		log.Println("MySQL.Sphere DB Init Failed")
	} else {
		log.Println("MySQL.Sphere DB Init Success")
	}
	defer db.Close()
	http.HandleFunc("/addRecord", addRecord)
	http.HandleFunc("/index.html", renderIndex)
	http.HandleFunc("/", renderIndex)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
