package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Request struct {
	Client        *http.Client
	Request       *http.Request
	Response      *http.Response
	Url           string
	Method        string
	Body          []byte
	ResponseBody  []byte
	Timeout       time.Duration
	Authorization string
	Retry         int
	Debug         bool
}

func (o *Request) Log(format string, v ...interface{}) {
	if o.Debug {
		log.Printf(format, v...)
	}
}

func (o *Request) Make() (err error) {
	o.Request, err = http.NewRequest(o.Method, o.Url, bytes.NewBuffer(o.Body))
	if err != nil {
		o.Log("%v", err)
		return err
	}
	o.Request.Header.Set("Content-Type", "application/json")
	if o.Authorization != "" {
		o.Request.Header.Set("Authorization", o.Authorization)
	}

	transport := &http.Transport{
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
	}
	o.Client = &http.Client{Transport: transport}
	o.Response, err = o.Client.Do(o.Request)
	if err != nil {
		o.Log("%v", err)
		return err
	}
	var prettyJSON bytes.Buffer
	if o.Method == "POST" || o.Method == "PATCH" || o.Method == "PUT" {
		err = json.Indent(&prettyJSON, o.Body, "", "\t")
		if err != nil {
			o.Log("JSON parse error: ", err)
			return err
		}
		o.Log("API Request: %s %d %s %s", o.Method, o.Response.StatusCode, o.Url, prettyJSON.String())
	}
	o.Log("API Request: %s %d %s", o.Method, o.Response.StatusCode, o.Url)
	o.ResponseBody, err = ioutil.ReadAll(o.Response.Body)
	_ = o.Response.Body.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	if len(o.ResponseBody) == 0 && o.Method == "DELETE" {
		return nil
	}

	firstSymbol := string(o.ResponseBody[0:1])
	if firstSymbol != "{" && firstSymbol != "[" {
		return fmt.Errorf("API response not in json format:\n%s", o.ResponseBody)
	}

	var prettyJSONResp bytes.Buffer
	err = json.Indent(&prettyJSONResp, o.ResponseBody, "", "\t")
	if err != nil {
		o.Log("JSON parse error: ", err)
		return err
	}
	o.Log("API Response: %s", string(prettyJSONResp.Bytes()))
	return nil
}
