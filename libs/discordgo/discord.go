// Discordgo - Discord bindings for Go
// Available at https://github.com/bwmarrin/discordgo

// Copyright 2015-2016 Bruce Marriner <bruce@sqls.net>.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains high level helper functions and easy entry points for the
// entire discordgo package.  These functions are being developed and are very
// experimental at this point.  They will most likely change so please use the
// low level functions if that's a problem.

// Package discordgo provides Discord binding for Go
package discordgo

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
)

// VERSION of DiscordGo, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.27.1"

// New creates a new Discord session with provided token.
// If the token is for a bot, it must be prefixed with "Bot "
//
//	e.g. "Bot ..."
//
// Or if it is an OAuth2 token, it must be prefixed with "Bearer "
//
//	e.g. "Bearer ..."
func New() (s *Session) {

	// Create an empty Session interface.
	s = &Session{
		State:                              NewState(),
		Ratelimiter:                        NewRatelimiter(),
		StateEnabled:                       true,
		Compress:                           true,
		ShouldReconnectOnError:             true,
		ShouldReconnectVoiceOnSessionError: true,
		ShouldRetryOnRateLimit:             true,
		ShardID:                            0,
		ShardCount:                         1,
		MaxRestRetries:                     3,
		Client:                             &http.Client{Timeout: (20 * time.Second)},
		Dialer:                             websocket.DefaultDialer,
		UserAgent:                          "DiscordBot (https://github.com/bwmarrin/discordgo, v" + VERSION + ")",
		sequence:                           new(int64),
		LastHeartbeatAck:                   time.Now().UTC(),
	}

	// Initialize the Identify Package with defaults
	// These can be modified prior to calling Open()
	s.Identify.Compress = true
	s.Identify.LargeThreshold = 250
	s.Identify.Properties.OS = runtime.GOOS
	s.Identify.Properties.Browser = "DiscordGo v" + VERSION
	s.Identify.Intents = IntentsAllWithoutPrivileged

	return
}

func (s *Session) SetToken(token string) {
	s.Identify.Token = token
	s.Token = token
}

func (s *Session) Login(email, password string) (ticket string, err error) {
	payload := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{email, password}

	response, err := s.RequestWithBucketID("POST", EndpointLogin, payload, EndpointLogin)
	if err != nil {
		return "", err
	}

	data := struct {
		Token  string `json:"token"`
		Ticket string `json:"ticket"`
		MFA    bool   `json:"mfa"`
	}{}

	err = unmarshal(response, &data)
	if err != nil {
		return
	}

	s.Token = data.Token
	s.MFA = data.MFA
	ticket = data.Ticket
	return
}

func (s *Session) TOTP(ticked string, code string) error {
	payload := struct {
		Ticket string `json:"ticket"`
		Code   string `json:"code"`
	}{ticked, code}

	response, err := s.RequestWithBucketID("POST", EndpointTOTP, payload, EndpointTOTP)
	if err != nil {
		return err
	}

	data := struct {
		Token string `json:"token"`
	}{}

	err = unmarshal(response, &data)
	if err != nil {
		return err
	}
	s.Token = data.Token
	return nil
}
