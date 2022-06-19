package handler

import (
	"errors"
	"go-account-command/errs"
	"go-account-command/router"
	"go-account-command/services"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type accountEventHandler struct {
	accountEventSrv services.AccountEventService
}

func NewAccountEventHandler(accountEventSrv services.AccountEventService) accountEventHandler {
	return accountEventHandler{accountEventSrv: accountEventSrv}
}

func (a accountEventHandler) AccountEventCreate(c *router.Context) {
	input := services.AccountEventCreateRequest{}

	if err := c.ShouldBind(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
			}
			ValidationErrors(c, out)
			c.Abort()
			return
		}
	}

	err := a.accountEventSrv.NewAccountEvent(input)
	if err != nil {
		handleError(c, errs.NewBadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (a accountEventHandler) AccountEventTransfer(c *router.Context) {
	input := services.AccountEventTransferRequest{}

	if err := c.ShouldBind(&input); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ErrorMsg, len(ve))
			for i, fe := range ve {
				out[i] = ErrorMsg{fe.Field(), getErrorMsg(fe)}
			}
			ValidationErrors(c, out)
			c.Abort()
			return
		}
	}

	err := a.accountEventSrv.AccountEventTransfer(input)
	if err != nil {
		handleError(c, errs.NewBadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, nil)
}

func (a accountEventHandler) ClearAccount(c *router.Context) {

	err := a.accountEventSrv.ClearAccount()

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"message": "success",
	})
}
