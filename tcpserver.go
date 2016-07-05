package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	log "github.com/cihub/seelog"
	"github.com/kaepa3/btext"
	"github.com/kaepa3/cmdbk"
	"github.com/kaepa3/cycle"
)

var sendQue = list.New()
var config ConfigData

type HealthFile struct {
	Time int    `json:"time"`
	File string `json:"file"`
}

// ConfigData is json format
type ConfigData struct {
	Port   string     `json:"port"`
	Health HealthFile `json:"health"`
}

func initLogger() {
	logConfig := `
	<seelog type="adaptive" mininterval="200000000" maxinterval="1000000000" critmsgcount="5">
		<formats>
		    <format id="main" format="Time:%Date(2006/01/02) %Time	file:%File	func:%FuncShort	line:%Line	level:%LEV	msg:%Msg%n" />
		    <format id="con" format="%Msg%n" />
		</formats>
		<outputs formatid="main">
			<rollingfile filename="log.log" type="size" maxsize="102400" maxrolls="1" formatid = "main"/>
			<console formatid = "con"/>
		</outputs>
	</seelog>`
	logger, err := log.LoggerFromConfigAsBytes([]byte(logConfig))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	log.ReplaceLogger(logger)
}
func logging(text string) {
	log.Info(text)
}

func main() {
	initLogger()
	config = settingGet("config.json")
	tcpAddr, err := net.ResolveTCPAddr("tcp", config.Port)
	checkError(err)
	listner, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}
}

func settingGet(configPath string) ConfigData {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	var datasets ConfigData
	jsonErr := json.Unmarshal(file, &datasets)
	if jsonErr != nil {
		fmt.Println("Format Error: ", jsonErr)
		panic(err)
	}
	return datasets
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	fmt.Println("client accept!")
	obj := cycle.CycleProc{Time: config.Health.Time, Flg: true, Action: addFileWrapper}
	obj.Action()
	cycle.DoProcess(obj)
	cmdbk.Start(callBack)
	for {
		revcivePacket(conn)
		sendPacket(conn)
		time.Sleep(10 * time.Millisecond)
	}
}
func callBack(text string) {
	addFile(text)
}

func addFileWrapper() {
	addFile(config.Health.File)
}

func addFile(text string) {
	_, err := os.Stat(text)
	if err != nil {
		fmt.Println(err, text)
		return
	}
	contents := btext.BParseFile(text)
	if len(contents) != 0 {
		sendQue.PushBack(contents)
	}
}

func revcivePacket(conn net.Conn) {
	messageBuf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(10 * time.Millisecond))
	messageLen, err := conn.Read(messageBuf)
	if 0 == revcheckErr(err) {
		message := string(btext.TParseAry(messageBuf[:messageLen]))
		logging("[rev]->\n" + message)
	}
}

func sendPacket(conn net.Conn) {
	if sendQue.Len() != 0 {
		message := sendQue.Remove(sendQue.Front())
		conn.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
		switch buff := message.(type) {
		case []byte:
			conn.Write(buff)
			logging("[Send]->\n" + btext.TParseAry(buff))

		}
	}
}

func revcheckErr(err error) (retVal int) {
	retVal = 0
	if err != nil {
		if strings.Index(err.Error(), "timeout") == -1 {
			checkError(err)
		}
		retVal = -1
	}
	return
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: error: %s", err.Error())
		os.Exit(1)
	}
}
