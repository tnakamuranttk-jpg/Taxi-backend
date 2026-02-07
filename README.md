# Taxi Backend

タクシー配車アプリのバックエンドシステムです。
Go言語（Ginフレームワーク）とMySQLを使用し、**クリーンアーキテクチャ**に基づいて構築されています。

## プロジェクト概要
このプロジェクトは、タクシーの配車管理、ユーザー管理、走行状態の追跡などを行うAPIサーバーの最小構成を提供します。

## 技術スタック
- **言語**: Go (1.25.7)
- **フレームワーク**: Gin
- **ORM**: GORM
- **データベース**: MySQL 8.0
- **インフラ**: Docker / Docker Compose
- **ID生成**: ULID

## ディレクトリ構成（クリーンアーキテクチャ）

```
├── cmd/
│   └── server/
│       └── main.go           # エントリーポイント
├── internal/
│   ├── config/               # 設定管理
│   │   └── config.go
│   ├── domain/               # ドメイン層（エンティティ、リポジトリインターフェース）
│   │   ├── user.go
│   │   └── errors.go
│   ├── usecase/              # ユースケース層（ビジネスロジック）
│   │   ├── user_usecase.go
│   │   └── user_usecase_test.go
│   ├── interface/            # インターフェース層（HTTP ハンドラー、ルーター）
│   │   ├── handler/
│   │   │   ├── user_handler.go
│   │   │   └── dto/
│   │   │       └── user_dto.go
│   │   └── router/
│   │       └── router.go
│   └── infrastructure/       # インフラ層（DB接続、リポジトリ実装）
│       └── persistence/
│           ├── db.go
│           └── user_repository.go
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## タクシー配車アプリの最小MVPロジック
本システムが目指す最小限のMVP（Minimum Viable Product）ロジックは以下の通りです。

1.  **ユーザー管理**: 乗客（Passenger）とドライバー（Driver）の情報を保持。
2.  **車両・ドライバーの状態管理**: ドライバーの現在地、空車/実車状態の管理。
3.  **配車リクエスト**: 乗客が現在地と目的地を指定してタクシーを呼ぶ機能。
4.  **マッチング**: 近くの空車タクシーに依頼を送り、ドライバーが承諾。
5.  **走行管理**: 乗車から降車までのステータス（Requested -> Accepted -> Ongoing -> Completed）の更新。

## セットアップ手順

### 1. データベースの起動
Docker Composeを使用してMySQLを起動します。

```bash
docker-compose up -d
```

### 2. 環境変数の設定
`.env` ファイルが作成されていることを確認してください。

```bash
DB_USER=user_name
DB_PASSWORD=user_password
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=my_database
SERVER_PORT=8080
```

### 3. アプリケーションの実行
```bash
go run cmd/server/main.go
```
サーバーはデフォルトで `http://localhost:8080` で起動します。

### 4. テストの実行
```bash
go test ./... -v
```

---

## APIエンドポイント

### バージョン付きエンドポイント（推奨）

| メソッド | エンドポイント | 説明 |
|:---|:---|:---|
| `GET` | `/api/v1/users` | 全ユーザー情報を取得 |
| `GET` | `/api/v1/users/:id` | 特定のユーザー情報を取得 |
| `POST` | `/api/v1/users` | 新規ユーザーを登録 |
| `GET` | `/health` | ヘルスチェック |

### レスポンス例

#### ユーザー取得成功
```json
{
  "id": "01JXXXXXXXXXXXXXXXXXXXXXX",
  "name": "田中太郎",
  "email": "tanaka@example.com"
}
```

#### バリデーションエラー
```json
{
  "error": "validation_error",
  "message": "Name cannot be empty"
}
```

#### 重複エラー
```json
{
  "error": "conflict",
  "message": "User with this email already exists"
}
```

---

## データの作成 (API)

本プロジェクトでは ID に **ULID** を採用しており、ユーザー作成時に自動生成されます。

### ユーザーの登録
ターミナルから以下の `curl` コマンドで新しいユーザーを作成できます。

```bash
curl -X POST http://localhost:8080/api/v1/users \
     -H "Content-Type: application/json" \
     -d '{"name":"鈴木一郎", "email":"suzuki@example.com"}'
```

成功すると、生成された ULID を含むユーザー情報が返ってきます。

### ユーザー一覧の取得
```bash
curl http://localhost:8080/api/v1/users
```

### 特定ユーザーの取得
```bash
curl http://localhost:8080/api/v1/users/{id}
```

---

## DBの中身を確認する

コンテナ内で実行されているMySQLに接続して、データを確認する方法です。

### 1. MySQLにログインしてインタラクティブに操作する
```bash
docker exec -it mysql-container mysql -u user_name -puser_password my_database --default-character-set=utf8mb4
```
ログイン後、SQLを実行できます（例: `SELECT * FROM users;`）。終了時は `exit` と入力します。

### 2. コマンド一発で中身を表示する
```bash
docker exec mysql-container mysql -u user_name -puser_password my_database --default-character-set=utf8mb4 -e "SELECT * FROM users;"
```

---

## アーキテクチャの特徴

### クリーンアーキテクチャの採用
- **依存関係の方向**: 外側の層（Handler, Repository）から内側の層（Domain, Usecase）への一方向のみ
- **テスタビリティ**: インターフェースを活用し、モックによる単体テストが容易
- **拡張性**: 新しいエンティティ（Driver, Ride など）を追加しても既存コードへの影響が最小

### カスタムエラー型
ドメイン固有のエラーを定義し、適切な HTTP ステータスコードを返却：
- `ErrUserNotFound` → 404 Not Found
- `ErrUserAlreadyExists` → 409 Conflict
- `ErrInvalidInput` → 400 Bad Request

### DTO（Data Transfer Object）
- リクエスト/レスポンス用の構造体をドメインモデルから分離
- API 仕様の変更がドメインモデルに影響しない設計
