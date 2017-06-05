package main

import (
	//"time"

	// сторонние пакеты
	//tarantool "github.com/tarantool/go-tarantool"

	// наши пакеты
	"mtproxy/config"
	"mtproxy/server"
	. "mtproxy/warn"
	//	// pprof
	//	"log"
	//	"net/http"
	//	_ "net/http/pprof"
	//	"runtime"
)

func main() {
	//	runtime.SetBlockProfileRate(1)
	//	runtime.SetMutexProfileFraction(1)
	//
	//	go func() {
	//		log.Println(http.ListenAndServe("192.168.1.101:8000", nil))
	//	}()

	//	opts := tarantool.Opts{
	//		Timeout:       50 * time.Millisecond,
	//		Reconnect:     100 * time.Millisecond,
	//		MaxReconnects: 3,
	//		User:          "test",
	//		Pass:          "test",
	//	}
	//
	//	tarantoolConn, err := tarantool.Connect("tcp://127.0.0.1:3301", opts)
	//	if err != nil {
	//		Warn("Failed to connect to tarantool server:", err)
	//		return
	//	}
	//	defer tarantoolConn.Close()
	//
	//	//resp, err := tarantoolConn.Eval("box.sql.execute('CREATE TABLE table1 (column1 INTEGER PRIMARY KEY, column2 VARCHAR(100));')", []interface{}{})
	//	//resp, err := tarantoolConn.Eval("box.sql.execute", []string{"CREATE TABLE table1 (column1 INTEGER PRIMARY KEY, column2 VARCHAR(100));"})
	//	resp, err := tarantoolConn.Call("box.sql.execute", []string{"CREATE TABLE table1 (column1 INTEGER PRIMARY KEY, column2 VARCHAR(100));"})
	//	if err != nil {
	//		Warn(1, err)
	//		return
	//	}
	//
	//	Warn(2, resp)
	//
	//	return

	// загружаем конфиг и аргументы командной строки
	conf, err := config.LoadConfig()
	if err != nil {
		Fatal(err)
	}

	// создаём на основании данных конфига объект server
	server, err := server.New(conf)
	if err != nil {
		Fatal(err)
	}
	defer server.Close()

	// обрабатываем команды от клиентов
	Fatal(server.ListenAndServe())
}
