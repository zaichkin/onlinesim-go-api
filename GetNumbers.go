package onlinesim

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type GetNumbers struct {
	client    *Onlinesim
}


func (c *Onlinesim) GetNumbers() *GetNumbers {
	return &GetNumbers{
		client:    c,
	}
}

type PriceResponse struct {
	Response string `json:"response"`
	Price    int    `json:"price"`
}

func (c *GetNumbers) price(service string) (error, int) {
	m := make(map[string]string)
	m["service"] = service
	result := c.client.get("getPrice", m)

	response := PriceResponse{}
	err := json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), 0
	}

	err = c.client.checkResponse(response.Response)
	if err != nil {
		return fmt.Errorf("%w", err), 0
	}
	return nil, response.Price
}

type GetResponse struct {
	Response string `json:"response"`
	Tzid     int `json:"tzid"`
}

func (c *GetNumbers) get(service string, country int) (error, int) {
	m := make(map[string]string)
	m["service"] = service
	m["country"] = strconv.Itoa(country)
	result := c.client.get("getNum", m)

	response := GetResponse{}
	err := json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), 0
	}

	err = c.client.checkResponse(response.Response)
	if err != nil {
		return fmt.Errorf("%w", err), 0
	}
	return nil, response.Tzid
}

type Order string
const (
	ASC Order = "ASC"
	DESC Order = "DESC"
)


type StateResponse []State

type State struct {
	Tzid   int    `json:"tzid"`
	Form   string `json:"form"`
	Time   int    `json:"time"`
	Number string `json:"number"`
	Msg    []struct {
		Service string `json:"service"`
		Msg     string `json:"msg"`
	} `json:"msg,omitempty"`
	Service  string `json:"service"`
	Country  int    `json:"country"`
	Response string `json:"response"`
	Sum      int    `json:"sum,omitempty"`
}

func (c *GetNumbers) state(message_to_code int, orderby Order) (error, StateResponse) {
	m := make(map[string]string)
	m["message_to_code"] = strconv.Itoa(message_to_code)
	m["orderby"] = string(orderby)
	m["msg_list"] = "1"
	m["clean"] = "0"
	m["type"] = "index"
	result := c.client.get("getState", m)

	err :=c.client.checkEmptyResponse(result)
	if err != nil {
		return fmt.Errorf("%w", err), nil
	}

	response := StateResponse{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), nil
	}

	return nil, response
}

func (c *GetNumbers) stateOne(tzid int, message_to_code int) (error, State) {
	m := make(map[string]string)
	m["message_to_code"] = strconv.Itoa(message_to_code)
	m["tzid"] = strconv.Itoa(tzid)
	m["msg_list"] = "1"
	m["clean"] = "0"
	m["type"] = "index"
	result := c.client.get("getState", m)

	err :=c.client.checkEmptyResponse(result)
	if err != nil {
		return fmt.Errorf("%w", err), State{}
	}

	response := StateResponse{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), State{}
	}

	return nil, response[0]
}

func (c *GetNumbers) next(tzid int) (error, bool) {
	m := make(map[string]string)
	m["tzid"] = strconv.Itoa(tzid)
	result := c.client.get("setOperationRevise", m)

	response := Default{}
	err := json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), false
	}

	err = c.client.checkResponse(response.Response)
	if err != nil {
		return fmt.Errorf("%w", err), false
	}
	return nil, true
}

func (c *GetNumbers) close(tzid int) (error, bool) {
	m := make(map[string]string)
	m["tzid"] = strconv.Itoa(tzid)
	result := c.client.get("setOperationOk", m)

	response := Default{}
	err := json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), false
	}

	err = c.client.checkResponse(response.Response)
	if err != nil {
		return fmt.Errorf("%w", err), false
	}
	return nil, true
}

type TariffsResponse struct {
	Name     string `json:"name"`
	Position int    `json:"position"`
	Code     int    `json:"code"`
	Other    interface{}    `json:"other"`
	New      bool   `json:"new"`
	Enabled  bool   `json:"enabled"`
	Services map[string]Service `json:"services"`
}

type Service struct {
	Count   interface{}    `json:"count"`
	Popular bool   `json:"popular"`
	Code    int    `json:"code"`
	Price   int    `json:"price"`
	ID      int    `json:"id"`
	Service string `json:"service"`
	Slug    interface{} `json:"slug"`
}

func (c *GetNumbers) tariffs() (error, map[string]TariffsResponse) {
	m := make(map[string]string)
	m["country"] = "all"
	result := c.client.get("getNumbersStats", m)
	err :=c.client.checkEmptyResponse(result)
	if err != nil {
		return fmt.Errorf("%w", err), map[string]TariffsResponse{}
	}

	response := map[string]TariffsResponse{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), map[string]TariffsResponse{}
	}

	return nil, response
}

func (c *GetNumbers) tariffsOne(country int) (error, TariffsResponse) {
	m := make(map[string]string)
	m["country"] = strconv.Itoa(country)
	result := c.client.get("getNumbersStats", m)

	err :=c.client.checkEmptyResponse(result)
	if err != nil {
		return fmt.Errorf("%w", err), TariffsResponse{}
	}

	response := TariffsResponse{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), TariffsResponse{}
	}

	return nil, response
}

type ServiceResponse struct {
	Service  []string `json:"service"`
	Response string   `json:"response"`
}

func (c *GetNumbers) service() (error, []string) {
	m := make(map[string]string)
	result := c.client.get("getService", m)

	response := ServiceResponse{}
	err := json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), nil
	}

	err = c.client.checkResponse(response.Response)
	if err != nil {
		return fmt.Errorf("%w", err), nil
	}

	return nil, response.Service
}

type ServiceNumberResponse struct {
	Number  []string `json:"number"`
	Response string   `json:"response"`
}

func (c *GetNumbers) serviceNumber(service string) (error, []string) {
	m := make(map[string]string)
	m["service"] = service
	result := c.client.get("getServiceNumber", m)

	response := ServiceNumberResponse{}
	err := json.Unmarshal(result, &response)
	if err != nil {
		return fmt.Errorf("%w", err), nil
	}

	err = c.client.checkResponse(response.Response)
	if err != nil {
		return fmt.Errorf("%w", err), nil
	}

	return nil, response.Number
}