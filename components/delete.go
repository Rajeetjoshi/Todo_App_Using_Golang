package components //

import (
	"go_todo_application2/database" //package that holds DB connection
	"net/http"                      //HTTP constants that holds req and response

	"github.com/gin-gonic/gin" //GIN framework for handling HTTP req
)

func RegisterDeleteRoutes(r *gin.Engine) { //declares fn RegisterDeleteRoutes, fn that receives GIN router 'r'
	r.GET("/delete/:id", deleteTodo) //id is  placeholder
}

func deleteTodo(c *gin.Context) { //declares fn deleteToDo, GIN passes *gin.Conntext 'c'
	id := c.Param("id")        //reads the parameter named id frm URL
	user := c.Query("user_id") //reads query named user frm URL

	database.DB.Exec("DELETE FROM todos WHERE id=?", id)   //
	c.Redirect(http.StatusSeeOther, "/home?user_id="+user) //
}
