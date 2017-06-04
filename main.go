package main

import (

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
