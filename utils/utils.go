package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"

	"errors"
	"log/slog"

	"golang.org/x/crypto/bcrypt"

	"github.com/arthurmarkus2013/simple-rest-server/database"
)

type Role string

const (
	NONE  Role = ""
	USER  Role = "user"
	ADMIN Role = "admin"
)

const tokenSecret = "374gfech98h93f9cmkvbztggfh09pw374ggf"

func GenerateToken(username, password string) (string, error) {
	result, role, err := ValidateCredentials(username, password)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return "", err
	}

	if !result {
		return "", errors.New("invalid credentials")
	}

	claims := jwt.MapClaims{
		"sub":   username,
		"roles": role,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return "", errors.New("failed to generate token")
	}

	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO tokens (token, ttl) VALUES (?, ?)")

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return "", err
	}

	defer stmt.Close()

	_, err = stmt.Exec(tokenString, claims["exp"])

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return "", err
	}

	return tokenString, nil
}

func ValidateToken(token string) (bool, Role, error) {
	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("SELECT token FROM tokens WHERE token = ?")

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return false, NONE, err
	}

	defer stmt.Close()

	var tokenString string

	err = stmt.QueryRow(token).Scan(&tokenString)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return false, NONE, err
	}

	if token != tokenString {
		return false, NONE, errors.New("invalid token")
	}

	result, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS512.Alg()}), jwt.WithExpirationRequired())

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return false, NONE, err
	}

	claims := result.Claims.(jwt.MapClaims)
	role := Role(claims["roles"].(string))

	return result.Valid, role, nil
}

func InvalidateToken(token string) bool {
	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("DELETE FROM tokens WHERE token = ?")

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return false
	}

	defer stmt.Close()

	result, err := stmt.Exec(token)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return false
	}

	affectedRows, _ := result.RowsAffected()

	return affectedRows == 1
}

func ValidateCredentials(username, password string) (bool, Role, error) {
	db := database.OpenDatabase()

	defer db.Close()

	stmt, err := db.Prepare("SELECT password, role FROM users WHERE username = ?")

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return false, NONE, err
	}

	defer stmt.Close()

	var passwordHash, role string

	err = stmt.QueryRow(username).Scan(&passwordHash, &role)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return false, NONE, err
	}

	return CheckPassword(password, passwordHash), Role(role), nil
}

var salt = "j403fjJ)FJ3jf9j))!Fj9f!IR9xxss07hh"

func HashPassword(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(salt+password), bcrypt.DefaultCost)

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())

		return "", err
	}

	return string(hashBytes), nil
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(salt+password))

	if err != nil {
		slog.Error("something went wrong", "error", err.Error())
	}

	return err == nil
}

func PurgeExpiredTokens(interval time.Duration) {
	timer := time.NewTicker(interval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			go func() {
				db := database.OpenDatabase()

				defer db.Close()

				stmt, err := db.Prepare("DELETE FROM tokens WHERE ttl < ?")

				if err != nil {
					slog.Error("something went wrong", "error", err.Error())

					return
				}

				defer stmt.Close()

				result, err := stmt.Exec(time.Now().Unix())

				if err != nil {
					slog.Error("something went wrong", "error", err.Error())

					return
				}

				affectedRows, _ := result.RowsAffected()

				slog.Info("purged expired tokens", "count", affectedRows)
			}()
		}
	}
}
