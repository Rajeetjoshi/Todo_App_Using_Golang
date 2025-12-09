package components

import (
	"go_todo_application2/database" //to import db conncetion
	"net/http"                      //used for HTTP request and responses
	"strings"

	"github.com/gin-gonic/gin" //GIN fw. for handling requests
)

func RegisterHomeRoutes(r *gin.Engine) { //fn to register all /home routes

	r.GET("/home", showHome) //when the new user opens /home it calls showHome fn
}

func showHome(c *gin.Context) { //this fn runs when we open /home

	name, _ := c.Cookie("user_name")           //for reading user's name
	email, _ := c.Cookie("user_email")         //for reading user's email
	birthdate, _ := c.Cookie("user_birthdate") //for reading user's birthdate

	if email == "" { //if no email is inputted.. no user logged in
		c.Redirect(http.StatusSeeOther, "/") //redirect user back to home/login pg
		return                               //stop the fn running
	}

	rows, _ := database.DB.Query( //runs below SQL query to get user's todo list
		"SELECT id, task FROM todos WHERE user_id = (SELECT id FROM users WHERE username = ?)",
		email, //pass email into SQL query
	) //this closing bracket for closing the function if we keep the function open it'll summon the error immediately

	var todos []map[string]any //create an empty list to store todos

	for rows.Next() { //loop thro each row
		var id int                                           //store todo id
		var task string                                      //store todo next
		rows.Scan(&id, &task)                                //read the row data into id and task
		todos = append(todos, gin.H{"id": id, "task": task}) //add todo item in list
	}

	c.HTML(http.StatusOK, "home.html", gin.H{ //render home.html page
		"user":      name, // send user's name to HTML
		"Name":      name,
		"Email":     email,
		"Birthdate": birthdate,
		"todos":     todos, //send all todo items to HTML

	})
}
func RegisterProfileRoutes(r *gin.Engine) {

	// 1️⃣ Show the Edit Name page
	r.GET("/edit-name", func(c *gin.Context) {

		email, err := c.Cookie("user_email")
		if err != nil || email == "" {
			c.String(400, "User not logged in")
			return
		}

		// get current name
		var currentName string
		err = database.DB.QueryRow("SELECT name FROM users WHERE username = ?", email).Scan(&currentName)
		if err != nil {
			c.String(500, "Error fetching name: "+err.Error())
			return
		}

		// render HTML page
		c.HTML(200, "edit_name.html", gin.H{
			"CurrentName": currentName,
		})
	})

	// 2️⃣ Handle Save button
	r.POST("/update-name", func(c *gin.Context) {

		newName := c.PostForm("new_name")
		if newName == "" {
			c.String(400, "Name cannot be empty")
			return
		}

		email, err := c.Cookie("user_email")
		if err != nil {
			c.String(400, "Not logged in")
			return
		}

		_, err = database.DB.Exec("UPDATE users SET name=? WHERE username=?", newName, email)
		if err != nil {
			c.String(500, "DB error: "+err.Error())
			return
		}

		// update cookie
		maxAge := 60 * 60 * 24 * 7
		c.SetCookie("user_name", newName, maxAge, "/", "", false, true)

		// redirect back to home
		c.Redirect(303, "/home")
	})

	// keep your existing update profile route
	r.POST("/update-profile", UpdateProfile)
}

func UpdateProfile(c *gin.Context) {

	oldEmail, _ := c.Cookie("user_email")

	newName := c.PostForm("new_name")
	newEmail := c.PostForm("new_email")
	newBirthdate := c.PostForm("birthdate")

	var exists int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = ?", newEmail).Scan(&exists)

	if err == nil && newEmail != oldEmail {
		c.String(http.StatusBadRequest, "This email already exists. Try another.")
		return
	}

	_, err = database.DB.Exec(`
        UPDATE users
        SET username = ?, name = ?, birthdate = ?
        WHERE username = ?
    `, newEmail, newName, newBirthdate, oldEmail)

	if err != nil {
		c.String(http.StatusInternalServerError, "Error updating profile: "+err.Error())
		return
	}

	maxAge := 60 * 60 * 24 * 7
	c.SetCookie("user_email", newEmail, maxAge, "/", "", false, true)
	c.SetCookie("user_name", newName, maxAge, "/", "", false, true)
	c.SetCookie("user_birthdate", newBirthdate, maxAge, "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/home")
}

// GET /edit-name -> separate page showing current name + input for new name
func showEditName(c *gin.Context) {
	email, _ := c.Cookie("user_email")
	if email == "" {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	name, _ := c.Cookie("user_name")

	c.HTML(http.StatusOK, "edit_name.html", gin.H{
		"CurrentName": name,
	})
}

// POST /update-name -> update DB + cookie then redirect back home
func handleUpdateName(c *gin.Context) {
	newName := strings.TrimSpace(c.PostForm("new_name"))
	if newName == "" {
		c.Redirect(http.StatusSeeOther, "/edit-name?error=empty")
		return
	}

	email, err := c.Cookie("user_email")
	if err != nil || email == "" {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	_, err = database.DB.Exec("UPDATE users SET name = ? WHERE username = ?", newName, email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error updating name: "+err.Error())
		return
	}

	maxAge := 60 * 60 * 24 * 7
	c.SetCookie("user_name", newName, maxAge, "/", "", false, true)

	c.Redirect(http.StatusSeeOther, "/home?updated=name")
}
