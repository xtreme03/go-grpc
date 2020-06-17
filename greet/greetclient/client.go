package main

import ("fmt"
	"google.golang.org/grpc"
	"go-grpc/greet/greetpb"
	"log"
	"context"
	"io"

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
	serverStreamAPI(c)
	
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