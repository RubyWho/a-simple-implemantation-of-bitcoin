package main

import(
	"fmt"
	"database/sql"
	"strconv"
	_"github.com/mattn/go-sqlite3"//使用“_”是只是用其init函数而不使用其内部的变量和函数 
)
var db *sql.DB

//创建Block表：此版本的表只用于测试，若真是发应用则需重新设计该数据库
const sql_CreateBlock_t = `CREATE TABLE IF NOT EXISTS Block_t(
		id INTEGER PRIMARY KEY NOT NULL,
		BLockSize INTEGER NOT NULL,
		Version VARCHAR(64) NOT NULL,
		PreHash VARCHAR(64) NOT NULL,
		Timestamp VARCHAR(64) NOT NULL,
		hash VARCHAR(64) NOT NULL,
		NumOfTrans INTEGER NOT NULL);`


//插入数据
func insertData(block *Block_t)bool {
	//values(?)表示占位符
	sql := `INSERT INTO Block_t(id, BLockSize, Version, PreHash, Timestamp, hash, NumOfTrans) values(?,?,?,?,?,?,?)`
	//预处理sql语句
	//Prepare创建一个准备好的状态用于之后的查询和命令，返回值可以同时执行多个查询和命令。返回值：Stmt，error
	//stmt是准备好的状态，stmt可以安全的被多个go程同时使用
	stmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer stmt.Close()
//	preHash, _:= fmt.Printf("%x",block.PreHash)
//	hash, _:= fmt.Printf("%x",block.hash)
	//正式插入数据
	_, err = stmt.Exec(block.id, block.BLockSize, block.Version, block.PreHash, block.Timestamp, block.hash, block.body.NumOfTrans)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

//查询数据->结果直接放到链上：blockChain
func queryData() (err error) {
	//查询语句
	sql := `SELECT * FROM Block_t`
	//开始查询，查询结果放在rows中
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	//最后要关闭rowsClose
	defer rows.Close()
	//循环遍历rows
	for rows.Next(){
		var newBlock *Block_t
		newBlock = new(Block_t)
		err = rows.Scan(&newBlock.id, &newBlock.BLockSize, &newBlock.Version,
						&newBlock.PreHash, &newBlock.Timestamp,
						&newBlock.hash, &newBlock.body.NumOfTrans)
		if err != nil {
			fmt.Println(err)
			return
		}
		(&blockChain).appendBlock(newBlock)
	}
	return
}
//根据id查询数据
func queryDataByID(id int) (block *Block_t, err error) {
	//查询语句
	sql := `SELECT * FROM Block_t WHERE id = ` + strconv.Itoa(id)
	//开始查询，查询结果放在rows中
	rows, err := db.Query(sql)
	if err != nil {
		fmt.Println(err)
		return
	}
	//最后要关闭rows
	defer rows.Close()
	//由于使用id查询，结果只有1个或者没有，所以下面使用if进行判断rows
	if rows.Next(){
		block = new(Block_t)
		err = rows.Scan(&block.id, &block.BLockSize, &block.Version,
						&block.PreHash, &block.Timestamp,
						&block.hash, &block.body.NumOfTrans)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	return
}

//打开数据库，如果不存在则创建
//注意：在这里没有关闭数据库，所以在调用这函数后记得调用一下：db.Close()
func openBlockChainDB() (err error){
	//打开数据库
	db, err = sql.Open("sqlite3","blockChain.db")
	if err != nil {
		fmt.Println(err)
		return
	}

//	//创建表，如果不存在的话创建，存在则跳过
	db.Exec(sql_CreateBlock_t)
	
	//创世块的编号为0
	block, err := queryDataByID(0)
//	//如果表不存在，则此表为第一次创建的新表，需要把创世块添加到书库链中
	if err == nil && block == nil{
		//创建创世块
		blcok := CreateGenesisBlock()
		//插入到数据库
		insertData(blcok)
	}
	//第一次更新数据库
	queryData()
	return
}