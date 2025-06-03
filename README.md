# BlockChain API
## Feature
Dưới đây là danh sách các chức năng (API endpoints và logic liên quan) đã được xây dựng cho hệ thống truy xuất nguồn gốc nông sản sử dụng Blockchain Ethereum (Sepolia testnet):

1. Submit dữ liệu truy xuất nguồn gốc (POST /submit)
* Nhận một JSON gồm id, idEditor và data (thông tin sản phẩm).
* Tính hash SHA256 của data.
* Gửi thông tin gồm ID, ngày tạo và hash lên mạng Ethereum bằng giao dịch.
* Lưu thông tin ID ↔ tx\_hash để phục vụ truy vấn.
* Trả về SHA256 và tx\_hash của giao dịch.

2. Truy xuất thông tin theo tx\_hash (GET /trace?tx=...)
* Nhận tx\_hash từ client.
* Gọi Etherscan API hoặc RPC để lấy nội dung giao dịch từ Ethereum.
* Phân tích lại payload (data của transaction) và trả về thông tin ban đầu đã ghi.

3. Truy vấn theo ID (GET /query?id=...)
* Nhận ID sản phẩm từ client.
* Truy tìm toàn bộ các tx\_hash có liên quan đến ID đó (các lần submit trước đó).
* Trả về danh sách tx\_hash và số lần thay đổi dữ liệu.

4. Danh sách tất cả giao dịch đã ghi (GET /list)
* Liệt kê toàn bộ các tx\_hash và các ID tương ứng.
* Dùng để kiểm tra toàn cục dữ liệu đã ghi lên blockchain.

5. Cấu hình hệ thống:
* Sử dụng file .env để cấu hình các biến như:
  * ETHEREUM\_RPC
  * PRIVATE\_KEY
  * ETHERSCAN\_API\_KEY
* Có module config nội bộ để load thông tin này từ file .env.

6. Cấu trúc thư mục chuẩn:
* internal/config: Load config từ .env
* internal/model: Định nghĩa các struct dữ liệu như InputData, Response,...
* internal/service: Logic xử lý dữ liệu (hash, gọi Ethereum, lưu tx,...)
* internal/eth: Tương tác với Ethereum Sepolia (gửi tx, lấy tx,...)
* internal/storage: Lưu mapping ID ↔ tx\_hash (có thể thay bằng Etherscan query hoặc Redis, BoltDB sau này)
* cmd/server: Main HTTP server

7. Tích hợp Etherscan API
* Truy xuất nội dung giao dịch bằng tx\_hash từ Etherscan để tránh phụ thuộc RPC.
* Chuẩn hóa truy vấn và xử lý nội dung trả về.

8. Hỗ trợ xác thực 2FA giữa frontend và API server (ý tưởng đã lên, triển khai sau):
* Chỉ cho phép truy vấn hoặc submit khi cả frontend và backend xác thực cùng một mã TOTP hợp lệ.


## Example
### Example Data:
```Json
{
  "id": "sp001",
  "idEditor": "editor123",
  "data": {
    "name": "Sầu riêng Ri6",
    "origin": "Đắk Lắk",
    "certified": true,
    "farmer": "Nguyễn Văn A",
    "harvest_date": "2025-06-01",
    "weight_kg": 2.5
  }
}
```
### Example Submit:
curl:
* [** POST **]
```bash
curl -X POST http://localhost:8080/submit \
    -H "Content-Type: application/json" \
    -H "X-API-Key: myfrontend" \
    -H "X-TOTP-Code: 123456" \
    -d '{
          "id": "sp001",
          "idEditor": "editor123",
          "data": {
            "name": "Sầu riêng Ri6",
            "origin": "Đắk Lắk",
            "certified": true,
            "farmer": "Nguyễn Văn A",
            "weight_kg": 2.5
          }
        }'
```
output:
```json
{
  "id": "sp001",
  "idEditor": "editor123",
  "time": "03/06/2025",
  "sha256": "57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
  "tx_hash": "0x5e80a3049438793da62e4a0adcddfcd987b688536b817fba7164624acf7f31b4"
}
```
### Example Trace:
curl:
* [** GET **]
```bash
curl "http://localhost:8080/trace?tx=0xfdbfd99963aa9d2d05584ed5712090f95124b6ddf3b0740315f98baa4162fc73" \
    -H "Content-Type: application/json" \
    -H "X-API-Key: myfrontend" \
    -H "X-TOTP-Code: 123456"
```
output:
```json
{
  "data": "ID: sp001\nDATE: 03/06/2025\nDATA: 9bf6d6bb85d68d0f39c2a4e23061f5c3155ca24cb8cb3179381e3b041b3ff83b",
  "tx_hash": "0xfdbfd99963aa9d2d05584ed5712090f95124b6ddf3b0740315f98baa4162fc73"
}
```
### Example Query:
curl:
* [** GET **]
```bash
curl http://localhost:8080/query?id=sp001 \
    -H "Content-Type: application/json" \
    -H "X-API-Key: myfrontend" \
    -H "X-TOTP-Code: 123456"
```
output:
```json
{
  "id": "sp001",
  "tx_count": 4,
  "tx_hashes": [
    "0xfdbfd99963aa9d2d05584ed5712090f95124b6ddf3b0740315f98baa4162fc73",
    "0x28a7bb37ed1f9fc0326b99b726a472e871613f15e512e3f582293e18bf8f6eb8",
    "0xb34127e6499881792b876358e95b60ed138f7ab776073f450ad5d783936527ce",
    "0x01ca9542416d37197c5a4fd8f4a68e69a6fdc6fc2f1692eb499ad62870e7f072"
  ]
}
```