package goprismic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Api struct {
	URL         string
	AccessToken string
	Data        ApiData
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
func Get(u, accessToken string) (*Api, error) {
	api := &Api{AccessToken: accessToken, URL: u}
	api.Data.Refs = make([]Ref, 0, 128)
	err := api.call(api.URL, map[string]string{}, &api.Data)
	return api, err
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

func (a *Api) call(u string, data map[string]string, res interface{}) error {
	callurl, err0 := url.Parse(u)
	if err0 != nil {
		return err0
	}
	values := callurl.Query()
	for k, v := range data {
		values.Add(k, v)
	}
	callurl.RawQuery = values.Encode()

	//fmt.Printf("call %s\n", callurl.String())
	req, err := http.NewRequest("GET", callurl.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Accept", "application/json")
	resp, err2 := http.DefaultClient.Do(req)
	defer resp.Body.Close()
	if err2 != nil {
		return err2
	}
	encoded, err3 := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(encoded))
	if err3 != nil {
		return err3
	}
	err4 := json.Unmarshal(encoded, res)
	//fmt.Printf("\n%+v\n", res)
	return err4
}
