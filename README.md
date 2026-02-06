# gin-fleamarket

Go（Gin フレームワーク）で構築されたフリーマーケット風の商品管理 REST API です。

## システム概要

商品（Item）の CRUD 操作を提供するシンプルな Web API サーバーです。  
現在はインメモリ（メモリ上のスライス）でデータを管理しており、サーバー再起動でデータはリセットされます。

### 技術スタック

| 項目 | 内容 |
|------|------|
| 言語 | Go 1.25 |
| Web フレームワーク | [Gin](https://github.com/gin-gonic/gin) v1.11.0 |
| ホットリロード | [Air](https://github.com/air-verse/air)（`.air.toml` で設定済み） |
| データストア | インメモリ（スライス） |

### アーキテクチャ

レイヤードアーキテクチャを採用しており、責務ごとにパッケージが分離されています。

```
gin-fleamarket/
├── main.go                # エントリーポイント・ルーティング定義
├── controllers/           # HTTP リクエスト/レスポンスの処理
│   └── item_controller.go
├── services/              # ビジネスロジック
│   └── item_service.go
├── repositories/          # データアクセス層
│   └── item_repository.go
├── models/                # データモデル定義
│   └── item.go
├── dto/                   # リクエスト/レスポンス用の DTO（Data Transfer Object）
│   └── item_dto.go
├── .air.toml              # Air（ホットリロード）設定
├── go.mod
└── go.sum
```

### データモデル（Item）

| フィールド | 型 | 説明 |
|-----------|------|------|
| ID | uint | 商品ID（自動採番） |
| Name | string | 商品名 |
| Price | uint | 価格 |
| Description | string | 説明 |
| SoldOut | bool | 売り切れフラグ |

## API エンドポイント

| メソッド | パス | 説明 |
|---------|------|------|
| GET | `/items` | 全商品一覧を取得 |
| GET | `/items/:id` | 指定IDの商品を取得 |
| POST | `/items` | 新規商品を作成 |
| PUT | `/items/:id` | 指定IDの商品を更新 |
| DELETE | `/items/:id` | 指定IDの商品を削除 |

### リクエスト/レスポンス例

#### 全商品取得

```bash
curl http://localhost:8080/items
```

```json
{
  "data": [
    { "ID": 1, "Name": "商品1", "Price": 1000, "Description": "説明1", "SoldOut": false },
    { "ID": 2, "Name": "商品2", "Price": 2000, "Description": "説明2", "SoldOut": true },
    { "ID": 3, "Name": "商品3", "Price": 3000, "Description": "説明3", "SoldOut": false }
  ]
}
```

#### 商品作成

```bash
curl -X POST http://localhost:8080/items \
  -H "Content-Type: application/json" \ 
  -d '{"name": "新商品", "price": 1500, "description": "新しい商品です"}'
```

```json
{
  "data": { "ID": 4, "Name": "新商品", "Price": 1500, "Description": "新しい商品です", "SoldOut": false }
}
```

#### 商品更新（部分更新対応）

```bash
curl -X PUT http://localhost:8080/items/1 \
  -H "Content-Type: application/json" \ 
  -d '{"name": "更新商品", "price": 2500}'
```

#### 商品削除

```bash
curl -X DELETE http://localhost:8080/items/1
```

### バリデーション

- **name**: 必須、2文字以上（作成時）
- **price**: 必須、1〜999999（作成時）
- **description**: 任意
- **soldOut**: 任意（更新時のみ）

更新時はポインタ型を使用しており、送信されたフィールドのみが更新されます（部分更新）。

## 実行方法

### 前提条件

- Go 1.25 以上がインストールされていること

### 通常起動

```bash
# 依存関係のダウンロード
go mod download

# サーバー起動
go run main.go
```

サーバーは `localhost:8080` で起動します。

### ホットリロード（開発時）

[Air](https://github.com/air-verse/air) を使用してホットリロードで開発できます。

```bash
# Air のインストール
go install github.com/air-verse/air@latest

# ホットリロードで起動
air
```

`.air.toml` に設定が定義されており、ソースコード変更時に自動でビルド＆再起動されます。

## 初期データ

起動時に以下の3つのサンプル商品がメモリ上にロードされます。

| ID | Name | Price | Description | SoldOut |
|----|------|-------|-------------|---------|
| 1 | 商品1 | 1000 | 説明1 | false |
| 2 | 商品2 | 2000 | 説明2 | true |
| 3 | 商品3 | 3000 | 説明3 | false |