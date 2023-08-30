package fiberext

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/grip211/order/pkg/xerror"
)

const defaultErrorMessage = "Internal Server Error"

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	if err != nil {
		code := fiber.StatusInternalServerError
		value := defaultErrorMessage

		var e *xerror.HTTPError
		var w *xerror.WrapError
		var f *fiber.Error

		if errors.As(err, &f) {
			code = f.Code
			value = f.Message
		} else if errors.As(err, &e) {
			code = e.Code()
			value = e.Message()
		} else {
			value = err.Error()
		}

		if errors.As(err, &w) {
			var pub *xerror.PublicError
			if errors.As(err, &pub) {
				value = fmt.Sprintf("%s: %v", w.Context, pub.Error())
			}
		}

		return ctx.Status(code).JSON(fiber.Map{
			"status":  code,
			"message": value,
		})
	}

	return nil
}
