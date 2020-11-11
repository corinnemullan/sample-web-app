// Corinne's first Go web app
//
// The homepage displays a random quote from the database each time it is loaded.
// Quotes can be added to and deleted from the database.

package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var db *gorm.DB

type Quote struct {
	gorm.Model
	Saying string
	Author string
}

type PageVariables struct {
	Message string
	Person  string
	Id uint
}

func homePage(w http.ResponseWriter, r *http.Request) {
	var quote Quote
	var count int64
	var randomOffset int
	var pageVariables PageVariables

	rand.Seed(time.Now().UnixNano())
	db.Model(&Quote{}).Count(&count)

	if count > 0 {
		randomOffset = rand.Intn(int(count))
		db.Offset(randomOffset).First(&quote)

		pageVariables = PageVariables{
			Message: quote.Saying,
			Person: quote.Author,
			Id: quote.ID,
		}
	}

	t, err := template.ParseFiles("homepage.html")
	if err != nil {
		log.Print("Template parsing error: ", err)
	}
	err = t.Execute(w, pageVariables)
	if err != nil {
		log.Print("Template execute error: ", err)
	}
}

func addQuote(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
		case "GET":
			http.ServeFile(w, r, "add_quote.html")

		case "POST":
			err := r.ParseForm()
			if err != nil {
				fmt.Fprintf(w, "Error parsing form: %v", err)
				return
			}

			saying := r.FormValue("saying")
			author := r.FormValue("author")

			quote := Quote{Saying: saying, Author: author}
			result := db.Create(&quote)

			if result.Error != nil {
				fmt.Fprintf(w, "Error adding quote to database: %v", result.Error)
				return
			}
			http.Redirect(w, r, "/", http.StatusFound)
	}
}

func deleteQuote(w http.ResponseWriter, r *http.Request) {
	var quote Quote
	var pageVariables PageVariables
	var id uint
	var id64 uint64

	id64, _ = strconv.ParseUint(r.URL.Path[len("/delete/"):], 10, 64)
	id = uint(id64)

	switch r.Method {
		case "GET":
			db.First(&quote, id)

			pageVariables = PageVariables{
				Message: quote.Saying,
				Person: quote.Author,
				Id: id,
			}

			t, err := template.ParseFiles("delete_quote.html")
			if err != nil {
				log.Print("Template parsing error: ", err)
			}
			err = t.Execute(w, pageVariables)
			if err != nil {
				log.Print("Template execute error: ", err)
			}
		case "POST":
			db.Delete(&Quote{}, id)
			http.Redirect(w, r, "/", http.StatusFound)
		}
}

func main() {
	var err error 
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Error connecting to database")
	}

	// Migrate the schema
	db.AutoMigrate(&Quote{})

	http.HandleFunc("/", homePage)
	http.HandleFunc("/add/", addQuote)
	http.HandleFunc("/delete/", deleteQuote)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
