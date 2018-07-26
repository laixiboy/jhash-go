package main 

import ( 
	"fmt"
	"encoding/binary"
	"strconv"
)

// An arbitrary initial parameter
const JHASH_INITVAL uint32 = 0xdeadbeef;

func Rol32(word uint32, shift uint) uint32 {
	return (word << shift) | (word >> ((-shift) & 31))
}

func JHashMix(a *uint32, b *uint32, c *uint32) {
	*a -= *c;  *a ^= Rol32(*c, 4);  *c += *b;	
	*b -= *a;  *b ^= Rol32(*a, 6);  *a += *c;
	*c -= *b;  *c ^= Rol32(*b, 8);  *b += *a;
	*a -= *c;  *a ^= Rol32(*c, 16); *c += *b;
	*b -= *a;  *b ^= Rol32(*a, 19); *a += *c;
	*c -= *b;  *c ^= Rol32(*b, 4);  *b += *a;
}

func JHashFinal(a *uint32, b *uint32, c *uint32) {
	*c ^= *b; *c -= Rol32(*b, 14);
	*a ^= *c; *a -= Rol32(*c, 11);
	*b ^= *a; *b -= Rol32(*a, 25);
	*c ^= *b; *c -= Rol32(*b, 16);
	*a ^= *c; *a -= Rol32(*c, 4);
	*b ^= *a; *b -= Rol32(*a, 14);
	*c ^= *b; *c -= Rol32(*b, 24);
}

/* jhash - hash an arbitrary key
 * @k: sequence of bytes as key
 * @length: the length of the key
 * @initval: the previous hash, or an arbitray value
 *
 * The generic version, hashes an arbitrary sequence of bytes.
 * No alignment or length assumptions are made about the input key.
 *
 * Returns the hash value of the key. The result depends on endianness.
 */
func JHash(key []byte, length uint32, initval uint32) uint32 {
	var a, b, c uint32
	k := []uint8(key)
	
	// Set up the internal state
	i := JHASH_INITVAL + length + initval
	a, b, c = i, i, i

	// All but the last block: affect some 32 bits of (a,b,c)	
	for {
		if length > 12 {
			a += binary.LittleEndian.Uint32(k)
			b += binary.LittleEndian.Uint32(k[4:8])
			c += binary.LittleEndian.Uint32(k[8:12])
			JHashMix(&a, &b, &c)
			length -= 12
			k = k[12:]
		} else {
			break
		}
	}

	// Last block: affect all 32 bits of (c) 
	// All the case statements fall through 
	switch (length) {
		case 12: c += uint32(k[11])<<24
			fallthrough
		case 11: c += uint32(k[10])<<16
			fallthrough
		case 10: c += uint32(k[9])<<8
			fallthrough
		case 9:  c += uint32(k[8])
			fallthrough
		case 8:  b += uint32(k[7])<<24
			fallthrough
		case 7:  b += uint32(k[6])<<16
			fallthrough
		case 6:  b += uint32(k[5])<<8
			fallthrough
		case 5:  b += uint32(k[4])
			fallthrough
		case 4:  a += uint32(k[3])<<24
			fallthrough
		case 3:  a += uint32(k[2])<<16
			fallthrough
		case 2:  a += uint32(k[1])<<8
			fallthrough
		case 1:  a += uint32(k[0])
			JHashFinal(&a, &b, &c)
			fallthrough
		case 0: // Nothing left to add
			break
	}
	
	return c 
}

func Uint32ToBytesLittleEndian(val uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, val)
	return b
}

func main() {
	var hash uint32 = 0
	const bucket_size uint32 = 16;
	var buckets [bucket_size]uint32
	fmt.Println("---------test number hash---------");
	for i := uint32(100000000); i < 100010000; i++ {
		hash = JHash(Uint32ToBytesLittleEndian(i), 4, hash)
		buckets[hash % bucket_size]++;
	}

	for k, v := range buckets {
		fmt.Println(v);
		buckets[k] = 0;
	}
	hash = 0
	
	fmt.Println("---------test string hash---------");
	prefix := "stringprefixabcdefghijk"
	for i := 100000000; i < 100010000; i++ {
		comb := prefix + strconv.Itoa(i)
		hash = JHash([]byte(comb), uint32(len(comb)), hash)
		buckets[hash % bucket_size]++;
	}

	for _, val := range buckets {
		fmt.Println(val);
	}
}
