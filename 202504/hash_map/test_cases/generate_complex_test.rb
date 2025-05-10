#!/usr/bin/env ruby
# ハッシュマップの複雑なテストケースを生成するスクリプト
# 使用方法: ruby generate_complex_test.rb [出力ディレクトリ] [操作数]

require 'json'
require 'fileutils'

# コマンドライン引数の処理
output_dir = ARGV[0] || './test_cases/case7'
operations_count = (ARGV[1] || '5000').to_i

# 出力ディレクトリの作成
FileUtils.mkdir_p(output_dir)

# ハッシュ衝突を発生させやすいキーのパターン
def generate_collision_prone_keys(base, count)
  keys = []
  count.times do |i|
    keys << "#{base}_#{i}"
  end
  keys
end

# アナグラムキーを生成（ハッシュ衝突を起こしやすい）
def generate_anagram_keys(base, count)
  keys = []
  count.times do |i|
    # 文字列をシャッフルして異なるアナグラムを生成
    chars = base.chars
    4.times do
      chars.shuffle!
      keys << "anagram_#{chars.join}"
    end
  end
  keys.uniq[0...count]
end

# 同じハッシュ値を持つ可能性が高いキーを生成
def generate_same_hash_keys(count)
  keys = []
  count.times do |i|
    # FNV-1aハッシュ関数で衝突しやすいキーパターン
    keys << "hash_#{i * 16777619 % 1000}"
  end
  keys
end

# 操作の種類
operations = ['put', 'get', 'remove']

# 入力データの生成
input_data = []
expected_data = {}

# 衝突しやすいキーを生成
collision_keys = generate_collision_prone_keys('collision', 100)
anagram_keys = generate_anagram_keys('abcdefgh', 100)
same_hash_keys = generate_same_hash_keys(100)

all_keys = collision_keys + anagram_keys + same_hash_keys

# 操作の生成
operations_count.times do |i|
  op = operations[i % 3] # put, get, removeを順番に実行
  
  case op
  when 'put'
    key = all_keys[i % all_keys.length]
    value = i * 100 + (key.hash % 100)
    input_data << { 'action' => 'put', 'key' => key, 'value' => value }
    expected_data[key] = value
  when 'get'
    # 直前のputで追加したキーを取得
    if i > 0 && input_data[i-1]['action'] == 'put'
      key = input_data[i-1]['key']
      input_data << { 'action' => 'get', 'key' => key }
    else
      # 存在するキーをランダムに選択
      existing_keys = expected_data.keys
      if existing_keys.any?
        key = existing_keys.sample
        input_data << { 'action' => 'get', 'key' => key }
      else
        # キーがない場合はputを実行
        key = all_keys[i % all_keys.length]
        value = i * 100
        input_data << { 'action' => 'put', 'key' => key, 'value' => value }
        expected_data[key] = value
      end
    end
  when 'remove'
    # 存在するキーをランダムに選択して削除
    existing_keys = expected_data.keys
    if existing_keys.any?
      key = existing_keys.sample
      input_data << { 'action' => 'remove', 'key' => key }
      expected_data.delete(key)
    else
      # キーがない場合はputを実行
      key = all_keys[i % all_keys.length]
      value = i * 100
      input_data << { 'action' => 'put', 'key' => key, 'value' => value }
      expected_data[key] = value
    end
  end
end

# 最後に全エントリを取得する操作を追加
input_data << { 'action' => 'all_entries' }

# 入力ファイルの書き込み
File.open(File.join(output_dir, 'input.txt'), 'w') do |f|
  f.write(JSON.pretty_generate(input_data))
end

# 期待される出力ファイルの書き込み
File.open(File.join(output_dir, 'expected.txt'), 'w') do |f|
  f.write(JSON.pretty_generate(expected_data))
end

puts "生成完了: #{operations_count}個の操作を含むテストケースを#{output_dir}に作成しました"
