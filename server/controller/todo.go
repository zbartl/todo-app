package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gobridge-kr/todo-app/server/database"
	"net/http"
	"strings"
)

// TodoController is a MVC controller to handle todo requests
type TodoController struct {
	database *database.Database
}

func (t *TodoController) ServeHTTP(
	w http.ResponseWriter,
	r *http.Request,
) {
	parts := strings.Split(r.URL.Path, "/")
	route := parts[1]
	if route != "todo" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	
	hasID := len(parts) > 2
	var id string
	if hasID {
		id = parts[2]
	}

	switch r.Method {
	case "GET":
		if hasID {
			t.GetOne(w, r, id)
		} else {
			t.GetAll(w, r)
		}
	case "POST":
		if hasID {
			t.PostOne(w, r, id)
		} else {
			t.PostAll(w, r)
		}
	case "PATCH":
		if hasID {
			t.PatchOne(w, r, id)
		} else {
			t.PatchAll(w, r)
		}
	case "DELETE":
		if hasID {
			t.DeleteOne(w, r, id)
		} else {
			t.DeleteAll(w, r)
		}
	case "OPTIONS":
		t.Options(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
} 

// GetOne retrieves and shows a todo by given ID
func (t *TodoController) GetOne(w http.ResponseWriter, r *http.Request, id string) {
	todo, err := t.database.GetTodo(id)
	if err != nil {
		if err == database.ErrItemNotFound {
			http.Error(w, "Item Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

// GetAll lists all todos
func (t *TodoController) GetAll(w http.ResponseWriter, r *http.Request) {
	todos := t.database.GetTodos()
	json.NewEncoder(w).Encode(todos)
}

// PostOne requests(POSTs with an ID) is not allowed
func (t *TodoController) PostOne(w http.ResponseWriter, r *http.Request, id string) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

// PostAll creates a todo
func (t *TodoController) PostAll(w http.ResponseWriter, r *http.Request) {
	var params map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	todo, err := t.database.AddTodo(params)
	if err != nil {
		if err == database.ErrBadRequest {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

// PatchOne updates a todo with given ID
func (t *TodoController) PatchOne(w http.ResponseWriter, r *http.Request, id string) {
	var params map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	todo, err := t.database.UpdateTodo(id, params)
	if err != nil {
		if err == database.ErrItemNotFound {
			http.Error(w, "Item Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

// PatchAll requests(PATCHs without an ID) is not allowed
func (t *TodoController) PatchAll(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Allowed", http.StatusMethodNotAllowed)
}

// DeleteOne deletes a todo by given ID
func (t *TodoController) DeleteOne(w http.ResponseWriter, r *http.Request, id string) {
	todo, err := t.database.DeleteTodo(id)
	if err != nil {
		if err == database.ErrItemNotFound {
			http.Error(w, "Item Not Found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(todo)
}

// DeleteAll deletes all todos
func (t *TodoController) DeleteAll(w http.ResponseWriter, r *http.Request) {
	t.database.DeleteTodos()
}

// Options returns server options data
func (t *TodoController) Options(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "")
}

// Todo is a constructor for TodoControllers
func Todo(database *database.Database) *TodoController {
	return &TodoController{
		database: database,
	}
}
