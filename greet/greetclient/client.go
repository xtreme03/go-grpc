package main

import ("fmt"
	"google.golang.org/grpc"
	"go-grpc/greet/greetpb"
	"log"
	"context"
	"io"
	"time"

)

func main(){
	fmt.Println("Hello I'm a client")
	conn,err:=grpc.Dial("localhost:50051",grpc.WithInsecure() )
	if err != nil {
		log.Fatalf("could not connect %v",err)

	}
	defer conn.Close()
	c:=greetpb.NewGreetServiceClient(conn)
	//uninaryAPI(c)
	//serverStreamAPI(c)
	//clientStreamingAPI(c)
	doBiDiStreaming(c)
	
}
func clientStreamingAPI(c greetpb.GreetServiceClient){
	fmt.Println("Client Streaming")
	requests:=[]*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:"Pall",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:"Pall1",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:"Pall2",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName:"Pall3",
			},
		},

	}

	stream,err := c.LongGreet(context.Background())
	if err != nil{
		log.Fatalf("Error while stream : %v",err)
	}
	for _,req :=range requests{
		fmt.Printf("Sending:%v",req)
		stream.Send(req)
		//time.Sleep(1000*time.Millisecond)
	}

	res,err:= stream.CloseAndRecv()
	if err != nil{
		log.Fatalf("Error : %v",err)
	}
	log.Printf("Response:%v",res)
}

func uninaryAPI(c greetpb.GreetServiceClient){
	req := &greetpb.GreetRequest{
		Greeting : &greetpb.Greeting{
			FirstName:"Pallab",
			LastName:"Nag",
		},

	}
	//fmt.Printf("Client %f",c)
	res,err :=c.Greet(context.Background(),req)
	if err != nil{
		log.Fatal("Server Error")
	}
	fmt.Printf("Response : %v",res)
}
func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephane",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
			},
		},
	}

	waitc := make(chan struct{})
	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, req := range requests {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving: %v", err)
				break
			}
			fmt.Printf("Received: %v\n", res.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}

func serverStreamAPI(c greetpb.GreetServiceClient){
	fmt.Println("Call server stream Api")
	req:=&greetpb.GreetManyTimesRequest{
		Greeting : &greetpb.Greeting{
			FirstName:"Pallab",
			LastName:"Nag",
		},
	}
	resStream,err :=c.GreetManyTimes(context.Background(),req)
	if err != nil {
		log.Fatalf("Error server streaming %v",err)
	}
	for {
		msg,err:=resStream.Recv()
		if err == io.EOF{
			break
		}
		if err != nil {
			log.Fatalf("Error: %v",err)
		}
		fmt.Println(msg.GetResult())


}
}