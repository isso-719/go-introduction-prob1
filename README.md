# 設問1 Section 5に以下の変更を加えよ
- 家計簿データファイルをプログラム引数で指定できるようにする

## 回答
以下のコードを main 関数内に記述することで解決した。

```go
	// コマンドライン引数を取得
	args := flag.Args()

	// AccountBook の作成
	var ab *AccountBook

	// args が 0 であれば accountbook.txt, そうでなければ args[0] 家計簿ファイルとする。
	if len(args) == 0 {
		// AccountBookをNewAccountBookを使って作成
		ab = NewAccountBook("accountbook.txt")
	} else {
		// 	// AccountBookを家計簿データファイルを使って作成
		ab = NewAccountBook(args[0])
	}
```

## 実行結果

コマンド引数指定なしは `accountbook.txt` が読み込まれる。
```bash
% go run .           
[1]入力 [2]最新10件 [3]集計 [4]終了
>2
===========
[0001] コーヒー:120円
[0002] パン:100円
===========
[1]入力 [2]最新10件 [3]集計 [4]終了
>
```

コマンド引数指定ありは `go run .` の後で指定したファイル読み込まれる。
```bash
% go run . testAB.txt
[1]入力 [2]最新10件 [3]集計 [4]終了
>2
===========
[0001] Switch:32970円
[0002] PS5:80500円
[0003] Xbox Series S:32978円
===========
[1]入力 [2]最新10件 [3]集計 [4]終了
>
```

ファイルが存在しなければエラー
```bash
% go run . hoge.txt  
エラー： open hoge.txt: no such file or directory
exit status 1
```

# 設問2 Section 5のプログラムの改善点を考察せよ
特に以下の観点を考慮すること。

- 読みやすいのか？
- ソフトウェアとしての品質は？

改善するには具体的にどのようにすればよいか示せ。

## 静的解析ツールの実行
- まず、以下の 2 つの静的解析ツールを実行する。
  - `go fmt ./...` 
  - 何も表示されなかった。
```bash
% go fmt ./...
```

  - `golangci-lint run ./...` (golangci-lint がインストールされていることが前提)
  - エラーを補足した。
```bash
% golangci-lint run ./...

accountbook.go:159:10: Error return value of `cw.Write` is not checked (errcheck)
        cw.Write(header)
```

ここで、`accountbook.go` の 159 行目で `cw.Write` が返す error チェックがされていない。

これではソフトウェアの信頼性が損なわれる。

よって、以下のように書き換えるべきである。

```go
// Before
cw.Write(header)

// After
if err := cw.Write(header); err != nil {
		fmt.Fprintln(os.Stderr, "エラー：", err)
		os.Exit(1)
	}
```

## main には最低限置く
main パッケージはあくまで実行するためのものだけ置く方が良いと思う。
だから、最低限のコードのみ置き、他は別パッケージに置く。

## main のパッケージは 1 つ
main パッケージはあくまで起動するためのものであるため、
1 つだけ置くのが理想的である。

## model と repository は分ける
model と repository は分けることで、わかりやすくなる。
model は、オブジェクトを表すものだけ(struct)を置く。
reporsitory は、メソッドの定義(interface)だけを置く。
これら 2 つは domain の責務を果たす。

## infra で実行する
repository で定義したメソッドは、infra で実行するようにする。

## usecase の作成
domain のメソッドを呼び出すときは、usecase を使う。

## 実際に実行するのは handler
main で実行しているようなユーザと直接繋がるようなところは handler パッケージで行うのが良いと思う。

## 結局
結局以下のような構成になると理想的である。

今回はここまで分けなくても良いと思うが、ある程度はファイルや package を分けなければならないと思う。
```
step05
├── cmd
│   └── main.go
├── db
│   └── accountbook.txt
└── pkg
    ├── domain
    │   ├── model
    │   │   ├── accountbook.go
    │   │   ├── item.go
    │   │   └── summary.go
    │   └── repository
    │       ├── accountbook.go
    │       ├── item.go
    │       └── summary.go
    ├── handle
    │   └── handle.go
    ├── infra
    │   ├── accountbook.go
    │   ├── item.go
    │   └── summary.go
    └── usecase
        ├── accountbook.go
        ├── item.go
        └── summary.go
```

上記のようなモデルをクリーンアーキテクチャという。

私は MVC なども良いとは思うが、Go はクリーンアーキテクチャで開発されているケースが多いと感じる部分がある。

## クリーンアーキテクチャの利点
クリーンアーキテクチャは、アプリケーションを構築する際に、アプリケーションの複雑さを抑えることができる思想である。

アプリケーションが複雑であると、テストが難しくなる。

しかし、複雑さを克服することができれば、テストが容易になり、
テストを行う試行回数も増えるため、ソフトウェアとしての信頼性が上がる。

つまり、クリーンアーキテクチャ的な思想でフォルダ設計すると必然的にコードが読みやすくなり、ソフトウェアの品質も上がる。
