# gin-fleamarket

Go (Gin) + PostgreSQL で構築したフリマ風 REST API です。
商品（Item）の CRUD 操作を提供します。

## 技術スタック

| カテゴリ           | 技術                                    |
| ------------------ | --------------------------------------- |
| 言語               | Go 1.25                                 |
| Web フレームワーク | [Gin](https://github.com/gin-gonic/gin) |
| ORM                | [GORM](https://gorm.io)                 |
| データベース       | PostgreSQL 16                           |
| ホットリロード     | [Air](https://github.com/air-verse/air) |
| コンテナ           | Docker Compose                          |

## ディレクトリ構成

```
.
├── main.go                  # エントリポイント（ルーティング・DI）
├── controllers/
│   └── item_controller.go   # HTTPリクエスト/レスポンス処理
├── services/
│   └── item_service.go      # ビジネスロジック
├── repositories/
│   └── item_repository.go   # データアクセス層（GORM / メモリ）
├── models/
│   └── item.go              # データモデル定義
├── dto/
│   └── item_dto.go          # リクエストバリデーション用DTO
├── infra/
│   ├── db.go                # DB接続設定
│   └── initializer.go       # 環境変数読み込み
├── migrations/
│   └── migration.go         # DBマイグレーション
├── docker-compose.yaml      # PostgreSQL / pgAdmin 定義
└── .air.toml                # Air（ホットリロード）設定
```

## セットアップ

### 前提条件

- Go 1.25+
- Docker / Docker Compose
- [Air](https://github.com/air-verse/air)（ホットリロードを使う場合）

### 1. リポジトリをクローン

```bash
git clone https://github.com/yone-al/gin-fleamarket.git
cd gin-fleamarket
```

### 2. 環境変数を設定

`.env` ファイルをプロジェクトルートに作成します。

```env
ENV=prod
DB_HOST=localhost
DB_USER=ginuser
DB_PASSWORD=ginpassword
DB_NAME=fleamarket
DB_PORT=5432
```

### 3. PostgreSQL を起動

```bash
docker compose up -d
```

PostgreSQL が `localhost:5432` で起動します。
pgAdmin には `http://localhost:81` でアクセスできます（Email: `gin@example.com` / Password: `ginpassword`）。

### 4. マイグレーション実行

```bash
go run migrations/migration.go
```

### 5. アプリケーション起動

```bash
# 通常起動
go run main.go

# ホットリロード（Air）
air
```

サーバーが `http://localhost:8080` で起動します。

## API エンドポイント

### 商品一覧取得

```
GET /items
```

**レスポンス例:**

```json
{
  "data": [
    {
      "ID": 1,
      "Name": "商品1",
      "Price": 1000,
      "Description": "説明1",
      "SoldOut": false,
      "CreatedAt": "2026-02-10T12:00:00Z",
      "UpdatedAt": "2026-02-10T12:00:00Z",
      "DeletedAt": null
    }
  ]
}
```

### 商品取得（ID指定）

```
GET /items/:id
```

### 商品作成

```
POST /items
Content-Type: application/json
```

**リクエストボディ:**

```json
{
  "name": "商品名",
  "price": 1500,
  "description": "商品の説明"
}
```

| フィールド  | 型     | 必須 | バリデーション |
| ----------- | ------ | ---- | -------------- |
| name        | string | Yes  | 2文字以上      |
| price       | uint   | Yes  | 1〜999999      |
| description | string | No   | -              |

### 商品更新

```
PUT /items/:id
Content-Type: application/json
```

**リクエストボディ（部分更新可）:**

```json
{
  "name": "新しい商品名",
  "price": 2000,
  "soldOut": true
}
```

| フィールド  | 型       | 必須 | バリデーション    |
| ----------- | -------- | ---- | ----------------- |
| name        | \*string | No   | 指定時は2文字以上 |
| price       | \*uint   | No   | 指定時は1〜999999 |
| description | \*string | No   | -                 |
| soldOut     | \*bool   | No   | -                 |

### 商品削除

```
DELETE /items/:id
```

## アーキテクチャ

3層アーキテクチャ（Presentation / Business Logic / Data Access）を採用しています。
各層はインターフェースで疎結合になっており、テストやリポジトリの差し替えが容易です。

```
[Presentation Layer]     controllers/   … HTTPリクエスト/レスポンス処理
        ↓
[Business Logic Layer]   services/      … ビジネスロジック
        ↓
[Data Access Layer]      repositories/  … データアクセスの抽象化
        ↓
        DB (PostgreSQL / GORM)
```

| 層             | パッケージ     | 責務                                                           |
| -------------- | -------------- | -------------------------------------------------------------- |
| Presentation   | `controllers`  | HTTPリクエストのパース・バリデーション・レスポンス返却         |
| Business Logic | `services`     | ビジネスロジック（部分更新処理など）                           |
| Data Access    | `repositories` | `IItemRepository` インターフェースによるデータアクセスの抽象化 |

Repository 層は実装を差し替え可能です:

- `ItemRepository` — PostgreSQL（GORM）を使った本番用実装
- `ItemMemoryRepository` — メモリ上のデータを使った開発/テスト用実装
