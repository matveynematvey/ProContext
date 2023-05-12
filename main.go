package main

import (
	"ProContext/lib"
	"ProContext/valute"
	"bytes"
	"encoding/xml"
	"fmt"
	"golang.org/x/net/html/charset"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	cbHostAPI = "https://www.cbr.ru/scripts/XML_daily_eng.asp"
	daysLimit = 90
)

type Client struct {
	Host   string
	Client http.Client
}

func NewClient() Client {
	return Client{
		Host:   cbHostAPI,
		Client: http.Client{Timeout: 10 * time.Second},
	}
}

func main() {
	client := NewClient()

	currencies := valute.NewValutes()

	start, end := time.Now(), time.Now().AddDate(0, 0, -daysLimit+1)
	fmt.Println("Program started...")
	for d := start; d.Before(end) == false; d = d.AddDate(0, 0, -1) {

		date := lib.FormDate(d)
		err := client.requestCurrencies(date, currencies)
		if err != nil {
			panic(err)
		}
	}

	lib.PrintResult(start, end, daysLimit, currencies)

}

func (c *Client) requestCurrencies(date string, valutes valute.Valutes) error {

	req, err := c.formReq(date)
	if err != nil {
		return err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return fmt.Errorf("do request error %w", err)
	}

	err = XMLHandler(date, body, valutes)
	if err != nil {
		return fmt.Errorf("cant unpack xml result: %w", err)
	}

	return nil
}

func XMLHandler(date string, data []byte, valutes valute.Valutes) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = charset.NewReaderLabel

	for t, err := decoder.Token(); t != nil; t, _ = decoder.Token() {

		if err != nil {
			return fmt.Errorf("token error %w", err)
		}

		switch tok := t.(type) {
		case xml.StartElement:
			if tok.Name.Local == "Valute" {

				valuteXML := new(valute.Valute)
				decoder.DecodeElement(&valuteXML, &tok)

				if _, ok := valutes[valuteXML.Name]; !ok {
					valutes[valuteXML.Name] = valute.New(date, valuteXML.ParseFloat())
				} else {
					valutes[valuteXML.Name].Add(date, valuteXML.ParseFloat())
				}
			}
		}
	}
	return nil
}

func (c *Client) formReq(date string) (*http.Request, error) {
	q := url.Values{}
	q.Add("date_req", date) //?

	req, err := http.NewRequest(http.MethodGet, c.Host+"?"+q.Encode(), nil)
	req.Header.Add("User-Agent", ``)
	if err != nil {
		return nil, fmt.Errorf("can't form request %w", err)
	}

	return req, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request error %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return body, nil
}
