package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

var argAssetZipFile string = "Azure_Public_Service_Icons_V3.zip"
var argConvertSvgWithInkscape = false
var argIconsFile string = "tex/icons.tex"
var argStyleFile string = "sty/azureicons.sty"
var argInkscapeBinPath string = "/usr/bin/inkscape"
var argInkscapeHelpArgs []string = []string{"--version"}
var argPdfTexOutDir string = "icons_tex"

type Entry struct {
	zipPathName     string
	pdfFileName     string
	pdfFileNamePath string
	numberName      string
	macroName       string
	englishName     string
	category        string
}

var entries map[string][]*Entry

func main() {
	iconsFile, err := os.Create(argIconsFile)
	if err != nil {
		log.Panic(err)
	}

	wicons := bufio.NewWriter(iconsFile)
	defer func() {
		if err = wicons.Flush(); err != nil {
			log.Panic(err)
		}
		if err := iconsFile.Close(); err != nil {
			log.Panic(err)
		}
	}()

	entries, err := processAzureIcons(argAssetZipFile)
	if err != nil {
		log.Fatal(err)
	}
	printEntries(wicons, entries)

	processTemplate(entries)

}

func processAzureIcons(src string) (map[string][]*Entry, error) {
	rz, err := zip.OpenReader(src)
	if err != nil {
		return nil, err
	}
	defer rz.Close()

	entries := make(map[string][]*Entry)

	for _, f := range rz.File {
		if f.FileInfo().IsDir() || !strings.HasSuffix(f.Name, ".svg") {
			continue
		}
		entry := new(Entry)
		log.Println("")
		entry.zipPathName = f.Name
		log.Println(f.Name)
		parts := strings.Split(f.Name, "/")
		//log.Println(parts)
		entry.category = parts[2]
		//log.Println("Category:", parts[2])
		//log.Println(parts[3])

		fparts := strings.Split(parts[3], "-")
		entry.numberName = fparts[0]
		entry.macroName = replaceNumberWithStrings(contractNames(fixName(strings.TrimSuffix(strings.Join(fparts[3:len(fparts)], ""), ".svg"))))
		log.Println("Number:", fparts[0])
		entry.englishName = fixName(strings.TrimSuffix(strings.Join(fparts[3:len(fparts)], " "), ".svg"))
		entry.pdfFileName = entry.numberName + "_" + fixName(strings.TrimSuffix(strings.Join(fparts[3:len(fparts)], ""), "svg")+"pdf")
		entry.pdfFileNamePath = argPdfTexOutDir + "/" + entry.pdfFileName
		log.Println("PDFFileName:", entry.pdfFileName)

		// category exists already?
		var ok bool
		var categoryList []*Entry
		if categoryList, ok = entries[entry.category]; ok {
			categoryList = append(categoryList, entry)
		} else {
			categoryList = make([]*Entry, 1)
			categoryList[0] = entry
		}
		entries[entry.category] = categoryList

		if argConvertSvgWithInkscape {
			var rc io.ReadCloser

			// Open file (in zip stream)
			rc, err = f.Open()
			if err != nil {
				log.Fatal(err)
			}

			e, err := runCommandReadStdin(rc, argInkscapeBinPath, "--file=-", "-D", "-z", "--export-latex", "--export-pdf="+entry.pdfFileNamePath)
			if err != nil {
				log.Println(e)
				log.Fatal(err)
			}
			// Close the file without defer to close before next iteration of loop
			err = rc.Close()

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return entries, nil
}

var fileNameBadChars = []string{"'", "(", ")"}

func fixName(s string) string {
	for i := 0; i < len(fileNameBadChars); i++ {
		s = strings.ReplaceAll(s, fileNameBadChars[i], "")
	}
	return s
}

func printEntries(w io.Writer, entries map[string][]*Entry) {
	keys := make([]string, 0, len(entries))
	for k, _ := range entries {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := entries[k]
		sort.Sort(ByName(v))
		fmt.Fprintln(w, "\\subsection{"+k+"}")
		for _, e := range v {
			//fmt.Fprintln(w, strings.ReplaceAll(e.pdfFileName, "_", "\\_"))
			pdfFileNameTeX := strings.ReplaceAll(e.pdfFileName, "_", "\\_")
			include := "\\includegraphics[width=1cm,valign=t]{" + e.pdfFileName + "}"
			fmt.Fprintf(w, "\\gxs{%s}{%s}{%s}{%s}{%s}\n\n", e.englishName, include, pdfFileNameTeX, e.macroName, e.category)
			//fmt.Fprintln(w, e.englishName)
			//fmt.Fprintln(w, e.macroName)
			//fmt.Fprintln(w, "\\includegraphics[width=1cm,valign=t]{"+e.pdfFileName+"}\\ \\  ")
		}

	}
}

func contractNames(s string) string {
	for i := 0; i < len(macroSubs); i += 2 {
		s = strings.ReplaceAll(s, macroSubs[i], macroSubs[i+1])
	}
	return s
}

// LaTeX does not allow numbers in macro names!!!!
// http://www.texfaq.org/FAQ-linmacnames
var numMap = map[string]string{"0": "Zero", "1": "One", "2": "Two", "3": "Three", "4": "Four", "5": "Five", "6": "Six", "7": "Seven", "8": "Eight", "9": "Nine"}

// To shorten macro names
var macroSubs = []string{

	"Alternate", "Altern",
	"Application", "App",
	"Archive", "Archv",
	"Attribute", "Att",
	"Attributes", "Atts",
	"Aurora", "Aur",
	"Authority", "Auth",
	"Azure", "Azr",
	"Account", "Accnt",
	"Bucket", "Bkt",
	"Certificate", "Cert",
	"Compute", "Cmput",
	"Credential", "Cred",
	"Database", "DB",
	"Database", "DB",
	"Default", "Dflt",
	"Development", "Dev",
	"Device", "Dev",
	"Directory", "Dir",
	"Distribution", "Distr",
	"Elemental", "Elem",
	"Emergency", "Emerg",
	"Encrypted", "Encrd",
	"Encryption", "Encr",
	"Enterprise", "Entrprs",
	"Function", "Func",
	"General", "Gen",
	"Global", "Glbl",
	"Governance", "Govern",
	"Groups", "Grp",
	"Internet", "INet",
	"Identity", "Ident",
	"Image", "Img",
	"Image", "Img",
	"Instance", "Inst",
	"Integration", "Integ",
	"Intelligent", "Intell",
	"Kinesis", "Kin",
	"Kubernetes", "Kuber",
	"License", "Licns",
	"LoadBalancer", "LB",
	"Managed", "Mngd",
	"Management", "Mngmt",
	"Manager", "Mngr",
	"Medical", "Med",
	"Microsoft", "MS",
	"Migration", "Migrat",
	"Mobile", "Mob",
	"Multiple", "Mult",
	"Namespaces", "NS",
	"Notification", "Notif",
	"Object", "Obj",
	"Organizational", "Org",
	"Page", "Pg",
	"Points", "Pts",
	"Providers", "Prvdr",
	"Profile", "Prof",
	"Recovery", "Recov",
	"Recovery", "Recov",
	"Registry", "Reg",
	"Rendering", "Rend",
	"Replication", "Repli",
	"Resource", "Res",
	"Search", "Srch",
	"Source", "Src",
	"Security", "Sec",
	"Service", "Svc",
	"Simulation", "Simul",
	"Storage", "Stor",
	"Stream", "Strm",
	"SystemManager", "SysMgr",
	"Thing", "Thng",
	"ThinkBox", "TB",
	"TimeSeries", "TS",
	"TrainingCertification", "TrainCert",
	"Universal", "Univer",
	"Vault", "Vlt",
	"Virtual", "Virt",
	"Volume", "Vol",
	"Windows", "Win",
	"Warning", "Warn",
}

func replaceNumberWithStrings(s string) string {
	for k, v := range numMap {
		s = strings.ReplaceAll(s, k, v)
	}
	return s
}
