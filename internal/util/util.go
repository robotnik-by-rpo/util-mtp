package util

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/solworktech/md2pdf/v2"
	"util/md/utils/exist"
)

// Struct for storing flag from CLI
type mtp struct{
	SavePath string
	Input string
	Output string
	Title string
	PageSize string
	Orientation string
	Theme string
	Author string
	IsFooter bool
	IsNewPage bool
	Font string
	FontSize float64
	Ru bool
}


// Function initing of mtp struct
func initMTP() (*mtp,error){
	newMTP := &mtp{}
	newMTP.SavePath = os.Getenv("SAVE_PATH")
	flag.StringVar(&newMTP.Input,"i","","Input md file")
	flag.StringVar(&newMTP.Output, "o","./","Output name of PDF file")
	flag.StringVar(&newMTP.Title,"t","","Title for pdf file")
	flag.StringVar(&newMTP.PageSize,"ps","A4","Page size of PDF file")
	flag.StringVar(&newMTP.Orientation,"or","portrait","Orientation of PDF file: portrait or landscape")
	flag.StringVar(&newMTP.Theme,"th","default","Theme of PDF file")
	flag.StringVar(&newMTP.Author,"a","","You can point author for footer")
	flag.BoolVar(&newMTP.IsFooter,"if",false,"Turn on footer")
	flag.BoolVar(&newMTP.IsNewPage,"ip",false,"Create new page after ---")
	flag.StringVar(&newMTP.Font, "f","Helvetica","Font fot text in PDF file")
	flag.Float64Var(&newMTP.FontSize,"fs",14.0,"Size for text in PDF file")
	flag.BoolVar(&newMTP.Ru,"r",false,"Turn on russion language")
	flag.Parse()

	if newMTP.Input == ""{
		log.Println("Flag -i must have meaning.")
		return nil, errors.New("Flag input wasn't pointed.")
	}
	if !strings.HasSuffix(newMTP.Output, ".pdf") {
		newMTP.Output = newMTP.Output + ".pdf"
	}
	
	return newMTP, nil
}

// Function for running CLI. It returns path to PDF file and error
func RunCLI() (string, error){
	data, err := initMTP()
	if err != nil{
		log.Printf("Error getting data for CLI: %v\n",err)
		return "",fmt.Errorf("Error runing CLI: %v\n",err)
	}

	if !exist.IsExistFile(data.Input){
		log.Printf("File doesn't exist")
		return "",errors.New("File doesn't exist in directory")
	}

	content, err := os.ReadFile(data.Input)
	if err != nil{
		log.Printf("Error reading file: %v\n",err)
		return "",fmt.Errorf("Error running CLI: %v",err)
	}

	fullOutputPath := data.Output
	if data.SavePath != "" {
		fullOutputPath = filepath.Join(data.SavePath, data.Output)
	}
	var opts []mdtopdf.RenderOption

	if data.Ru {
		opts = append(opts, mdtopdf.WithUnicodeTranslator("cp1251"))
	}

	if data.IsNewPage {
		opts = append(opts, mdtopdf.IsHorizontalRuleNewPage(true))
	}

	theme := mdtopdf.LIGHT
	customTheme := ""
	switch strings.ToLower(data.Theme) {
	case "dark":
		theme = mdtopdf.DARK
	case "light", "default", "":
		theme = mdtopdf.LIGHT
	default:
		theme = mdtopdf.CUSTOM
		customTheme = data.Theme
	}
	
	params := mdtopdf.PdfRendererParams{
		Orientation:     data.Orientation,
		Papersz:         data.PageSize,
		PdfFile:         fullOutputPath,
		Opts:            opts,
		Theme:           theme,
		CustomThemeFile: customTheme,
	}
	pdfRender := mdtopdf.NewPdfRenderer(params)

	pdfRender.SetStyle(mdtopdf.Styler{
		FontFamily: data.Font,
		Size:       data.FontSize,
	})

	if data.Title != "" {
		pdfRender.Pdf.SetTitle(data.Title, true)
	}
	if data.Author != "" {
		pdfRender.Pdf.SetAuthor(data.Author, true)
	}

	// Footer
	if data.IsFooter {
		pdfRender.Pdf.SetFooterFunc(func() {
			pdfRender.Pdf.SetY(-15)
			pdfRender.Pdf.SetFont(data.Font, "I", data.FontSize)
			footer := ""
			if data.Author != "" {
				footer += data.Author + "  "
			}
			if data.Title != "" {
				footer += data.Title + "  "
			}
			footer += fmt.Sprintf("Page %d", pdfRender.Pdf.PageNo())
			pdfRender.Pdf.CellFormat(0, 10, footer, "", 0, "C", false, 0, "")
		})
	}

	if err := pdfRender.Process(content); err != nil {
		return "", fmt.Errorf("error converting to PDF: %w", err)
	}
	log.Printf("Successful convertation in PDF file - %s\n",fullOutputPath)
	return fullOutputPath,nil
}
