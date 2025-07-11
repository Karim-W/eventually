//
//  eventually.go
//  eventually
//
//  Created by karim-w on 11/07/2025.
//

package eventually

import (
	"context"
)

type Observable[T any] interface {
	Emit(value T)
	Subscribe(ctx context.Context, buffer int) (<-chan T, func())
}

func MakeObservable[T any]() *observable[T] {
	return MakeObservableWithMode[T](Block)
}

func MakeObservableWithMode[T any](mode BackpressureMode) *observable[T] {
	return &observable[T]{
		mode:        mode,
		subscribers: make(map[chan T]struct{}),
	}
}
