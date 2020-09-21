package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"log"
	"time"
)

var db *sqlx.DB

func init() {
	// connect to database
	var err error
	db, err = sqlx.Open("postgres", "user=autochrone password=autochrone dbname=achr sslmode=disable")
	if err != nil {
		panic("could not connect to database")
	}
}

// Project an autochrone project
type Project struct {
	ID           int       `db:"id"`
	Slug         string    `db:"slug"`
	Name         string    `db:"name"`
	Description  string    `db:"description"`
	CreationDate time.Time `db:"creation_date"`
	TargetDate   time.Time `db:"target_date"`
}

// User an autochrone user
type User struct {
	ProjectID    int       `db:"project_id"`
	ID           int       `db:"id"`
	IsAdmin      bool      `db:"is_admin"`
	Slug         string    `db:"slug"`
	Name         string    `db:"name"`
	CreationDate time.Time `db:"creation_date"`
}

// NewProject creates a new project in the database
func NewProject(slug, name, description string, targetDate time.Time) *Project {
	creationDate := time.Now()
	row := db.QueryRowx("insert into projects (slug, name, description, creation_date, target_date) values ($1, $2, $3, $4, $5) returning id", slug, name, description, creationDate, targetDate)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Print(err)
		return nil
	}
	return &Project{
		ID:           id,
		Slug:         slug,
		Name:         name,
		Description:  description,
		CreationDate: creationDate,
		TargetDate:   targetDate,
	}
}

// AddUser adds a user to an existing project
func (p *Project) NewUser(name string, isAdmin bool) *User {
	slug := "anursiet"
	creationDate := time.Now()
	row := db.QueryRowx("insert into users (project_id, is_admin, slug, name, creation_date) values ($1, $2, $3, $4, $5) returning id", p.ID, isAdmin, slug, name, creationDate)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Print(err)
		return nil
	}
	return &User{
		ProjectID:    p.ID,
		ID:           id,
		IsAdmin:      isAdmin,
		Slug:         slug,
		Name:         name,
		CreationDate: creationDate,
	}
}

// GetProjectBySlug returns an existing project
func GetProjectBySlug(slug string) *Project {
	row := db.QueryRowx("select id, slug, name, description, creation_date, target_date from projects where slug = $1", slug)
	p := &Project{}
	if err := row.StructScan(p); err != nil {
		log.Print(err)
		return nil
	}
	return p
}
