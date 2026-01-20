// internal/session/manager.go
package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

const (
	sessionCookieName = "threefive_session"
	sessionIDLength   = 32
)

type Store interface {
	Find(ctx context.Context, token string) (b []byte, exists bool, err error)
	Commit(ctx context.Context, token string, b []byte, expiry time.Time) error
	Delete(ctx context.Context, token string) error
	Cleanup(ctx context.Context) error
}

type SessionManager struct {
	store    Store
	lifetime time.Duration
	logger   *slog.Logger
}

type sessionData struct {
	Values map[string]interface{}
}

func init() {
	// Register types for gob encoding
	gob.Register(map[string]interface{}{})
}

func NewManager(store Store, lifetime time.Duration, logger *slog.Logger) *SessionManager {
	return &SessionManager{
		store:    store,
		lifetime: lifetime,
		logger:   logger.With("component", "session_manager"),
	}
}

// Load retrieves the session data for the given request.
func (s *SessionManager) Load(ctx context.Context, c *echo.Context) (map[string]interface{}, error) {
	// Get session token from cookie
	cookie, err := c.Request().Cookie(sessionCookieName)
	if err != nil {
		return make(map[string]interface{}), nil
	}

	token := cookie.Value
	if token == "" {
		return make(map[string]interface{}), nil
	}

	// Load session data from store
	b, exists, err := s.store.Find(ctx, token)
	if err != nil {
		s.logger.Error("failed to load session",
			"error", err,
			"token", token[:10]+"...",
		)
		return make(map[string]interface{}), err
	}

	if !exists {
		return make(map[string]interface{}), nil
	}

	// Decode session data
	data, err := s.decode(b)
	if err != nil {
		s.logger.Error("failed to decode session data",
			"error", err,
		)
		return make(map[string]interface{}), err
	}

	s.logger.Debug("session loaded",
		"token", token[:10]+"...",
	)

	return data.Values, nil
}

// Save persists the session data and writes the session cookie.
func (s *SessionManager) Save(ctx context.Context, c *echo.Context, values map[string]interface{}) error {
	// Get existing token or generate new one
	token, err := s.getToken(c)
	if err != nil {
		token, err = s.generateToken()
		if err != nil {
			s.logger.Error("failed to generate session token",
				"error", err,
			)
			return err
		}
	}

	// Encode session data
	data := &sessionData{Values: values}
	b, err := s.encode(data)
	if err != nil {
		s.logger.Error("failed to encode session data",
			"error", err,
		)
		return err
	}

	// Calculate expiry
	expiry := time.Now().Add(s.lifetime)

	// Save to store
	err = s.store.Commit(ctx, token, b, expiry)
	if err != nil {
		return err
	}

	// Write cookie
	s.writeSessionCookie(c, token, expiry)

	s.logger.Debug("session saved",
		"token", token[:10]+"...",
		"expiry", expiry,
	)

	return nil
}

// Destroy deletes the session data and expires the session cookie.
func (s *SessionManager) Destroy(ctx context.Context, c *echo.Context) error {
	token, err := s.getToken(c)
	if err != nil {
		return nil
	}

	err = s.store.Delete(ctx, token)
	if err != nil {
		return err
	}

	// Clear cookie
	s.clearSessionCookie(c)

	s.logger.Debug("session destroyed",
		"token", token[:10]+"...",
	)

	return nil
}

// RenewToken updates the session token while preserving the session data.
func (s *SessionManager) RenewToken(ctx context.Context, c *echo.Context) error {
	// Load existing session data
	values, err := s.Load(ctx, c)
	if err != nil {
		return err
	}

	// Delete old session
	oldToken, _ := s.getToken(c)
	if oldToken != "" {
		s.store.Delete(ctx, oldToken)
	}

	// Generate new token
	newToken, err := s.generateToken()
	if err != nil {
		return err
	}

	// Encode session data
	data := &sessionData{Values: values}
	b, err := s.encode(data)
	if err != nil {
		return err
	}

	// Save with new token
	expiry := time.Now().Add(s.lifetime)
	err = s.store.Commit(ctx, newToken, b, expiry)
	if err != nil {
		return err
	}

	// Write new cookie
	s.writeSessionCookie(c, newToken, expiry)

	s.logger.Info("session token renewed",
		"old_token", oldToken[:10]+"...",
		"new_token", newToken[:10]+"...",
	)

	return nil
}

// Cleanup is exposed to allow the session store to cleanup expired sessions
func (s *SessionManager) Cleanup(ctx context.Context) error {
	return s.store.Cleanup(ctx)
}

// Helper functions

func (s *SessionManager) getToken(c *echo.Context) (string, error) {
	cookie, err := c.Request().Cookie(sessionCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (s *SessionManager) generateToken() (string, error) {
	b := make([]byte, sessionIDLength)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

func (s *SessionManager) writeSessionCookie(c *echo.Context, token string, expiry time.Time) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiry,
		MaxAge:   int(time.Until(expiry).Seconds()),
		Secure:   false, // Set to true in production with HTTPS
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	c.SetCookie(cookie)
}

func (s *SessionManager) clearSessionCookie(c *echo.Context) {
	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	}
	c.SetCookie(cookie)
}

func (s *SessionManager) encode(data *sessionData) ([]byte, error) {
	var b []byte
	buf := &bufferedWriter{b: &b}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (s *SessionManager) decode(b []byte) (*sessionData, error) {
	dec := gob.NewDecoder(&bufferedReader{b: b})
	data := &sessionData{}
	err := dec.Decode(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type bufferedWriter struct {
	b *[]byte
}

func (w *bufferedWriter) Write(p []byte) (int, error) {
	*w.b = append(*w.b, p...)
	return len(p), nil
}

type bufferedReader struct {
	b []byte
	i int
}

func (r *bufferedReader) Read(p []byte) (int, error) {
	if r.i >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.i:])
	r.i += n
	return n, nil
}
