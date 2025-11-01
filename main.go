package main

import (
	"fmt"
	"net"
	"strings"
	//"os"
)

// In-Memory database is the actual "Server" in this case.
// we need a client to test out our code

// How redis works>

// we store data in form of key-pair value
// SET/GET  ->  RESP[REddis Searlization protocol]

//How RESP works?

//  EG : SET admin ahmed
//  *3  $3  set  $5 admin  $5  admin

// 1) * -> size of array
// 2) $ -> length of next wore
//

func main() {

	// Writing my server

	fmt.Println("Starting the server !")
	fmt.Println("Listening on port :6379 ......")

	//Start a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	aof.Read(func(value Value) {
		if len(value.array) == 0 {
			return
		}
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]

		if !ok {
			fmt.Println("Invalid Commands", command)
			return
		}

		handler(args)
	})

	//Accepting requests
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	//delayed close
	defer conn.Close()

	for {

		//made a new read function for RESP
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request , expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected length of array > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		// created an instance of our Writer struct and used it to write to our memory
		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid Command", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}
		result := handler(args)
		writer.Write(result)
	}

}
