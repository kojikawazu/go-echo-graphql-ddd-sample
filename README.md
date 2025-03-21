# Go/Echo/GraphQLのDDD設計用サンプルコード

## Summary

- Go/Echoを使用する。
- DDD設計を行う。
- 通信方式はGraphQLを使用する。

## Directory

以下ディレクトリ構成とする。

```bash
/project
  ├── cmd/                # エントリーポイント
  ├── config/             # 設定ファイル
  ├── internal/
  │   ├── domain/         # ドメイン層（エンティティ、リポジトリ、VO）
  │   ├── usecase/        # ユースケース層（アプリケーションサービス）
  │   ├── infrastructure/ # インフラ層（DB, API クライアント, リポジトリ実装）
  │   ├── interfaces/     # インターフェース層（HTTPハンドラ, gRPC, CLI）
  │   ├── middleware/     # Echoのミドルウェア
  │   ├── pkg/            # 汎用ユーティリティ
  │   ├── router/         # ルーティング
  │   ├── test/           # ユニット・統合テスト
  ├── Dockerfile
  ├── go.mod
  ├── go.sum
  └── README.md
```
