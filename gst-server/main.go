package main

//go:generate bash -c "mkdir -p $GOPATH/src/github.com/vasili-v/grpc-stream-test/stream && protoc -I $GOPATH/src/github.com/vasili-v/grpc-stream-test/ $GOPATH/src/github.com/vasili-v/grpc-stream-test/stream.proto --go_out=plugins=grpc:$GOPATH/src/github.com/vasili-v/grpc-stream-test/stream && ls $GOPATH/src/github.com/vasili-v/grpc-stream-test/stream"

import (
	"fmt"
	"io"
	"net"

	"google.golang.org/grpc"

	pb "github.com/vasili-v/grpc-stream-test/stream"
)

func handler(in *pb.Message) *pb.Message {
	return &pb.Message{
		Payload: in.Payload,
	}
}

type server struct{}

func (s *server) New(stream pb.Stream_NewServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Printf("receiving error: %s\n", err)
			return err
		}

		err = stream.Send(handler(in))
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

	opts := []grpc.ServerOption{}
	if maxStreams > 0 {
		opts = append(opts,
			grpc.MaxConcurrentStreams(uint32(maxStreams)),
		)
	}

	p := grpc.NewServer(opts...)
	pb.RegisterStreamServer(p, &server{})
	err = p.Serve(ln)
	if err != nil {
		panic(err)
	}
}
