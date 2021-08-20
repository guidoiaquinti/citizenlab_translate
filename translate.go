package main

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/awserr"
	translate "github.com/aws/aws-sdk-go/service/translate"
)

func translateText(tc translatorClient, input string) (string, error) {
	var ti translate.TextInput

	ti.SetSourceLanguageCode("en")
	ti.SetTargetLanguageCode(tc.cliArgs.outputLanguage)
	ti.SetText(input)

	result, err := tc.awsTranslate.Text(&ti)
	if err != nil {
		aerr, ok := err.(awserr.Error)
		if ok {
			switch errorCode := aerr.Code(); errorCode {
			case "MissingRegion":
				err = errors.New("Please specify the AWS region to use")
			case "NoCredentialProviders":
				err = errors.New("Please specify valid AWS credentials")
			}
		}
		return "", err
	}

	return *result.TranslatedText, nil
}
