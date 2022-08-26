package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

type Checktype struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Image        string      `json:"image"`
	Timeout      int         `json:"timeout"`
	Options      string      `json:"options,omitempty"`
	RequiredVars interface{} `json:"required_vars"`
	Assets       []string    `json:"assets"`
}

type Checktypes struct {
	Checktype []Checktype `json:"checktypes"`
}

// MergeChecktypes merges Checktypes first parameter into Checktypes second
// parameter returning a merged list of Checktypes.
func MergeChecktypes(ctda, ctdb Checktypes) Checktypes {
	// There are no performance enhancements of any kind.
	res := Checktypes{}
	for _, ctb := range ctdb.Checktype {
		res.Checktype = append(res.Checktype, ctb)
	}
	for _, cta := range ctda.Checktype {
		added := false
		for i := 0; i < len(res.Checktype); i++ {
			if res.Checktype[i].Name == cta.Name {
				res.Checktype[i] = cta
				added = true
				continue
			}
		}
		if !added {
			res.Checktype = append(res.Checktype, cta)
		}
	}
	return res
}

func FetchChecktypesFromURL(url string) (Checktypes, error) {
	httpClient := http.Client{
		Timeout: time.Second * 5,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Checktypes{}, err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return Checktypes{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return Checktypes{}, errors.New("unexpected status code fetching checktypes from URL")
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return Checktypes{}, err
	}
	ct := Checktypes{}
	err = json.Unmarshal(body, &ct)

	return ct, err
}
