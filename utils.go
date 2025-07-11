//
//  utils.go
//  eventually
//
//  Created by karim-w on 11/07/2025.
//

package eventually

import (
	"context"
	"sync"
)

func Map[T, R any](src Observable[T], fn func(T) R) Observable[R] {
	out := MakeObservable[R]()
	go func() {
		ctx := context.Background()
		ch, _ := src.Subscribe(ctx, 16)
		for val := range ch {
			out.Emit(fn(val))
		}
	}()
	return out
}

func Filter[T any](src Observable[T], predicate func(T) bool) Observable[T] {
	out := MakeObservable[T]()
	go func() {
		ctx := context.Background()
		ch, _ := src.Subscribe(ctx, 16)
		for val := range ch {
			if predicate(val) {
				out.Emit(val)
			}
		}
	}()
	return out
}

func Reduce[T, R any](src Observable[T], initial R, fn func(R, T) R) Observable[R] {
	out := MakeObservable[R]()
	go func() {
		ctx := context.Background()
		ch, _ := src.Subscribe(ctx, 16)
		accumulator := initial
		for val := range ch {
			accumulator = fn(accumulator, val)
			out.Emit(accumulator)
		}
	}()
	return out
}

func CombineLatest[T any](srcs ...Observable[T]) Observable[[]T] {
	out := MakeObservable[[]T]()
	go func() {
		ctx := context.Background()
		channels := make([]<-chan T, len(srcs))
		for i, src := range srcs {
			ch, _ := src.Subscribe(ctx, 16)
			channels[i] = ch
		}

		values := make([]T, len(srcs))
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for i, ch := range channels {
					select {
					case val, ok := <-ch:
						if !ok {
							return
						}
						values[i] = val
					default:
					}
				}
				out.Emit(values)
			}
		}
	}()
	return out
}

func Merge[T any](srcs ...Observable[T]) Observable[T] {
	out := MakeObservable[T]()
	go func() {
		ctx := context.Background()
		var wg sync.WaitGroup

		for _, src := range srcs {
			wg.Add(1)
			go func(s Observable[T]) {
				defer wg.Done()
				ch, _ := s.Subscribe(ctx, 16)
				for val := range ch {
					out.Emit(val)
				}
			}(src)
		}

		wg.Wait()
	}()
	return out
}

func Concat[T any](srcs ...Observable[T]) Observable[T] {
	out := MakeObservable[T]()
	go func() {
		ctx := context.Background()
		for _, src := range srcs {
			ch, _ := src.Subscribe(ctx, 16)
			for val := range ch {
				out.Emit(val)
			}
		}
	}()
	return out
}

func FromChannel[T any](ch <-chan T) Observable[T] {
	out := MakeObservable[T]()
	go func() {
		for val := range ch {
			out.Emit(val)
		}
	}()
	return out
}

func ToChannel[T any](src Observable[T]) <-chan T {
	ch := make(chan T)
	go func() {
		ctx := context.Background()
		subCh, _ := src.Subscribe(ctx, 16)
		for val := range subCh {
			ch <- val
		}
		close(ch)
	}()
	return ch
}

func ForEach[T any](src Observable[T], fn func(T)) {
	go func() {
		ctx := context.Background()
		ch, _ := src.Subscribe(ctx, 16)
		for val := range ch {
			fn(val)
		}
	}()
}
