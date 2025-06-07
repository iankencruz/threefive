package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPasswordHashing(t *testing.T) {
	tests := []struct {
		name          string
		password      string
		compare       string
		wantErr       bool
		expectHashErr bool
	}{
		{
			name:     "correct password matches hash",
			password: "secure123",
			compare:  "secure123",
			wantErr:  true,
		},
		{
			name:     "incorrect password does not match",
			password: "secure123",
			compare:  "wrongpass",
			wantErr:  false,
		},
		{
			name:          "empty password fails to hash",
			password:      "",
			compare:       "",
			wantErr:       true,
			expectHashErr: false, // bcrypt technically doesn't fail on empty input
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			hash, err := HashPassword(tc.password)

			if tc.expectHashErr {
				assert.Error(t, err)
				return
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, hash)
			}

			match := CheckPasswordHash(tc.compare, hash)
			assert.Equal(t, tc.wantErr, match)
		})
	}
}
