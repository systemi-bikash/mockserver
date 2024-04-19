package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Struct to represent the data in the JSON file
type Data struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}

func main() {
	// Defining HTTP handler functions
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/", fileHandler)

	// Start the HTTP server
	fmt.Println("Mock server running on port 8080...")
	http.ListenAndServe(":8080", nil)
}

// Handler function to create a JSON file with provided data
func createHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var requestData Data
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the filename from the request
	filename := "data.json"
	if requestData.Status != "" {
		filename = requestData.Status + ".json"
	}

	// Construct path to the data folder
	dataFolderPath := "data"
	dataFilePath := filepath.Join(dataFolderPath, filename)

	// Create the data folder if it doesn't exist
	err = os.MkdirAll(dataFolderPath, 0755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or overwrite the JSON file
	file, err := os.Create(dataFilePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Encode and write data to the JSON file
	err = json.NewEncoder(file).Encode(requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "JSON file created successfully: %s", filename)
}

// Handler function to handle requests for JSON files
func fileHandler(w http.ResponseWriter, r *http.Request) {
	// Extract requested filename from the URL path
	filename := strings.TrimPrefix(r.URL.Path, "/")
	if filename == "" {
		// If no filename is provided, return an error
		http.Error(w, "Please provide a filename", http.StatusBadRequest)
		return
	}

	// Construct path to the data folder
	dataPath := filepath.Join("data", filename)

	// Read data from the requested JSON file
	file, err := os.ReadFile(dataPath)
	if err != nil {
		// If file not found or other error, return 404 Not Found
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Write file content to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}
