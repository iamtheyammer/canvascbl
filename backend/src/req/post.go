package req

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

func MakePostRequestWithBody(url string, body []byte) (*http.Response, string, error) {
	client := http.Client{}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	return resp, string(respBody), nil
}

func MakePostRequest(url string) (*http.Response, string, error) {
	return MakePostRequestWithBody(url, nil)
}
