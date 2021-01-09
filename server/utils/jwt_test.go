package jwtea

import (
	"fmt"
	"testing"
	"time"
)


func TestValidate(tg *testing.T) {
	provider := NewProvider(&Configuration{
		Secret:   "testing_secret",
		Issuer:   "Foo",
		Audience: "Foo",
	})
	now := time.Now().UTC()
	anHourAgo := now.Add(time.Hour * -1).Unix()
	anHourFromNow := now.Add(time.Hour * 1).Unix()
	
	tg.Parallel()
	
	tg.Run("happy path", func(t *testing.T) {
		body := &Body{
			User: "test.m.tester",
			Role: "everyone",
			Jti:  "token-1",
			Iss:  "Foo",
			Aud:  "Foo",
			Exp:  anHourFromNow,
			Nbf:  anHourAgo,
			Iat:  anHourAgo,
		}
		token := provider.generate(body)
		err := provider.Validate(fmt.Sprintf("Bearer %s", token))
		if err != nil {
			t.Fatalf("expected Validate to pass, but received: %s", err)
		}
	})

	var issAudTestData = []struct {
		scenario string
		issuer string
		audience string
	} {
		{ "invalid iss", "Bar", "Foo" },
		{ "invalid aud", "Foo", "Bar" },
		{ "invalid iss and aud", "Bar", "Bar" },
	}

	for _, test := range issAudTestData {
		tg.Run(test.scenario, func(t *testing.T) {
			body := &Body{
				User: "test.m.tester",
				Role: "everyone",
				Jti:  "token-1",
				Iss:  test.issuer,
				Aud:  test.audience,
				Exp:  anHourFromNow,
				Nbf:  anHourAgo,
				Iat:  anHourAgo,
			}
			token := provider.generate(body)
			err := provider.Validate(fmt.Sprintf("Bearer %s", token))
			if err == nil {
				t.Fatalf("expected Validate to fail, but passed with iss, aud: %s,%s", test.issuer, test.audience)
			}
		})
	}
	
	var expTestData = []struct {
		scenario string
		expiration int64
		notBefore int64
	} {
		{ "invalid exp", anHourAgo, anHourAgo },
		{ "invalid nbf", anHourFromNow, anHourFromNow },
		{ "invalid exp and nbf", anHourAgo, anHourFromNow },
	}

	for _, test := range expTestData {
		tg.Run(test.scenario, func(t *testing.T) {
			body := &Body{
				User: "test.m.tester",
				Role: "everyone",
				Jti:  "token-1",
				Iss:  "Foo",
				Aud:  "Foo",
				Exp:  test.expiration,
				Nbf:  test.notBefore,
				Iat:  anHourAgo,
			}
			token := provider.generate(body)
			err := provider.Validate(fmt.Sprintf("Bearer %s", token))
			if err == nil {
				t.Fatalf("expected Validate to fail, but passed with exp, nbf: %d,%d", test.expiration, test.notBefore)
			}
		})
	}

	tg.Run("invalid signature", func(t *testing.T) {
		anotherProvider := NewProvider(&Configuration{
			Secret:   "another_secret",
			Issuer:   "Foo",
			Audience: "Foo",
		})
		body := &Body{
			User: "test.m.tester",
			Role: "everyone",
			Jti:  "token-1",
			Iss:  "Foo",
			Aud:  "Foo",
			Exp:  anHourFromNow,
			Nbf:  anHourAgo,
			Iat:  anHourAgo,
		}
		token := anotherProvider.generate(body)
		err := provider.Validate(fmt.Sprintf("Bearer %s", token))
		if err == nil {
			t.Fatal("expected Validate to fail for bad signature")
		}
	})
}