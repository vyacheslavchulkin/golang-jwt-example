package controllers

import (
	m "../models"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"time"
)

type Token struct {
	UserUUID string
	jwt.StandardClaims
}

func Auth(w http.ResponseWriter, r *http.Request) {
	filter := bson.D{primitive.E{Key: "_id", Value: r.FormValue("uuid")}}
	user, _ := m.FindOneUser(filter)
	if user == nil || !checkPassword(user, r.FormValue("password")) {
		invalidPassword(w)
		return
	}
	token, err := createTokens(user)
	if err != nil {
		serverError(w)
		log.Fatal(err)
		return
	}
	response := responseJson{
		Status:   "success",
		Message:  "New tokens was created",
		Response: token,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

func AuthRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken := r.FormValue("refresh_token")

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_TOKEN_SECRET")), nil
	})

	if err != nil {
		badRequest(w)
		return
	}

	_, ok := token.Claims.(jwt.Claims)
	if !ok && !token.Valid {
		unauthorized(w)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		unauthorized(w)
		return
	}

	refreshUuid, ok := claims["refresh_uuid"].(string)
	if !ok {
		unprocessableEntity(w)
		return
	}

	userUUID, ok := claims["user_uuid"].(string)
	if !ok {
		unprocessableEntity(w)
		return
	}

	filter := bson.D{primitive.E{Key: "_id", Value: userUUID}}
	user, err := m.FindOneUser(filter)
	if err != nil {
		unprocessableEntity(w)
		return
	}

	if !checkRefreshToken(user, refreshUuid) {
		unauthorized(w)
		return
	}

	cleanRefreshToken(user)

	newTokens, err := createTokens(user)
	if err != nil {
		serverError(w)
		log.Fatal(err)
		return
	}

	response := responseJson{
		Status:   "success",
		Message:  "New tokens was created",
		Response: newTokens,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(response)
}

func checkPassword(user *m.User, password string) bool {
	encryptionErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return encryptionErr == nil
}

func checkRefreshToken(user *m.User, refreshToken string) bool {
	encryptionErr := bcrypt.CompareHashAndPassword([]byte(user.RefreshToken), []byte(refreshToken))
	return encryptionErr == nil
}

func createTokens(user *m.User) (*responseToken, error) {
	accessSecret := []byte(os.Getenv("ACCESS_TOKEN_SECRET"))
	refreshSecret := []byte(os.Getenv("REFRESH_TOKEN_SECRET"))

	accessClaims := jwt.MapClaims{}
	accessClaims["access_token_uuid"] = uuid.New().String()
	accessClaims["admin"] = user.Admin
	accessClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	accessClaims["name"] = user.Name
	accessClaims["user_uuid"] = user.ID
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims).SignedString(accessSecret)
	if err != nil {
		return nil, err
	}

	refreshUUID := uuid.New().String()
	refreshClaims := jwt.MapClaims{}
	refreshClaims["refresh_uuid"] = refreshUUID
	refreshClaims["user_uuid"] = user.ID
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims).SignedString(refreshSecret)
	if err != nil {
		return nil, err
	}

	token := &responseToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	setRefreshTokenToUser(user, refreshUUID)
	return token, nil
}

func setRefreshTokenToUser(user *m.User, refreshToken string) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	filter := bson.M{"_id": user.ID}
	update := bson.D{{"refresh_uuid", string(hash)}}
	m.FindOneUserAndUpdate(filter, update)
}

func cleanRefreshToken(user *m.User) {
	setRefreshTokenToUser(user, "")
}
