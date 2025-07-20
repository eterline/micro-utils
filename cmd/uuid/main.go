package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

func main() {

	count := int64(0)
	ctx, stop := context.WithTimeout(context.Background(), time.Second)
	defer stop()

	for range 20 {
		go func() {
			for {
				uuid.New()
				atomic.AddInt64(&count, 1)
			}
		}()
	}

	<-ctx.Done()
	fmt.Println(atomic.LoadInt64(&count))
}
