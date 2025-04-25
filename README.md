# 多裝置與伺服器溫度監測 gRPC 模擬程式

這個專案使用 Go 語言和 gRPC 框架，模擬多個裝置向一個伺服器傳輸溫度監測數據的場景。

## 檔案結構

專案包含以下檔案：

1.  **`temperature.proto`**: 定義了 gRPC 服務的 Protocol Buffer 描述，包括 `TemperatureReading` 消息和 `TemperatureMonitor` 服務。
2.  **`server.go`**: 實現了 gRPC 伺服器，用於接收來自裝置的溫度數據並回覆確認訊息。
3.  **`device.go`**: 模擬一個溫度監測裝置，定期向伺服器發送隨機生成的溫度數據並接收確認訊息。
4.  **`Dockerfile`**: 用於構建包含所有必要依賴的 Docker 鏡像。
5.  **`docker-compose.yaml`**: 用於使用 Docker Compose 方便地啟動和管理服務容器。

## 如何使用 Docker 進行開發與測試

本專案提供 Dockerfile 和 docker-compose.yaml 文件，方便您在容器化的環境中進行開發和測試。

### 1. 建構 Docker 鏡像

在專案根目錄下（包含 `Dockerfile`）執行以下命令來建構 Docker 鏡像。`--no-cache` 選項可以確保每次都從頭開始構建鏡像，避免緩存問題。

```bash
docker build -t go_grpc . --no-cache
```

這個命令將會下載 Golang 基礎鏡像，安裝必要的依賴（包括 `protoc` 和 gRPC Go 外掛），複製程式碼到容器中，並設定工作目錄。

### 2. 運行 gRPC 服務容器

在專案根目錄下（包含 `docker-compose.yaml`）執行以下命令來啟動 gRPC 伺服器容器：

```bash
docker-compose up
```

這將會使用 `docker-compose.yaml` 中定義的 `go_grpc` 服務來啟動一個容器。伺服器將會在容器內的 `50051` 端口監聽，並且宿主機的 `50051` 端口也會映射到這個容器端口。

### 3. 生成 gRPC 程式碼 (在容器內)

進入運行的 `go_grpc` 容器並執行 `protoc` 命令：

```bash
docker ps # 找出運行容器ID
docker exec -it <go_grpc_container_id_or_name> sh
```

然後在容器的 `/app` 目錄下執行：

```bash
mkdir temperature

protoc --go_out=./temperature --go_opt=paths=source_relative --go-grpc_out=./temperature --go-grpc_opt=paths=source_relative ./temperature.proto
```
將在 /app 目錄下，建立 `temperature` 子目錄，並生成 `temperature.pb.go` 和 `temperature_grpc.pb.go` 兩個檔案


### 4. 測試運行

進入運行的 `go_grpc` 容器：

```bash
docker ps # 找出運行容器ID
docker exec -it <go_grpc_container_id_or_name> sh
```

初始化產生go.mod，容器內執行以下命令
```bash
/app # go mod init go_grpc
/app # go mod tidy
```

先運行 Server 端：
```bash
go run server.go
```

然後在不同的終端機中運行多個 Client 端 (模擬多個裝置)：
在不同的終端機中進入運行的 `go_grpc` 容器：

```bash
docker ps # 找出運行容器ID
docker exec -it <go_grpc_container_id_or_name> sh
```

在不同的 `go_grpc` 容器，運行 多個不同的 Client 端：
```bash
go run device.go
```

將看到 Server 端接收來自不同裝置的溫度數據，並回覆確認訊息。
每個 Client 端 (模擬裝置) 會定期生成並發送隨機溫度數據。


## 程式說明

### `temperature.proto`

* 定義了 `TemperatureReading` 消息，包含 `device_id` (字串)、`temperature` (浮點數) 和 `timestamp` (64 位整數)。
* 定義了 `TemperatureMonitor` 服務，包含一個雙向串流 RPC 方法 `SendTemperature`，裝置可以向伺服器發送 `TemperatureReading` 串流，伺服器可以回覆 `Ack` 串流。
* 定義了 `Ack` 消息，包含一個確認訊息 `message` (字串)。

### `server.go`

* 啟動一個 gRPC 伺服器並監聽在 `localhost:50051`。
* 實現了 `TemperatureMonitor` 服務的 `SendTemperature` 方法。
* 在接收到來自裝置的 `TemperatureReading` 後，會記錄相關資訊並回覆一個包含確認訊息的 `Ack`。

### `device.go`

* 模擬一個溫度監測裝置，生成一個隨機的 `device_id`。
* 連接到在 `localhost:50051` 運行的 gRPC 伺服器。
* 定期（每 2 秒）生成一個隨機溫度和目前的時間戳。
* 將 `TemperatureReading` 消息發送到伺服器。
* 接收並記錄伺服器回覆的 `Ack` 訊息。
