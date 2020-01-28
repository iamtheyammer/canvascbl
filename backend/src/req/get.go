package req

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var client = http.Client{}

func MakeAuthenticatedGetRequest(url string, token string) (*http.Response, string, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return resp, string(body), nil
}
