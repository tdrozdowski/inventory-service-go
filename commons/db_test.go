package commons

import (
	"github.com/stretchr/testify/assert"
	"os"
	"sync"
	"testing"
)

func TestGetDB_BadUri(t *testing.T) {
	once = sync.Once{}
	err := os.Setenv("DATABASE_URL", "bad uri")
	if err != nil {
		assert.FailNow(t, "Couldn't set DATABASE_URI")
	}

	assert.Panics(t, func() {
		_ = GetDB()
	})

}

func TestGetDB_GetDB(t *testing.T) {
	once = sync.Once{}
	err := os.Setenv("DATABASE_URL", "postgres://localhost:5432/test?sslmode=disable")
	if err != nil {
		assert.FailNow(t, "Couldn't set DATABASE_URI")
	}
	assert.NotPanics(t, func() {
		p := GetDB()
		assert.NotNil(t, p)
	})
}
