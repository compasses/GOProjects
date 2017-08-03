package common

//the first 2 bytes are ASCII Fa
//the next 6 bytes are the directory block height.
//the next 32 bytes are the KeyMR  of the directory block at that height

type DirectoryBlockAnchorInfo struct {
	KeyMR    *Hash
	DBHeight uint32
}

type DBlockHeaderForAnchor struct {
	DBHeight uint32 `json:"dbheight"`
}

type DBlockForAnchor struct {
	KeyMR  string                `json:"keymr"`
	Header DBlockHeaderForAnchor `json:"header"`
}

type EthTxReceipt struct {
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	From             string `json:"from"`
	Root             string `json:"root"`
	TransactionHash  string `json:"transactionHash"`
	TransactionIndex string `json:"transactionIndex"`
}
