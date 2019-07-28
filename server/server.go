package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net"
	"os"
	"time"
)

//TODO: postgres settings needs to be configured
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "****"
	dbname   = "picusbc"
)

//JSON structures.
type RegisterRequest struct {
	msg   	string
	peer_id int
	token 	int
	server_ip string
}

type RegisterResponse struct {
	msg   	string
	peer_id string
	token 	string
	status 	string
}

type Job struct {
	name   	string
	os 		string
	distro 	string
	command string
}

type SendJob struct {
	msg   	string
	ticket	string
	peer_id string
	token 	string
	job		Job
}

type Result struct {
	msg   	string
	ticket	string
	peer_id string
	token 	string
	result	string
}

type OnlinePeer struct {
	peer_id string
	token 	string
}


var (
	err 	error //postgres connection, ops and error handling
	db 		sql.DB
	peers	[]OnlinePeer
)
func isPeerValid(x string, token string) bool {
	for _, k := range peers {
		if x == k.peer_id || token == k.token {
			return false
		}
	}
	return true
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		dec := json.NewDecoder(bufio.NewReader(c))
		enc := json.NewEncoder(c)

		var p map[string]interface{} //incoming packet from peer.
		var answer interface{} //outgoing packet to peer.

		if err := dec.Decode(&p); err != nil { //read the incoming packet
			log.Println(err)
			return
		}

		if p["msg"] == "register" {
			var isAccepted bool
			if isPeerValid(p["peer_id"].(string), p["token"].(string)) {
				answer = map[string]interface{}{
					"msg": "register",
					"peer_id":  p["peer_id"].(string),
					"token":  p["token"].(string),
					"status":  "accepted",
				}
				isAccepted = true
				peers = append(peers, OnlinePeer{p["peer_id"].(string), p["token"].(string)})
			} else {
				answer = map[string]interface{}{
					"msg": "register",
					"peer_id":  p["peer_id"],
					"token":  p["token"],
					"status":  "rejected",
				}
				isAccepted = false
			}
			persistAuthLog(p["peer_id"].(string), p["token"].(string), isAccepted)
		}

		if err := enc.Encode(&answer); err != nil {
			log.Println(err)
		}
	}
	c.Close()
}

func createTables() {
	//DB ops
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStatement := `
		CREATE TABLE results (
		  id SERIAL PRIMARY KEY,
		  peer_id TEXT,
		  timestamp TIMESTAMP,
		  token TEXT,
		  ticket VARCHAR ,
		  job JSON,
		  result JSON
		);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		//panic(err)
	}

	sqlStatement = `
		CREATE TABLE auth_logs (
		  id SERIAL PRIMARY KEY,
		  peer_id TEXT,
		  timestamp TIMESTAMP,
		  token TEXT,
		  isRegisteringAccepted BOOLEAN
		);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		//panic(err)
	}

	sqlStatement = `
		CREATE TABLE shutdown_logs (
		  id SERIAL PRIMARY KEY,
		  peer_id_list TEXT,
		  timestamp TIMESTAMP,
		  reason TEXT
		);`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		//panic(err)
	}
}

func persistAuthLog(peer_id string, token string, isRegistered bool) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	sqlStatement := `
	INSERT INTO auth_logs (peer_id, timestamp, token, isregisteringaccepted)
	VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(sqlStatement, peer_id, time.Now(), token, isRegistered)
	if err != nil {
		panic(err)
	}
}

func main() {
	createTables()

	//TCP ops
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide a port number!")
		return
	}

	PORT := ":" + arguments[1]
	listen, err := net.Listen("tcp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()

	//listen for incoming connections.
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go handleConnection(conn)
	}
}