# git-download

GitリポジトリをGit管理情報（`.git`）なしで取得するためのGitサブコマンド。

## 概要

`git download` は、指定したGitリポジトリのファイルをローカルに取得する。

通常の `git clone` とは異なり、取得先はGitリポジトリとして管理しない。

```text
git clone
└── repository/
    ├── .git/
    ├── src/
    └── README.md

git download
└── repository/
    ├── src/
    └── README.md
```

主な用途は、Gitリポジトリのソースコードや設定ファイルなどを、Git管理を必要とせず取得することである。

## コマンド

```text
git download <repository> [directory]
```

### 引数

#### `<repository>`

取得元のGitリポジトリ。

URLまたはGitがサポートするリポジトリ指定形式を受け付ける。

例：

```text
https://github.com/example/example.git
git@github.com:example/example.git
```

#### `[directory]`

取得先ディレクトリ。

省略した場合は、リポジトリ名からディレクトリ名を決定する。

```powershell
git download https://github.com/example/example.git
```

の場合、

```text
./example/
```

に取得する。

明示的に指定することもできる。

```powershell
git download https://github.com/example/example.git my-project
```

この場合、

```text
./my-project/
```

に取得する。

## 基本動作

```text
1. リポジトリの存在・アクセス可否を確認
2. リポジトリの内容を取得
3. 指定されたブランチまたはデフォルトブランチのファイルを展開
4. `.git` ディレクトリを作成しない
5. 取得したファイルを指定ディレクトリに配置
```

最終的な取得先は通常のディレクトリとして扱えるものとする。

```text
project/
├── src/
├── public/
├── package.json
└── README.md
```

## Git管理情報

`git download` によって取得されたディレクトリには、Git管理情報を残さない。

特に以下を取得先に残してはならない。

```text
.git/
```

また、取得後に `git status` や `git log` などのGit操作を行って、取得元リポジトリの履歴を参照できる状態にする必要はない。

## ブランチ

デフォルトでは、リポジトリのデフォルトブランチを取得する。

将来的にブランチ指定をサポートする。

例：

```powershell
git download https://github.com/example/example.git --branch develop
```

## タグ・コミット指定

将来的に特定のタグやコミットを指定できるようにする。

例：

```powershell
git download https://github.com/example/example.git --tag v1.0.0
```

```powershell
git download https://github.com/example/example.git --commit abcdef123456
```

ただし、初期実装では必須としない。

## 既存ディレクトリ

取得先ディレクトリが既に存在する場合の動作を明確にする。

初期仕様では、安全のため以下とする。

* 空ではないディレクトリへの取得を拒否する
* 既存ファイルを上書きしない
* `--force` などによる強制上書きは将来的な機能とする

例：

```text
Error: destination directory 'my-project' already exists and is not empty.
```

## エラー

取得に失敗した場合は、エラー内容を標準エラー出力に表示し、0以外の終了コードを返す。

主なエラー：

* リポジトリが存在しない
* リポジトリへのアクセス権がない
* 認証に失敗した
* ネットワークエラー
* 取得先ディレクトリを作成できない
* 取得先ディレクトリが既に存在する
* リポジトリの取得に失敗した

## 終了コード

```text
0    成功
1    一般的なエラー
```

将来的にエラー種別ごとの終了コードを追加してもよい。

## CLI例

### GitHubから取得

```powershell
git download https://github.com/example/example.git
```

結果：

```text
Downloading https://github.com/example/example.git...
Downloaded to ./example/
```

### ディレクトリを指定

```powershell
git download https://github.com/example/example.git ./my-project
```

### SSH URL

```powershell
git download git@github.com:example/example.git
```

## 設計方針

### `git` サブコマンドとして提供する

実行形式は以下とする。

```text
git download
```

Gitは `git-<command>` 形式の外部サブコマンドをサポートしているため、実装は以下のような名前で提供する。

```text
git-download
```

PATH上に `git-download` を配置することで、

```powershell
git download ...
```

として実行できるようにする。

### `git clone` との差別化

|              | `git clone` | `git download` |
| ------------ | ----------- | -------------- |
| ファイル取得       | ○           | ○              |
| `.git`       | あり          | なし             |
| Git管理        | あり          | なし             |
| `git log`    | ○           | ×              |
| `git status` | ○           | ×              |
| 主用途          | 開発          | ファイル取得         |

## 初期実装

初期実装では、実装の単純化を優先して以下の方式を許容する。

```text
git clone
    ↓
取得先を確認
    ↓
.git を削除
```

ただし、これは内部実装上の都合であり、`git download` の仕様として `.git` の作成を要求するものではない。

将来的には、Gitの機能を利用して`.git`を作成せずにファイルのみ取得する方式へ変更してもよい。

## 非目標

初期バージョンでは以下を対象外とする。

* Gitリポジトリとしての管理
* コミット履歴の取得
* リモートとの同期
* `git pull` 相当の更新
* 複数リポジトリの一括取得
* 特定コミット・タグ・ブランチの指定
* ダウンロード済みファイルとの差分更新

これらは必要に応じて将来追加する。

## 目標

最終的には、以下の操作だけでGitリポジトリの内容を通常のディレクトリとして取得できることを目標とする。

```powershell
git download <repository>
```

利用者はGitの管理機能を意識することなく、リポジトリに含まれるファイルだけを取得できるものとする。
