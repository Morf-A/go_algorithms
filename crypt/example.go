package crypt

import "fmt"

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

func ExampleRSA() {
	pub, priv := RSAGenKeys()
	t := 'A'
	e := RSA{Key: pub}
	d := RSA{Key: priv}
	fmt.Println("plain:", string(t))
	cipher := e.Exp(uint64(t))
	fmt.Println("encrypted: ", cipher)
	t1 := d.Exp(cipher)
	fmt.Println("decrypted:", string(t1))
}
