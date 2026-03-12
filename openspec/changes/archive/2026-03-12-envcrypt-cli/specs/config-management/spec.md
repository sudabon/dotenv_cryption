## ADDED Requirements

### Requirement: dotenv.yaml の読み込み
CLI SHALL load `dotenv.yaml` from the current directory at startup and apply its configuration.

#### Scenario: 設定ファイルが存在する
- **WHEN** カレントディレクトリに `dotenv.yaml` が存在する状態で CLI を実行する
- **THEN** ファイルの内容に基づいて cloud プロバイダ、暗号化アルゴリズム、ファイルプレフィックスが設定される

#### Scenario: 設定ファイルが存在しない
- **WHEN** カレントディレクトリに `dotenv.yaml` が存在しない状態で CLI を実行する
- **THEN** エラーメッセージ「dotenv.yaml not found」を表示して終了コード1で終了する

### Requirement: cloud フィールドによるプロバイダ選択
The `cloud` field in `dotenv.yaml` SHALL accept `gcp` or `aws` and determine which Secret Manager provider is used.

#### Scenario: GCP プロバイダの選択
- **WHEN** `dotenv.yaml` に `cloud: gcp` が設定されている
- **THEN** `gcp.project_id` と `gcp.secret_id` を使用して GCP Secret Manager に接続する

#### Scenario: AWS プロバイダの選択
- **WHEN** `dotenv.yaml` に `cloud: aws` が設定されている
- **THEN** `aws.region` と `aws.secret_id` を使用して AWS Secrets Manager に接続する

#### Scenario: 不正なプロバイダ指定
- **WHEN** `dotenv.yaml` に `cloud: azure` のようなサポート外の値が設定されている
- **THEN** エラーメッセージ「unsupported cloud provider」を表示して終了コード1で終了する

### Requirement: 設定フィールドのバリデーション
CLI SHALL validate the presence of required fields at startup.

#### Scenario: GCP の必須フィールドが不足
- **WHEN** `cloud: gcp` で `project_id` が未設定の状態で CLI を実行する
- **THEN** エラーメッセージで不足フィールドを明示して終了コード1で終了する

#### Scenario: AWS の必須フィールドが不足
- **WHEN** `cloud: aws` で `region` が未設定の状態で CLI を実行する
- **THEN** エラーメッセージで不足フィールドを明示して終了コード1で終了する
