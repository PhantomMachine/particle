package particle

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type ParticleAuthError struct {
	Error     string `json:"error"`
	ErrorDesc string `json:"error_description"`
	MFAToken  string `json:"mfa_token"`
}

type ParticleAuthResponse struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func Authorize(id, secret, username, password string) (*ParticleAuthResponse, error) {

	v := url.Values{}
	v.Add("grant_type", "password")
	v.Add("username", username)
	v.Add("password", password)
	v.Add("client_id", id)
	v.Add("client_secret", secret)

	client := http.Client{}
	do := func(v url.Values) (*http.Response, error) {
		req, _ := http.NewRequest("POST", "https://api.particle.io/oauth/token", strings.NewReader(v.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		basicAuth := base64.RawURLEncoding.EncodeToString([]byte(username + ":" + password))
		req.Header.Set("Authentication", "Bearer "+basicAuth)
		return client.Do(req)
	}

	res, err := do(v)
	if err != nil {
		log.Fatal(err)
	}

	switch res.StatusCode {
	case http.StatusOK:
		break
	case http.StatusForbidden:
		var perr ParticleAuthError
		err := json.NewDecoder(res.Body).Decode(&perr)
		if err != nil {
			return nil, err
		}

		if perr.Error == "mfa_required" {
			code, _ := ReadFromStdin("Enter otp: ")

			v.Add("mfa_token", perr.MFAToken)
			v.Add("otp", code)
			v.Set("grant_type", "urn:custom:mfa-otp")

			res, err = do(v)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("particle auth error - %s: %s", perr.Error, perr.ErrorDesc)
		}
	default:
		return nil, fmt.Errorf("unhandled response code: %s", res.Status)
	}

	par := &ParticleAuthResponse{}
	err = json.NewDecoder(res.Body).Decode(par)
	if err != nil {
		return nil, err
	}

	return par, nil
}

func ReadFromStdin(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	code, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.Trim(code, "\n"), nil
}
