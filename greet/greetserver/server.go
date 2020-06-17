package main

import (
	"strconv"
	"fmt"
	"context"
	"net"
	"log"
	"google.golang.org/grpc"
	"go-grpc/greet/greetpb"
	"time"
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