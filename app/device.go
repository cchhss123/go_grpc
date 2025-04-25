package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"
	"io"

	"google.golang.org/grpc"

	pb "go_grpc/temperature"
)

const (
	serverAddress = "localhost:50051"
	deviceIdPrefix = "device-"
)

func main() {
	// 模擬裝置 ID
	rand.Seed(time.Now().UnixNano())
	deviceId := fmt.Sprintf("%s%d", deviceIdPrefix, rand.Intn(100))

	// 連接到 gRPC 伺服器
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure()) // 注意：生產環境應使用 TLS
	if err != nil {
		log.Fatalf("無法連接到伺服器: %v", err)
	}
	defer conn.Close()

	client := pb.NewTemperatureMonitorClient(conn)

	// 創建串流
	stream, err := client.SendTemperature(context.Background())
	if err != nil {
		log.Fatalf("創建串流失敗: %v", err)
	}
	defer stream.CloseSend()

	log.Printf("裝置 %s 開始傳輸溫度數據...\n", deviceId)

	// 模擬定期發送溫度數據
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		temperature := 20.0 + rand.Float32()*15.0 // 模擬 20-35 °C 的隨機溫度
		timestamp := time.Now().Unix()

		reading := &pb.TemperatureReading{
			DeviceId:    deviceId,
			Temperature: temperature,
			Timestamp:   timestamp,
		}

		log.Printf("裝置 %s 發送溫度: %.2f °C，時間戳: %d\n", deviceId, temperature, timestamp)

		if err := stream.Send(reading); err != nil {
			log.Printf("發送數據失敗: %v", err)
			break
		}

		// 接收伺服器的確認訊息 (非阻塞)
		if ack, err := stream.Recv(); err == nil {
			log.Printf("裝置 %s 收到確認: %s\n", deviceId, ack.Message)
		} else if err != io.EOF {
			log.Printf("接收確認訊息錯誤: %v", err)
		}
	}

	log.Printf("裝置 %s 停止傳輸。\n", deviceId)
}