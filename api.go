package goprismic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type work struct {
	u string
	data map[string]string
	res interface{}

	ret chan error
}

type Api struct {
	URL         string
	AccessToken string
	Data        ApiData

	queue chan work
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

// Api entry point
func Get(u, accessToken string, workers int) (*Api, error) {
	api := &Api{
		AccessToken: accessToken,
		URL: u,
		queue: make(chan work),
	}
	api.Data.Refs = make([]Ref, 0, 128)
	err := api.call(api.URL, map[string]string{}, &api.Data)
	if err != nil {
		return nil, err
	}
	if workers <= 0 {
		panic("Cannot run with no worker")
	}
	for workers > 0 {
		go api.loopWorker()
		workers--
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
		u: u,
		data: data,
		res: res,
		ret: make(chan error),
	}
	a.queue <- w
	return <- w.ret
}

func (a *Api) loopWorker() {
	for {
		w := <- a.queue
		err := a.call(w.u, w.data, w.res)
		w.ret <- err
	}
}

func (a *Api) call(u string, data map[string]string, res interface{}) error {
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

	resp, errdo := http.DefaultClient.Do(req)
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
