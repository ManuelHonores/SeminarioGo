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
		c.JSON(http.StatusOK, gin.H{
			"persons": s.FindAll(),
		})
	}
}

// FindPersonByID ...
func FindPersonByID(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		c.JSON(http.StatusOK, gin.H{
			"person": s.FindByID(id),
		})
	}
}

// InsertPerson ...
func InsertPerson(s Service) gin.HandlerFunc {
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
		err = s.AddPerson(Person)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

// DeletePerson ...
func DeletePerson(s Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = s.DeletePerson(id)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"Message": "Deleted",
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
		err = s.UpdatePerson(id, Person)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"Message": "Modified",
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
