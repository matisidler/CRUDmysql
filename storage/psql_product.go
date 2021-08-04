//ACÁ ESTAN LAS QUERYS

package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/matisidler/CRUDmysql/pkg/product"
)

//Creamos una constante (como mi variable "q") para ejecutar las querys.
const (
	//CONSTAINT: por defecto se pone asi: nombreTabla_nombreColumna_primaryKey/foreignKey
	MigrateProduct = `CREATE TABLE IF NOT EXISTS products(
		id SERIAL NOT NULL,
		name VARCHAR(25) NOT NULL,
		observations VARCHAR(100),
		price INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT products_id_pk PRIMARY KEY (id)

	) `
	//returning id nos devuelve el ID despues de realizar la inserción.
	psqlCreateProduct  = `INSERT INTO products(name,observations,price,created_at) VALUES($1,$2,$3,$4) RETURNING id`
	psqlGetAllProduct  = `SELECT id, name, observations, price, created_at, updated_at FROM products`
	psqlGetProductById = `SELECT * FROM products WHERE ID = $1`
	psqlUpdateProduct  = `UPDATE products SET name = $1, observations = $2, price = $3, updated_at = now() WHERE id = $4`
	psqlDeleteProduct  = `DELETE FROM products WHERE id = $1`
)

var obsNull = sql.NullString{}
var updatedAtNull = sql.NullTime{}

//PsqlProduct nos genera la variable db para interactuar con la base de datos.
type PsqlProduct struct {
	db *sql.DB
}

//NewPsqlProduct retorna un nuevo puntero de PsqlProduct
func newPsqlProduct(db *sql.DB) *PsqlProduct {
	return &PsqlProduct{db}
}

//Migrate crea la tabla Products en la base de datos
func (p *PsqlProduct) Migrate() error {
	stmt, err := p.db.Prepare(MigrateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	fmt.Println("Migración de Producto ejecutada correctamente")
	return nil
}

func (p *PsqlProduct) Create(m *product.Model) error {
	stmt, err := p.db.Prepare(psqlCreateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	m.CreatedAt = time.Now()
	err = stmt.QueryRow(m.Name, stringToNull(m.Observaciones), m.Price, m.CreatedAt).Scan(&m.ID)
	if err != nil {
		return err
	}

	fmt.Println("Se creó el producto correctamente.")
	fmt.Printf("%+v\n", m)
	return nil
}

func (p *PsqlProduct) GetAll() (product.Models, error) {
	stmt, err := p.db.Prepare(psqlGetAllProduct)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	models := make(product.Models, 0)
	for rows.Next() {
		m := &product.Model{}
		//Controlamos los datos Nulos.
		err := rows.Scan(&m.ID, &m.Name, &obsNull, &m.Price, &m.CreatedAt, &updatedAtNull)
		if err != nil {
			return nil, err
		}
		m.Observaciones = obsNull.String
		m.UpdatedAt = updatedAtNull.Time
		models = append(models, m)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	fmt.Println("ID	/	NOMBRE	/	OBSERVACION		/	PRECIO	/	FECHA_CREACION	/	FECHA_ACTUALIZACIÓN	")
	return models, nil

}

func (p *PsqlProduct) GetById(i uint) (*product.Model, error) {
	stmt, err := db.Prepare(psqlGetProductById)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	m := &product.Model{}
	err = stmt.QueryRow(i).Scan(&m.ID, &m.Name, &obsNull, &m.Price, &m.CreatedAt, &updatedAtNull)
	if err != nil {
		return nil, err
	}
	m.Observaciones = obsNull.String
	m.UpdatedAt = updatedAtNull.Time
	return m, nil
}

func (p *PsqlProduct) Update(m *product.Model) error {
	stmt, err := db.Prepare(psqlUpdateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if m.Observaciones == "" {
		obsNull.Valid = false
	} else {
		obsNull.Valid = true
		obsNull.String = m.Observaciones
	}
	res, err := stmt.Exec(&m.Name, &obsNull, &m.Price, &m.ID)
	if err != nil {
		return err
	}
	if modifiedRows, _ := res.RowsAffected(); modifiedRows != 1 {
		return errors.New("error: more than 1 (or 0) rows modified")
	}
	fmt.Println("Product was updated.")
	return nil
}

func (p *PsqlProduct) Delete(id uint) error {
	stmt, err := db.Prepare(psqlDeleteProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec(id)
	if err != nil {
		return err
	}
	if RowsAffected, _ := res.RowsAffected(); RowsAffected != 1 {
		return errors.New("error: more than 1 (or 0) rows modified")
	}
	fmt.Println("deleted correctly.")
	return nil
}
