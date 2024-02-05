package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SherClockHolmes/webpush-go"
)

type JSONSubscription struct {
	Subscriptions []webpush.Subscription `json:"subscriptions"`
}

type JSONSubscriptionManager struct {
	subFilePath string
}

type SubscriptionManager interface {
	Add(sub webpush.Subscription) error
	Set(subs []webpush.Subscription) error
	GetAll() ([]webpush.Subscription, error)
}

func NewJSONSubscriptionManager() SubscriptionManager {
	cwd, _ := os.Getwd()
	subFilePath := cwd + "/assets/subscribers.json"
	return &JSONSubscriptionManager{subFilePath: subFilePath}
}

func (s *JSONSubscriptionManager) Add(sub webpush.Subscription) error {
	subscriptions, err := s.GetAll()
	if err != nil {
		return err
	}
	subscriptions = append(subscriptions, sub)
	return s.Set(subscriptions)

}

func (s *JSONSubscriptionManager) Set(subs []webpush.Subscription) error {
	jsonSub := JSONSubscription{
		Subscriptions: subs,
	}

	jsonString, err := json.Marshal(jsonSub)
	if err != nil {
		return fmt.Errorf("unable to marshal subscription to json, error %w", err)
	}

	return replaceFile(s.subFilePath, string(jsonString))
}

func (s *JSONSubscriptionManager) GetAll() ([]webpush.Subscription, error) {
	data, err := readAllFile(s.subFilePath)
	if err != nil {
		return nil, fmt.Errorf("unable to get subscription data, error %w", err)
	}

	jsonSub := JSONSubscription{}
	err = json.Unmarshal([]byte(data), &jsonSub)
	if err != nil {
		return nil, fmt.Errorf("unexpected subscription data, error %w", err)
	}

	return jsonSub.Subscriptions, nil
}

func readAllFile(path string) (string, error) {
	readFile, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("unable to open file %v, error %w", path, err)
	}

	return string(readFile), nil
}

func replaceFile(path string, data string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to open file %v, error %w", path, err)
	}
	defer f.Close()
	f.Seek(0, 0)

	_, err = fmt.Fprintf(f, "%v", data)
	return err
}
