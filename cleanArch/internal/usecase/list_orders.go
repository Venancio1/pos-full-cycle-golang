package usecase

import (
	"github.com/pos-full-cycle-golang/cleanArch/internal/entity"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	OrderRepository entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: OrderRepository,
	}
}

func (l *ListOrdersUseCase) ListOrders() ([]*entity.Order, error) {
	orders, err := l.OrderRepository.List()
	if err != nil {
		return nil, err
	}

	return orders, nil
}
