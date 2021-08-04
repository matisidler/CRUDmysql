package invoiceheader

import (
	"database/sql"
	"time"
)

//Modelo de invoiceheader
type Model struct {
	ID        uint
	Client    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Storage interface {
	Migrate() error
	CreateTx(*sql.Tx, *Model) error
}

//Servicio de producto
type Service struct {
	storage Storage
}

//Retorna un puntero de Service
func NewService(s Storage) *Service {
	return &Service{s}
}

//Migrate se usa para migrar producto. Es decir, crear la tabla producto.
func (s *Service) Migrate() error {
	return s.storage.Migrate()
}

func (s *Service) CreateTx(t *sql.Tx, m *Model) error {
	return s.storage.CreateTx(t, m)
}
