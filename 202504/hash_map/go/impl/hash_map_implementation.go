package impl

import (
	"fmt"
)

// デフォルト値
const (
	DefaultBucketSize = 256  // 初期バケットサイズを増やす
	DefaultLoadFactor = 0.65 // 負荷係数を少し小さくして衰突を減らす
)

// entry はキーと値のペアを表す
type entry struct {
	key   string
	hash  uint32 // ハッシュ値をキャッシュして文字列比較を減らす
	value interface{}
}

// bucket はチェイニング方式で衝突を解決するためのバケット
type bucket struct {
	entries []*entry
}

// HashMapImplementation はハッシュマップの基本実装を提供する
type HashMapImplementation struct {
	buckets    []*bucket // バケット配列
	size       int       // 現在の要素数
	bucketSize int       // バケット配列のサイズ
	loadFactor float64   // リサイズのトリガーとなる負荷係数
}

// NewHashMap は新しいHashMapを作成する
func NewHashMap(bucketSize int) *HashMapImplementation {
	// バケットサイズが0以下の場合はデフォルト値を使用
	if bucketSize <= 0 {
		bucketSize = DefaultBucketSize
	}

	// バケット配列を初期化
	buckets := make([]*bucket, bucketSize)
	for i := 0; i < bucketSize; i++ {
		buckets[i] = &bucket{entries: make([]*entry, 0)}
	}

	return &HashMapImplementation{
		buckets:    buckets,
		size:       0,
		bucketSize: bucketSize,
		loadFactor: DefaultLoadFactor,
	}
}

// computeHash はキーのハッシュ値を計算する
func (h *HashMapImplementation) computeHash(key string) uint32 {
	// FNV-1aハッシュ関数を最適化
	const (
		offsetBasis = 2166136261
		fnvPrime    = 16777619
	)

	hash := uint32(offsetBasis)
	bytes := []byte(key)

	// 4バイトずつ処理して高速化
	for i := 0; i < len(bytes); i += 4 {
		// 4バイトずつ処理（可能な場合）
		if i+4 <= len(bytes) {
			k := uint32(bytes[i]) | uint32(bytes[i+1])<<8 | uint32(bytes[i+2])<<16 | uint32(bytes[i+3])<<24
			hash ^= k
			hash *= fnvPrime
			hash ^= hash >> 16
		} else {
			// 残りのバイトを処理
			for j := i; j < len(bytes); j++ {
				hash ^= uint32(bytes[j])
				hash *= fnvPrime
			}
		}
	}

	return hash
}

// convertKey はキーをハッシュマップで使用する文字列に変換する
func (h *HashMapImplementation) convertKey(key interface{}) string {
	// 文字列型の場合はそのまま返す
	if strKey, ok := key.(string); ok {
		return strKey
	}
	// その他の型はfmt.Sprintfを使用
	return fmt.Sprintf("%v", key)
}

// Put はキーと値のペアを格納する
func (h *HashMapImplementation) Put(key, value interface{}) {
	// キーを文字列に変換（最適化された方法）
	strKey := h.convertKey(key)

	// 負荷係数をチェックし、必要に応じてリサイズ
	if float64(h.size)/float64(h.bucketSize) >= h.loadFactor {
		h.resize()
	}

	// ハッシュ値を計算
	hashValue := h.computeHash(strKey)
	// ハッシュ値からインデックスを計算
	index := int(hashValue & uint32(h.bucketSize-1))

	// バケットへの参照を一度だけ取得
	bucket := h.buckets[index]
	entries := bucket.entries

	// バケット内でキーを検索
	for i, e := range entries {
		// まずハッシュ値を比較
		if e.hash == hashValue && e.key == strKey {
			// 既存のキーの場合は値を更新
			entries[i].value = value
			return
		}
	}

	// 新しいエントリを追加する前にスライスの容量を確認
	if cap(entries) == len(entries) {
		// 容量を2倍に増やす（頻繁な再割り当てを避ける）
		newEntries := make([]*entry, len(entries), max(4, len(entries)*2))
		copy(newEntries, entries)
		entries = newEntries
		bucket.entries = entries
	}

	// 新しいキーの場合はエントリを追加
	bucket.entries = append(entries, &entry{key: strKey, hash: hashValue, value: value})
	h.size++
}

// Get はキーに対応する値を取得する
func (h *HashMapImplementation) Get(key interface{}) (interface{}, bool) {
	// キーを文字列に変換
	strKey := h.convertKey(key)

	// ハッシュ値を計算
	hashValue := h.computeHash(strKey)
	// ハッシュ値からインデックスを計算
	index := int(hashValue & uint32(h.bucketSize-1))

	// バケットへの参照を直接取得
	entries := h.buckets[index].entries

	// バケット内でキーを検索（長さに関わらず同じ検索方法）
	for _, e := range entries {
		// まずハッシュ値を比較し、次に文字列を比較
		if e.hash == hashValue && e.key == strKey {
			return e.value, true
		}
	}

	// キーが見つからない場合
	return nil, false
}

// Remove はキーに対応するエントリを削除する
func (h *HashMapImplementation) Remove(key interface{}) bool {
	// キーを文字列に変換
	strKey := h.convertKey(key)

	// ハッシュ値を計算
	hashValue := h.computeHash(strKey)
	// ハッシュ値からインデックスを計算
	index := int(hashValue & uint32(h.bucketSize-1))

	// バケットへの参照を直接取得
	bucket := h.buckets[index]
	entries := bucket.entries

	// バケット内でキーを検索
	for i, e := range entries {
		// まずハッシュ値を比較し、次に文字列を比較
		if e.hash == hashValue && e.key == strKey {
			// 最後の要素と入れ替えて削除（O(1)の操作）
			lastIdx := len(entries) - 1
			if i != lastIdx {
				// 削除する要素が最後でない場合は入れ替え
				entries[i] = entries[lastIdx]
			}
			// スライスを縮小（再割り当てなし）
			bucket.entries = entries[:lastIdx]
			h.size--
			return true
		}
	}

	// キーが見つからない場合
	return false
}

// resize はバケットサイズを拡張する
func (h *HashMapImplementation) resize() {
	// 古いバケットを保存
	oldBuckets := h.buckets
	oldSize := h.bucketSize

	// バケットサイズを2倍に拡張
	h.bucketSize *= 2

	// 新しいバケットの初期化
	// 平均的なバケットサイズに基づいてメモリを事前に確保
	averageEntriesPerBucket := max(2, h.size/oldSize)
	h.buckets = make([]*bucket, h.bucketSize)

	for i := 0; i < h.bucketSize; i++ {
		// 各バケットに十分な初期容量を事前に確保
		h.buckets[i] = &bucket{entries: make([]*entry, 0, averageEntriesPerBucket)}
	}

	// 古いバケットからすべてのエントリを新しいバケットに再配置
	for _, bucket := range oldBuckets {
		for _, e := range bucket.entries {
			// 既に計算されたハッシュ値を使用
			// ビットマスクを使用して高速にインデックスを計算
			index := int(e.hash & uint32(h.bucketSize-1))

			// 新しいバケットに直接追加
			h.buckets[index].entries = append(h.buckets[index].entries, e)
		}
	}

	// サイズは変更しない（エントリ数は同じ）
}

// Size は現在の要素数を取得する
func (h *HashMapImplementation) Size() int {
	return h.size
}

// GetAllEntries は全てのエントリを取得する（テスト用）
func (h *HashMapImplementation) GetAllEntries() map[string]interface{} {
	result := make(map[string]interface{})

	// すべてのバケットを走査
	for _, bucket := range h.buckets {
		for _, e := range bucket.entries {
			result[e.key] = e.value
		}
	}

	return result
}
