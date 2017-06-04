package config

import (
	"flag"
	"net"
	"os"
	"path/filepath"

	// наши пакеты
	. "mtproxy/warn"
)

type Config struct {
	ServerAddr string
	ServerNet  string

	TarantoolAddr string
	TarantoolNet  string
	TarantoolUser string
	TarantoolPass string
}

// LoadConfig читает значения по-умолчанию, аргументы командной строки, парсит и загружает файл конфигурации. Возвращает конфигурацию в виде структуры.
// Address could be specified in following ways:
//
// TCP connections:
// - tcp://192.168.1.1:3013
// - tcp://my.host:3013
// - tcp:192.168.1.1:3013
// - tcp:my.host:3013
// - 192.168.1.1:3013
// - my.host:3013
// Unix socket:
// - unix:///abs/path/tnt.sock
// - unix:path/tnt.sock
// - /abs/path/tnt.sock  - first '/' indicates unix socket
// - ./rel/path/tnt.sock - first '.' indicates unix socket
// - unix/:path/tnt.sock  - 'unix/' acts as a "host" and "/path..." as a port

func LoadConfig() (*Config, error) {

	// длинные параметры командной строки удобнее набирать, когда они пишутся через тире.
	serverAddr := flag.String("server-addr", "localhost:3000", `The address for mysql client requests. Example: -server-addr="192.168.1.1:3000".`)
	tarantoolAddr := flag.String("tarantool-addr", "localhost:3001", `The address of tarantool server. Example: -tarantool-addr="192.168.1.1:3013".`)
	// ToDo добавить логин и пароль для прокси и тарантула

	flag.Parse()

	var c Config
	var err error

	c.ServerNet, c.ServerAddr, err = parseAddr(*serverAddr)
	if err != nil {
		return nil, Errorln("Error in address for mysql client requests.", err)
	}

	c.TarantoolNet, c.TarantoolAddr, err = parseAddr(*tarantoolAddr)

	if err != nil {
		return nil, Errorln("Error in address of tarantool server.", err)
	}

	return &c, nil

}

func parseAddr(address string) (string, string, error) {
	if address == "" {
		return "", "", Errorln("empty address")
	}

	network := "tcp"

	// Unix socket connection
	if address[0] == '.' || address[0] == '/' {
		network = "unix"

	} else if len(address) > 7 && address[0:7] == "unix://" {
		network = "unix"
		address = address[7:]

	} else if len(address) > 5 && address[0:5] == "unix:" {
		network = "unix"
		address = address[5:]

	} else if len(address) > 6 && address[0:6] == "unix/:" {
		network = "unix"
		address = address[6:]

	} else if len(address) > 6 && address[0:6] == "tcp://" {
		address = address[6:]

	} else if len(address) > 4 && address[0:4] == "tcp:" {
		address = address[4:]
	}

	var err error

	// try to resolve address
	if network == "tcp" {
		_, err = net.ResolveTCPAddr(network, address)

	} else if network == "unix" {
		address, err = filepath.Abs(address)
		if err == nil {
			var f *os.File
			f, err = os.Open(address)
			if err == nil {
				f.Close()
			}
		}
	}

	return network, address, err

}
