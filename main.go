// Copyright (c) Luke Atherton 2015

package main

import (
	"log"
	"net/http"

	. "github.com/lukeatherton/hivebase-authenticator/app"
)

func main() {
	config := BuildConfig()
	publisher := NewAmpqPublisher(config.GetExchangeAddress(), config.GetAmpqUsername(), config.GetAmpqPassword(), config.GetTopic())
	repo := NewMongoRepo(config.GetDbHosts(), config.GetAuthDb(), config.GetDbUsername(), config.GetDbPassword())

	auth := BuildAuthenticator(repo, config.GetPrivateKeyPath(), config.GetPublicKeyPath())

	if err := http.ListenAndServe(":8001", NewRouter(publisher, repo, auth)); err != nil {
		log.Fatal(err)
	}
}
