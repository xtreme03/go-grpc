package main

import ("fmt"
	"google.golang.org/grpc"
	"go-grpc/calculator/calculatorpb"
	"log"
	"context"
	"io"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

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
	//serverStreamAPI(c)
	doErrorUnary(c)
	
}
func doErrorUnary(c calculatorpb.CalculateServiceClient) {
	fmt.Println("Starting to do a SquareRoot Unary RPC...")

	// correct call
	doErrorCall(c, 10)

	// error call
	doErrorCall(c, -2)
}

func doErrorCall(c calculatorpb.CalculateServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})

	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Printf("Error message from server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number!")
				return
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot: %v", err)
			return
		}
	}
	fmt.Printf("Result of square root of %v: %v\n", n, res.GetNumberRoot())
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