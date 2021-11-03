package main

import (
	"time"
	"fmt"
	"sort"

	"amazon.com/review/services/datasource"
	"amazon.com/review/services/review"
	"encoding/json"
	"io/ioutil"
)

type Response struct {
	ID string;
	ASIN string; 
	Username string; 
	Avatar string; 
	Rate string; 
	Title string;
	TranslatedTitle string;
	Text string;
	TranslatedText string;
	Date string;
	TranslatedDate string;
	Size string;
	TranslatedSize string;
	Purchase string;
	TranslatedPurchase string;
	Helpful string;
	TranslatedHelpful string;
	Weight int;
}

func main() {
	var responces []Response

	rows, err := datasource.New()
    if err != nil {
		panic(err)
    }

	c := make(chan review.ReviewTranslate)

	for index, row := range rows {
		if index == 0 {
			continue
		} 
		
		go review.Translate(row[0], index, &c)
	}
	
	for index, _ := range rows {
		if index == 0 {
			continue
		} 
		
		select {
		case review:= <-c:
			responces = append(responces, Response{
				ID: review.ID,
				ASIN: review.ASIN,
				Username: review.Username,
				Avatar: review.Avatar,
				Rate: review.Rate,
				Title: review.Title,
				TranslatedTitle: review.TranslatedTitle,
				Text: review.Text,
				TranslatedText: review.TranslatedText,
				Date: review.Date,
				TranslatedDate: review.TranslatedDate,
				Size: review.Size,
				TranslatedSize: review.TranslatedSize,
				Purchase: review.Purchase,
				TranslatedPurchase: review.TranslatedPurchase,
				Helpful: review.Helpful,
				TranslatedHelpful: review.TranslatedHelpful,
				Weight: review.Weight,
			})
		case <-time.After(2*time.Minute):
			fmt.Println("Ne répond pas")
		}
	}


	sort.SliceStable(responces, func(i, j int) bool {
		return responces[i].Weight < responces[j].Weight
	})

	file, _ := json.MarshalIndent(responces, "", " ")
	_ = ioutil.WriteFile("output.json", file, 0644)
	return 
}