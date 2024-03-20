package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	ecpb "example.com/testy/ec_site_grpc/pkg/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type ecCartServer struct {
	ecpb.UnimplementedCartServiceServer
}

func (*ecCartServer) GetProducts(ctx context.Context, req *ecpb.GetProductsRequest) (*ecpb.GetProductsResponse, error) {
	// リクエストから値を取り出す
	// user_id := req.GetUserId

	// 該当ユーザーのカート内商品IDのリストを返す(固定値を返す)
	return &ecpb.GetProductsResponse{
		Product: []*ecpb.GetProductsResponse_Product{
			{
				ProductId: "A0001",
				Qty:       1,
			},
			{
				ProductId: "A0002",
				Qty:       5,
			},
			{
				ProductId: "B0001",
				Qty:       1,
			},
		},
	}, nil
}

func NewEcCartServer() *ecCartServer {
	return &ecCartServer{}
}

type ecProductServer struct {
	ecpb.UnimplementedProductServiceServer
}

func (*ecProductServer) GetInfo(ctx context.Context, req *ecpb.GetInfoRequest) (*ecpb.GetInfoResponse, error) {
	// リクエストから値を取り出す
	product_id := req.GetId()

	// 商品データべた書き
	products := map[string]ecpb.GetInfoResponse{
		"A0001": {
			Id:    "A0001",
			Name:  "世界一流エンジニアの思考法",
			Genre: ecpb.Genre_BOOK,
			Price: 1700,
		},
		"A0002": {
			Id:    "A0002",
			Name:  "Clean Architecture 達人に学ぶソフトウェアの構造と設計",
			Genre: ecpb.Genre_BOOK,
			Price: 3168,
		},
		"B0001": {
			Id:    "B0001",
			Name:  "映画『THE FIRST SLAM DUNK』STANDARD EDITION",
			Genre: ecpb.Genre_DVD,
			Price: 3523,
		},
	}

	// 商品IDにマッチしたデータを返す
	return &ecpb.GetInfoResponse{
		Id:    products[product_id].Id,
		Name:  products[product_id].Name,
		Genre: products[product_id].Genre,
		Price: products[product_id].Price,
	}, nil
}

func NewEcProductServer() *ecProductServer {
	return &ecProductServer{}
}

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer()

	// 3. gRPCサーバーにGreetingServiceを登録
	ecpb.RegisterCartServiceServer(s, NewEcCartServer())
	ecpb.RegisterProductServiceServer(s, NewEcProductServer())

	// リフレクションの設定
	reflection.Register(s)

	// 3. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 4.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
