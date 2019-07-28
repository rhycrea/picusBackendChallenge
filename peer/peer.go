package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"
)

const MIN = 1
const MAX = 1000000

func random() string {
	return strconv.Itoa(rand.Intn(MAX-MIN) + MIN)
}

//func handleConnection(c net.Conn) {
//	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
//	for {
//		netData, err := bufio.NewReader(c).ReadString('\n')
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//
//		temp := strings.TrimSpace(string(netData))
//		if temp == "shutdown" {
//			break
//		}
//
//		result := strconv.Itoa(random()) + "\n"
//		c.Write([]byte(string(result)))
//	}
//	c.Close()
//}

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Please provide an IP and a port number respectively.")
		return
	}
	rand.Seed(time.Now().Unix())

	//send server a request
	conn, err := net.Dial("tcp4", args[1] + ":" + args[2])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()


	enc := json.NewEncoder(conn)

	register_request := map[string]interface{}{
		"msg": "register",
		"peer_id":  random(),
		"token":  random(),
		"server_ip":  args[1],
	}

	if err := enc.Encode(&register_request); err != nil {
		log.Println(err)
	}
	conn.Close()
}