// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package timing

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/DataDog/datadog-agent/pkg/trace/teststatsd"
)

func TestTiming(t *testing.T) {
	assert := assert.New(t)
	Stop() // https://github.com/DataDog/datadog-agent/issues/13934
	stats := &teststatsd.Client{}

	t.Run("report", func(t *testing.T) {
		set := newSet(stats)
		set.Since("counter1", time.Now().Add(-2*time.Second))
		set.Since("counter1", time.Now().Add(-3*time.Second))
		set.report()

		counts := stats.GetCountSummaries()
		assert.Equal(1, len(counts))
		assert.Contains(counts, "counter1.count")

		gauges := stats.GetGaugeSummaries()
		assert.Equal(2, len(gauges))
		assert.Contains(gauges, "counter1.avg")
		assert.Contains(gauges, "counter1.max")
	})

	t.Run("autoreport", func(t *testing.T) {
		stats.Reset()
		set := newSet(stats)
		set.Since("counter1", time.Now().Add(-1*time.Second))
		set.autoreport(time.Millisecond)
		if runtime.GOOS == "windows" {
			time.Sleep(5 * time.Second)
		}
		time.Sleep(10 * time.Millisecond)
		Stop()
		assert.Contains(stats.GetCountSummaries(), "counter1.count")
	})

	t.Run("panic", func(t *testing.T) {
		Start(stats)
		Stop()
		Stop()
	})

	t.Run("race", func(t *testing.T) {
		stats.Reset()
		set := newSet(stats)
		var wg sync.WaitGroup
		for i := 0; i < 150; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				set.Since("counter1", time.Now().Add(-time.Second))
				set.Since(fmt.Sprintf("%d", rand.Int()), time.Now().Add(-time.Second))
			}()
		}
		for i := 0; i < 150; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				set.report()
			}()
		}
		wg.Wait()
	})
}
