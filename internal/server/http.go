package server

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/grip211/order/domain"
	"github.com/grip211/order/internal/repo"
	"github.com/grip211/order/pkg/commands"
	"github.com/grip211/order/pkg/nats"
)

type cache interface {
	Set(ctx context.Context, key string, value domain.Order, ttl time.Duration) error
	Get(ctx context.Context, key string) (domain.Order, bool, error)
	Has(_ context.Context, key string) bool
}

type saver interface {
	Save(ctx context.Context, command commands.SaveCommand) (int64, error)
	Find(ctx context.Context, id string) (domain.Order, error)
	FindAll(ctx context.Context) ([]domain.Order, error)
}

type HTTPHandler struct {
	saver saver
	nats  *nats.Client
	cache cache
}

func (h *HTTPHandler) All(ctx *fiber.Ctx) error {
	all, err := h.saver.FindAll(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.Render("index", fiber.Map{
		"orders": all,
	})
}

func (h *HTTPHandler) AllREST(ctx *fiber.Ctx) error {
	all, err := h.saver.FindAll(ctx.Context())
	if err != nil {
		return err
	}
	return ctx.JSON(all)
}

func (h *HTTPHandler) Save(ctx *fiber.Ctx) error {
	request := &domain.Order{}
	if err := ctx.BodyParser(request); err != nil {
		return err
	}

	affected, err := h.saver.Save(ctx.Context(), commands.SaveCommand{Request: request})
	if err != nil {
		return err
	}

	return ctx.JSON(map[string]interface{}{
		"rows_affected": affected,
	})
}

func (h *HTTPHandler) Get(ctx *fiber.Ctx) error {
	key := ctx.Query("id")
	if key == "" {
		return ctx.JSON(map[string]string{
			"error": "empty key",
		})
	}
	if h.cache.Has(ctx.Context(), key) {
		order, _, _ := h.cache.Get(ctx.Context(), key)
		return ctx.JSON(map[string]interface{}{
			"order": order,
		})
	}

	order, err := h.saver.Find(ctx.Context(), key)
	if err != nil {
		return err
	}

	_ = h.cache.Set(ctx.Context(), key, order, time.Hour)
	return ctx.JSON(map[string]interface{}{
		"order": order,
	})
}

func (h *HTTPHandler) Publish(ctx *fiber.Ctx) error {
	request := &domain.Order{}
	if err := ctx.BodyParser(request); err != nil {
		return err
	}

	bytes, err := json.Marshal(request)
	if err != nil {
		return err
	}

	if err = h.nats.Publish("test", bytes); err != nil {
		return err
	}

	return ctx.JSON(map[string]interface{}{
		"error": "",
	})
}

func NewHTTPHandler(
	repo *repo.Repo,
	nats *nats.Client,
	cache cache,
) *HTTPHandler {
	h := &HTTPHandler{
		saver: repo,
		nats:  nats,
		cache: cache,
	}
	return h
}
