package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func loginToAdmin(password string) ([]*http.Cookie, error) {
	form := url.Values{}
	form.Add("isTest", "false")
	form.Add("goformId", "LOGIN")
	form.Add("password", base64.URLEncoding.EncodeToString([]byte(password)))

	res, err := request(http.MethodPost, setProcessURL, bytes.NewBufferString(form.Encode()), nil)
	if err != nil {
		return nil, err
	}

	body, err := decodeBody(res.Body)
	if err != nil {
		return nil, err
	}

	if body["result"] != "0" {
		return nil, errors.New("unable to login to the admin interface")
	}

	return res.Cookies(), nil
}

// getBalance makes a USSD request to get the data balance,
// checks every second for the status of the first request, and
// then makes a final request to get the balance.
func getBalance(cookies []*http.Cookie) (int, error) {
	if err := startBalanceRequest(cookies); err != nil {
		return 0, err
	}

	done, errs := make(chan bool), make(chan error)

	go func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case <-ticker.C:
				completed, err := checkBalanceRequest(cookies)
				if err != nil {
					errs <- err
					return
				}

				if !completed {
					break
				}

				done <- true
				return
			}
		}
	}()

	for {
		select {
		case err := <-errs:
			log.Fatal(err)
		case <-done:
			return finishBalanceRequest(cookies)
		}
	}
}

func startBalanceRequest(cookies []*http.Cookie) error {
	form := url.Values{}
	form.Add("isTest", "false")
	form.Add("goformId", "USSD_PROCESS")
	form.Add("notCallback", "true")
	form.Add("USSD_operator", "ussd_send")
	form.Add("USSD_send_number", "*461*4#")

	res, err := request(http.MethodPost, setProcessURL, bytes.NewBufferString(form.Encode()), cookies)
	if err != nil {
		return err
	}

	_, err = decodeBody(res.Body)
	return err
}

func checkBalanceRequest(cookies []*http.Cookie) (bool, error) {
	getURL, err := url.Parse(getProcessURL)
	if err != nil {
		return false, err
	}

	params := url.Values{}
	params.Add("cmd", "ussd_write_flag")
	params.Add("_", strconv.Itoa(int(time.Now().Unix())))
	getURL.RawQuery = params.Encode()

	res, err := request(http.MethodGet, getURL.String(), nil, cookies)
	if err != nil {
		return false, err
	}

	body, err := decodeBody(res.Body)
	if err != nil {
		return false, err
	}

	switch body["ussd_write_flag"] {
	case "15":
		return false, nil
	case "16":
		return true, nil
	default:
		return false, errors.New("invalid ussd_write_flag value")
	}
}

func finishBalanceRequest(cookies []*http.Cookie) (int, error) {
	getURL, err := url.Parse(getProcessURL)
	if err != nil {
		return 0, err
	}

	params := url.Values{}
	params.Add("cmd", "ussd_data_info")
	params.Add("_", strconv.Itoa(int(time.Now().Unix())))
	getURL.RawQuery = params.Encode()

	res, err := request(http.MethodGet, getURL.String(), nil, cookies)
	if err != nil {
		return 0, err
	}

	body, err := decodeBody(res.Body)
	if err != nil {
		return 0, err
	}

	return balanceFromResponse(body)
}

func balanceFromResponse(response map[string]string) (int, error) {
	msg, err := hex.DecodeString(response["ussd_data"])
	if err != nil {
		return 0, err
	}

	// Remove irrelevant zero bytes that distort
	// the string value of the message.
	cleanMsg := make([]byte, 0, len(msg))
	for _, b := range msg {
		if b != 0 {
			cleanMsg = append(cleanMsg, b)
		}
	}

	log.Printf("Received message from MTN: %s", string(cleanMsg))

	matches := balanceRegex.FindStringSubmatch(string(cleanMsg))
	if len(matches) == 0 {
		return 0, errors.New(fmt.Sprintf("received invalid message: %s", string(cleanMsg)))
	}

	return strconv.Atoi(matches[1])
}

func decodeBody(body io.Reader) (map[string]string, error) {
	var b map[string]string
	return b, json.NewDecoder(body).Decode(&b)
}

func request(method string, url string, body io.Reader, cookies []*http.Cookie) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate")
	req.Header.Add("Accept-Language", "en-US,en;q=0.9")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Host", "192.168.0.1")
	req.Header.Add("Origin", "http://192.168.0.1")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Referer", "http://192.168.0.1/index.html")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	return http.DefaultClient.Do(req)
}
