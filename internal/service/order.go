package service

import (
	"context"

	"github.com/grip211/order/pkg/database"
	"github.com/grip211/order/pkg/database/postgres"
	"github.com/grip211/order/pkg/drop"
	"github.com/grip211/order/pkg/nats"
)

type Service struct {
	*drop.Impl
	Pool database.Pool
	Nats *nats.Client
}

func New(ctx context.Context, opt *Options) (*Service, error) {
	var err error
	s := &Service{}
	s.Impl = drop.NewContext(ctx)

	s.Pool, err = postgres.NewPool(s.Context(), opt.Database)
	if err != nil {
		return nil, err
	}
	s.AddDropper(s.Pool.(*postgres.Pool))

	s.Nats, err = nats.New(opt.NatsURL)
	if err != nil {
		return nil, err
	}

	s.AddDropper(s.Nats)

	return s, nil
}
