package particle

import (
	"encoding/json"
	"net/http"
)

type DeviceInfo struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Connected    bool              `json:"connected"`
	Online       bool              `json:"online"`
	Variables    map[string]string `json:"variables"`
	SerialNumber string            `json:"serial_number"`
}

func ListDevices(accessToken string) ([]DeviceInfo, error) {
	client := http.Client{}
	req, _ := http.NewRequest("GET", "https://api.particle.io/v1/devices", nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	var devices []DeviceInfo
	err = json.NewDecoder(res.Body).Decode(&devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}
