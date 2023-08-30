package commands

import "github.com/grip211/order/domain"

type SaveCommand struct {
	Request *domain.Order
}
