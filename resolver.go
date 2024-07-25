package main

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"

	//	"errors"
	"io"
	"log"
)

const (
	separator = "\r\n"
)

const (
	COMMAND = "COMMAND"
	DOCS = "DOCS"
	GET = "GET"
	SET = "SET"
	PING = "PING"
)

type Resp struct{
	writer *bufio.Writer
}

func newWriter(wr io.Writer)*Resp{
	return &Resp{writer: bufio.NewWriter(wr)}
}

func (w *Resp) ok(args []Value,data *InMemoryData){
	w.writeResponse("+OK\r\n")
}

func (w *Resp) set(args []Value,data *InMemoryData){
	data.save(args[0].str,args[1].str)
	w.writeResponse("+OK\r\n")
}

func (w *Resp) get(args []Value,data *InMemoryData){
	r,err := data.get(args[0].str)
	if err != nil{
		errMsg := "-"+err.Error()+separator
		w.writeResponse(errMsg)
	}else{
		r,ok := r.(string)
		if ok {
			response := createBulkResponse(r)
			w.writeResponse(response)
		}else{
			err := errors.New("unable to convert value to string")
			response := "-"+err.Error()+separator
			w.writeResponse(response)
		}
	}
}
func (w *Resp) hget(args []Value,data *InMemoryData){
	key,field:= args[0].str,args[1].str
	m,err := data.get(key)
	if err != nil{
		response := "-"+err.Error()+separator
		w.writeResponse(response)
	}else{
		m,ok := m.(map[string]interface{})
		if ok{
			value,ok := m[field].(string)
			if ok{
				value := createBulkResponse(value)
				w.writeResponse(value)
			}else{
				err := errors.New("unable to convert the value from the hash")
				response := "-"+err.Error()+separator
				w.writeResponse(response)

			}
		}
	}
}

func (w *Resp) hset(args []Value,data *InMemoryData){
	key,field,value := args[0].str,args[1].str,args[2].str
	r,err := data.get(key)	
	if err != nil{
		newMap := make(map[string]interface{})
		newMap[field]=value	
		data.saveStruct(key,newMap)
	}
	m,ok := r.(map[string]interface{})
	if ok{
		m[field] = value
	}else{
		err := errors.New("Unable to convert the hash from the database.\nNot map structure")
		response := "-"+err.Error()+separator
		w.writeResponse(response)
	}
}

func (w *Resp) mapper(command string,args []Value,data *InMemoryData){
	funcMapper := map[string]func(args []Value, data *InMemoryData){
			"COMMAND":w.ok,
			"DOCS":w.ok,
			"SET":w.set,
			"GET":w.get,
			"HSET":w.hget,
			"HGET":w.hget,
		}
		resolveFunction := funcMapper[command]
		resolveFunction(args,data)
}

func (w *Resp) ResolveRequest(req Value,data *InMemoryData){
	command := req.array[0].str
	w.mapper(command,req.array[1:],data)
	w.writer.Flush()
}	

func createBulkResponse(value string)string{
	lenght := strconv.Itoa(len(value))
	return "$"+lenght+separator+value+separator
}

func (w *Resp) writeResponse(response string)error{
	buff := []byte(response)
	_,err := w.writer.Write(buff)
	if err !=nil{
		return err
	}
	return nil
}
