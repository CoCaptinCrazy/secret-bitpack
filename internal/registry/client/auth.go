package client

import "net/http"

func (c *Client) Authorize(request *http.Request, token string) {
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
}
