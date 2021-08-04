package main

import (
	"fmt"

	"github.com/matisidler/CRUDmysql/storage"
)

func main() {
	storage.NewConnection(storage.Postgres)
	serviceProduct, err := storage.DAOProduct(storage.Postgres)
	if err != nil {
		fmt.Println(err)
	}
	ms, err := serviceProduct.GetAll()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(ms)
	}

}

//Lo que quiere el programa es que en vez de pasarle un storage.PsqlInvoiceItem, le pase un invoiceitem.Storage
