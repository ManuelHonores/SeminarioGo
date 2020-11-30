package person

import (
	"SeminarioGo/internal/config"

	"github.com/jmoiron/sqlx"
)

// Person ...
type Person struct {
	ID       int64
	Name     string
	Lastname string
	Age      int64
}

// Service ...
type Service interface {
	AddPerson(Person) error
	FindAll() []*Person
	UpdatePerson(int, Person) error
	DeletePerson(int) error
	FindByID(int) *Person
}

type service struct {
	db   *sqlx.DB
	conf *config.Config
}

// New ...
func New(db *sqlx.DB, c *config.Config) (Service, error) {
	return service{db, c}, nil
}

func (s service) FindAll() []*Person {
	var list []*Person
	if err := s.db.Select(&list, "SELECT * FROM persons"); err != nil {
		panic(err)
	}
	return list
}

func (s service) AddPerson(p Person) error {
	query := "INSERT INTO persons (name, lastname, age) VALUES (?,?,?)"
	statementInsert, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statementInsert.Exec(p.Name, p.Lastname, p.Age)
	if err != nil {
		return err
	}
	return nil
}

func (s service) FindByID(ID int) *Person {
	var Person Person
	query := "SELECT * FROM persons WHERE ID = ?"
	if err := s.db.Get(&Person, query, ID); err != nil { //Get es analogo a QueryRow
		return nil
	}
	return &Person
}

func (s service) DeletePerson(ID int) error {
	query := "DELETE FROM persons WHERE ID = ?"
	statementDelete, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statementDelete.Exec(ID)
	if err != nil {
		return err
	}
	return nil
}

func (s service) UpdatePerson(ID int, p Person) error {
	query := "UPDATE persons SET name = ?, lastname = ?, age = ? WHERE id = :id"
	statementUpdate, err := s.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statementUpdate.Exec(p.Name, p.Lastname, p.Age, ID)
	if err != nil {
		return err
	}
	return nil
}
