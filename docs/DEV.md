# development guide

* 想定するgo runtimeのバージョン等は `${PROJECT_ROOT}/go.mod` を参照してください。

## 動作確認方法

### local

1. webサーバの起動

* `cmd/main.go` がエントリーポイントです。

```
$ cd ${PROJECT_ROOT}
$ APP_ENV=MOCK go run cmd/main.go
```

* [NOTE] 環境変数APP_ENVに `MOCK` が設定されている場合、DBをモック化してオンメモリに情報を更新/参照します。

2. ブラウザアクセス

* `http://localhost:5000/` にアクセス。
    * [NOTE] AWS Elastic Beanstalkのデフォルトに合わせて5000番ポートで待ち受けています。

### deploy

1. パッケージング

```
$ ce ${PROJECT_ROOT}
$ zip app.zip -r *
```

2. デプロイ

* AWS ConsoleからElastic Beanstalkの対象環境を選択し、前手順で作成したファイルをアップロードします。
    * [NOTE] `${PROJECT_ROOT}/Buildfile` にビルド処理、`${PROJECT_ROOT}/Procfile` に起動処理が記述されています。
