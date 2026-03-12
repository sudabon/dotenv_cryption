## Why

`.env` ファイルには機密情報が含まれるが、チーム間での共有やGitでの管理が難しい。暗号化してリポジトリにコミットし、必要時に復号できるCLIツールがあれば、セキュアかつ簡便な運用が可能になる。マスターキーをクラウド(GCP/AWS)のSecret Managerに保管することで、ローカルに鍵を置かない安全な鍵管理を実現する。

## What Changes

- Go製CLIツール `envcrypt` を新規作成
- `encrypt` コマンド: `.env` ファイルをAES-256-GCMで暗号化し、暗号化ファイルを出力
- `decrypt` コマンド: 暗号化ファイルを復号し、`.env` を復元
- Google Cloud Secret Manager / AWS Secret Manager からマスターキーを取得
- ファイルごとにランダムなデータキーを生成し、マスターキーでラップする鍵階層
- 設定ファイル `dotenv.yaml` による構成管理
- 独自バイナリフォーマット（ENVC マジックバイト付き）で暗号化ファイルを保存

## Capabilities

### New Capabilities
- `cli-commands`: cobra を使った encrypt / decrypt コマンドの実装
- `config-management`: viper による `dotenv.yaml` の読み込みと設定管理
- `encryption-engine`: AES-256-GCM による暗号化・復号およびデータキーのラップ/アンラップ
- `secret-provider`: GCP / AWS Secret Manager からのマスターキー取得
- `file-format`: ENVC バイナリフォーマットの読み書き

### Modified Capabilities
(なし - 新規プロジェクト)

## Impact

- 新規Goモジュールの作成（依存: cobra, viper, cloud.google.com/go/secretmanager, aws-sdk-go-v2）
- GCP IAM ロール `roles/secretmanager.secretAccessor` が必要
- AWS IAM ポリシー `secretsmanager:GetSecretValue` が必要
- `.gitignore` に `.env` を追加し、暗号化ファイルのみコミット対象とする運用
