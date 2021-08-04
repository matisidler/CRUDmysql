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
	mySQLMigrateProduct = `CREATE TABLE IF NOT EXISTS products(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		name VARCHAR(25) NOT NULL,
		observations VARCHAR(100),
		price INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP	) `
	mySQLCreateProduct = `INSERT INTO products(name, observations, price) VALUES(?,?,?)`
	mySQLGetAll        = `SELECT * FROM products`
	mySQLGetById       = `SELECT * FROM products WHERE id = ?`
	mySQLUpdate        = `UPDATE products SET name = ?, observations = ?, price = ?, updated_at = now() WHERE ID = ?`
	mySQLDelete        = `DELETE FROM products WHERE id = ?`
)

/* var obsNull = sql.NullString{}
var updatedAtNull = sql.NullTime{} */

//mySqlProduct nos genera la variable db para interactuar con la base de datos.
type mySqlProduct struct {
	db *sql.DB
}

//NewmySqlProduct retorna un nuevo puntero de mySqlProduct
func newMySqlProduct(db *sql.DB) *mySqlProduct {
	return &mySqlProduct{db}
}

//Migrate crea la tabla Products en la base de datos
func (p *mySqlProduct) Migrate() error {
	stmt, err := p.db.Prepare(mySQLMigrateProduct)
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

func (p *mySqlProduct) Create(m *product.Model) error {
	if m.Name == "" {
		return errors.New("the name can't be empty")
	}
	stmt, err := p.db.Prepare(mySQLCreateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	m.CreatedAt = time.Now()
	res, err := stmt.Exec(m.Name, stringToNull(m.Observaciones), m.Price)
	if err != nil {
		return err
	}
	r, _ := res.RowsAffected()
	if r != 1 {
		return errors.New("error: more than 1 (or 0) rows affected")
	}
	id, _ := res.LastInsertId()
	m.ID = uint(id)
	fmt.Printf("The product was created succesfully. ID: %d\n", m.ID)
	fmt.Printf("%+v\n", m)
	return nil
}

func (p *mySqlProduct) GetAll() (product.Models, error) {
	stmt, err := p.db.Prepare(mySQLGetAll)
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

func (p *mySqlProduct) GetById(i uint) (*product.Model, error) {
	stmt, err := db.Prepare(mySQLGetById)
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

func (p *mySqlProduct) Update(m *product.Model) error {

	stmt, err := db.Prepare(mySQLUpdate)
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

func (p *mySqlProduct) Delete(id uint) error {
	stmt, err := db.Prepare(mySQLDelete)
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
