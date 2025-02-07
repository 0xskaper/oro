package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type OroAPIResponse struct {
	isSuccess bool         `json:"success"`
	Data      interface{}  `json:"data,omitempty"`
	Error     *OroAPIError `json:"error,omitempty"`
}

type OroAPIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Task struct {
	ID          string         `json:"id" gorm"primaryKey;type:uuid;default:gen_random_uuid()"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	isCompleted bool           `json:"is_completed" gorm:"default:false"`
	isImportant bool           `json:"is_important" gorm:"default:false"`
	DueDate     *time.Time     `json:"due_date"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type TaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	isImportant bool       `json:"is_important"`
	DueDate     *time.Time `json:"due_date"`
}

var oroDB *gorm.DB

func initializeOroDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=oro port=5432 sslmode=disable TimeZone=UTC"
	var err error
	oroDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = oroDB.AutoMigrate(&Task{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}

func main() {
	initializeOroDB()

	r := mux.NewRouter()
	oroApi := r.PathPrefix("/oro/v1").Subrouter()

	oroApi.HandleFunc("/tasks", createTask).Methods("POST")
	oroApi.HandleFunc("/tasks", loadTasks).Methods("GET")
	oroApi.HandleFunc("/tasks/{id}", loadTask).Methods("GET")
	oroApi.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	oroApi.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	oroApi.HandleFunc("/tasks/{id}/complete", toggleTaskCompletion).Methods("POST")
	oroApi.HandleFunc("/tasks/{id}/important", toggleTaskImportance).Methods("POST")

	r.Use(corsMiddleware)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server initialized at post: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
