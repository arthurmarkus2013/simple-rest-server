package routes

type User struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type Credentials struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Movie struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	ReleaseYear int    `json:"release_year" binding:"required"`
}
