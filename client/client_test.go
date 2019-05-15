package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	client, err := NewClient(&Config{
		VaultAddress: "http://localhost:8200",
	})
	assert.NoError(t, err)

	client.SetToken("s.STh7XGnk3krBo18GqLepPmVs")

	t.Run("Read", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			result := client.Read("/secret/data/a", "b")
			assert.NoError(t, result.Error)

			assert.Equal(t, result.Value, "c")
		})

		t.Run("WrongPath", func(t *testing.T) {
			result := client.Read("/secret/a", "b")

			assert.Error(t, result.Error)
			assert.Len(t, result.Warnings, 1)
		})
	})
}
