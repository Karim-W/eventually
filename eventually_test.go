//
//  eventually_test.go
//  eventually
//
//  Created by karim-w on 11/07/2025.
//

package eventually_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/karim-w/eventually"
)

func TestMakeObservable(t *testing.T) {
	o := eventually.MakeObservable[int]()
	if o == nil {
		t.Fatal("MakeObservable returned nil")
	}

	go func() {
		time.Sleep(1 * time.Second)
		for i := 0; i < 10; i++ {
			o.Emit(i)
		}
	}()

	ch, unsub := o.Subscribe(nil, 5)
	if ch == nil {
		t.Fatal("Subscribe returned nil channel")
	}
	defer unsub()

	for i := 0; i < 10; i++ {
		select {
		case val := <-ch:
			fmt.Printf("Received value: %d\n", val)
			if val != i {
				t.Errorf("Expected %d, got %d", i, val)
			}
		case <-time.After(5 * time.Second):
			t.Error("Timeout waiting for value")
		}
	}
}
