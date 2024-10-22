/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package observability

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvents(t *testing.T) {

	t.Run("stats are capped", func(t *testing.T) {
		// starting from a fresh stats state
		s := newEventState()

		// when adding items below the cap
		for n := 0; n < 3; n++ {
			s.StartedDiscoveryJob(uint32(n))
		}
		// then exactly the items are tracked
		jobs := s.GetDiscoveryJobs()
		jobIds := []uint32{jobs[0].JobId, jobs[1].JobId, jobs[2].JobId}
		assert.ElementsMatch(t, jobIds, []uint32{0, 1, 2})

		// when adding above the cap
		for n := 0; n < maxItemsToKeep*2; n++ {
			s.StartedDiscoveryJob(uint32(n))
		}

		assert.Len(t, s.GetDiscoveryJobs(), maxItemsToKeep, "only the maximum items to keep should be shown")
		assert.Equal(t, 3+maxItemsToKeep*2, s.GetDiscoveryJobsCount())

	})

	t.Run("stats are thread-safe", func(t *testing.T) {
		// starting from a fresh stats state
		ResetGlobalEvents()

		// given some stats added in parallel
		const goroutines = 2
		const count = maxItemsToKeep/2 - 1
		addSomeStats(goroutines, count)
		expectedCount := count * goroutines

		fetchedJobs := GlobalEvents().GetDiscoveryJobs()

		assert.Len(t, fetchedJobs, expectedCount)
	})

	t.Run("all stats fields", func(t *testing.T) {
		// starting from a fresh stats state
		s := newEventState()

		for n := 0; n < 3; n++ {
			s.StartedDiscoveryJob(uint32(n))
		}

		assert.Len(t, s.GetDiscoveryJobs(), 3)

	})
}

func addSomeStats(numGoroutines int, count int) {

	var wg sync.WaitGroup
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func() {
			for n := 0; n < count; n++ {
				s := GlobalEvents()
				s.StartedDiscoveryJob(uint32(n))
			}
			defer wg.Done()
		}()

	}
	wg.Wait()

}
