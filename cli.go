package main

import "github.com/urfave/cli/v2"

type cliArgs struct {
	inputType      string
	outputLanguage string
	outputDirPath  string
	useCache       bool
	debug          bool
}

func getApp() cli.App {
	return cli.App{
		Name:   "citizenlab_translator",
		Usage:  "fetch and translate 'github.com/CitizenLabDotCo/citizenlab' source language files using the 'AWS Translate' service",
		Action: citizenlabTranslator,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "type",
				Aliases:  []string{"t"},
				Usage:    "Which file to translate (options: 'json' or 'yml')",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "lang",
				Aliases:  []string{"l"},
				Usage:    "Output language code (example: 'it')",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "output",
				Aliases:  []string{"o"},
				Usage:    "Output directory path where to save the translated file (example: 'output')",
				Required: true,
			},
			&cli.BoolFlag{
				Name:    "no-cache",
				Aliases: []string{"n"},
				Usage:   "Don't check the destination file for items already translated",
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Usage:   "Use debug logging",
			},
		},
	}
}

func validateCliArgs(cli *cliArgs) error {
	if cli.inputType != "json" && cli.inputType != "yml" {
		return ErrInvalidInputType
	}
	return nil
}

func parseCliArgs(ctx *cli.Context) *cliArgs {

	args := cliArgs{
		inputType:      ctx.String("type"),
		outputLanguage: ctx.String("lang"),
		outputDirPath:  ctx.String("output"),
		useCache:       !ctx.Bool("nocache"),
		debug:          ctx.Bool("debug"),
	}

	return &args
}
