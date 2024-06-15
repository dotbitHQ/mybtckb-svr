package contract

import "github.com/nervosnetwork/ckb-sdk-go/types"

type TypeContractInfo struct {
	TxHash   string
	Index    uint
	CodeHash string //hash(type script) or hash(data)
	HashType types.ScriptHashType
}
type Env struct {
	Xudt, Spore, SporeCluster, RGBPP, UniqueCell *TypeContractInfo
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
		TxHash:   "0x5e8d2a517d50fd4bb4d01737a7952a1f1d35c8afc77240695bb569cd7d9d5a1f",
		Index:    0,
		CodeHash: "0x685a60219309029d01310311dba953d67029170ca4848a4ff638e57002130a0d",
		HashType: types.HashTypeData1,
	},
	SporeCluster: &TypeContractInfo{
		TxHash:   "0xcebb174d6e300e26074aea2f5dbd7f694bb4fe3de52b6dfe205e54f90164510a",
		Index:    0,
		CodeHash: "0x0bbe768b519d8ea7b96d58f1182eb7e6ef96c541fbd9526975077ee09f049058",
		HashType: types.HashTypeData1,
	},
	UniqueCell: &TypeContractInfo{
		TxHash:   "0xff91b063c78ed06f10a1ed436122bd7d671f9a72ef5f5fa28d05252c17cf4cef",
		Index:    0,
		CodeHash: "0x8e341bcfec6393dcd41e635733ff2dca00a6af546949f70c57a706c0f344df8b",
		HashType: types.HashTypeType,
	},
	RGBPP: &TypeContractInfo{
		TxHash:   "",
		Index:    0,
		CodeHash: "0x61ca7a4796a4eb19ca4f0d065cb9b10ddcf002f10f7cbb810c706cb6bb5c3248",
		HashType: types.HashTypeType,
	},
}
