# Taxi-backend API仕様書

**バージョン**: v1  
**ベースURL**: `http://localhost:8080`  
**最終更新日**: 2026-02-10

---

## 目次

- [共通仕様](#共通仕様)
- [ヘルスチェック](#ヘルスチェック)
- [ユーザー API](#ユーザー-api)
- [ドライバー API](#ドライバー-api)
- [配車 API](#配車-api)
- [エラー定義](#エラー定義)
- [データモデル](#データモデル)

---

## 共通仕様

### リクエスト形式

- Content-Type: `application/json`
- 文字コード: UTF-8

### レスポンス形式

- Content-Type: `application/json`
- 文字コード: UTF-8

### ID 形式

すべてのリソースIDには **ULID**（Universally Unique Lexicographically Sortable Identifier）を使用しています。

- 長さ: 26文字
- 例: `01KGWG153K2TDKE9AR81E5S048`
- 特性: 時系列ソート可能、大文字小文字を区別しない

### 共通エラーレスポンス形式

```json
{
  "error": "error_type",
  "message": "Human-readable error description"
}
```

| フィールド | 型 | 説明 |
|:---|:---|:---|
| `error` | string | エラーの種別（`validation_error`, `not_found`, `conflict`, `bad_request`, `internal_error`）|
| `message` | string | エラーの詳細メッセージ |

### HTTPステータスコード一覧

| コード | 説明 | 使用場面 |
|:---:|:---|:---|
| `200` | OK | 取得・更新成功 |
| `201` | Created | リソース作成成功 |
| `400` | Bad Request | バリデーションエラー、不正な入力値 |
| `404` | Not Found | リソースが見つからない |
| `409` | Conflict | リソースの競合（重複、状態不整合） |
| `500` | Internal Server Error | サーバー内部エラー |

---

## ヘルスチェック

### `GET /health`

サーバーの稼働状態を確認します。

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | サーバー正常稼働 |

```json
{
  "status": "ok"
}
```

---

## ユーザー API

### `POST /api/v1/users`

新しいユーザーを登録します。

**リクエストボディ**

| フィールド | 型 | 必須 | バリデーション | 説明 |
|:---|:---|:---:|:---|:---|
| `name` | string | ✅ | 1〜100文字 | ユーザー名 |
| `email` | string | ✅ | メールアドレス形式 | メールアドレス（一意） |

**リクエスト例**

```json
{
  "name": "山田太郎",
  "email": "yamada@example.com"
}
```

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `201` | ユーザー作成成功 |
| `400` | バリデーションエラー（名前が空、メール形式不正） |
| `409` | メールアドレスが既に登録済み |
| `500` | サーバー内部エラー |

**成功レスポンス例** (`201`)

```json
{
  "id": "01KGWG153K2TDKE9AR81E5S048",
  "name": "山田太郎",
  "email": "yamada@example.com"
}
```

**エラーレスポンス例** (`400`)

```json
{
  "error": "validation_error",
  "message": "Name cannot be empty"
}
```

**エラーレスポンス例** (`409`)

```json
{
  "error": "conflict",
  "message": "User with this email already exists"
}
```

---

### `GET /api/v1/users`

登録されている全ユーザーを取得します。

**パラメータ**: なし

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | ユーザー一覧取得成功 |
| `500` | サーバー内部エラー |

**成功レスポンス例** (`200`)

```json
[
  {
    "id": "01KGWG153K2TDKE9AR81E5S048",
    "name": "山田太郎",
    "email": "yamada@example.com"
  },
  {
    "id": "01KGWG1D0KTE28JC38RSX0V6RA",
    "name": "佐藤花子",
    "email": "sato@example.com"
  }
]
```

> **注意**: ユーザーが0件の場合は空配列 `[]` を返します。

---

### `GET /api/v1/users/:id`

指定したIDのユーザー情報を取得します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | ユーザーID（ULID） |

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | ユーザー取得成功 |
| `400` | IDが未指定 |
| `404` | ユーザーが見つからない |
| `500` | サーバー内部エラー |

**成功レスポンス例** (`200`)

```json
{
  "id": "01KGWG153K2TDKE9AR81E5S048",
  "name": "山田太郎",
  "email": "yamada@example.com"
}
```

---

### `GET /api/v1/users/:passenger_id/rides`

指定した乗客の配車履歴を取得します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `passenger_id` | string | ✅ | 乗客（ユーザー）ID |

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 配車履歴取得成功 |
| `400` | IDが未指定 |
| `500` | サーバー内部エラー |

**成功レスポンス例** (`200`)

```json
[
  {
    "id": "01KGXYZ12345ABCDE67890FGH",
    "passenger_id": "01KGWG153K2TDKE9AR81E5S048",
    "driver_id": "01KGWG1D0KTE28JC38RSX0V6RA",
    "pickup_latitude": 35.6762,
    "pickup_longitude": 139.6503,
    "dropoff_latitude": 35.6895,
    "dropoff_longitude": 139.6917,
    "status": "completed",
    "created_at": "2026-02-08T01:47:45.267591+09:00",
    "updated_at": "2026-02-08T02:15:30.123456+09:00"
  }
]
```

---

## ドライバー API

### `POST /api/v1/drivers`

新しいドライバーを登録します。

**リクエストボディ**

| フィールド | 型 | 必須 | バリデーション | 説明 |
|:---|:---|:---:|:---|:---|
| `name` | string | ✅ | 1〜100文字 | ドライバー名 |
| `email` | string | ✅ | メールアドレス形式 | メールアドレス（一意） |
| `license_number` | string | ✅ | 1〜50文字 | 免許番号（一意） |

**リクエスト例**

```json
{
  "name": "佐藤運転手",
  "email": "sato-driver@example.com",
  "license_number": "LICENSE-001"
}
```

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `201` | ドライバー作成成功 |
| `400` | バリデーションエラー |
| `409` | メールアドレスまたは免許番号が既に登録済み |
| `500` | サーバー内部エラー |

**成功レスポンス例** (`201`)

```json
{
  "id": "01KGWG153K2TDKE9AR81E5S048",
  "name": "佐藤運転手",
  "email": "sato-driver@example.com",
  "license_number": "LICENSE-001",
  "status": "offline",
  "latitude": 0,
  "longitude": 0,
  "created_at": "2026-02-08T01:47:45.267591+09:00",
  "updated_at": "2026-02-08T01:47:45.267591+09:00"
}
```

> **注意**: 新規作成時のステータスは必ず `offline` になります。

---

### `GET /api/v1/drivers`

登録されている全ドライバーを取得します。

**パラメータ**: なし

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | ドライバー一覧取得成功 |
| `500` | サーバー内部エラー |

**成功レスポンス例** (`200`)

```json
[
  {
    "id": "01KGWG153K2TDKE9AR81E5S048",
    "name": "佐藤運転手",
    "email": "sato-driver@example.com",
    "license_number": "LICENSE-001",
    "status": "available",
    "latitude": 35.6762,
    "longitude": 139.6503,
    "created_at": "2026-02-08T01:47:45.267591+09:00",
    "updated_at": "2026-02-08T02:00:00.000000+09:00"
  }
]
```

---

### `GET /api/v1/drivers/available`

配車可能（空車状態）なドライバーのみを取得します。

**パラメータ**: なし

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 配車可能ドライバー取得成功 |
| `500` | サーバー内部エラー |

> **返却条件**: `status` が `available` のドライバーのみ返されます。

---

### `GET /api/v1/drivers/:id`

指定したIDのドライバー情報を取得します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | ドライバーID（ULID） |

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | ドライバー取得成功 |
| `400` | IDが未指定 |
| `404` | ドライバーが見つからない |
| `500` | サーバー内部エラー |

---

### `PUT /api/v1/drivers/:id/location`

ドライバーの現在地（緯度・経度）を更新します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | ドライバーID（ULID） |

**リクエストボディ**

| フィールド | 型 | 必須 | バリデーション | 説明 |
|:---|:---|:---:|:---|:---|
| `latitude` | number | ✅ | -90 〜 90 | 緯度 |
| `longitude` | number | ✅ | -180 〜 180 | 経度 |

**リクエスト例**

```json
{
  "latitude": 35.6762,
  "longitude": 139.6503
}
```

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 位置更新成功 |
| `400` | バリデーションエラー（緯度/経度範囲外） |
| `404` | ドライバーが見つからない |
| `500` | サーバー内部エラー |

**成功レスポンス例** (`200`)

```json
{
  "id": "01KGWG153K2TDKE9AR81E5S048",
  "name": "佐藤運転手",
  "email": "sato-driver@example.com",
  "license_number": "LICENSE-001",
  "status": "available",
  "latitude": 35.6762,
  "longitude": 139.6503,
  "created_at": "2026-02-08T01:47:45.267591+09:00",
  "updated_at": "2026-02-08T02:10:00.000000+09:00"
}
```

---

### `PUT /api/v1/drivers/:id/status`

ドライバーの稼働ステータスを更新します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | ドライバーID（ULID） |

**リクエストボディ**

| フィールド | 型 | 必須 | バリデーション | 説明 |
|:---|:---|:---:|:---|:---|
| `status` | string | ✅ | `available`, `busy`, `offline` のいずれか | 新しいステータス |

**リクエスト例**

```json
{
  "status": "available"
}
```

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | ステータス更新成功 |
| `400` | バリデーションエラー（不正なステータス値） |
| `404` | ドライバーが見つからない |
| `500` | サーバー内部エラー |

---

### `GET /api/v1/drivers/:driver_id/rides`

指定したドライバーの配車履歴を取得します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `driver_id` | string | ✅ | ドライバーID |

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 配車履歴取得成功 |
| `400` | IDが未指定 |
| `500` | サーバー内部エラー |

> **並び順**: `created_at` 降順（新しい順）

---

## 配車 API

### ステータス遷移図

```
                         ┌──────────────────┐
                         │                  │
   ┌────────────┐   accept    ┌──────────┐  arrive  ┌─────────┐  start   ┌─────────┐  complete  ┌───────────┐
   │ requested  │ ──────────→ │ accepted │ ───────→ │ arrived │ ──────→ │ ongoing │ ────────→ │ completed │
   └────────────┘             └──────────┘          └─────────┘         └─────────┘           └───────────┘
        │                          │                     │                   │
        │        cancel            │      cancel         │     cancel        │
        └─────────────┬────────────┴─────────────────────┴───────────────────┘
                      │
                      ▼
               ┌────────────┐
               │ cancelled  │
               └────────────┘
```

**ステータス定義**

| ステータス | 説明 | 遷移先 |
|:---|:---|:---|
| `requested` | 配車依頼中（乗客がリクエスト送信済み） | `accepted`, `cancelled` |
| `accepted` | ドライバーが承諾済み | `arrived`, `cancelled` |
| `arrived` | ドライバーが乗車地点に到着 | `ongoing`, `cancelled` |
| `ongoing` | 乗車中（走行中） | `completed`, `cancelled` |
| `completed` | 乗車完了 | _(終端ステータス)_ |
| `cancelled` | キャンセル済み | _(終端ステータス)_ |

---

### `POST /api/v1/rides`

新しい配車リクエストを作成します。

**リクエストボディ**

| フィールド | 型 | 必須 | バリデーション | 説明 |
|:---|:---|:---:|:---|:---|
| `passenger_id` | string | ✅ | 存在するユーザーID | 乗客のID |
| `pickup_latitude` | number | ✅ | -90 〜 90 | 乗車地点の緯度 |
| `pickup_longitude` | number | ✅ | -180 〜 180 | 乗車地点の経度 |
| `dropoff_latitude` | number | ✅ | -90 〜 90 | 降車地点の緯度 |
| `dropoff_longitude` | number | ✅ | -180 〜 180 | 降車地点の経度 |

**リクエスト例**

```json
{
  "passenger_id": "01KGWG153K2TDKE9AR81E5S048",
  "pickup_latitude": 35.6762,
  "pickup_longitude": 139.6503,
  "dropoff_latitude": 35.6895,
  "dropoff_longitude": 139.6917
}
```

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `201` | 配車リクエスト作成成功 |
| `400` | バリデーションエラー（緯度/経度範囲外） |
| `404` | 乗客（ユーザー）が見つからない |
| `409` | 乗客に進行中の配車が既に存在 |
| `500` | サーバー内部エラー |

**成功レスポンス例** (`201`)

```json
{
  "id": "01KGXYZ12345ABCDE67890FGH",
  "passenger_id": "01KGWG153K2TDKE9AR81E5S048",
  "pickup_latitude": 35.6762,
  "pickup_longitude": 139.6503,
  "dropoff_latitude": 35.6895,
  "dropoff_longitude": 139.6917,
  "status": "requested",
  "created_at": "2026-02-08T01:47:45.267591+09:00",
  "updated_at": "2026-02-08T01:47:45.267591+09:00"
}
```

> **注意**:
> - 作成時の `status` は必ず `requested`
> - `driver_id` はマッチング前のため省略されます
> - 1人の乗客が同時に持てる進行中の配車は1件のみ

---

### `GET /api/v1/rides/:id`

指定したIDの配車情報を取得します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | 配車ID（ULID） |

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 配車情報取得成功 |
| `400` | IDが未指定 |
| `404` | 配車が見つからない |
| `500` | サーバー内部エラー |

---

### `PUT /api/v1/rides/:id/accept`

ドライバーが配車リクエストを承諾します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | 配車ID（ULID） |

**リクエストボディ**

| フィールド | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `driver_id` | string | ✅ | 承諾するドライバーのID |

**リクエスト例**

```json
{
  "driver_id": "01KGWG1D0KTE28JC38RSX0V6RA"
}
```

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 承諾成功 |
| `400` | 配車が `requested` ステータスでない |
| `404` | 配車またはドライバーが見つからない |
| `409` | ドライバーが配車可能でない / ドライバーに進行中の配車が既に存在 |
| `500` | サーバー内部エラー |

**副作用**:
- 配車のステータスが `requested` → `accepted` に変更
- 配車に `driver_id` が設定される
- ドライバーのステータスが `available` → `busy` に自動変更

> **前提条件**:
> - 配車のステータスが `requested` であること
> - ドライバーのステータスが `available` であること
> - ドライバーに進行中の配車がないこと

---

### `PUT /api/v1/rides/:id/arrive`

ドライバーが乗車地点に到着したことを記録します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | 配車ID（ULID） |

**リクエストボディ**: なし

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 到着記録成功 |
| `400` | 配車が `accepted` ステータスでない |
| `404` | 配車が見つからない |
| `500` | サーバー内部エラー |

**副作用**:
- 配車のステータスが `accepted` → `arrived` に変更

---

### `PUT /api/v1/rides/:id/start`

乗車を開始します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | 配車ID（ULID） |

**リクエストボディ**: なし

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 乗車開始成功 |
| `400` | 配車が `arrived` ステータスでない |
| `404` | 配車が見つからない |
| `500` | サーバー内部エラー |

**副作用**:
- 配車のステータスが `arrived` → `ongoing` に変更

---

### `PUT /api/v1/rides/:id/complete`

乗車を完了します。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | 配車ID（ULID） |

**リクエストボディ**: なし

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | 乗車完了成功 |
| `400` | 配車が `ongoing` ステータスでない |
| `404` | 配車が見つからない |
| `500` | サーバー内部エラー |

**副作用**:
- 配車のステータスが `ongoing` → `completed` に変更
- ドライバーのステータスが `busy` → `available` に自動変更

---

### `PUT /api/v1/rides/:id/cancel`

配車をキャンセルします。

**パスパラメータ**

| パラメータ | 型 | 必須 | 説明 |
|:---|:---|:---:|:---|
| `id` | string | ✅ | 配車ID（ULID） |

**リクエストボディ**: なし

**レスポンス**

| ステータス | 説明 |
|:---:|:---|
| `200` | キャンセル成功 |
| `400` | 配車が `completed` または `cancelled` ステータス（キャンセル不可） |
| `404` | 配車が見つからない |
| `500` | サーバー内部エラー |

**副作用**:
- 配車のステータスが `cancelled` に変更
- ドライバーが割り当てられている場合、ドライバーのステータスが `busy` → `available` に自動変更

> **注意**: `completed` および `cancelled` ステータスの配車はキャンセルできません。

---

## エラー定義

### エラー種別一覧

| error フィールド | HTTPステータス | 説明 |
|:---|:---:|:---|
| `validation_error` | 400 | 入力値のバリデーションに失敗 |
| `bad_request` | 400 | リクエスト内容が不正 |
| `not_found` | 404 | リソースが見つからない |
| `conflict` | 409 | リソースの競合 |
| `internal_error` | 500 | サーバー内部エラー |

### エラーメッセージ一覧

#### ユーザー関連

| メッセージ | ステータス | トリガー |
|:---|:---:|:---|
| `Name cannot be empty` | 400 | 名前が空 |
| `Invalid email format` | 400 | メールアドレスの形式不正 |
| `User with this email already exists` | 409 | メールアドレスの重複 |
| `User not found` | 404 | 指定IDのユーザーが存在しない |
| `User ID is required` | 400 | IDパラメータが未指定 |

#### ドライバー関連

| メッセージ | ステータス | トリガー |
|:---|:---:|:---|
| `Name cannot be empty` | 400 | 名前が空 |
| `Invalid email format` | 400 | メールアドレスの形式不正 |
| `License number cannot be empty` | 400 | 免許番号が空 |
| `Driver with this email or license number already exists` | 409 | メールまたは免許番号の重複 |
| `Driver not found` | 404 | 指定IDのドライバーが存在しない |
| `Latitude must be between -90 and 90` | 400 | 緯度範囲外 |
| `Longitude must be between -180 and 180` | 400 | 経度範囲外 |
| `Invalid driver status` | 400 | 不正なステータス値 |

#### 配車関連

| メッセージ | ステータス | トリガー |
|:---|:---:|:---|
| `Passenger not found` | 404 | 乗客IDに該当するユーザーが存在しない |
| `Passenger already has an active ride` | 409 | 乗客に進行中の配車がある |
| `Invalid latitude` | 400 | 緯度範囲外 |
| `Invalid longitude` | 400 | 経度範囲外 |
| `Ride not found` | 404 | 指定IDの配車が存在しない |
| `Driver is not available` | 409 | ドライバーが空車でない |
| `Driver already has an active ride` | 409 | ドライバーに進行中の配車がある |
| `Ride cannot be accepted in current status` | 400 | 承諾不可のステータス |
| `Invalid ride status transition` | 400 | 不正なステータス遷移 |

---

## データモデル

### UserResponse

| フィールド | 型 | 説明 |
|:---|:---|:---|
| `id` | string | ユーザーID（ULID, 26文字） |
| `name` | string | ユーザー名 |
| `email` | string | メールアドレス |

### DriverResponse

| フィールド | 型 | 説明 |
|:---|:---|:---|
| `id` | string | ドライバーID（ULID, 26文字） |
| `name` | string | ドライバー名 |
| `email` | string | メールアドレス |
| `license_number` | string | 免許番号 |
| `status` | string | 稼働ステータス（`available` / `busy` / `offline`） |
| `latitude` | number | 現在地の緯度 |
| `longitude` | number | 現在地の経度 |
| `created_at` | string | 作成日時（ISO 8601） |
| `updated_at` | string | 更新日時（ISO 8601） |

### DriverStatus（列挙値）

| 値 | 説明 |
|:---|:---|
| `available` | 空車（配車可能） |
| `busy` | 実車（乗客対応中） |
| `offline` | オフライン（稼働停止） |

### RideResponse

| フィールド | 型 | 説明 |
|:---|:---|:---|
| `id` | string | 配車ID（ULID, 26文字） |
| `passenger_id` | string | 乗客のID |
| `driver_id` | string | ドライバーのID（承諾前は省略） |
| `pickup_latitude` | number | 乗車地点の緯度 |
| `pickup_longitude` | number | 乗車地点の経度 |
| `dropoff_latitude` | number | 降車地点の緯度 |
| `dropoff_longitude` | number | 降車地点の経度 |
| `status` | string | 配車ステータス |
| `created_at` | string | 作成日時（ISO 8601） |
| `updated_at` | string | 更新日時（ISO 8601） |

### RideStatus（列挙値）

| 値 | 説明 |
|:---|:---|
| `requested` | 配車依頼中 |
| `accepted` | ドライバーが承諾 |
| `arrived` | ドライバーが乗車地点に到着 |
| `ongoing` | 乗車中（走行中） |
| `completed` | 乗車完了 |
| `cancelled` | キャンセル済み |

### ErrorResponse

| フィールド | 型 | 説明 |
|:---|:---|:---|
| `error` | string | エラー種別 |
| `message` | string | エラーの詳細メッセージ（省略される場合あり） |

---

## 使用例（curl）

### 配車フロー全体の実行例

```bash
# 1. 乗客を作成
curl -s -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"山田太郎", "email":"yamada@example.com"}' | jq .

# 2. ドライバーを作成
curl -s -X POST http://localhost:8080/api/v1/drivers \
  -H "Content-Type: application/json" \
  -d '{"name":"佐藤運転手", "email":"sato@example.com", "license_number":"LICENSE-001"}' | jq .

# 3. ドライバーを空車にする
curl -s -X PUT http://localhost:8080/api/v1/drivers/{DRIVER_ID}/status \
  -H "Content-Type: application/json" \
  -d '{"status":"available"}' | jq .

# 4. ドライバーの位置を更新
curl -s -X PUT http://localhost:8080/api/v1/drivers/{DRIVER_ID}/location \
  -H "Content-Type: application/json" \
  -d '{"latitude":35.6762, "longitude":139.6503}' | jq .

# 5. 配車リクエスト作成
curl -s -X POST http://localhost:8080/api/v1/rides \
  -H "Content-Type: application/json" \
  -d '{
    "passenger_id":"{USER_ID}",
    "pickup_latitude":35.6762,
    "pickup_longitude":139.6503,
    "dropoff_latitude":35.6895,
    "dropoff_longitude":139.6917
  }' | jq .

# 6. ドライバーが承諾
curl -s -X PUT http://localhost:8080/api/v1/rides/{RIDE_ID}/accept \
  -H "Content-Type: application/json" \
  -d '{"driver_id":"{DRIVER_ID}"}' | jq .

# 7. ドライバーが到着
curl -s -X PUT http://localhost:8080/api/v1/rides/{RIDE_ID}/arrive | jq .

# 8. 乗車開始
curl -s -X PUT http://localhost:8080/api/v1/rides/{RIDE_ID}/start | jq .

# 9. 乗車完了
curl -s -X PUT http://localhost:8080/api/v1/rides/{RIDE_ID}/complete | jq .

# 10. 配車履歴を確認
curl -s http://localhost:8080/api/v1/users/{USER_ID}/rides | jq .
curl -s http://localhost:8080/api/v1/drivers/{DRIVER_ID}/rides | jq .
```
