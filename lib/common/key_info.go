package common

type CoinType string // EIP-155

// https://github.com/satoshilabs/slips/blob/master/slip-0044.md
const (
	CoinTypeEth      CoinType = "60"
	CoinTypeTrx      CoinType = "195"
	CoinTypeBNB      CoinType = "714"
	CoinTypeBSC      CoinType = "9006"
	CoinTypeMatic    CoinType = "966"
	CoinTypeDogeCoin CoinType = "3"
	CoinTypeCKB      CoinType = "309"
)

type ChainId string //BIP-44

const (
	ChainIdEthMainNet     ChainId = "1"
	ChainIdBscMainNet     ChainId = "56"
	ChainIdPolygonMainNet ChainId = "137"

	ChainIdEthTestNet     ChainId = "5" // Goerli
	ChainIdBscTestNet     ChainId = "97"
	ChainIdPolygonTestNet ChainId = "80001"
)
