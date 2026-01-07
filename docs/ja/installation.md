# インストール

以下のいずれかの方法で Kex をインストールできます。

## バイナリのダウンロード

[GitHub Releases](https://github.com/mew-ton/kex/releases) ページから、最新のコンパイル済みバイナリをダウンロードできます。

1. OS とアーキテクチャに合ったアーカイブをダウンロードしてください (例: `_linux_amd64.tar.gz`, `_darwin_arm64.tar.gz`)。
2. アーカイブを解凍します。
3. `kex` バイナリを `PATH` の通ったディレクトリ (例: `/usr/local/bin`) に移動します。

```bash
# Linux amd64 の例
tar -xzf kex_linux_amd64.tar.gz
sudo mv kex /usr/local/bin/
```

## Go Install

Go (1.25以上) がインストールされている場合、以下のコマンドで直接インストールできます:

```bash
go install github.com/mew-ton/kex/cmd/kex@latest
```

`$(go env GOPATH)/bin` が `PATH` に含まれていることを確認してください。
