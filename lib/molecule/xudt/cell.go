// Generated by Molecule 0.7.5
// Generated by Moleculec-Go 0.1.11

package xudt

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

type GovernanceMembersBuilder struct {
	parent_id     Bytes
	lock_args     Bytes
	multisig_args Bytes
	members       BytesVec
}

func (s *GovernanceMembersBuilder) Build() GovernanceMembers {
	b := new(bytes.Buffer)

	totalSize := HeaderSizeUint * (4 + 1)
	offsets := make([]uint32, 0, 4)

	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.parent_id.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.lock_args.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.multisig_args.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.members.AsSlice()))

	b.Write(packNumber(Number(totalSize)))

	for i := 0; i < len(offsets); i++ {
		b.Write(packNumber(Number(offsets[i])))
	}

	b.Write(s.parent_id.AsSlice())
	b.Write(s.lock_args.AsSlice())
	b.Write(s.multisig_args.AsSlice())
	b.Write(s.members.AsSlice())
	return GovernanceMembers{inner: b.Bytes()}
}

func (s *GovernanceMembersBuilder) ParentId(v Bytes) *GovernanceMembersBuilder {
	s.parent_id = v
	return s
}

func (s *GovernanceMembersBuilder) LockArgs(v Bytes) *GovernanceMembersBuilder {
	s.lock_args = v
	return s
}

func (s *GovernanceMembersBuilder) MultisigArgs(v Bytes) *GovernanceMembersBuilder {
	s.multisig_args = v
	return s
}

func (s *GovernanceMembersBuilder) Members(v BytesVec) *GovernanceMembersBuilder {
	s.members = v
	return s
}

func NewGovernanceMembersBuilder() *GovernanceMembersBuilder {
	return &GovernanceMembersBuilder{parent_id: BytesDefault(), lock_args: BytesDefault(), multisig_args: BytesDefault(), members: BytesVecDefault()}
}

type GovernanceMembers struct {
	inner []byte
}

func GovernanceMembersFromSliceUnchecked(slice []byte) *GovernanceMembers {
	return &GovernanceMembers{inner: slice}
}
func (s *GovernanceMembers) AsSlice() []byte {
	return s.inner
}

func GovernanceMembersDefault() GovernanceMembers {
	return *GovernanceMembersFromSliceUnchecked([]byte{36, 0, 0, 0, 20, 0, 0, 0, 24, 0, 0, 0, 28, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0})
}

func GovernanceMembersFromSlice(slice []byte, compatible bool) (*GovernanceMembers, error) {
	sliceLen := len(slice)
	if uint32(sliceLen) < HeaderSizeUint {
		errMsg := strings.Join([]string{"HeaderIsBroken", "GovernanceMembers", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}

	totalSize := unpackNumber(slice)
	if Number(sliceLen) != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "GovernanceMembers", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}

	if uint32(sliceLen) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "GovernanceMembers", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	offsetFirst := unpackNumber(slice[HeaderSizeUint:])
	if uint32(offsetFirst)%HeaderSizeUint != 0 || uint32(offsetFirst) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"OffsetsNotMatch", "GovernanceMembers", strconv.Itoa(int(offsetFirst % 4)), "!= 0", strconv.Itoa(int(offsetFirst)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	if sliceLen < int(offsetFirst) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "GovernanceMembers", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(offsetFirst))}, " ")
		return nil, errors.New(errMsg)
	}

	fieldCount := uint32(offsetFirst)/HeaderSizeUint - 1
	if fieldCount < 4 {
		return nil, errors.New("FieldCountNotMatch")
	} else if !compatible && fieldCount > 4 {
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

	_, err = BytesFromSlice(slice[offsets[2]:offsets[3]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesVecFromSlice(slice[offsets[3]:offsets[4]], compatible)
	if err != nil {
		return nil, err
	}

	return &GovernanceMembers{inner: slice}, nil
}

func (s *GovernanceMembers) TotalSize() uint {
	return uint(unpackNumber(s.inner))
}
func (s *GovernanceMembers) FieldCount() uint {
	var number uint = 0
	if uint32(s.TotalSize()) == HeaderSizeUint {
		return number
	}
	number = uint(unpackNumber(s.inner[HeaderSizeUint:]))/4 - 1
	return number
}
func (s *GovernanceMembers) Len() uint {
	return s.FieldCount()
}
func (s *GovernanceMembers) IsEmpty() bool {
	return s.Len() == 0
}
func (s *GovernanceMembers) CountExtraFields() uint {
	return s.FieldCount() - 4
}

func (s *GovernanceMembers) HasExtraFields() bool {
	return 4 != s.FieldCount()
}

func (s *GovernanceMembers) ParentId() *Bytes {
	start := unpackNumber(s.inner[4:])
	end := unpackNumber(s.inner[8:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *GovernanceMembers) LockArgs() *Bytes {
	start := unpackNumber(s.inner[8:])
	end := unpackNumber(s.inner[12:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *GovernanceMembers) MultisigArgs() *Bytes {
	start := unpackNumber(s.inner[12:])
	end := unpackNumber(s.inner[16:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *GovernanceMembers) Members() *BytesVec {
	var ret *BytesVec
	start := unpackNumber(s.inner[16:])
	if s.HasExtraFields() {
		end := unpackNumber(s.inner[20:])
		ret = BytesVecFromSliceUnchecked(s.inner[start:end])
	} else {
		ret = BytesVecFromSliceUnchecked(s.inner[start:])
	}
	return ret
}

func (s *GovernanceMembers) AsBuilder() GovernanceMembersBuilder {
	ret := NewGovernanceMembersBuilder().ParentId(*s.ParentId()).LockArgs(*s.LockArgs()).MultisigArgs(*s.MultisigArgs()).Members(*s.Members())
	return *ret
}

type TickBuilder struct {
	tick_type    Byte
	token_id     Bytes
	value        Uint128
	merchant     Script
	coin_type    Bytes
	tx_hash      Bytes
	receipt_addr Bytes
}

func (s *TickBuilder) Build() Tick {
	b := new(bytes.Buffer)

	totalSize := HeaderSizeUint * (7 + 1)
	offsets := make([]uint32, 0, 7)

	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.tick_type.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.token_id.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.value.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.merchant.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.coin_type.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.tx_hash.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.receipt_addr.AsSlice()))

	b.Write(packNumber(Number(totalSize)))

	for i := 0; i < len(offsets); i++ {
		b.Write(packNumber(Number(offsets[i])))
	}

	b.Write(s.tick_type.AsSlice())
	b.Write(s.token_id.AsSlice())
	b.Write(s.value.AsSlice())
	b.Write(s.merchant.AsSlice())
	b.Write(s.coin_type.AsSlice())
	b.Write(s.tx_hash.AsSlice())
	b.Write(s.receipt_addr.AsSlice())
	return Tick{inner: b.Bytes()}
}

func (s *TickBuilder) TickType(v Byte) *TickBuilder {
	s.tick_type = v
	return s
}

func (s *TickBuilder) TokenId(v Bytes) *TickBuilder {
	s.token_id = v
	return s
}

func (s *TickBuilder) Value(v Uint128) *TickBuilder {
	s.value = v
	return s
}

func (s *TickBuilder) Merchant(v Script) *TickBuilder {
	s.merchant = v
	return s
}

func (s *TickBuilder) CoinType(v Bytes) *TickBuilder {
	s.coin_type = v
	return s
}

func (s *TickBuilder) TxHash(v Bytes) *TickBuilder {
	s.tx_hash = v
	return s
}

func (s *TickBuilder) ReceiptAddr(v Bytes) *TickBuilder {
	s.receipt_addr = v
	return s
}

func NewTickBuilder() *TickBuilder {
	return &TickBuilder{tick_type: ByteDefault(), token_id: BytesDefault(), value: Uint128Default(), merchant: ScriptDefault(), coin_type: BytesDefault(), tx_hash: BytesDefault(), receipt_addr: BytesDefault()}
}

type Tick struct {
	inner []byte
}

func TickFromSliceUnchecked(slice []byte) *Tick {
	return &Tick{inner: slice}
}
func (s *Tick) AsSlice() []byte {
	return s.inner
}

func TickDefault() Tick {
	return *TickFromSliceUnchecked([]byte{118, 0, 0, 0, 32, 0, 0, 0, 33, 0, 0, 0, 37, 0, 0, 0, 53, 0, 0, 0, 106, 0, 0, 0, 110, 0, 0, 0, 114, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 53, 0, 0, 0, 16, 0, 0, 0, 48, 0, 0, 0, 49, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func TickFromSlice(slice []byte, compatible bool) (*Tick, error) {
	sliceLen := len(slice)
	if uint32(sliceLen) < HeaderSizeUint {
		errMsg := strings.Join([]string{"HeaderIsBroken", "Tick", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}

	totalSize := unpackNumber(slice)
	if Number(sliceLen) != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "Tick", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}

	if uint32(sliceLen) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "Tick", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	offsetFirst := unpackNumber(slice[HeaderSizeUint:])
	if uint32(offsetFirst)%HeaderSizeUint != 0 || uint32(offsetFirst) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"OffsetsNotMatch", "Tick", strconv.Itoa(int(offsetFirst % 4)), "!= 0", strconv.Itoa(int(offsetFirst)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	if sliceLen < int(offsetFirst) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "Tick", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(offsetFirst))}, " ")
		return nil, errors.New(errMsg)
	}

	fieldCount := uint32(offsetFirst)/HeaderSizeUint - 1
	if fieldCount < 7 {
		return nil, errors.New("FieldCountNotMatch")
	} else if !compatible && fieldCount > 7 {
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

	_, err = ByteFromSlice(slice[offsets[0]:offsets[1]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[1]:offsets[2]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = Uint128FromSlice(slice[offsets[2]:offsets[3]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = ScriptFromSlice(slice[offsets[3]:offsets[4]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[4]:offsets[5]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[5]:offsets[6]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesFromSlice(slice[offsets[6]:offsets[7]], compatible)
	if err != nil {
		return nil, err
	}

	return &Tick{inner: slice}, nil
}

func (s *Tick) TotalSize() uint {
	return uint(unpackNumber(s.inner))
}
func (s *Tick) FieldCount() uint {
	var number uint = 0
	if uint32(s.TotalSize()) == HeaderSizeUint {
		return number
	}
	number = uint(unpackNumber(s.inner[HeaderSizeUint:]))/4 - 1
	return number
}
func (s *Tick) Len() uint {
	return s.FieldCount()
}
func (s *Tick) IsEmpty() bool {
	return s.Len() == 0
}
func (s *Tick) CountExtraFields() uint {
	return s.FieldCount() - 7
}

func (s *Tick) HasExtraFields() bool {
	return 7 != s.FieldCount()
}

func (s *Tick) TickType() *Byte {
	start := unpackNumber(s.inner[4:])
	end := unpackNumber(s.inner[8:])
	return ByteFromSliceUnchecked(s.inner[start:end])
}

func (s *Tick) TokenId() *Bytes {
	start := unpackNumber(s.inner[8:])
	end := unpackNumber(s.inner[12:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *Tick) Value() *Uint128 {
	start := unpackNumber(s.inner[12:])
	end := unpackNumber(s.inner[16:])
	return Uint128FromSliceUnchecked(s.inner[start:end])
}

func (s *Tick) Merchant() *Script {
	start := unpackNumber(s.inner[16:])
	end := unpackNumber(s.inner[20:])
	return ScriptFromSliceUnchecked(s.inner[start:end])
}

func (s *Tick) CoinType() *Bytes {
	start := unpackNumber(s.inner[20:])
	end := unpackNumber(s.inner[24:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *Tick) TxHash() *Bytes {
	start := unpackNumber(s.inner[24:])
	end := unpackNumber(s.inner[28:])
	return BytesFromSliceUnchecked(s.inner[start:end])
}

func (s *Tick) ReceiptAddr() *Bytes {
	var ret *Bytes
	start := unpackNumber(s.inner[28:])
	if s.HasExtraFields() {
		end := unpackNumber(s.inner[32:])
		ret = BytesFromSliceUnchecked(s.inner[start:end])
	} else {
		ret = BytesFromSliceUnchecked(s.inner[start:])
	}
	return ret
}

func (s *Tick) AsBuilder() TickBuilder {
	ret := NewTickBuilder().TickType(*s.TickType()).TokenId(*s.TokenId()).Value(*s.Value()).Merchant(*s.Merchant()).CoinType(*s.CoinType()).TxHash(*s.TxHash()).ReceiptAddr(*s.ReceiptAddr())
	return *ret
}

type AuthBuilder struct {
	inner [21]Byte
}

func NewAuthBuilder() *AuthBuilder {
	return &AuthBuilder{inner: [21]Byte{ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault(), ByteDefault()}}
}

func (s *AuthBuilder) Build() Auth {
	b := new(bytes.Buffer)
	len := len(s.inner)
	for i := 0; i < len; i++ {
		b.Write(s.inner[i].AsSlice())
	}
	return Auth{inner: b.Bytes()}
}

func (s *AuthBuilder) Set(v [21]Byte) *AuthBuilder {
	s.inner = v
	return s
}

func (s *AuthBuilder) Nth0(v Byte) *AuthBuilder {
	s.inner[0] = v
	return s
}

func (s *AuthBuilder) Nth1(v Byte) *AuthBuilder {
	s.inner[1] = v
	return s
}

func (s *AuthBuilder) Nth2(v Byte) *AuthBuilder {
	s.inner[2] = v
	return s
}

func (s *AuthBuilder) Nth3(v Byte) *AuthBuilder {
	s.inner[3] = v
	return s
}

func (s *AuthBuilder) Nth4(v Byte) *AuthBuilder {
	s.inner[4] = v
	return s
}

func (s *AuthBuilder) Nth5(v Byte) *AuthBuilder {
	s.inner[5] = v
	return s
}

func (s *AuthBuilder) Nth6(v Byte) *AuthBuilder {
	s.inner[6] = v
	return s
}

func (s *AuthBuilder) Nth7(v Byte) *AuthBuilder {
	s.inner[7] = v
	return s
}

func (s *AuthBuilder) Nth8(v Byte) *AuthBuilder {
	s.inner[8] = v
	return s
}

func (s *AuthBuilder) Nth9(v Byte) *AuthBuilder {
	s.inner[9] = v
	return s
}

func (s *AuthBuilder) Nth10(v Byte) *AuthBuilder {
	s.inner[10] = v
	return s
}

func (s *AuthBuilder) Nth11(v Byte) *AuthBuilder {
	s.inner[11] = v
	return s
}

func (s *AuthBuilder) Nth12(v Byte) *AuthBuilder {
	s.inner[12] = v
	return s
}

func (s *AuthBuilder) Nth13(v Byte) *AuthBuilder {
	s.inner[13] = v
	return s
}

func (s *AuthBuilder) Nth14(v Byte) *AuthBuilder {
	s.inner[14] = v
	return s
}

func (s *AuthBuilder) Nth15(v Byte) *AuthBuilder {
	s.inner[15] = v
	return s
}

func (s *AuthBuilder) Nth16(v Byte) *AuthBuilder {
	s.inner[16] = v
	return s
}

func (s *AuthBuilder) Nth17(v Byte) *AuthBuilder {
	s.inner[17] = v
	return s
}

func (s *AuthBuilder) Nth18(v Byte) *AuthBuilder {
	s.inner[18] = v
	return s
}

func (s *AuthBuilder) Nth19(v Byte) *AuthBuilder {
	s.inner[19] = v
	return s
}

func (s *AuthBuilder) Nth20(v Byte) *AuthBuilder {
	s.inner[20] = v
	return s
}

type Auth struct {
	inner []byte
}

func AuthFromSliceUnchecked(slice []byte) *Auth {
	return &Auth{inner: slice}
}
func (s *Auth) AsSlice() []byte {
	return s.inner
}

func AuthDefault() Auth {
	return *AuthFromSliceUnchecked([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}

func AuthFromSlice(slice []byte, _compatible bool) (*Auth, error) {
	sliceLen := len(slice)
	if sliceLen != 21 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "Auth", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(21)}, " ")
		return nil, errors.New(errMsg)
	}
	return &Auth{inner: slice}, nil
}

func (s *Auth) RawData() []byte {
	return s.inner
}

func (s *Auth) Nth0() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[0:1])
	return ret
}

func (s *Auth) Nth1() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[1:2])
	return ret
}

func (s *Auth) Nth2() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[2:3])
	return ret
}

func (s *Auth) Nth3() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[3:4])
	return ret
}

func (s *Auth) Nth4() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[4:5])
	return ret
}

func (s *Auth) Nth5() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[5:6])
	return ret
}

func (s *Auth) Nth6() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[6:7])
	return ret
}

func (s *Auth) Nth7() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[7:8])
	return ret
}

func (s *Auth) Nth8() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[8:9])
	return ret
}

func (s *Auth) Nth9() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[9:10])
	return ret
}

func (s *Auth) Nth10() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[10:11])
	return ret
}

func (s *Auth) Nth11() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[11:12])
	return ret
}

func (s *Auth) Nth12() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[12:13])
	return ret
}

func (s *Auth) Nth13() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[13:14])
	return ret
}

func (s *Auth) Nth14() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[14:15])
	return ret
}

func (s *Auth) Nth15() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[15:16])
	return ret
}

func (s *Auth) Nth16() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[16:17])
	return ret
}

func (s *Auth) Nth17() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[17:18])
	return ret
}

func (s *Auth) Nth18() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[18:19])
	return ret
}

func (s *Auth) Nth19() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[19:20])
	return ret
}

func (s *Auth) Nth20() *Byte {
	ret := ByteFromSliceUnchecked(s.inner[20:21])
	return ret
}

func (s *Auth) AsBuilder() AuthBuilder {
	t := NewAuthBuilder()
	t.Nth0(*s.Nth0())
	t.Nth1(*s.Nth1())
	t.Nth2(*s.Nth2())
	t.Nth3(*s.Nth3())
	t.Nth4(*s.Nth4())
	t.Nth5(*s.Nth5())
	t.Nth6(*s.Nth6())
	t.Nth7(*s.Nth7())
	t.Nth8(*s.Nth8())
	t.Nth9(*s.Nth9())
	t.Nth10(*s.Nth10())
	t.Nth11(*s.Nth11())
	t.Nth12(*s.Nth12())
	t.Nth13(*s.Nth13())
	t.Nth14(*s.Nth14())
	t.Nth15(*s.Nth15())
	t.Nth16(*s.Nth16())
	t.Nth17(*s.Nth17())
	t.Nth18(*s.Nth18())
	t.Nth19(*s.Nth19())
	t.Nth20(*s.Nth20())
	return *t
}

type IdentityOptBuilder struct {
	isNone bool
	inner  Bytes
}

func NewIdentityOptBuilder() *IdentityOptBuilder {
	return &IdentityOptBuilder{isNone: true, inner: BytesDefault()}
}
func (s *IdentityOptBuilder) Set(v Bytes) *IdentityOptBuilder {
	s.isNone = false
	s.inner = v
	return s
}
func (s *IdentityOptBuilder) Build() IdentityOpt {
	var ret IdentityOpt
	if s.isNone {
		ret = IdentityOpt{inner: []byte{}}
	} else {
		ret = IdentityOpt{inner: s.inner.AsSlice()}
	}
	return ret
}

type IdentityOpt struct {
	inner []byte
}

func IdentityOptFromSliceUnchecked(slice []byte) *IdentityOpt {
	return &IdentityOpt{inner: slice}
}
func (s *IdentityOpt) AsSlice() []byte {
	return s.inner
}

func IdentityOptDefault() IdentityOpt {
	return *IdentityOptFromSliceUnchecked([]byte{})
}

func IdentityOptFromSlice(slice []byte, compatible bool) (*IdentityOpt, error) {
	if len(slice) == 0 {
		return &IdentityOpt{inner: slice}, nil
	}

	_, err := BytesFromSlice(slice, compatible)
	if err != nil {
		return nil, err
	}
	return &IdentityOpt{inner: slice}, nil
}

func (s *IdentityOpt) IntoBytes() (*Bytes, error) {
	if s.IsNone() {
		return nil, errors.New("No data")
	}
	return BytesFromSliceUnchecked(s.AsSlice()), nil
}
func (s *IdentityOpt) IsSome() bool {
	return len(s.inner) != 0
}
func (s *IdentityOpt) IsNone() bool {
	return len(s.inner) == 0
}
func (s *IdentityOpt) AsBuilder() IdentityOptBuilder {
	var ret = NewIdentityOptBuilder()
	if s.IsSome() {
		ret.Set(*BytesFromSliceUnchecked(s.AsSlice()))
	}
	return *ret
}

type OmniLockWitnessLockBuilder struct {
	signature     BytesOpt
	omni_identity IdentityOpt
	preimage      BytesOpt
}

func (s *OmniLockWitnessLockBuilder) Build() OmniLockWitnessLock {
	b := new(bytes.Buffer)

	totalSize := HeaderSizeUint * (3 + 1)
	offsets := make([]uint32, 0, 3)

	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.signature.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.omni_identity.AsSlice()))
	offsets = append(offsets, totalSize)
	totalSize += uint32(len(s.preimage.AsSlice()))

	b.Write(packNumber(Number(totalSize)))

	for i := 0; i < len(offsets); i++ {
		b.Write(packNumber(Number(offsets[i])))
	}

	b.Write(s.signature.AsSlice())
	b.Write(s.omni_identity.AsSlice())
	b.Write(s.preimage.AsSlice())
	return OmniLockWitnessLock{inner: b.Bytes()}
}

func (s *OmniLockWitnessLockBuilder) Signature(v BytesOpt) *OmniLockWitnessLockBuilder {
	s.signature = v
	return s
}

func (s *OmniLockWitnessLockBuilder) OmniIdentity(v IdentityOpt) *OmniLockWitnessLockBuilder {
	s.omni_identity = v
	return s
}

func (s *OmniLockWitnessLockBuilder) Preimage(v BytesOpt) *OmniLockWitnessLockBuilder {
	s.preimage = v
	return s
}

func NewOmniLockWitnessLockBuilder() *OmniLockWitnessLockBuilder {
	return &OmniLockWitnessLockBuilder{signature: BytesOptDefault(), omni_identity: IdentityOptDefault(), preimage: BytesOptDefault()}
}

type OmniLockWitnessLock struct {
	inner []byte
}

func OmniLockWitnessLockFromSliceUnchecked(slice []byte) *OmniLockWitnessLock {
	return &OmniLockWitnessLock{inner: slice}
}
func (s *OmniLockWitnessLock) AsSlice() []byte {
	return s.inner
}

func OmniLockWitnessLockDefault() OmniLockWitnessLock {
	return *OmniLockWitnessLockFromSliceUnchecked([]byte{16, 0, 0, 0, 16, 0, 0, 0, 16, 0, 0, 0, 16, 0, 0, 0})
}

func OmniLockWitnessLockFromSlice(slice []byte, compatible bool) (*OmniLockWitnessLock, error) {
	sliceLen := len(slice)
	if uint32(sliceLen) < HeaderSizeUint {
		errMsg := strings.Join([]string{"HeaderIsBroken", "OmniLockWitnessLock", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint))}, " ")
		return nil, errors.New(errMsg)
	}

	totalSize := unpackNumber(slice)
	if Number(sliceLen) != totalSize {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "OmniLockWitnessLock", strconv.Itoa(int(sliceLen)), "!=", strconv.Itoa(int(totalSize))}, " ")
		return nil, errors.New(errMsg)
	}

	if uint32(sliceLen) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"TotalSizeNotMatch", "OmniLockWitnessLock", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	offsetFirst := unpackNumber(slice[HeaderSizeUint:])
	if uint32(offsetFirst)%HeaderSizeUint != 0 || uint32(offsetFirst) < HeaderSizeUint*2 {
		errMsg := strings.Join([]string{"OffsetsNotMatch", "OmniLockWitnessLock", strconv.Itoa(int(offsetFirst % 4)), "!= 0", strconv.Itoa(int(offsetFirst)), "<", strconv.Itoa(int(HeaderSizeUint * 2))}, " ")
		return nil, errors.New(errMsg)
	}

	if sliceLen < int(offsetFirst) {
		errMsg := strings.Join([]string{"HeaderIsBroken", "OmniLockWitnessLock", strconv.Itoa(int(sliceLen)), "<", strconv.Itoa(int(offsetFirst))}, " ")
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

	_, err = BytesOptFromSlice(slice[offsets[0]:offsets[1]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = IdentityOptFromSlice(slice[offsets[1]:offsets[2]], compatible)
	if err != nil {
		return nil, err
	}

	_, err = BytesOptFromSlice(slice[offsets[2]:offsets[3]], compatible)
	if err != nil {
		return nil, err
	}

	return &OmniLockWitnessLock{inner: slice}, nil
}

func (s *OmniLockWitnessLock) TotalSize() uint {
	return uint(unpackNumber(s.inner))
}
func (s *OmniLockWitnessLock) FieldCount() uint {
	var number uint = 0
	if uint32(s.TotalSize()) == HeaderSizeUint {
		return number
	}
	number = uint(unpackNumber(s.inner[HeaderSizeUint:]))/4 - 1
	return number
}
func (s *OmniLockWitnessLock) Len() uint {
	return s.FieldCount()
}
func (s *OmniLockWitnessLock) IsEmpty() bool {
	return s.Len() == 0
}
func (s *OmniLockWitnessLock) CountExtraFields() uint {
	return s.FieldCount() - 3
}

func (s *OmniLockWitnessLock) HasExtraFields() bool {
	return 3 != s.FieldCount()
}

func (s *OmniLockWitnessLock) Signature() *BytesOpt {
	start := unpackNumber(s.inner[4:])
	end := unpackNumber(s.inner[8:])
	return BytesOptFromSliceUnchecked(s.inner[start:end])
}

func (s *OmniLockWitnessLock) OmniIdentity() *IdentityOpt {
	start := unpackNumber(s.inner[8:])
	end := unpackNumber(s.inner[12:])
	return IdentityOptFromSliceUnchecked(s.inner[start:end])
}

func (s *OmniLockWitnessLock) Preimage() *BytesOpt {
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

func (s *OmniLockWitnessLock) AsBuilder() OmniLockWitnessLockBuilder {
	ret := NewOmniLockWitnessLockBuilder().Signature(*s.Signature()).OmniIdentity(*s.OmniIdentity()).Preimage(*s.Preimage())
	return *ret
}
