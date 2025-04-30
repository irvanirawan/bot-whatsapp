package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

//go:embed static/*
var embeddedFiles embed.FS

func main() {
	// Load environment variables
	godotenv.Load()

	// Initialize DatabaseHandler
	dbHandler := &DatabaseHandler{}
	err := dbHandler.Initialize()
	if err != nil {
		panic(err)
	}

	// Initialize WhatsMeowController
	meowController := &WhatsMeowController{
		DatabaseHandler: dbHandler,
	}
	// Start the WhatsMeow client
	go meowController.StartClient()

	staticFiles, err := fs.Sub(embeddedFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	fileServer := http.FileServer(http.FS(staticFiles))

	http.Handle("/", http.StripPrefix("/", fileServer))

	// Define HTTP handler functions
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		meowController.Connect()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Connected successfully."}`)
	})

	http.HandleFunc("/disconnect", func(w http.ResponseWriter, r *http.Request) {
		meowController.Disconnect()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Disconnected successfully."}`)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		meowController.Logout()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Logged out successfully."}`)
	})

	http.HandleFunc("/qrcode", func(w http.ResponseWriter, r *http.Request) {
		qr := meowController.GetQr()
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"QRCode": "%s", "message": "QR Code retrieved successfully."}`, qr)
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		isConnected, isLoggedIn, err := meowController.GetStatus()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"message": "Failed to get status."}`)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Status retrieved successfully.", "data": {"Connected": %t, "LoggedIn": %t}}`, isConnected, isLoggedIn)
	})

	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"message": "Method not allowed."}`)
			return
		}

		var body map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"message": "Invalid request body."}`)
			return
		}
		phone, ok := body["Phone"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"message": "Invalid phone number."}`)
			return
		}
		message, ok := body["Message"].(string)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"message": "Invalid message."}`)
			return
		}
		meowController.SendTextMessage(phone, message)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"message": "Message sent successfully."}`)
	})

	// Set up graceful shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000" // Default port if not specified
	}
	go func() {
		if err := http.ListenAndServe(":"+port, nil); err != nil {
			fmt.Println("Error starting server:", err)
		}
	}()

	// Beri jeda waktu sejenak untuk memastikan server siap
	time.Sleep(1 * time.Second) // Jeda 1 detik

	// Buka browser ke alamat targetURL
	err = openbrowser("http://localhost:" + port)
	if err != nil {
		log.Printf("Gagal membuka browser: %v", err)
		// Aplikasi akan terus berjalan tanpa browser terbuka
	}

	// Graceful shutdown on interrupt
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\nShutting down gracefully...")
}

// openbrowser mencoba membuka URL di browser default
func openbrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin": // macOS
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("sistem operasi tidak didukung untuk membuka browser: %s", runtime.GOOS)
	}
	return err
}
