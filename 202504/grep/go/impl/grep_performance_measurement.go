package impl

import (
	"fmt"
	"io/ioutil"
	"strings"

	utils "study-session-architect/utils/go"
)

// loadGrepTestData は入力ファイルと期待値ファイルを読み込む
func loadGrepTestData(fileDir string) (string, string, []string, error) {
	// 入力データの読み込み
	inputData, err := ioutil.ReadFile(strings.Join([]string{fileDir, "input.txt"}, "/"))
	if err != nil {
		return "", "", nil, fmt.Errorf("入力ファイルの読み込みに失敗しました: %v", err)
	}

	// 入力データのパース
	inputLines := strings.Split(string(inputData), "\n")
	if len(inputLines) < 2 {
		return "", "", nil, fmt.Errorf("入力データが不足しています")
	}

	filePath := strings.Join([]string{fileDir, strings.TrimSpace(inputLines[0])}, "/")
	pattern := strings.TrimSpace(inputLines[1])

	// 期待値の読み込み
	expectedData, err := ioutil.ReadFile(strings.Join([]string{fileDir, "expected.txt"}, "/"))
	if err != nil {
		return filePath, pattern, nil, fmt.Errorf("期待値ファイルの読み込みに失敗しました: %v", err)
	}

	// 期待値のパース
	var expectedOutput []string
	for _, line := range strings.Split(string(expectedData), "\n") {
		if line != "" {
			expectedOutput = append(expectedOutput, line)
		}
	}

	return filePath, pattern, expectedOutput, nil
}

// MeasureGrepPerformance はGrepの性能と正当性を計測する
func MeasureGrepPerformance(fileDir string, iterations int) map[string]interface{} {
	var err error
	filePath, pattern, expectedOutput, err := loadGrepTestData(fileDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	grep := &GrepImplementation{}

	fmt.Printf("Grep実装のパフォーマンス計測と正当性検証:\n")
	fmt.Printf("ファイル: %s\n", filePath)
	fmt.Printf("検索パターン: %s\n", pattern)
	fmt.Printf("繰り返し回数: %d\n", iterations)

	var matchingLines []string

	// 処理時間とメモリ使用量を計測
	results := utils.MeasurePerformance("Grep", func() {
		for i := 0; i < iterations; i++ {
			matchingLines = grep.Search(filePath, pattern)
			if iterations == 1 {
				fmt.Printf("ヒット数: %d\n", len(matchingLines))
			}
		}
	})

	// 正当性検証
	valid := utils.VerifyResult("Grep", matchingLines, expectedOutput)
	results["valid"] = valid

	return results
}
