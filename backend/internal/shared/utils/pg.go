// backend/internal/shared/utils/pg.go
package utils

import (
	"github.com/google/uuid"
	"github.com/iankencruz/threefive/internal/shared/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// ============================================
// UUID Conversions
// ============================================

func UUIDToPg(id *uuid.UUID) pgtype.UUID {
	if id == nil {
		return pgtype.UUID{Valid: false}
	}
	return pgtype.UUID{Bytes: *id, Valid: true}
}

func PgToUUID(pu pgtype.UUID) *uuid.UUID {
	if !pu.Valid {
		return nil
	}
	id := uuid.UUID(pu.Bytes)
	return &id
}

// ============================================
// String/Text Conversions
// ============================================

func StrToPg(s *string) pgtype.Text {
	if s == nil {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func PgToStr(pt pgtype.Text) *string {
	if !pt.Valid {
		return nil
	}
	return &pt.String
}

func PtrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// ============================================
// Integer Conversions
// ============================================

func IntToPg(i *int) pgtype.Int4 {
	if i == nil {
		return pgtype.Int4{Valid: false}
	}
	return pgtype.Int4{Int32: int32(*i), Valid: true}
}

func PgToInt(pi pgtype.Int4) *int {
	if !pi.Valid {
		return nil
	}
	i := int(pi.Int32)
	return &i
}

// ============================================
// Boolean Conversions
// ============================================

func BoolToPg(b *bool, defaultValue bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Bool: defaultValue, Valid: true}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func BoolToPgNull(b *bool) pgtype.Bool {
	if b == nil {
		return pgtype.Bool{Valid: false}
	}
	return pgtype.Bool{Bool: *b, Valid: true}
}

func PgToBool(pb pgtype.Bool) *bool {
	if !pb.Valid {
		return nil
	}
	return &pb.Bool
}

// ============================================
// Custom SQLC Type Conversions
// ============================================

func StatusToPg(status *string) sqlc.NullPageStatus {
	if status == nil {
		return sqlc.NullPageStatus{Valid: false}
	}
	return sqlc.NullPageStatus{
		PageStatus: sqlc.PageStatus(*status),
		Valid:      true,
	}
}

func ProjectStatusToPg(status *string) sqlc.NullProjectStatus {
	if status == nil {
		return sqlc.NullProjectStatus{Valid: false}
	}
	return sqlc.NullProjectStatus{
		ProjectStatus: sqlc.ProjectStatus(*status),
		Valid:         true,
	}
}
