package main

import (
	"fmt"

	"github.com/matisidler/CRUDpqv2/storage"
)

func main() {
	storage.NewConnection(storage.MySQL)
	serviceProduct, err := storage.DAOProduct(storage.MySQL)
	if err != nil {
		fmt.Println(err)
	}
	serviceProduct.GetAll()

}

//Lo que quiere el programa es que en vez de pasarle un storage.PsqlInvoiceItem, le pase un invoiceitem.Storage
