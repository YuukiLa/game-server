package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

var key []byte = []byte("gamegamegamegame")

//密码加密
func HashPwd(pwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
	}
	return string(hash)
}

//密码比较
func EqualsPwd(fPwd, dPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dPwd), []byte(fPwd))
	return err == nil
}

//使用PKCS5进行填充

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//使用PKCS7进行填充
func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//aes加密，填充秘钥key的16位，24,32分别对应AES-128, AES-192, or AES-256.
func AesCBCEncrypt(rawData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	//填充原文
	blockSize := block.BlockSize()

	rawData = PKCS5Padding(rawData, blockSize)

	//初始向量IV必须是唯一，但不需要保密
	cipherText := make([]byte, blockSize+len(rawData))

	//block大小 16
	iv := cipherText[:blockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	//block大小和初始向量大小一定要一致
	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(cipherText[blockSize:], rawData)

	return cipherText, nil
}

func AesCBCDncrypt(encryptData []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	blockSize := block.BlockSize()

	if len(encryptData) < blockSize {
		panic("ciphertext too short")
	}
	iv := encryptData[:blockSize]

	encryptData = encryptData[blockSize:]
	fmt.Println("encryptData=", encryptData)
	fmt.Println(len(encryptData) % blockSize)

	// CBC mode always works in whole blocks.
	if len(encryptData)%blockSize != 0 {
		panic("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(encryptData, encryptData)
	//解填充
	encryptData = PKCS5UnPadding(encryptData)
	return encryptData, nil
}

func Encrypt(rawData string) (string, error) {
	rawDataByte := []byte(fmt.Sprintf("%v_%v", rawData, time.Now().Unix()))
	data, err := AesCBCEncrypt(rawDataByte)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func Dncrypt(rawData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(rawData)
	if err != nil {
		return "", err
	}
	dnData, err := AesCBCDncrypt(data)
	if err != nil {
		return "", err
	}
	arr := strings.Split(string(dnData), "_")
	genTime, _ := strconv.ParseInt(arr[1], 10, 64)
	if time.Now().Unix()-genTime > 7200 {
		return "", errors.New("token过期")
	}
	return arr[0], nil
}
