package zkLibs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewStoreUnsupported(t *testing.T) {
	client := "localhost:9999"

	kv, err := NewStore(
		"unsupported",
		[]string{client},
		&Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	assert.Error(t, err)
	assert.Nil(t, kv)
	assert.Equal(t, "Backend storage not supported yet, please choose one of ", err.Error())
}