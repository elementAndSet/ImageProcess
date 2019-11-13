package main

import (
	"encoding/json"
	"fmt"
	img "myimage"
	mp "mypath"
	"os"
	fp "path/filepath"
	"strconv"
	"strings"
)

func main() {
	var commandsList [][]string
	help := func() {
		fmt.Println("    -f0 : create new folder and make result file in that")
		fmt.Println(" 1Ddiff [-f0]  : get 1D diff image")
		fmt.Println(" 2Ddiff [-f0]  : get 2D diff info image")
		fmt.Println(" 2Dsimple [-f0]  : get 2D diff info clean image ")
		fmt.Println(" info [-f0]  : get info CSV and JSON")
		//fmt.Println(" jpegnum : numbering jpeg format file")
		//fmt.Println(" pngnum : numbering png format file")
		fmt.Println(" num : numbering jpeg andpng format file")
		fmt.Println(" graph [-f0] : make Red-Green, Green-Blue, Red-Blue Graph")
	}
	if len(os.Args) < 2 {
		fmt.Println("put your command")
		help()
	} else {
		for mut := 1; mut < len(os.Args); mut++ {
			if strings.HasPrefix(os.Args[mut], "-") == false {
				var commands = []string{os.Args[mut]}
				commandsList = append(commandsList, commands)
			} else {
				commandsList[len(commandsList)-1] = append(commandsList[len(commandsList)-1], os.Args[mut])
			}
		}
		for _, v := range commandsList {
			fun := v[0]
			parameter := append([]string(nil), v[1:]...)
			switch fun {
			case "help":
				help()
			case "num":
				imgNum("jpg", "jpeg", "JPG", "png", "PNG")
			case "info":
				imgInfo(parameter)
			case "1Ddiff":
				imgDiff1D(parameter)
			case "2Ddiff":
				imgDiff2D(parameter)
			case "2Dsimple":
				diff2DSimple(parameter)
			case "graph":
				imgGraph(parameter)
			case "test":
				flags(parameter)
			default:
				fmt.Println(" wrong command!")
			}
		}
	}
}

func flags(op []string) (bool, bool) {
	var f0, f1 = false, false
	for _, v := range op {
		switch {
		case v == "-f0":
			f0 = true
		case v == "-f1":
			f1 = true
		default:
			fmt.Println("not correct")
		}
	}
	return f0, f1
}

func imgInfo(options []string) {
	f0, _ := flags(options)
	list := mp.GetFilePathNameSlice("jpg", "jpeg", "JPG", "png")
	var format string
	var imgInfoJSONSlice []img.ImgAver
	for _, v := range list {
		_, extsn := mp.GetFileNameAndExtsn(v)
		if extsn == ".jpg" {
			format = "jpeg"
		} else if extsn == ".png" {
			format = "png"
		}

		baseName := fp.Base(v)
		RGBAInfo := img.GetRGBAInfo(v, format)
		R, G, B, A := img.GetRGBAAverage(RGBAInfo)
		xdiff, ydiff := img.GetDiffInfo(RGBAInfo)
		XDIFFAVE, YDIFFAVE := img.GetDiffAverage(xdiff, ydiff)
		imgAveInfo := img.ImgAver{
			Name:  baseName,
			R:     R,
			G:     G,
			B:     B,
			A:     A,
			XDiff: XDIFFAVE,
			YDiff: YDIFFAVE,
		}
		imgInfoJSONSlice = append(imgInfoJSONSlice, imgAveInfo)
	}
	JSONIndent, _ := json.MarshalIndent(imgInfoJSONSlice, " ", "  ")
	JSONIndentString := string(JSONIndent)

	cp, _ := os.Getwd()
	var infoFilesFolder, infoJSONFilePath, infoCSVFilePath string
	if f0 {
		infoFilesFolder = mp.MkDir("InfoFiles")
		infoJSONFilePath = cp + "\\" + infoFilesFolder + "\\" + "infoJSON.json"
		infoCSVFilePath = cp + "\\" + infoFilesFolder + "\\" + "infoCSV.CSV"
	} else {
		infoJSONFilePath = "infoJSON.json"
		infoCSVFilePath = "infoCSV.CSV"
	}

	fJSON, fJerr := os.Create(infoJSONFilePath)
	if fJerr != nil {
		panic(fJerr)
	}
	defer fJSON.Close()
	fJSON.WriteString(JSONIndentString)

	fCSV, fCerr := os.Create(infoCSVFilePath)
	if fCerr != nil {
		panic(fCerr)
	}
	defer fCSV.Close()
	for _, v := range imgInfoJSONSlice {
		each := "Name:" + v.Name + ", R:" + strconv.Itoa(v.R) + ", "
		each = each + "G:" + strconv.Itoa(v.G) + ", B:" + strconv.Itoa(v.B) + ", "
		each = each + "XDIIF:" + fmt.Sprintf("%f", v.XDiff) + ", " + "YDIFF:" + fmt.Sprintf("%f", v.YDiff) + "\n"
		fCSV.WriteString(each)
	}
}

func imgDiff2D(option []string) {
	var format string
	f0, f1 := flags(option)
	path, _ := os.Getwd()
	var newPath string
	if f0 || f1 {
		newPath = path + "\\" + mp.MkDir("2Ddiff")
	}
	for _, v := range mp.GetFilePathNameSlice("jpg", "jpeg", "JPG", "png") {
		name, extsn := mp.GetFileNameAndExtsn(v)
		if extsn == ".jpg" || extsn == ".jpeg" || extsn == ".JPG" {
			format = "jpeg"
		} else if extsn == ".png" {
			format = "png"
		} else {
			format = extsn
			continue
		}
		xDiff, yDiff := img.ReturnDiffInfo(img.ReturnNRGBA(v, format))
		xDiffSq, yDiffSq, averSq := img.ReturnDiffSquareInfo(xDiff, yDiff)
		img.DrawFromNRGBA(newPath+"\\"+name+"_graph", format, img.ReturnIrregDiffImg(xDiffSq, yDiffSq, averSq))
	}
}

func diff2DSimple(option []string) {
	var format string
	path, _ := os.Getwd()
	f0, f1 := flags(option)
	var newPath string
	if f0 || f1 {
		newPath = path + "\\" + mp.MkDir("2DdiffSimplify")
	}
	for _, v := range mp.GetFilePathNameSlice("jpg", "jpeg", "JPG", "png") {
		name, extsn := mp.GetFileNameAndExtsn(v)
		if extsn == ".jpg" || extsn == ".jpeg" || extsn == ".JPG" {
			format = "jpeg"
		} else if extsn == ".png" {
			format = "png"
		} else {
			format = extsn
			continue
		}
		xDiff, yDiff := img.ReturnDiffInfo(img.ReturnNRGBA(v, format))
		xDiffSq, yDiffSq, averSq := img.ReturnDiffSquareInfo(xDiff, yDiff)
		xDiffSqSimple, yDiffSqSimple := img.ReturnSimplify(xDiffSq, averSq), img.ReturnSimplify(yDiffSq, averSq)
		img.DrawFromNRGBA(newPath+"\\"+name+"_graph", format, img.ReturnIrregDiffImg(xDiffSqSimple, yDiffSqSimple, averSq))
	}
}

func imgDiff1D(option []string) {
	f0, _ := flags(option)
	var format string
	Path, _ := os.Getwd()
	var newFolder, xDiffNewName, yDiffNewName string
	if f0 {
		newFolder = mp.MkDir("1Ddiff")
	}
	for _, v := range mp.GetFilePathNameSlice("jpg", "jpeg", "JPG", "png") {
		Name, extsn := mp.GetFileNameAndExtsn(v)
		if extsn == ".jpg" || extsn == ".jpeg" || extsn == ".JPG" {
			format = "jpeg"
		} else if extsn == ".png" {
			format = "png"
		} else {
			format = extsn
			continue
		}
		RGBAInfo := img.GetRGBAInfo(v, format)
		xDiff, yDiff := img.GetDiffInfo(RGBAInfo)
		xDiffAve, yDiffAve := img.GetDiffAverage(xDiff, yDiff)
		xIrregDiff := img.GetIrregDiff(xDiff, xDiffAve)
		yIrregDiff := img.GetIrregDiff(yDiff, yDiffAve)
		XIrregImg := img.MakeImgFromXIrreg(xIrregDiff)
		YIrregImg := img.MakeImgFromYIrreg(yIrregDiff)
		if f0 {
			xDiffNewName = Path + "\\" + newFolder + "\\" + Name + "_XDIFF" + "." + format
			yDiffNewName = Path + "\\" + newFolder + "\\" + Name + "_YDIFF" + "." + format
		} else {
			xDiffNewName = Name + "_XDIFF" + "." + format
			yDiffNewName = Name + "_YDIFF" + "." + format
		}
		fmt.Println(v, " ", Name, " ", format)
		fmt.Println(xDiffNewName, " ", yDiffNewName)
		if format == "jpeg" {
			img.DrawJPEGImg(xDiffNewName, XIrregImg)
			img.DrawJPEGImg(yDiffNewName, YIrregImg)

		} else if format == "png" {
			img.DrawPNGImg(xDiffNewName, XIrregImg)
			img.DrawPNGImg(yDiffNewName, YIrregImg)
		}
	}
}

func imgNum(extsn ...string) {
	var names []string
	for _, v := range extsn {
		names = append(names, mp.GetFilePathNameSlice(v)...)
	}
	mp.FileNumbering(names)
}

func imgGraph(option []string) {
	var newName, np string
	f0, f1 := flags(option)
	cp, _ := os.Getwd()
	if f0 || f1 {
		np = cp + "\\" + mp.MkDir("GraphImg")
	}
	for _, EXTSN := range []string{"png", "jpg"} {
		var EXTSN2 string
		if EXTSN == "jpg" || EXTSN == "jpeg" || EXTSN == "JPG" {
			EXTSN2 = "jpeg"
		} else if EXTSN == "png" {
			EXTSN2 = "png"
		}
		for _, v := range mp.GetFilePathNameSlice(EXTSN) {
			nameAndExtsn := fp.Base(v)
			fmt.Println("NameAndExtsn : ", nameAndExtsn)
			RG, GB, RB := img.MakeGraphBlueprint(nameAndExtsn, img.GetRGBAInfo(nameAndExtsn, EXTSN2))
			name, _ := mp.GetFileNameAndExtsn(nameAndExtsn)
			if f0 || f1 {
				newName = np + "\\" + name
			} else {
				newName = name
			}
			img.DrawFromNRGBA(newName+"_RGgraph", EXTSN2, RG)
			img.DrawFromNRGBA(newName+"_GBgraph", EXTSN2, GB)
			img.DrawFromNRGBA(newName+"_RBgraph", EXTSN2, RB)
		}
	}
}

/*
func fileCp(src, dst string) {

	in, iErr := os.Open(src)
	if iErr != nil {
		panic(iErr)
	}
	defer in.Close()
	out, oErr := os.Create(dst)
	if oErr != nil {
		panic(oErr)
	}
	defer out.Close()
	_, cErr := io.Copy(in, out)
	if cErr != nil {
		panic(cErr)
	}
}
*/
