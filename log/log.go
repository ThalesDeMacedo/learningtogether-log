package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net"
	"net/http"
)

var client = http.Client{}

func Info(msg string, requestId uuid.UUID)  {
	sendLog(new("info", msg, requestId))
}

func Error(msg string, requestId uuid.UUID)  {
	sendLog(new("error", msg, requestId))
	panic(fmt.Errorf(msg))
}

func getIpMachine() net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			return ip
		}
	}
	return nil
}

func sendLog(l log)  {
	marshal, err := json.Marshal(l)
	if err != nil {
		marshal = []byte("Erro ao fazer parse do log")
	}

	r, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:5555/gelf", bytes.NewReader(marshal))

	_, err = client.Do(r)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func new(level string, msg string, requestId uuid.UUID) log {
	return log{
		Level: level,
		Message: msg,
		Host:    getIpMachine().String(),
		RequestId: requestId.String(),
	}
}

type log struct {
	Message   string `json:"message"`
	Host      string `json:"host"`
	RequestId string `json:"request_id"`
	Level     string `json:"level"`
	Application string `json:"application"`
	Language string `json:"language"`
	framework string `json:"framework"`
	sequence string `json:"sequence"`
	instance string `json:"instance"`
}