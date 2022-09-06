package mainb

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

var noteStore = make(map[string]Note)

var id int

func GetNoteHandler(w http.ResponseWriter, r *http.Request) {

	var notes []Note

	for _, v := range noteStore {
		notes = append(notes, v)
	}

	w.Header().Set("Content-Type", "application/json")

	j, err := json.Marshal(notes)

	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(j)

}

func PostNoteHandler(w http.ResponseWriter, r *http.Request) {

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	note.CreatedAt = time.Now()
	id++
	k := strconv.Itoa(id)
	noteStore[k] = note

	w.Header().Set("Content-type", "application/json")
	j, err := json.Marshal(note)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

func PutNoteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]
	var noteUpdate Note

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
		MaxHEaderBytes: 1 << 20,
	}
	log.Println("Listening at 8080 ...")
	server.ListenAndServe()
}
