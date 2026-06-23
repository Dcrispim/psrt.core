package webconnector

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeOriginURL(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("empty origin")
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("allowed_origin must be an absolute URL, got %q", raw)
	}
	if u.Fragment != "" || u.RawQuery != "" {
		return "", fmt.Errorf("allowed_origin must not contain query or fragment")
	}
	host := strings.ToLower(u.Host)
	return u.Scheme + "://" + host, nil
}

func normalizeAllowedOriginsList(raw string) ([]string, error) {
	parts := strings.Split(raw, ",")
	seen := make(map[string]bool)
	var out []string
	for _, part := range parts {
		n, err := normalizeOriginURL(part)
		if err != nil {
			if strings.TrimSpace(part) == "" {
				continue
			}
			return nil, err
		}
		if seen[n] {
			continue
		}
		seen[n] = true
		out = append(out, n)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("allowed_origin must list at least one URL")
	}
	return out, nil
}

func joinAllowedOrigins(origins []string) string {
	return strings.Join(origins, ",")
}

func parseAllowedOriginsField(raw string) (string, error) {
	list, err := normalizeAllowedOriginsList(raw)
	if err != nil {
		return "", err
	}
	return joinAllowedOrigins(list), nil
}

func originsMatch(requestOrigin, allowed string) bool {
	req, err := normalizeOriginURL(requestOrigin)
	if err != nil {
		return false
	}
	for _, a := range strings.Split(allowed, ",") {
		if req == strings.TrimSpace(a) {
			return true
		}
	}
	return false
}
