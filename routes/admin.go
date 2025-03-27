package routes

import "github.com/gin-gonic/gin"

func Admin(r *gin.Engine) {
	// TODO: implement RBAC checks for all routes here
	r.GET("/admin/users/:slug", func(c *gin.Context) {
		userSlug := c.Param("slug")

	})
}
