package libcoin

import (
	"encoding/binary"
)

// MarshalNBO serializes an Amount into the 24-byte TALER_AmountNBO format:
//
//	[uint64 value big-endian] [uint32 fraction big-endian] [char[12] currency zero-padded]
func (receiver Amount) MarshalNBO() [24]byte {
	var buf [24]byte
	binary.BigEndian.PutUint64(buf[0:8], uint64(receiver.Value))
	binary.BigEndian.PutUint32(buf[8:12], uint32(receiver.Fraction))
	copy(buf[12:24], receiver.Currency)
	return buf
}
