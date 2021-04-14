package app

import (
	c "./controllers"
	m "./models"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	m.Connect()
}

func App() {
	http.HandleFunc("/", c.Home)
	http.HandleFunc("/api/auth/login/", c.Auth)
	http.HandleFunc("/api/auth/refresh/", c.AuthRefreshToken)
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_SERVER"), nil))
}

func CreateTestUser() {
	password := []byte("admin")
	hash, _ := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	newUser := m.User{ // example user
		ID:        "48093dc5-3d2d-4a36-83a8-23f6f87177f8",
		Name:      "Admin",
		Admin:     true,
		Password:  string(hash),
		CreatedAt: time.Now(),
	}

	err := m.CreateUser(newUser)
	if err != nil {
		log.Fatal("Couldn't create a new user")
	} else {
		log.Println("Success, a new user was created")
	}
}
