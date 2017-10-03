package main

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream && protoc -I $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/ $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream.proto --go_out=plugins=grpc:$GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream && ls $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream"

import (
	"fmt"
	"os"
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
	c, err := grpc.Dial(server, grpc.WithInsecure())
	if err != nil {
		panic(fmt.Errorf("dialing error: %s", err))
	}
	defer c.Close()

	pairs := make([]*pair, total)
	for i := range pairs {
		pairs[i] = &pair{
			req: &pb.Request{
				Id: uint32(i),
				Attributes: []*pb.Attribute{
					{
						Id:    "test",
						Type:  "string",
						Value: "test",
					},
				},
			},
		}
	}

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

	dump(pairs, "")
}
