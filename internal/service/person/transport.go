package person

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

// HTTPService ...
type HTTPService interface {
	Register(*gin.Engine)
}

type endpoint struct {
	method   string
	path     string
	function gin.HandlerFunc
}

type httpService struct {
	endpoints []*endpoint
}

// NewHTTPTransport ...
func NewHTTPTransport(s Service) HTTPService {
	endpoints := makeEndpoints(s)
	return httpService{endpoints}
}

func makeEndpoints(s Service) []*endpoint {
	list := []*endpoint{}

	list = append(list,
		&endpoint{
			method:   "GET",
			path:     "/person",
			function: FindPersons(s),
		},
		&endpoint{
			method:   "GET",
			path:     "/person/:id",
			function: FindPersonByID(s),
		},
		&endpoint{
			method:   "POST",
			path:     "/person",
			function: InsertPerson(s),
		},
		&endpoint{
			method:   "PUT",
			path:     "/person/:id",
			function: UpdatePerson(s),
		},
		&endpoint{
			method:   "DELETE",
			path:     "/person/:id",
			function: DeletePerson(s),
		},
	)
	return list
}

// FindPersons ...
func FindPersons(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		persons, err := s.FindAll()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Message": "No se pudo procesar su consulta",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"persons": persons,
			})
		}
	}
}

// FindPersonByID ...
func FindPersonByID(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))
		res, err := s.FindByID(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"person": err,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"person": res,
			})
		}
	}
}

// InsertPerson ...
func InsertPerson(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Person Person
		data, _ := ioutil.ReadAll(c.Request.Body)
		if err := json.Unmarshal(data, &Person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Message": err,
			})
		} else {
			res, err := s.AddPerson(Person)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"Message": "La persona no pudo ser agregada a la tabla",
				})
			} else {
				c.JSON(http.StatusOK, gin.H{
					"Message": "Persona agregada correctamente",
					"person":  res,
				})
			}
		}

	}
}

// DeletePerson ...
func DeletePerson(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Param("id"))

		res, err := s.DeletePerson(id)
		fmt.Println("Error en deleteTrasport: ", err)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Message": "No se pudo eliminar a la persona seleccionada",
				"ID":      res,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"Message": "La persona ha sido eliminada",
				"ID":      res,
			})
		}
	}
}

// UpdatePerson ...
func UpdatePerson(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		var Person Person
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		if err = json.Unmarshal(data, &Person); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res, err := s.UpdatePerson(id, Person)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"Message": "No se pudo actualizar la persona seleccionada",
				"ID":      res,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"Message": "La persona pudo ser modificada",
				"ID":      res,
			})
		}
	}
}

// Register ...
func (s httpService) Register(r *gin.Engine) {
	for _, e := range s.endpoints {
		r.Handle(e.method, e.path, e.function)
	}
}
