# BlockChain API
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