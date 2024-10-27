package telegram

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"tgBot/lib/errorsLib"
)

type Client struct {
	host     string
	basePath string
	client   *http.Client
}

func NewClient(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   &http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset, limit int) (updates []Update, err error) {
	defer func() {
		err = errorsLib.Wrap("error request:", err)
	}()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.DoRequest(getUpdates, q)
	if err != nil {
		return nil, err
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatId int, text string) error {
	q := url.Values{}
	q.Add("chatId", strconv.Itoa(chatId))
	q.Add("text", text)

	_, err := c.DoRequest(sendMessage, q)
	if err != nil {
		return errorsLib.Wrap("error send message: ", err)
	}

	return nil
}

func (c *Client) DoRequest(method string, query url.Values) (data []byte, err error) {
	defer func() {
		err = errorsLib.Wrap("error request:", err)
	}()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	return body, nil
}
