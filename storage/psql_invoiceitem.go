//ACÁ ESTAN LAS QUERYS DE INVOICE HEADER

package storage

import (
	"database/sql"
	"fmt"

	"github.com/matisidler/CRUDmysql/pkg/invoiceitem"
)

//Creamos una constante (como mi variable "q") para ejecutar las querys.
const (
	//CONSTAINT: por defecto se pone asi: nombreTabla_nombreColumna_primaryKey/foreignKey
	MigrateInvoiceItem = `CREATE TABLE IF NOT EXISTS invoice_items(
		id SERIAL NOT NULL,
		invoice_header_id INT NOT NULL,
		product_id INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT invoice_items_id_pk PRIMARY KEY (id),
		CONSTRAINT invoice_items_product_id_fk FOREIGN KEY (product_id) REFERENCES products (id))`
	psqlCreateInvoiceItem = `INSERT INTO invoice_items(invoice_header_id, product_id) VALUES ($1, $2) RETURNING id, created_at`
)

//PsqlInvoiceItem nos genera la variable db para interactuar con la base de datos. CONSTRAINT invoice_items_invoice_header_id_fk FOREIGN KEY (invoice_header_id) REFERENCES invoice_headers (id),
type PsqlInvoiceItem struct {
	db *sql.DB
}

//NewInvoiceItem retorna un nuevo puntero de InvoiceItem
func PsqlNewInvoiceItem(db *sql.DB) *PsqlInvoiceItem {
	return &PsqlInvoiceItem{db}
}

//Migrate crea la tabla Invoice Header en la base de datos
func (p *PsqlInvoiceItem) Migrate() error {
	stmt, err := p.db.Prepare(MigrateInvoiceItem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	fmt.Println("Migración de Invoice Item ejecutada correctamente")
	return nil
}

func (p *PsqlInvoiceItem) CreateTx(tx *sql.Tx, headerID uint, models []*invoiceitem.Model) error {
	stmt, err := p.db.Prepare(psqlCreateInvoiceItem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range models {
		err = stmt.QueryRow(headerID, item.ProductID).Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}
