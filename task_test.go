// Package schedule
package schedule

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultTask_Run(t *testing.T) {
	var mark bool
	d := NewDefaultTask(func(ctx context.Context) {
		mark = true
	})
	d.Run(context.Background())
	assert.True(t, mark)
}
