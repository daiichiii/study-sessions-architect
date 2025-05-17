package impl

import (
	"bufio"
	"log"
	"os"
	"strings"
	"sync"
)

// GrepImplementation はGrepの基本実装を提供する
type GrepImplementation struct{}

// Search はファイルから特定のパターンを検索する
func (g *GrepImplementation) Search(filePath, pattern string) []string {
	// ファイルを開く
	f, err := os.Open(filePath)
	if err != nil {
		log.Printf("failed to open file %s: %v", filePath, err)
		return nil
	}
	defer f.Close()

	// まず全ての行を読み込む
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Printf("failed to scan file %s: %v", filePath, err)
		return nil
	}

	// 並列処理のためのワーカー数を定義
	workerCount := 4 // 並列ワーカー数（CPUコア数などに応じて調整可能）

	// 結果を格納するスライス
	results := make([]string, 0, len(lines)/2) // 容量は適当に見積もる
	var mutex sync.Mutex

	// 処理を分割
	chunkSize := (len(lines) + workerCount - 1) / workerCount
	if chunkSize < 1 {
		chunkSize = 1
	}

	// ワーカーの終了を待つためのWaitGroup
	var wg sync.WaitGroup

	// 各ワーカーに作業を分配
	for i := 0; i < len(lines); i += chunkSize {
		wg.Add(1)
		go func(start int) {
			defer wg.Done()
			end := start + chunkSize
			if end > len(lines) {
				end = len(lines)
			}

			// このワーカーの担当範囲の行を処理
			var localMatches []string
			for j := start; j < end; j++ {
				if strings.Contains(lines[j], pattern) {
					localMatches = append(localMatches, lines[j])
				}
			}

			// 結果をマージ（スレッドセーフに）
			mutex.Lock()
			results = append(results, localMatches...)
			mutex.Unlock()
		}(i)
	}

	// すべてのワーカーが完了するのを待つ
	wg.Wait()

	return results
}
