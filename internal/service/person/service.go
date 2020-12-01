package person

import (
	"SeminarioGo/internal/config"
	"fmt"

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
	AddPerson(Person) (*Person, error)
	FindAll() ([]*Person, error)
	UpdatePerson(int, Person) (int, error)
	DeletePerson(int) (int, error)
	FindByID(int) (*Person, error)
}

type service struct {
	db   *sqlx.DB
	conf *config.Config
}

// New ...
func New(db *sqlx.DB, c *config.Config) (Service, error) {
	return service{db, c}, nil
}

func (s service) FindAll() ([]*Person, error) {
	var list []*Person
	if err := s.db.Select(&list, "SELECT * FROM persons"); err != nil {
		return nil, err
	}
	return list, nil
}

func (s service) AddPerson(p Person) (*Person, error) {
	var pers Person
	pers.Name = p.Name
	pers.Lastname = p.Lastname
	pers.Age = p.Age
	fmt.Println(pers)
	query := "INSERT INTO persons (name, lastname, age) VALUES (?,?,?)"
	statementInsert, err := s.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	_, err = statementInsert.Exec(p.Name, p.Lastname, p.Age)
	if err != nil {
		return nil, err
	}
	return &pers, nil
}

func (s service) FindByID(ID int) (*Person, error) {
	var Person Person
	query := "SELECT * FROM persons WHERE ID = ?"
	if err := s.db.Get(&Person, query, ID); err != nil { //Get es analogo a QueryRow
		return nil, err
	}
	return &Person, nil
}

func (s service) DeletePerson(ID int) (int, error) {
	query := "DELETE FROM persons WHERE ID = ?"
	statementDelete, err := s.db.Prepare(query)
	if err != nil {
		return -1, err
	}
	_, err = statementDelete.Exec(ID)
	fmt.Println("Error en delete: ", err)
	if err != nil {
		return -1, err
	}
	return ID, nil
}

func (s service) UpdatePerson(ID int, p Person) (int, error) {
	query := "UPDATE persons SET name = ?, lastname = ?, age = ? WHERE id = :id"
	statementUpdate, err := s.db.Prepare(query)
	if err != nil {
		return -1, err
	}
	res, err := statementUpdate.Exec(p.Name, p.Lastname, p.Age, ID)
	fmt.Println(res)
	fmt.Println(err)
	if err != nil {
		return -1, err
	}
	return ID, err
}
