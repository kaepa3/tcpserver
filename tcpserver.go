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

func main() {
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
	obj := cycle.CycleProc{Time: config.Health.Time, Flg: true, Action: addFile}
	obj.Action()
	cycle.DoProcess(obj)
	for {
		revcivePacket(conn)
		sendPacket(conn)
		time.Sleep(100 * time.Millisecond)
	}
}

func addFile() {
	contents, err := ioutil.ReadFile(config.Health.File) // ReadFileの戻り値は []byte
	if err == nil {
		sendQue.PushBack(contents)
	}
}

func revcivePacket(conn net.Conn) {
	messageBuf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	messageLen, err := conn.Read(messageBuf)
	if 0 == revcheckErr(err) {
		message := string(messageBuf[:messageLen])
		fmt.Println("comming -", message)
	}
}

func sendPacket(conn net.Conn) {
	if sendQue.Len() != 0 {
		message := sendQue.Remove(sendQue.Front())
		conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
		switch buff := message.(type) {
		case []byte:
			conn.Write(buff)
			fmt.Println("send -", buff)
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
