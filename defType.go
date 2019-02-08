package main

var blockChain BlockChain

//区块链信息结构体
type BlockChain []*Block_t

//区块类型信息结构体
type Block_t struct {
	/*header */
	BlockHeader_t //区块头
	body BlockBody_t   //区块体
}

//区块头
type BlockHeader_t struct {
	id int				//id:可作为每一个区块的唯一标识
	BLockSize int  //区块字节大小
	Version   string //版本号
	PreHash   string //上一块的哈希值
	Timestamp string  //时间戳
	hash      string //本区块的hash
	// MerkleHash string //本区块记录交易的默克尔树根哈希值
	// DegreeDif  float32 //难度系数
	// RandomNum  int   //随机数
}

//区块内容
type BlockBody_t struct {
	NumOfTrans int            //交易数量
	AllTrans   []*Transaction_t //所有交易
}

//交易Transaction
type Transaction_t struct {
	NumberOfTrans  int       //交易数量
	NumOfSrcFunds  int       //资金来源数量
	NumOfAimsFunds int       //资金去向数量
	TransTime      string      //交易时间
	SrcFunds       SrcFunds_t  //资金来源
	AimsFunds      AimsFunds_t //资金去向
}

//资金来源
type SrcFunds_t struct {
	HashOfTrsans string //资金来源的交易哈希值
}

//资金去向
type AimsFunds_t struct {
	NumOfFunds string //资金的数量（资金转出的个数）
}
