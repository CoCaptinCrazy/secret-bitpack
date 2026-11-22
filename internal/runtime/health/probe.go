package health

import (
	"context"
	"net/http"
	"time"
)

func HTTP(ctx context.Context, url string, timeout time.Duration) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return &StatusError{Code: response.StatusCode}
	}
	return nil
}
