package components

import (
	"go_todo_application2/database" //connect to db
	"net/http"                      //for handling HTTP req. and responses

	"github.com/gin-gonic/gin" //GIN fw. for routing btwn web pgs
)

func RegisterTimelineRoutes(r *gin.Engine) { //fn that registers timelineroutes
	r.GET("/timeline", showTimeline)       //when user visits /timeline it shows timeline pg
	r.POST("/timeline/edit", editTask)     //to edit a task frm a timeline
	r.POST("/timeline/delete", deleteTask) //delete a task frm timeline
	r.POST("/timeline/pin", pinTask)       //to pin a task in timeline
}
func showTimeline(c *gin.Context) { //fn that shows timeline pg

	email, _ := c.Cookie("user_email") //read user's mail stored in cookie

	if email == "" {
		c.Redirect(http.StatusSeeOther, "/") //it no email found user is not logged in and send them back to login pg
		return                               //stop the fn
	}

	rows, _ := database.DB.Query(`
        SELECT id, task, created_at, pinned 
        FROM todos 
        WHERE user_id = (SELECT id FROM users WHERE username = ?)
        ORDER BY pinned DESC, created_at DESC
    `, email)

	var data []map[string]any //create variables to temporarily store each col

	for rows.Next() { //loop thro each task in row
		var id, pinned int                 //task id, pinned or not
		var t, time string                 //task text, time created at
		rows.Scan(&id, &t, &time, &pinned) //put row data into variables
		data = append(data, gin.H{         //add this row into our list
			"id":     id,
			"task":   t,
			"time":   time,
			"pinned": pinned,
		})
	}

	c.HTML(http.StatusOK, "timeline.html", gin.H{"timeline": data}) //show timeline.html and pass all tasks inside it
}

func editTask(c *gin.Context) { //fn for editing a task
	email, _ := c.Cookie("user_email") //check if user is logged in
	if email == "" {
		c.Redirect(http.StatusSeeOther, "/") //if no user is logged i.e no email is put redirect user to login pg
		return                               //stop this fn
	}

	id := c.PostForm("id")     //which task is to edit
	task := c.PostForm("task") //updated task text

	// Update only if the task belongs to logged in user
	res, err := database.DB.Exec(
		`UPDATE todos SET task = ? WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)`,
		task, id, email, //new text, task id, user email
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error updating task") //if smtg went wrong it'll give error in updating the task
		return                                                          //stop the fn
	}
	// Optionally check rows affected to ensure the update happened (for safety)
	_ = res
	c.Redirect(http.StatusSeeOther, "/timeline") //after editing reload tmeline
}

func deleteTask(c *gin.Context) { //fn for deleting task
	email, _ := c.Cookie("user_email") //check if logged in
	if email == "" {
		c.Redirect(http.StatusSeeOther, "/") //if the email is empty.. redirect user to login pg
		return                               //stop fn
	}

	id := c.PostForm("id")      //task ID to delete
	_, err := database.DB.Exec( //delete only if this task belongs to this user
		`DELETE FROM todos WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)`,
		id, email, //which task? of which user to delete??
	)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error deleting task") //if error happens in deleting task
		return                                                          //stop fn
	}
	c.Redirect(http.StatusSeeOther, "/timeline") //refresh timeline after deleting
}

func pinTask(c *gin.Context) { //fn to pin task
	email, _ := c.Cookie("user_email") //check login
	if email == "" {
		c.Redirect(http.StatusSeeOther, "/") //again if email is nil...
		return                               //stop fn running
	}

	id := c.PostForm("id") //task id

	var count int // count pinned tasks only for this user
	database.DB.QueryRow("SELECT COUNT(*) FROM todos WHERE pinned=1 AND user_id = (SELECT id FROM users WHERE username = ?)", email).Scan(&count)

	var current int //check if this task is already pinned or not
	database.DB.QueryRow("SELECT pinned FROM todos WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", id, email).Scan(&current)

	if current == 1 { //if task is already pinned unpin it
		database.DB.Exec("UPDATE todos SET pinned=0 WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", id, email)
	} else {
		if count >= 3 { //user can pin only 3 tasks if more than 3 show this popup
			c.String(http.StatusBadRequest, "You can pin a maximum of 3 tasks only.")
			return
		}
		database.DB.Exec("UPDATE todos SET pinned=1 WHERE id = ? AND user_id = (SELECT id FROM users WHERE username = ?)", id, email)
	}
	c.Redirect(http.StatusSeeOther, "/timeline") //reload timeline after pin/unpin
}
