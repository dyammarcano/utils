package ulid

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"io"
	"math/rand"
	"sync"
	"time"
)

// MonotonicReader is an interface that should yield monotonically increasing
// entropy into the provided slice for all calls with the same ms parameter. If
// a MonotonicReader is provided to the New constructor, its MonotonicRead
// method will be used instead of Read.
type MonotonicReader interface {
	io.Reader
	MonotonicRead(ms uint64, p []byte) error
}

// New returns an ULID with the given Unix milliseconds timestamp and an
// optional entropy source. Use the Timestamp function to convert
// a time.Time to Unix milliseconds.
//
// ErrBigTime is returned when passing a timestamp bigger than MaxTime.
// Reading from the entropy source may also return an error.
//
// Safety for concurrent use is only dependent on the safety of the
// entropy source.
func New(ms uint64, entropy io.Reader) (*ULID, error) {
	id := &ULID{}
	var err error
	if err = id.SetTime(ms); err != nil {
		return id, err
	}

	switch e := entropy.(type) {
	case nil:
		return id, err
	case MonotonicReader:
		err = e.MonotonicRead(ms, id[6:])
	default:
		_, err = io.ReadFull(e, id[6:])
	}

	return id, err
}

// MustNew is a convenience function equivalent to New that panics on failure
// instead of returning an error.
func MustNew(ms uint64, entropy io.Reader) *ULID {
	id, err := New(ms, entropy)
	if err != nil {
		panic(err)
	}
	return id
}

// DefaultEntropy returns a thread-safe per process monotonically increasing
// entropy source.
func DefaultEntropy() io.Reader {
	entropyOnce.Do(func() {
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		entropy = &LockedMonotonicReader{
			MonotonicReader: Monotonic(rng, 0),
		}
	})
	return entropy
}

// Make returns an ULID with the current time in Unix milliseconds and
// monotonically increasing entropy for the same millisecond.
// It is safe for concurrent use, leveraging a sync.Pool underneath for minimal
// contention.
func Make() (id *ULID) {
	// NOTE: MustNew can't panic since DefaultEntropy never returns an error.
	return MustNew(Now(), DefaultEntropy())
}

// Parse parses an encoded ULID, returning an error in case of failure.
//
// ErrDataSize is returned if the len(ulid) is different from an encoded
// ULID's length. Invalid encodings produce undefined ULIDs. For a version that
// returns an error instead, see ParseStrict.
func Parse(ulid string) (id *ULID, err error) {
	return id, parse([]byte(ulid), false, id)
}

// ParseStrict parses an encoded ULID, returning an error in case of failure.
//
// It is like Parse, but additionally validates that the parsed ULID consists
// only of valid base32 characters. It is slightly slower than Parse.
//
// ErrDataSize is returned if the len(ulid) is different from an encoded
// ULID's length. Invalid encodings return ErrInvalidCharacters.
func ParseStrict(ulid string) (id ULID, err error) {
	return id, parse([]byte(ulid), true, &id)
}

func parse(v []byte, strict bool, id *ULID) error {
	// Check if a base32 encoded ULID is the right length.
	if len(v) != EncodedSize {
		return ErrDataSize
	}

	// Check if all the characters in a base32 encoded ULID are part of the
	// expected base32 character set.
	if strict &&
		(dec[v[0]] == 0xFF ||
			dec[v[1]] == 0xFF ||
			dec[v[2]] == 0xFF ||
			dec[v[3]] == 0xFF ||
			dec[v[4]] == 0xFF ||
			dec[v[5]] == 0xFF ||
			dec[v[6]] == 0xFF ||
			dec[v[7]] == 0xFF ||
			dec[v[8]] == 0xFF ||
			dec[v[9]] == 0xFF ||
			dec[v[10]] == 0xFF ||
			dec[v[11]] == 0xFF ||
			dec[v[12]] == 0xFF ||
			dec[v[13]] == 0xFF ||
			dec[v[14]] == 0xFF ||
			dec[v[15]] == 0xFF ||
			dec[v[16]] == 0xFF ||
			dec[v[17]] == 0xFF ||
			dec[v[18]] == 0xFF ||
			dec[v[19]] == 0xFF ||
			dec[v[20]] == 0xFF ||
			dec[v[21]] == 0xFF ||
			dec[v[22]] == 0xFF ||
			dec[v[23]] == 0xFF ||
			dec[v[24]] == 0xFF ||
			dec[v[25]] == 0xFF) {
		return ErrInvalidCharacters
	}

	// Check if the first character in a base32 encoded ULID will overflow. This
	// happens because the base32 representation encodes 130 bits, while the
	// ULID is only 128 bits.
	//
	// See https://github.com/oklog/ulid/issues/9 for details.
	if v[0] > '7' {
		return ErrOverflow
	}

	// Use an optimized unrolled loop (from https://github.com/RobThree/NUlid)
	// to decode a base32 ULID.

	// 6 bytes timestamp (48 bits)
	(*id)[0] = (dec[v[0]] << 5) | dec[v[1]]
	(*id)[1] = (dec[v[2]] << 3) | (dec[v[3]] >> 2)
	(*id)[2] = (dec[v[3]] << 6) | (dec[v[4]] << 1) | (dec[v[5]] >> 4)
	(*id)[3] = (dec[v[5]] << 4) | (dec[v[6]] >> 1)
	(*id)[4] = (dec[v[6]] << 7) | (dec[v[7]] << 2) | (dec[v[8]] >> 3)
	(*id)[5] = (dec[v[8]] << 5) | dec[v[9]]

	// 10 bytes of entropy (80 bits)
	(*id)[6] = (dec[v[10]] << 3) | (dec[v[11]] >> 2)
	(*id)[7] = (dec[v[11]] << 6) | (dec[v[12]] << 1) | (dec[v[13]] >> 4)
	(*id)[8] = (dec[v[13]] << 4) | (dec[v[14]] >> 1)
	(*id)[9] = (dec[v[14]] << 7) | (dec[v[15]] << 2) | (dec[v[16]] >> 3)
	(*id)[10] = (dec[v[16]] << 5) | dec[v[17]]
	(*id)[11] = (dec[v[18]] << 3) | dec[v[19]]>>2
	(*id)[12] = (dec[v[19]] << 6) | (dec[v[20]] << 1) | (dec[v[21]] >> 4)
	(*id)[13] = (dec[v[21]] << 4) | (dec[v[22]] >> 1)
	(*id)[14] = (dec[v[22]] << 7) | (dec[v[23]] << 2) | (dec[v[24]] >> 3)
	(*id)[15] = (dec[v[24]] << 5) | dec[v[25]]

	return nil
}

// MustParse is a convenience function equivalent to Parse that panics on failure
// instead of returning an error.
func MustParse(ulid string) *ULID {
	id, err := Parse(ulid)
	if err != nil {
		panic(err)
	}
	return id
}

// MustParseStrict is a convenience function equivalent to ParseStrict that
// panics on failure instead of returning an error.
func MustParseStrict(ulid string) *ULID {
	id, err := ParseStrict(ulid)
	if err != nil {
		panic(err)
	}
	return &id
}

// Bytes returns bytes slice representation of ULID.
func (id *ULID) Bytes() []byte {
	return id[:]
}

// String returns a lexicographically sortable string encoded ULID
// (26 characters, non-standard base 32) e.g. 01AN4Z07BY79KA1307SR9X4MV3.
// Format: tttttttttteeeeeeeeeeeeeeee where t is time and e is entropy.
func (id *ULID) String() string {
	ulid := make([]byte, EncodedSize)
	_ = id.MarshalTextTo(ulid)
	return string(ulid)
}

// MarshalBinary implements the encoding.BinaryMarshaler interface by
// returning the ULID as a byte slice.
func (id *ULID) MarshalBinary() ([]byte, error) {
	ulid := make([]byte, len(id))
	return ulid, id.MarshalBinaryTo(ulid)
}

// MarshalBinaryTo writes the binary encoding of the ULID to the given buffer.
// ErrBufferSize is returned when the len(dst) != 16.
func (id *ULID) MarshalBinaryTo(dst []byte) error {
	if len(dst) != len(id) {
		return ErrBufferSize
	}

	copy(dst, id[:])
	return nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface by
// copying the passed data and converting it to an ULID. ErrDataSize is
// returned if the data length is different from ULID length.
func (id *ULID) UnmarshalBinary(data []byte) error {
	if len(data) != len(*id) {
		return ErrDataSize
	}

	copy((*id)[:], data)
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface by
// returning the string encoded ULID.
func (id *ULID) MarshalText() ([]byte, error) {
	ulid := make([]byte, EncodedSize)
	return ulid, id.MarshalTextTo(ulid)
}

// MarshalTextTo writes the ULID as a string to the given buffer.
// ErrBufferSize is returned when the len(dst) != 26.
func (id *ULID) MarshalTextTo(dst []byte) error {
	// Optimized unrolled loop ahead.
	// From https://github.com/RobThree/NUlid

	if len(dst) != EncodedSize {
		return ErrBufferSize
	}

	// 10 byte timestamp
	dst[0] = Encoding[(id[0]&224)>>5]
	dst[1] = Encoding[id[0]&31]
	dst[2] = Encoding[(id[1]&248)>>3]
	dst[3] = Encoding[((id[1]&7)<<2)|((id[2]&192)>>6)]
	dst[4] = Encoding[(id[2]&62)>>1]
	dst[5] = Encoding[((id[2]&1)<<4)|((id[3]&240)>>4)]
	dst[6] = Encoding[((id[3]&15)<<1)|((id[4]&128)>>7)]
	dst[7] = Encoding[(id[4]&124)>>2]
	dst[8] = Encoding[((id[4]&3)<<3)|((id[5]&224)>>5)]
	dst[9] = Encoding[id[5]&31]

	// 16 bytes of entropy
	dst[10] = Encoding[(id[6]&248)>>3]
	dst[11] = Encoding[((id[6]&7)<<2)|((id[7]&192)>>6)]
	dst[12] = Encoding[(id[7]&62)>>1]
	dst[13] = Encoding[((id[7]&1)<<4)|((id[8]&240)>>4)]
	dst[14] = Encoding[((id[8]&15)<<1)|((id[9]&128)>>7)]
	dst[15] = Encoding[(id[9]&124)>>2]
	dst[16] = Encoding[((id[9]&3)<<3)|((id[10]&224)>>5)]
	dst[17] = Encoding[id[10]&31]
	dst[18] = Encoding[(id[11]&248)>>3]
	dst[19] = Encoding[((id[11]&7)<<2)|((id[12]&192)>>6)]
	dst[20] = Encoding[(id[12]&62)>>1]
	dst[21] = Encoding[((id[12]&1)<<4)|((id[13]&240)>>4)]
	dst[22] = Encoding[((id[13]&15)<<1)|((id[14]&128)>>7)]
	dst[23] = Encoding[(id[14]&124)>>2]
	dst[24] = Encoding[((id[14]&3)<<3)|((id[15]&224)>>5)]
	dst[25] = Encoding[id[15]&31]

	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface by
// parsing the data as string encoded ULID.
//
// ErrDataSize is returned if the len(v) is different from an encoded
// ULID's length. Invalid encodings produce undefined ULIDs.
func (id *ULID) UnmarshalText(v []byte) error {
	return parse(v, false, id)
}

// Time returns the Unix time in milliseconds encoded in the ULID.
// Use the top level Time function to convert the returned value to
// a time.Time.
func (id *ULID) Time() uint64 {
	return uint64(id[5]) | uint64(id[4])<<8 |
		uint64(id[3])<<16 | uint64(id[2])<<24 |
		uint64(id[1])<<32 | uint64(id[0])<<40
}

// MaxTime returns the maximum Unix time in milliseconds that
// can be encoded in an ULID.
func MaxTime() uint64 { return maxTime }

// Now is a convenience function that returns the current
// UTC time in Unix milliseconds. Equivalent to:
//
//	Timestamp(time.Now().UTC())
func Now() uint64 { return Timestamp(time.Now().UTC()) }

// Timestamp converts a time.Time to Unix milliseconds.
//
// Because of the way ULID stores time, times from the year
// 10889 produces undefined results.
func Timestamp(t time.Time) uint64 {
	return uint64(t.Unix())*1000 +
		uint64(t.Nanosecond()/int(time.Millisecond))
}

// Time converts Unix milliseconds in the format
// returned by the Timestamp function to a time.Time.
func Time(ms uint64) time.Time {
	s := int64(ms / 1e3)
	ns := int64((ms % 1e3) * 1e6)
	return time.Unix(s, ns)
}

// SetTime sets the time component of the ULID to the given Unix time
// in milliseconds.
func (id *ULID) SetTime(ms uint64) error {
	if ms > maxTime {
		return ErrBigTime
	}

	(*id)[0] = byte(ms >> 40)
	(*id)[1] = byte(ms >> 32)
	(*id)[2] = byte(ms >> 24)
	(*id)[3] = byte(ms >> 16)
	(*id)[4] = byte(ms >> 8)
	(*id)[5] = byte(ms)

	return nil
}

// Entropy returns the entropy from the ULID.
func (id *ULID) Entropy() []byte {
	e := make([]byte, 10)
	copy(e, id[6:])
	return e
}

// SetEntropy sets the ULID entropy to the passed byte slice.
// ErrDataSize is returned if len(e) != 10.
func (id *ULID) SetEntropy(e []byte) error {
	if len(e) != 10 {
		return ErrDataSize
	}

	copy((*id)[6:], e)
	return nil
}

// Compare returns an integer comparing id and other lexicographically.
// The result will be 0 if id==other, -1 if id < other, and +1 if id > other.
func (id *ULID) Compare(*ULID) int {
	other := &ULID{}
	return bytes.Compare(id[:], other[:])
}

// Scan implements the sql.Scanner interface. It supports scanning
// a string or byte slice.
func (id *ULID) Scan(src interface{}) error {
	switch x := src.(type) {
	case nil:
		return nil
	case string:
		return id.UnmarshalText([]byte(x))
	case []byte:
		return id.UnmarshalBinary(x)
	}

	return ErrScanValue
}

// Value implements the sql/driver.Valuer interface, returning the ULID as a
// slice of bytes, by invoking MarshalBinary. If your use case requires a string
// representation instead, you can create a wrapper type that calls String()
// instead.
//
//	type stringValuer ulid.ULID
//
//	func (v stringValuer) Value() (driver.Value, error) {
//	    return ulid.ULID(v).String(), nil
//	}
//
//	// Example usage.
//	db.Exec("...", stringValuer(id))
//
// All valid ULIDs, including zero-value ULIDs, return a valid Value with a nil
// error. If your use case requires zero-value ULIDs to return a non-nil error,
// you can create a wrapper type that special-cases this behavior.
//
//	var zeroValueULID ulid.ULID
//
//	type invalidZeroValuer ulid.ULID
//
//	func (v invalidZeroValuer) Value() (driver.Value, error) {
//	    if ulid.ULID(v).Compare(zeroValueULID) == 0 {
//	        return nil, fmt.Errorf("zero value")
//	    }
//	    return ulid.ULID(v).Value()
//	}
//
//	// Example usage.
//	db.Exec("...", invalidZeroValuer(id))
func (id *ULID) Value() (driver.Value, error) {
	return id.MarshalBinary()
}

type rng interface{ Int63n(n int64) int64 }

// LockedMonotonicReader wraps a MonotonicReader with a sync.Mutex for
// safe concurrent use.
type LockedMonotonicReader struct {
	mu sync.Mutex
	MonotonicReader
}

// MonotonicRead synchronizes calls to the wrapped MonotonicReader.
func (r *LockedMonotonicReader) MonotonicRead(ms uint64, p []byte) (err error) {
	r.mu.Lock()
	err = r.MonotonicReader.MonotonicRead(ms, p)
	r.mu.Unlock()
	return err
}

type uint80 struct {
	Hi uint16
	Lo uint64
}

func (u *uint80) SetBytes(bs []byte) {
	u.Hi = binary.BigEndian.Uint16(bs[:2])
	u.Lo = binary.BigEndian.Uint64(bs[2:])
}

func (u *uint80) AppendTo(bs []byte) {
	binary.BigEndian.PutUint16(bs[:2], u.Hi)
	binary.BigEndian.PutUint64(bs[2:], u.Lo)
}

func (u *uint80) Add(n uint64) (overflow bool) {
	lo, hi := u.Lo, u.Hi
	if u.Lo += n; u.Lo < lo {
		u.Hi++
	}
	return u.Hi < hi
}

func (u *uint80) IsZero() bool {
	return u.Hi == 0 && u.Lo == 0
}
