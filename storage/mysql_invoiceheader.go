//ACÁ ESTAN LAS QUERYS DE INVOICE HEADER

package storage

import (
	"database/sql"
	"fmt"

	"github.com/matisidler/CRUDpqv2/pkg/invoiceheader"
)

//Creamos una constante (como mi variable "q") para ejecutar las querys.
const (
	//CONSTAINT: por defecto se pone asi: nombreTabla_nombreColumna_primaryKey/foreignKey
	MySqlMigrateInvoiceHeader = `CREATE TABLE IF NOT EXISTS invoice_headers(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		client VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP
	) `

	mySQLCreateInvoiceHeader = `INSERT INTO invoice_headers(client) VALUES (?)`
)

//mySqlInvoiceHeader nos genera la variable db para interactuar con la base de datos.
type mySqlInvoiceHeader struct {
	db *sql.DB
}

//NewInvoiceHeader retorna un nuevo puntero de InvoiceHeader
func MysqlNewInvoiceHeader(db *sql.DB) *mySqlInvoiceHeader {
	return &mySqlInvoiceHeader{db}
}

//Migrate crea la tabla Invoice Header en la base de datos
func (p *mySqlInvoiceHeader) Migrate() error {
	stmt, err := p.db.Prepare(MySqlMigrateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	fmt.Println("Migración de Invoice Header ejecutada correctamente")
	return nil
}

func (p *mySqlInvoiceHeader) CreateTx(tx *sql.Tx, m *invoiceheader.Model) error {
	stmt, err := tx.Prepare(mySQLCreateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(m.Client)
	if err != nil {
		return err
	}
	id, _ := res.LastInsertId()
	m.ID = uint(id)
	return nil
}
