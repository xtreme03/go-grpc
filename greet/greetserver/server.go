package main

import (
	"strconv"
	"fmt"
	"context"
	"net"
	"log"
	"google.golang.org/grpc"
	"go-grpc/greet/greetpb"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
	"time"
	"io"
)
type server struct{}
func (*server) Greet ( ctx context.Context, req * greetpb.GreetRequest  ) (* greetpb.GreetResponse,error){
	fmt.Println("Invoked")
	firstname:=req.GetGreeting().GetFirstName()
	result := "Hello"+ firstname
	res := &greetpb.GreetResponse{
		Result : result,
	}
	return res,nil
}

func (*server) GreetManyTimes (req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("Invoked Greet ManyTimes")
	firstname:=req.GetGreeting().GetFirstName()
	for i:=0;i<10;i++{
		result:="hello"+firstname+"number"+strconv.Itoa(i)
		res:=&greetpb.GreetManyTimesResponse{
			Result:result,
		}
		stream.Send(res)
		time.Sleep(1000*time.Millisecond)
	}
	return nil
}

// func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer)  error {
// 	fmt.Println("Long greet function was invoked")
// 	result:="Hello"
//  for {
// 	 req,err :=stream.Recv()
// 	 if err == io.EOF{
// 		stream.SendAndClose(&greetpb.LongGreetResponse{
// 			Result :result,
// 		})
// 	 }
// 	if err != nil {
// 		 log.Fatalf("Error while reading client stream: %v",err)
// 	 }
// 	 firstName:=req.GetGreeting().GetFirstName()
// 	 result += firstName +"!    "
//  }
// }
func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Printf("LongGreet function was invoked with a streaming request\n")
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		firstName := req.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "
	}
}


func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Printf("GreetEveryone function was invoked with a streaming request\n")

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		firstName := req.GetGreeting().GetFirstName()
		result := "Hello " + firstName + "! "

		sendErr := stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

}

func (*server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Printf("GreetWithDeadline function was invoked with %v\n", req)
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.DeadlineExceeded {
			// the client canceled the request
			fmt.Println("The client canceled the request!")
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(1 * time.Second)
	}
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello " + firstName
	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}
	return res, nil
}

func main(){
	fmt.Println("hello World")
	lis, err :=net.Listen("tcp", "0.0.0.0:50051")
	if err != nil{
		log.Fatalf("Failed to listen : %v",err)
	}
	s:= grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s,&server{})

	if err := s.Serve(lis); err!= nil{
		log.Fatalf("Failed to serve %v", err)
	}
	
}