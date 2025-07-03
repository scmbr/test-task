package notifier

import (
	"net/http"
	"time"
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
