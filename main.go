package main //main package declared

import ( //import pulls intrnal and external modules need for ur project
	"go_todo_application2/components" //this imports all component's file
	"go_todo_application2/database"   //this imports db.go file

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin" //this imports GIN framework very popular for building web servers
)

func main() {
	database.InitDB()             //calls fn frm ur db.
	r := gin.Default()            //'r' is router.. creates a new GIN router
	r.LoadHTMLGlob("templates/*") //tells GIN to load all .html files frm templates folder

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("sessions", store))

	components.RegisterLoginRoutes(r)  //calls a fn. defined in each file inside components
	components.RegisterSignupRoutes(r) //'r' in this is Router as declared abv
	components.RegisterHomeRoutes(r)   //
	components.RegisterAddRoutes(r)    //"hey rpoter pls handle request to /add a task"
	components.RegisterDeleteRoutes(r)
	components.RegisterTimelineRoutes(r)
	components.RegisterProfileRoutes(r)
	components.RegisterLogoutRoutes(r)

	r.Run(":8080") //starts ur web server on port 8080
}
