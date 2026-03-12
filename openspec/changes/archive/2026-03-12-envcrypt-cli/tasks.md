## 1. プロジェクト初期化

- [x] 1.1 Go モジュールの初期化 (`go mod init`) と基本ディレクトリ構造の作成 (cmd/, internal/, pkg/)
- [x] 1.2 依存パッケージの追加 (cobra, viper, cloud.google.com/go/secretmanager, aws-sdk-go-v2)
- [x] 1.3 main.go のエントリーポイント作成と cobra ルートコマンドの設定

## 2. 設定管理 (config-management)

- [x] 2.1 `dotenv.yaml` の構造体定義と viper による読み込み実装
- [x] 2.2 cloud フィールドのバリデーション (gcp / aws のみ許可)
- [x] 2.3 プロバイダ別の必須フィールドバリデーション (GCP: project_id, secret_id / AWS: region, secret_id)
- [x] 2.4 設定管理のユニットテスト

## 3. 暗号化エンジン (encryption-engine)

- [x] 3.1 AES-256-GCM による暗号化関数の実装 (ノンス生成 + 暗号化)
- [x] 3.2 AES-256-GCM による復号関数の実装
- [x] 3.3 データキー生成関数の実装 (crypto/rand で32バイト)
- [x] 3.4 データキーのラップ/アンラップ関数の実装 (マスターキーによる AES-GCM 暗号化)
- [x] 3.5 暗号化エンジンのユニットテスト (正常系・異常系: 不正キー長、改ざん検知、不正マスターキー)

## 4. ファイルフォーマット (file-format)

- [x] 4.1 ENVC バイナリフォーマットの書き込み関数の実装 (MAGIC + VERSION + NONCE_LEN + WRAPPED_KEY_LEN + NONCE + WRAPPED_KEY + CIPHERTEXT)
- [x] 4.2 ENVC バイナリフォーマットの読み込み・解析関数の実装
- [x] 4.3 フォーマット検証の実装 (マジックバイト、バージョン、データサイズ)
- [x] 4.4 ファイルフォーマットのユニットテスト

## 5. Secret Provider (secret-provider)

- [x] 5.1 SecretProvider インターフェースの定義 (`GetMasterKey() ([]byte, error)`)
- [x] 5.2 GCP Secret Manager プロバイダの実装
- [x] 5.3 AWS Secrets Manager プロバイダの実装
- [x] 5.4 プロバイダファクトリの実装 (cloud フィールドに基づく切り替え)
- [x] 5.5 Secret Provider のユニットテスト (モック使用)

## 6. CLI コマンド (cli-commands)

- [x] 6.1 `encrypt` サブコマンドの実装 (--file フラグ、デフォルト .env)
- [x] 6.2 `decrypt` サブコマンドの実装 (--file フラグ、デフォルト .env.enc)
- [x] 6.3 出力ファイル名の生成ロジック (encrypted_prefix 対応)
- [x] 6.4 encrypt/decrypt のエンドツーエンドフローの結合 (設定読み込み → キー取得 → 暗号化/復号 → ファイル出力)
- [x] 6.5 CLI コマンドの統合テスト

## 7. 仕上げ

- [x] 7.1 .gitignore に `.env` を追加
- [x] 7.2 エラーメッセージの統一と認証ガイダンスの追加
- [x] 7.3 go vet / staticcheck による静的解析パス
