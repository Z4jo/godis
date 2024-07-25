package main

import (
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"
	"sync"
)

type InMemoryData struct{
	mutex sync.Mutex	
	mainMap map[string]interface{}
}

func (d *InMemoryData) save(key string, value interface{}){
	d.mutex.Lock()
	d.mainMap[key] = value
	d.mutex.Unlock()
}

func (d *InMemoryData) saveStruct(key string, value interface{}){
	d.mutex.Lock()
	d.mainMap[key] = value
	d.mutex.Unlock()
}

func (d *InMemoryData) get(key string)(interface{},error){
	r,ok := d.mainMap[key]
	if ok {
		return r.(string),nil
	}
	return "",errors.New("this key is not existing in the database")
}

func handleConnection(con net.Conn,data *InMemoryData){
	req := newReader(con)
	resp := newWriter(con)
	for {
		value,err := req.ReadType()
		fmt.Printf("%+v\n",value)	
		if err != nil{
			if err.Error() == "EOF"	{
				con.Close()
				break;
			}else{
				log.Fatal()
			}
		}
		fmt.Printf("%+v\n",value)	
		resp.ResolveRequest(value, data)	
		if err !=nil{
			log.Fatal(err)
		}
	}
}

func main(){
	data := &InMemoryData{mainMap: make(map[string]interface{})}	
	l,err := net.Listen("tcp",":6969")
	defer l.Close()
	if err != nil{
		slog.Error("tcp listener failed")
		log.Fatal(err)
	}
	slog.Info("NOw we are listening to 6969")
	for{
		con,err := l.Accept()
		if err != nil{
			slog.Error("connection failed")
			log.Fatal(err)
		}
		defer con.Close()
		go handleConnection(con,data)
	}
}
