package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/phantommachine/particle"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	cnf := &Config{}

	cnf.ClientID = os.Getenv("PARTICLE_CLIENT_ID")
	cnf.ClientSecret = os.Getenv("PARTICLE_CLIENT_SECRET")
	cnf.Username = os.Getenv("PARTICLE_USERNAME")
	cnf.Password = os.Getenv("PARTICLE_PASSWORD")
	cnf.Topic = os.Getenv("PARTICLE_EVENT_TOPIC")

	par, err := particle.Authorize(cnf.ClientID, cnf.ClientSecret, cnf.Username, cnf.Password)
	if err != nil {
		log.Fatal(err)
	}

	accessToken := par.AccessToken

	devices, err := particle.ListDevices(accessToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%#v\n", devices)

	for _, device := range devices {
		for variable := range device.Variables {
			err := particle.GetVariable(accessToken, device.ID, variable)
			if err != nil {
				panic(err)
			}
		}
	}

	h := func(e, d string) {
		measurements := struct {
			Temperature float32 `json:"temperature"`
			Humidity    float32 `json:"humidity"`
		}{}
		r := strings.NewReader(d)
		err := json.NewDecoder(r).Decode(&measurements)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %s\n\t%#v\n", e, d, measurements)
	}

	particle.Subscribe(cnf.Topic, accessToken, h)

	// loop indefinitely, and handle SSE events
	for {
		runtime.Gosched()
	}
}

type Config struct {
	ClientID     string `env:"PARTICLE_CLIENT_ID"`
	ClientSecret string `env:"PARTICLE_CLIENT_SECRET"`
	Username     string `env:"PARTICLE_USERNAME"`
	Password     string `env:"PARTICLE_PASSWORD"`
	Topic        string `env:"PARTICLE_EVENT_TOPIC"`
}
