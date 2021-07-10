package main

import (
	"bufio"
	"fmt"
	"github.com/kesselborn/go-getopt"
	"github.com/linuzilla/go-gloria-mailer/config"
	"github.com/linuzilla/go-gloria-mailer/excel"
	"github.com/linuzilla/go-gloria-mailer/flags"
	"github.com/linuzilla/go-gloria-mailer/mailer"
	"github.com/linuzilla/go-gloria-mailer/merger"
	mime_composer "github.com/linuzilla/go-gloria-mailer/mime-composer"
	"github.com/linuzilla/go-gloria-mailer/mimemail"
	"log"
	"os"
	"strconv"
	"strings"
)

const VERSION = "gloria-mailer version 0.0.1"
const ConfigFileEnv = "MAILER_CONFIG"
const DefaultConfigFile = "settings.conf"

func main() {
	optionDefinition := getopt.Options{
		Description: VERSION,
		Definitions: getopt.Definitions{
			{"debug|d", "debug mode: on", getopt.Optional | getopt.Flag | getopt.NoEnvHelp, false},
			{"no-debug", "debug mode: off", getopt.Optional | getopt.Flag | getopt.NoEnvHelp, false},
			{"send-email", "send email: on", getopt.Optional | getopt.Flag | getopt.NoEnvHelp, false},
			{"no-send-email", "send email: off", getopt.Optional | getopt.Flag | getopt.NoEnvHelp, false},
			{"excel", "excel file", getopt.Optional | getopt.NoEnvHelp, ``},
			{"template", "template file", getopt.Optional | getopt.NoEnvHelp, ``},
			{"verbose|v", "verbose mode", getopt.Optional | getopt.Flag | getopt.NoEnvHelp, false},
			{"config|c|" + ConfigFileEnv, "config file", getopt.IsConfigFile | getopt.ExampleIsDefault, DefaultConfigFile},
		},
	}

	options, _, _, e := optionDefinition.ParseCommandLine()

	help, wantsHelp := options["help"]
	exitCode := 0

	if e != nil || wantsHelp {
		switch {
		case wantsHelp && help.String == "usage":
			fmt.Print(optionDefinition.Usage())
		case wantsHelp && help.String == "help":
			fmt.Print(optionDefinition.Help())
		default:
			fmt.Println("**** Error: ", e.Error(), "\n", optionDefinition.Help())
			exitCode = e.ErrorCode
		}
	} else {
		start(options["config"].String, options)
	}
	os.Exit(exitCode)
}

func start(configurationFile string, options map[string]getopt.OptionValue) {
	settings := config.New(configurationFile)

	flags.Initialize(settings, options)

	fmt.Println(VERSION)
	fmt.Printf("\nExcel: [ %s ], Email Template: [ %s ]\n", flags.ExcelFile, flags.TemplateFile)
	fmt.Printf("Debugging: %s, Send REAL email: %s\n", strconv.FormatBool(flags.Debug), strconv.FormatBool(flags.SendEmail))
	fmt.Printf("SMTP Server: %s, port: %d, authentication: %s\n\n", settings.Smtp.Host, settings.Smtp.Port, strconv.FormatBool(settings.Smtp.Auth))

	reader, err := excel.New(settings.Excel.File)
	if err != nil {
		log.Fatal(err)
	}

	if flags.Verbose {
		fmt.Printf("Excel file: [ %s ], SheetName: [ %s ]\n", flags.ExcelFile, reader.SheetName())
		for _, fieldName := range reader.Fields() {
			fmt.Printf(">> Field: [ %s ]\n", fieldName)
		}
	}

	dataMerger := merger.New(reader.Fields())

	parser := mimemail.New(flags.TemplateFile)

	if flags.Verbose {
		fmt.Println("Subject: [", parser.Subject(), "]")
		fmt.Println("Media-Type: ", parser.MediaType())

		for _, part := range parser.Parts() {
			fmt.Println("   >>", part.ContentType)
			//fmt.Println(part.Body)
		}
	}

	//err = reader.Each(func(dataMap map[string]string) bool {
	//	fmt.Printf("Email: [ %s ], Name: [ %s ]\n",
	//		dataMap[settings.Excel.EmailColumn], dataMap[settings.Excel.NameColumn])
	//	return true
	//})
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	client := mailer.New(settings.Smtp)

	client.SetFrom(settings.Main.SenderEmail, settings.Main.SenderName)

	all := false
	err = reader.Each(func(dataMap map[string]string) bool {
		if flags.Verbose {
			for _, fieldName := range reader.Fields() {
				fmt.Printf(">> %s : [ %s ]\n", fieldName, dataMap[fieldName])
			}
			fmt.Println()
		}

		fmt.Printf("Email: [ %s ], Name: [ %s ]\n",
			dataMap[settings.Excel.EmailColumn], dataMap[settings.Excel.NameColumn])

		recipient := dataMap[settings.Excel.EmailColumn]

		if flags.Debug {
			recipient = settings.Debug.SendTo
		}

		if !all {
			fmt.Printf("\nDo you want to send to [ %s ] ?  [y(es)/N(o)/a(ll)/q(uit)]: ", recipient)

			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')

			fmt.Println()

			switch strings.TrimSpace(answer) {
			case "a", "A":
				all = true
			case "q", "Q":
				return false
			default:
				return true
			case "y", "Y":
			}
		}

		err := client.SendMailTo(
			recipient,
			dataMap[settings.Excel.NameColumn],
			parser.Subject(),
			func(composer mime_composer.MimeComposer) {
				for _, part := range parser.Parts() {
					composer.AddMultipartAlternative(part.ContentType,
						dataMerger.Merge(part.ContentType, part.Body, dataMap))
				}
			})

		if err != nil {
			log.Println(err)
		}
		return true
	})

	//reader.
}
