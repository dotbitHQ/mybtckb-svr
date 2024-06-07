package contract

import "github.com/nervosnetwork/ckb-sdk-go/types"

type TypeContractInfo struct {
	TxHash   string
	Index    uint
	CodeHash string //hash(type script) or hash(data)
	HashType types.ScriptHashType
}
type Env struct {
	Xudt, Spore, SporeCluster, UniqueCell *TypeContractInfo
}

var EnvMain = Env{
	Xudt: &TypeContractInfo{
		TxHash:   "0xc07844ce21b38e4b071dd0e1ee3b0e27afd8d7532491327f39b786343f558ab7",
		Index:    0,
		CodeHash: "0x50bd8d6680b8b9cf98b73f3c08faf8b2a21914311954118ad6609be6e78a1b95",
		HashType: types.HashTypeData1,
	},
	Spore: &TypeContractInfo{
		TxHash:   "0x96b198fb5ddbd1eed57ed667068f1f1e55d07907b4c0dbd38675a69ea1b69824",
		Index:    0,
		CodeHash: "0x4a4dce1df3dffff7f8b2cd7dff7303df3b6150c9788cb75dcf6747247132b9f5",
		HashType: types.HashTypeData1,
	},
	SporeCluster: &TypeContractInfo{
		TxHash:   "0xe464b7fb9311c5e2820e61c99afc615d6b98bdefbe318c34868c010cbd0dc938",
		Index:    0,
		CodeHash: "0x7366a61534fa7c7e6225ecc0d828ea3b5366adec2b58206f2ee84995fe030075",
		HashType: types.HashTypeData1,
	},
	UniqueCell: &TypeContractInfo{
		TxHash:   "0x67524c01c0cb5492e499c7c7e406f2f9d823e162d6b0cf432eacde0c9808c2ad",
		Index:    0,
		CodeHash: "0x2c8c11c985da60b0a330c61a85507416d6382c130ba67f0c47ab071e00aec628",
		HashType: types.HashTypeData1,
	},
}

var EnvTest = Env{
	Xudt: &TypeContractInfo{
		TxHash:   "0xbf6fb538763efec2a70a6a3dcb7242787087e1030c4e7d86585bc63a9d337f5f",
		Index:    0,
		CodeHash: "0x25c29dc317811a6f6f3985a7a9ebc4838bd388d19d0feeecf0bcd60f6c0975bb",
		HashType: types.HashTypeType,
	},
	Spore: &TypeContractInfo{
		TxHash:   "0x06995b9fc19461a2bf9933e57b69af47a20bf0a5bc6c0ffcb85567a2c733f0a1",
		Index:    0,
		CodeHash: "0x5e063b4c0e7abeaa6a428df3b693521a3050934cf3b0ae97a800d1bc31449398",
		HashType: types.HashTypeData1,
	},
	SporeCluster: &TypeContractInfo{
		TxHash:   "0xfbceb70b2e683ef3a97865bb88e082e3e5366ee195a9c826e3c07d1026792fcd",
		Index:    0,
		CodeHash: "0x7366a61534fa7c7e6225ecc0d828ea3b5366adec2b58206f2ee84995fe030075",
		HashType: types.HashTypeData1,
	},
	UniqueCell: &TypeContractInfo{
		TxHash:   "0xff91b063c78ed06f10a1ed436122bd7d671f9a72ef5f5fa28d05252c17cf4cef",
		Index:    0,
		CodeHash: "0x8e341bcfec6393dcd41e635733ff2dca00a6af546949f70c57a706c0f344df8b",
		HashType: types.HashTypeType,
	},
}
