// Corinne's first Go web app
// Just to try stuff out
package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"log"
	"math/rand"
	"net/http"
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
}

func homePage(w http.ResponseWriter, r *http.Request) {
	var quote Quote
	var count int64
	var randomId int
	var homePageVariables PageVariables

	rand.Seed(time.Now().UnixNano())
	db.Model(&Quote{}).Count(&count)

	// Select a random entry from the Quote table to display
	if count == 0 {
		homePageVariables = PageVariables{
			Message: "Please add some quotes to your database!",
			Person: "Snarky Developer",
		}
	} else {
		if count == 1 {
			randomId = 1
		} else {
			randomId = rand.Intn(int(count)) + 1
		}

		db.First(&quote, randomId)

		homePageVariables = PageVariables{
			Message: quote.Saying,
			Person: quote.Author,
		}
	}

	t, err := template.ParseFiles("homepage.html")
	if err != nil {
		log.Print("Template parsing error: ", err)
	}
	err = t.Execute(w, homePageVariables)
	if err != nil {
		log.Print("Template execute error: ", err)
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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
