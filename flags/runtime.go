package flags

import (
	"github.com/kesselborn/go-getopt"
	"github.com/linuzilla/go-gloria-mailer/config"
)

var Debug = false
var Verbose = false
var SendEmail = false
var ExcelFile = ``
var TemplateFile = ``

func Initialize(settings config.Settings, options map[string]getopt.OptionValue) {
	if settings.Debug.Debugging {
		Debug = true
	}

	if settings.Main.SendEmail {
		SendEmail = true
	}

	ExcelFile = settings.Excel.File
	TemplateFile = settings.Main.Template

	if options[`send-email`].Bool {
		SendEmail = true
	}

	if options[`no-send-email`].Bool {
		SendEmail = false
	}

	if options[`no-debug`].Bool {
		Debug = false
	}

	if options[`debug`].Bool {
		Debug = true
	}

	if options[`verbose`].Bool {
		Verbose = true
	}

	if options[`excel`].Set {
		ExcelFile = options[`excel`].String
	}

	if options[`template`].Set {
		TemplateFile = options[`template`].String
	}
}
