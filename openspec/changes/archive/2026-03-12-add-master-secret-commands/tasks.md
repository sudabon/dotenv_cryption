## 1. OpenSpec と CLI 追加

- [x] 1.1 `cli-commands` / `secret-provider` の change artifact を作成する
- [x] 1.2 `create master` / `delete master` の Cobra コマンドを追加する

## 2. Provider 拡張

- [x] 2.1 SecretProvider インターフェースとテスト用スタブを create/delete 対応に更新する
- [x] 2.2 GCP provider にマスター鍵の作成・削除を実装する
- [x] 2.3 AWS provider にマスター鍵の作成・削除を実装する

## 3. 検証とドキュメント

- [x] 3.1 CLI/Provider テストを追加・更新する
- [x] 3.2 README のセットアップとコマンド使用例を更新する
- [x] 3.3 テストを実行して change の tasks を完了にする
