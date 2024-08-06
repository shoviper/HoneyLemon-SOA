package models

type ClientInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	BirthDate string `json:"birthDate"`
}

type RegisterClient struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Address   string `json:"address"`
	BirthDate string `json:"birthDate"`
	Password  string `json:"password"`
}

type LoginClient struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type DeleteClient struct {
	ID string `json:"id"`
}
