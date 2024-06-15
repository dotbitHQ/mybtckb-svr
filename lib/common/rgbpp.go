package common

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

func HexToLittleEndianHex(bytes []byte) (string, error) {
	// 反转字节数组以符合小端字节序
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	// 将反转后的字节数组编码为十六进制字符串
	littleEndianHex := hex.EncodeToString(bytes)
	return littleEndianHex, nil
}
func GetOutpointByargs(args []byte) (index uint32, txHash string) {
	fmt.Println(args)
	// 将字节数组转换为整数
	index = binary.LittleEndian.Uint32(args[:4])
	fmt.Println(index)

	txHash, err := HexToLittleEndianHex(args[4:])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	fmt.Println(txHash)
	return
}
