package components

import (
	"go_todo_application2/database" //connects to db
	"net/http"                      //handles HTTP req and response
	"time"                          //go library helps u work w current time, date while adding the tasks

	"github.com/gin-gonic/gin" //GIN fw. for routing btwn tasks
)

func RegisterAddRoutes(r *gin.Engine) { //fn for registering /add routes w ur GIN router
	r.POST("/add", addTodo) //when user submits add task form, this fn will run
}

func addTodo(c *gin.Context) { //fn for handling adding a new todo task into db
	email, err := c.Cookie("user_email") //reads logged in user's email
	if err != nil || email == "" {       //if coolkie is missing or  error occurs it tells user is not logged in
		c.Redirect(http.StatusSeeOther, "/") //send user to login pg
		return
	}

	task := c.PostForm("task")            //read the task text frm HTML form input
	taskDate := c.PostForm("task_date")   // reads the task date.. "YYYY-MM-DD"
	taskTime := c.PostForm("task_time")   // reads the task time.. "HH:MM"
	var createdAt string                  //this var will store final timestamp of task
	if taskDate != "" && taskTime != "" { //if the date and time isn't empty
		createdAt = taskDate + " " + taskTime + ":00" //if user selected both date & time, use them to form full datetime
	} else { //if user didn't select date/time
		loc, _ := time.LoadLocation("Asia/Kolkata")
		createdAt = time.Now().In(loc).Format("2006-01-02 15:04:05") //if user not selected date/time user current indian time
	}

	// Insert using user_id by selecting from users table using username/email
	_, err = database.DB.Exec(
		"INSERT INTO todos(user_id, task, created_at) VALUES((SELECT id FROM users WHERE username = ?), ?, ?)",
		email, task, createdAt,
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error adding task: "+err.Error()) //if any db error occurs, print that the problem arisen
		return
	}

	c.Redirect(http.StatusSeeOther, "/home?added=1") //after adding the task successfully redirect user back to /home.. "?added=1" is used to show popup "Added successfully!"
}
