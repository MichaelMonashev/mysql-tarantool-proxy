// Принимает и обрабатывает запросы от клиентов.
package server

import (
	"time"

	// сторонние пакеты
	tarantool "github.com/tarantool/go-tarantool"
	"github.com/youtube/vitess/go/mysql"
	//"github.com/youtube/vitess/go/sqltypes"
	//"github.com/youtube/vitess/go/vt/servenv/grpcutils"

	// наши пакеты
	"mtproxy/config"
	. "mtproxy/warn"
)

type Server struct {
	Addr string
	Net  string
	User string
	Pass string

	TarantoolAddr string
	TarantoolNet  string
	TarantoolUser string
	TarantoolPass string
}

func New(c *config.Config) (*Server, error) {

	s := &Server{
		Addr: c.ServerAddr,
		Net:  c.ServerNet,
		User: c.ServerUser,
		Pass: c.ServerPass,

		TarantoolAddr: c.TarantoolAddr,
		TarantoolNet:  c.TarantoolNet,
		TarantoolUser: c.TarantoolUser,
		TarantoolPass: c.TarantoolPass,
	}

	return s, nil
}

func (s *Server) ListenAndServe() error {

	opts := tarantool.Opts{
		Timeout:       50 * time.Millisecond,
		Reconnect:     100 * time.Millisecond,
		MaxReconnects: 3,
		User:          s.TarantoolUser,
		Pass:          s.TarantoolPass,
	}

	addr := s.TarantoolNet + "://" + s.TarantoolAddr

	tarantoolConn, err := tarantool.Connect(addr, opts)
	if err != nil {
		return Errorln("Failed to connect to the Tarantool server on", addr, ":", err)
	}
	defer tarantoolConn.Close()
	Warn("Successfully connected to the Tarantool server on", tarantoolConn.RemoteAddr())

	handler := &Handler{
		tarantoolConn: tarantoolConn,
	}

	//authServer := &mysql.AuthServerNone{}
	authServer := mysql.NewAuthServerStatic()
	authServer.Entries[s.User] = &mysql.AuthServerStaticEntry{
		Password: s.Pass,
		UserData: "",
	}

	// create a Listener.
	listener, err := mysql.NewListener(s.Net, s.Addr, authServer, handler)
	if err != nil {
		return Errorln("Failed to listen to incoming connection:", err)
	} else {
		Warn("Start listening on", listener.Addr())
	}
	defer listener.Close()

	//	if *mysqlSslCert != "" && *mysqlSslKey != "" {
	//		listener.TLSConfig, err = grpcutils.TLSServerConfig(*mysqlSslCert, *mysqlSslKey, *mysqlSslCa)
	//		if err != nil {
	//			log.Fatalf("grpcutils.TLSServerConfig failed: %v", err)
	//			return
	//		}
	//	}
	//	listener.AllowClearTextWithoutTLS = *mysqlAllowClearTextWithoutTLS

	// starts listening.
	listener.Accept()

	return nil
}

func (s *Server) Close() {
	return
}
