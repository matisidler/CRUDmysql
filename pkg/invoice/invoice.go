package invoice

import (
	"github.com/matisidler/CRUDmysql/pkg/invoiceheader"
	"github.com/matisidler/CRUDmysql/pkg/invoiceitem"
)

type Model struct {
	Header *invoiceheader.Model
	Items  []*invoiceitem.Model
}

type Storage interface {
	Create(*Model) error
}

type Service struct {
	storage Storage
}

func NewService(s Storage) *Service {
	return &Service{s}
}

func (s *Service) Create(m *Model) error {
	return s.storage.Create(m)
}
