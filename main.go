package main

import (
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"

	"html/template"
	"log"
	"net/http"
	"time"
)

type Blog struct {
	Id      string
	Author  string
	Title   string
	Content string
	Date    string
}

var db *sql.DB

func Checkerr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:abeokuta101@tcp(127.0.0.1:3306)/deji")
	Checkerr(err)
	defer db.Close()
	fmt.Println("Successfully connected to mysql Database")

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Get("/", Home)
	router.Post("/post", PostBlog)
	router.Get("/edit/{Id}", EditPage)
	router.Get("/delete/{Id}", DeletePage)
	router.Post("/update/{Id}", Update)

	fmt.Println("listening")
	log.Fatal(http.ListenAndServe(":2022", router))
}

func Home(w http.ResponseWriter, req *http.Request) {
	rows, err := db.Query("SELECT * FROM deji.Blogposts")
	Checkerr(err)
	defer rows.Close()

	var data Blog
	var Blogposts []Blog
	for rows.Next() {
		err = rows.Scan(&data.Id, &data.Title, &data.Author, &data.Content, &data.Date)
		Checkerr(err)
		Blogposts = append(Blogposts, data)
	}

	temp := template.Must(template.ParseFiles("Home.html"))
	err = temp.Execute(w, Blogposts)
	Checkerr(err)

}

func PostBlog(w http.ResponseWriter, req *http.Request) {
	InputAuthor := req.FormValue("author")
	InputTitle := req.FormValue("title")
	InputContent := req.FormValue("content")

	date := time.Now().Format("Mon Jan 02 15:04:05")

	ins, err := db.Prepare("INSERT INTO `deji`.`Blogposts` (`Id`,`Author`,`Title`,`Content`,`Date`) VALUES (?,?,?,?,?);")
	Checkerr(err)
	defer ins.Close()

	res, err := ins.Exec(uuid.NewString(), InputAuthor, InputTitle, InputContent, date)
	rowsAffec, _ := res.RowsAffected()
	if err != nil || rowsAffec != 1 {
		fmt.Println("error inserting row: ", err)
	}

	http.Redirect(w, req, "/", 301)
}

func EditPage(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "Id")
	row, err := db.Query("SELECT * FROM deji.Blogposts WHERE Id = ? ;", id)
	Checkerr(err)

	data := Blog{}
	for row.Next() {
		var Id, Author, Title, Content, Date string
		err := row.Scan(&Id, &Title, &Author, &Content, &Date)
		Checkerr(err)
		data.Id = Id
		data.Author = Author
		data.Title = Title
		data.Content = Content
		data.Date = Date

	}
	temp := template.Must(template.ParseFiles("edit.html"))
	err = temp.Execute(w, data)

}

func DeletePage(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "Id")
	del, err := db.Prepare("DELETE FROM `deji`.`Blogposts` WHERE (`Id`= ?);")
	Checkerr(err)
	defer del.Close()
	var res sql.Result
	res, err = del.Exec(id)
	rowsAff, _ := res.RowsAffected()
	fmt.Println("rowsAff:", rowsAff)
	Checkerr(err)

	http.Redirect(w, req, "/", 302)

}

func Update(w http.ResponseWriter, req *http.Request) {

	id := chi.URLParam(req, "Id")
	//InputAuthor := req.FormValue("author")
	InputTitle := req.FormValue("title")
	InputContent := req.FormValue("content")

	//now := time.Now()
	//date := now.Format("Mon Jan 02 15:04:05")
	//upStmt := "UPDATE `testdb`.`products` SET `name` = ?, `price` = ?, `description` = ? WHERE (`idproducts` = ?);"
	statement, err := db.Prepare("UPDATE Blogposts SET Title = ?, Content= ? WHERE (Id= ?);")

	Checkerr(err)
	defer statement.Close()
	var res sql.Result
	res, err = statement.Exec(InputTitle, InputContent, id)
	rowsAff, _ := res.RowsAffected()
	if err != nil || rowsAff != 1 {
		fmt.Println(err)
	}

	Checkerr(err)
	http.Redirect(w, req, "/", 301)

}
