package leb128

import (
	"math/big"
	"sync"
)

// LEB128 code based on the sample here: https://en.wikipedia.org/wiki/LEB128.

// Typed wrapper around a sync.Pool
type bigPool struct {
	pool sync.Pool
}

// Retrieves a big.Int from the pool. Note that the returned big.Int can have any value
// and so must be manually initialized after retrieval.
func (b *bigPool) Get() *big.Int {
	out := b.pool.Get()
	return out.(*big.Int)
}

// Puts a big.Int into the pool.
func (b *bigPool) Put(val *big.Int) {
	b.pool.Put(val)
}

// Shared instance of bigPool for use in the methods below.
var pool = &bigPool{
	pool: sync.Pool{
		New: func() interface{} {
			return big.NewInt(0)
		},
	},
}

// Cached big.Int of 1.
var one = big.NewInt(1)

// Cached big.Int of 128.
var oneTwentyEight = big.NewInt(128)

// FromUInt64 encodes n with LEB128 and returns the encoded bytes.
func FromUInt64(n uint64) (out []byte) {
	more := true
	for more {
		b := byte(n & 0x7F)
		n >>= 7
		if n == 0 {
			more = false
		} else {
			b = b | 0x80
		}
		out = append(out, b)
	}
	return
}

// ToUInt64 decodes LEB128-encoded bytes into a uint64.
func ToUInt64(encoded []byte) uint64 {
	var result uint64
	var shift, i uint
	for {
		b := encoded[i]
		result |= (uint64(0x7F & b)) << shift
		if b&0x80 == 0 {
			break
		}
		shift += 7
		i++
	}
	return result
}

// FromBigInt encodes the signed big integer n in two's complement,
// LEB128-encodes it, and returns the encoded bytes.
func FromBigInt(n *big.Int) (out []byte) {
	size := n.BitLen()
	negative := n.Sign() < 0
	if negative {
		// big.Int stores integers using sign and magnitude. Returns a copy
		// as the code below is destructive.
		n = twosComplementBigInt(n)
	} else {
		// The code below is destructive so make a copy.
		n = big.NewInt(0).Set(n)
	}

	more := true
	for more {
		// Note that bBigInt below is set to whatever the
		// modulus result of DivMod is, and thus doesn't
		// need to be initialized manually.
		bBigInt := pool.Get()
		n.DivMod(n, oneTwentyEight, bBigInt) // This does the mask and the shift.
		b := uint8(bBigInt.Int64())
		pool.Put(bBigInt)

		// We just logically right-shifted the bits of n so we need to sign extend
		// if n is negative (this simulates an arithmetic shift).
		if negative {
			signExtend(n, size)
		}

		if (n.Sign() == 0 && b&0x40 == 0) ||
			(negative && equalsNegativeOne(n, size) && b&0x40 > 0) {
			more = false
		} else {
			b = b | 0x80
		}
		out = append(out, b)
	}
	return
}

// ToBigInt decodes the signed big integer found in the given bytes.
func ToBigInt(encoded []byte) *big.Int {
	result := big.NewInt(0)
	var shift, i int
	var b byte
	size := len(encoded) * 8

	for {
		b = encoded[i]
		for bitPos := uint(0); bitPos < 7; bitPos++ {
			result.SetBit(result, 7*i+int(bitPos), uint((b>>bitPos)&0x01))
		}
		shift += 7
		if b&0x80 == 0 {
			break
		}
		i++
	}

	if b&0x40 > 0 {
		// Sign extend.
		for ; shift < size; shift++ {
			result.SetBit(result, shift, 1)
		}
		result = twosComplementBigInt(result)
		result.Neg(result)
	}
	return result
}

func twosComplementBigInt(n *big.Int) *big.Int {
	absValBytes := n.Bytes()
	for i, b := range absValBytes {
		absValBytes[i] = ^b
	}
	bitsFlipped := big.NewInt(0).SetBytes(absValBytes)
	return bitsFlipped.Add(bitsFlipped, one)
}

func signExtend(n *big.Int, size int) {
	bitPos := size - 7
	max := size
	if bitPos < 0 {
		bitPos = 0
		max = 7
	}
	for ; bitPos < max; bitPos++ {
		n.SetBit(n, bitPos, 1)
	}
}

// equalsNegativeOne is a poor man's check that n, which
// is encoded in two's complement, is all 1's.
func equalsNegativeOne(n *big.Int, size int) bool {
	for i := 0; i < size; i++ {
		if !(n.Bit(i) == 1) {
			return false
		}
	}
	return true
}
