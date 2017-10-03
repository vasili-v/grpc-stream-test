package main

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream && protoc -I $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/ $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream.proto --go_out=plugins=grpc:$GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream && ls $GOPATH/src/github.com/vasili-v/grpc-stream-test/gst-server/stream"

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"

	pb "github.com/vasili-v/grpc-stream-test/gst-server/stream"
)

type server struct{}

func (s *server) Test(stream pb.Stream_TestServer) error {
	fmt.Println("got new stream")

	ch := make(chan pb.Response)
	th := make(chan int, limit)
	go func() {
		defer close(ch)

		var wg sync.WaitGroup
		defer wg.Wait()

		for {
			in, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("stream depleted")
				return
			}

			if err != nil {
				fmt.Printf("receiving error: %s\n", err)
				return
			}

			wg.Add(1)
			th <- 0
			go func(in *pb.Request) {
				defer func() {
					<-th
					wg.Done()
				}()

				time.Sleep(2 * time.Microsecond)

				ch <- pb.Response{
					Id:          in.Id,
					Status:      time.Now().Format(time.RFC3339Nano),
					Obligations: in.Attributes,
				}
			}(in)
		}
	}()

	for out := range ch {
		err := stream.Send(&out)
		if err != nil {
			fmt.Printf("sending error: %s\n", err)
			return err
		}
	}

	return nil
}

func main() {
	ln, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}

	p := grpc.NewServer()
	pb.RegisterStreamServer(p, &server{})
	err = p.Serve(ln)
	if err != nil {
		panic(err)
	}
}
