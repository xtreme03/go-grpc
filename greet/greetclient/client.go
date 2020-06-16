package main

import ("fmt"
	"google.golang.org/grpc"
	"go-grpc/greet/greetpb"
	"log"
	"context"

)

func main(){
	fmt.Println("Hello I'm a client")
	conn,err:=grpc.Dial("localhost:50051",grpc.WithInsecure() )
	if err != nil {
		log.Fatalf("could not connect %v",err)

	}
	defer conn.Close()
	c:=greetpb.NewGreetServiceClient(conn)
	uninaryAPI(c)
	
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