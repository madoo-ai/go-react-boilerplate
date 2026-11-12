package xmodule

import "testing"

func TestConfigValidate(t *testing.T) {
	valid := Config{
		DSN: "postgres://localhost/db", AppURL: "http://localhost:5173",
		AppLocale: "id", SessionTTL: 1,
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("valid config rejected: %v", err)
	}

	cases := []struct {
		name string
		cfg  Config
	}{
		{"missing DSN", Config{AppURL: "http://localhost", AppLocale: "id", SessionTTL: 1}},
		{"relative URL", Config{DSN: "x", AppURL: "/app", AppLocale: "id", SessionTTL: 1}},
		{"bad locale", Config{DSN: "x", AppURL: "http://localhost", AppLocale: "fr", SessionTTL: 1}},
		{"bad TTL", Config{DSN: "x", AppURL: "http://localhost", AppLocale: "en"}},
	}
	for _, test := range cases {
		if err := test.cfg.Validate(); err == nil {
			t.Errorf("%s: expected an error", test.name)
		}
	}
}

func TestResolveCookieSecure(t *testing.T) {
	t.Setenv("COOKIE_SECURE", "")
	secure := Config{AppURL: "https://example.com"}
	secure.ResolveCookieSecure()
	if !secure.CookieSecure {
		t.Fatal("https APP_URL should enable secure cookies")
	}
	plain := Config{AppURL: "http://example.com"}
	plain.ResolveCookieSecure()
	if plain.CookieSecure {
		t.Fatal("http APP_URL should not enable secure cookies")
	}

	t.Setenv("COOKIE_SECURE", "true")
	forced := Config{AppURL: "http://example.com"}
	forced.ResolveCookieSecure()
	if !forced.CookieSecure {
		t.Fatal("COOKIE_SECURE=true should win")
	}
}

func TestEnvHelpers(t *testing.T) {
	t.Setenv("BOOL_TEST", "true")
	if !envBool("BOOL_TEST") {
		t.Fatal("expected true")
	}
	t.Setenv("BOOL_TEST", "yes")
	if envBool("BOOL_TEST") {
		t.Fatal("only explicit true values should enable a flag")
	}
	if got := splitList(" http://a , http://b/ "); len(got) != 2 || got[1] != "http://b" {
		t.Fatalf("splitList = %#v", got)
	}
}
