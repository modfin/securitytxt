package main

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v7"
	"github.com/modfin/henry/slicez"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (t *Date) UnmarshalText(text []byte) error {
	tt, err := time.Parse("2006-01-02", string(text))
	if err != nil {
		tt, err = time.Parse(time.RFC3339, string(text))
	}
	*t = Date{tt}
	return err
}

type config struct {
	Port int `env:"PORT" envDefault:"8080"`

	Comment string `env:"COMMENT"`

	Contact []string `env:"CONTACT_URIS" envSeparator:" "`
	Expires Date     `env:"EXPIRES"`

	Acknowledgments    []string `env:"ACKNOWLEDGMENT_URIS" envSeparator:" "`
	Canonical          []string `env:"CANONICAL_URIS" envSeparator:" "`
	Encryption         []string `env:"ENCRYPTION_URIS" envSeparator:" "`
	Hiring             []string `env:"HIRING_URIS" envSeparator:" "`
	PreferredLanguages []string `env:"PREFERRED_LANGUAGES" envSeparator:" "`
	Policy             []string `env:"POLICY_URIS" envSeparator:" "`
	Extensions         []string `env:"EXTENSIONS" envSeparator:" "`

	RAW string `env:"RAW_SECURITY_TXT"`
}

func main() {

	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	data, err := mkSecurityText(cfg)
	if err != nil {
		fmt.Println("could not load required values")
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("== serving security.txt as follows ==\n%s\n", string(data))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		_, _ = w.Write(data)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}

func mkSecurityText(c config) ([]byte, error) {
	if c.RAW != "" {
		return []byte(c.RAW), nil
	}

	if len(c.Contact) == 0 {
		return nil, errors.New("contact must be provided")
	}

	var lines []string

	prefix := func(prefix string, slice []string) []string {
		return slicez.Map(slice, func(a string) string { return fmt.Sprintf("%s %s", prefix, a) })
	}

	if len(c.Comment) > 0 {
		comments := prefix("#", strings.Split(c.Comment, "\n"))
		lines = append(lines, comments...)
		lines = append(lines, "")
	}

	lines = append(lines, fmt.Sprintf("Expires: %s", c.Expires.Format(time.RFC3339)))
	lines = append(lines, prefix("Contact:", c.Contact)...)
	lines = append(lines, "")

	lines = append(lines, prefix("Acknowledgments:", c.Acknowledgments)...)

	lines = append(lines, prefix("Canonical:", c.Canonical)...)

	lines = append(lines, prefix("Encryption:", c.Encryption)...)

	lines = append(lines, prefix("Hiring:", c.Hiring)...)

	if len(c.PreferredLanguages) > 0 {
		lines = append(lines, fmt.Sprintf("Preferred-Languages: %s", strings.Join(c.PreferredLanguages, ", ")))
	}

	lines = append(lines, prefix("Policy:", c.Policy)...)

	if len(c.Extensions) > 0 {
		lines = append(lines, "")
		lines = append(lines, c.Extensions...)
	}

	return []byte(strings.Join(lines, "\n")), nil
}
