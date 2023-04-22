package adguardhome

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Record struct {
	Domain string `json:"domain"`
	Answer string `json:"answer"`
}

type Records []Record

func (c *DnsCfg) RewriteList() (*Records, error) {
	url := c.url + "/control/rewrite/list"

	mimeTypeJson := "application/json"

	tr := &http.Transport{
		IdleConnTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Add("Content-Type", mimeTypeJson)
	req.Header.Add("Accept", mimeTypeJson)
	req.SetBasicAuth(c.username, c.password)

	response, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("Please contact an administrator if the problem persists (method: GetDnsConfig, status code: %d)", response.StatusCode)
	}

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	data := new(Records)

	if err := json.Unmarshal(respBody, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (c *DnsCfg) RewriteAdd(payload *Record) error {
	url := c.url + "/control/rewrite/add"

	mimeTypeJson := "application/json"

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		IdleConnTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Add("Content-Type", mimeTypeJson)
	req.Header.Add("Accept", mimeTypeJson)
	req.SetBasicAuth(c.username, c.password)

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Please contact an administrator if the problem persists (method: SaveDnsConfig, status code: %d)", response.StatusCode)
	}

	return nil
}

func (c *DnsCfg) RewriteDelete(payload *Record) error {
	url := c.url + "/control/rewrite/delete"

	mimeTypeJson := "application/json"

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	tr := &http.Transport{
		IdleConnTimeout: 5 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout: time.Second * 10,
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Add("Content-Type", mimeTypeJson)
	req.Header.Add("Accept", mimeTypeJson)
	req.SetBasicAuth(c.username, c.password)

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Please contact an administrator if the problem persists (method: SaveDnsConfig, status code: %d)", response.StatusCode)
	}

	return nil
}
