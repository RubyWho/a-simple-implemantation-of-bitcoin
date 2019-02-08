package main

import(
	"net"
	"fmt"
	"strconv"
	"strings"
)

//处理接收到的数据
func handleData(senderAddr string, data []byte,size int){
	//转换为Block
	block, success := stringToBlock(string(data[:size]))
	if success{
		if (&blockChain).appendBlock(block){
			//插入到数据库中
			insertData(block)
		}
	}
}

//监听端口函数
func listenPort(){
	
	//创建监听
	listener, err := net.ListenUDP("udp",&net.UDPAddr{IP: net.IPv4zero, Port: 9999})
	//创建监听失败
	if err != nil {
		fmt.Println(err)
		return
	}
	//
	for{
		data := make([]byte, 256)
		//读取数据:size=数据大小，senderAddr=发送方的地址
		size, senderAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
		}
		//处理数据
		go handleData(senderAddr.String(),data,size)
	} 
}

//广播
func broadcast(data []byte) {
	//设置广播地址：根据自己的网段设置广播地址
	ip := net.ParseIP("192.168.1.255")
	//源地址设为0,端口号也设为0
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	//目标地址设为广播地址，广播端口为9999
	dstAddr := &net.UDPAddr{IP: ip, Port: 9999}
	//通过ListenUDP创建一个unconnected的 *UDPConn
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		fmt.Println(err)
	}
	//发送数据
	n, err := conn.WriteToUDP(data, dstAddr)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("send size:",n)
}

//将Block_t转换为Byte，字段之间用';'隔开
func blockToByte(block *Block_t) (BlockBytes []byte){
	//拼接字符串
	blockStr := "*;" + strconv.Itoa(block.id) + ";" + strconv.Itoa(block.BLockSize) + ";" + block.Version + ";" + block.Timestamp + ";" + block.PreHash + ";" + block.hash + ";" + strconv.Itoa(block.body.NumOfTrans) + ";#"
	//转换为[]byte
	BlockBytes = []byte(blockStr)
	return
}

//将string转换为Block_t
func stringToBlock(blockStr string) (block *Block_t, isSuccess bool){
	isSuccess = true
	//分割字符串
	strs := strings.Split(blockStr,";")
	//判断数据是否已经接受完毕
	if strs[0] != "*" || strs[len(strs)-1] != "#" {
		fmt.Println("data is failed receive")
		isSuccess = false
		return
	}
	//转换为Block
	block = new(Block_t)
	block.id, _ = strconv.Atoi(strs[1])
	block.BLockSize, _ = strconv.Atoi(strs[2])
	block.Version = strs[3]
	block.Timestamp = strs[4]
	block.PreHash = strs[5]
	block.hash = strs[6]
	block.body.NumOfTrans, _= strconv.Atoi(strs[7])
	return
}