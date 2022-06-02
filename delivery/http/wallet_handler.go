package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/alramdein/e-wallet/event"
	"github.com/alramdein/e-wallet/models"
	"github.com/alramdein/e-wallet/pb"
	"github.com/alramdein/e-wallet/usecase"
	"github.com/lovoo/goka"
	"google.golang.org/protobuf/encoding/protojson"

	echo "github.com/labstack/echo/v4"
)

// var handler *WalletHandler

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type WalletHandler struct {
	WalletUsecase models.WalletUsecase
	Emitter       *goka.Emitter
}

func NewWalletHandler(e *echo.Echo, wc models.WalletUsecase, emitter *goka.Emitter) {
	handler := &WalletHandler{
		WalletUsecase: wc,
		Emitter:       emitter,
	}

	e.GET("/wallet/:id/details", handler.FetchBalance)
	e.POST("/deposit", handler.Deposit)
}

func (w *WalletHandler) FetchBalance(c echo.Context) error {
	id := c.Param("id")
	walletID, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "invalid id")
	}

	wallet, err := w.WalletUsecase.GetWalletByID(c.Request().Context(), int64(walletID))
	switch err {
	case nil:
		break
	case usecase.ErrNotFound:
		return c.JSON(http.StatusNotFound, &Response{
			Message: "wallet not found",
			Data:    nil,
		})
	default:
		fmt.Printf("error when GetWalletByID:%v", err)
		return c.JSON(http.StatusInternalServerError, &Response{
			Message: "internal server error",
			Data:    nil,
		})
	}

	return c.JSON(http.StatusOK, &Response{
		Message: "success fetching data",
		Data:    wallet,
	})
}

func (w *WalletHandler) Deposit(c echo.Context) error {
	walletID, err := strconv.Atoi(c.FormValue("wallet_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Response{
			Message: "invalid walletID",
			Data:    nil,
		})
	}
	amount, err := strconv.Atoi(c.FormValue("amount"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, &Response{
			Message: "invalid amount",
			Data:    nil,
		})
	}

	key := "deposit-" + fmt.Sprint(walletID)

	deposit := &pb.Deposit{
		WalletId: int64(walletID),
		Amount:   int64(amount),
	}
	message, err := protojson.Marshal(deposit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &Response{
			Message: "failed encode deposit data",
			Data:    deposit,
		})
	}

	event.Send(w.Emitter, key, string(message))
	return c.JSON(http.StatusOK, &Response{
		Message: "successfully deposit money",
		Data:    deposit,
	})
}
