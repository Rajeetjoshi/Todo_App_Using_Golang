package components

import (
	"go_todo_application2/database" //import db connection file
	"net/http"                      //used for HTTP responses and request
	"strings"                       //it contains useful fns for working w text.. remove spaces, split text, join text, replacee text

	"github.com/gin-gonic/gin" //GIN fw. for routes, req, responses
)

func RegisterLoginRoutes(r *gin.Engine) { //fn to register all login related routes
	r.GET("/", showLogin)         //when user goes to "/" show login page
	r.GET("/login", showLogin)    //when user goes to "/login" also show login pg
	r.POST("/login", handleLogin) //when user submits login form call handleLogin()
}

func showLogin(c *gin.Context) { //thisfn shows the login pg.. *gin.Context means it points to gin.Context, in gin, Context holds everything abt current HTTP req and response
	c.HTML(http.StatusOK, "index.html", nil) //render index.html page
}

func handleLogin(c *gin.Context) { //this fn runs when user submits login form
	user := strings.TrimSpace(strings.ToLower(c.PostForm("email"))) //read email frm login form, convertit to lowercase and remove extra spaces
	pass := c.PostForm("password")                                  //read the pw. exactly as entered

	var id int                                                                                                //to store user id
	var name string                                                                                           //to store user name frm db
	row := database.DB.QueryRow("SELECT id, name FROM users WHERE username = ? AND password = ?", user, pass) //query that checks if a user with this email and pw. exists
	err := row.Scan(&id, &name)                                                                               //scan the result
	if err != nil {
		c.String(http.StatusUnauthorized, "Invalid credentials") //if scan fails it means invalid email/pw.
		return                                                   //stop the fn
	}

	// if we reach here means login successful
	maxAge := 60 * 60 * 24 * 7 //how long a cookie shud stay in the browser (7 days)
	c.SetCookie("user_email", user, maxAge, "/", "", false, false)
	c.SetCookie("user_name", name, maxAge, "/", "", false, false)

	c.Redirect(http.StatusSeeOther, "/home") //after login pg. redirect to /home pg
}
