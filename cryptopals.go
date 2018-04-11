package main

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
)

const (
	charSet = " !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
)

var (
	challenge1a    = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	challenge1test = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	challenge2a    = "1c0111001f010100061a024b53535009181c"
	challenge2b    = "686974207468652062756c6c277320657965"
	challenge2test = "746865206b696420646f6e277420706c6179"
	challenge3a    = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
)

func main() {
	fmt.Println("Challenge 1")
	challenge1()

	fmt.Println("Challenge 2")
	challenge2()

	fmt.Println("Challenge 3")
	challenge3()
}

func exit(success bool) {
	status := 0
	if !success {
		fmt.Println("FAIL")
		status = 1
	}
	os.Exit(status)
}

func test(result, test string) {
	if result == test {
		fmt.Printf("Success! %s\n", result)
	} else {
		fmt.Printf("Fail! %s\n", result)
		exit(false)
	}
}

func challenge1() {
	res, err := hexToBase64([]byte(challenge1a))
	if err != nil {
		fmt.Println(err.Error())
		exit(false)
	}
	test(string(res), challenge1test)
}

func challenge2() {
	left, err := hex.DecodeString(challenge2a)
	right, err := hex.DecodeString(challenge2b)
	res, err := xorBuffer(left, right)
	if err != nil {
		fmt.Println(err.Error())
		exit(false)
	}

	test(hex.EncodeToString(res), challenge2test)
}

func challenge3() {
	input, err := hex.DecodeString(challenge3a)
	if err != nil {
		fmt.Println(err.Error())
		exit(false)
	}

	results := []string{}
	for _, char := range charSet {
		result := xorWithChar(input, byte(char))
		results = append(results, hex.EncodeToString(result))
	}

	for _, res := range results {
		a, _ := hex.DecodeString(res)
		fmt.Printf("%s\n", a)
	}
}

// hexToBase64 returns a []byte which is the base64 encoded string containing
// the hex input.
func hexToBase64(hexStr []byte) ([]byte, error) {
	decoded, err := hex.DecodeString(string(hexStr))
	if err != nil {
		return nil, err
	}
	encoded := base64.StdEncoding.EncodeToString(decoded)
	return []byte(encoded), nil
}

// xorBuffer takes 2 buffers, xors each of their elements, and returns the resulting
// buffer.
func xorBuffer(a, b []byte) ([]byte, error) {
	resp := []byte{}
	if len(a) != len(b) {
		return resp, errors.New("xorBuffer: inputs must be the same length")
	}
	for idx := range a {
		resp = append(resp, a[idx]^b[idx])
	}
	return resp, nil
}

// xorWithChar takes a buffer, xors each element with the given byte, and
// returns the resulting buffer.
func xorWithChar(a []byte, b byte) []byte {
	resp := []byte{}
	for idx := range a {
		resp = append(resp, a[idx]^b)
	}
	return resp
}
