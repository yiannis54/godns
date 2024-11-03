package dns

import (
	"errors"
	"io"
	"net"
	"net/http"
	"strings"
)

func getIP() (string, error) {
	req, err := http.NewRequest(http.MethodGet, "http://ifconfig.co/", http.NoBody)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res, _ := strings.CutSuffix(string(body), "\n")
	if ip := net.ParseIP(string(res)); ip != nil {
		return ip.String(), nil
	}

	return "", errors.New("invalid ip format")
}
