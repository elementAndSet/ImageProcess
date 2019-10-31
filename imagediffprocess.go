package main

import (
	"fmt"
	img "myimage"
	mp "mypath"
	"os"
)

func main() {
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
