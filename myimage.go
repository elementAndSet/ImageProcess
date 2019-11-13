package myimage

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	fp "path/filepath"
)

//ImgAver :
type ImgAver struct {
	Name  string
	R     int
	G     int
	B     int
	A     int
	XDiff float64
	YDiff float64
}

//ImgInfo : RGBA
type ImgInfo struct {
	R int
	G int
	B int
	A int
}

//RGBA :
type RGBA struct {
	R int
	G int
	B int
	A int
}

//DIFFTHRESH  :
const DIFFTHRESH float64 = 0.7

// GetImgFileNames : Get Slice
func GetImgFileNames() []string {
	base, _ := os.Getwd()
	jpgTag := "\\*.jpg"
	jpegTag := "\\*.jpeg"
	JPGTag := "\\*.JPG"
	pngTag := "\\*.png"
	jpgFileFilter := base + jpgTag
	jpegFileFilter := base + jpegTag
	JPGFileFilter := base + JPGTag
	pngFileFilter := base + pngTag
	jpgFileNames, _ := fp.Glob(jpgFileFilter)
	jpegFileNames, _ := fp.Glob(jpegFileFilter)
	JPGFileNames, _ := fp.Glob(JPGFileFilter)
	pngFileNames, _ := fp.Glob(pngFileFilter)
	var filePaths []string
	filePaths = append(append(append(filePaths, jpgFileNames...), jpegFileNames...), JPGFileNames...)
	filePaths = append(filePaths, pngFileNames...)
	var fileNames []string
	for _, v := range filePaths {
		fileNames = append(fileNames, fp.Base(v))
	}
	return fileNames
}

// GetImgFilePaths : Get Slice
func GetImgFilePaths() []string {
	base, _ := os.Getwd()
	jpgTag := "\\*.jpg"
	jpegTag := "\\*.jpeg"
	JPGTag := "\\*.JPG"
	pngTag := "\\*.png"
	jpgFileFilter := base + jpgTag
	jpegFileFilter := base + jpegTag
	JPGFileFilter := base + JPGTag
	pngFileFilter := base + pngTag
	jpgFileNames, _ := fp.Glob(jpgFileFilter)
	jpegFileNames, _ := fp.Glob(jpegFileFilter)
	JPGFileNames, _ := fp.Glob(JPGFileFilter)
	pngFileNames, _ := fp.Glob(pngFileFilter)
	var filePaths []string
	filePaths = append(append(append(filePaths, jpgFileNames...), jpegFileNames...), JPGFileNames...)
	filePaths = append(filePaths, pngFileNames...)
	return filePaths
}

//GetImgTint : get Average RGBA from Filename
func GetImgTint(fn, fe string) (float64, float64, float64) {
	f, fErr := os.Open(fn)
	if fErr != nil {
		fmt.Println("os Err")
		log.Fatal(fErr)
	}
	defer f.Close()
	var imgInfo image.Image
	var imgErr error
	if fe == "png" {
		imgInfo, imgErr = png.Decode(f)
	} else if fe == "jpeg" {
		imgInfo, imgErr = jpeg.Decode(f)
	} else {
	}
	if imgErr != nil {
		fmt.Println("img Err")
		log.Fatal(imgErr)

	}

	var r, g, b, Ry, Gy, By, R, G, B float64
	var red, green, blue uint32
	xSize := imgInfo.Bounds().Max.X
	ySize := imgInfo.Bounds().Max.Y

	for i := 0; i < xSize; i++ {
		for j := 0; j < ySize; j++ {
			red, green, blue, _ = imgInfo.At(i, j).RGBA()
			red = red / 256
			green = green / 256
			blue = blue / 256
			r = (float64)(red)
			g = (float64)(green)
			b = (float64)(blue)
			Ry, Gy, By = Ry+r, Gy+g, By+b
		}
		Ry = Ry / (float64)(ySize)
		Gy = Gy / (float64)(ySize)
		By = By / (float64)(ySize)
		R, G, B = R+Ry, G+Gy, B+By
		Ry, Gy, By = 0, 0, 0
	}
	R = R / (float64)(xSize)
	G = G / (float64)(xSize)
	B = B / (float64)(xSize)
	return R, G, B
}

//GetRGBAInfo : get each pixel RGBA infomation
func GetRGBAInfo(fn, fm string) [][]ImgInfo {
	f, fErr := os.Open(fn)
	if fErr != nil {
		log.Fatal(fErr)
	}
	defer f.Close()
	var imgInfo image.Image
	var imgErr error
	if fm == "png" {
		imgInfo, imgErr = png.Decode(f)
	} else if fm == "jpeg" {
		imgInfo, imgErr = jpeg.Decode(f)
	} else {
	}
	if imgErr != nil {
		fmt.Println("img Err")
		log.Fatal(imgErr)

	}
	xSize, ySize := imgInfo.Bounds().Max.X, imgInfo.Bounds().Max.Y
	var RGBA [][]ImgInfo
	for j := 0; j < ySize; j++ {
		var RGBA0 []ImgInfo
		for i := 0; i < xSize; i++ {
			r, g, b, a := imgInfo.At(i, j).RGBA()
			var rgba ImgInfo
			rgba.R, rgba.G, rgba.B, rgba.A = (int)(r), (int)(g), (int)(b), (int)(a)
			RGBA0 = append(RGBA0, rgba)
		}
		RGBA = append(RGBA, RGBA0)
	}
	return RGBA
}

//GetRGBAAverage : get Average RGBA information from [][]ImgInfo structure
func GetRGBAAverage(info [][]ImgInfo) (int, int, int, int) {
	var Info []ImgInfo
	for _, v := range info {
		Info = append(Info, v...)
	}
	var r, g, b, a float64
	LEN := (float64)(len(Info))
	for _, v := range Info {
		r, g, b, a = r+(float64)(v.R), g+(float64)(v.G), b+(float64)(v.B), a+(float64)(v.A)
	}
	r, g, b, a = r/LEN, g/LEN, b/LEN, a/LEN
	return (int)(r / 256), (int)(g / 256), (int)(b / 256), (int)(a / 256)
}

//GetDiffInfo : get horizontal diff infomation and vertical diff infomation
func GetDiffInfo(info [][]ImgInfo) ([][]float64, [][]float64) {
	var xDiff, yDiff [][]float64
	xLen, yLen := len(info[0]), len(info)

	//var xDiffSq float64
	for _, v := range info {
		var xDiffXAxis []float64
		for i := 0; i < xLen-1; i++ {
			r, g := v[i].R-v[i+1].R, v[i].G-v[i+1].G
			b, a := v[i].B-v[i+1].B, v[i].A-v[i+1].A
			xDiffSq := math.Sqrt((float64)(r*r + g*g + b*b + a*a))
			xDiffXAxis = append(xDiffXAxis, xDiffSq)
		}
		xDiff = append(xDiff, xDiffXAxis)
	}
	//var yDiffSq float64
	for i := 0; i < xLen; i++ {
		var yDiffYAxis []float64
		for j := 0; j < yLen-1; j++ {
			r, g := info[j][i].R-info[j+1][i].R, info[j][i].G-info[j+1][i].G
			b, a := info[j][i].B-info[j+1][i].B, info[j][i].A-info[j+1][i].A
			yDiffSq := math.Sqrt((float64)(r*r + g*g + b*b + a*a))
			yDiffYAxis = append(yDiffYAxis, yDiffSq)
		}
		yDiff = append(yDiff, yDiffYAxis)
	}
	//fmt.Println("xDiff x : ", len(xDiff))
	//fmt.Println("yDiff x : ", len(yDiff))
	return xDiff, yDiff
}

//GetDiffAverage : get Average diff degree from horizontal or vertical
func GetDiffAverage(xD, yD [][]float64) (float64, float64) {
	xD1, xD2 := len(xD[0]), len(xD)
	yD1, yD2 := len(yD[0]), len(yD)
	var xDAver, yDAver float64
	for _, v := range xD {
		for _, v := range v {
			xDAver = xDAver + v
		}
	}
	for _, v := range yD {
		for _, v := range v {
			yDAver = yDAver + v
		}
	}
	return xDAver / (float64)(xD1*xD2), yDAver / (float64)(yD1*yD2)
}

//GetIrregDiff : check diff degree over or not over your Threshold
func GetIrregDiff(D [][]float64, Threshold float64) [][]float64 {
	var irregDiff [][]float64
	for _, v := range D {
		var axisParallel []float64
		for _, v := range v {
			if v > Threshold {
				axisParallel = append(axisParallel, 1.0)
			} else {
				axisParallel = append(axisParallel, 0.0)
			}
		}
		irregDiff = append(irregDiff, axisParallel)
	}
	return irregDiff
}

//DrawFromNRGBA :
func DrawFromNRGBA(name, extsn string, info *image.NRGBA) {
	nameAndExtsn := name + "." + extsn
	f, ferr := os.Create(nameAndExtsn)
	if ferr != nil {
		panic(ferr)
	}
	if extsn == "jpeg" {
		jpegOption := jpeg.Options{
			Quality: 100,
		}
		jpeg.Encode(f, info, &jpegOption)
	} else if extsn == "png" {
		png.Encode(f, info)
	}
}

// DrawJPEGImg :
func DrawJPEGImg(fn string, info [][]ImgInfo) {
	f, fErr := os.Create(fn)
	if fErr != nil {
		log.Fatal(fErr)
	}
	defer f.Close()
	width, height := len(info[0]), len(info)
	fmt.Println("Draw width : ", width, " height : ", height)
	NRGBAImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			NRGBAImg.Set(i, j, color.NRGBA{
				R: (uint8)(info[j][i].R),
				G: (uint8)(info[j][i].G),
				B: (uint8)(info[j][i].B),
				A: (uint8)(info[j][i].A),
			})

		}
	}
	jpegOption := jpeg.Options{
		Quality: 100,
	}
	jpeg.Encode(f, NRGBAImg, &jpegOption)
}

// DrawPNGImg :
func DrawPNGImg(fn string, info [][]ImgInfo) {
	f, fErr := os.Create(fn)
	if fErr != nil {
		log.Fatal(fErr)
	}
	defer f.Close()
	width, height := len(info[0]), len(info)
	fmt.Println("Draw width : ", width, " height : ", height)
	NRGBAImg := image.NewNRGBA(image.Rect(0, 0, width, height))
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			NRGBAImg.Set(i, j, color.NRGBA{
				R: (uint8)(info[j][i].R),
				G: (uint8)(info[j][i].G),
				B: (uint8)(info[j][i].B),
				A: (uint8)(info[j][i].A),
			})

		}
	}
	png.Encode(f, NRGBAImg)
}

//SweepNoise : sweep one pixel noise
func SweepNoise(xIrreg [][]float64) [][]float64 {
	var xIrregCopy [][]float64
	for _, v := range xIrreg {
		xLine := append([]float64(nil), v...)
		xIrregCopy = append(xIrregCopy, xLine)
	}
	yLen, xLen := len(xIrreg), len(xIrreg[0])
	for i := 0; i < yLen-2; i++ {
		for j := 0; j < xLen-2; j++ {
			if xIrregCopy[i][j] < 1 && xIrregCopy[i][j+1] < 1 && xIrregCopy[i][j+2] < 1 &&
				xIrregCopy[i+1][j] < 1 && xIrregCopy[i+1][j+2] < 1 &&
				xIrregCopy[i+2][j] < 1 && xIrregCopy[i+2][j+1] < 1 && xIrregCopy[i+2][j+2] < 1 {
				xIrregCopy[i+1][j+1] = 0.0
			}
		}
	}
	return xIrregCopy
}

//SweepNoise2 : sweep isolated four pixel noise
func SweepNoise2(xIrreg [][]float64) [][]float64 {
	var xIrregCopy [][]float64
	for _, v := range xIrreg {
		xLine := append([]float64(nil), v...)
		xIrregCopy = append(xIrregCopy, xLine)
	}
	yLen, xLen := len(xIrreg), len(xIrreg[0])
	for i := 0; i < yLen-3; i++ {
		for j := 0; j < xLen-3; j++ {
			if xIrregCopy[i][j] < 1 && xIrregCopy[i][j+1] < 1 && xIrregCopy[i][j+2] < 1 && xIrregCopy[i][j+3] < 1 &&
				xIrregCopy[i+1][j] < 1 && xIrregCopy[i+1][j+3] < 1 &&
				xIrregCopy[i+2][j] < 1 && xIrregCopy[i+2][j+3] < 1 &&
				xIrregCopy[i+3][j] < 1 && xIrregCopy[i+3][j+1] < 1 && xIrregCopy[i+3][j+2] < 1 && xIrregCopy[i+3][j+3] < 1 {
				xIrregCopy[i+1][j+1], xIrregCopy[i+1][j+2] = 0.0, 0.0
				xIrregCopy[i+2][j+1], xIrregCopy[i+2][j+2] = 0.0, 0.0
			}
		}
	}
	return xIrregCopy
}

//MakeImgFromXIrreg :
func MakeImgFromXIrreg(Irreg [][]float64) [][]ImgInfo {
	var line []ImgInfo
	var linestack [][]ImgInfo
	width := len(Irreg[0]) + 1
	height := len(Irreg)
	for i := 0; i < width; i++ {
		baseInfo := ImgInfo{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
		line = append(line, baseInfo)
	}
	for j := 0; j < height; j++ {
		adhoc := append([]ImgInfo(nil), line...)
		linestack = append(linestack, adhoc)
	}
	for ve := 0; ve < height-1; ve++ {
		for ho := 0; ho < width-2; ho++ {
			if Irreg[ve][ho] > DIFFTHRESH {
				linestack[ve][ho].R = 255
				linestack[ve][ho].A = 255
			}
		}
	}
	return linestack
}

//MakeImgFromYIrreg :
func MakeImgFromYIrreg(Irreg [][]float64) [][]ImgInfo {
	var line []ImgInfo
	var linestack [][]ImgInfo
	width := len(Irreg)
	height := len(Irreg[0]) + 1
	for i := 0; i < width; i++ {
		baseInfo := ImgInfo{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}
		line = append(line, baseInfo)
	}
	for j := 0; j < height; j++ {
		adhoc := append([]ImgInfo(nil), line...)
		linestack = append(linestack, adhoc)
	}
	for ho := 0; ho < width-1; ho++ {
		for ve := 0; ve < height-2; ve++ {
			if Irreg[ho][ve] > DIFFTHRESH {
				linestack[ve][ho].B = 255
				linestack[ve][ho].A = 255
			}
		}
	}
	return linestack
}

//MakeGraphBlueprint :
func MakeGraphBlueprint(fn string, info [][]ImgInfo) (*image.NRGBA, *image.NRGBA, *image.NRGBA) {
	//path, _ := os.Getwd()
	/*
		name, extsn := mp.GetFileNameAndExtsn(fn)
		fmt.Println("name : ", name)
		fmt.Println("extsn : ", extsn)
		newNameRG := name + "_RG_graph" + extsn
		newNameGB := name + "_GB_graph" + extsn
		newNameRB := name + "_RB_graph" + extsn

			fRG, errRG := os.Create(newNameRG)
			if errRG != nil {
				panic(errRG)
			}
			defer fRG.Close()
			fGB, errGB := os.Create(newNameGB)
			if errGB != nil {
				panic(errGB)
			}
			defer fGB.Close()
			fRB, errRB := os.Create(newNameRB)
			if errRB != nil {
				panic(errRB)
			}
			defer fRB.Close()
	*/

	RGgraph := image.NewNRGBA(image.Rect(0, 0, 255, 255))
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			RGgraph.Set(i, j, color.NRGBA{
				R: (uint8)(i),
				G: 255 - (uint8)(j),
				B: 0,
				A: 192,
			})
		}
	}
	GBgraph := image.NewNRGBA(image.Rect(0, 0, 255, 255))
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			GBgraph.Set(i, j, color.NRGBA{
				R: 0,
				G: (uint8)(j),
				B: 255 - (uint8)(j),
				A: 192,
			})
		}
	}
	RBgraph := image.NewNRGBA(image.Rect(0, 0, 255, 255))
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j++ {
			RBgraph.Set(i, j, color.NRGBA{
				R: (uint8)(i),
				G: 0,
				B: 255 - (uint8)(j),
				A: 192,
			})
		}
	}
	for _, line := range info {
		for _, pixel := range line {
			RGgraph.Set(pixel.R/256, 255-pixel.G/256, color.NRGBA{
				R: 255,
				G: 255,
				B: 0,
				A: 255,
			})
			GBgraph.Set(pixel.G/256, 255-pixel.B/256, color.NRGBA{
				R: 0,
				G: 255,
				B: 255,
				A: 255,
			})
			RBgraph.Set(pixel.R/256, 255-pixel.B/256, color.NRGBA{
				R: 255,
				G: 0,
				B: 255,
				A: 255,
			})
		}
	}

	return RGgraph, GBgraph, RGgraph
}

//ReturnNRGBA : get each pixel RGBA infomation
func ReturnNRGBA(fn, fm string) *image.NRGBA {
	f, fErr := os.Open(fn)
	if fErr != nil {
		log.Fatal(fErr)
	}
	defer f.Close()
	var imgInfo image.Image
	var imgErr error
	if fm == "png" {
		imgInfo, imgErr = png.Decode(f)
	} else if fm == "jpeg" {
		imgInfo, imgErr = jpeg.Decode(f)
	} else {
	}
	if imgErr != nil {
		fmt.Println("img Err")
		log.Fatal(imgErr)

	}
	xSize, ySize := imgInfo.Bounds().Max.X, imgInfo.Bounds().Max.Y
	NRGBA := image.NewNRGBA(image.Rect(0, 0, xSize, ySize))
	for j := 0; j < ySize; j++ {
		for i := 0; i < xSize; i++ {
			NRGBA.Set(i, j, imgInfo.At(i, j))
		}
	}
	return NRGBA
}

//ReturnRGBAAverage : get Average RGBA information from [][]ImgInfo structure
func ReturnRGBAAverage(NRGBA *image.NRGBA) (int, int, int, int) {
	xSize, ySize := NRGBA.Bounds().Max.X, NRGBA.Bounds().Max.Y
	fullSize := xSize * ySize
	var R, G, B, A float64
	for i := 0; i < xSize; i++ {
		for j := 0; j < ySize; j++ {
			r, g, b, a := NRGBA.At(i, j).RGBA()
			R, G, B, A = R+(float64)(r/256), G+(float64)(g/256), B+(float64)(b/256), A+(float64)(a/256)
		}
	}
	R, G, B, A = R/(float64)(fullSize), G/(float64)(fullSize), B/(float64)(fullSize), A/(float64)(fullSize)
	return (int)(R), (int)(G), (int)(B), (int)(A)
}

//ReturnDiffInfo : get horizontal diff infomation and vertical diff infomation
func ReturnDiffInfo(info *image.NRGBA) ([][]RGBA, [][]RGBA) {
	MaxX, MaxY := info.Bounds().Max.X, info.Bounds().Max.Y
	var xDiff, yDiff [][]RGBA

	for j := 0; j < MaxY; j++ {
		var xDiffHorizontal []RGBA
		for i := 0; i < MaxX-1; i++ {
			r0, g0, b0, a0 := info.At(i, j).RGBA()
			r, g, b, a := info.At(i+1, j).RGBA()
			R := (int)(r0/256) - (int)(r/256)
			G := (int)(g0/256) - (int)(g/256)
			B := (int)(b0/256) - (int)(b/256)
			A := (int)(a0/256) - (int)(a/256)
			rgba := RGBA{R, G, B, A}
			xDiffHorizontal = append(xDiffHorizontal, rgba)
		}
		xDiff = append(xDiff, xDiffHorizontal)
	}

	for i := 0; i < MaxX; i++ {
		var yDiffVertical []RGBA
		for j := 0; j < MaxY-1; j++ {
			r0, g0, b0, a0 := info.At(i, j).RGBA()
			r, g, b, a := info.At(i, j+1).RGBA()
			R := (int)(r0/256) - (int)(r/256)
			G := (int)(g0/256) - (int)(g/256)
			B := (int)(b0/256) - (int)(b/256)
			A := (int)(a0/256) - (int)(a/256)
			rgba := RGBA{R, G, B, A}
			yDiffVertical = append(yDiffVertical, rgba)
		}
		yDiff = append(yDiff, yDiffVertical)
	}

	return xDiff, yDiff
}

//ReturnDiffSquareInfo :
func ReturnDiffSquareInfo(xDiff [][]RGBA, yDiff [][]RGBA) ([][]float64, [][]float64, float64) {
	var SX, SY, S float64
	var SqX, SqY [][]float64
	Xs, Ys := len(xDiff)*len(xDiff[0]), len(yDiff)*len(yDiff[0])
	for _, V := range xDiff {
		var xHorizon []float64
		for _, v := range V {
			S = math.Sqrt((float64)(v.R*v.R + v.G*v.G + v.B*v.B + v.A*v.A))
			SX = SX + S
			xHorizon = append(xHorizon, S)
			//SqX[mut1][mut2] = S
			//mut2++
		}
		SqX = append(SqX, xHorizon)
	}
	for _, V := range yDiff {
		var yVertical []float64
		for _, v := range V {
			S = math.Sqrt((float64)(v.R*v.R + v.G*v.G + v.B*v.B + v.A*v.A))
			SY = SY + S
			yVertical = append(yVertical, S)
		}
		SqY = append(SqY, yVertical)
	}
	return SqX, SqY, (SX / (float64)(Xs)) + (SY / (float64)(Ys))
}

// ReturnSimplify :
func ReturnSimplify(info [][]float64, Threshold float64) [][]float64 {
	var mut, VLen int
	for _, V := range info {
		VLen = len(V)
		for index := 0; index < VLen-1; index++ {
			if V[index] > Threshold {
				mut = 0
				for V[index+mut] > Threshold && index+mut < VLen-1 {
					V[index+mut] = 0.0
					mut++
				}
				V[index+(mut/2)] = Threshold + 1.0
			}
		}
	}
	return info
}

//ReturnIrregDiffImg : check diff degree over or not over your Threshold
func ReturnIrregDiffImg(xInfo, yInfo [][]float64, Threshold float64) *image.NRGBA {

	_, xYLen, yXLen, _ := len(xInfo[0]), len(xInfo), len(yInfo), len(yInfo[0])
	diffImg := image.NewNRGBA(image.Rect(0, 0, yXLen, xYLen))
	var mutX, mutY int

	mutX, mutY = 0, 0
	for _, V := range yInfo {
		for _, v := range V {
			if v > Threshold {
				diffImg.Set(mutX, mutY, color.NRGBA{
					B: 255,
					A: 255,
				})
			}
			mutY++
		}
		mutY = 0
		mutX++

	}
	mutX, mutY = 0, 0
	for _, V := range xInfo {
		for _, v := range V {
			if v > Threshold {
				diffImg.Set(mutX, mutY, color.NRGBA{
					R: 255,
					A: 255,
				})
			}
			mutX++
		}
		mutX = 0
		mutY++
	}
	return diffImg
}
