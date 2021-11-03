package review

import (
	"fmt"
	"context"
	"strings"
    "regexp"
    "net/url"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/translate"
	"github.com/go-rod/rod"
	// "github.com/go-rod/rod/lib/utils"
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
	Size string;
	Purchase string;
	Helpful string;	
}

type ReviewTranslate struct {
    *Review
    TranslatedTitle string
    TranslatedText string
    TranslatedDate string
    TranslatedSize string
    TranslatedPurchase string
    TranslatedHelpful string
}

func Fetch (review *Review) {
	browser := rod.New().MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(review.URL)

	// sleep for 0.5 seconds before every retry
	// sleeper := func() utils.Sleeper {
	// 	return func(context.Context) error {
	// 		time.Sleep(time.Second / 2)
	// 		return nil
	// 	}
	// }

    page.MustWaitLoad().MustScreenshot("screenshots/" + review.ID + ".png")
	review.Username = page.MustElement(".a-profile-content > span").MustText()
	review.Avatar = *page.MustElement(".a-profile-avatar img").MustAttribute("src")
	// review.Rate = page.MustElement(".review-rating").MustHTML()
	review.Title = page.MustElement(".review-title").MustText()
	review.Text = page.MustElement(".review-text-content > span").MustText()
	review.Date = page.MustElement(".review-date").MustText()
	// review.Size = page.Sleeper(sleeper).MustElement(".review-data > a").MustText() + " | " + *page.MustElement(".review-data > a").MustAttribute("href") 
	// review.Purchase = page.Sleeper(sleeper).MustElement(".review-data > span a span").MustText() + " | " + *page.MustElement(".review-data > span a").MustAttribute("href")
	// review.Helpful = page.Sleeper(sleeper).MustElement(".cr-vote-text").MustText()
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

	review := Review{id, asin, reviewUrl, wieght,"", "", "", "", "", "", "", "", "",}
	Fetch(&review)

	fmt.Println("----------------------------------")
	fmt.Println(review); 
	sourceLanguageCode := aws.String("en")
	 
	isFrench, _ := regexp.MatchString("amazon\\.fr", review.URL)
    if isFrench == true {	
		sourceLanguageCode = aws.String("fr")
	}

	if err != nil {	
		fmt.Println(review.URL); 
		panic(err)
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile("default"));
	if err != nil {	
		fmt.Println(review.URL); 
		panic(err)
	}
	client := translate.NewFromConfig(cfg)
	title, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: aws.String("es"),
		Text:               &review.Title,
	})
	if err != nil {
		fmt.Println(review.URL); 
		panic(err)
	}

	time.Sleep(20 * time.Second)
	
	text, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: aws.String("es"),
		Text:               &review.Text,
	})
	if err != nil {
		fmt.Println(review.URL); 
		panic(err)
	}

	time.Sleep(20 * time.Second)

	date, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
		SourceLanguageCode: sourceLanguageCode,
		TargetLanguageCode: aws.String("es"),
		Text:               &review.Date,
	})
	if err != nil {
		fmt.Println(review.URL); 
		panic(err)
	}

	time.Sleep(20 * time.Second)

	// size, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
	// 	SourceLanguageCode: sourceLanguageCode,
	// 	TargetLanguageCode: aws.String("es"),
	// 	Text:               &review.Size,
	// })
	// if err != nil {
	// 	fmt.Println(review.URL); 
	// 	panic(err)
	// }

	// purchase, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
	// 	SourceLanguageCode: sourceLanguageCode,
	// 	TargetLanguageCode: aws.String("es"),
	// 	Text:               &review.Purchase,
	// })
	// if err != nil {
	// 	fmt.Println(review.URL); 
	// 	panic(err)
	// }

	// helpful, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
	// 	SourceLanguageCode: sourceLanguageCode,
	// 	TargetLanguageCode: aws.String("es"),
	// 	Text:               &review.Helpful,
	// })
	// if err != nil {
	// 	panic(err)
	// }

	*c<- ReviewTranslate{
		Review: &review,
		TranslatedTitle: *title.TranslatedText,
		TranslatedText: *text.TranslatedText,
		TranslatedDate: *date.TranslatedText,
		// TranslatedSize: *size.TranslatedText,
		// TranslatedPurchase: *purchase.TranslatedText,
		// TranslatedHelpful: *helpful.TranslatedText,
	}
	return
}

