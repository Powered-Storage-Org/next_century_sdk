package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Powered-Storage-Org/next_century_sdk/core/schema"
)

// next century meters API endpoint
const (
	Login = "/Login"

	// GET /api/Properties/:id/DailyReads/?date&from&to // &from=%s&to=%s is optional
	DailyReads = "/api/Properties/%s/DailyReads/"

	// GET /api/Properties/:id/Units"
	Units = "/api/Properties/%s/Units"

	// TODO: add more API endpoints
)

type client struct {
	apiURL string

	apiKey string
}

func New(addr, email, password string) *client {
	loginUrl, err := url.JoinPath(addr, Login)
	if err != nil {
		panic(err)
	}

	apiKey := genNewToken(email, password, loginUrl)

	return &client{
		apiURL: addr,
		apiKey: apiKey,
	}
}

func genNewToken(email, password, loginUrl string) string {
	postData := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: nil, // Use default configuration which validates certificates
		},
	}

	req, err := http.NewRequest("POST", loginUrl, strings.NewReader(postData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer closRespBody(resp)

	body := struct {
		Token string `json:"token"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(&body); err != nil {
		log.Fatal(err)
	}

	return body.Token
}

func (c *client) GetDailyReadsWithCustomJsonPars(propertyID string, timeReq schema.TimeRequest, structParser any) error {
	DailyReadsUrl, err := url.JoinPath(c.apiURL, fmt.Sprintf(DailyReads, propertyID))
	if err != nil {
		return err
	}

	if timeReq.Date.IsZero() {
		return fmt.Errorf("date is required")
	}

	DailyReadsUrl += fmt.Sprintf("?date=%s", timeReq.Date.Format("2006-01-02"))

	if !timeReq.From.IsZero() {
		DailyReadsUrl += fmt.Sprintf("&from=%s", timeReq.From.Format("2006-01-02"))
	}

	if !timeReq.To.IsZero() {
		DailyReadsUrl += fmt.Sprintf("&to=%s", timeReq.To.Format("2006-01-02"))
	}

	req, err := http.NewRequest("GET", DailyReadsUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("authorization", c.apiKey)
	req.Header.Set("version", "2")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer closRespBody(resp)

	// if status is StatusUnauthorized gen new token
	if resp.StatusCode == http.StatusUnauthorized {
		log.Fatal("Unauthorized, pleas rerun the program")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get daily reads: %s", resp.Status)
	}

	if err = json.NewDecoder(resp.Body).Decode(structParser); err != nil {
		return err
	}

	return nil
}

func (c *client) GetDailyReads(propertyID string, timeReq schema.TimeRequest) ([]*schema.MeterData, error) {
	var meterData []*schema.MeterData
	if err := c.GetDailyReadsWithCustomJsonPars(propertyID, timeReq, &meterData); err != nil {
		return nil, err
	}

	return meterData, nil
}

func (c *client) GetUnits(propertyID string) ([]*schema.Unit, error) {
	UnitsUrl, err := url.JoinPath(c.apiURL, fmt.Sprintf(Units, propertyID))
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", UnitsUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer closRespBody(resp)

	// if status is StatusUnauthorized gen new token
	if resp.StatusCode == http.StatusUnauthorized {
		log.Fatal("Unauthorized, pleas rerun the program")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get units: %s", resp.Status)
	}

	var units []*schema.Unit
	if err = json.NewDecoder(resp.Body).Decode(&units); err != nil {
		return nil, err
	}

	return units, nil
}
