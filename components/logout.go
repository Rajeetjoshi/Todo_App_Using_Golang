package components

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterLogoutRoutes(r *gin.Engine) {
	r.GET("/logout", func(c *gin.Context) {

		// delete cookies
		c.SetCookie("user_email", "", -1, "/", "", false, true)
		c.SetCookie("user_name", "", -1, "/", "", false, true)
		c.SetCookie("user_birthdate", "", -1, "/", "", false, true)

		// redirect to login page
		c.Redirect(http.StatusSeeOther, "/")
	})
}
