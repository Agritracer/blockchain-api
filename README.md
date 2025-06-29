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
  "sha256": [
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece",
    "ID: sp001\nIDEditor:editor123\nSTATUS: insert\nDATE: 07/06/2025\nDATA: 57459f66a1ad10050e79d9fad5bcb81e4268afcf863b7ad1fad67d3470df6ece"
  ],
  "tx_count": 35,
  "tx_hashes": [
    "0x8548d0796cd35b382a92843238b0de1182284f62a90cd5db5a25e77717b032d3",
    "0x2596417cc3b673a11ba647919cc5947451adec153f2b9e69e357f5cfd2a7d6a8",
    "0x42d746317e92cedbc90dd06721bf28db9b276430d540ca47d6b0ac2329d41509",
    "0xd1e5740a864651d31b4dc7b36bb52fad79de1762205ba629c95c8f2d52f714ee",
    "0xadf7fd5449b76310b70a0d6dbd8e19616602f88a836529dd8a22c96ec147bf41",
    "0x755bd88d950a668626949f3c1812c0b77bd69d24cbbd05c7df990b39eec83b90",
    "0x3652fdd29c901f502181932bfb9ffd970c8133e3cd2b6df801a4aac78c914ee8",
    "0xac28169091455410cedafa62bde50b53e9ed35bff9bffb2ca105c0ddb6b84ea1",
    "0x18dbca59b38d74d93457a783c59c57e7101b4a322e18ce6798b877a94f060cc0",
    "0xf61602b999be2b66bcd3e0d4e3467fb3394a3aac30066a6c774a713247217515",
    "0x79705b3a1a5ef763114d8d1d9686f4cabf4b65bb90e10e93ff6005435a2b6d92",
    "0xf04a01b320837606561a088c7fede7a38666c3d01e1585ad36ae684548006dc5",
    "0xc0981385a58871980bb5d5d5712587d32c518a26d14621f27b84bd3c7aedbafd",
    "0x0a10ea8e38c257b451eb053e8f1279c579b4e6d144e3d45daf8b94f1d2f4aa53",
    "0xd0ae4c035e4a2bc94912b929e56ff821af9cd4cfbfed209aafd3772950952bce",
    "0x0a625efe92158af2dfd7afe21fc5c4eeb802b9aec5548e57911582070d29db53",
    "0xb164b2261fb9aefe5e191f241748988f381007034f56dc2d4b982ec6e4512c68",
    "0x5e80a3049438793da62e4a0adcddfcd987b688536b817fba7164624acf7f31b4",
    "0x01dba29086035bfacf29907f01f2f4cd6bd5ae5eea53beb4881f77af30a05ea9",
    "0xb03df98f901b2b3baed11f708f2d79065e058c2f73ab5b89a5bf2e78d8e2b371",
    "0xb48b959bd288076552fdadd4c53785a7e5d16332e1a892b7b7e04cf2f90c429e",
    "0xec125ce7e7ee1ad634961ced7190f18e9012bc3b86669c7ede14c576b5aae0d5",
    "0xa089d5f37d1dc72b875de6eb85363cb2d6adfbfceca80e759aefb5f99bb1913c",
    "0x6f93dde350c63f18851fb7477aae5b7555695e7bdbdc7b54789e517ca9e66232",
    "0x6e6d1f35943b189130899ab3abf10be2701e89e63c21e298438a4bbfa6cb964e",
    "0x1d4cb21b9cfc363137e0b5e686cce8f7e61577b8c85e5eeed5f288c1729343d4",
    "0x4e9a8379e181391768e8cafa02ffe437582d9f8e88093fe205b3c1e9d5f5bf99",
    "0x4beeaf45b251feba257e0d3ce30243c7f353e8a9d3f2ebea7ea081b823e2f829",
    "0x233833ee611cae247c794b3d3ee9c496e46d24f9e51e51d0e208ceee1b0bae96",
    "0xb52ac937f06fd3b71863187f9c238287d990b93256da02644629630bf021cd6d",
    "0x7eddf8c86296d0bd6ab26deeed15a82d3ce1f17edd660e02aab858557fee7322",
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