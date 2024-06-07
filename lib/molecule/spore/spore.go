// Generated by Molecule 0.7.5
// Generated by Moleculec-Go 0.1.11

package molecule

import (
	"bytes"
	"encoding/binary"
	"errors"
	"strconv"
	"strings"
)

type Number uint32

const HeaderSizeUint uint32 = 4

// Byte is the primitive type
type Byte [1]byte

func NewByte(b byte) Byte {
	return Byte([1]byte{b})
}
func ByteDefault() Byte {
	return Byte([1]byte{0})
}
func ByteFromSliceUnchecked(slice []byte) *Byte {
	b := new(Byte)
	b[0] = slice[0]
	return b
}
func (b *Byte) AsSlice() []byte {
	return b[:]
}
func ByteFromSlice(slice []byte, _compatible bool) (*Byte, error) {
	if len(slice) != 1 {
		return nil, errors.New("TotalSizeNotMatch")
	}
	b := new(Byte)
	b[0] = slice[0]
	return b, nil
}
func unpackNumber(b []byte) Number {
	bytesBuffer := bytes.NewBuffer(b)
	var x Number
	binary.Read(bytesBuffer, binary.LittleEndian, &x)
	return x
}
func packNumber(num Number) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, uint32(num))
	return b
}

type BytesBuilder struct {
	inner []Byte
}

func (s *BytesBuilder) Build() Bytes {
	size := packNumber(Number(len(s.inner)))

	b := new(bytes.Buffer)

	b.Write(size)
	len := len(s.inner)
	for i := 0; i < len; i++ {
		b.Write(s.inner[i].AsSlice())
	}

	sb := Bytes{inner: b.Bytes()}

	return sb
}

func (s *BytesBuilder) Set(v []Byte) *BytesBuilder {
	s.inner = v
	return s
}
func (s *BytesBuilder) Push(v Byte) *BytesBuilder {
	s.inner = append(s.inner, v)
	return s
}
func (s *BytesBuilder) Extend(iter []Byte) *BytesBuilder {
	for i := 0; i < len(iter); i++ {
		s.inner = append(s.inner, iter[i])
	}
	return s
}
func (s *BytesBuilder) Replace(index uint, v Byte) *Byte {
	if uint(len(s.inner)) > index {
		a := s.inner[index]
		s.inner[index] = v
		return &a
	}
	return nil
}

func NewBytesBuilder() *BytesBuilder {
	return &BytesBuilder{[]Byte{}}
}

type Bytes struct {
	inner []byte
}

func BytesFromSliceUnchecked(slice []byte) *Bytes {
	return &Bytes{inner: slice}
}
func (s *Bytes) AsSlice() []byte {
	return s.inner
}

func BytesDefault() Bytes {
	return *BytesFromSliceUnchecked([]byte{0, 0, 0, 0})
}

func BytesFromSlice(slice []byte, _compatible bool) (*Bytes, error) {
	sliceLen := len(slice)
	if sliceLen < int(HeaderSizeUint) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "Bytes", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}
	itemCount := unpackNumber(slice)
	if itemCount == 0 {
		if sliceLen != int(HeaderSizeUint) {
			errMsg := strings.Join([]string{"TotalSizeNotMatch", "Bytes", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(HeaderSizeUint))}, " ")
			return nil, errors.New(errMsg)
		}
		return &Bytes{inner: slice}, nil
	}
	totalSize := int(HeaderSizeUint) + int(1*itemCount)
	if sliceLen != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "Bytes", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}
	return &Bytes{inner: slice}, nil
}

func (s *Bytes) TotalSize() uint {
	return uint(HeaderSizeUint) + 1*s.ItemCount()
}
func (s *Bytes) ItemCount() uint {
	number := uint(unpackNumber(s.inner))
	return number
}
func (s *Bytes) Len() uint {
	return s.ItemCount()
}
func (s *Bytes) IsEmpty() bool {
	return s.Len() == 0
}

// if *Byte is nil, index is out of bounds
func (s *Bytes) Get(index uint) *Byte {
	var re *Byte
	if index < s.Len() {
		start := uint(HeaderSizeUint) + 1*index
		end := start + 1
		re = ByteFromSliceUnchecked(s.inner[start:end])
	}
	return re
}

func (s *Bytes) RawData() []byte {
	return s.inner[HeaderSizeUint:]
}

func (s *Bytes) AsBuilder() BytesBuilder {
	size := s.ItemCount()
	t := NewBytesBuilder()
	for i := uint(0); i < size; i++ {
		t.Push(*s.Get(i))
	}
	return *t
}

type BytesOptBuilder struct {
	isNone bool
	inner  Bytes
}

func NewBytesOptBuilder() *BytesOptBuilder {
	return &BytesOptBuilder{isNone: true, inner: BytesDefault()}
}
func (s *BytesOptBuilder) Set(v Bytes) *BytesOptBuilder {
	s.isNone = false
	s.inner = v
	return s
}
func (s *BytesOptBuilder) Build() BytesOpt {
	var ret BytesOpt
	if s.isNone {
		ret = BytesOpt{inner: []byte{}}
	} else {
		ret = BytesOpt{inner: s.inner.AsSlice()}
	}
	return ret
}

type BytesOpt struct {
	inner []byte
}

func BytesOptFromSliceUnchecked(slice []byte) *BytesOpt {
	return &BytesOpt{inner: slice}
}
func (s *BytesOpt) AsSlice() []byte {
	return s.inner
}

func BytesOptDefault() BytesOpt {
	return *BytesOptFromSliceUnchecked([]byte{})
}

func BytesOptFromSlice(slice []byte, compatible bool) (*BytesOpt, error) {
	if len(slice) == 0 {
		return &BytesOpt{inner: slice}, nil
	}

	_, err := BytesFromSlice(slice, compatible)
	if err != nil {
		return nil, err
	}
	return &BytesOpt{inner: slice}, nil
}

func (s *BytesOpt) IntoBytes() (*Bytes, error) {
	if s.IsNone() {
		return nil, errors.New("No data")
	}
	return BytesFromSliceUnchecked(s.AsSlice()), nil
}
func (s *BytesOpt) IsSome() bool {
	return len(s.inner) != 0
}
func (s *BytesOpt) IsNone() bool {
	return len(s.inner) == 0
}
func (s *BytesOpt) AsBuilder() BytesOptBuilder {
	var ret = NewBytesOptBuilder()
	if s.IsSome() {
		ret.Set(*BytesFromSliceUnchecked(s.AsSlice()))
	}
	return *ret
}

type SporeDataBuilder struct {
	content_type Bytes
	content      Bytes
	cluster_id   BytesOpt
}

func (s *SporeDataBuilder) Build() SporeData {
	b := new(bytes.Buffer)

	totalSize := HeaderSizeUint * (3 + 1)
	offsets := make([]uint32, 0, 3)

	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.content_type.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.content.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.cluster_id.AsSlice()))

	b.Write(packNumber(Number(totalSize)))

	for i := 0; i < len(offsets); i++ {
		b.Write(packNumber(Number(offsets[i])))
	}

	b.Write(s.content_type.AsSlice())
	b.Write(s.content.AsSlice())
	b.Write(s.cluster_id.AsSlice())
	return SporeData{inner: b.Bytes()}
}

func (s *SporeDataBuilder) ContentType(v Bytes) *SporeDataBuilder {
	s.content_type = v
	return s
}

func (s *SporeDataBuilder) Content(v Bytes) *SporeDataBuilder {
	s.content = v
	return s
}

func (s *SporeDataBuilder) ClusterId(v BytesOpt) *SporeDataBuilder {
	s.cluster_id = v
	return s
}

func NewSporeDataBuilder() *SporeDataBuilder {
	return &SporeDataBuilder{content_type: BytesDefault(), content: BytesDefault(), cluster_id: BytesOptDefault()}
}

type SporeData struct {
	inner []byte
}

func SporeDataFromSliceUnchecked(slice []byte) *SporeData {
	return &SporeData{inner: slice}
}
func (s *SporeData) AsSlice() []byte {
	return s.inner
}

func SporeDataDefault() SporeData {
	return *SporeDataFromSliceUnchecked([]byte{24, 0, 0, 0, 16, 0, 0, 0, 20, 0, 0, 0, 24, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func SporeDataFromSlice(slice []byte, compatible bool) (*SporeData, error) {
	sliceLen := len(slice)
	if uint32(sliceLen) < HeaderSizeUint {
		errMsg := strings.Join([]string{"HeaderIsBroken", "SporeData", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}

	totalSize := unpackNumber(slice)
	if Number(sliceLen) != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "SporeData", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}

	if uint32(sliceLen) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "SporeData", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	offsetFirst := unpackNumber(slice[HeaderSizeUint:])
	if uint32(offsetFirst)%HeaderSizeUint != 0 || uint32(offsetFirst) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"OffsetsNotMatch", "SporeData", strconv.Itoa(int(offsetFirst % 4)), "!= 0", strconv.Itoa(int(offsetFirst)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	if sliceLen < int(offsetFirst) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "SporeData", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(offsetFirst))}, " ")
		return nil, errors.New(errMsg)
	}

	fieldCount := uint32(offsetFirst)/HeaderSizeUint - 1
	if fieldCount < 3 {
		return nil, errors.New("FieldCountNotMatch")
	} else if !compatible && fieldCount > 3 {
		return nil, errors.New("FieldCountNotMatch")
	}

	offsets := make([]uint32, fieldCount)

	for i := 0; i < int(fieldCount); i++ {
		offsets[i] = uint32(unpackNumber(slice[HeaderSizeUint:][int(HeaderSizeUint)*i:]))
	}
	offsets = append(offsets, uint32(totalSize))

	for i := 0; i < len(offsets); i++ {
		if i&1 != 0 && offsets[i-1] > offsets[i] {
			return nil, errors.New("OffsetsNotMatch")
		}
	}

	var err error

	_, err = BytesFromSlice(slice[offsets[0]:offsets[1]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[1]:offsets[2]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesOptFromSlice(slice[offsets[2]:offsets[3]], compatible)
	if err != nil {
		return nil, err
	}

	return &SporeData{inner: slice}, nil
}

func (s *SporeData) TotalSize() uint {
	return uint(unpackNumber(s.inner))
}
func (s *SporeData) FieldCount() uint {
	var number uint = 0
	if uint32(s.TotalSize()) == HeaderSizeUint {
		return number
	}
	number = uint(unpackNumber(s.inner[HeaderSizeUint:]))/4 - 1
	return number
}
func (s *SporeData) Len() uint {
	return s.FieldCount()
}
func (s *SporeData) IsEmpty() bool {
	return s.Len() == 0
}
func (s *SporeData) CountExtraFields() uint {
	return s.FieldCount() - 3
}

func (s *SporeData) HasExtraFields() bool {
	return 3 != s.FieldCount()
}

func (s *SporeData) ContentType() *Bytes {
	start := unpackNumber(s.inner[4:])
	end := unpackNumber(s.inner[8:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *SporeData) Content() *Bytes {
	start := unpackNumber(s.inner[8:])
	end := unpackNumber(s.inner[12:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *SporeData) ClusterId() *BytesOpt {
	var ret *BytesOpt
	start := unpackNumber(s.inner[12:])
	if s.HasExtraFields() {
		end := unpackNumber(s.inner[16:])
		ret = BytesOptFromSliceUnchecked(s.inner[start:end])
	} else {
		ret = BytesOptFromSliceUnchecked(s.inner[start:])
	}
	return ret
}

func (s *SporeData) AsBuilder() SporeDataBuilder {
	ret := NewSporeDataBuilder().ContentType(*s.ContentType()).Content(*s.Content()).ClusterId(*s.ClusterId())
	return *ret
}

type ClusterDataBuilder struct {
	name        Bytes
	description Bytes
}

func (s *ClusterDataBuilder) Build() ClusterData {
	b := new(bytes.Buffer)

	totalSize := HeaderSizeUint * (2 + 1)
	offsets := make([]uint32, 0, 2)

	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.name.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.description.AsSlice()))

	b.Write(packNumber(Number(totalSize)))

	for i := 0; i < len(offsets); i++ {
		b.Write(packNumber(Number(offsets[i])))
	}

	b.Write(s.name.AsSlice())
	b.Write(s.description.AsSlice())
	return ClusterData{inner: b.Bytes()}
}

func (s *ClusterDataBuilder) Name(v Bytes) *ClusterDataBuilder {
	s.name = v
	return s
}

func (s *ClusterDataBuilder) Description(v Bytes) *ClusterDataBuilder {
	s.description = v
	return s
}

func NewClusterDataBuilder() *ClusterDataBuilder {
	return &ClusterDataBuilder{name: BytesDefault(), description: BytesDefault()}
}

type ClusterData struct {
	inner []byte
}

func ClusterDataFromSliceUnchecked(slice []byte) *ClusterData {
	return &ClusterData{inner: slice}
}
func (s *ClusterData) AsSlice() []byte {
	return s.inner
}

func ClusterDataDefault() ClusterData {
	return *ClusterDataFromSliceUnchecked([]byte{20, 0, 0, 0, 12, 0, 0, 0, 16, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func ClusterDataFromSlice(slice []byte, compatible bool) (*ClusterData, error) {
	sliceLen := len(slice)
	if uint32(sliceLen) < HeaderSizeUint {
		errMsg := strings.Join([]string{"HeaderIsBroken", "ClusterData", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}

	totalSize := unpackNumber(slice)
	if Number(sliceLen) != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "ClusterData", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}

	if uint32(sliceLen) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "ClusterData", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	offsetFirst := unpackNumber(slice[HeaderSizeUint:])
	if uint32(offsetFirst)%HeaderSizeUint != 0 || uint32(offsetFirst) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"OffsetsNotMatch", "ClusterData", strconv.Itoa(int(offsetFirst % 4)), "!= 0", strconv.Itoa(int(offsetFirst)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	if sliceLen < int(offsetFirst) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "ClusterData", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(offsetFirst))}, " ")
		return nil, errors.New(errMsg)
	}

	fieldCount := uint32(offsetFirst)/HeaderSizeUint - 1
	if fieldCount < 2 {
		return nil, errors.New("FieldCountNotMatch")
	} else if !compatible && fieldCount > 2 {
		return nil, errors.New("FieldCountNotMatch")
	}

	offsets := make([]uint32, fieldCount)

	for i := 0; i < int(fieldCount); i++ {
		offsets[i] = uint32(unpackNumber(slice[HeaderSizeUint:][int(HeaderSizeUint)*i:]))
	}
	offsets = append(offsets, uint32(totalSize))

	for i := 0; i < len(offsets); i++ {
		if i&1 != 0 && offsets[i-1] > offsets[i] {
			return nil, errors.New("OffsetsNotMatch")
		}
	}

	var err error

	_, err = BytesFromSlice(slice[offsets[0]:offsets[1]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[1]:offsets[2]], compatible)
	if err != nil {
		return nil, err
	}

	return &ClusterData{inner: slice}, nil
}

func (s *ClusterData) TotalSize() uint {
	return uint(unpackNumber(s.inner))
}
func (s *ClusterData) FieldCount() uint {
	var number uint = 0
	if uint32(s.TotalSize()) == HeaderSizeUint {
		return number
	}
	number = uint(unpackNumber(s.inner[HeaderSizeUint:]))/4 - 1
	return number
}
func (s *ClusterData) Len() uint {
	return s.FieldCount()
}
func (s *ClusterData) IsEmpty() bool {
	return s.Len() == 0
}
func (s *ClusterData) CountExtraFields() uint {
	return s.FieldCount() - 2
}

func (s *ClusterData) HasExtraFields() bool {
	return 2 != s.FieldCount()
}

func (s *ClusterData) Name() *Bytes {
	start := unpackNumber(s.inner[4:])
	end := unpackNumber(s.inner[8:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *ClusterData) Description() *Bytes {
	var ret *Bytes
	start := unpackNumber(s.inner[8:])
	if s.HasExtraFields() {
		end := unpackNumber(s.inner[12:])
		ret = BytesFromSliceUnchecked(s.inner[start:end])
	} else {
		ret = BytesFromSliceUnchecked(s.inner[start:])
	}
	return ret
}

func (s *ClusterData) AsBuilder() ClusterDataBuilder {
	ret := NewClusterDataBuilder().Name(*s.Name()).Description(*s.Description())
	return *ret
}