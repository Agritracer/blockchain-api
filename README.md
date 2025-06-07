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
  "status": "Active",
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
curl -X POST http://localhost:8080/api/submit \
    -H "Content-Type: application/json" \
    -H "X-API-Key: myfrontend" \
    -d '{
          "id": "sp001",
          "idEditor": "editor123",
          "status": "Active",
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
curl "http://localhost:8080/api/trace?tx=0xfdbfd99963aa9d2d05584ed5712090f95124b6ddf3b0740315f98baa4162fc73" \
    -H "Content-Type: application/json" \
    -H "X-API-Key: myfrontend"
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
curl http://localhost:8080/api/query?id=sp001 \
    -H "Content-Type: application/json" \
    -H "X-API-Key: myfrontend"
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
### Example List:
curl:
* [** GET **]
```bash
curl http://localhost:8080/list \
    -H "Content-Type: application/json" \
    -H "X-API-Key: myfrontend"
```
output:
```json
{
  "tx_count": 39,
  "tx_hashes": [
    "0x94ea15517480a6edf35c895a0ae03539bafedca4525382d87b09bc11cb59333e",
    "0xd01ae08644a73b46065804544b91af588b280eac6bc7a82fcee2d848b88ca5a2",
    "0x192810514f2d8fa65bdefa0750042a79800d7c94ec100a5ab1dadda688fa3cf8",
    "0x624dd475af4c0ada1d259591b0df31246751d1fb7a5db61e14035bfe72eb8db4",
    "0x01d849586883a03539bc2539b84b33a04b0f0c136f71d3079ddea62c4c86575c",
    "0x3958729bbf132fafbbed452abe6f2b3ac277be884d8469ec65d0355cbbabcefd",
    "0xaf5f32fe4022c863cbc9f4c24f0ebe9a44712a933f43bec81cf328a9e805a1cd",
    "0x337644b86a610ec489daea636a34b61118cfe9353720e0194f818d1e5229d6ae",
    "0xddb2a6c4cbc1b782e71fa0425f56b9e49a14ffd6c8cac58d955cd8d51f878fdb",
    "0x231f1e833b0a74925ad9ea6efcd7379e5815cb613b79d47750ae53446312b224",
    "0xcfeb8ca434078e4fb885a8f163614826d3d5707e54e57c51d0defb813b4da7df",
    "0x763016efc46ac812731f970c8ed83ced4673a9125857064a33af96c3f18c8e1e",
    "0x752fcc36953629f1e412d8e2a6d7e045c862622fda23021c4b36375130d0886c",
    "0x9197010e80f429e0c4b07c9c285e1db8555186c0756d0ef0b8f3bebad1b2150d",
    "0xdaf63cb524998ef8adb50be1144e50f0ea3281547774cfb3edc4e2ba6f1b9ce6",
    "0xa829fd1917e116d2526ab2a16ab81d236b9f8b468f83c9c8f797dfaf5624aa2c",
    "0xc7ef470da30dbe3894cf32b18d32b4680467934f1648c0a7492fb30f92d66cde",
    "0xa9dd614942aa4e9bddce392765096a08bcb8070eacf68be6a656f2cd11838434",
    "0x25b18c06a9d51bf83b8d1a7af7636561de2e0df28666e3829f0db8a97354b376",
    "0x985168269361704e2ea22075771f33f74e565fddd5fb2a0b734d3d0edc7c83a5",
    "0xe0063bd792f1c8a924ba5b0d8df196466c5885e03ac3a374fa1b0308a5bd37c0",
    "0x085d1976f9b22fa0043437142273a08268d71dfdd0654d45f535e470d3bcab40",
    "0x7d81e3e162db71ceac351b023a7768e860b3d42d1172f560ebc6b402791f7926",
    "0x6f6a928745b3e62e12e845842b873911911dea76de8988dfe4d4622e6e5ee850",
    "0xd104c6a7a23a3d79927b4e2f7a5926c8ef223bd410bb4595e345db8f1eddfd8e",
    "0x8a2f656d417a396d8191d3611e0d353f35dcd59e03d2c19ad74e87c3b6932395",
    "0xa0cb1d1b57dead581605224922a8960f432296f5d32ff7a604b3b2e0c30b4de0",
    "0xc1494ba54452fc4cc49c14668d97f6205296d66202652ccee5a3b27df63e8a77",
    "0x072446f0ca2d45148a344c756a8c70f1bd0d70ea8eb97d58664535e62476f837",
    "0x4fd7f774681a4b9a60349cefd3262a77511afc9b7f034deba87cb9c28d51102b",
    "0x324522250205a86ee0231c15b32ac32ae2301c43acd6d84cd35246296c8484a7",
    "0xf3d7285d341b1341f991872215c6c976e26aff71e66f5b236379e3b20beb1e07",
    "0xab4c363467c6d9e10572e39c356977873381245aaeba5ae8bb1753976c8ad348",
    "0x29a788bee6cb754b22648a9df478e2fd32c9f1b13d33dd5d7a10e20f85355c16",
    "0xcec49455295e6a4fb000679acfde3b42330e2abc447aee64b8c9face26facc01",
    "0xb0baabb642471e3d0d630a11a4f249697bcd7f8a8a4dd3cab8c034b64c5087b8",
    "0xf30e41d9037ffa48239f942b5b72bf46a2afc794df35c2be3da9ef0928cc2058",
    "0x9e103418d42194e9b98352a83a324cf4e29ff71f7f6e11fc95d87a3008a87e2a",
    "0xc3ea4f3bdd9c63e0c477ccfc6f89cd7d2549914cfdb76b4b8e19694daffb9ef6"
  ],
  "wallet": "0xbC66be87455FA3AfB000fBbD009C9E403BE1122D"
}
```