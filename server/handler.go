// Принимает и обрабатывает запросы от клиентов.
package server

import (
	"bytes"
	"fmt"

	// сторонние пакеты
	tarantool "github.com/tarantool/go-tarantool"
	"github.com/youtube/vitess/go/mysql"
	"github.com/youtube/vitess/go/sqltypes"
	//"github.com/youtube/vitess/go/vt/servenv/grpcutils"

	// наши пакеты
	. "mtproxy/warn"
)

type Handler struct {
	tarantoolConn *tarantool.Connection
}

func (h *Handler) NewConnection(c *mysql.Conn) {
	//Warn(1)
}

func (h *Handler) ConnectionClosed(c *mysql.Conn) {
	//Warn(2)
}

func (h *Handler) ComQuery(c *mysql.Conn, query []byte) (*sqltypes.Result, error) {

	Warn(string(query), query)

	//strings.ToLower(string(query))

	if bytes.Equal(query, []byte("select @@version_comment limit 1")) {
		Warn("version_comment")

		return &sqltypes.Result{}, nil
	}

	//resp, err := h.tarantoolConn.Eval("box.sql.execute", string(query))
	resp, err := h.tarantoolConn.Call17("box.sql.execute", []string{string(query)}) // ToDo убрать string() ?

	if err != nil {
		return nil, mysql.NewSQLErrorFromError(err)
	}

	Warn(resp)

	if resp.Code != tarantool.OkCode {
		Warn("err==nil, resp.Code != tarantool.OkCode")
		return nil, Errorln(fmt.Sprintf("<%d ERR 0x%x %s>", resp.RequestId, resp.Code, resp.Error))
	}

	Warn(33)
	Warn(len(resp.Data))
	for i, v := range resp.Data {
		Warn(i, v)
		for ii, vv := range v.([]interface{}) {
			Warn(ii, vv)
		}
	}
	Warn(33)
	return &sqltypes.Result{
	//		RowsAffected: 123,
	//		InsertID:     123456789,
	//		Fields: []*querypb.Field{
	//			{
	//				Name: "user",
	//				Type: querypb.Type_VARCHAR,
	//			},
	//			{
	//				Name: "user_data",
	//				Type: querypb.Type_VARCHAR,
	//			},
	//		},
	//		Rows: [][]sqltypes.Value{
	//			{
	//				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte(c.User)),
	//				sqltypes.MakeTrusted(querypb.Type_VARCHAR, []byte(c.UserData.Get().Username)),
	//			},
	//		},
	}, nil
}
