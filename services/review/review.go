package review

import (
	"net/url"
	"strings"

	"github.com/go-rod/rod"
)

type Review struct{ ID string; Asin string; Username string; Content string;  }

type ReviewTranslate Review && struct {
    *Review
    ContentTranslate string
}

func Fetch (reviewUrl string) (Review, error) {
	u, err := url.Parse(reviewUrl)
    if err != nil {
		return Review{}, err
    }

	// get Asin from url
    m, _ := url.ParseQuery(u.RawQuery)
    Asin := m["ASIN"][0];

	// get id from url
	s := strings.Split(u.Path, "/")
	id := s[3] 

	page := rod.New().MustConnect().MustPage(reviewUrl)
    page.MustWaitLoad().MustScreenshot(id + ".png")
	content := page.MustElement(".review-text-content > span")
	username := page.MustElement(".a-profile-content > span")

	return Review{id, Asin, username.MustText(), content.MustText()}, nil
}



func Translate (c *chan review) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile("default"));
	
}

