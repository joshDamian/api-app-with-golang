package news

import (
    "net/http"
    "encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
    "time"
)

var IndexDomains string = "cnn.com,channelstv.com,punchng.com,techchrunch.com,thenextweb.com,bloomberg.com"

type Client struct {
	http     *http.Client
	key      string
	PageSize int
}

type Article struct {
	Source struct {
		ID   interface{} `json:"id"`
		Name string      `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}


type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}


func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{httpClient, key, pageSize}
}


func (c *Client) FetchAllNews(page string, domains string) (*Results, error) {
    if domains == "" {
        domains = IndexDomains;
    }
    endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?domains=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", domains, c.PageSize, page, c.key)
    resp, err := c.http.Get(endpoint)
   
    
    if err != nil {
        return nil, err
    }
    
    defer resp.Body.Close();
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf(string(body))
    }
    
    res := &Results{}
    return res, json.Unmarshal(body, res)
}


func (c *Client) FetchByQuery(query, page string) (*Results, error) {
	endpoint := fmt.Sprintf("https://newsapi.org/v2/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(query), c.PageSize, page, c.key)
	resp, err := c.http.Get(endpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{}
	return res, json.Unmarshal(body, res)
}


func (a *Article) FormatPublishedDate() string {
	year, month, day := a.PublishedAt.Date()
	return fmt.Sprintf("%v %d, %d", month, day, year)
}