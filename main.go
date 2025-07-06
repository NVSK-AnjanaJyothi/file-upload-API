package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	initDB()
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/files", listFiles).Methods("GET")
	router.HandleFunc("/download/{filename}", downloadFile).Methods("GET")

	// Protected routes
	secured := router.PathPrefix("/").Subrouter()
	secured.Use(authMiddleware)
	secured.HandleFunc("/upload", uploadFile).Methods("POST")
	secured.HandleFunc("/delete/{filename}", deleteFile).Methods("DELETE")
	secured.HandleFunc("/file-info", getFileInfo).Methods("GET")

	fmt.Println("Server running at http://localhost:8083")
	log.Fatal(http.ListenAndServe(":8083", router))
}

// Upload file handler
func uploadFile(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	saveFile, err := os.Create("uploads/" + handler.Filename)
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer saveFile.Close()

	size, err := io.Copy(saveFile, file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	err = saveFileToDB(handler.Filename, size)
	if err != nil {
		http.Error(w, "Failed to save file info to database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded and saved to database successfully.\n", handler.Filename)
}

// Download file handler
func downloadFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filename := params["filename"]
	filePath := "uploads/" + filename

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	http.ServeFile(w, r, filePath)
}

// List uploaded file names
func listFiles(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir("uploads")
	if err != nil {
		http.Error(w, "Unable to read uploads folder", http.StatusInternalServerError)
		return
	}

	var fileNames []string
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileNames)
}

// Delete a file
func deleteFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	filename := params["filename"]
	filePath := "uploads/" + filename

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s deleted successfully.", filename)
}

// Show full file info from database
type FileRecord struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	UploadedAt string `json:"uploaded_at"`
}

func getFileInfo(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, size, uploaded_at FROM files")
	if err != nil {
		http.Error(w, "Failed to fetch file info", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var files []FileRecord
	for rows.Next() {
		var file FileRecord
		err := rows.Scan(&file.ID, &file.Name, &file.Size, &file.UploadedAt)
		if err != nil {
			http.Error(w, "Error reading record", http.StatusInternalServerError)
			return
		}
		files = append(files, file)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}
