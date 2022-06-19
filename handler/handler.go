package handler

import (
	"fmt"
	"go-account-command/errs"
	"go-account-command/router"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "oneof":
		joinString := strings.Join(strings.Fields(fe.Param()), ", ")
		return fmt.Sprintf("This field condition is [ %v ]", joinString)
	}
	return "Unknown error"
}

func handleError(c *router.Context, err error) {
	switch e := err.(type) {
	case errs.AppError:
		c.JSON(e.Code, map[string]interface{}{
			"error": err.Error(),
		})
	case error:
		c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func ValidationErrors(c *router.Context, err []ErrorMsg) {
	c.JSON(http.StatusUnprocessableEntity, err)
}
