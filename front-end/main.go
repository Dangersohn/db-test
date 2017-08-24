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
	Nummer     string
	Staffel    string
}

func main() {

	router := httprouter.New()
	router.GET("/index", index)
	log.Fatal(http.ListenAndServe(":4000", router))
}

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	fmt.Print(serie[0].Folge[1])
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
