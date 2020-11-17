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

//check that t < n
func (r *RSA) Encrypt(t uint16) (uint32, error) {
	if err := r.CheckT(uint32(t)); err != nil {
		return 0, err
	}
	return r.exp(uint32(t)), nil
}

func (r *RSA) CheckT(t uint32) error {
	if t >= r.Key.N {
		return fmt.Errorf("t can`t be bigger than n %d %d", t, r.Key.N)
	}
	return nil
}

func (r *RSA) Decrypt(t uint32) (uint16, error) {
	if err := r.CheckT(t); err != nil {
		return 0, err
	}
	return uint16(r.exp(t)), nil
}

func (r *RSA) exp(t uint32) uint32 {
	return uint32(ModExp(uint64(t), uint64(r.Key.E), uint64(r.Key.N)))
}

type RSAKey struct {
	E uint32
	N uint32
}

func GetRandInt(random io.Reader, nbites int) uint64 {
	if nbites > 63 || nbites < 0 {
		panic(fmt.Sprintf("bad nbites param (%d)", nbites))
	}
	fullBytes := nbites / 8
	rBites := nbites % 8
	if rBites > 0 {
		fullBytes++
	}
	bytes := make([]byte, fullBytes, 8)
	_, err := io.ReadFull(random, bytes)
	if err != nil {
		panic(err)
	}
	pad := make([]byte, 8-fullBytes)
	bytes = append(bytes, pad...)
	res := binary.LittleEndian.Uint64(bytes)
	if rBites > 0 {
		res >>= (8 - rBites)
	}
	return res
}

func GetRandPrime(random io.Reader, bytes int) uint64 {
	for {
		p := GetRandInt(random, bytes)
		if p > 1 && IsPrime(p) {
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
	return ((z * z) % n * t) % n
}

func GetCoprime(random io.Reader, r uint64, nbits int) uint64 {
	for {
		e := GetRandPrime(random, nbits)
		//e - prime and p not divisible by e => greatest common factor == 1 => e coprime r
		if r%e != 0 {
			return e
		}
	}

}

//we use Euclidean algorithm
func GetMultInverse(e, r int64) int64 {
	g, i, _ := Euclid(e, r)
	if g != 1 {
		panic("expected greatest common factor == 1")
	}
	if i < 0 {
		return r - int64(-1*i)%r
	}
	return int64(i) % r
}

func Euclid(a, b int64) (int64, int64, int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, i, j := Euclid(b, a%b)
	return g, j, i - a/b*j
}

type RSADigits struct {
	P uint16
	Q uint16
	N uint32
	R uint32
	E uint32
	D uint32
}

func (r RSADigits) GetKeyPair() (RSAKey, RSAKey) {
	return RSAKey{E: r.E, N: r.N}, RSAKey{E: r.D, N: r.N}
}

func (r RSADigits) Debug() {
	fmt.Printf("p: %016b\n", r.P)
	fmt.Printf("q: %016b\n", r.Q)
	fmt.Printf("n: %032b\n", r.N)
	fmt.Printf("r: %032b\n", r.R)
	fmt.Printf("e: %032b\n", r.E)
	fmt.Printf("d: %032b\n", r.D)
}

func RSAGenDigits(random io.Reader) RSADigits {
	var p, q uint16
	var n uint32
	for {
		p = uint16(GetRandPrime(random, 16))
		q = uint16(GetRandPrime(random, 15))
		n = uint32(p) * uint32(q)
		if n > math.MaxUint16 {
			break
		}
	}

	r := uint32(p-1) * uint32(q-1)
	e := uint32(GetCoprime(random, uint64(r), 16))
	d := uint32(GetMultInverse(int64(e), int64(r))) // e * d % r == 1
	return RSADigits{
		P: p,
		Q: q,
		N: n,
		R: r,
		E: e,
		D: d,
	}
}
