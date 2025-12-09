package components

import (
	"database/sql"
	"go_todo_application2/database"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterSignupRoutes(r *gin.Engine) {
	r.POST("/signup", handleSignup)
	r.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
}

func handleSignup(c *gin.Context) {
	user := strings.TrimSpace(strings.ToLower(c.PostForm("email")))
	pass := c.PostForm("password")
	name := strings.TrimSpace(c.PostForm("name"))

	// Empty field validation
	if user == "" || pass == "" || name == "" {
		c.String(http.StatusBadRequest, "Username, name or password cannot be empty.")
		return
	}

	// Email validation
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(user) {
		c.HTML(http.StatusBadRequest, "index.html", gin.H{
			"error": "Please enter a valid email address.",
		})
		return
	}

	// Check if the email already exists
	var id int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = ?", user).Scan(&id)
	if err == nil {
		c.HTML(http.StatusOK, "index.html", gin.H{"error": "User already exists! Please login instead."})
		return
	} else if err != sql.ErrNoRows && err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "Database error: " + err.Error()})
		return
	}

	// Store user
	_, err = database.DB.Exec("INSERT INTO users(username, password, name) VALUES(?, ?, ?)", user, pass, name)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "index.html", gin.H{"error": "Error creating account: " + err.Error()})
		return
	}

	// Cookies
	maxAge := 60 * 60 * 24 * 7
	c.SetCookie("user_email", user, maxAge, "/", "", false, true)
	c.SetCookie("user_name", name, maxAge, "/", "", false, true)

	// Redirect to home
	c.Redirect(http.StatusSeeOther, "/home?user_id="+user)
}
