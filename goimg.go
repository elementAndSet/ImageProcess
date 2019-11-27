package main

import (
	"encoding/json"
	"fmt"
	"image"
	img "myimage"
	mp "mypath"
	"os"
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

type fEnum struct {
	f1 func(commonInfo, string, string, string, string) (commonInfo, img.ImgAver)
	f2 func(commonInfo, string, string, string, string) commonInfo
}

//Flags :
type Flags struct {
	f0 bool //new folder
	f1 bool //new folder, origin file
}

//with pipe
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
			case "info":
				objects = imgInfo(parameter, objects)
			case "1Ddiff":
				objects = imgDiff1D(parameter, objects)
			case "2Ddiff":
				objects = imgDiff2D(parameter, objects)
			case "2Dsimple":
				objects = imgDiff2DSimple(parameter, objects)
			case "graph":
				objects = imgGraph(parameter, objects)
			case "j2p":
				objects = imgJ2P(parameter, objects)
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

func imgInfo(option []string, objects []commonInfo) []commonInfo {
	f := func(object commonInfo, nameNoPointExtsn string, extsnWithPoint string, deleExtsn string, noUse string) (commonInfo, img.ImgAver) {
		var info commonInfo
		info = extractInfo(object, nameNoPointExtsn+extsnWithPoint, deleExtsn)
		if object.averR == 0 && object.averG == 0 && object.averB == 0 && object.averA == 0 {
			info.averR, info.averG, info.averB, info.averA = img.ReturnRGBAAverage(info.imageNRGBA)
			object.averR, object.averG, object.averB, object.averA = info.averR, info.averG, info.averB, info.averA
		} else {
			info.averR, info.averG, info.averB, info.averA = object.averR, object.averG, object.averB, object.averA
		}
		imgAveInfo := img.ImgAver{
			Name:  info.filename,
			R:     info.averR,
			G:     info.averG,
			B:     info.averB,
			A:     info.averA,
			XDiff: info.xDiffSqAver,
			YDiff: info.yDiffSqAver,
		}
		return info, imgAveInfo
	}
	fn := fEnum{f, nil}
	return commonIter(option, objects, fn, "imgInfo")
}

func imgDiff1D(option []string, objects []commonInfo) []commonInfo {
	f := func(object commonInfo, nameNoPointExtsn string, extsnWithPoint string, deleExtsn string, newNameNoPointExtsn string) commonInfo {
		var info commonInfo
		info = extractInfo(object, nameNoPointExtsn+extsnWithPoint, deleExtsn)
		xIrregDiffImg := img.ReturnIrregDiffImg(info.xDiffSq, [][]float64{}, (info.xDiffSqAver+info.yDiffSqAver)/2.0)
		yIrregDiffImg := img.ReturnIrregDiffImg([][]float64{}, info.yDiffSq, (info.xDiffSqAver+info.yDiffSqAver)/2.0)
		img.DrawFromNRGBA(newNameNoPointExtsn+"_XDIFF", deleExtsn, xIrregDiffImg)
		img.DrawFromNRGBA(newNameNoPointExtsn+"_YDIFF", deleExtsn, yIrregDiffImg)
		return info
	}
	fn := fEnum{nil, f}
	return commonIter(option, objects, fn, "1Ddiff")
}

func imgDiff2D(option []string, objects []commonInfo) []commonInfo {
	f := func(object commonInfo, nameNoPointExtsn string, extsnWithPoint string, deleExtsn string, newNameNoPointExtsn string) commonInfo {
		var info commonInfo
		info = extractInfo(object, nameNoPointExtsn+extsnWithPoint, deleExtsn)
		IrregDiffImg := img.ReturnIrregDiffImg(info.xDiffSq, info.yDiffSq, (info.xDiffSqAver+info.yDiffSqAver)/2.0)
		img.DrawFromNRGBA(newNameNoPointExtsn+"_2Ddiff", deleExtsn, IrregDiffImg)
		return info
	}
	fn := fEnum{nil, f}
	return commonIter(option, objects, fn, "2Ddiff")
}

func imgDiff2DSimple(option []string, objects []commonInfo) []commonInfo {
	f := func(object commonInfo, nameNoPointExtsn string, extsnWithPoint string, deleExtsn string, newNameNoPointExtsn string) commonInfo {
		var info commonInfo
		info = extractInfo(object, nameNoPointExtsn+extsnWithPoint, deleExtsn)
		averSq := (info.xDiffSqAver + info.yDiffSqAver) / 2.0
		xDiffSqSimple, yDiffSqSimple := img.ReturnSimplify(info.xDiffSq, averSq), img.ReturnSimplify(info.yDiffSq, averSq)
		img.DrawFromNRGBA(newNameNoPointExtsn+"_2Dsimple", deleExtsn, img.ReturnIrregDiffImg(xDiffSqSimple, yDiffSqSimple, averSq))
		return info
	}
	fn := fEnum{nil, f}
	return commonIter(option, objects, fn, "2DdiffSimple")
}

func imgNum(option []string, object []commonInfo) {

}

/*
func imgNum(extsn ...string) {
	var names []string
	for _, v := range extsn {
		names = append(names, mp.GetFilePathNameSlice(v)...)
	}
	mp.FileNumbering(names)
}
*/

func imgGraph(option []string, objects []commonInfo) []commonInfo {
	f := func(object commonInfo, nameNoPointExtsn string, extsnWithPoint string, deleExtsn string, newNameNoPointExtsn string) commonInfo {
		var info commonInfo
		RGBAInfo := img.GetRGBAInfo(nameNoPointExtsn+extsnWithPoint, deleExtsn)
		RG, GB, RB := img.MakeGraphBlueprint("", RGBAInfo)
		img.DrawFromNRGBA(newNameNoPointExtsn+"_RGgraph", deleExtsn, RG)
		img.DrawFromNRGBA(newNameNoPointExtsn+"_GBgraph", deleExtsn, GB)
		img.DrawFromNRGBA(newNameNoPointExtsn+"_RBgraph", deleExtsn, RB)
		return info
	}
	fn := fEnum{nil, f}
	return commonIter(option, objects, fn, "ImgGraph")
}

func imgJ2P(option []string, objects []commonInfo) []commonInfo {
	f := func(object commonInfo, nameNoPointExtsn string, extsnWithPoint string, deleExtsn string, newNameNoPointExtsn string) commonInfo {
		var info commonInfo
		var newNRGBA *image.NRGBA
		if deleExtsn == "png" {
			newNRGBA = img.ReturnNRGBA(nameNoPointExtsn+extsnWithPoint, deleExtsn)
		} else if deleExtsn == "jpeg" {
			newNRGBA = img.ReturnNRGBA(nameNoPointExtsn+extsnWithPoint, deleExtsn)
		}
		img.DrawFromNRGBA(newNameNoPointExtsn, "png", newNRGBA)
		info.imageNRGBA = newNRGBA
		return info
	}
	fn := fEnum{nil, f}
	return commonIter(option, objects, fn, "J2Pfolder")
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

func extractInfo(object commonInfo, nameWithPointExtsn string, deleExtsn string) commonInfo {
	var info commonInfo
	info.filename = nameWithPointExtsn
	info.imageNRGBA = img.ReturnNRGBA(info.filename, deleExtsn)
	if len(object.xDiff) < 1 && len(object.yDiff) < 1 {
		info.xDiff, info.yDiff = img.ReturnDiffInfo(info.imageNRGBA)
	} else {
		info.xDiff, info.yDiff = object.xDiff, object.yDiff
	}
	if object.xDiffSqAver < 0.1 && object.yDiffSqAver < 0.1 {
		info.xDiffSq, info.yDiffSq, info.xDiffSqAver, info.yDiffSqAver = img.ReturnDiffSquareInfo(info.xDiff, info.yDiff)
	} else {
		info.xDiffSq, info.yDiffSq = object.xDiffSq, object.yDiffSq
		info.xDiffSqAver, info.yDiffSqAver = object.xDiffSqAver, object.yDiffSqAver
	}
	return info
}

func commonIter(option []string, objects []commonInfo, fn fEnum, folderName string) []commonInfo {
	var newPath, newNameNoPointExtsn string
	flag := returnFlags(option)
	cp, _ := os.Getwd()
	if flag.f0 || flag.f1 {
		newPath = cp + "\\" + mp.MkDir(folderName) + "\\"
	}
	var infoSlice []img.ImgAver
	for _, v := range objects {
		nameNoPointExtsn, extsnWithPoint := mp.GetFileNameAndExtsn(v.filename)

		if flag.f0 || flag.f1 {
			newNameNoPointExtsn = newPath + nameNoPointExtsn
		} else {
			newNameNoPointExtsn = nameNoPointExtsn
		}
		var deleExtsn string
		if extsnWithPoint == ".jpg" || extsnWithPoint == ".jpeg" || extsnWithPoint == ".JPG" {
			deleExtsn = "jpeg"
		} else if extsnWithPoint == ".png" || extsnWithPoint == ".PNG" {
			deleExtsn = "png"
		}
		if fn.f1 != nil {
			_, element := fn.f1(v, nameNoPointExtsn, extsnWithPoint, deleExtsn, newNameNoPointExtsn)
			infoSlice = append(infoSlice, element)
		} else {
			fn.f2(v, nameNoPointExtsn, extsnWithPoint, deleExtsn, newNameNoPointExtsn)
		}
		fmt.Println(" " + nameNoPointExtsn + " executed")
	}
	if fn.f1 != nil {
		fCSV, fCerr := os.Create(newPath + "infoCSV.CSV")
		if fCerr != nil {
			panic(fCerr)
		}
		defer fCSV.Close()
		fJSON, fJerr := os.Create(newPath + "infoJSON.json")
		if fJerr != nil {
			panic(fJerr)
		}
		defer fJSON.Close()

		JSONIndent, _ := json.MarshalIndent(infoSlice, " ", "  ")
		JSONIndentString := string(JSONIndent)
		fJSON.WriteString(JSONIndentString)
		for _, v := range infoSlice {
			each := "Name:" + v.Name + ", R:" + strconv.Itoa(v.R) + ", "
			each = each + "G:" + strconv.Itoa(v.G) + ", B:" + strconv.Itoa(v.B) + ", "
			each = each + "XDIIF:" + fmt.Sprintf("%f", v.XDiff) + ", " + "YDIFF:" + fmt.Sprintf("%f", v.YDiff) + "\n"
			fCSV.WriteString(each)
		}
	}
	return objects
}
