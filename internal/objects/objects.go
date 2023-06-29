package objects

// Request/Response Objects
type RequestBodyChirp struct {
	Body string `json:"body"`
}

type RequestBodyUser struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RequestBodyLogin struct {
	Password       string `json:"password"`
	Email          string `json:"email"`
	ExpiresSeconds int    `json:"expires_in_seconds"`
}

type ResponseUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type ResponseUserLogon struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

// Data Objects
type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

// Database Structure
type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}
