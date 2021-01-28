package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type ToDo struct {
	gorm.Model
	Task string
	Due  time.Time
}

type todosResource struct {
	db *gorm.DB
}

// Routes creates a REST router for the todos resource
func (rs todosResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", rs.List)    // GET /todos - read a list of todos
	r.Post("/", rs.Create) // POST /todos - create a new todo and persist it
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a todos map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", rs.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func (rs *todosResource) List(w http.ResponseWriter, r *http.Request) {
	var t []*ToDo
	result := rs.db.Find(&t)
	if result.Error != nil {
		JSONError(w, result.Error)
		return
	}

	JSON(w, t)
}

func (rs *todosResource) Create(w http.ResponseWriter, r *http.Request) {
	t := &ToDo{}
	if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		log.Printf("could not decode: %+v", err)
		JSONError(w, fmt.Errorf("couldn't decode json"))
		return
	}

	result := rs.db.Create(t)
	if result.Error != nil {
		JSONError(w, result.Error)
		return
	}

	JSON(w, t)
}

func (rs *todosResource) Get(w http.ResponseWriter, r *http.Request) {
	u, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		log.Printf("could not parse id: %v", err)
		JSONError(w, err)
		return
	}
	t := &ToDo{}
	result := rs.db.First(t, uint(u))
	if result.Error != nil {
		JSONError(w, result.Error)
		return
	}

	JSON(w, t)
}

func (rs *todosResource) Update(w http.ResponseWriter, r *http.Request) {
	u, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		log.Printf("could not parse id: %v", err)
		JSONError(w, err)
		return
	}

	t := &ToDo{}
	if err := json.NewDecoder(r.Body).Decode(t); err != nil {
		log.Printf("could not decode: %+v", err)
		JSONError(w, fmt.Errorf("couldn't decode json"))
		return
	}

	t.ID = uint(u)
	result := rs.db.Save(t)
	if result.Error != nil {
		JSONError(w, result.Error)
		return
	}

	JSON(w, t)
}

func (rs *todosResource) Delete(w http.ResponseWriter, r *http.Request) {
	u, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		log.Printf("could not parse id: %v", err)
		JSONError(w, err)
		return
	}

	result := rs.db.Delete(&ToDo{}, uint(u))
	if result.Error != nil {
		JSONError(w, result.Error)
		return
	}

	JSON(w, map[string]string{"status": "success"})
}
