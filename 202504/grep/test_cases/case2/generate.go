package main

import (
	"fmt"
	"math/rand"
	"os"
)

func main() {
	const totalLines = 1000000
	const pattern = "banana"
	f, _ := os.Create("textfile.txt")
	defer f.Close()
	for i := 0; i < totalLines; i++ {
		if rand.Intn(100) == 0 {
			fmt.Fprintf(f, "This line contains %s\n", pattern)
		} else {
			fmt.Fprintf(f, "This is a normal line %d\n", i)
		}
	}

	// input.txt を作成（textfile.txt, pattern名）
	fin, _ := os.Create("input.txt")
	defer fin.Close()
	fmt.Fprintln(fin, "textfile.txt")
	fmt.Fprintln(fin, pattern)
}
