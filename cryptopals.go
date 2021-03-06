package main

import (
	"bufio"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
)

const (
	charSet       = " !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[]^_`abcdefghijklmnopqrstuvwxyz{|}~"
	lowerAlphaSet = "abcdefghijklmnopqrstuvwxyz"
	hexSet        = "0123456789ABCDEF"
	iceKey        = "ICE"
	maxKeysize    = 40
)

var (
	challenge1a     = "49276d206b696c6c696e6720796f757220627261696e206c696b65206120706f69736f6e6f7573206d757368726f6f6d"
	challenge1test  = "SSdtIGtpbGxpbmcgeW91ciBicmFpbiBsaWtlIGEgcG9pc29ub3VzIG11c2hyb29t"
	challenge2a     = "1c0111001f010100061a024b53535009181c"
	challenge2b     = "686974207468652062756c6c277320657965"
	challenge2test  = "746865206b696420646f6e277420706c6179"
	challenge3a     = "1b37373331363f78151b7f2b783431333d78397828372d363c78373e783a393b3736"
	challenge5a     = "Burning 'em, if you ain't quick and nimble\nI go crazy when I hear a cymbal"
	challenge5atest = "0b3637272a2b2e63622c2e69692a23693a2a3c6324202d623d63343c2a26226324272765272a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"

	challenge5b     = "I go crazy when I hear a cymbal"
	challenge5btest = "a282b2f20430a652e2c652a3124333a653e2b2027630c692b20283165286326302e27282f"
)

// engFreq represents the frequency of letter usage in the english language.
// http://en.algoritmy.net/article/40379/Letter-frequency-English
var engFreq = map[rune]float32{
	'a': 0.08167,
	'b': 0.01492,
	'c': 0.02782,
	'd': 0.04253,
	'e': 0.12702,
	'f': 0.02228,
	'g': 0.02015,
	'h': 0.06094,
	'i': 0.06966,
	'j': 0.00153,
	'k': 0.00772,
	'l': 0.04025,
	'm': 0.02406,
	'n': 0.06749,
	'o': 0.07507,
	'p': 0.01929,
	'q': 0.00095,
	'r': 0.05987,
	's': 0.06327,
	't': 0.09056,
	'u': 0.02758,
	'v': 0.00978,
	'w': 0.02360,
	'x': 0.00150,
	'y': 0.01974,
	'z': 0.00074,
}

type letterCounter map[rune]int

func newCounter() letterCounter {
	c := make(map[rune]int)
	for _, char := range lowerAlphaSet {
		c[char] = 0
	}
	return c
}

func main() {
	fmt.Println("Challenge 1")
	challenge1()

	fmt.Println("Challenge 2")
	challenge2()

	fmt.Println("Challenge 3")
	challenge3()

	fmt.Println("Challenge 4")
	challenge4()

	fmt.Println("Challenge 5")
	challenge5()

	fmt.Println("Challenge 6")
	challenge6()
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
	results, err := hexStringsValid([]string{string(challenge3a)})
	if err != nil {
		fmt.Println(err.Error())
		exit(false)
	}

	winner := scoreStrings(results)
	fmt.Printf("%s\n", winner)
}

func challenge4() {
	inputs, err := readStringsFromFilename("files/4.txt")
	if err != nil {
		fmt.Println(err.Error())
		exit(false)
	}

	results, err := hexStringsValid(inputs)
	if err != nil {
		fmt.Println(err.Error())
		exit(false)
	}

	winner := scoreStrings(results)
	fmt.Printf("%s", winner)
}

func challenge5() {
	encrypted := repeatKeyEncrypt([]byte(challenge5a), []byte(iceKey))
	test(hex.EncodeToString(encrypted), challenge5atest)
}

func challenge6() {
	file, err := readWholeFile("files/6.txt")
	if err != nil {
		fmt.Println(err.Error())
		exit(false)
	}
	file, err = base64.StdEncoding.DecodeString(string(file))
	if err != nil {
		fmt.Println(err.Error())
		exit(false)
	}

	curBest := 0
	bestNormal := 9999.00
	keySize := 1
	for keySize = 1; keySize <= maxKeysize; keySize++ {
		set1 := file[:keySize]
		set2 := file[keySize : keySize*2]
		set3 := file[keySize*2 : keySize*3]
		set4 := file[keySize*3 : keySize*4]
		set5 := file[keySize*4 : keySize*5]
		set6 := file[keySize*5 : keySize*6]
		set7 := file[keySize*6 : keySize*7]
		set8 := file[keySize*7 : keySize*8]
		set9 := file[keySize*8 : keySize*9]
		set10 := file[keySize*9 : keySize*10]
		set11 := file[keySize*10 : keySize*11]
		totalDist := 0
		for i := 0; i < keySize; i++ {
			totalDist += hammingDistance(set1[i], set2[i])
			totalDist += hammingDistance(set2[i], set3[i])
			totalDist += hammingDistance(set3[i], set4[i])
			totalDist += hammingDistance(set4[i], set5[i])
			totalDist += hammingDistance(set5[i], set6[i])
			totalDist += hammingDistance(set6[i], set7[i])
			totalDist += hammingDistance(set7[i], set8[i])
			totalDist += hammingDistance(set8[i], set9[i])
			totalDist += hammingDistance(set9[i], set10[i])
			totalDist += hammingDistance(set10[i], set11[i])
		}

		avgDist := totalDist / 10
		normalized := float64(avgDist) / float64(keySize)
		if normalized < bestNormal {
			bestNormal = normalized
			curBest = keySize
		}
	}

	// Split up the ciphertext into blocks of len(keysize).
	keySize = curBest
	split := [][]byte{}
	for i := 0; i < len(file); i += keySize {
		split = append(split, file[i:i+keySize])
	}

	transposed := [][]byte{}
	for i := 0; i < keySize; i++ {
		curB := []byte{}
		for _, b := range split {
			curB = append(curB, b[i])
		}
		transposed = append(transposed, curB)
	}

	// Score each string against the potential rotating xor byte.
	var fullkey []byte
	for i := range transposed {
		var topScore float32 = 9999.00
		var topKey byte
		for _, key := range charSet {
			enc := repeatKeyEncrypt(transposed[i], []byte{byte(key)})
			score := chiScore(string(enc))
			if score < topScore {
				topScore = score
				topKey = byte(key)
			}
		}
		fullkey = append(fullkey, topKey)
	}

	fileDecrypted := repeatKeyEncrypt(file, fullkey)
	fmt.Printf("key: %s, cleartext: %s\n", fullkey, fileDecrypted)
}

// hammingDistance calculates the bit difference between two bytes.
func hammingDistance(a, b byte) int {
	c := 0
	v := uint8(a ^ b)
	for c = 0; v != 0; c++ {
		v &= v - 1
	}
	return c
}

// repeatKeyEncrypt takes a byte buffer and a key, and encrypts the buffer with
// the given key with a rotating key xor operation. The encrypted buffer is
// returned.
func repeatKeyEncrypt(toEncrypt []byte, key []byte) []byte {
	result := []byte{}
	for idx, b := range toEncrypt {
		result = append(result, b^key[idx%len(key)])
	}
	return result
}

// hexStringsValid takes an array of hex-encoded strings, and returns the decoded
// strings only if they contain valid (readable) characters.
func hexStringsValid(inputs []string) ([]string, error) {
	results := []string{}
	for _, str := range inputs {
		decoded, err := hex.DecodeString(str)
		if err != nil {
			fmt.Println(err.Error())
			exit(false)
		}

		for _, char := range charSet {
			result := xorWithChar(decoded, byte(char))

			// Only append the result if all of the characters were readable.
			valid := true
			for _, c := range result {
				if !charIsValid(rune(c)) {
					valid = false
					break
				}
			}
			if valid {
				results = append(results, string(result))
			}
		}
	}
	return results, nil
}

// scoreStrings returns the top-scoring string according to the chiScore function.
func scoreStrings(toScore []string) string {
	// Score each sentence, and keep the top (lowest) one.
	var topScore float32 = 9999.00
	var topScored string
	for _, res := range toScore {
		score := chiScore(res)
		if score < topScore {
			topScore = score
			topScored = res
		}
	}
	return topScored
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

// filterValidChars returns a copy of input with the invalid (as determined
// by charIsValid()) filtered out.
func filterValidChars(input []byte) []byte {
	out := []byte{}
	for _, v := range input {
		if charIsValid(rune(v)) {
			fmt.Println(v)
			out = append(out, v)
		}
	}
	return out
}

// charIsValid returns true if the provided rune is uppercase, lowercase, a digit,
// or punctuation.
func charIsValid(a rune) bool {
	if unicode.IsUpper(a) {
		return true
	} else if unicode.IsLower(a) {
		return true
	} else if unicode.IsDigit(a) {
		return true
	} else if unicode.IsSpace(a) {
		return true
	} else if unicode.IsPunct(a) {
		return true
	}
	return false
}

// chiScore performs the chi-square test for the string, returning the result
// when compared to the frequency of letters contained in the string to the
// english language.
func chiScore(a string) float32 {
	// First count the number of occurrences of each character. We'll always
	// be working with lowercase strings.
	counter := newCounter()
	b := strings.ToLower(a)
	notCounted := 0
	for _, char := range b {
		if _, ok := counter[char]; ok {
			counter[char]++
		} else {
			notCounted++
		}
	}

	// Compare each of the counts against the frequency expected in the
	// population.
	var score float32 = 0.00
	length := len(a) - notCounted
	for char, freq := range engFreq {
		observed := float32(counter[char])
		expected := freq * float32(length)
		diff := observed - expected
		score += diff * diff / expected
	}

	return score
}

// readStringsFromFilename returns an array of strings, one for each line in the
// file.
func readStringsFromFilename(name string) ([]string, error) {
	resp := []string{}
	data, err := os.Open(name)
	defer data.Close()
	if err != nil {
		return resp, err
	}

	dScanner := bufio.NewScanner(data)
	for dScanner.Scan() {
		if err != nil {
			return resp, err
		}
		resp = append(resp, dScanner.Text())
	}
	return resp, nil
}

// readWholeFile returns the entire contents of a file in an array.
func readWholeFile(name string) ([]byte, error) {
	data, err := ioutil.ReadFile(name)
	if err != nil {
		return data, err
	}
	return data, nil
}
