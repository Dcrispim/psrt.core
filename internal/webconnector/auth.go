package webconnector

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const pairCodeTTL = 5 * time.Minute

type Auth struct {
	mu        sync.Mutex
	pairCode  string
	pairUntil time.Time
	tokens    map[string]time.Time
	audit     *Audit
}

func NewAuth(audit *Audit) *Auth {
	a := &Auth{
		tokens: make(map[string]time.Time),
		audit:  audit,
	}
	a.regeneratePairCode()
	return a
}

func (a *Auth) regeneratePairCode() {
	b := make([]byte, 3)
	_, _ = rand.Read(b)
	n := (int(b[0])<<16 | int(b[1])<<8 | int(b[2])) % 1000000
	a.pairCode = fmt.Sprintf("%06d", n)
	a.pairUntil = time.Now().Add(pairCodeTTL)
}

func (a *Auth) PairCodeForDisplay() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	if time.Now().After(a.pairUntil) || a.pairCode == "" {
		a.regeneratePairCode()
	}
	return a.pairCode
}

func (a *Auth) RefreshPairCode() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.regeneratePairCode()
	code := a.pairCode
	a.audit.PairCodeRefreshed()
	return code
}

func (a *Auth) PairUntil() time.Time {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.pairUntil
}

func (a *Auth) Pair(code, origin, remote string) (string, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if time.Now().After(a.pairUntil) || code != a.pairCode {
		a.audit.PairFailure(origin, remote)
		return "", false
	}
	a.pairCode = ""
	a.pairUntil = time.Time{}
	token, err := newToken()
	if err != nil {
		return "", false
	}
	a.tokens[token] = time.Now()
	a.audit.PairSuccess(origin, remote)
	return token, true
}

func (a *Auth) ValidateToken(token string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	_, ok := a.tokens[token]
	return ok
}

func newToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func bearerToken(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if !strings.HasPrefix(h, "Bearer ") {
		return ""
	}
	return strings.TrimSpace(h[7:])
}

func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := bearerToken(r)
		if token == "" || !s.auth.ValidateToken(token) {
			s.audit.AuthFailure(r.URL.Path, r.RemoteAddr)
			WriteErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized"))
			return
		}
		next(w, r)
	}
}
