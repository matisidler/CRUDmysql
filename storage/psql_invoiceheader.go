//ACÁ ESTAN LAS QUERYS DE INVOICE HEADER

package storage

import (
	"database/sql"
	"fmt"

	"github.com/matisidler/CRUDmysql/pkg/invoiceheader"
)

//Creamos una constante (como mi variable "q") para ejecutar las querys.
const (
	//CONSTAINT: por defecto se pone asi: nombreTabla_nombreColumna_primaryKey/foreignKey
	MigrateInvoiceHeader = `CREATE TABLE IF NOT EXISTS invoice_headers(
		id SERIAL NOT NULL,
		client VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT invoice_headers_id_pk PRIMARY KEY (id)
	) `
	psqlCreateInvoiceHeader = `INSERT INTO invoice_headers(client) VALUES ($1) RETURNING id, created_at`
)

//PsqlInvoiceHeader nos genera la variable db para interactuar con la base de datos.
type PsqlInvoiceHeader struct {
	db *sql.DB
}

//NewInvoiceHeader retorna un nuevo puntero de InvoiceHeader
func PsqlNewInvoiceHeader(db *sql.DB) *PsqlInvoiceHeader {
	return &PsqlInvoiceHeader{db}
}

//Migrate crea la tabla Invoice Header en la base de datos
func (p *PsqlInvoiceHeader) Migrate() error {
	stmt, err := p.db.Prepare(MigrateInvoiceHeader)
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

func (p *PsqlInvoiceHeader) CreateTx(tx *sql.Tx, m *invoiceheader.Model) error {
	stmt, err := tx.Prepare(psqlCreateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()
	return stmt.QueryRow(m.Client).Scan(&m.ID, &m.CreatedAt)
}
