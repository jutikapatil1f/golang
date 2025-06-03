package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	ID       int    `json:"id"`
	Task     string `json:"task"`
	Complete bool   `json:"complete"`
}

var todos = []Todo{
	{ID: 1, Task: "Write your first Task!! Edit and Delete as you wish ", Complete: false},
}

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next(w, r)
	}
}

// func enableCors(w http.ResponseWriter) {
// 	w.Header().Set("Access-Control-Allow-Origin", "*")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// }

func handleTodos(w http.ResponseWriter, r *http.Request) {
	// enableCors(w)
	//set content type
	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(todos) // encode all todos as json (GET)
		return
	}

	if r.Method == http.MethodPost {
		var newTodo Todo
		err := json.NewDecoder(r.Body).Decode(&newTodo) // decode newTodo from json Body
		if err != nil {
			fmt.Print("error", err)
			return
		}

		newTodo.ID = len(todos) + 1
		todos = append(todos, newTodo)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTodo) // encode that new todo (POST)

		return

	}

	if r.Method == http.MethodPut {
		var updatedTodo Todo
		err := json.NewDecoder(r.Body).Decode(&updatedTodo)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		for i, todo := range todos {
			if todo.ID == updatedTodo.ID {
				todos[i].Task = updatedTodo.Task
				todos[i].Complete = updatedTodo.Complete

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(todos[i])
				return
			}
		}

		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	if r.Method == http.MethodDelete {
		var deleteTodo Todo
		err := json.NewDecoder(r.Body).Decode(&deleteTodo)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		for i, todo := range todos {
			if todo.ID == deleteTodo.ID {
				todos = append(todos[:i], todos[i+1:]...)

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(todo)
				return
			}
		}

		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World! Server is running on port 8080.")

}

func main() {
	http.HandleFunc("/", withCORS(homeHandler)) //define route
	http.HandleFunc("/todos", withCORS(handleTodos))
	fmt.Print("Server is starting on port 8080.")
	err := http.ListenAndServe(":8080", nil) // to start server

	if err != nil {
		fmt.Print("error", err)
	}

}
