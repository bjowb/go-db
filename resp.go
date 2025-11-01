package main

import (
	"bufio"
	"fmt"
	"io"
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

// initialize struct using bufio library
type Resp struct {
	reader *bufio.Reader
}

// made new object of Resp struct for using
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// used function to readLine in our command
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {

		//Read byte by byte
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		//increment pointer to next byte and add it to array
		n += 1
		line = append(line, b)

		//break when last 2 elements are CRLF :- \r\f
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}

	//editted out last 2 elements cause they are CRLF
	return line[:len(line)-2], n, nil
}

// used to readInteger in our command
func (r *Resp) readInteger() (x int, n int, err error) {

	//read the string till \r
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	//convert the number given into int and return it
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(i64), n, nil
}

// # parsing starts from here

// 1st function call
func (r *Resp) Read() (Value, error) {

	// read the type of data
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	//on basis of read, choose function to call
	//rn, we implemented only 2
	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Println("Unknown type : %v", string(_type))
		return Value{}, nil
	}
}

// #1)  readArray :-
//    Skip the first byte because we have already read it in the Read method.
//    Read the integer that represents the number of elements in the array.
//    Iterate over the array and for each line, call the Read method to parse the type according to the character at the beginning of the line.
//    With each iteration, append the parsed value to the array in the Value object and return it.

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	//read length of array
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	//for each line, parse and read the value
	v.array = make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// add parsed value to array
		v.array[i] = val
	}

	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	if length < 0 {
		return v, nil
	}

	buff := make([]byte, length)
	_, err = io.ReadFull(r.reader, buff)
	if err != nil {
		return v, err
	}

	v.bulk = string(buff)

	//read the trailing CRLF :- 

	r.readLine()

	return v, nil
}

// Now we convert the given format of data into bytes to store in memory
// this process is called marshalling
// here we convert different data type to bytes
func (v Value) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}
	}
}

func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes

}

func (v Value) marshalArray() []byte {
	len := len(v.array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshalError() []byte {
	var bytes []byte

	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshalNull() []byte {
	return []byte("$-1\r\n")
}

// Now we write a Writer struct which will help us write data to memoryy
type Writer struct {
	writer io.Writer
}

// create instance of upper class
func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// now write all the data
func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
