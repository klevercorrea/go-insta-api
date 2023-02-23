package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/klevercorrea/go-insta-api/pkg/instagram"
	"github.com/klevercorrea/go-insta-api/pkg/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("INSTAGRAM_USERNAME")
	password := os.Getenv("INSTAGRAM_PASSWORD")
	ctx := context.Background()

	sessionString, err := instagram.Login(ctx, username, password)
	if err != nil {
		panic(err)
	}

	sessionBytes, err := base64.StdEncoding.DecodeString(sessionString)
	if err != nil {
		panic(err)
	}

	insta, err := utils.ImportFromBytes(sessionBytes)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Logged in as %s\n", insta.Account.Username)

	utils.ShortSleep()

	profile, err := instagram.GetProfileData(ctx, insta, insta.Account.Username)
	if err != nil {
		panic(err)
	}

	utils.ShortSleep()

	err = instagram.GetFollowers(ctx, insta, profile.FollowerCount)
	if err != nil {
		fmt.Println("Error getting followers: ", err)
		panic(err)
	}
}
