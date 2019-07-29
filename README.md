# Current Features
  - server creates database tables
  - multiple peers can register to server.
  - server accepts or rejects a peer.
  - server logs results of registering requests to database.

# Run
1. postgres settings needs to be configured in server.go

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "****"
	dbname   = "picusbc"
)

2. server.go [portno]
3. peer.go [serverip] [serverport] 
