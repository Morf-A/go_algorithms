package crypt

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

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

func (r *RSA) Exp(t uint64) uint64 {
	return ModExp(t, r.Key.E, r.Key.N)
}

type RSAKey struct {
	E uint64
	N uint64
}

func GetRandInt(random io.Reader, nbytes int) uint64 {
	if nbytes > 8 {
		panic("integer length overflow")
	}
	bytes := make([]byte, nbytes, 8)
	_, err := io.ReadFull(random, bytes)
	if err != nil {
		panic(err)
	}
	pad := make([]byte, 8-nbytes)
	bytes = append(bytes, pad...)
	return binary.LittleEndian.Uint64(bytes)
}

func GetRandPrime(random io.Reader, bytes int, min uint64) uint64 {
	for {
		p := GetRandInt(random, bytes)
		if p >= min && IsPrime(p) {
			return p
		}
	}
}

//we use Fermat's little theorem, but in real projects Miller–Rabin primality test is better
func IsPrime(p uint64) bool {
	return ModExp(2, p-1, p) == 1
}

func IsEven(x uint64) bool {
	return x%2 == 0
}

func ModExp(t, e, n uint64) uint64 {
	if e == 0 {
		return 1
	}
	if IsEven(e) {
		z := ModExp(t, e/2, n)
		return (z * z) % n
	}
	z := ModExp(t, (e-1)/2, n)
	return (z * z * t) % n
}

func GetCoprime(random io.Reader, r uint64) uint64 {
	for {
		e := GetRandPrime(random, 4, math.MaxUint8)
		//e - prime and p not divisible by e => greatest common factor == 1 => e coprime r
		if r%e != 0 {
			return e
		}
	}

}

//we use Euclidean algorithm
func GetMultInverse(e, r uint64) uint64 {
	g, i, j := Euclid(e, r)
	if g != 1 {
		panic("expected greatest common factor == 1")
	}

	fmt.Println("g, i, j: ", g, i, j)
	return uint64(j) % r
}

func Euclid(a, b uint64) (uint64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, i, j := Euclid(b, a%b)
	return g, j, i - int64(a/b)*j
}

func GenKeys(random io.Reader) (RSAKey, RSAKey) {
	p := GetRandPrime(random, 4, math.MaxUint16)
	q := GetRandPrime(random, 4, math.MaxUint16)
	n := p * q
	r := (p - 1) * (q - 1)
	e := GetCoprime(random, r)
	fmt.Println(n, e, r)
	d := GetMultInverse(e, r) // e * d % r == 1

	fmt.Println("e: ", e)
	fmt.Println("d: ", d)
	fmt.Println(e * d % r)

	// fmt.Printf("p: %064b\n", p)
	// fmt.Printf("q: %064b\n", q)
	// fmt.Printf("n: %064b\n", n)
	// fmt.Printf("r: %064b\n", r)
	// fmt.Printf("e: %064b\n", e)
	// fmt.Printf("d: %064b\n", d)

	return buildKeys(e, d, n)
}

func buildKeys(e, d, n uint64) (RSAKey, RSAKey) {
	return RSAKey{E: e, N: n}, RSAKey{E: d, N: n}
}

func RSAGenKeys() (RSAKey, RSAKey) {
	return buildKeys(5, 269, 493)
}
