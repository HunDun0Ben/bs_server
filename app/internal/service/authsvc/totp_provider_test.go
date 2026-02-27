package authsvc

import (
	"context"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
)

func TestTOTPProvider_Verify(t *testing.T) {
	provider := NewTOTPProvider()
	secret := "JBSWY3DPEHPK3PXP" // Example secret

	t.Run("Valid Code", func(t *testing.T) {
		code, err := totp.GenerateCode(secret, time.Now())

		valid, err := provider.Verify(context.Background(), secret, code)
		assert.NoError(t, err)
		assert.True(t, valid)
	})

	t.Run("Invalid Code", func(t *testing.T) {
		valid, err := provider.Verify(context.Background(), secret, "000000")
		assert.Error(t, err)
		assert.False(t, valid)
		assert.Equal(t, "invalid TOTP code", err.Error())
	})
}
