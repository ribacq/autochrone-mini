package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"log"
	//"os"
	"time"
)

var db *sqlx.DB

func init() {
	// connect to database
	var err error
	// for local tests
	db, err = sqlx.Open("postgres", "user=autochrone password=autochrone dbname=achr sslmode=disable")
	//db, err = sqlx.Open("postgres", os.Getenv("DATABASE_URL"))
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
	Users        map[int]*User
	Measures     map[int]*Measure
	Tasks        map[int]*Task
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

// Task an autochrone project task
type Task struct {
	ID           int       `db:"id"`
	ProjectID    int       `db:"project_id"`
	Description  string    `db:"description"`
	CreationDate time.Time `db:"creation_date"`
	EndDate      time.Time `db:"end_date"`
	Cancelled    bool      `db:"cancelled"`
	Color        string    `db:"color"`
}

// Measure an autochrone measure
type Measure struct {
	ID            int    `db:"id"`
	ProjectID     int    `db:"project_id"`
	Code          string `db:"code"`
	Name          string `db:"name"`
	Unit          string `db:"unit"`
	GoalDirection string `db:"goal_direction"`
	Goal          int    `db:"goal"`
}

// HasMin m.GoalDirection == "min"
func (m *Measure) HasMin() bool { return m.GoalDirection == "min" }

// HasMax m.GoalDirection == "max"
func (m *Measure) HasMax() bool { return m.GoalDirection == "max" }

// Note an autochrone note
type Note struct {
	ID             int       `db:"id"`
	ProjectID      int       `db:"project_id"`
	UserID         int       `db:"user_id"`
	TaskID         int       `db:"task_id"`
	CreationDate   time.Time `db:"creation_date"`
	Comment        string    `db:"comment"`
	MeasuresValues map[int]int
}

// NewProject creates a new project in the database
func NewProject(slug, name, description string, targetDate time.Time) *Project {
	creationDate := time.Now().UTC()
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
		Users:        make(map[int]*User),
	}
}

// NewUser adds a user to an existing project
func (p *Project) NewUser(name string, isAdmin bool) *User {
	slug := RandomSlug()
	creationDate := time.Now().UTC()
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
	p.Users[u.ID] = u
	return u
}

// NewMeasure adds a measure to an existing project
func (p *Project) NewMeasure(code, name, unit, goalDirection string, goal int) *Measure {
	row := db.QueryRowx("insert into measures (project_id, code, name, unit, goal_direction, goal) values ($1, $2, $3, $4, $5, $6) returning id", p.ID, code, name, unit, goalDirection, goal)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Print(err)
		return nil
	}
	m := &Measure{
		ID:            id,
		ProjectID:     p.ID,
		Code:          code,
		Name:          name,
		Unit:          unit,
		GoalDirection: goalDirection,
		Goal:          goal,
	}
	p.Measures[m.ID] = m
	return m
}

// NewNote adds a note to an existing project
func (p *Project) NewNote(userID int, comment string, measuresValues map[int]int) *Note {
	creationDate := time.Now().UTC()
	row := db.QueryRowx("insert into notes (project_id, user_id, creation_date, comment) values ($1, $2, $3, $4) returning id", p.ID, userID, creationDate, comment)
	var id int
	if err := row.Scan(&id); err != nil {
		log.Print(err)
		return nil
	}
	n := &Note{
		ID:             id,
		ProjectID:      p.ID,
		CreationDate:   creationDate,
		Comment:        comment,
		MeasuresValues: make(map[int]int),
	}
	for measureID, measureValue := range measuresValues {
		if _, ok := p.Measures[measureID]; ok {
			row = db.QueryRowx("insert into notes_measures_values (note_id, measure_id, value) values ($1, $2, $3) returning id", n.ID, measureID, measureValue)
			var id int
			if err := row.Scan(&id); err != nil {
				log.Print(err)
				return nil
			}
			n.MeasuresValues[measureID] = measureValue
		}
	}
	p.Notes = append([]*Note{n}, p.Notes...)
	return n
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
	if err := p.FetchTasks(); err != nil {
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
	p.Users = make(map[int]*User)
	for rows.Next() {
		u := &User{}
		if err = rows.StructScan(u); err != nil {
			return err
		}
		p.Users[u.ID] = u
	}
	return nil
}

// FetchMeasures loads measures for a project
func (p *Project) FetchMeasures() error {
	rows, err := db.Queryx("select id, project_id, code, name, unit, goal_direction, goal from measures where project_id = $1", p.ID)
	if err != nil {
		return err
	}
	p.Measures = make(map[int]*Measure)
	for rows.Next() {
		m := &Measure{}
		if err := rows.StructScan(m); err != nil {
			return err
		}
		p.Measures[m.ID] = m
	}
	return nil
}

// FetchTasks loads tasks for a project
func (p *Project) FetchTasks() error {
	rows, err := db.Queryx("select id, project_id, description, creation_date, end_date, cancelled, color from tasks where project_id = $1", p.ID)
	if err != nil {
		return err
	}
	p.Tasks = make(map[int]*Task)
	for rows.Next() {
		t := &Task{}
		if err := rows.StructScan(t); err != nil {
			return err
		}
		p.Tasks[t.ID] = t
	}
	return nil
}

// FetchNotes loads notes for a project
func (p *Project) FetchNotes() error {
	rows, err := db.Queryx("select id, project_id, user_id, task_id, creation_date, comment from notes where project_id = $1 order by creation_date desc", p.ID)
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

// GetUserBySlug returns the project’s user with the given slug
func (p *Project) GetUserBySlug(slug string) *User {
	for _, user := range p.Users {
		if user.Slug == slug {
			return user
		}
	}
	return nil
}
