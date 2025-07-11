//
//  observables.go
//  eventually
//
//  Created by karim-w on 11/07/2025.
//

package eventually

import (
	"context"
	"log"
)

func (o *observable[T]) Emit(value T) {
	o.subscribersMu.RLock()
	defer o.subscribersMu.RUnlock()

	for ch := range o.subscribers {
		ch := ch // capture
		switch o.mode {
		case Drop:
			select {
			case ch <- value:
			default:
				log.Println("Dropped due to full channel")
			}
		case Block:
			ch <- value
		case Async:
			go func() {
				ch <- value
			}()
		}
	}
}

func (o *observable[T]) Subscribe(ctx context.Context, buffer int) (<-chan T, func()) {
	// for the assholes who don't pass a context
	if ctx == nil {
		ctx = context.TODO()
	}

	ch := make(chan T, buffer)

	o.subscribersMu.Lock()
	o.subscribers[ch] = struct{}{}
	o.subscribersMu.Unlock()

	unsub := func() {
		o.subscribersMu.Lock()
		defer o.subscribersMu.Unlock()
		if _, ok := o.subscribers[ch]; ok {
			delete(o.subscribers, ch)
			close(ch)
		}
	}

	go func() {
		<-ctx.Done()
		unsub()
	}()

	return ch, unsub
}
