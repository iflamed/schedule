// Package schedule
package schedule

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestScheduler_Timezone(t *testing.T) {
	s := NewScheduler(context.Background(), time.UTC)
	prc, err := time.LoadLocation("Asia/Shanghai")
	assert.NoError(t, err)
	hour := s.now.Hour()
	s.Timezone(prc)
	assert.NotEqual(t, hour, s.now.Hour())
}
