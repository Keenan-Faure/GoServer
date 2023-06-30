package objects

import "time"

// Request/Response Objects
type RequestBodyChirp struct {
	Body string `json:"body"`
}

type RequestBodyUser struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type RequestBodyLogin struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type ResponseUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type ResponseUserLogon struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
type ResponseRefreshToken struct {
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
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RevokedTokens map[time.Time]string `json:"revoked_tokens"`
}
