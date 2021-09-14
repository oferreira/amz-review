package main

import (
	"context"
	"fmt"

	"amazon.com/review/services/datasource"
	"amazon.com/review/services/review"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	//"github.com/aws/aws-lambda-go/lambda"
)



func translate () 

func main() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile("default"));
	
	rows, err := datasource.New()
    if err != nil {
		panic(err)
    }

	for index, row := range rows {
		if index == 0 {
			continue
		} 

		result, err := 	review.Fetch(row[0])
		if err != nil {	
			panic(err)
		}

		client := translate.NewFromConfig(cfg)

		response, err := client.TranslateText(context.Background(), &translate.TranslateTextInput{
			SourceLanguageCode: aws.String("en"),
			TargetLanguageCode: aws.String("fr"),
			Text:               &result.Content,
		})
		if err != nil {
			panic(err)
		}
		fmt.Println(*response.TranslatedText);
		panic("dd")
	}
	
	if(err != nil) {
		fmt.Println(err);
	}

	

	// client.TranslateText(context.Background(), )
	
}