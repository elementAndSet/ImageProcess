package main

import (
	"encoding/json"
	"fmt"
	"image"
	img "myimage"
	mp "mypath"
	"os"
	fp "path/filepath"
	"strconv"
	"strings"
)

type commonInfo struct {
	filename    string
	RGBAInfo    img.ImgInfo
	averR       int
	averG       int
	averB       int
	averA       int
	xDiffSq     [][]float64
	yDiffSq     [][]float64
	yDiffSqAver float64
	xDiffSqAver float64
	imageNRGBA  *image.NRGBA
	xDiff       [][]img.RGBA
	yDiff       [][]img.RGBA
}

//Flags :
type Flags struct {
	f0 bool //new folder
	f1 bool //new folder, origin file
}

func pseudomain() {

	var commandsList [][]string
	help := func() {
		fmt.Println("    -f0 : create new folder and make result file in that")
		fmt.Println(" 1Ddiff [-f0]  : get 1D diff image")
		fmt.Println(" 2Ddiff [-f0]  : get 2D diff info image")
		fmt.Println(" 2Dsimple [-f0]  : get 2D diff info clean image ")
		fmt.Println(" info [-f0]  : get info CSV and JSON")
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
	}
}

func main() {
	var commandsList [][]string
	help := func() {
		fmt.Println("    -f0 : create new folder and make result file in that")
		fmt.Println(" 1Ddiff [-f0]  : get 1D diff image")
		fmt.Println(" 2Ddiff [-f0]  : get 2D diff info image")
		fmt.Println(" 2Dsimple [-f0]  : get 2D diff info clean image ")
		fmt.Println(" info [-f0]  : get info CSV and JSON")
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
				returnFlags(parameter)
			default:
				fmt.Println(" wrong command!")
			}
		}
	}
}

func returnFlags(op []string) Flags {
	var Flag Flags
	for _, v := range op {
		switch {
		case v == "-f0":
			Flag.f0 = true
		case v == "-f1":
			Flag.f1 = true
		default:
			fmt.Println("not correct")
		}
	}
	return Flag
}

func imgInfo(options []string) []commonInfo {
	flag := returnFlags(options)
	list := mp.GetFilePathNameSlice("jpg", "jpeg", "JPG", "png")
	infoList := []commonInfo{}
	var format string
	var imgInfoJSONSlice []img.ImgAver
	for _, v := range list {
		_, extsn := mp.GetFileNameAndExtsn(v)
		if extsn == ".jpg" || extsn == "jpge" {
			format = "jpeg"
		} else if extsn == ".png" {
			format = "png"
		}

		var info commonInfo
		info.filename = fp.Base(v)
		info.imageNRGBA = img.ReturnNRGBA(v, format)
		info.averR, info.averG, info.averB, info.averA = img.ReturnRGBAAverage(info.imageNRGBA)
		info.xDiff, info.yDiff = img.ReturnDiffInfo(info.imageNRGBA)
		_, _, info.xDiffSqAver, info.yDiffSqAver = img.ReturnDiffSquareInfo(info.xDiff, info.yDiff)
		imgAveInfo := img.ImgAver{
			Name:  info.filename,
			R:     info.averR,
			G:     info.averG,
			B:     info.averB,
			A:     info.averA,
			XDiff: info.xDiffSqAver,
			YDiff: info.yDiffSqAver,
		}
		imgInfoJSONSlice = append(imgInfoJSONSlice, imgAveInfo)
		infoList = append(infoList, info)
		fmt.Println(info.filename, " processed")
	}
	JSONIndent, _ := json.MarshalIndent(imgInfoJSONSlice, " ", "  ")
	JSONIndentString := string(JSONIndent)

	cp, _ := os.Getwd()
	var infoFilesFolder, infoJSONFilePath, infoCSVFilePath string
	if flag.f0 {
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
	return infoList
}

func imgDiff1D(option []string) {
	flag := returnFlags(option)
	var format string
	Path, _ := os.Getwd()
	var newFolder, xDiffNewName, yDiffNewName string
	if flag.f0 {
		newFolder = mp.MkDir("1Ddiff")
	}
	for _, v := range mp.GetFilePathNameSlice("jpg", "jpeg", "JPG", "png") {
		var info commonInfo
		Name, extsn := mp.GetFileNameAndExtsn(v)
		if extsn == ".jpg" || extsn == ".jpeg" || extsn == ".JPG" {
			format = "jpeg"
		} else if extsn == ".png" {
			format = "png"
		} else {
			format = extsn
			continue
		}

		info.imageNRGBA = img.ReturnNRGBA(v, format)
		info.xDiff, info.yDiff = img.ReturnDiffInfo(info.imageNRGBA)
		info.xDiffSq, info.yDiffSq, info.xDiffSqAver, info.yDiffSqAver = img.ReturnDiffSquareInfo(info.xDiff, info.yDiff)

		if flag.f0 {
			xDiffNewName = Path + "\\" + newFolder + "\\" + Name + "_XDIFF" + "."
			yDiffNewName = Path + "\\" + newFolder + "\\" + Name + "_YDIFF" + "."
		} else {
			xDiffNewName = Name + "_XDIFF" + "."
			yDiffNewName = Name + "_YDIFF" + "."
		}
		img.DrawFromNRGBA(xDiffNewName, format, img.ReturnIrregDiffImg(info.xDiffSq, [][]float64{}, info.xDiffSqAver))
		img.DrawFromNRGBA(yDiffNewName, format, img.ReturnIrregDiffImg([][]float64{}, info.yDiffSq, info.xDiffSqAver))

		fmt.Println(v, " ", Name, " ", format)
		fmt.Println(xDiffNewName, " ", yDiffNewName)
	}
}

func imgDiff2D(option []string) {
	var format string
	flag := returnFlags(option)
	path, _ := os.Getwd()
	var newPath string
	if flag.f0 || flag.f1 {
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
		xDiffSq, yDiffSq, xAverSq, yAverSq := img.ReturnDiffSquareInfo(xDiff, yDiff)
		img.DrawFromNRGBA(newPath+"\\"+name+"_graph", format, img.ReturnIrregDiffImg(xDiffSq, yDiffSq, (xAverSq+yAverSq/2.0)))
	}
}

func diff2DSimple(option []string) {
	var format string
	path, _ := os.Getwd()
	flag := returnFlags(option)
	var newPath string
	if flag.f0 || flag.f1 {
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
		xDiffSq, yDiffSq, xAverSq, yAverSq := img.ReturnDiffSquareInfo(xDiff, yDiff)
		averSq := (xAverSq + yAverSq) / 2.0
		xDiffSqSimple, yDiffSqSimple := img.ReturnSimplify(xDiffSq, averSq), img.ReturnSimplify(yDiffSq, averSq)
		img.DrawFromNRGBA(newPath+"\\"+name+"_graph", format, img.ReturnIrregDiffImg(xDiffSqSimple, yDiffSqSimple, averSq))
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
	flag := returnFlags(option)
	cp, _ := os.Getwd()
	if flag.f0 || flag.f1 {
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
			if flag.f0 || flag.f1 {
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
