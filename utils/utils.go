package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/arthurmarkus2013/simple-rest-server/database"
)

type Role string

const (
	USER  Role = "user"
	ADMIN Role = "admin"
)

func GenerateToken(username, password string, role Role) (string, error) {
	result, err := ValidateCredentials(username, password)

	if err != nil {
		return "", err
	}

	if !result {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"username":  username,
		"role":      role,
		"createdAt": time.Now(),
		"expiresAt": time.Now().Add(time.Hour),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte("secret"))

	if err != nil {
		return "", errors.New("failed to generate token")
	}

	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tokens (token) VALUES (?)")

	if err != nil {
		return "", err
	}

	defer stmt.Close()

	_, err = stmt.Exec(tokenString)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(token string) (bool, error) {
	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("SELECT token FROM tokens WHERE token = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	var tokenString string

	err = stmt.QueryRow(token).Scan(&tokenString)

	if err != nil {
		return false, err
	}

	if token != tokenString {
		return false, errors.New("invalid token")
	}

	result, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte("secret"), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}), jwt.WithExpirationRequired())

	if err != nil {
		return false, err
	}

	return result.Valid, nil
}

func InvalidateToken(token string) bool {
	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM tokens WHERE token = ?")

	if err != nil {
		return false
	}

	defer stmt.Close()

	result, err := stmt.Exec(token)

	if err != nil {
		return false
	}

	affectedRows, _ := result.RowsAffected()

	return affectedRows == 1
}

func ValidateCredentials(username, password string) (bool, error) {
	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("SELECT password FROM users WHERE username = ?")

	if err != nil {
		return false, err
	}

	defer stmt.Close()

	var passwordHash string

	err = stmt.QueryRow(username).Scan(&passwordHash)

	if err != nil {
		return false, err
	}

	return CheckPassword(password, passwordHash), nil
}

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), 64)

	if err != nil {
		return "", err
	}

	return string(hashBytes), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
