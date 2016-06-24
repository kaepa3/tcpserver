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

// ConfigData is json format
type ConfigData struct {
	Port string `json:"port"`
}

func main() {
	setting := settingGet("config.json")
	tcpAddr, err := net.ResolveTCPAddr("tcp", setting.Port)
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
	que := list.New()
	fmt.Println("client accept!")
	obj := cycle.CycleProc{Time: 1000, Flg: true, Action: addFile}
	cycle.DoProcess(obj)
	for {
		revcivePacket(conn)
		sendPacket(conn, que)
		time.Sleep(100 * time.Millisecond)
	}
}

func addFile() {
	fmt.Println("call back")
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

func sendPacket(conn net.Conn, que *list.List) {
	if que.Len() != 0 {
		message := que.Remove(que.Front())
		conn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
		switch buff := message.(type) {
		case []byte:
			conn.Write(buff)
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
