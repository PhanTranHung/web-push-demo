package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"webpush/config"
	"webpush/util"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

type NotiPayload struct {
	Title   string `json:"title"`
	Body    string `json:"body"`
	Icon    string `json:"icon"`
	Vibrate []int  `json:"vibrate"`
}

type WebPushHandler struct {
	configs             config.Configs
	subscriptionManager util.SubscriptionManager
}

type WebPushHandlerInterface interface {
	GetVapidPubKey(c echo.Context) error
	SendNotification(c echo.Context) error
	SaveSubscription(c echo.Context) error
}

func NewWebPushHandler(configs config.Configs, subscriptionManager util.SubscriptionManager) WebPushHandlerInterface {
	return &WebPushHandler{configs: configs, subscriptionManager: subscriptionManager}
}

func (h *WebPushHandler) GetVapidPubKey(c echo.Context) error {
	return c.String(http.StatusOK, h.configs.GetENVConfigs().VapidPublicKey)
}

func (h *WebPushHandler) SendNotification(c echo.Context) error {

	req := NotiPayload{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad body").SetInternal(err)
	}
	subs, err := h.subscriptionManager.GetAll()
	if err != nil {
		return fmt.Errorf("unable to get all subscription, error %w", err)
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("unable to marshal payload, error %w", err)
	}

	envcfg := h.configs.GetENVConfigs()
	for _, sub := range subs {
		tr := http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		}
		client := http.Client{Transport: &tr}

		resp, err := webpush.SendNotification(payload, &sub, &webpush.Options{
			Subscriber:      envcfg.VapidContact,
			VAPIDPublicKey:  envcfg.VapidPublicKey,
			VAPIDPrivateKey: envcfg.VapidPrivateKey,
			TTL:             2419200,
			Topic:           "Universal events",
			Urgency:         webpush.UrgencyHigh,
			HTTPClient:      &client,
		})

		if err != nil {
			slog.Error("unable to send notification to subscriber, error %w", err)
		}

		resBody := bytes.Buffer{}
		_, err = io.Copy(&resBody, resp.Body)
		fmt.Println(resBody.String())

		defer resp.Body.Close()
	}

	return c.String(http.StatusOK, "OK, Notification sent")
}

func (h *WebPushHandler) SaveSubscription(c echo.Context) error {

	req := webpush.Subscription{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Bad body").SetInternal(err)
	}
	err := h.subscriptionManager.Add(req)
	if err != nil {
		return err
	}

	return c.String(http.StatusOK, "OK, Subscription saved")
}
