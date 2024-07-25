package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Req struct{
	reader *bufio.Reader
}

func newReader(rd io.Reader) *Req{
	return &Req{reader: bufio.NewReader(rd)}
}

func (r *Req) readSize()(int,error){
	size,err := r.reader.ReadBytes('\r')
	if err != nil{
		return 0,err
	}
	_,err = r.reader.ReadByte()
	if err != nil{
		return 0,err
	}
	size = bytes.TrimSuffix(size,[]byte("\r"))
	reformedString := string(size)
	return strconv.Atoi(reformedString)
}

func (r *Req) readArray()(Value,error){
	nElements,err := r.readSize()
	if err != nil{
		return Value{},err
	}
	var values []Value
	for i := 0;i<nElements;i++{
		v,err := r.ReadType()
		if err != nil{
			return Value{},err
		}
		values = append(values, v)
	}
	return Value{array: values},nil
}

func (r *Req) readString()(Value,error){
	value,err := r.extractValue()
	if err != nil{
		return Value{},err
	}
	//WARN: check 0 len if true throw error specific
		return Value{str: string(value)},nil	
	}

func (r *Req) readBulkString()(Value,error){
	size,err := r.readSize()
	if err != nil{
		return Value{},err
	}
	buf := make([]byte,size+2)

	n,err := r.reader.Read(buf)
	if n != size {
		//TODO: return error 
	}
	if err != nil{
		return Value{},err
	}
	value := string(buf[:len(buf)-2])
	return Value{str:value},nil
}

func (r *Req) ReadType()(Value,error){
	typeOfStructure,err := r.reader.ReadByte()

	if err != nil{
		log.Fatal(err)
	}
	var ret Value
	switch string(typeOfStructure){
		case "*": 		
			ret,err = r.readArray()			
			
		case "+":
			ret,err = r.readString()
		case ":":
		case "$":
			ret,err = r.readBulkString()
		case "-":
	}
	return ret,err
} 

func (r *Req) extractValue()([]byte,error){
	value,err := r.reader.ReadBytes('\r')			
	if err !=nil {
		return nil,err
	}
	_,err =r.reader.ReadByte()
	if err !=nil {
		return nil,err
	}
	value = bytes.TrimSuffix(value,[]byte("\r"))
	return value,nil
}
