package main

import (
	"fmt"
	img "mysrc/image"
)

func makeImgInfoFromIrreg(xIrreg, yIrreg [][]float64) ([][]img.ImgInfo, [][]img.ImgInfo) {
	width, height := len(yIrreg), len(xIrreg)
	fmt.Println("width : ", width, " heith : ", height)
	var hLine, vLine []img.ImgInfo
	var hDiff, vDiff [][]img.ImgInfo
	for i := 0; i < width; i++ {
		baseInfo := img.ImgInfo{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
		hLine = append(hLine, baseInfo)
		vLine = append(vLine, baseInfo)
	}
	for j := 0; j < height; j++ {
		var hStack, vStack []img.ImgInfo
		hStack = append(hStack, hLine...)
		vStack = append(vStack, vLine...)
		hDiff = append(hDiff, hStack)
		vDiff = append(vDiff, vStack)
	}

	for Xy := 0; Xy < height; Xy++ {
		for Xx := 0; Xx < width-1; Xx++ {
			if xIrreg[Xy][Xx] > 0.9 {
				hDiff[Xy][Xx].R = 255
				hDiff[Xy][Xx].A = 255
				//fmt.Println(" x y : ", Xx, " ", Xy)
				//vStack[Xy][Xx+1].R = 255
			}
		}
	}

	for Yx := 0; Yx < width; Yx++ {
		for Yy := 0; Yy < height-1; Yy++ {
			if yIrreg[Yx][Yy] > 0.9 {
				vDiff[Yy][Yx].B = 255
				vDiff[Yy][Yx].A = 255
			}
		}
	}

	return hDiff, vDiff
}

func main() {
	fmt.Println("image Name")
	names := img.GetImgFileNames()
	for _, v := range names {
		format := img.GetImgFileFormat(v)
		//r, g, b := img.GetImgTint(v, format)
		//fmt.Println(v, " ", r, " ", g, " ", b)
		imgInfo := img.GetRGBAInfo(v, format)
		r, g, b, a := img.GetRGBAAverage(imgInfo)
		fmt.Println(v, " ", r, " ", g, " ", b, " ", a)
		xDiff, yDiff := img.GetDiffInfo(imgInfo)
		xDAv, yDAv := img.GetDiffAverage(xDiff, yDiff)
		fmt.Println("xDAv : ", xDAv, " yDAv : ", yDAv)
		xIrregDiff, yIrregDiff := img.GetIrregDiff(xDiff, xDAv), img.GetIrregDiff(yDiff, yDAv)
		//for _, v := range xIrregDiff {
		//	for k, v := range v {
		//		fmt.Println(" xIrregDiff : ", k, " ", v)
		//	}
		//}
		hoNewImg, veNewImg := makeImgInfoFromIrreg(xIrregDiff, yIrregDiff)
		redFN, greenFN := "NewHo"+v, "NewVe"+v
		img.DrawImg(redFN, hoNewImg)
		img.DrawImg(greenFN, veNewImg)

	}
}
