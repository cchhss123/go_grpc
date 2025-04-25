package main

import (
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "go_grpc/temperature"
)

const (
	port = ":50051"
)

// server 實現了 TemperatureMonitor 服務
type server struct {
	pb.UnimplementedTemperatureMonitorServer
}

// SendTemperature 處理來自裝置的溫度數據流
func (s *server) SendTemperature(stream pb.TemperatureMonitor_SendTemperatureServer) error {
	for {
		reading, err := stream.Recv()
		if err == io.EOF {
			return nil // 裝置已關閉串流
		}
		if err != nil {
			log.Printf("接收數據錯誤: %v", err)
			return err
		}

		log.Printf("收到來自裝置 %s 的溫度數據: %.2f °C，時間戳: %d\n", reading.DeviceId, reading.Temperature, reading.Timestamp)

		// 在這裡您可以將接收到的溫度數據儲存到資料庫或其他地方

		// 回覆確認訊息
		ack := &pb.Ack{Message: fmt.Sprintf("Server 收到來自裝置 %s 的數據", reading.DeviceId)}
		if err := stream.Send(ack); err != nil {
			log.Printf("發送確認訊息錯誤: %v", err)
			return err
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("監聽失敗: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTemperatureMonitorServer(s, &server{})

	// 在 gRPC 伺服器上註冊 reflection 服務，允許客戶端查詢服務的元數據
	reflection.Register(s)

	log.Printf("伺服器監聽於 %s\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("伺服器啟動失敗: %v", err)
	}
}