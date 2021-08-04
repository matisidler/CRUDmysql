package invoiceitem

import (
	"database/sql"
	"time"
)

//Modelo de invoice item.
type Model struct {
	ID              uint
	InvoiceHeaderID uint
	ProductID       uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Storage interface {
	Migrate() error
	CreateTx(*sql.Tx, uint, []*Model) error
	/* GetAll() ([]*Model, error) */
}

//Servicio de InvoiceItem
type Service struct {
	storage Storage
}

//Retorna un puntero de Service
func NewService(s Storage) *Service {
	return &Service{s}
}

/* func (m *Model) String() string {
	return fmt.Sprintf("%02d | %02d | %02d | %10s | %10s\n",
		m.ID, m.InvoiceHeaderID, m.ProductID, m.CreatedAt.Format("2006-01-02"), m.UpdatedAt.Format("2006-01-02"))
} */

//Migrate se usa para migrar producto. Es decir, crear la tabla producto.
func (s *Service) Migrate() error {
	return s.storage.Migrate()
}

func (s *Service) CreateTx(t *sql.Tx, i uint, m []*Model) error {
	return s.storage.CreateTx(t, i, m)
}

/* func (s *Service) GetAll() ([]*Model, error) {
	return s.storage.GetAll()
}
*/
