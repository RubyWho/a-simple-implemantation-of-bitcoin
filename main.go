package main

import (
	"fmt"
	"time"
	"crypto/sha256"
	"strconv"
)

//将block转换为[]byte HashMag
func (block *Block_t) convToByted() (hashMsg []byte) {
	//hash的消息有：id、blocksize块太小、Timestamp时间戳、PreHash上一个hash值、NumOfTrans所有的交易数量
	//注：hash的消息内容可根据要求随时调整
//	msgStr := fmt.Sprint("%d%d%s%s%d",block.id, block.BLockSize, block.Timestamp, block.PreHash, block.body.NumOfTrans)
	msgStr := strconv.Itoa(block.id) + strconv.Itoa(block.BLockSize) + block.Timestamp + block.PreHash + strconv.Itoa(block.body.NumOfTrans)
	fmt.Printf("%s\n",msgStr)
	//转换为切片
	hashMsg = []byte(msgStr)
	return
}

//计算hash值
func calcHash(block *Block_t) (hash string){
	//将要参与计算hash的消息转换为[]byte格式
	hashMsg := block.convToByted()
	//func Sum256(data []byte) [Size]byte
	//返回数据的SHA256校验和。
	hashByte := sha256.Sum256(hashMsg)
	//拷贝hash成切片到block.header.hash
	//string(hash[:])转换为string类型
	hash = string(hashByte[:])
	fmt.Printf("hash:%x\n",hash)
	return hash
}

//获取链上的最后一区块
func getLatestBlock(blockChain BlockChain)(latestBlock Block_t, isNull bool){
	isNull = true
	if len(blockChain) != 0 {
		//这里使用的是拷贝值的方式进行返回，不使用地址拷贝（指针）的原因是防止该区块在返回之后被人串改
		latestBlock = *blockChain[len(blockChain)-1]
		isNull = false
	}
	return 
}

//创建新块
func generateBlock(blockBody BlockBody_t, preHash string)(newBlock *Block_t){
	
	//创建新区块，并初始化：块大小、版本号、上一块的hash值、时间戳、blockBody
	//这些值有些在计算hash值时会用到
	newBlock = new(Block_t)
	//获取最后一块区块
	latestBlock,isNull := getLatestBlock(blockChain)
	if isNull{
		newBlock.id = 0
	}else{
		newBlock.id = latestBlock.id + 1
	}
	newBlock.BLockSize = 4
	newBlock.Version = "V1.0"
	newBlock.PreHash = preHash
	newBlock.Timestamp = time.Now().String()
	newBlock.body = blockBody
	//计算本区块的hash值
	newBlock.hash = calcHash(newBlock)
	return 
}

//创建创世块
func CreateGenesisBlock() *Block_t {
	//创建新块作为创世块
	genesisBlock := generateBlock(BlockBody_t{0,make([]*Transaction_t,0,0)},"")
	//返回
	return genesisBlock
}

//创建链
func CreateBlockChain() (blockChain BlockChain){
	//创建链
	blockChain = make(BlockChain,0)
	//返回链
	return
}


//验证区块的合法性
//此步骤在上链前执行，不在产生区块时执行
func (blockChain BlockChain) isValidBlock(block *Block_t) bool{
	//此处简单判断其合法性
	//判断新区块的hash值是否正确 
	if block.hash != calcHash(block) {
		fmt.Println("hash error")
		return false
	}
	
	//如果是创世块，取消hash值验证
	if len(blockChain) == 0 {
		fmt.Println("创世块")
		return true
	}
	if block.id == 0 {
		return true
	}
	//非创世块：判断此链的最后一块的hash是否为新区块的上一块hash值
	if  blockChain[len(blockChain)-1].hash != block.PreHash {
		fmt.Println("PreHash error")
		return false
	}
	return true
}

//将新的区块上链
func (blockChain *BlockChain)appendBlock(block *Block_t) bool {
	//如果该区块合法
	if blockChain.isValidBlock(block){
		//上链
		*blockChain = append(*blockChain,block)
		return true
	}
	return false
}

//菜单
func menu() {
	fmt.Println("***********BlockChain Menu**************")
	fmt.Print("  ******** 1:Show the BlockChain ********\n")
	fmt.Print("  ******** 2:Query the Database *********\n")
	fmt.Print("  ******** 3:Append Block Chain *********\n")
	fmt.Print("  ******** 0:Exit BlockChain Sys *********\n")
	fmt.Println("****************************************")
	fmt.Print("***choose:")
}


//测试专用打印机
func (blockChain BlockChain) printBlockChain(){
	for index:= 0; index < len(blockChain); index++{
		fmt.Printf("id:%d\n",blockChain[index].id)
		fmt.Printf("BLockSize:%d\n",blockChain[index].BLockSize)
		fmt.Printf("Version:%s\n",blockChain[index].Version)
		fmt.Printf("PreHash:%x\n",blockChain[index].PreHash)
		fmt.Printf("Timestamp:%s\n",blockChain[index].Timestamp)
		fmt.Printf("hash:%x\n",blockChain[index].hash)
		fmt.Printf("NumOfTrans:%d\n",blockChain[index].body.NumOfTrans)
		fmt.Println()
	}
}

//Block测试打印机
func printBlock(block *Block_t){
	fmt.Printf("id:%d\n",block.id)
	fmt.Printf("BLockSize:%d\n",block.BLockSize)
	fmt.Printf("Version:%s\n",block.Version)
	fmt.Printf("PreHash:%x\n",block.PreHash)
	fmt.Printf("Timestamp:%s\n",block.Timestamp)
	fmt.Printf("hash:%x\n",block.hash)
	fmt.Printf("NumOfTrans:%d\n",block.body.NumOfTrans)
}

func main() {
	//创建链
	blockChain = CreateBlockChain()
	//打开数据库
	openBlockChainDB()
	//程序最后关闭数据库
	defer db.Close()
	go listenPort()
Stop_Lable:
	for{
		menu()
		var choose int
		fmt.Scanf("%d",&choose);
		switch choose{
			case 0:
				break Stop_Lable
			case 1:
				blockChain.printBlockChain()
			case 2:
				queryData()
			case 3:
				blockBody := BlockBody_t{1,make([]*Transaction_t,0,0)}
				block := generateBlock(blockBody,blockChain[len(blockChain)-1].hash)
				//appendBlock(&blockChain,block) //配合将新的区块上链2使用
//				(&blockChain).appendBlock(block)
				data := blockToByte(block)
				broadcast(data)
			default:
		}
	}
}
