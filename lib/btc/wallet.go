package btc

/*
Power by https://github.com/wenweih/bitcoin_address_protocol
*/

import (
	"crypto/ecdsa"
    "fmt"
	"math/big"
    "golang.org/x/crypto/ripemd160"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
)

const privKeyBytesLen = 32
const addressChecksumLen = 4

// Generate new private key & public key.
func NewPair() ([]byte, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
	    panic(err)
	}
	d := private.D.Bytes()
	b := make([]byte, 0, privKeyBytesLen)
	priKey := paddedAppend(privKeyBytesLen, b, d)
	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return priKey,pubKey
}

func byteString(b []byte) (s string) {
    s = ""
    for i := 0; i < len(b); i++ {
        s += fmt.Sprintf("%02X", b[i])
    }
    return s
}

// b58encode encodes a byte slice b into a base-58 encoded string.
func b58encode(b []byte) (s string) {
    /* See https://en.bitcoin.it/wiki/Base58Check_encoding */

    const BITCOIN_BASE58_TABLE = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

    /* Convert big endian bytes to big int */
    x := new(big.Int).SetBytes(b)

    /* Initialize */
    r := new(big.Int)
    m := big.NewInt(58)
    zero := big.NewInt(0)
    s = ""

    /* Convert big int to string */
    for x.Cmp(zero) > 0 {
        /* x, r = (x / 58, x % 58) */
        x.QuoRem(x, m, r)
        /* Prepend ASCII character */
        s = string(BITCOIN_BASE58_TABLE[r.Int64()]) + s
    }

    return s
}

// b58checkencode encodes version ver and byte slice b into a base-58 check encoded string.
func b58checkencode(ver uint8, b []byte) (s string) {
    /* Prepend version */
    bcpy := append([]byte{ver}, b...)

    /* Create a new SHA256 context */
    sha256H := sha256.New()

    /* SHA256 Hash #1 */
    sha256H.Reset()
    sha256H.Write(bcpy)
    hash1 := sha256H.Sum(nil)

    /* SHA256 Hash #2 */
    sha256H.Reset()
    sha256H.Write(hash1)
    hash2 := sha256H.Sum(nil)

    /* Append first four bytes of hash */
    bcpy = append(bcpy, hash2[0:4]...)

    /* Encode base58 string */
    s = b58encode(bcpy)
        /* For number of leading 0's in bytes, prepend 1 */
    for _, v := range bcpy {
        if v != 0 {
            break
        }
        s = "1" + s
    }

    return s
}

// paddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func paddedAppend(size uint, dst, src []byte) []byte {
    for i := 0; i < int(size)-len(src); i++ {
        dst = append(dst, 0)
    }
    return append(dst, src...)
}

// GetAddress returns wallet address
func GetAddress(pub_bytes []byte) (address string) {
    /* See https://en.bitcoin.it/wiki/Technical_background_of_Bitcoin_addresses */

    /* Convert the public key to bytes */

    /* SHA256 Hash */
    sha256_h := sha256.New()
    sha256_h.Reset()
    sha256_h.Write(pub_bytes)
    pub_hash_1 := sha256_h.Sum(nil)

    /* RIPEMD-160 Hash */
    ripemd160_h := ripemd160.New()
    ripemd160_h.Reset()
    ripemd160_h.Write(pub_hash_1)
    pub_hash_2 := ripemd160_h.Sum(nil)
    /* Convert hash bytes to base58 check encoded sequence */
    address = b58checkencode(0x00, pub_hash_2)

    return address
}
