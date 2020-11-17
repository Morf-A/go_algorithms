package crypt

import (
	"fmt"
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

func ExampleRSA() {

	p := GetRandInt(ReadOne{}, 2)
	p >>= 1
	fmt.Printf("%d: %016b\n", p, p)

	p2 := GetRandInt(ReadOne{}, 2)
	pp := p * p2
	fmt.Printf("%d: %064b\n", pp, pp)

	pppp := pp * pp
	fmt.Printf("%d: %064b\n", pppp, pppp)

	// pub, priv := RSAGenKeys(rand.Reader)
	// t := 'A'
	// e := RSA{Key: pub}
	// d := RSA{Key: priv}
	// fmt.Println("plain:", string(t))
	// cipher := e.Exp(int64(t))
	// fmt.Println("encrypted: ", cipher)
	// t1 := d.Exp(cipher)
	// fmt.Println("decrypted:", string(t1))
}
