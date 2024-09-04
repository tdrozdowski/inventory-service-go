package commons

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetDB_BadUri(t *testing.T) {
	err := os.Setenv("DATABASE_URL", "bad uri")
	if err != nil {
		assert.FailNow(t, "Couldn't set DB_URI")
	}

	assert.Panics(t, func() {
		_ = GetDB()
	})

}
