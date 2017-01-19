package main // import "github.com/nutmegdevelopment/nutcracker-cli"

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/nutmegdevelopment/nutcracker/secrets"
)

const (
	version  = "0.0.4"
	base64ID = "$base64$"
)

var (
	server string
	key    string
	id     string
	debug  bool
)

func getSecret(c *cli.Context) error {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	u := url.URL{
		Host:   server,
		Scheme: "https",
		Path:   "/secrets/view",
	}

	reqBody, err := json.Marshal(&map[string]string{"name": c.String("name")})
	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Request URI: ", u.String())
	req, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Secret-ID", id)
	req.Header.Add("X-Secret-Key", key)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}

	secret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Decode secret if base64 encoded
	if len(secret) > 8 && bytes.Compare(secret[0:8], []byte(base64ID)) == 0 {
		log.Debug("Decoding base64 secret")
		decoded := make([]byte, base64.StdEncoding.DecodedLen(len(secret)-8))
		n, err := base64.StdEncoding.Decode(decoded, secret[8:])
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("%s\n", decoded[:n])
	} else {
		fmt.Printf("%s\n", secret)
	}

	return err
}

func listSecrets(c *cli.Context) error {
	if debug {
		log.SetLevel(log.DebugLevel)
	}

	u := url.URL{
		Host:   server,
		Scheme: "https",
	}

	if c.Bool("all") {
		u.Path = "/secrets/list/secrets"
	} else {
		u.Path = fmt.Sprintf("/secrets/list/key/%s", id)
	}

	log.Debug("Request URI: ", u.String())
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("X-Secret-ID", id)
	req.Header.Add("X-Secret-Key", key)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(resp.Status)
	}

	dec := json.NewDecoder(resp.Body)

	for dec.More() {
		var recv []secrets.Secret
		err = dec.Decode(&recv)
		if err != nil {
			log.Fatal(err)
		}
		for i := range recv {
			log.Debug(recv[i])
			fmt.Println(recv[i].Name)
		}
	}

	err = resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func main() {
	app := cli.NewApp()
	app.Name = "nutcracker-cli"
	app.Usage = "CLI interface for nutcracker"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "server, s",
			Usage:       "Nutcracker server.  e.g localhost:443",
			Destination: &server,
		},
		cli.StringFlag{
			Name:        "id, i",
			Usage:       "Nutcracker API ID",
			Destination: &id,
		},
		cli.StringFlag{
			Name:        "key, k",
			Usage:       "Nutcracker API key",
			Destination: &key,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Debug output",
			Destination: &debug,
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "get",
			Aliases: []string{"g"},
			Usage:   "get a secret",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "name of the secret",
				},
			},
			Action: getSecret,
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list available secrets",
			Action:  listSecrets,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "all, a",
					Usage: "list all secrets",
				},
			},
		},
	}

	app.Run(os.Args)

}
