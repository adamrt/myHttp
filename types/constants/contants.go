package constants

const (
	methodName methodId = iota
	PING                // Server check
	CAST                // Send a message
	OBTAIN              // Request Data
	VALID               // Auth
	OMIT                // Delete
	BYE                 // End the connection
)

type methodId int

// VXP/1.0 METHOD URL

const (
	headerName headerId = iota
	requestID
	clientID
	contentType
	connection
)

type headerId int

//Valid headers will be
// RequestID:

type StateId int

const (
	StateName StateId = iota
	STATE_REQLINE
	STATE_HEADERS
	STATE_BODY
	STATE_DONE
)
