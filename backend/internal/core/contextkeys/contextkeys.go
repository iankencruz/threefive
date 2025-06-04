// internal/core/context/contextkeys.go
package contextkeys

type ctxKey string

var (
	SessionID ctxKey = "sessionID"
	User      ctxKey = "user"
	CompanyID ctxKey = "companyID"
	RequestID ctxKey = "requestID"
	Flash     ctxKey = "flash"
)
