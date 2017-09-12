package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"goji.io/pat"

	"goji.io"
	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"
)

const (
	//MongoDBHost ist der Hostname
	MongoDBHost = "127.0.0.1:27017"
	//DBName ist der Name der DB
	DBName = "test"
	//CollectionName ist der Name der Collection
	CollectionName = "people"
)

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Typ", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q", message)

}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(json)
}

// Serie Speichert serien
type Serie struct {
	Name  string
	Cover string
	Folge []Folge
}
type Folge struct {
	FolgenName string
	Nummer     string
	Staffel    string
	Gesehen    bool
}

func main() {
	mongoSession, err := mgo.Dial(MongoDBHost)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}
	mongoSession.SetMode(mgo.Monotonic, true)

	mux := goji.NewMux()
	mux.HandleFunc(pat.Get("/serien"), allEntrys(mongoSession))
	mux.HandleFunc(pat.Post("/serien"), addEntry(mongoSession))
	mux.HandleFunc(pat.Get("/find/:name"), findEntry(mongoSession))
	http.ListenAndServe("localhost:8080", mux)
}

func allEntrys(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DBName).C(CollectionName)

		var serie []Serie
		err := c.Find(bson.M{}).All(&serie)
		if err != nil {
			ErrorWithJSON(w, "Database error",
				http.StatusInternalServerError)
			log.Println("Faild get all person: ", err)
			return
		}

		respBody, err := json.MarshalIndent(serie, "", " ")
		if err != nil {
			log.Fatal(err)
		}
		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func findEntry(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		c := session.DB(DBName).C(CollectionName)

		name := pat.Param(r, "name")

		var serie []Serie

		err := c.Find(bson.M{"name": name}).All(&serie)
		if err != nil {
			ErrorWithJSON(w, "Database error",
				http.StatusInternalServerError)
			log.Println("Faild get all person: ", err)
			return
		}

		respBody, err := json.MarshalIndent(serie, "", " ")
		if err != nil {
			log.Fatal(err)
		}
		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}

func addEntry(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		session := s.Copy()
		defer session.Close()

		var serie Serie
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&serie)
		if err != nil {
			ErrorWithJSON(w, "Incorrect body", http.StatusBadRequest)
			return
		}
		c := session.DB(DBName).C(CollectionName)
		_, err = c.Upsert(bson.M{"name": serie.Name}, serie)
		if err != nil {
			log.Printf("Update : Error : %s\n", err)
		}

		w.Header().Set("Content-Typ", "application/json")
		w.Header().Set("Location", r.URL.Path+"/"+serie.Name)
		w.WriteHeader(http.StatusCreated)
	}
}

//DelEntrys ermögtlich das Löschen von Einträgen
func DelEntrys(mongoSession *mgo.Session) {
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	collection := sessionCopy.DB(DBName).C(CollectionName)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter name: ")
	name, _ := reader.ReadString('\n')

	name = strings.Replace(name, "\n", "", -1)

	err := collection.Remove(bson.M{"name": name})
	if err != nil {
		log.Printf("DelEntrys : ERROR : %s\n", err)
	}

}

// ShowEntrys is a function that is launched as a goroutine to perform
// the MongoDB work.
func ShowEntrys(mongoSession *mgo.Session) {

	// Request a socket connection from the session to process our query.
	// Close the session when the goroutine exits and put the connection back
	// into the pool.
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	// Get a collection to execute the query against
	collection := sessionCopy.DB(DBName).C("people")

	//Ausgabe von Person
	var serie []Serie
	err := collection.Find(nil).All(&serie)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}

	i := 1
	for i < len(serie) {
		fmt.Printf("Name: %s \t Number: %s\n", serie[i].Name, serie[i].Cover)
		i++
	}
}
