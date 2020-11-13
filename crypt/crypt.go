package crypt

import "math"

type Caesar struct {
	Key byte
}

func NewCaesar(key byte) *Caesar {
	return &Caesar{
		Key: key,
	}
}

func (r *Caesar) encryptBlock(p byte) byte {
	return p + r.Key
}

func (r *Caesar) decryptBlock(c byte) byte {
	return c - r.Key
}

func (r *Caesar) Encrypt(plain []byte) []byte {
	//1 block = 1 byte
	var x byte //iv = 0
	cipher := make([]byte, 0, len(plain))
	for _, p := range plain {
		c := r.encryptBlock(p ^ x)
		cipher = append(cipher, c)
		x = c
	}
	return cipher
}

func (r *Caesar) Decrypt(cipher []byte) []byte {
	//1 block = 1 byte
	var x byte //iv = 0
	plain := make([]byte, 0, len(cipher))
	for _, c := range cipher {
		p := r.decryptBlock(c) ^ x
		plain = append(plain, p)
		x = c
	}
	return plain
}

type RSA struct {
	Key RSAKey
}

func (r *RSA) Exp(t int) int {
	var res float64 = 1
	tf := float64(t)
	for i := 1; i <= r.Key.e; i++ {
		res = math.Mod(res*tf, float64(r.Key.n))
	}
	return int(res)
}

type RSAKey struct {
	e int
	n int
}

func RSAGenKeys() (RSAKey, RSAKey) {
	e := 5
	d := 269
	n := 493
	return RSAKey{e: e, n: n}, RSAKey{e: d, n: n}
}
