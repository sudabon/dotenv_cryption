# cli-commands Specification

## Purpose
TBD - created by archiving change envcrypt-cli. Update Purpose after archive.
## Requirements
### Requirement: encrypt コマンド
CLI SHALL provide an `envcrypt encrypt --file <path>` command that encrypts the specified `.env` file and generates an output file. If the `--file` flag is omitted, the current directory's `.env` SHALL be used by default.

#### Scenario: ファイル指定で暗号化
- **WHEN** `envcrypt encrypt --file .env` を実行する
- **THEN** 暗号化されたファイル `.env.enc` が生成される

#### Scenario: デフォルトファイルで暗号化
- **WHEN** `envcrypt encrypt` をファイル指定なしで実行する
- **THEN** カレントディレクトリの `.env` が暗号化され `.env.enc` が生成される

#### Scenario: 指定ファイルが存在しない
- **WHEN** 存在しないファイルを `--file` で指定して `envcrypt encrypt` を実行する
- **THEN** エラーメッセージを表示して終了コード1で終了する

### Requirement: decrypt コマンド
CLI SHALL provide an `envcrypt decrypt --file <path>` command that decrypts an encrypted file and restores `.env`. If the `--file` flag is omitted, the current directory's `.env.enc` SHALL be used by default.

#### Scenario: ファイル指定で復号
- **WHEN** `envcrypt decrypt --file .env.enc` を実行する
- **THEN** 復号された `.env` ファイルが生成される

#### Scenario: デフォルトファイルで復号
- **WHEN** `envcrypt decrypt` をファイル指定なしで実行する
- **THEN** カレントディレクトリの `.env.enc` が復号され `.env` が生成される

#### Scenario: 不正なフォーマットのファイル
- **WHEN** ENVC マジックバイトを持たないファイルを指定して `envcrypt decrypt` を実行する
- **THEN** エラーメッセージ「invalid file format」を表示して終了コード1で終了する

### Requirement: 出力ファイル名のカスタマイズ
The encrypted file prefix SHALL be configurable via `dotenv.yaml` `files.encrypted_prefix`.

#### Scenario: カスタムプレフィックスでの暗号化
- **WHEN** `dotenv.yaml` に `encrypted_prefix: enc.` が設定されている状態で `envcrypt encrypt --file .env` を実行する
- **THEN** 出力ファイル名は `enc..env` となる

#### Scenario: デフォルトプレフィックスでの暗号化
- **WHEN** `dotenv.yaml` に `encrypted_prefix` が未設定の状態で暗号化を実行する
- **THEN** 出力ファイル名はデフォルトの `.env.enc` サフィックス形式となる

### Requirement: create master コマンド
CLI SHALL provide an `envcrypt create master` command that generates a new 32-byte master key and stores it in the Secret Manager configured by `dotenv.yaml`.

#### Scenario: GCP でマスター鍵を作成
- **WHEN** `cloud: gcp` が設定された `dotenv.yaml` で `envcrypt create master` を実行する
- **THEN** `gcp.project_id` と `gcp.secret_id` のシークレットに新しい 32 バイト鍵が登録される

#### Scenario: AWS でマスター鍵を作成
- **WHEN** `cloud: aws` が設定された `dotenv.yaml` で `envcrypt create master` を実行する
- **THEN** `aws.region` と `aws.secret_id` のシークレットに新しい 32 バイト鍵が登録される

#### Scenario: 既存シークレットがある
- **WHEN** 同名のシークレットが既に存在する状態で `envcrypt create master` を実行する
- **THEN** シークレットを上書きせず、既に存在することを示すエラーで終了する

### Requirement: delete master コマンド
CLI SHALL provide an `envcrypt delete master` command that deletes the master key secret configured by `dotenv.yaml`.

#### Scenario: 設定済みシークレットを削除
- **WHEN** `dotenv.yaml` に設定されたシークレットが存在する状態で `envcrypt delete master` を実行する
- **THEN** 対応する Secret Manager からそのシークレットが削除される

#### Scenario: シークレットが存在しない
- **WHEN** 設定されたシークレットが存在しない状態で `envcrypt delete master` を実行する
- **THEN** シークレット名を含む not found エラーで終了する

