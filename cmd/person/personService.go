package main

import (
	"SeminarioGo/internal/config"
	"SeminarioGo/internal/database"
	"SeminarioGo/internal/service/person"
	"flag"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func main() {
	configFile := flag.String("config", "./config.yaml", "this is the service config")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFile)
	fmt.Println("Este es el cfd: ", cfg)
	db, err := database.NewDatabase(cfg)
	defer db.Close()

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if err := createSchema(db); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	service, _ := person.New(db, cfg)
	httpService := person.NewHTTPTransport(service)

	r := gin.Default()
	httpService.Register(r)
	r.Run()
}

func createSchema(db *sqlx.DB) error {
	fmt.Println("llego al schema")
	schema := `CREATE TABLE IF NOT EXISTS persons (
		id integer primary key autoincrement,
		name varchar(50) NOT NULL,
		lastname varchar(50) NOT NULL,
		age integer NOT NULL);`
	fmt.Println(schema)
	// execute a query on the server
	_, err := db.Exec(schema)
	if err != nil {
		fmt.Println("Error creando schema")
		return err
	}
	fmt.Println("Llego")
	return nil
}
