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
	Users        []*User
	Measures     []*Measure
	Notes        []*Note
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

// Measure an autochrone measure
type Measure struct {
	ID          int    `db:"id"`
	ProjectID   int    `db:"project_id"`
	Code        string `db:"code"`
	Name        string `db:"name"`
	Unit        string `db:"unit"`
	HasGoalLow  bool
	GoalLow     int
	HasGoalHigh bool
	GoalHigh    int
}

// Note an autochrone note
type Note struct {
	ID             int       `db:"id"`
	ProjectID      int       `db:"project_id"`
	UserID         int       `db:"user_id"`
	CreationDate   time.Time `db:"creation_date"`
	Comment        string    `db:"comment"`
	MeasuresValues map[int]int
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
		Users:        []*User{},
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
	u := &User{
		ProjectID:    p.ID,
		ID:           id,
		IsAdmin:      isAdmin,
		Slug:         slug,
		Name:         name,
		CreationDate: creationDate,
	}
	p.Users = append(p.Users, u)
	return u
}

// GetProjectBySlug returns an existing project
func GetProjectBySlug(slug string) *Project {
	row := db.QueryRowx("select id, slug, name, description, creation_date, target_date from projects where slug = $1", slug)
	p := &Project{}
	if err := row.StructScan(p); err != nil {
		log.Print(err)
		return nil
	}
	if err := p.FetchUsers(); err != nil {
		log.Print(err)
		return nil
	}
	if err := p.FetchMeasures(); err != nil {
		log.Print(err)
		return nil
	}
	if err := p.FetchNotes(); err != nil {
		log.Print(err)
		return nil
	}
	return p
}

// FetchUsers loads users for a project
func (p *Project) FetchUsers() error {
	rows, err := db.Queryx("select id, project_id, is_admin, slug, name, creation_date from users where project_id = $1", p.ID)
	if err != nil {
		return err
	}
	p.Users = []*User{}
	for rows.Next() {
		u := &User{}
		if err = rows.StructScan(u); err != nil {
			return err
		}
		p.Users = append(p.Users, u)
	}
	return nil
}

// FetchMeasures loads measures for a project
func (p *Project) FetchMeasures() error {
	rows, err := db.Queryx("select id, project_id, code, name, unit from measures where project_id = $1", p.ID)
	if err != nil {
		return err
	}
	p.Measures = []*Measure{}
	for rows.Next() {
		m := &Measure{}
		if err := rows.StructScan(m); err != nil {
			return err
		}
		row := db.QueryRowx("select goal_low, goal_high from measures_goals where measure_id = $1", m.ID)
		var gL, gH *int
		if err := row.Scan(&gL, &gH); err == nil {
			if gL != nil {
				m.HasGoalLow = true
				m.GoalLow = *gL
			}
			if gH != nil {
				m.HasGoalHigh = true
				m.GoalHigh = *gH
			}
		}
		p.Measures = append(p.Measures, m)
	}
	return nil
}

// FetchNotes loads notes for a project
func (p *Project) FetchNotes() error {
	rows, err := db.Queryx("select id, project_id, user_id, creation_date, comment from notes where project_id = $1", p.ID)
	if err != nil {
		return err
	}
	p.Notes = []*Note{}
	for rows.Next() {
		n := &Note{
			MeasuresValues: make(map[int]int),
		}
		if err := rows.StructScan(n); err != nil {
			return err
		}
		valuesRows, err := db.Queryx("select measure_id, value from notes_measures_values where note_id = $1", n.ID)
		if err != nil {
			return err
		}
		var measureId, value int
		for valuesRows.Next() {
			if err := valuesRows.Scan(&measureId, &value); err != nil {
				return err
			}
			n.MeasuresValues[measureId] = value
		}
		p.Notes = append(p.Notes, n)
	}
	return nil
}
