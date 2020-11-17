package crypt

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math/big"
)

func ExampleCBC() {
	caesar := NewCaesar(16)

	input1 := []byte("ABC")
	fmt.Println("input1:", input1, string(input1))
	c1 := caesar.Encrypt(input1)
	fmt.Println("c1:", c1, string(c1))
	p1 := caesar.Decrypt(c1)
	fmt.Println("p1:", p1, string(p1))

	input2 := []byte("AAAAAAAAAA")
	fmt.Println("input2:", input2, string(input2))
	c2 := caesar.Encrypt(input2)
	fmt.Println("c2:", c2, string(c2))
	p2 := caesar.Decrypt(c2)
	fmt.Println("p2:", p2, string(p2))
}

type ReadOne struct{}

func (r ReadOne) Read(p []byte) (int, error) {
	for i := 0; i < len(p); i++ {
		p[i] = 255
	}
	return len(p), nil
}

//can make mistakes
func IsPrime(p uint64) bool {
	return ModExp(2, p-1, p) == 1
}

func IsPrimeMR(p uint64) bool {
	bigP := new(big.Int)
	bigP.SetUint64(p)
	return bigP.ProbablyPrime(50)
}

func ExampleRSA() {
	digits := RSAGenDigits(rand.Reader, IsPrimeMR)
	digits.Debug()
	pub, priv := digits.GetKeyPair()
	t := "HI"
	encryptor := RSA{Key: pub}
	decryptor := RSA{Key: priv}
	fmt.Println("plain:", t)
	tint := binary.BigEndian.Uint16([]byte(t))
	fmt.Println("tint: ", tint)
	cipher, err := encryptor.Encrypt(tint)
	if err != nil {
		panic(err)
	}
	fmt.Println("encrypted: ", cipher)
	t1int, err := decryptor.Decrypt(cipher)
	if err != nil {
		panic(err)
	}
	fmt.Println("t1int: ", t1int)
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, t1int)
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted:", buf.String())
	if buf.String() != t {
		panic("not equal")
	}
}
