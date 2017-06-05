package config

import (
	"flag"
	"net"
	"os"
	"path/filepath"

	// сторонние пакеты
	toml "github.com/pelletier/go-toml"

	// наши пакеты
	. "mtproxy/warn"
)

const (
	defaultServerAddr    = "127.0.0.1:3000"
	defaultTarantoolAddr = "127.0.0.1:3301"
)

type Config struct {
	ServerAddr string
	ServerNet  string
	ServerUser string
	ServerPass string

	TarantoolAddr string
	TarantoolNet  string
	TarantoolUser string
	TarantoolPass string
}

// LoadConfig читает значения по-умолчанию, аргументы командной строки, парсит и загружает файл конфигурации. Возвращает конфигурацию в виде структуры.
// Address could be specified in following ways:
//
// TCP connections:
// - tcp://192.168.1.1:3301
// - tcp://my.host:3301
// - tcp:192.168.1.1:3301
// - tcp:my.host:3301
// - 192.168.1.1:3301
// - my.host:3301
// Unix socket:
// - unix:///abs/path/tnt.sock
// - unix:path/tnt.sock
// - /abs/path/tnt.sock  - first '/' indicates unix socket
// - ./rel/path/tnt.sock - first '.' indicates unix socket
// - unix/:path/tnt.sock  - 'unix/' acts as a "host" and "/path..." as a port

func LoadConfig() (*Config, error) {

	flag := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// длинные параметры командной строки удобнее набирать, когда они пишутся через тире.
	serverAddr := flag.String("server-addr", "", `Address for mysql client requests. Default: `+defaultServerAddr+`. Example: -server-addr="127.0.0.1:3000".`)
	serverUser := flag.String("server-user", "", `User name for mysql client requests. Example: -server-user="admin".`)
	serverPass := flag.String("server-pass", "", `Password for mysql client requests. Example: -server-pass="1234567".`)
	tarantoolAddr := flag.String("tarantool-addr", "", `Address of tarantool server. Default: `+defaultTarantoolAddr+`. Example: -tarantool-addr="127.0.0.1:3301".`)
	tarantoolUser := flag.String("tarantool-user", "", `User for login to tarantool server. Example: -tarantool-user="admin".`)
	tarantoolPass := flag.String("tarantool-pass", "", `Password to use when connecting to tarantool server. Example: -tarantool-pass="1234567".`)
	config := flag.String("config", "", `Path to TOML configuration file. Example: -config="/etc/mtproxy/config.toml".`)

	flag.Parse(os.Args[1:])

	var c Config
	var err error

	if *config != "" {
		// load TOML data by filename
		tree, err := toml.LoadFile(*config)
		if err != nil {
			return nil, Errorln("Can not load configuration from file", config, ":", err)
		}

		// load configuration from TOML file
		c.ServerNet, c.ServerAddr, err = parseAddr(tree.Get("server.addr").(string))
		if err != nil {
			return nil, Errorln("Error in configuration file in address for mysql client requests.", err)
		}
		c.ServerUser = tree.Get("server.user").(string)
		c.ServerPass = tree.Get("server.pass").(string)

		c.TarantoolNet, c.TarantoolAddr, err = parseAddr(tree.Get("tarantool.addr").(string))
		if err != nil {
			return nil, Errorln("Error in configuration file in address of tarantool server.", err)
		}
		c.TarantoolUser = tree.Get("tarantool.user").(string)
		c.TarantoolPass = tree.Get("tarantool.pass").(string)
	}

	// set defaults
	if c.ServerAddr == "" {
		c.ServerNet, c.ServerAddr, err = parseAddr(defaultServerAddr)
		if err != nil {
			return nil, Errorln("Error in default address for mysql client requests.", err)
		}
	}

	if c.TarantoolAddr == "" {
		c.TarantoolNet, c.TarantoolAddr, err = parseAddr(defaultTarantoolAddr)
		if err != nil {
			return nil, Errorln("Error in default address of tarantool server.", err)
		}
	}

	// override configuration from command line

	if *serverAddr != "" {
		c.ServerNet, c.ServerAddr, err = parseAddr(*serverAddr)
		if err != nil {
			return nil, Errorln("Error command line args in address for mysql client requests.", err)
		}
	}
	if *serverUser != "" {
		c.ServerUser = *serverUser
	}
	if *serverPass != "" {
		c.ServerPass = *serverPass
	}

	if *tarantoolAddr != "" {
		c.TarantoolNet, c.TarantoolAddr, err = parseAddr(*tarantoolAddr)
		if err != nil {
			return nil, Errorln("Error command line args in address of tarantool server.", err)
		}
	}
	if *tarantoolUser != "" {
		c.TarantoolUser = *tarantoolUser
	}
	if *tarantoolPass != "" {
		c.TarantoolPass = *tarantoolPass
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
