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
	Name          string  `json:"name"`
	HashrateID    string  `json:"hashrate_id"`
	HashrateUnit  string  `json:"hashrate_unit"`
	HashrateValue float64 `json:"hashrate_value"`
	PowerID       string  `json:"power_id"`
	PowerValue    float64 `json:"power_value"`
	PowerUnit     string  `json:"power_unit"`
}

// WtmConfig config a url
type WtmConfig struct {
	URL string
}

// WtmFactor is the interface
type WtmFactor interface {
	GetFactorName() map[string]Algorithm
	GetFactorKey() map[string]map[string]float64
	GetFactory() map[string]Algorithm
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

func (w *wtmFactor) GetFactorKey() map[string]map[string]float64 {
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

	data := map[string]map[string]float64{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		log.Println(err)
	}
	return data
}

func (w *wtmFactor) GetFactorName() map[string]Algorithm {
	client := NewClientRequest()
	resp := client.Core(w.Config.URL, "GET")
	defer resp.Body.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	algor := map[string]Algorithm{}
	doc.Find(".form-row .py-1").Each(func(i int, s *goquery.Selection) {

		algoName := strings.TrimSpace(s.Find("label.ck-button span.btn.btn-default.btn-block.mb-1").Text())
		hr := s.Find("div.input-group.input-group-sm")

		hashrateID := ""
		hashrateUnit := ""
		powerID := ""
		powerUnit := ""
		algoKey := ""

		hr.Each(func(i int, s *goquery.Selection) {
			if id, idEx := s.Find("input.form-control").Attr("id"); idEx {
				hrPref := "factor_"
				hrSuff := "_hr"
				pPref := "factor_"
				pSuff := "_p"

				unit := strings.TrimSpace(s.Find("span.input-group-text").Text())

				if strings.Contains(id, hrPref) && strings.Contains(id, hrSuff) {
					hashrateID = id
					hashrateUnit = unit
					algoKey = id[0 : len(id)-3]
				} else if strings.Contains(id, pPref) && strings.Contains(id, pSuff) {
					powerID = id
					powerUnit = unit
					algoKey = id[0 : len(id)-2]
				}
			}
		})

		algor[algoKey] = Algorithm{
			Name:         algoName,
			HashrateID:   hashrateID,
			HashrateUnit: hashrateUnit,
			PowerID:      powerID,
			PowerUnit:    powerUnit,
		}

	})

	return algor
}

func (w *wtmFactor) GetFactory() map[string]Algorithm {
	algorithms := w.GetFactorName()
	keys := w.GetFactorKey()

	vga := "adapt_1080Ti"
	cardAmount := 1.0

	H := 0.0 // adapt_380
	B := 0.0 // adapt_fury
	R := 0.0 // adapt_470
	M := 0.0 // adapt_480
	F := 0.0 // adapt_570
	W := 0.0 // adapt_580
	U := 0.0 // adapt_vega56
	z := 0.0 // adapt_vega64
	V := 0.0 // adapt_vii
	C := 0.0 // adapt_1050Ti
	k := 0.0 // adapt_10606
	S := 0.0 // adapt_1070
	A := 0.0 // adapt_1070Ti
	O := 0.0 // adapt_1080
	D := 0.0 // adapt_1080Ti
	N := 0.0 // adapt_1660
	j := 0.0 // adapt_1660Ti
	q := 0.0 // adapt_2060
	L := 0.0 // adapt_2070
	I := 0.0 // adapt_2080
	P := 0.0 // adapt_2080Ti

	if vga == "adapt_380" {
		H = cardAmount
	} else if vga == "adapt_fury" {
		B = cardAmount
	} else if vga == "adapt_470" {
		R = cardAmount
	} else if vga == "adapt_480" {
		M = cardAmount
	} else if vga == "adapt_570" {
		F = cardAmount
	} else if vga == "adapt_580" {
		W = cardAmount
	} else if vga == "adapt_vega56" {
		U = cardAmount
	} else if vga == "adapt_vega64" {
		z = cardAmount
	} else if vga == "adapt_vii" {
		V = cardAmount
	} else if vga == "adapt_1050Ti" {
		C = cardAmount
	} else if vga == "adapt_10606" {
		k = cardAmount
	} else if vga == "adapt_1070" {
		S = cardAmount
	} else if vga == "adapt_1070Ti" {
		A = cardAmount
	} else if vga == "adapt_1080" {
		O = cardAmount
	} else if vga == "adapt_1080Ti" {
		D = cardAmount
	} else if vga == "adapt_1660" {
		N = cardAmount
	} else if vga == "adapt_1660Ti" {
		j = cardAmount
	} else if vga == "adapt_2060" {
		q = cardAmount
	} else if vga == "adapt_2070" {
		L = cardAmount
	} else if vga == "adapt_2080" {
		I = cardAmount
	} else if vga == "adapt_2080Ti" {
		P = cardAmount
	}

	for key := range keys["v"] {
		at := H * keys["m"][key] // adapt_380
		ut := B * keys["x"][key] // adapt_fury
		st := R * keys["g"][key] // adapt_470
		ct := M * keys["v"][key] // adapt_480
		lt := F * keys["y"][key] // adapt_570
		ft := W * keys["b"][key] // adapt_580
		ht := U * keys["w"][key] // adapt_vega56
		pt := z * keys["T"][key] // adapt_vega64
		dt := V * keys["E"][key] // adapt_vii
		K := C * keys["o"][key]  // adapt_1050Ti
		Y := k * keys["i"][key]  // adapt_10606
		Q := S * keys["a"][key]  // adapt_1070
		X := A * keys["s"][key]  // adapt_1070Ti
		J := O * keys["c"][key]  // adapt_1080
		Z := D * keys["l"][key]  // adapt_1080Ti
		tt := N * keys["f"][key] // adapt_1660
		et := j * keys["u"][key] // adapt_1660Ti
		nt := q * keys["h"][key] // adapt_2060
		rt := L * keys["p"][key] // adapt_2070
		ot := I * keys["d"][key] // adapt_2080
		it := P * keys["_"][key] // adapt_2080Ti

		result := at + ut + st + ct + lt + ft + ht + pt + dt + K + Y + Q + X + J + Z + tt + et + nt + rt + ot + it

		hrPref := "factor_"
		hrSuff := "_hr"
		pPref := "factor_"
		pSuff := "_p"
		if strings.Contains(key, hrPref) && strings.Contains(key, hrSuff) {
			algoKey := key[1 : len(key)-3]
			alg := algorithms[algoKey]
			alg.HashrateValue = result
			algorithms[algoKey] = alg
		} else if strings.Contains(key, pPref) && strings.Contains(key, pSuff) {
			algoKey := key[1 : len(key)-2]
			alg := algorithms[algoKey]
			alg.PowerValue = result
			algorithms[algoKey] = alg
		}

	}

	return algorithms
}
