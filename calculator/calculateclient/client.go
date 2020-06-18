package main

import ("fmt"
	"google.golang.org/grpc"
	"go-grpc/calculator/calculatorpb"
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
	c:=calculatorpb.NewCalculateServiceClient(conn)
	//uninaryAPI(c)
	serverStreamAPI(c)
	
}

func uninaryAPI(c calculatorpb.CalculateServiceClient){
	req := &calculatorpb.CalculateRequest{
		Value : &calculatorpb.Calculate{
			Num1:12,
			Num2:13,
		},

	}


	
	//fmt.Printf("Client %f",c)
	res,err :=c.Sum(context.Background(),req)
	if err != nil{
		log.Fatal("Server Error")
	}
	fmt.Printf("Response : %v",res)
}

func serverStreamAPI(c calculatorpb.CalculateServiceClient){
	fmt.Println("Call server stream Api")
	req:=&calculatorpb.PrimeDeRequest{
		Num: 150,
	}
	resStream,err :=c.PrimeDecomposition(context.Background(),req)
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