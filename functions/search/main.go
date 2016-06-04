package main

import (
	"encoding/json"
	"os"
	"reflect"

	"github.com/apex/go-apex"

	"gopkg.in/olivere/elastic.v2"
)

type restaurant struct {
	RestaurantID string `json:"restaurant_id"`
	Name         string `json:"name"`
	NameAlphabet string `json:"name_alphabet"`
	NameKana     string `json:"name_kana"`
	Address      string `json:"address"`
	Description  string `json:"description"`
	Purpose      string `json:"purpose"`
	Category     string `json:"category"`
	PhotoCount   string `json:"photo_count"`
	MenuCount    string `json:"menu_count"`
	AccessCount  string `json:"access_count"`
	Closed       string `json:"closed"`
	Location     string `json:"location"`
}

type input struct {
	SearchWord string `json:"search_word"`
}

type output struct {
	Restaurants []restaurant `json:"result"`
}

func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		client, err := elastic.NewClient(
			elastic.SetSniff(false),
			elastic.SetURL(os.Getenv("ES_HOST")),
		)
		if err != nil {
			panic(err)
		}

		var in input
		if err := json.Unmarshal(event, &in); err != nil {
			panic(err)
		}

		q := elastic.NewQueryStringQuery(in.SearchWord)
		q = q.DefaultField("name")
		searchResult, err := client.Search().
			Index("ldgourmet").
			Query(q).
			Sort("name", true).
			From(0).Size(30).
			Pretty(true).
			Do()
		if err != nil {
			panic(err)
		}

		var out output
		var ttyp restaurant
		for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
			if t, ok := item.(restaurant); ok {
				out.Restaurants = append(out.Restaurants, t)
			}
		}

		return out, nil
	})
}
