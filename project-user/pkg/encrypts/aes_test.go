package encrypts

import (
	"fmt"
	"testing"

	"test.com/project-user/config"
)

func TestEncrypt(t *testing.T) {
	plain := "3"
	// AES 规定有3种长度的key: 16, 24, 32分别对应AES-128, AES-192, or AES-256
	//key := "abcdefgehjhijkmlkjjwwoew"
	// 加密
	cipherByte, err := Encrypt(plain, config.AESKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s ==> %s\n", plain, cipherByte)
	// 解密
	plainText, err := Decrypt(cipherByte, config.AESKey)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s ==> %s\n", cipherByte, plainText)
	if plainText == plain {
		fmt.Println("解密成功")
	}
}
