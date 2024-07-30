package Management

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
)

type Aes struct {
	Key       []byte
	BlockSize int
}

func NewAes(key []byte, blockSize int) *Aes {
	return &Aes{Key: key, BlockSize: blockSize}
}

func (at *Aes) padding(src []byte) []byte {
	//填充个数
	padding := aes.BlockSize - len(src)%aes.BlockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, paddingText...)
}

func (at *Aes) unPadding(src []byte) []byte {
	size := len(src)
	return src[:(size - int(src[size-1]))]
}

func (at *Aes) Encrypt(src []byte) ([]byte, error) {
	//key只能是 16 24 32长度
	block, err := aes.NewCipher(at.Key)
	if err != nil {
		return nil, err
	}
	//填充
	src = at.padding(src)
	//返回加密结果
	encryptData := make([]byte, len(src))
	//存储每次加密的数据
	tmpData := make([]byte, at.BlockSize)

	//分组分块加密
	for index := 0; index < len(src); index += at.BlockSize {
		block.Encrypt(tmpData, src[index:index+at.BlockSize])
		copy(encryptData[index:index+at.BlockSize], tmpData)
	}
	return encryptData, nil
}

func (at *Aes) Decrypt(src []byte) ([]byte, error) {
	//key只能是 16 24 32长度
	block, err := aes.NewCipher(at.Key)
	if err != nil {
		return nil, err
	}
	//返回加密结果
	decryptData := make([]byte, len(src))
	//存储每次加密的数据
	tmpData := make([]byte, at.BlockSize)

	//分组分块加密
	for index := 0; index < len(src); index += at.BlockSize {
		block.Decrypt(tmpData, src[index:index+at.BlockSize])
		copy(decryptData[index:index+at.BlockSize], tmpData)
	}
	return at.unPadding(decryptData), nil
}

func (at *Aes) EncryptString(src string) string {

	b, err := at.Encrypt([]byte(src))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b)
}

func (at *Aes) DecryptString(src string) ([]byte, error) {

	decodeBytes, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return nil, err
	}
	return at.Decrypt(decodeBytes)
}
