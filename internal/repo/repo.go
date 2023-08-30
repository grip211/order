package repo

import (
	"context"
	"database/sql"

	builder "github.com/doug-martin/goqu/v9"

	"github.com/grip211/order/domain"
	"github.com/grip211/order/pkg/commands"
	"github.com/grip211/order/pkg/database"
	"github.com/grip211/order/pkg/models"
)

type Repo struct {
	database database.Pool
}

func (r *Repo) FindAll(ctx context.Context) ([]domain.Order, error) {
	query := r.database.Builder().
		Select(
			builder.L("*"),
		).
		From("public.order").
		LeftJoin(
			builder.T("payment").Schema("public"),
			builder.On(builder.Ex{"payment.transaction_id": builder.I("order.order_uid")}),
		)

	var response []models.OrderReadModel
	if err := query.ScanStructsContext(ctx, &response); err != nil {
		return nil, err
	}

	itemsQuery := r.database.Builder().Select(builder.L("*")).From("order_item")

	var items []models.Item
	if err := itemsQuery.ScanStructsContext(ctx, &items); err != nil {
		return nil, err
	}

	mapItems := make(map[string][]models.Item, len(items))
	for _, item := range items {
		if _, ok := mapItems[item.TrackNumber]; !ok {
			mapItems[item.TrackNumber] = make([]models.Item, 0, 10)
		}
		mapItems[item.TrackNumber] = append(mapItems[item.TrackNumber], item)
	}

	orders := make([]domain.Order, 0, len(response))
	for _, order := range response {
		orders = append(orders, order.ToEntity(mapItems[order.Order.TrackNumber]))
	}

	return orders, nil
}

func (r *Repo) Find(ctx context.Context, id string) (domain.Order, error) {
	query := r.database.Builder().
		Select(
			builder.L("*"),
		).
		From("public.order").
		LeftJoin(
			builder.T("payment").Schema("public"),
			builder.On(builder.Ex{"payment.transaction_id": builder.I("order.order_uid")}),
		).
		Where(
			builder.L("order_uid").Eq(id),
		)

	var response models.OrderReadModel
	if _, err := query.ScanStructContext(ctx, &response); err != nil {
		return domain.Order{}, err
	}

	itemsQuery := r.database.Builder().
		Select(builder.L("*")).
		From("order_item").
		Where(builder.C("track_number").Eq(response.Order.TrackNumber))

	var items []models.Item
	if err := itemsQuery.ScanStructsContext(ctx, &items); err != nil {
		return domain.Order{}, err
	}

	return response.ToEntity(items), nil
}

func (r *Repo) Save(ctx context.Context, command commands.SaveCommand) (int64, error) {
	order, items, payment := models.OrderModelsFromEntity(command.Request)

	var affected int64
	err := r.database.Builder().WithTx(func(tx *builder.TxDatabase) error {
		var (
			err    error
			result sql.Result
		)

		result, err = tx.Insert("public.order").Rows(order).Executor().ExecContext(ctx)
		if err != nil {
			return err
		}

		affected, err = result.RowsAffected()

		_, err = tx.Insert("public.order_item").Rows(items).Executor().ExecContext(ctx)
		if err != nil {
			return err
		}

		_, err = tx.Insert("public.payment").Rows(payment).Executor().ExecContext(ctx)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return affected, nil
}

func New(db database.Pool) *Repo {
	return &Repo{database: db}
}

func (r *Repo) Pool() database.Pool {
	return r.database
}
