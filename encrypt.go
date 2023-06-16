package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

func GenerateBCryptPasswordUtil(hashedPwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(hashedPwd, bcrypt.DefaultCost)
	if err != nil {
		log.Printf("encrypting pwd %s", hashedPwd)
	}
	return string(hash)
}

func CompareHashAndPasswordUtil(hashedPWD []byte, PWD []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPWD, PWD)
	return err == nil
}
