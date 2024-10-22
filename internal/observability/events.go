/*
 * SPDX-FileCopyrightText: 2024 Siemens AG
 *
 * SPDX-License-Identifier: MIT
 *
 */

package observability

import (
	"sync"
	"time"
)

const maxItemsToKeep = 20

func timestampNow() string {
	currentTime := time.Now().UTC()
	return currentTime.Format(time.RFC3339Nano)

}

type Events struct {
	sync.RWMutex
	DiscoveryCount     int
	StartedDiscoveries []DiscoveryEvent
}

var globalEvents = newEventState()

func newEventState() *Events {
	return &Events{
		StartedDiscoveries: []DiscoveryEvent{},
	}
}

func GlobalEvents() *Events {
	return globalEvents
}

func ResetGlobalEvents() {
	globalEvents = newEventState()
}

// Discovery Events
type DiscoveryEvent struct {
	Timestamp string `json:"timestamp"`
	JobId     uint32 `json:"job_id"`
}

func (e *Events) GetDiscoveryJobsCount() int {
	e.Lock()
	defer e.Unlock()

	return e.DiscoveryCount

}

func (e *Events) StartedDiscoveryJob(id uint32) {
	e.Lock()
	defer e.Unlock()
	e.StartedDiscoveries = append(e.StartedDiscoveries, DiscoveryEvent{
		Timestamp: timestampNow(),
		JobId:     id,
	})
	e.DiscoveryCount++
	e.StartedDiscoveries = keepMaxSize(e.StartedDiscoveries)

}

func (e *Events) GetDiscoveryJobs() []DiscoveryEvent {
	e.Lock()
	defer e.Unlock()

	return e.StartedDiscoveries
}

// Keep only a certain amount of events as history
func keepMaxSize(items []DiscoveryEvent) []DiscoveryEvent {
	count := len(items)
	if count > maxItemsToKeep {
		items = items[count-maxItemsToKeep:]
	}
	return items
}
