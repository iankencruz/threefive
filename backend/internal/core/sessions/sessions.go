package sessions

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/core/contextkeys"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Manager struct {
	DB              *pgxpool.Pool
	Logger          *slog.Logger
	cleanupInterval time.Duration
	stopCleanup     chan struct{}
}

func New(db *pgxpool.Pool, logger *slog.Logger) *Manager {
	return NewWithCleanupInterval(db, logger, 5*time.Minute) // default: 5min
}

func NewWithCleanupInterval(db *pgxpool.Pool, logger *slog.Logger, interval time.Duration) *Manager {
	m := &Manager{
		DB:              db,
		Logger:          logger,
		cleanupInterval: interval,
		stopCleanup:     make(chan struct{}),
	}
	if interval > 0 {
		go m.startCleanup()
	}
	return m
}

const (
	sessionCookieName = "user_session"
	sessionLifespan   = 5 * time.Hour // 5 Hours
)

func (m *Manager) SetUserID(w http.ResponseWriter, r *http.Request, userID uuid.UUID) error {
	sessionToken := generateSessionToken()
	expiry := time.Now().Add(sessionLifespan)
	userAgent := r.UserAgent()

	args := pgx.NamedArgs{
		"token":      sessionToken,
		"user_id":    userID,
		"expires_at": expiry,
		"user_agent": userAgent,
	}

	_, err := m.DB.Exec(r.Context(), `
		INSERT INTO sessions (token, user_id, expires_at, user_agent)
		VALUES (@token, @user_id, @expires_at, @user_agent)
		ON CONFLICT (token) DO UPDATE SET 
			expires_at = EXCLUDED.expires_at,
			user_agent = EXCLUDED.user_agent
	`, args)

	if err != nil {
		return fmt.Errorf("failed to insert session: %w", err)
	}

	cookie := &http.Cookie{
		Name:     sessionCookieName,
		Value:    sessionToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // set to true in production behind HTTPS
		SameSite: http.SameSiteLaxMode,
		Expires:  expiry,
	}
	http.SetCookie(w, cookie)
	return nil
}

func (m *Manager) GetUserID(r *http.Request) (uuid.UUID, error) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return uuid.Nil, err
	}

	var userID uuid.UUID
	var expiresAt time.Time

	args := pgx.NamedArgs{
		"token": cookie.Value,
		"now":   time.Now(),
	}

	err = m.DB.QueryRow(r.Context(), `
		SELECT user_id, expires_at FROM sessions
		WHERE token = @token AND expires_at > @now
	`, args).Scan(&userID, &expiresAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, nil
		}
		return uuid.Nil, fmt.Errorf("query session failed: %w", err)
	}

	// Refresh session if nearing expiry
	newExpiry := time.Now().Add(sessionLifespan)
	if newExpiry.Sub(expiresAt) > (10 * time.Minute) {
		m.refreshExpiry(r.Context(), cookie.Value, newExpiry)
	}

	return userID, nil
}

func (m *Manager) refreshExpiry(ctx context.Context, token string, newExpiry time.Time) {
	_, _ = m.DB.Exec(ctx, `
		UPDATE sessions SET expires_at = @expires_at WHERE token = @token
	`, pgx.NamedArgs{
		"token":      token,
		"expires_at": newExpiry,
	})
}

func (m *Manager) Clear(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return nil // No session to clear
	}

	_, _ = m.DB.Exec(r.Context(), `
		DELETE FROM sessions WHERE token = @token
	`, pgx.NamedArgs{"token": cookie.Value})

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})
	return nil
}

func (m *Manager) Put(ctx context.Context, token string, key string, value string) error {
	_, err := m.DB.Exec(ctx, `
		INSERT INTO session_data (token, key, value)
		VALUES (@token, @key, @value)
		ON CONFLICT (token, key) DO UPDATE SET value = EXCLUDED.value
	`, pgx.NamedArgs{
		"token": token,
		"key":   key,
		"value": value,
	})
	return err
}

func (m *Manager) GetString(ctx context.Context, token string, key string) (string, error) {
	var val string
	err := m.DB.QueryRow(ctx, `
		SELECT value FROM session_data
		WHERE token = @token AND key = @key
	`, pgx.NamedArgs{
		"token": token,
		"key":   key,
	}).Scan(&val)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return val, nil
}

func generateSessionToken() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36)
}

func (m *Manager) Exists(r *http.Request, key string) bool {
	cookie, err := r.Cookie("user_session")
	if err != nil {
		return false
	}
	return cookie.Value != ""
}

// RenewToken destroys the current session and creates a new one with the same data.

func (m *Manager) RenewToken(ctx context.Context) error {
	token, ok := ctx.Value(contextkeys.SessionID).(string)
	if !ok || token == "" {
		return errors.New("no session token found in context")
	}

	// Destroy the existing session
	if err := m.Destroy(ctx); err != nil {
		return fmt.Errorf("failed to destroy session: %w", err)
	}

	// Create a new session
	newToken := generateSessionToken()
	newExpiry := time.Now().Add(sessionLifespan)

	// You can add metadata later if needed (IP, user-agent, etc.)
	_, err := m.DB.Exec(ctx, `
		INSERT INTO sessions (token, expires_at)
		VALUES (@token, @expires_at)
	`, pgx.NamedArgs{
		"token":      newToken,
		"expires_at": newExpiry,
	})

	if err != nil {
		return fmt.Errorf("failed to create new session: %w", err)
	}

	return nil
}

func (m *Manager) Create(ctx context.Context) (string, error) {
	sessionID := generateSessionToken()

	expiry := time.Now().Add(sessionLifespan)

	_, err := m.DB.Exec(ctx, `
		INSERT INTO sessions (token, user_id, expires_at, user_agent)
		VALUES (@token, NULL, @expires_at, '')
	`, pgx.NamedArgs{
		"token":      sessionID,
		"expires_at": expiry,
	})
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
func (m *Manager) Destroy(ctx context.Context) error {
	cookie, ok := ctx.Value("user_session").(string)
	if !ok || cookie == "" {
		return errors.New("no session ID in context")
	}

	_, err := m.DB.Exec(ctx, `
		DELETE FROM sessions WHERE token = @token
	`, pgx.NamedArgs{"token": cookie})
	return err
}

func (m *Manager) startCleanup() {
	ticker := time.NewTicker(m.cleanupInterval)
	defer ticker.Stop()

	m.Logger.Info("session cleanup started", "interval", m.cleanupInterval.String())

	for {
		select {
		case <-ticker.C:
			m.Logger.Debug("running session cleanup")
			_, err := m.DB.Exec(context.Background(), `
				DELETE FROM sessions WHERE expires_at < now()
			`)
			if err != nil {
				m.Logger.Error("session cleanup error", "error", err)
			}
		case <-m.stopCleanup:
			m.Logger.Info("session cleanup stopped")
			return
		}
	}
}

func (m *Manager) StopCleanup() {
	if m.cleanupInterval > 0 {
		close(m.stopCleanup)
	}
}
