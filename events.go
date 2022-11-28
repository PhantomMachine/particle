package particle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/r3labs/sse"
)

type ParticleEvent struct {
	Data        string    `json:"data"`
	TTL         int       `json:"ttl"`
	PublishedAt time.Time `json:"published_at"`
	CoreID      string    `json:"coreid"`
}

type Event struct {
	*sse.Event
}

type ParticleEventHandler func(event, data string)

func sseHandler(h ParticleEventHandler) func(*sse.Event) {
	return func(msg *sse.Event) {
		if len(msg.Data) == 0 {
			return
		}

		evt := &ParticleEvent{}
		err := json.NewDecoder(bytes.NewReader(msg.Data)).Decode(evt)
		if err != nil {
			panic(err)
		}

		h(string(msg.Event), evt.Data)
	}
}

func Subscribe(topic, token string, h ParticleEventHandler) {
	segments := strings.Split(topic, "/")
	url := fmt.Sprintf("https://api.particle.io/v1/events/%s?access_token=%s", segments[0], token)
	client := sse.NewClient(url)

	handler := h

	h = func(event, data string) {
		if event != topic {
			return
		}

		handler(event, data)
	}

	err := client.Subscribe(filepath.Join(segments[1:]...), sseHandler(h))
	if err != nil {
		panic(err)
	}
}
