package storage

import (
	"database/sql"
	"fmt"

	"github.com/matisidler/CRUDmysql/pkg/invoice"
	"github.com/matisidler/CRUDmysql/pkg/invoiceheader"
	"github.com/matisidler/CRUDmysql/pkg/invoiceitem"
)

type mySqlInvoice struct {
	db            *sql.DB
	storageHeader invoiceheader.Storage
	storageItem   invoiceitem.Storage
}

func NewMySqlInvoice(db *sql.DB, h invoiceheader.Storage, i invoiceitem.Storage) *mySqlInvoice {
	return &mySqlInvoice{
		db:            db,
		storageHeader: h,
		storageItem:   i,
	}
}

func (p *mySqlInvoice) Create(m *invoice.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	err = p.storageHeader.CreateTx(tx, m.Header)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Header: %v", err)
	}
	fmt.Println("Creado con exito.")
	err = p.storageItem.CreateTx(tx, m.Header.ID, m.Items)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("Item: %v", err)
	}
	fmt.Println("Creado con exito.")

	return tx.Commit()

}
