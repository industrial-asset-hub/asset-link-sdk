/*
 * SPDX-FileCopyrightText: 2023 Siemens AG
 *
 * SPDX-License-Identifier:
 *
 */

package observability

import (
  "sync"
  "time"
)

const maxItemToKeep = 20

func timestampNow() string {
  currentTime := time.Now().UTC()
  return currentTime.Format(time.RFC3339Nano)

}

type Events struct {
  sync.RWMutex
  StartedDiscoveries []DiscoveryEvent
}

var globalEvents = New()

func New() *Events {
  return &Events{
    StartedDiscoveries: []DiscoveryEvent{},
  }
}

func GlobalEvents() *Events {
  return globalEvents
}

func Reset() {
  globalEvents = New()
}

// Discovery Events
type DiscoveryEvent struct {
  Timestamp string `json:"timestamp"`
  JobId     uint32 `json:"job_id"`
}

func (e *Events) StartedDiscoveryJob(id uint32) {
  e.Lock()
  defer e.Unlock()
  e.StartedDiscoveries = append(e.StartedDiscoveries, DiscoveryEvent{
    Timestamp: timestampNow(),
    JobId:     id,
  })
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
  if count > maxItemToKeep {
    items = items[count-maxItemToKeep:]
  }
  return items
}
