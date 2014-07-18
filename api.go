package goprismic

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

type work struct {
	u    string
	data map[string]string
	res  interface{}

	ret chan error
}

type Api struct {
	URL         string
	AccessToken string
	Data        ApiData

	Config Config

	client    http.Client
	queue     chan work
	curSec    int
	reqCurSec int
}

type ApiData struct {
	Forms         map[string]Form   `json:"forms"`
	Refs          []Ref             `json:"refs"`
	Bookmarks     map[string]string `json:"bookmarks"`
	Tags          []string          `json:"tags"`
	Types         map[string]string `json:"types"`
	OAuthInitiate string            `json:"oauth_initiate"`
	OAuthToken    string            `json:"oauth_token"`
}

type Config struct {
	// Number of workers (simultaneous connections)
	Workers int
	// Timeout for HTTP requests
	Timeout time.Duration
	// Max requests per second (default: 0 - unlimited)
	ReqPerSec int
	// Debug mode
	Debug bool
}

// Default configuration
var DefaultConfig = Config{
	Workers: 3,
	Timeout: 5 * time.Second,
	ReqPerSec: 5,
}

// Api entry point
// Use Get(url, accessToken, DefaultConfig) to use the default config
func Get(u, accessToken string, cfg Config) (*Api, error) {
	api := &Api{
		AccessToken: accessToken,
		URL:         u,
		Config:      cfg,
		queue:       make(chan work),
		client:      http.Client{Timeout: cfg.Timeout},
		curSec:      -1,
	}
	api.Data.Refs = make([]Ref, 0, 128)
	err := api.call(api.URL, map[string]string{}, &api.Data)
	if err != nil {
		return nil, err
	}

	if cfg.Workers <= 0 {
		panic("Cannot run with no worker")
	}
	for workers := 0; workers < cfg.Workers; workers++ {
		go api.loopWorker()
	}
	return api, nil
}

// Refreshes the Api data
func (a *Api) Refresh() error {
	return a.call(a.URL, map[string]string{}, &a.Data)
}

// Fetches the master ref
func (a *Api) Master() *SearchForm {
	for _, r := range a.Data.Refs {
		if r.IsMasterRef {
			return a.createSearchForm(r)
		}
	}
	return &SearchForm{err: fmt.Errorf("Master ref not found !?!")}
}

// Fetch another ref
func (a *Api) Ref(label string) *SearchForm {
	for _, r := range a.Data.Refs {
		if r.Label == label {
			return a.createSearchForm(r)
		}
	}
	return &SearchForm{err: fmt.Errorf("No ref found with label '%s'", label)}
}

func (a *Api) createSearchForm(r Ref) *SearchForm {
	f := &SearchForm{api: a, ref: r}
	f.data = make(map[string]string)
	return f
}

func (a *Api) work(u string, data map[string]string, res interface{}) error {
	w := work{
		u:    u,
		data: data,
		res:  res,
		ret:  make(chan error),
	}
	a.queue <- w
	return <-w.ret
}

func (a *Api) loopWorker() {
	for {
		w := <-a.queue
		err := a.call(w.u, w.data, w.res)
		w.ret <- err
	}
}

func (a *Api) call(u string, data map[string]string, res interface{}) error {
	// test the number of requests per second
	if a.Config.ReqPerSec > 0 {
		curSec := time.Now().Second()
		if curSec != a.curSec {
			a.curSec = curSec
			a.reqCurSec = 0
		}
		a.reqCurSec++
		if a.reqCurSec > a.Config.ReqPerSec {
			if a.Config.Debug {
				log.Printf("Prismic - Too many requests")
			}
			return errors.New("Too many requests")
		}
	}
	// call
	callurl, errparse := url.Parse(u)
	if errparse != nil {
		return errparse
	}

	values := callurl.Query()
	for k, v := range data {
		values.Set(k, v)
	}
	if a.AccessToken != "" {
		values.Set("access_token", a.AccessToken)
	}
	callurl.RawQuery = values.Encode()
	req, errreq := http.NewRequest("GET", callurl.String(), nil)
	if errreq != nil {
		return errreq
	}
	req.Header.Add("Accept", "application/json")

	if a.Config.Debug {
		log.Printf("Prismic - requesting %s", callurl.String())
	}
	resp, errdo := a.client.Do(req)
	if errdo != nil {
		return errdo
	}

	defer resp.Body.Close()
	encoded, errread := ioutil.ReadAll(resp.Body)
	if errread != nil {
		return errread
	}
	if resp.StatusCode != 200 {
		err := new(PrismicError)
		errjson := json.Unmarshal(encoded, err)
		if errjson != nil {
			return errjson
		} else {
			return err
		}
	}

	errjson := json.Unmarshal(encoded, res)
	return errjson
}
