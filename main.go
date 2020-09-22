// autochrone-mini: mini collaborative web goal and time tracker
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"prettyDate": func(t time.Time) string {
			return t.Format("2 Jan. 2006")
		},
		"formDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
	})

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
	name := c.PostForm("name")
	slug := c.PostForm("slug")
	admin := c.PostForm("admin")
	description := c.PostForm("description")
	targetDate, err := time.Parse("2006-01-02", c.PostForm("target-date"))

	if name == "" || slug == "" || admin == "" || len(name) > 140 || len(name) > 140 || len(admin) > 140 || len(description) > 1000 || err != nil {
		c.Redirect(http.StatusBadRequest, "/")
		return
	}

	project := NewProject(slug, name, description, targetDate)
	if project == nil {
		c.Redirect(http.StatusInternalServerError, "/")
		return
	}
	user := project.NewUser(admin, true)
	if user == nil {
		c.Redirect(http.StatusInternalServerError, "/")
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/%s", project.Slug))
}

// ProjectGET read-only project or user write access (add notes, delete own notes, with ?auth=str)
func ProjectGET(c *gin.Context) {
	project := GetProjectBySlug(c.Param("pslug"))
	c.HTML(http.StatusOK, "project-read-only-page", gin.H{
		"Project": project,
	})
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
