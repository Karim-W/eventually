//
//  model.go
//  eventually
//
//  Created by karim-w on 11/07/2025.
//

package eventually

import "sync"

type BackpressureMode int

const (
	Drop BackpressureMode = iota
	Block
	Async
)

type observable[T any] struct {
	mode          BackpressureMode
	subscribersMu sync.RWMutex
	subscribers   map[chan T]struct{}
}
