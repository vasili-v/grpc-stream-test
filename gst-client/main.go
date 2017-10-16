package main

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream && protoc -I $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/ $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream.proto --go_out=plugins=grpc:$GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream && ls $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream"

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/vasili-v/grpc-stream-test/gst-server/stream"
)

type pair struct {
	req *pb.Request

	sent time.Time
	recv *time.Time
	dup  int
}

func main() {
	pairs := newPairs(total, size)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	if compression {
		opts = append(opts,
			grpc.WithCompressor(grpc.NewGZIPCompressor()),
			grpc.WithDecompressor(grpc.NewGZIPDecompressor()),
		)
	}
	c, err := grpc.Dial(server, opts...)
	if err != nil {
		panic(fmt.Errorf("dialing error: %s", err))
	}
	defer c.Close()

	if clients > 1 {
		var wg sync.WaitGroup

		step := len(pairs) / clients
		if step <= 0 {
			step = 1
		}

		client := pb.NewStreamClient(c)

		for i := 0; i*step < len(pairs); i++ {
			s, err := client.Test(context.Background())
			if err != nil {
				panic(fmt.Errorf("rpc error: %s", err))
			}

			start := i * step
			end := start + step
			if end > len(pairs) {
				end = len(pairs)
			}

			wg.Add(1)
			go func(s pb.Stream_TestClient, chunk []*pair) {
				defer wg.Done()

				testStreamSync(s, chunk, start)
			}(s, pairs[start:end])
		}

		wg.Wait()
	}

	dump(pairs, "")
}

func testStreamSync(s pb.Stream_TestClient, pairs []*pair, start int) {
	miss := 0

	for i, p := range pairs {
		p.sent = time.Now()
		err := s.Send(p.req)
		if err != nil {
			panic(fmt.Errorf("sending error: %s", err))
		}

		res, err := s.Recv()
		if err != nil {
			panic(fmt.Errorf("receiving error at %d: %s", i+1, err))
		}

		j := int(res.Id) - start
		if j >= 0 && j < len(pairs) {
			p := pairs[j]
			if p.req.Id == res.Id {
				if p.recv == nil {
					t := time.Now()
					p.recv = &t
				} else {
					p.dup++
				}
			} else {
				miss++
			}
		} else {
			miss++
		}
	}

	err := s.CloseSend()
	if err != nil {
		panic(fmt.Errorf("closing error: %s", err))
	}

	if miss > 0 {
		panic(fmt.Errorf("got %d messages with invalid ids", miss))
	}

	dup := 0
	for _, p := range pairs {
		dup += p.dup
	}

	if dup > 0 {
		panic(fmt.Errorf("got %d duplicates", dup))
	}
}

func testStream(c *grpc.ClientConn, pairs []*pair) {
	s, err := pb.NewStreamClient(c).Test(context.Background())
	if err != nil {
		panic(fmt.Errorf("rpc error: %s", err))
	}

	miss := 0
	ch := make(chan int)

	count := len(pairs)
	go func() {
		defer func() { close(ch) }()

		for count > 0 {
			out, err := s.Recv()
			if err != nil {
				panic(fmt.Errorf("receiving error at %d: %s", len(pairs)-count+1, err))
			}

			if out.Id < uint32(len(pairs)) {
				p := pairs[out.Id]
				if p.recv == nil {
					t := time.Now()
					p.recv = &t
					count--
				} else {
					p.dup++
				}
			} else {
				miss++
			}
		}
	}()

	for _, p := range pairs {
		p.sent = time.Now()
		err := s.Send(p.req)
		if err != nil {
			panic(fmt.Errorf("sending error: %s", err))
		}
	}

	if count > 0 {
		fmt.Fprintf(os.Stderr, "waiting for %d responses\n", count)
	}

	select {
	case <-ch:
	case <-time.After(timeout):
	}

	err = s.CloseSend()
	if err != nil {
		panic(fmt.Errorf("closing error: %s", err))
	}

	if count > 0 {
		panic(fmt.Errorf("couldn't receive %d responses", count))
	}

	if miss > 0 {
		panic(fmt.Errorf("got %d messages with invalid ids", miss))
	}

	dup := 0
	for _, p := range pairs {
		dup += p.dup
	}

	if dup > 0 {
		panic(fmt.Errorf("got %d duplicates", dup))
	}
}

func newPairs(n, size int) []*pair {
	out := make([]*pair, n)

	if size > 0 {
		fmt.Fprintf(os.Stderr, "making messages to send:\n")
	}

	for i := range out {
		if size > 0 {
			buf := make([]byte, size)
			if random {
				for j := range buf {
					buf[j] = byte(rand.Intn(256))
				}
			} else {
				for j := range buf {
					buf[j] = 0xaa
				}
			}

			if i < 3 {
				fmt.Fprintf(os.Stderr, "\t%d: % x\n", i, buf)
			} else if i == 3 {
				fmt.Fprintf(os.Stderr, "\t%d: ...\n", i)
			}

			out[i] = &pair{
				req: &pb.Request{
					Id:      uint32(i),
					Payload: buf,
				},
			}
		} else {
			out[i] = &pair{}
		}
	}

	return out
}
