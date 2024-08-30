package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if os.Getenv("MODE") != "PROD" {
		godotenv.Load()
	}
}

func main() {
	r := mux.NewRouter()
	/*
		Handle users
	*/

	/*
		Handle image
	*/

	/*
		Handle Appointement
	*/
	r.HandleFunc("/appointment", handleCreateAppointment)
	r.HandleFunc("/appointment/{id}", handleDeleteAppointment)
	r.HandleFunc("/appointments", handleWatchAppointment)
	r.HandleFunc("/appointment/{id}/edit", handleUpdateAppointment)
	connectTodataBase()

	port := getListeningPort()
	log.Printf("Server is Listening on port %s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))

}

func getListeningPort() string {
	const defaultPort = "7000"

	port := os.Getenv("PORT")

	if port == "" {
		return defaultPort
	}

	return port
}
