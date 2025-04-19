package main

import (
	"fmt"
	"os"

	impl "study-session-architect/grep/go/impl"
)

// ===============================================
// メイン関数
// ===============================================

func main() {
	fmt.Println("==============================")
	fmt.Println("Grep性能計測と正当性検証")
	fmt.Println("==============================")

	// Grep計測と検証
	fmt.Println("Grep実装のテスト")
	fileDir := os.Args[1]
	grepResults := impl.MeasureGrepPerformance(fileDir, 1)

	// 検証結果の要約
	fmt.Println("\n==============================")
	fmt.Println("テスト結果サマリー")
	fmt.Println("==============================")

	grepValid := false
	if grepResults != nil {
		grepValid, _ = grepResults["valid"].(bool)
	}

	fmt.Printf("Grep: %s\n", boolToCheckmark(grepValid))
}

// boolToCheckmark はブール値をチェックマーク文字列に変換
func boolToCheckmark(b bool) string {
	if b {
		return "成功 ✓"
	}
	return "失敗 ✗"
}
