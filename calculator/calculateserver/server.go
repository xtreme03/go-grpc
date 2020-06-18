package main

import (
	"fmt"
	"context"
	"net"
	"log"
	"google.golang.org/grpc"
	"go-grpc/calculator/calculatorpb"
)
type server struct{}
func (*server) Sum ( ctx context.Context, req * calculatorpb.CalculateRequest  ) (* calculatorpb.CalculateResponse,error){
	fmt.Println("Invoked")
	num1:=req.GetValue().GetNum1()
	num2:=req.GetValue().GetNum2()
	result := num1+num2
	res := &calculatorpb.CalculateResponse{
		Result : result,
	}
	return res,nil
}

func (*server) PrimeDecomposition (req * calculatorpb.PrimeDeRequest,stream calculatorpb.CalculateService_PrimeDecompositionServer)  (error){
	fmt.Println("Prime decomposition")
	num:=req.GetNum()
	fmt.Println(num)
	k:=int32(2)
	for num>1{
		if num%k==0{
			res:=&calculatorpb.PrimeDeResponse{
				Result:k,
			}
			num=num/k
			stream.Send(res)

		} else {
			k=k+1

		}
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
	calculatorpb.RegisterCalculateServiceServer(s,&server{})

	if err := s.Serve(lis); err!= nil{
		log.Fatalf("Failed to serve %v", err)
	}
	
}