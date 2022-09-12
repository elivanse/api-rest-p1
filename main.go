package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"cratedat"`
}

// defino un map noteStore de tipo Note
var noteStore = make(map[string]Note)

var id int

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {

	//creo un array temporal tipo Note para recorrer
	var notes []Note

	//recorro el noteStore, cada Note mapeado lo apendo al array temporal
	for _, v := range noteStore {
		notes = append(notes, v)
	}

	// defino la cabecera de la respuesta en formato jason
	w.Header().Set("Content-Type", "application/json")

	// marshaleo jasonita el array
	// si hay error que cunda el panico
	j, err := json.Marshal(notes)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func PostNoteHandler(w http.ResponseWriter, r *http.Request) {
	//note de tipo Note
	var note Note
	// del requerimiento request r
	// hago el Decode a ver si el formato jason
	// es correcto
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	// asigno la fecha actual al campo Createdat del note
	note.CreatedAt = time.Now()
	// muevo el id despues de usado
	id++
	// asigno el id en formato string a k para
	// ubicarlo en el noteStore recordando map[string]
	k := strconv.Itoa(id)
	noteStore[k] = note

	// defino cabecera de jason y lo armo
	w.Header().Set("Content-type", "application/json")
	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}

	// autputeo la jasonita armada
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	if _, ok := noteStore[k]; ok {
		delete(noteStore, k)
	} else {
		log.Printf("No enontramos el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)

}

func PutNoteHandler(w http.ResponseWriter, r *http.Request) {

	// paso a vars lo que resuelve gorilla
	// crea un slice con id en strings
	vars := mux.Vars(r)

	// paso a k temp el id
	k := vars["id"]

	// defino noteUpdate tipo Note y la formateo jasonita
	var noteUpdate Note
	err := json.NewDecoder(r.Body).Decode(&noteUpdate)
	if err != nil {
		panic(err)
	}

	// la asignacion de un map nos devuelve
	// adicionalmente un booleano que nos dice
	// si efectivamente hay algo
	// por eso si noteupdate y si ok sigue
	if note, ok := noteStore[k]; ok {
		noteUpdate.CreatedAt = note.CreatedAt
		delete(noteStore, k)
		noteStore[k] = noteUpdate
	} else {
		log.Printf("No se encontro ningun elemento con el id %s", k)
	}
	w.WriteHeader(http.StatusNoContent)

}

func main() {

	r := mux.NewRouter().StrictSlash(false)

	r.HandleFunc("/api/notes", GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes", PostNoteHandler).Methods("POST")
	r.HandleFunc("/api/notes{id}", PutNoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes{id}", DeleteNoteHandler).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Listening at 8080 ...")
	server.ListenAndServe()
}
