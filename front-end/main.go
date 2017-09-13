package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Serie struct {
	Name  string
	Cover string
	Folge []Folge
}
type Folge struct {
	FolgenName string
	Nummer     int
	Staffel    int
	Gesehen    bool
}

func main() {

	router := httprouter.New()
	router.GET("/serien", serien)
	router.GET("/find/:name", findEntry)
	router.POST("/suchen", suche)
	router.GET("/add", add)
	log.Fatal(http.ListenAndServe(":4000", router))
}

func suche(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	http.Redirect(w, r, "/find/"+r.FormValue("name"), http.StatusFound)
}

func add(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := template.ParseGlob("template/*.html")
	if err != nil {
		fmt.Print(err)
	}
	// fürt das Template von bereich Contetn aus
	err = t.ExecuteTemplate(w, "addentry", "test")
	if err != nil {
		panic(err)
	}
}

func findEntry(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	url := "http://localhost:8080/find/" + ps.ByName("name")
	//Holt Contetn von einer Seite
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	//Liest den Body einer seite ein
	body, err := ioutil.ReadAll(res.Body)
	// Packt Json in eine Struck
	var serie []Serie
	json.Unmarshal(body, &serie)
	//HTML
	t, err := template.ParseGlob("template/*.html")
	if err != nil {
		fmt.Print(err)
	}
	// fürt das Template von bereich Contetn aus
	err = t.ExecuteTemplate(w, "content", serie)
	if err != nil {
		panic(err)
	}
}

func serien(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	url := "http://localhost:8080/serien"

	//Holt Contetn von einer Seite
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	//Liest den Body einer seite ein
	body, err := ioutil.ReadAll(res.Body)
	// Packt Json in eine Struck
	var serie []Serie
	json.Unmarshal(body, &serie)
	//HTML
	t, err := template.ParseGlob("template/*.html")
	if err != nil {
		fmt.Print(err)
	}
	// fürt das Template von bereich Contetn aus
	err = t.ExecuteTemplate(w, "content", serie)
	if err != nil {
		panic(err)
	}
}
