package review

import (
	"fmt"
	"context"
	"strings"
    "regexp"
    "net/url"
	

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/go-rod/rod"
)

type Review struct{ 
	ID string;
	ASIN string;
	URL string; 
	Weight int;
	Username string; 
	Avatar string; 
	Rate string; 
	Title string;
	Text string;
	Date string;
	Data string;
	Helpful string;	
}

type ReviewTranslate struct {
    *Review
    TranslatedTitle string
    TranslatedText string
    TranslatedDate string
    TranslatedData string
    TranslatedHelpful string
}

func Fetch (review *Review) {
	fmt.Println(review.URL);
	page := rod.New().MustConnect().MustPage(review.URL)
    page.MustWaitLoad().MustScreenshot("screenshots/" + review.ID + ".png")
	review.Username = page.MustElement(".a-profile-content > span").MustText()
	review.Avatar = page.MustElement(".a-profile-avatar").MustHTML()
	review.Rate = page.MustElement(".review-rating").MustHTML()
	review.Title = page.MustElement(".review-title").MustText()
	review.Text = page.MustElement(".review-text-content > span").MustText()
	review.Date = page.MustElement(".review-date").MustHTML()
	review.Data = page.MustElement(".review-data").MustHTML()
	review.Helpful = page.MustElement(".cr-vote-text").MustText()

	fmt.Println("----------------------------------")
	fmt.Println(review.Text);
}


func Translate (reviewUrl string,wieght int, c *chan ReviewTranslate) {
	u, err := url.Parse(reviewUrl)
    if err != nil {
		panic(err)
    }

	// get ASIN from url
    m, _ := url.ParseQuery(u.RawQuery)
    asin := m["ASIN"][0];

	// get id from url
	s := strings.Split(u.Path, "/")
	id := s[3] 

	review := Review{id, asin, reviewUrl, 1,"", "", "", "", "", "", "", "",}
	Fetch(&review)

	fmt.Println("----------------------------------")
	fmt.Println(review); 
	sourceLanguageCode := aws.String("en")
	 
	isFrench, _ := regexp.MatchString("amazon\\.fr", review.URL)
    if isFrench == true {	
		sourceLanguageCode = aws.String("fr")
	}

	if err != nil {	
		panic(err)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile("default"));
	if err != nil {	
		panic(err)
	}
	client := translate.NewFromConfig(cfg)
	title, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: aws.String("es"),
		Text:               &review.Title,
	})
	if err != nil {
		panic(err)
	}
	
	text, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: aws.String("es"),
		Text:               &review.Text,
	})
	if err != nil {
		panic(err)
	}

	
	date, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: aws.String("es"),
		Text:               &review.Date,
	})
	if err != nil {
		panic(err)
	}


	data, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: aws.String("es"),
		Text:               &review.Data,
	})
	fmt.Println("----------------------------------")
	fmt.Println(*data.TranslatedText);
	if err != nil {
		panic(err)
	}
	helpful, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: aws.String("es"),
		Text:               &review.Helpful,
	})
	if err != nil {
		panic(err)
	}

	*c<- ReviewTranslate{
		Review: &review,
		TranslatedTitle: *title.TranslatedText,
		TranslatedText: *text.TranslatedText,
		TranslatedDate: *date.TranslatedText,
		TranslatedData: *data.TranslatedText,
		TranslatedHelpful: *helpful.TranslatedText,
	}
	return
}

