package main

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/vasili-v/grpc-stream-test/stream && protoc -I $GOPATH/src/github.com/vasili-v/grpc-stream-test/ $GOPATH/src/github.com/vasili-v/grpc-stream-test/stream.proto --go_out=plugins=grpc:$GOPATH/src/github.com/vasili-v/grpc-stream-test/stream && ls $GOPATH/src/github.com/vasili-v/grpc-stream-test/stream"

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/vasili-v/grpc-stream-test/stream"
)

type pair struct {
	m *pb.Message

	sent time.Time
	recv *time.Time
}

func main() {
	pairs := newPairs(total, size)

	opts := []grpc.DialOption{grpc.WithInsecure()}
	c, err := grpc.Dial(server, opts...)
	if err != nil {
		panic(fmt.Errorf("dialing error: %s", err))
	}
	defer c.Close()

	client := pb.NewStreamClient(c)
	ss := make([]pb.Stream_NewClient, streams)
	for i := range ss {
		s, err := client.New(context.Background())
		if err != nil {
			panic(fmt.Errorf("rpc error: %s", err))
		}

		ss[i] = s
	}

	var wg sync.WaitGroup

	step := len(pairs) / streams
	rem := len(pairs) % streams

	start := 0

	for i, s := range ss {
		end := start + step
		if i < rem {
			end++
		}

		wg.Add(1)
		go func(s pb.Stream_NewClient, pairs []*pair) {
			defer wg.Done()

			for _, p := range pairs {
				p.sent = time.Now()
				err := s.Send(p.m)
				if err != nil {
					continue
				}

				_, err = s.Recv()
				if err != nil {
					continue
				}

				now := time.Now()
				p.recv = &now
			}

			err := s.CloseSend()
			if err == nil {
				s.Recv()
			}
		}(s, pairs[start:end])

		start = end
		if start >= len(pairs) {
			break
		}
	}

	wg.Wait()

	dump(pairs, "")
}

func newPairs(n, size int) []*pair {
	out := make([]*pair, n)

	fmt.Fprintf(os.Stderr, "making messages to send:\n")

	for i := range out {
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

		if i < 2 || i >= len(out)-1 {
			fmt.Fprintf(os.Stderr, "\t%d: % x\n", i, buf)
		} else if i == 2 {
			fmt.Fprintf(os.Stderr, "\t%d: ...\n", i)
		}

		out[i] = &pair{m: &pb.Message{Payload: buf}}
	}

	return out
}
