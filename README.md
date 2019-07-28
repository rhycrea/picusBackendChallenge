# Current Features
  - server creates database tables
  - multiple peers can register to server.
  - server accepts or rejects a peer.
  - server logs results of registering requests to database.

# Picus Backend Challenge 

We need to create a system which has a server and a minimum of two peers. The peers should connect and register to the server by a token. The server should start sending a bunch of job to the peer after a peer registering. The peer that was received the jobs should send results to the server. The server should write results to PostgreSQL DB. Meanwhile, we need a registration log for login and logout actions on PostgreSQL, so we can trace how long the peers remain login. The system must support multiple peers at the same time. If we stop the service by pressing Control + C keys combination, you must send shutdown message to all peers. You'd better write this shutdown to the log table with shutdown reason.

##Examples of Messaging

Below, you can see some JSON message examples contents.

###Register (Peer to Server )

```json
{ "msg": "register", "peer_id": "...", "token": "...", "server_ip": "..." }
```

###Registering was Accepted (Server to Peer)
```json
{ "msg": "register", "peer_id": "...", "token": "...", "status": "accepted" }
```

###Registering was Rejected (Server to Peer)

```json
{ "msg": "register", "peer_id": "...", "token": "...", "status": "rejected" }
```

###Send Jobs (Server to Peer)

```json
{ "msg": "register",  "ticket": "126743", "peer_id": "...", "token": "...", "job": { "name":"List Files", "os":"linux", "distro": "debian", "command": "ls -al" } }
```

```json
{ "msg": "register",  "ticket": "864678", "peer_id": "...", "token": "...", "job": {  "name":"Free Memory Space", "os":"linux", "distro": "centos", "command": "free -h" } }
```


###Result (Peer to Server)
```json
{ "msg": "result",  "ticket": "126743", "peer_id": "...", "token": "...", "result": "file1 file2 file3"}
```
```json
{ "msg": "result",  "ticket": "864678", "peer_id": "...", "token": "...", "result": "total:7.4G used: 514M free:1.6G"}
```

Shutdown (Server to All Peers)
```json
{ "msg": "shutdown", "token": "...", "reason":"Upgrade server"}
```
