package wtmfactor

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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
	GetFactor() []Algorithm
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

func (w *wtmFactor) GetFactor() []Algorithm {
	req, err := http.NewRequest("GET", w.Config.URL, nil)
	client := &http.Client{
		Timeout: time.Second * 60,
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
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
