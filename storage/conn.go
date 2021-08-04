//ACÁ NOS CONECTAMOS A LA BD

package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/matisidler/CRUDpqv2/pkg/product"
)

//Creamos la conexión a la BD.
//Utilizamos el Patrón Singleton para que solo se ejecute una vez.

//Creamos dos variables que van a poder ser usadas por todos los archivos del paquete storage.
//Con once hacemos el patrón Singleton para que se ejecute una sola vez.
var (
	db   *sql.DB
	once sync.Once
)

type Driver string

const (
	MySQL    Driver = "MYSQL"
	Postgres Driver = "POSTGRES"
)

func NewConnection(d Driver) {
	switch d {
	case MySQL:
		newMySqlDB()
	case Postgres:
		newPqDB()
	default:
		log.Fatal("Parametro no valido")
	}

}

func newPqDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("postgres", "postgres://postgres:password@localhost:5432/gocrud?sslmode=disable")
		if err != nil {
			log.Fatalf("can't open DB %v", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("can't do ping: %v", err)
		}
		fmt.Println("Conectado a MySQL.")
	})
	return db
}

func newMySqlDB() *sql.DB {
	once.Do(func() {
		var err error
		db, err = sql.Open("mysql", "root:password@tcp(localhost:3306)/gocrud?parseTime=true")
		if err != nil {
			log.Fatalf("can't open DB %v", err)
		}

		err = db.Ping()
		if err != nil {
			log.Fatalf("can't do ping: %v", err)
		}
		fmt.Println("Conectado a MySQL.")
	})
	return db
}

func stringToNull(s string) sql.NullString {
	var nullString sql.NullString
	if s == "" {
		nullString.Valid = false
	} else {
		nullString.Valid = true
		nullString.String = s
	}
	return nullString
}

func Pool() *sql.DB {
	return db
}

func DAOProduct(d Driver) (product.Storage, error) {
	switch d {
	case Postgres:
		return newPsqlProduct(db), nil
	case MySQL:
		return newMySqlProduct(db), nil
	default:
		return nil, errors.New("Driver no aceptado")
	}
}
