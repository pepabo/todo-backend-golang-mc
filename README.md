# Todo-Backend API server written in Go for "LOLIPOP! Managed Cloud"

[Todo-Backend](https://todobackend.com/)のTODOアプリのスペックに対応したバックエンドAPIサーバです。

ロリポップ！マネージドクラウドのGoコンテナに2コマンドでデプロイ可能です。

## ビルド方法

``` console
$ make build
```

## データベースの初期化方法

``` console
$ SSH_USER=test-test-123 SSH_PORT=33322 DB_NAME=dbname DB_USER=user DB_PASS=pass make initdb
```

## デプロイ方法

![img](mc.png)

``` console
$ SSH_USER=test-test-123 SSH_PORT=33322 DB_NAME=dbname DB_USER=user DB_PASS=pass make deploy
```
