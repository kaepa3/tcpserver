package main

import (
	"container/list"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/cihub/seelog"
	"github.com/kaepa3/btext"
	"github.com/kaepa3/cmdbk"
	"github.com/kaepa3/cycle"
	"github.com/kaepa3/tcpserver/localpackage"
)

var sendQue = list.New()
var config ConfigData
var m = new(sync.Mutex)

type HealthFile struct {
	Time int    `json:"time"`
	File string `json:"file"`
}

type ResponceConfig struct {
	Code string `json:"code"`
	File string `json:"file"`
}

// ConfigData is json format
type ConfigData struct {
	Port      string           `json:"port"`
	TimeGrant bool             `json:"time_grant"`
	Health    HealthFile       `json:"health"`
	Responce  []ResponceConfig `json:"responce"`
}

func initLogger() {
	logConfig := `
	<seelog type="adaptive" mininterval="200000000" maxinterval="1000000000" critmsgcount="5">
		<formats>
		    <format id="main" format="Time:%Date(2006/01/02) %Time	file:%File	func:%FuncShort	line:%Line	level:%LEV	msg:%Msg%n" />
		    <format id="con" format="%Msg%n" />
		</formats>
		<outputs formatid="main">
			<rollingfile filename="rev.log" type="size" maxsize="102400" maxrolls="1" formatid = "main"/>
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

func main() {
	initLogger()
	config = settingGet("config.json")
	fmt.Println("setting:", config)

	tcpAddr, err := net.ResolveTCPAddr("tcp", config.Port)
	if err != nil {
		log.Error(err.Error())
	}
	listner, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Error(err.Error())
	}

	log.Info("start" + strings.Repeat("*", 20))
	for {
		conn, err := listner.AcceptTCP()
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

func handleClient(conn *net.TCPConn) {
	fmt.Println("client accept!")

	obj := cycle.CycleProc{Time: config.Health.Time, Flg: true, Action: addFileWrapper}
	obj.Action()
	cycle.DoProcess(obj)
	cmdbk.Start(callBack)
	defer conn.Close()
	for {
		revcivePacket(conn)
		sendPacket(conn)
		time.Sleep(100 * time.Millisecond)
	}
}

func callBack(text string) {
	fmt.Println("file:", text)
	addFile(text)
}

func addFileWrapper() {
	addFile(config.Health.File)
}

func addFile(text string) {
	contents := btext.BParseFile(text)
	if len(contents) != 0 {
		m.Lock()
		defer m.Unlock()
		sendQue.PushBack(contents)
	}
}

func revcivePacket(conn net.Conn) {
	messageBuf := make([]byte, 4058)
	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	messageLen, err := conn.Read(messageBuf)
	if 0 == revcheckErr(err) {
		data := messageBuf[:messageLen]

		message := string(btext.TParseAry(data))
		code, err := dispatch.GetCode(data)
		log.Info("[rev:0x" + fmt.Sprintf("%x", code) + "]->\n" + message)
		if err == nil {
			insertFile(uint(code))
		} else {
			log.Info("エラー")
		}
	}
}

func insertFile(code uint) {
	for _, v := range config.Responce {
		vCode, err := exchangeCode(v.Code)
		if err == nil {
			if vCode == code {
				fmt.Println("responce->", v.File)
				addFile(v.File)
				break
			}
		} else {
			log.Info("変換エラー：" + v.Code)
		}
	}
}

func exchangeCode(codeStr string) (uint, error) {
	if len(codeStr) > 2 && codeStr[0:2] == "0x" {
		val, err := strconv.ParseUint(codeStr, 0, 32)
		if err == nil {
			return uint(val), nil
		}
	}
	if regexp.MustCompile(`[0-9]`).Match([]byte(codeStr)) {
		val, err := strconv.ParseUint(codeStr, 10, 32)
		if err == nil {
			return uint(val), nil
		}
	}
	return 0, fmt.Errorf("error occured")
}

func sendPacket(conn net.Conn) {
	m.Lock()
	defer m.Unlock()
	for sendQue.Len() != 0 {
		message := sendQue.Remove(sendQue.Front())
		switch buff := message.(type) {
		case []byte:
			if config.TimeGrant == true {
				grantTime(buff[36:44])
			}
			conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
			_, err := conn.Write(buff)
			if err != nil {
				println("Write to server failed:", err.Error())
				os.Exit(1)
			}
			log.Info("[Send]->\n" + btext.TParseAry(buff))
		}
	}
}

func revcheckErr(err error) (retVal int) {
	retVal = 0
	if err != nil {
		if strings.Index(err.Error(), "timeout") == -1 {
			log.Info("err!!!", err.Error())
		}
		retVal = -1
	}
	return
}

func grantTime(buffer []byte) {
	if len(buffer) < 8 {
		return
	}
	now := time.Now()

	hi, _ := strconv.ParseInt(strconv.Itoa(int(now.Year())/100), 16, 16)
	lo, _ := strconv.ParseInt(strconv.Itoa(int(now.Year())%100), 16, 16)
	buffer[0] = byte(hi)
	buffer[1] = byte(lo)

	buffer[2] = byte(now.Month())
	buffer[3] = byte(now.Day())
	buffer[4] = byte(now.Hour())
	buffer[5] = byte(now.Minute())
	buffer[6] = byte(now.Second())
}
