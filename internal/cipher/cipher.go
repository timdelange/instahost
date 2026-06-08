package cipher

import "encoding/binary"

const DefaultKey = "instahost"
const ChecksumBytes = 4

var crc32Table [256]uint32

func init() {
	for i := 0; i < 256; i++ {
		c := uint32(i)
		for j := 0; j < 8; j++ {
			if c&1 != 0 {
				c = 0xedb88320 ^ (c >> 1)
			} else {
				c >>= 1
			}
		}
		crc32Table[i] = c
	}
}

func CRC32(data []byte) uint32 {
	crc := uint32(0xffffffff)
	for _, b := range data {
		crc = crc32Table[(crc^uint32(b))&0xff] ^ (crc >> 8)
	}
	return (crc ^ 0xffffffff)
}

func SealWithChecksum(data []byte) []byte {
	out := make([]byte, len(data)+ChecksumBytes)
	copy(out, data)
	binary.BigEndian.PutUint32(out[len(data):], CRC32(data))
	return out
}

func OpenWithChecksum(data []byte) ([]byte, bool) {
	if len(data) <= ChecksumBytes {
		return nil, false
	}
	payload := data[:len(data)-ChecksumBytes]
	stored := binary.BigEndian.Uint32(data[len(data)-ChecksumBytes:])
	if stored != CRC32(payload) {
		return nil, false
	}
	return payload, true
}

func XorBytes(data []byte, key string) []byte {
	if key == "" {
		key = DefaultKey
	}
	keyBytes := []byte(key)
	out := make([]byte, len(data))
	for i, b := range data {
		out[i] = b ^ keyBytes[i%len(keyBytes)]
	}
	return out
}
