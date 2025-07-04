package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/scmbr/test-task/internal/dto"
	"github.com/sirupsen/logrus"
)

type IPNotifier struct {
	client     *http.Client
	webhookURL string
}

func NewIPNotifier(url string) *IPNotifier {
	return &IPNotifier{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		webhookURL: url,
	}
}
func (n *IPNotifier) NotifyChange(userGUID, oldIP, newIP string) error {
	payload := dto.IPChangeNotification{
		UserGUID: userGUID,
		OldIP:    oldIP,
		NewIP:    newIP,
		Time:     time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal notification payload: %w", err)
	}
	logrus.WithFields(logrus.Fields{
		"userGUID": userGUID,
		"oldIP":    oldIP,
		"newIP":    newIP,
		"time":     payload.Time,
		"payload":  string(jsonData),
	}).Info("Sending IP change notification")
	req, err := http.NewRequest(
		http.MethodPost,
		n.webhookURL,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "IPChangeNotifier/1.0")

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("webhook returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
