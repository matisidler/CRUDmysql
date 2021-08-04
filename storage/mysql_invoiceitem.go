//ACÁ ESTAN LAS QUERYS DE INVOICE HEADER

package storage

import (
	"database/sql"
	"fmt"

	"github.com/matisidler/CRUDpqv2/pkg/invoiceitem"
)

//Creamos una constante (como mi variable "q") para ejecutar las querys.
const (
	//CONSTAINT: por defecto se pone asi: nombreTabla_nombreColumna_primaryKey/foreignKey
	mySqlMigrateInvoiceItem = `CREATE TABLE IF NOT EXISTS invoice_items(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		invoice_header_id INT NOT NULL,
		product_id INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP,
		CONSTRAINT invoice_items_invoiceheaders_id_fk FOREIGN KEY (invoice_header_id) REFERENCES invoice_headers (id) ON UPDATE RESTRICT ON DELETE RESTRICT,		
		CONSTRAINT invoice_items_product_id_fk FOREIGN KEY (product_id) REFERENCES products (id) ON UPDATE RESTRICT ON DELETE RESTRICT)`
	mySQLCreateInvoiceItem = `INSERT INTO invoice_items(invoice_header_id, product_id) VALUES (?, ?)`
)

//mySqlInvoiceItem nos genera la variable db para interactuar con la base de datos. CONSTRAINT invoice_items_invoice_header_id_fk FOREIGN KEY (invoice_header_id) REFERENCES invoice_headers (id),
type mySqlInvoiceItem struct {
	db *sql.DB
}

//NewInvoiceItem retorna un nuevo puntero de InvoiceItem
func NewMySqlInvoiceItem(db *sql.DB) *mySqlInvoiceItem {
	return &mySqlInvoiceItem{db}
}

//Migrate crea la tabla Invoice Header en la base de datos
func (p *mySqlInvoiceItem) Migrate() error {
	stmt, err := p.db.Prepare(mySqlMigrateInvoiceItem)
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

func (p *mySqlInvoiceItem) CreateTx(tx *sql.Tx, headerID uint, models []*invoiceitem.Model) error {

	stmt, err := tx.Prepare(mySQLCreateInvoiceItem)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, item := range models {
		res, err := stmt.Exec(headerID, item.ProductID)
		if err != nil {
			return err
		}
		id, err := res.LastInsertId()
		if err != nil {
			return err
		}
		item.ID = uint(id)
	}
	return nil
}
