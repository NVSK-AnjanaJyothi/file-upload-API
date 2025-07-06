# 📁 File Upload API (Golang + JWT + SQLite)

This is a simple and secure file upload API built using **Go (Golang)**.  
It supports user login via **JWT tokens**, file uploading, storing metadata in **SQLite**, and retrieving uploaded files.

---

## 🚀 Features

- ✅ JWT-based user authentication
- 📤 File upload with multipart/form-data
- 🗃️ Save file metadata to SQLite
- 🔐 Protected API endpoints
- 📂 Download uploaded files
- 🧹 Delete uploaded files
- 💡 Clean, modular codebase (`main.go`, `db.go`, `auth.go`)

---

## 🛠 Tech Stack

- Go (Golang)
- Gorilla Mux (Router)
- SQLite (lightweight embedded DB)
- JWT (Authentication)
- Postman / curl (for testing)

---

## 📁 Project Structure

file-upload-api/
├── main.go # Handles routes and upload logic
├── db.go # Database connection & helper functions
├── auth.go # JWT generation & middleware
├── uploads/ # Folder where files are stored
├── go.mod / go.sum # Go module files


---

## 🔐 Authentication

1. **Login API**:
## 🔐 Authentication

1. **Login API**:
POST /login

css
Copy
Edit
Example Body:
```json
{
  "username": "admin",
  "password": "password"
}
Returns: JWT Token

Use the token in headers:

makefile
Copy
Edit
Authorization: Bearer <your-token>
📤 File Upload
Endpoint: POST /upload

Header: Authorization: Bearer <token>

Body Type: form-data

Key: file

Value: (select a file)

📥 Fetch Files
Get all files: GET /files

Get specific file: GET /file/{filename}

Delete file: DELETE /file/{filename}

🧪 Testing
You can test using:

Postman

curl

Default Credentials
Username	Password
admin	password

📌 Notes
Make sure uploads/ folder exists. If not, it will be created automatically.

File data is stored in SQLite (data.db)

