# Sphere Work Assessment #
## Go MySQL Web Application##
###Submitted by Jeet Thakkar###

 If you would prefer to work on our back-end, please create a Go app that displays a list of values from a MySQL or PostgreSQL database on a webpage and gives the ability to add values to the database with a simple input field and submit button on the same page. To that end, you will have to find out how to connect/save to the database and serve the webpage with Go.

 ##Instructions##
 1. Create the MySQL Database according to the commands listed in mysql-config in this folder. The lines that begin with `$` are entered into the terminal and the commands that begin with `mysql> ` are the ones that are entered into the MySQL application that is executed by running `mysql -u root -p` and entering the password.
 2. Create a new Database and Table using the SQL commands in mysql-config. The databse is named `sphere` and the table is named `CustomerData`. The 4 fields in the table should be named **CustomerId**, **CustomerName**, **CustomerAge** and **CustomerAddress**.
 3. Once, the MySQL database is setup, run `go build server.go && ./server`, and load [localhost:8080/index.html](localhost:8080/index.html) to view the applet.
 4. An empty database looks something like this: 