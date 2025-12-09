# 📝 ToDo Application (Go + Gin Framework)

A simple and modern ToDo web application built using **Golang**, **Gin Framework**, and **MySQL**.  
It supports **user login**, **task creation**, **timeline view**, and **profile management** with a clean UI.

---

## 🚀 Features

### 🔐 User Authentication
- Login using email + password  
- Cookies used for user session handling  
- Basic validation & error responses  

### 🧑‍💼 User Profile Management
- Update **Name**
- Update **Email**
- Update **Birthdate**
- Inline edit UI using modals (clean modern UX)

### 📝 Task Management
- Add tasks
- Select date & time for each task
- Tasks appear in timeline in proper order

### 📅 Timeline Page
- Shows all tasks grouped by:
  - Today  
  - Tomorrow  
  - Upcoming
- Clean card-based UI

### 🎨 UI/UX Features
- Glassmorphism style cards  
- Responsive design  
- Modern navbar  
- Profile modal  
- Smooth animations  

---

## 🏗️ Tech Stack

| Layer | Technology |
|------|-------------|
| Backend | **Golang**, **Gin Framework** |
| Database | **MySQL** |
| Templating | HTML, Go Templates |
| Styling | CSS, Gradient UI |
| Authentication | Cookies |

---

## 📂 Project Structure
project/
│── main.go
│── go.mod
│── database/
│ └── db.go
│── components/
│ ├── login.go
│ ├── home.go
│ ├── timeline.go
│ └── profile.go
│── templates/
│ ├── index.html
│ ├── home.html
│ ├── timeline.html
│ └── components...
│── static/
└── css /
💡 Author

Your Name
GitHub: https://github.com/Rajeetjoshii

