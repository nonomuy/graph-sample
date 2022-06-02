# Architecture guide

## システム構成

* 本番システム構成はAWS Elastic Beanstalk + AWS Neptuneを想定します。
* AWSに接続できない環境での開発を想定し、オンメモリに情報を保持するモックモードをサポートします。
    * refs.) [DEV.md](DEV.md)

## パッケージ構成

* レイヤードアーキテクチャを想定します。
* DomainとApplicationを剥がす程の複雑性をもたないため、Application層は割愛しています。
* `domain/repository` 層をinterfaceとして設けています。
    * 実装クラスは環境毎に差し替えることを想定します。

```
└ cmd
└ presentation
└ domain
  └ mopdel
  └ repository
└ infrastructure
  └ mockRepository
  └ neptureRepository
└ web
└ Buildfile
└ build.sh
└ Procfile
└ go.mod
```