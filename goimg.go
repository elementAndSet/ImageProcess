package main

import (
	"encoding/json"
	"fmt"
	"io"
	img "myimage"
	mp "mypath"
	"os"
	fp "path/filepath"
	"strconv"
)

func main() {
	argsLen := len(os.Args)
	if argsLen == 1 {
		fmt.Println("put your command")
		fmt.Println(" help : show commands")
		fmt.Println(" diff [-f0]  : get diff image")
		fmt.Println(" info [-fo]  : get info CSV and JSON")
		fmt.Println(" jpegnum : numbering jpeg format file")
		fmt.Println(" pngnum : numbering png format file")
		fmt.Println(" num : numbering jpeg andpng format file")
		fmt.Println(" graph [-fo] : make Red-Green, Green-Blue, Red-Blue Graph")

	} else if argsLen == 2 {
		switch os.Args[1] {
		case "help":
			fmt.Println("    -f0 : create new folder and make result file in that")
			fmt.Println(" diff [-f0]  : get diff image")
			fmt.Println(" info [-fo]  : get info CSV and JSON")
			fmt.Println(" jpegnum : numbering jpeg format file")
			fmt.Println(" pngnum : numbering png format file")
			fmt.Println(" num : numbering jpeg andpng format file")
			fmt.Println(" graph [-fo] : make Red-Green, Green-Blue, Red-Blue Graph")
		case "diff":
			imgDiffPro("Current")
		case "info":
			imgInfo("Current")
		case "jpegnum":
			imgNum("jpg", "jpeg", "JPG")
		case "pngnum":
			imgNum("png")
		case "num":
			imgNum("jpg", "jpeg", "JPG", "png")
		case "graph":
			imgGraph("Current")
		default:
			fmt.Println(" wrong command!")
		}
	} else if argsLen == 3 && os.Args[2] == "-f0" {
		switch os.Args[1] {
		case "diff":
			imgDiffPro("-f0")
		case "info":
			imgInfo("-f0")
		case "graph":
			imgGraph("-f0")
		default:
			fmt.Println(" wrong command!")
		}
		fmt.Println(" wrong command!")
	} else if argsLen == 3 && os.Args[2] == "-f1" {
		switch os.Args[1] {
		case "diff":
			imgDiffPro("-f1")
		case "graph":
			imgGraph("-f1")
		default:
			fmt.Println(" wrong command!")
		}
	} else {
		fmt.Println(" wrong command!")
	}
}

func imgInfo(option string) {
	list := mp.GetFilePathNameSlice("jpg", "png")
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
	if option == "-f0" {
		infoFilesFolder = mp.MkDir("InfoFiles")
		infoJSONFilePath = cp + "\\" + infoFilesFolder + "\\" + "infoJSON.json"
		infoCSVFilePath = cp + "\\" + infoFilesFolder + "\\" + "infoCSV.CSV"
	} else if option == "Current" {
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

func imgDiffPro(option string) {
	var format string
	Path, _ := os.Getwd()
	var newFolder, xDiffNewName, yDiffNewName string
	if option == "-f0" || option == "-f1" {
		newFolder = mp.MkDir("NewImg")
	}
	for _, v := range mp.GetFilePathNameSlice("jpg", "png") {
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
		if option == "-f0" || option == "-f1" {
			xDiffNewName = Path + "\\" + newFolder + "\\" + Name + "_XDIFF" + "." + format
			yDiffNewName = Path + "\\" + newFolder + "\\" + Name + "_YDIFF" + "." + format
		} else if option == "Current" {
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
		if option == "-f1" {
			fileCp(v, Path+"\\"+newFolder+Name+"."+format)
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

func imgGraph(option string) {
	var newName, np string
	cp, _ := os.Getwd()
	if option == "-f0" || option == "-f1" {
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
			if option == "-f0" || option == "-f1" {
				newName = np + "\\" + name
			} else if option == "Current" {
				newName = name
			}
			img.DrawFromNRGBA(newName+"_RGgraph", EXTSN2, RG)
			img.DrawFromNRGBA(newName+"_GBgraph", EXTSN2, GB)
			img.DrawFromNRGBA(newName+"_RBgraph", EXTSN2, RB)
			if option == "-f1" {
				newName = np + "\\" + nameAndExtsn
				fileCp(nameAndExtsn, newName)
			}
		}
	}
}

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
