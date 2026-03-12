# dotenv_cryption

環境変数ファイルの暗号化・復号ツール

`.env` ファイルを AES-256-GCM で暗号化し、Git リポジトリで安全に管理できる Go 製 CLI ツールです。マスターキーは Google Cloud Secret Manager または AWS Secrets Manager から取得し、ローカルに鍵を保持しません。

## 特徴

- AES-256-GCM による認証付き暗号化（改ざん検知付き）
- ファイルごとにランダムなデータキーを生成するエンベロープ暗号化
- GCP Secret Manager / AWS Secrets Manager 対応
- `dotenv.yaml` による宣言的な設定管理
- 独自バイナリフォーマット（ENVC ヘッダー）によるバージョン管理対応

## インストール

GitHub Releases から利用している環境に対応した tarball をダウンロードしてください。

- Apple Silicon Mac: `envcrypt_<version>_darwin_arm64.tar.gz`
- Intel Mac: `envcrypt_<version>_darwin_amd64.tar.gz`
- Linux x86_64: `envcrypt_<version>_linux_amd64.tar.gz`

Releases: `https://github.com/sudabon/dotenv_cryption/releases`

```bash
VERSION=v0.1.0
OS=darwin
ARCH=arm64

curl -LO "https://github.com/sudabon/dotenv_cryption/releases/download/${VERSION}/envcrypt_${VERSION}_${OS}_${ARCH}.tar.gz"
tar -xzf "envcrypt_${VERSION}_${OS}_${ARCH}.tar.gz"
install -m 0755 envcrypt /usr/local/bin/envcrypt
envcrypt version
```

## セットアップ

### 1. 設定ファイルを作成

プロジェクトルートに `dotenv.yaml` を作成します。

**GCP の場合:**

```yaml
cloud: gcp

gcp:
  project_id: my-project
  secret_id: envcrypt-master-key

crypto:
  algorithm: aes-256-gcm

files:
  encrypted_prefix: ""
```

**AWS の場合:**

```yaml
cloud: aws

aws:
  region: ap-northeast-1
  secret_id: envcrypt-master-key

crypto:
  algorithm: aes-256-gcm

files:
  encrypted_prefix: ""
```

### 2. クラウド認証を設定

**GCP の場合:**

```bash
gcloud auth application-default login
```

必要な IAM ロール: `roles/secretmanager.secretAccessor`

create/delete も使う場合は、Secret 作成・削除権限も付与してください。

**AWS の場合:**

```bash
aws configure
```

必要な IAM ポリシー: `secretsmanager:GetSecretValue`

create/delete も使う場合は、Secret 作成・削除権限も付与してください。

### 3. マスターキーを作成

```bash
envcrypt create master
```

`dotenv.yaml` に書かれた `secret_id` に対して、新しい 32 バイトのマスター鍵を登録します。対象シークレットが既に存在する場合は上書きせずにエラーで終了します。

### 4. .gitignore の設定

```gitignore
.env
```

## 使い方

### 暗号化

```bash
# デフォルト (.env → .env.enc)
envcrypt encrypt

# ファイル指定
envcrypt encrypt --file .env.production
```

### 復号

```bash
# デフォルト (.env.enc → .env)
envcrypt decrypt

# ファイル指定
envcrypt decrypt --file .env.production.enc
```

### マスターキー作成

```bash
envcrypt create master
```

### マスターキー削除

```bash
envcrypt delete master
```

### ワークフロー例

```bash
# 暗号化してコミット
envcrypt encrypt
git add .env.enc
git commit -m "Update encrypted env"

# 別の環境で復号
git pull
envcrypt decrypt
```

### マスターキーの再作成

```bash
envcrypt delete master
envcrypt create master
```

## 設定リファレンス

| フィールド | 説明 | 必須 |
|---|---|---|
| `cloud` | クラウドプロバイダ (`gcp` または `aws`) | Yes |
| `gcp.project_id` | GCP プロジェクト ID | cloud=gcp の場合 |
| `gcp.secret_id` | GCP Secret Manager のシークレット ID | cloud=gcp の場合 |
| `aws.region` | AWS リージョン | cloud=aws の場合 |
| `aws.secret_id` | AWS Secrets Manager のシークレット ID | cloud=aws の場合 |
| `crypto.algorithm` | 暗号化アルゴリズム (デフォルト: `aes-256-gcm`) | No |
| `files.encrypted_prefix` | 暗号化ファイルのプレフィックス | No |

### 出力ファイル名

- `encrypted_prefix` 未設定: `.env` → `.env.enc`（サフィックス形式）
- `encrypted_prefix: enc.`: `.env` → `enc..env`（プレフィックス形式）

## 暗号化ファイルフォーマット

```
MAGIC(4B)  VERSION(1B)  NONCE_LEN(1B)  WRAPPED_KEY_LEN(2B)  NONCE  WRAPPED_KEY  CIPHERTEXT
"ENVC"     0x01         12             variable              ...    ...          ...
```

## アーキテクチャ

```
.env
 ↓
Random Data Key (32B, ファイルごとに生成)
 ↓
AES-256-GCM Encrypt
 ↓
Data Key を Master Key でラップ
 ↓
ENVC バイナリフォーマットで保存 → .env.enc

Master Key は Cloud Secret Manager から取得（ローカル保存なし）
```

## 開発

```bash
# テスト実行
go test ./...

# ビルド
go build -o envcrypt .

# 静的解析
go vet ./...
```

## ライセンス

MIT
