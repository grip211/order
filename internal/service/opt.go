package service

import (
	"github.com/grip211/order/pkg/database"
)

type Options struct {
	Database *database.Opt
	NatsURL  string
}
