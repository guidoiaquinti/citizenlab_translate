package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/session"
	translate "github.com/aws/aws-sdk-go/service/translate"
	"github.com/urfave/cli/v2"
)

var urlMap = map[string]string{
	"json": "https://raw.githubusercontent.com/CitizenLabDotCo/citizenlab/master/front/app/translations/en.json",
	"yml":  "https://raw.githubusercontent.com/CitizenLabDotCo/citizenlab/master/back/config/locales/en.yml",
}

type translatorClient struct {
	awsTranslate *translate.Translate
	cliArgs      cliArgs
	log          *logrus.Logger
}

func main() {
	app := getApp()

	sort.Sort(cli.FlagsByName(app.Flags))

	err := app.Run(os.Args)
	if err != nil {
		fmt.Printf("\n%s\n", err)
	}
}

func citizenlabTranslator(ctx *cli.Context) error {
	var err error

	// Parse CLI args
	cliArgs := parseCliArgs(ctx)

	// Validate CLI args
	err = validateCliArgs(cliArgs)
	if err != nil {
		return err
	}

	// Setup logger
	log := logrus.New()
	if cliArgs.debug {
		log.Level = logrus.DebugLevel
	}

	// AWS client
	// TODO: find a way to early validate credentials
	// TODO: find a way to early validate if the target language is supported
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	awsTranslate := translate.New(sess)

	c := translatorClient{
		awsTranslate: awsTranslate,
		cliArgs:      *cliArgs,
		log:          log,
	}

	url := urlMap[cliArgs.inputType]
	log.Infof("Fetching source translation file from %q...", url)
	sourceFile, err := fetchSourceFile(url)
	if err != nil {
		log.Fatalf("Error: %s", err)
		panic(err)
	}

	log.Infof("Parsing file...")
	sourceFileParsed, err := parseFile(cliArgs.inputType, sourceFile)
	if err != nil {
		log.Fatalf("Error: %s", err)
		panic(err)
	}

	log.Infof("Translating the source file and saving it to the dir %q (with useCache: %t)...", cliArgs.outputDirPath, cliArgs.useCache)
	err = translateFile(c, sourceFileParsed)
	if err != nil {
		log.Fatalf("Error: %s", err)
		panic(err)
	}

	log.Infoln("Done!")

	return nil
}

func translateFile(c translatorClient, sourceFile map[string]interface{}) error {
	sourceFileFlatten := make(map[string]string)
	outputFileFlatten := make(map[string]string)
	var err error

	outputFilePath := filepath.Join(
		c.cliArgs.outputDirPath,
		fmt.Sprintf("%s.%s", c.cliArgs.outputLanguage, c.cliArgs.inputType),
	)

	// Let's flat the source file so that we don't go crazy traversing unstructured
	// and untrusted upstream inputs
	c.log.Infof("Flatting the source file...")
	sourceFileFlatten = flatInputContent(sourceFile)

	// If useCache, check if the output file exists and load its data
	if c.cliArgs.useCache {
		if fileExists(outputFilePath) {
			c.log.Debugf("Output file %q found, using it as cache...\n", outputFilePath)

			outputFile, err := os.Open(outputFilePath)
			if err != nil {
				return err
			}
			defer outputFile.Close()

			outputFileByteValue, _ := ioutil.ReadAll(outputFile)
			outputFileContent, err := parseFile(c.cliArgs.inputType, outputFileByteValue)
			if err != nil {
				return err
			}

			// Let's flat this file as well
			c.log.Infof("Flatting the output file...")
			outputFileFlatten = flatInputContent(outputFileContent)

		} else {
			c.log.Debugf("Output file %q to use as cache not found, starting from scratch...\n", outputFilePath)
		}
	}

	// For each entry in the source file
	for sourceKey, sourceValue := range sourceFileFlatten {

		if c.cliArgs.useCache {
			// Check if we already have an entry in the destination file and it's not null
			if val, ok := outputFileFlatten[sourceKey]; val != "" && ok {
				c.log.Debugf("Skip translating %q as it's already in the output file\n", sourceKey)
				continue
			}
		}

		// Translate item
		if sourceValue != "" {

			if len(sourceValue) < 5000 {
				c.log.Debugf("Translating %q...\n", sourceKey)
				translatedItem, err := translateText(c, sourceValue)
				if err != nil {
					return err
				}
				outputFileFlatten[sourceKey] = translatedItem
			} else {
				c.log.Debugf("Unable to translate %s as it's value is > 5000 chars...\n", sourceKey)
				// TODO: handle this
			}

		} else {
			outputFileFlatten[sourceKey] = ""
		}

		// Marshal and save the output at each iteration to don't waste
		// AWS Translate invocations in case of errors
		err = saveToFile(c.cliArgs.inputType, outputFileFlatten, outputFilePath)
		if err != nil {
			return err
		}
	}

	return nil
}
