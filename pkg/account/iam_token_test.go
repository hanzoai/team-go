// Copyright 2026 Hanzo AI, Inc. All rights reserved.
// SPDX-License-Identifier: MIT

package account

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/hanzoai/base/core"
	"github.com/hanzoai/base/tools/router"
	"github.com/hanzoai/team-go/pkg/token"
	"github.com/hanzoai/team-go/pkg/wsauth"
)

// rsToken is a token that SAYS RS256 in its header. Nothing here signs it,
// because nothing here checks it: the point of the routing under test is that
// this service never verifies an RS token itself — it reads the answer the
// platform's JWKS middleware already reached.
func rsToken(sub string) string {
	part := func(v any) string {
		b, _ := json.Marshal(v)
		return base64.RawURLEncoding.EncodeToString(b)
	}
	h := part(map[string]string{"typ": "JWT", "alg": "RS256"})
	p := part(map[string]string{"sub": sub, "owner": "acme"})
	return h + "." + p + ".notasignature"
}

func req(tok string) *core.RequestEvent {
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+tok)
	return &core.RequestEvent{Event: router.Event{Response: httptest.NewRecorder(), Request: r}}
}

// The whole point: a caller carrying an IAM token is the same person to this
// service as one carrying its own session, and the identity comes from the
// verification the platform already did.
func TestAnIAMTokenIdentifiesTheCaller(t *testing.T) {
	re := req(rsToken("z@hanzo.ai"))
	// What Base's JWKS middleware leaves behind once it has verified the token.
	re.Set("authSub", "z@hanzo.ai")

	g := &api{cfg: config{serverSecret: "s3cret"}}
	acct, _, _, err := g.account(re)
	if err != nil {
		t.Fatalf("a verified IAM caller was refused: %v", err)
	}
	if want := wsauth.AccountID("z@hanzo.ai"); acct != want {
		t.Errorf("account = %q, want %q — the id must be the one member rows carry", acct, want)
	}
}

// An RS token nobody verified is refused, and NOT retried on the HS256 path.
// The two are different authorities; falling between them is how an unchecked
// token gets treated as a session.
func TestAnUnverifiedIAMTokenIsRefused(t *testing.T) {
	re := req(rsToken("z@hanzo.ai")) // no authSub: the middleware never vouched for it

	g := &api{cfg: config{serverSecret: "s3cret"}}
	if _, _, _, err := g.account(re); err == nil {
		t.Fatal("an unverified RS token was accepted")
	} else if !strings.Contains(err.Error(), "not verified") {
		t.Errorf("refused for the wrong reason: %v — it must not read as a bad session", err)
	}
}

// The session this service mints for its own SPA keeps working. The IAM path is
// an addition, not a replacement, until the SPA stops carrying one.
func TestTheOwnSessionStillWorks(t *testing.T) {
	const acct = "6f1b7c62-6c1a-4a54-9a5e-3f5b1c2d4e70"
	tok, err := token.Generate(acct, "", map[string]any{"org": "acme"}, "s3cret")
	if err != nil {
		t.Fatal(err)
	}
	g := &api{cfg: config{serverSecret: "s3cret"}}
	got, org, _, err := g.account(req(tok))
	if err != nil {
		t.Fatalf("the service's own session was refused: %v", err)
	}
	if got != acct || org != "acme" {
		t.Errorf("account=%q org=%q, want %s/acme", got, org, acct)
	}
}

// A token signed HS256 by somebody else is still refused. Routing by `alg` must
// not become a way past the secret.
func TestAForgedSessionIsStillRefused(t *testing.T) {
	const acct = "6f1b7c62-6c1a-4a54-9a5e-3f5b1c2d4e70"
	tok, err := token.Generate(acct, "", nil, "the-wrong-secret")
	if err != nil {
		t.Fatal(err)
	}
	g := &api{cfg: config{serverSecret: "s3cret"}}
	if _, _, _, err := g.account(req(tok)); err == nil {
		t.Fatal("a token signed with another secret was accepted")
	}
}
