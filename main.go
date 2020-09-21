// autochrone-mini: mini collaborative web goal and time tracker
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("./templates/*.html")

	r.GET("/", RootGET)
	r.POST("/", RootPOST)
	r.GET("/:pslug", ProjectGET)
	r.POST("/:pslug", ProjectPOST)
	r.GET("/:pslug/admin", ProjectAdminGET)
	r.POST("/:pslug/admin", ProjectAdminPOST)

	r.Run(":8080")
}

// RootGET project creation form + about
func RootGET(c *gin.Context) {
	c.HTML(http.StatusOK, "home-page", nil)
}

// RootPOST project creation and redirect to project
func RootPOST(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

// ProjectGET read-only project or user write access (add notes, delete own notes, with ?auth=str)
func ProjectGET(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

// ProjectPOST add note
func ProjectPOST(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

// ProjectAdminGET project admin
func ProjectAdminGET(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}

// ProjectAdminGET update project
func ProjectAdminPOST(c *gin.Context) {
	c.Redirect(http.StatusFound, "/")
}
