package main

import (
	"log"
	"os"
	"sort"
	"text/template"
)

type TemplateInfo struct {
	Items []TemplateItem
}

type TemplateItem struct {
	MacroName   string
	PDFFileName string
}

func processTemplate(entries map[string][]*Entry) {
	t := template.Must(template.New("letter").Delims("<<", ">>").Parse(latexStyleTemplate))

	templateList := make(map[string]*TemplateItem, 0)
	for _, entry := range entries {
		for _, e := range entry {
			item := new(TemplateItem)
			item.MacroName = e.macroName
			item.PDFFileName = e.pdfFileName
			templateList[e.macroName] = item
		}
	}

	keys := make([]string, 0, len(templateList))
	for k, _ := range templateList {
		keys = append(keys, k)

	}
	sort.Strings(keys)
	items := new(TemplateInfo)
	items.Items = make([]TemplateItem, 0)

	for _, k := range keys {
		entry := templateList[k]
		ti := new(TemplateItem)
		ti.MacroName = entry.MacroName
		ti.PDFFileName = entry.PDFFileName
		items.Items = append(items.Items, *ti)
	}

	err := t.Execute(os.Stdout, items)
	if err != nil {
		log.Println("executing template:", err)
	}
}

const latexStyleTemplate = `
%%  Copyright (C) 2021 by Glen Newton%%  This file may be distributed and/or modified under the%  conditions of the LaTeX Project Public License, either%  version 1.3 of this license or (at your option) any later%  version. The latest version of this license is in:%%  http://www.latex-project.org/lppl.txt%%  and version 1.3 or later is part of all distributions of%  LaTeX version 2005/12/01 or later.%\NeedsTeXFormat{LaTeX2e}
\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{azureicons}[2020/01/21 AWS Architectural Icons]

\RequirePackage{graphicx}
\RequirePackage[export]{adjustbox}

\newcommand{\assetZipFile}{Azure\_Public\_Service\_Icons\_V3.zip}
\newcommand{\assetZipFileSplit}{\seqsplit{Azure\_Public\_Service\_Icons\_V3.zip}}
\newcommand{\inkscapeVersion}{Inkscape 0.92.4 (5da689c313, 2019-01-14)}

%%%%%%%%%%%%%%%%
<<range .Items>>
\newcommand{\<<.MacroName>>}[1]{\includegraphics[width=#1,valign=t]{<<.PDFFileName>>}}
<<end>>

`
