// Corinne's first Go web app
// Just to try stuff out
package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"html/template"
	"log"
	"net/http"
)

var db *gorm.DB

type Quote struct {
	gorm.Model
	Saying string
	Author string
}

type PageVariables struct {
	Message string
}

func homePage(w http.ResponseWriter, r *http.Request) {
	var quote Quote
	db.First(&quote)
	homePageVariables := PageVariables{
		Message: quote.Saying,		
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

	// Populate some data
	db.Create(&Quote{Saying: "I like cats CatS CAts CATS", Author: "Ria"})
	db.Create(&Quote{Saying: "What can I eat?", Author: "Logan"})

	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
