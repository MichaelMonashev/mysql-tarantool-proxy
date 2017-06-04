// Принимает и обрабатывает запросы от клиентов.
package server

import (
	// сторонние пакеты
	tarantool "github.com/mialinx/go-tarantool"
	"github.com/youtube/vitess/go/mysql"
	"github.com/youtube/vitess/go/sqltypes"

	// наши пакеты
	"mtproxy/config"
	. "mtproxy/warn"
)

type Server struct {
	Addr          string
	Net           string
	TarantoolAddr string
	TarantoolNet  string
}

func New(config *config.Config) (*Server, error) {

	s := &Server{
		Addr:          config.ServerAddr,
		Net:           config.ServerNet,
		TarantoolAddr: config.TarantoolAddr,
		TarantoolNet:  config.TarantoolNet,
	}

	return s, nil
}

func (s *Server) ListenAndServe() error {

	opts := tarantool.Opts{
		Timeout: 50 * time.Millisecond,
		Reconnect: 100 * time.Millisecond,
		MaxReconnects: 3,
		User: "test",
		Pass: "test",
	}

	tarantoolConn, err := tarantool.Connect(s.TarantoolNet+"://"+s.TarantoolAddr, opts)
	if err != nil {
		return Errorln("Failed to connect to tarantool server:", err)
	}
	defer tarantoolConn.Close()

	handler := &Handler{
		tarantoolConn: tarantoolConn,
	}

	//authServer := &mysql.AuthServerNone{}
	authServer := mysql.NewAuthServerStatic()

	// create a Listener.
	listener, err := mysql.NewListener(s.Net, s.Addr, authServer, handler)
	if err != nil {
		return Errorln("Failed to listen to incoming connection:", err)
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

type Handler struct {
	tarantoolConn *tarantool.Connection
}

func (h *Handler) NewConnection(c *mysql.Conn) {
}

func (h *Handler) ConnectionClosed(c *mysql.Conn) {

}

func (h *Handler) ComQuery(c *mysql.Conn, query []byte) (*sqltypes.Result, error) {

	resp, err := h.tarantoolConn.Eval("box.sql.execute", string(query)) // ToDo убрать string ?
	if err!=nil {
		return nil, err
	}
	if resp.Code != tarantool.OkCode {
		Warn("err==nil, resp.Code != tarantool.OkCode")
		return nil, Errorln(fmt.Sprintf("<%d ERR 0x%x %s>", resp.RequestId, resp.Code, resp.Error))
	}


	return nil, nil

}
