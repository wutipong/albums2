package types

type Album struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Assets []Asset `json:"assets"`
}

type Asset struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
}

type User struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

type Session struct {
	ID string `json:"id"`
	// expiresAt: Date;

	// createdAt: Date;
	// updatedAt: Date;
}
