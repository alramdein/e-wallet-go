package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alramdein/e-wallet/event"
	"github.com/alramdein/e-wallet/models"
	"github.com/alramdein/e-wallet/pb"
	"github.com/lovoo/goka"
	"google.golang.org/protobuf/proto"

	echo "github.com/labstack/echo/v4"
)

type WalletHandler struct {
	WalletUsecase models.WalletUsecase
	Emitter       goka.Emitter
}

func NewWalletHandler(e *echo.Echo, wc models.WalletUsecase, emitter *goka.Emitter) {
	handler := &WalletHandler{
		WalletUsecase: wc,
		Emitter:       *emitter,
	}
	e.GET("/wallet/:id/details", handler.FetchBalance)
	e.POST("/deposit", handler.Deposit)
}

func (w *WalletHandler) FetchBalance(c echo.Context) error {
	id := c.Param("id")
	walletID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusNotFound, "wallet not found")
	}
	w.WalletUsecase.GetWalletByID(c.Request().Context(), int64(walletID))
	return c.JSON(http.StatusOK, nil)
}

func (w *WalletHandler) Deposit(c echo.Context) error {
	walletID, err := strconv.Atoi(c.FormValue("wallet_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid walletID")
	}
	amount, err := strconv.Atoi(c.FormValue("amount"))
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid amount")
	}
	key := "deposit-" + fmt.Sprint(walletID)

	deposit := &pb.Deposit{
		WalletId: int64(walletID),
		Amount:   int64(amount),
	}
	message, err := proto.Marshal(deposit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed encode deposit data")
	}

	event.Send(&w.Emitter, key, message)
	return c.JSON(http.StatusOK, nil)
}
