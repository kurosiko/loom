# Secure P2P Backend (Go)

安全を最優先にした Go 言語のバックエンドサーバー。ビデオ通話などの P2P 通信のためのシグナリングサーバーを実装しています。

## セキュリティ機能

### 暗号化
- **Argon2id**: パスワードハッシュ化に最も安全な KDF を使用
- **AES-256-GCM**: 認証付き暗号化でデータ保護
- **TLS 1.3**: 強制された最新の TLS バージョン
- **HMAC-SHA256**: メッセージ整合性検証

### セキュリティヘッダー
- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block
- Strict-Transport-Security (HSTS)
- Content-Security-Policy
- Referrer-Policy
- Permissions-Policy

## 実装済み機能

### 認証・認可
- ユーザー登録（パスワードは Argon2id でハッシュ化）
- ログイン（セキュアなセッショントークン発行）
- パスワード検証

### P2P シグナリング（WebRTC）
- WebSocket によるリアルタイム通信
- ルーム管理（作成、参加、退出）
- SDP Offer/Answer の交換
- ICE Candidate の交換
- ピア間メッセージング

### API エンドポイント
```
POST /api/register   - ユーザー登録
POST /api/login      - ログイン
POST /api/room       - ルーム作成
GET  /api/rooms      - アクティブルーム一覧
WS   /ws             - WebSocket シグナリング
```

## ビルド方法

```bash
# 依存関係のインストール
go mod tidy

# ビルド
go build -o server ./cmd/server
```

## 実行方法

### 開発モード（TLS なし）
```bash
./server -addr :8080
```

### 本番モード（TLS あり）
```bash
./server -addr :443 -cert /path/to/cert.pem -key /path/to/key.pem
```

## WebSocket メッセージフォーマット

### ルーム参加
```json
{
  "type": "join_room",
  "payload": {
    "room_id": "room-uuid",
    "user_id": "user-uuid"
  }
}
```

### WebRTC シグナリング
```json
{
  "type": "offer",
  "payload": {
    "type": "offer",
    "room_id": "room-uuid",
    "sender_id": "peer-uuid",
    "target_id": "peer-uuid",
    "sdp": "v=0..."
  }
}
```

## アーキテクチャ

```
secure-p2p-backend/
├── cmd/
│   └── server/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── crypto/
│   │   └── crypto.go        # 暗号化処理（Argon2, AES-GCM）
│   ├── handlers/
│   │   └── server.go        # HTTP/WebSocket ハンドラー
│   ├── models/
│   │   └── models.go        # データモデル
│   └── signaling/
│       └── room_manager.go  # ルーム管理
├── go.mod
└── go.sum
```

## セキュリティベストプラクティス

1. **パスワード保存**: Argon2id でソルト付きハッシュ化
2. **通信経路**: TLS 1.3 で暗号化
3. **データ保護**: AES-256-GCM で認証付き暗号化
4. **セッション管理**: クリプトグラフィックに安全なトークン
5. **入力検証**: すべての入力をサニタイズ
6. **セキュリティヘッダー**: XSS、クリックジャッキング対策
7. **エラーハンドリング**: 情報漏洩を防ぐための適切なエラーメッセージ

## 注意事項

- 本番環境では必ず TLS 証明書を使用してください
- CORS ポリシーは適切に設定してください（現在は開発用に全て許可）
- レート制限や DDOS 対策を追加してください
- ログには機密情報を含まないようにしてください

## ライセンス

MIT
