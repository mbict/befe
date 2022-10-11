package dsl

import (
	"github.com/joho/godotenv"
	"log"
	"strings"
)

func DotEnv(filePath ...string) {
	err := godotenv.Load(filePath...)
	if err != nil {
		log.Default().Println("error loading from environment file from dotenv :", strings.Join(filePath, ","), "error :", err)
	}
}
