// autochrone-mini: mini collaborative web goal and time tracker
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"prettyDate": func(t time.Time) string {
			return t.Format("2 Jan. 2006")
		},
		"prettyDateTime": func(t time.Time) string {
			return t.Format("2 Jan. 2006 15h04")
		},
		"formDate": func(t time.Time) string {
			return t.Format("2006-01-02")
		},
		"prettyMinutes": func(minutes int) string {
			return fmt.Sprintf("%vh%02dm", minutes/60, minutes%60)
		},
	})

	r.LoadHTMLGlob("./templates/*.html")

	r.GET("/", RootGET)
	r.POST("/", RootPOST)
	r.GET("/:pslug", ProjectGET)
	r.POST("/:pslug", ProjectPOST)

	r.Run("autochrone.herokuapp.com:443")
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
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	project := NewProject(slug, name, description, targetDate)
	if project == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	user := project.NewUser(admin, true)
	if user == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/%s", project.Slug))
}

// ProjectGET read-only project or user write access (add notes, delete own notes, with ?auth=str)
func ProjectGET(c *gin.Context) {
	project := GetProjectBySlug(c.Param("pslug"))
	if project == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	c.HTML(http.StatusOK, "project-page", gin.H{
		"Project": project,
		"User":    project.GetUserBySlug(c.Query("auth")),
	})
}

// ProjectPOST add note, user or measure
func ProjectPOST(c *gin.Context) {
	project := GetProjectBySlug(c.Param("pslug"))
	if project == nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}
	user := project.GetUserBySlug(c.Query("auth"))

	switch c.PostForm("query") {
	case "new-user":
		username := c.PostForm("name")
		if user.IsAdmin && username != "" && len(username) <= 140 {
			project.NewUser(username, c.PostForm("is-admin") == "on")
		}
	case "new-measure":
		code := c.PostForm("code")
		name := c.PostForm("name")
		unit := c.PostForm("unit")
		goalDirection := c.PostForm("goal-direction")
		goal, err := strconv.Atoi(c.PostForm("goal"))

		if user.IsAdmin && code != "" && len(code) <= 42 && name != "" && len(name) <= 140 && unit != "" && len(unit) <= 42 {
			if goalDirection == "none" || ((goalDirection == "min" || goalDirection == "max") && err == nil) {
				project.NewMeasure(code, name, unit, goalDirection, goal)
			}
		}
	case "new-note":
		comment := c.PostForm("comment")
		if comment != "" && len(comment) <= 1000 {
			var err error
			measuresValues := make(map[int]int)
			for measureID := range project.Measures {
				measuresValues[measureID], err = strconv.Atoi(c.PostForm(fmt.Sprintf("measure-%v", measureID)))
				if err != nil {
					delete(measuresValues, measureID)
				}
			}
			if err == nil {
				project.NewNote(user.ID, comment, measuresValues)
			}
		}
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("/%s?auth=%s", project.Slug, user.Slug))
}
