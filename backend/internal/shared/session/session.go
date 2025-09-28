// backend/internal/shared/session/session.go
package session

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/netip"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds session configuration
type Config struct {
	Duration    time.Duration
	IdleTimeout time.Duration
	CookieName  string
	Domain      string
	Secure      bool
	HTTPOnly    bool
	SameSite    string
}

// DefaultConfig returns default session configuration
func DefaultConfig() Config {
	return Config{
		Duration:    24 * time.Hour, // 24 hours
		IdleTimeout: 2 * time.Hour,  // 2 hours of inactivity
		CookieName:  "session_token",
		Domain:      "",    // Will be set based on environment
		Secure:      true,  // HTTPS only in production
		HTTPOnly:    true,  // Prevent XSS
		SameSite:    "lax", // Good for SSR
	}
}

// Manager handles session operations
type Manager struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
	config  Config
}

// NewManager creates a new session manager
func NewManager(db *pgxpool.Pool, queries *sqlc.Queries, config Config) *Manager {
	return &Manager{
		db:      db,
		queries: queries,
		config:  config,
	}
}

// CreateSession creates a new session for a user
func (m *Manager) CreateSession(ctx context.Context, userID uuid.UUID, r *http.Request) (sqlc.Sessions, error) {
	// Generate secure session token
	token, err := m.generateSecureToken()
	if err != nil {
		return sqlc.Sessions{}, fmt.Errorf("failed to generate session token: %w", err)
	}

	// Calculate expiration time
	expiresAt := time.Now().Add(m.config.Duration)

	// Extract client info
	ipAddress := m.extractIPAddress(r)
	userAgent := m.extractUserAgent(r)

	// Create session in database
	session, err := m.queries.CreateSession(ctx, sqlc.CreateSessionParams{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
		IpAddress: ipAddress,
		UserAgent: userAgent,
	})
	if err != nil {
		return sqlc.Sessions{}, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// GetSession retrieves a session by token
func (m *Manager) GetSession(ctx context.Context, token string) (sqlc.GetSessionByTokenRow, error) {
	session, err := m.queries.GetSessionByToken(ctx, token)
	if err != nil {
		return sqlc.GetSessionByTokenRow{}, fmt.Errorf("session not found or expired: %w", err)
	}

	// Check if session needs refresh (idle timeout)
	if m.shouldRefreshSession(session) {
		if err := m.RefreshSession(ctx, token); err != nil {
			// Log error but don't fail - session is still valid
			// You might want to add logging here
		}
	}

	return session, nil
}

// RefreshSession extends the session expiration time
func (m *Manager) RefreshSession(ctx context.Context, token string) error {
	newExpiresAt := time.Now().Add(m.config.Duration)

	_, err := m.queries.UpdateSessionExpiry(ctx, sqlc.UpdateSessionExpiryParams{
		Token:     token,
		ExpiresAt: newExpiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to refresh session: %w", err)
	}

	return nil
}

// DeactivateSession marks a session as inactive
func (m *Manager) DeactivateSession(ctx context.Context, token string) error {
	err := m.queries.DeactivateSession(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to deactivate session: %w", err)
	}
	return nil
}

// DeactivateAllUserSessions deactivates all sessions for a user
func (m *Manager) DeactivateAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	err := m.queries.DeactivateAllUserSessions(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user sessions: %w", err)
	}
	return nil
}

// GetUserActiveSessions returns all active sessions for a user
func (m *Manager) GetUserActiveSessions(ctx context.Context, userID uuid.UUID) ([]sqlc.Sessions, error) {
	sessions, err := m.queries.GetActiveSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user sessions: %w", err)
	}
	return sessions, nil
}

// CleanupExpiredSessions removes expired sessions from database
func (m *Manager) CleanupExpiredSessions(ctx context.Context) error {
	err := m.queries.CleanupExpiredSessions(ctx)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired sessions: %w", err)
	}
	return nil
}

// SetSessionCookie sets the session cookie in the response
func (m *Manager) SetSessionCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     m.config.CookieName,
		Value:    token,
		Path:     "/",
		Domain:   m.config.Domain,
		Expires:  time.Now().Add(m.config.Duration),
		MaxAge:   int(m.config.Duration.Seconds()),
		HttpOnly: m.config.HTTPOnly,
		Secure:   m.config.Secure,
		SameSite: m.parseSameSite(m.config.SameSite),
	}
	http.SetCookie(w, cookie)
}

// ClearSessionCookie clears the session cookie
func (m *Manager) ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     m.config.CookieName,
		Value:    "",
		Path:     "/",
		Domain:   m.config.Domain,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: m.config.HTTPOnly,
		Secure:   m.config.Secure,
		SameSite: m.parseSameSite(m.config.SameSite),
	}
	http.SetCookie(w, cookie)
}

// GetSessionTokenFromRequest extracts session token from request
func (m *Manager) GetSessionTokenFromRequest(r *http.Request) string {
	// Try cookie first
	if cookie, err := r.Cookie(m.config.CookieName); err == nil {
		return cookie.Value
	}

	// Try Authorization header as fallback
	if auth := r.Header.Get("Authorization"); auth != "" {
		const prefix = "Bearer "
		if len(auth) > len(prefix) && auth[:len(prefix)] == prefix {
			return auth[len(prefix):]
		}
	}

	return ""
}

// StartCleanupRoutine starts a background routine to cleanup expired sessions
func (m *Manager) StartCleanupRoutine(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := m.CleanupExpiredSessions(ctx); err != nil {
					// Log error but continue
					// You might want to add proper logging here
				}
			}
		}
	}()
}

// Helper methods

func (m *Manager) generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 256-bit token
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (m *Manager) extractIPAddress(r *http.Request) *netip.Addr {
	// Check for forwarded IP first (for reverse proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip, err := netip.ParseAddr(xff); err == nil {
			return &ip
		}
	}

	// Check for real IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		if ip, err := netip.ParseAddr(xri); err == nil {
			return &ip
		}
	}

	// Fall back to remote address
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if host != "" {
		if ip, err := netip.ParseAddr(host); err == nil {
			return &ip
		}
	}

	return nil
}

func (m *Manager) extractUserAgent(r *http.Request) pgtype.Text {
	userAgent := r.UserAgent()
	if userAgent == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: userAgent, Valid: true}
}

func (m *Manager) shouldRefreshSession(session sqlc.GetSessionByTokenRow) bool {
	// Refresh if session will expire within the idle timeout period
	timeUntilExpiry := time.Until(session.ExpiresAt)
	return timeUntilExpiry < m.config.IdleTimeout
}

func (m *Manager) parseSameSite(sameSite string) http.SameSite {
	switch sameSite {
	case "strict":
		return http.SameSiteStrictMode
	case "lax":
		return http.SameSiteLaxMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
