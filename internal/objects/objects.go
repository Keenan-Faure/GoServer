package objects

import "time"

// Webhook Request Format
type WebhookRequest struct {
	Event string      `json:"event"`
	Data  WebhookUser `json:"data"`
}

type WebhookUser struct {
	UserID int `json:"user_id"`
}

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

type BaseResponse struct {
	Body string `json:"body"`
}

type ResponseUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type ResponseUserLogon struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}
type ResponseRefreshToken struct {
	Token string `json:"token"`
}

// Data Objects
type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	ID          int    `json:"id"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

// Database Structure
type DBStructure struct {
	Chirps        map[int]Chirp        `json:"chirps"`
	Users         map[int]User         `json:"users"`
	RevokedTokens map[time.Time]string `json:"revoked_tokens"`
}

// Docs Endpoints
type Endpoints struct {
	Status      bool             `json:"status"`
	Description string           `json:"description"`
	Routes      map[string]Route `json:"routes"`
	Version     string           `json:"version"`
	Time        time.Time        `json:"time"`
}

type Route struct {
	Description   string            `json:"description"`
	Supports      []string          `json:"supports"`
	Params        map[string]Params `json:"params"`
	AcceptsData   bool              `json:"accepts_data"`
	Format        any               `json:"format"`
	Authorization string            `json:"auth"`
}

type Params struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
