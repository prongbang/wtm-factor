package wtmfactor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ClientRequest is the interface
type ClientRequest interface {
	Get(url string) string
	Core(url string, method string) *http.Response
}

type clientRequest struct {
}

func (c *clientRequest) Get(url string) string {
	resp := c.Core(url, "GET")
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body)
}

func (c *clientRequest) Core(url string, method string) *http.Response {
	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	return resp
}

// NewClientRequest is new instance
func NewClientRequest() ClientRequest {
	return &clientRequest{}
}

// Algorithm is a model
type Algorithm struct {
	Name       string `json:"name"`
	HashrateID string `json:"hashrate_id"`
	PowerID    string `json:"power_id"`
}

// WtmConfig config a url
type WtmConfig struct {
	URL string
}

// WtmFactor is the interface
type WtmFactor interface {
	GetFactorName() []Algorithm
	GetFactorKey() map[string]interface{}
}

type wtmFactor struct {
	Config WtmConfig
}

// NewWtmFactor provide function
func NewWtmFactor(config WtmConfig) WtmFactor {
	return &wtmFactor{
		Config: config,
	}
}

func (w *wtmFactor) GetFactorKey() map[string]interface{} {
	client := NewClientRequest()

	html := client.Get(w.Config.URL)
	assets := "/assets/application-"
	prefixScript := fmt.Sprintf(`<script src="%s`, assets)
	suffixScript := `.js"></script>`
	prefixScriptIdx := strings.Index(html, prefixScript)
	html = html[prefixScriptIdx+len(prefixScript):]
	suffixScriptIdx := strings.Index(html, suffixScript)
	hash := html[0:suffixScriptIdx]

	jsURL := fmt.Sprintf("%s%s%s.js", w.Config.URL, assets, hash)
	log.Println("GET -->", jsURL)

	js := client.Get(jsURL)
	prefix := `m={"#factor`
	suffix := `r=Object.keys(v)`
	prefixIdx := strings.Index(js, prefix)
	js = js[prefixIdx:]
	suffixIdx := strings.Index(js, suffix)
	js = js[0:suffixIdx]

	jsonData := ""
	adapt := strings.Split(js, "},")
	for i, a := range adapt {
		aK := "={"
		aIdx := strings.Index(a, aK)
		if aIdx != -1 {
			key := a[0:aIdx]
			a = a[aIdx+len(aK):]
			comma := ","
			if i == len(adapt)-2 {
				comma = ""
			}
			jsonData += fmt.Sprintf(`"%s":{%s}%s`, key, a, comma)
		}
	}
	jsonData = strings.Replace(jsonData, ":.", ":0.", -1)
	jsonData = fmt.Sprintf("{%s}", jsonData)

	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Println(err)
	}
	return data
}

func (w *wtmFactor) GetFactorName() []Algorithm {
	client := NewClientRequest()
	resp := client.Core(w.Config.URL, "GET")
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	algor := []Algorithm{}
	doc.Find(".form-row .py-1").Each(func(i int, s *goquery.Selection) {

		algoName := strings.TrimSpace(s.Find("label.ck-button span.btn.btn-default.btn-block.mb-1").Text())
		hr := s.Find("div.input-group.input-group-sm")

		hashrateID := ""
		powerID := ""

		hr.Each(func(i int, s *goquery.Selection) {
			if id, idEx := s.Find("input.form-control").Attr("id"); idEx {
				hrPref := "factor_"
				hrSuff := "_hr"
				pPref := "factor_"
				pSuff := "_p"

				if strings.Contains(id, hrPref) && strings.Contains(id, hrSuff) {
					hashrateID = id
				} else if strings.Contains(id, pPref) && strings.Contains(id, pSuff) {
					powerID = id
				}
			}
		})

		algor = append(algor, Algorithm{
			Name:       algoName,
			HashrateID: hashrateID,
			PowerID:    powerID,
		})
	})

	return algor
}
