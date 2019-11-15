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
	imageNRGBA  *image.NRGBA
	RGBAInfo    img.ImgInfo
	averR       int
	averG       int
	averB       int
	averA       int
	xDiff       [][]img.RGBA
	yDiff       [][]img.RGBA
	xDiffSq     [][]float64
	yDiffSq     [][]float64
	yDiffSqAver float64
	xDiffSqAver float64
}

//Flags :
type Flags struct {
	f0 bool //new folder
	f1 bool //new folder, origin file
}

func main() {

	var commandsList [][]string
	var objects []commonInfo
	nameExtsn := mp.GetFilePathNameSlice("jpg", "jpeg", "JPG", "png", "PNG")
	for _, v := range nameExtsn {
		var obj commonInfo
		obj.filename = v
		objects = append(objects, obj)
	}
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
			case "1Ddiff":
				objects = imgDiff1D(parameter, objects)
			case "2Ddiff":
				objects = imgDiff2D(parameter, objects)
			case "2Dsimple":
				objects = imgDiff2DSimple(parameter, objects)
			case "graph":
				objects = imgGraph(parameter, objects)
			case "test":
				//n, e := mp.GetFileNameAndExtsn(objects[2].filename)
				//fmt.Println(n, " -- ", e)
				var info commonInfo
				fmt.Println(info)
			default:
				fmt.Println(" wrong command!")
			}
		}
	}
}

/*
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
*/

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

func imgDiff1D(option []string, objects []commonInfo) []commonInfo {
	flag := returnFlags(option)
	var format string
	Path, _ := os.Getwd()
	var newFolder, xDiffNewName, yDiffNewName string
	if flag.f0 {
		newFolder = mp.MkDir("1Ddiff")
	}
	for _, v := range objects {
		Name, extsn := mp.GetFileNameAndExtsn(v.filename)
		if extsn == ".jpg" || extsn == ".jpeg" || extsn == ".JPG" {
			format = "jpeg"
		} else if extsn == ".png" || extsn == ".PNG" {
			format = "png"
		} else {
			format = extsn
			continue
		}

		v.imageNRGBA = img.ReturnNRGBA(v.filename, format)
		v.xDiff, v.yDiff = img.ReturnDiffInfo(v.imageNRGBA)
		v.xDiffSq, v.yDiffSq, v.xDiffSqAver, v.yDiffSqAver = img.ReturnDiffSquareInfo(v.xDiff, v.yDiff)
		if flag.f0 {
			xDiffNewName = Path + "\\" + newFolder + "\\" + Name + "_XDIFF" + "."
			yDiffNewName = Path + "\\" + newFolder + "\\" + Name + "_YDIFF" + "."
		} else {
			xDiffNewName = Name + "_XDIFF" + "."
			yDiffNewName = Name + "_YDIFF" + "."
		}
		img.DrawFromNRGBA(xDiffNewName, format, img.ReturnIrregDiffImg(v.xDiffSq, [][]float64{}, v.xDiffSqAver))
		img.DrawFromNRGBA(yDiffNewName, format, img.ReturnIrregDiffImg([][]float64{}, v.yDiffSq, v.xDiffSqAver))

		fmt.Println(v.filename, " ", Name, " ", format)
		fmt.Println(xDiffNewName, " ", yDiffNewName)
	}
	return objects
}

func imgDiff2D(option []string, objects []commonInfo) []commonInfo {
	var format string
	flag := returnFlags(option)
	path, _ := os.Getwd()
	var newPath string
	if flag.f0 || flag.f1 {
		newPath = path + "\\" + mp.MkDir("2Ddiff") + "\\"
	}
	for _, v := range objects {
		name, extsn := mp.GetFileNameAndExtsn(v.filename)
		if extsn == ".jpg" || extsn == ".jpeg" || extsn == ".JPG" {
			format = "jpeg"
		} else if extsn == ".png" {
			format = "png"
		} else {
			format = extsn
			continue
		}
		if len(v.xDiff) == 0 || len(v.yDiff) == 0 {
			fmt.Println(" recycle diff")
			v.xDiff, v.yDiff = img.ReturnDiffInfo(img.ReturnNRGBA(v.filename, format))
		}
		if len(v.xDiffSq) == 0 || len(v.yDiffSq) == 0 || v.xDiffSqAver == 0.0 || v.yDiffSqAver == 0.0 {
			fmt.Println(" RECYCLE DIFF")
			v.xDiffSq, v.yDiffSq, v.xDiffSqAver, v.yDiffSqAver = img.ReturnDiffSquareInfo(v.xDiff, v.yDiff)
		}
		//v.xDiff, v.yDiff = img.ReturnDiffInfo(img.ReturnNRGBA(v.filename, format))
		//v.xDiffSq, v.yDiffSq, v.xDiffSqAver, v.yDiffSqAver = img.ReturnDiffSquareInfo(v.xDiff, v.yDiff)
		img.DrawFromNRGBA(newPath+name+"_graph", format, img.ReturnIrregDiffImg(v.xDiffSq, v.yDiffSq, (v.xDiffSqAver+v.yDiffSqAver/2.0)))
		fmt.Println(name + "_graph." + format)
	}
	return objects
}

func imgDiff2DSimple(option []string, objects []commonInfo) []commonInfo {
	var format string
	path, _ := os.Getwd()
	flag := returnFlags(option)
	var newPath string
	if flag.f0 || flag.f1 {
		newPath = path + "\\" + mp.MkDir("2DdiffSimplify") + "\\"
	}
	for _, v := range objects {
		name, extsn := mp.GetFileNameAndExtsn(v.filename)
		if extsn == ".jpg" || extsn == ".jpeg" || extsn == ".JPG" {
			format = "jpeg"
		} else if extsn == ".png" {
			format = "png"
		} else {
			format = extsn
			continue
		}
		v.xDiff, v.yDiff = img.ReturnDiffInfo(img.ReturnNRGBA(v.filename, format))
		v.xDiffSq, v.yDiffSq, v.xDiffSqAver, v.yDiffSqAver = img.ReturnDiffSquareInfo(v.xDiff, v.yDiff)
		averSq := (v.xDiffSqAver + v.yDiffSqAver) / 2.0
		xDiffSqSimple, yDiffSqSimple := img.ReturnSimplify(v.xDiffSq, averSq), img.ReturnSimplify(v.yDiffSq, averSq)
		img.DrawFromNRGBA(newPath+name+"_graph", format, img.ReturnIrregDiffImg(xDiffSqSimple, yDiffSqSimple, averSq))
		fmt.Println(name + "_graph." + format)
	}
	return objects
}

func imgNum(extsn ...string) {
	var names []string
	for _, v := range extsn {
		names = append(names, mp.GetFilePathNameSlice(v)...)
	}
	mp.FileNumbering(names)
}

func imgGraph(option []string, objects []commonInfo) []commonInfo {
	var newName, np string
	flag := returnFlags(option)
	cp, _ := os.Getwd()
	if flag.f0 || flag.f1 {
		np = cp + "\\" + mp.MkDir("GraphImg")
	}
	for _, v := range objects {
		var EXTSN string
		name, extsn := mp.GetFileNameAndExtsn(v.filename)
		if extsn == ".jpg" || extsn == ".jpeg" || extsn == ".JPG" {
			EXTSN = "jpeg"
		} else if extsn == ".png" {
			EXTSN = "png"
		}
		nameAndExtsn := fp.Base(v.filename)
		//fmt.Println("NameAndExtsn : ", nameAndExtsn)
		RG, GB, RB := img.MakeGraphBlueprint(nameAndExtsn, img.GetRGBAInfo(nameAndExtsn, EXTSN))
		//name, _ := mp.GetFileNameAndExtsn(nameAndExtsn)
		if flag.f0 || flag.f1 {
			newName = np + "\\" + name
		} else {
			newName = name
		}
		img.DrawFromNRGBA(newName+"_RGgraph", EXTSN, RG)
		img.DrawFromNRGBA(newName+"_GBgraph", EXTSN, GB)
		img.DrawFromNRGBA(newName+"_RBgraph", EXTSN, RB)
	}
	return objects
}
