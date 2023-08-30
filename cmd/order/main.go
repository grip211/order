package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/nats-io/nats.go"
	"github.com/urfave/cli/v2"

	"github.com/grip211/order/config"
	"github.com/grip211/order/domain"
	"github.com/grip211/order/internal/repo"
	"github.com/grip211/order/internal/server"
	"github.com/grip211/order/internal/service"
	"github.com/grip211/order/pkg/cache"
	"github.com/grip211/order/pkg/commands"
	"github.com/grip211/order/pkg/fiberext"
	"github.com/grip211/order/pkg/log"
	"github.com/grip211/order/pkg/signal"
)

func main() {
	application := cli.App{
		Name: "Order",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config-file",
				Required: true,
				Usage:    "YAML config filepath",
				EnvVars:  []string{"CONFIG_FILE"},
				FilePath: "/srv/web_secret/config_file",
			},
			&cli.StringFlag{
				Name:     "bind-address",
				Usage:    "IP и порт сервера, например: 0.0.0.0:3001",
				Required: false,
				Value:    "0.0.0.0:3001",
				EnvVars:  []string{"BIND_ADDRESS"},
			},
			&cli.StringFlag{
				Name:     "bind-socket",
				Usage:    "Путь к Unix сокет файлу",
				Required: false,
				Value:    "/tmp/order.sock",
				EnvVars:  []string{"BIND_SOCKET"},
			},
			&cli.IntFlag{
				Name:     "listener",
				Usage:    "Unix socket or TCP",
				Required: false,
				Value:    1,
				EnvVars:  []string{"LISTENER"},
			},
		},
		Action: Main,
		After: func(c *cli.Context) error {
			log.Info("stopped")
			return nil
		},
	}

	if err := application.Run(os.Args); err != nil {
		log.Error(err)
	}
}

// nolint:funlen // it's ok
func Main(ctx *cli.Context) error {
	appContext, cancel := context.WithCancel(ctx.Context)
	defer func() {
		cancel()
		log.Info("app context is canceled, Luntik is down!")
	}()

	cfg, err := config.New(ctx.String("config-file"))
	if err != nil {
		return err
	}

	orders, err := service.New(appContext, &service.Options{
		Database: &cfg.Database,
		NatsURL:  cfg.NatsURL,
	})
	if err != nil {
		return err
	}

	defer func() {
		orders.Shutdown(func(err error) {
			log.Warning(err)
		})
		orders.Stacktrace()
	}()

	await, stop := signal.Notifier(func() {
		log.Info("received a system signal to shutdown Luntik, start shutdown process..")
	})

	repository := repo.New(orders.Pool)

	handler := server.NewHTTPHandler(repository, orders.Nats, cache.NewInMemory())
	engine := html.New("./templates", ".html")
	app := fiber.New(fiber.Config{
		Views:        engine,
		ServerHeader: "Orders Server",
		Prefork:      cfg.Prefork,
		ErrorHandler: fiberext.ErrorHandler,
	})
	v1 := app.Group("/api/v1")
	app.Get("/all", handler.All)
	v1.Get("/all", handler.AllREST)
	v1.Get("/get", handler.Get)
	v1.Post("/save", handler.Save)
	v1.Post("/publish", handler.Publish)

	// данную логику нужно куда-то перенести в другое место, но пока пусть будет тут
	go func() {
		err := orders.Nats.Subscribe("test", func(msg *nats.Msg) {
			request := &domain.Order{}
			if err := json.Unmarshal(msg.Data, request); err != nil {
				log.Warning(err)
			} else {
				affected, err := repository.Save(orders.Context(), commands.SaveCommand{Request: request})
				if err != nil {
					log.Warning(err)
				}
				log.Info(fmt.Sprintf("Successfully save records: %d", affected))
			}
		})
		if err != nil {
			stop(err)
		}
	}()

	go func() {
		var ln net.Listener
		if ln, err = signal.Listener(
			orders.Context(),
			ctx.Int("listener"),
			ctx.String("bind-socket"),
			ctx.String("bind-address"),
		); err != nil {
			stop(err)
			return
		}
		if err = app.Listener(ln); err != nil {
			stop(err)
		}
	}()

	log.Info("Orders: service is launched")
	return await()
}
