package impl

import (
	"bufio"
	"log"
	"os"
	"strings"
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

	// 完全に異なるアプローチを使用: 単純なスキャンと直接マッチ
	var matches []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, pattern) {
			matches = append(matches, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("failed to scan file %s: %v", filePath, err)
	}

	return matches
}
