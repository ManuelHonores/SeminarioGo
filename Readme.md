# Seminario Golang

## Se creó una API REST con CRUD de personas

### Sentencia para ejecutar el programa:

~~~
go run cmd/person/personService.go -config ./config/config.yaml
~~~

### Pasos para probar CRUD (probado en Postman):

> Obtener todos las personas (GET): http://localhost:8080/person

> Agregar una nueva persona (POST): http://localhost:8080/person

```js
Ejemplo para utilizar de Body
{
    "name": "Persona1",
    "lastname": "Apellido1",
    "age": 29
}
```

> Modificar una persona (PUT): http://localhost:8080/person/1

(http://localhost:8080/person/{id})

```js
Ejemplo para utilizar de Body
{
    "name": "UpdatePersona1",
    "lastname": "UpdateApellido1",
    "age": 30
}
```

> Eliminar una persona (DELETE): http://localhost:8080/person/1

(http://localhost:8080/person/{id})


### Tambien se realizo la búsqueda de una persona por id

> Obtener una persona en particular: http://localhost:8080/person/1

(http://localhost:8080/person/{id})