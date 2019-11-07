package main

import (
	"encoding/json"
	"fmt"
	img "myimage"
	mp "mypath"
	"os"
	fp "path/filepath"
	"strconv"
)

func main() {
	argsLen := len(os.Args)
	if argsLen == 1 {
		fmt.Println("common")
	} else if argsLen == 2 && os.Args[1] == "help" {
		fmt.Println(" diff : get diff image")
		fmt.Println(" info : get info CSV and JSON")
	} else if argsLen == 2 && os.Args[1] == "diff" {
		imgDiffPro()
	} else if argsLen == 2 && os.Args[1] == "info" {
		imgInfoFolder()
	} else {
		fmt.Println(" wrong command!")
	}
}

func imgInfoFolder() {
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
		imgAveInfo := img.ImgAver{
			Name: baseName,
			R:    R,
			G:    G,
			B:    B,
			A:    A,
		}
		imgInfoJSONSlice = append(imgInfoJSONSlice, imgAveInfo)
	}
	JSONIndent, _ := json.MarshalIndent(imgInfoJSONSlice, " ", "  ")
	JSONIndentString := string(JSONIndent)
	//fmt.Println(JSONIndentString)
	cp, _ := os.Getwd()

	InfoFilesFolder := mp.MkDir("InfoFiles")
	infoJSONFilePath := cp + "\\" + InfoFilesFolder + "\\" + "infoJSON.json"
	fJSON, fJerr := os.Create(infoJSONFilePath)
	if fJerr != nil {
		panic(fJerr)
	}
	defer fJSON.Close()
	fJSON.WriteString(JSONIndentString)

	infoCSVFilePath := cp + "\\" + InfoFilesFolder + "\\" + "infoCSV.CSV"
	fCSV, fCerr := os.Create(infoCSVFilePath)
	if fCerr != nil {
		panic(fCerr)
	}
	defer fCSV.Close()
	for _, v := range imgInfoJSONSlice {
		each := "Name : " + v.Name + ", R : " + strconv.Itoa(v.R) + ", "
		each = each + "G : " + strconv.Itoa(v.G) + ", B : " + strconv.Itoa(v.B) + "\n"
		fCSV.WriteString(each)
	}
}

func imgDiffPro() {
	var format string
	Path, _ := os.Getwd()
	newFolder := mp.MkDir("NewImg")
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
		xDiffNewName := Path + "\\" + newFolder + "\\" + Name + "_XDIFF" + "." + format
		yDiffNewName := Path + "\\" + newFolder + "\\" + Name + "_YDIFF" + "." + format
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
