package main

import (
	"fmt"
)

// 查找表，用于存储每个8位字节的位计数
var popCountTable [256]int

func init() {
	// 初始化查找表
	for i := 0; i < 256; i++ {
		popCountTable[i] = popCount(byte(i))
	}
}

// 计算一个字节中的位1的个数
func popCount(b byte) int {
	count := 0
	for i := 0; i < 8; i++ {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

func main() {
	// 创建一个初始位图，长度为10，所有位都设置为0
	bitmap := make([]byte, 2)

	// 设置第 13 位为 1
	setBit(&bitmap, 13)

	// 查询第 13 位的值
	fmt.Printf("Bit at position 13: %d\n", getBit(bitmap, 13))

	// 清除第 13 位的值（设置为 0）
	clearBit(bitmap, 13)

	// 查询第 13 位的值
	fmt.Printf("Bit at position 13 after clearing: %d\n", getBit(bitmap, 13))

	// 自动扩容，设置第 20 位为 1
	setBit(&bitmap, 200)

	// 查询第 20 位的值
	fmt.Printf("Bit at position 20: %d\n", getBit(bitmap, 200))

	setBit(&bitmap, 201)
	setBit(&bitmap, 202)
	setBit(&bitmap, 203)

	fmt.Println(countBits(bitmap))

	fmt.Println(bitmap)

	fmt.Println(popCountTable)
}

// 设置指定位为 1
func setBit(bitmap *[]uint8, pos int) {
	bytePos := pos / 8
	bitPos := pos % 8
	if bytePos >= len(*bitmap) {
		// 自动扩容，添加新的字节
		newBytes := bytePos - len(*bitmap) + 1
		*bitmap = append(*bitmap, make([]uint8, newBytes)...)
	}
	(*bitmap)[bytePos] |= 1 << uint(bitPos)
}

// 查询指定位的值
func getBit(bitmap []uint8, pos int) int {
	bytePos := pos / 8
	bitPos := pos % 8
	if bytePos >= len(bitmap) {
		// 如果字节位置超出范围，返回 0
		return 0
	}
	return int((bitmap[bytePos] >> uint(bitPos)) & 1)
}

// 清除指定位的值（设置为 0）
func clearBit(bitmap []uint8, pos int) {
	bytePos := pos / 8
	bitPos := pos % 8
	if bytePos < len(bitmap) {
		bitmap[bytePos] &= ^(1 << uint(bitPos))
	}
}

// 计算位图中1的个数
func countBits(bitmap []uint8) int {
	count := 0
	for _, byteValue := range bitmap {
		count += popCountTable[byteValue]
	}
	return count
}
