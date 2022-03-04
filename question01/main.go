package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	flagOutputDir string
)

func init() {
	flag.StringVar(&flagOutputDir, "dir", ".", "output")
}

func main() {
	flag.Parse()

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

	for {
		// モードを選択して実行する
		var mode int
		fmt.Println("[1]入力 [2]最新10件 [3]集計 [4]終了")
		fmt.Printf(">")
		fmt.Scan(&mode)

		// モードによって処理を変える
		switch mode {
		case 1: // 入力
			var n int
			for {
				fmt.Print("何件入力しますか>")
				fmt.Scan(&n)
				if n > 0 {
					break
				}
				fmt.Fprintln(os.Stderr, "1以上を入力してください")
			}

			for i := 0; i < n; i++ {
				ab.AddItem(inputItem())
			}
		case 2: // 最新10件
			items := ab.GetItems(10)
			showItems(items)
		case 3: // 集計
			summaries := ab.GetSummaries()
			file := filepath.Join(flagOutputDir, "summary.csv")
			SaveSummary(file, summaries)
			fmt.Println(file, "を出力しました")
		case 4: // 終了
			fmt.Println("終了します")
			os.Exit(0)
		default:
			fmt.Fprintln(os.Stderr, "1〜4で入力してください")
		}

		if ab.NumItems() > 20 {
			fmt.Fprintln(os.Stderr, "これ以上家計簿を増やすことができません")
			fmt.Fprintln(os.Stderr, "終了します")
			os.Exit(1)
		}
	}
}

// Itemを入力し返す
func inputItem() *Item {
	var item Item

	fmt.Print("品目>")
	fmt.Scan(&item.Category)

	fmt.Print("値段>")
	fmt.Scan(&item.Price)

	return &item
}

// Itemの一覧を出力する
func showItems(items []*Item) {
	fmt.Println("===========")
	// itemsの要素を1つずつ取り出してitemに入れて繰り返す
	for i, item := range items {
		fmt.Printf("[%04d] %s:%d円\n", i+1, item.Category, item.Price)
	}
	fmt.Println("===========")
}
