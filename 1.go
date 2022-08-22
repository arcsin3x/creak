package main
 
import (
	"fmt"
	"os"
	"math/big"
	"log"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"crypto/ecdsa"
	"net/smtp"
	"github.com/scorredoira/email"
	"net/mail"
	"github.com/ethereum/go-ethereum/rpc"
	"strconv"
)
 
//这里我们需要定义两个回应的格式
type Info struct {
	From string	//info.From:发件地址或账号
	To string	//to:  收件地址
	Title string	//标题
	Body string	//body:邮件内容
	Host string	//info.Host:邮件服务器地址
	Password string	//info.Password:密码
}
 
 
var reply interface{}
 
func main() {
	//GoMail("测试测试")
	f, err := os.OpenFile("./key.txt", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		// 打开文件失败处理
		fmt.Println(err)
		return
	}
	for i:=0;i<10;i++  {
		go GetBlance(f)
	}
	select {
 
	}
 
	defer f.Close()
}
 
func GetBlance(_file *os.File) {
	client, err := rpc.Dial("https://mainnet.infura.io/v3/a3c59244a3d04a77bd0509816d8954c7")
	if err != nil {
		fmt.Println("错误:", err)
	}
	for ; ; {
		priv, addr := CreateKey()
		err2 := client.Call(&reply, "eth_getBalance", addr, "latest") //第一个是用来存放回复数据的格式，第二个是请求方法
		if err2 != nil {
			fmt.Println("错误:", err2)
		}
		fmt.Println(priv, addr)
		//这里得到的还是16进制的需要做个进制转换成10进制
		if reply != "0x0" {
			n := new(big.Int)
			n, _ = n.SetString(reply.(string)[2:], 16)
			//fmt.Println(n)
			// 查找文件末尾的偏移量
			n1, _ := _file.Seek(0, 2)
			value,_:=strconv.ParseFloat(n.String(),64)
			content := "私钥为:"+priv + "\t地址为:" + addr + "\t余额为:" + strconv.FormatFloat(value/1000000000000000000,'f',5,64) + "ETH \n"
			GoMail(content)
			Infos.Body=content
			// 从末尾的偏移量开始写入内容
			_, err = _file.WriteAt([]byte(content), n1)
		}
	}
	defer client.Close()
}
 
func CreateKey() (privs, addrs string) {
	//创建私钥
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	/*	//可通过此代码导入私钥
	privateKey,err=crypto.HexToECDSA("93d5d04256882aaad507ff09f510969f347758109793448aa79e1b4dbe5f6efa")
	if err != nil {
		log.Fatal(err)
	}
	*/
	privateKeyBytes := crypto.FromECDSA(privateKey)
	priv := hexutil.Encode(privateKeyBytes)[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	//fmt.Println(address)
	return priv, address
}
 
 
var Infos=Info{
	"arcsin3x@126.com",//发件邮箱
	"631461051@qq.com",//收件邮箱
	"碰撞成功",
	"",
	"smtp.126.com:25",
	"MWBEBKSULMMOFHML",//发件邮箱的授权码可以再邮箱后台获取
}
 
func GoMail(body string)  {
	m := email.NewMessage(Infos.Title, body)
	m.From = mail.Address{"来自我的电脑",Infos.From}
	m.To = []string{Infos.To}
	email.Send(Infos.Host, smtp.PlainAuth("", Infos.From, Infos.Password, "smtp.163.com"), m)
}