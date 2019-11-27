package mygraph

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

//Data :
type Data struct {
	V1 interface{}
	V2 int
}

//DataInt :
type DataInt struct {
	V1 int
	V2 int
}

//DotInt :
func DotInt(data []DataInt, filename string) {
	var maxX, maxY int
	maxX, maxY = 0, 0
	for _, v := range data {
		if maxX < v.V1 {
			maxX = v.V1
		}
		if maxY < v.V2 {
			maxY = v.V2
		}
	}

	NRGBAImg := image.NewNRGBA(image.Rect(0, 0, maxX+10, maxY+10))
	for i := 0; i < maxX; i++ {
		for j := 0; j < 5; j++ {
			NRGBAImg.Set(i, j, color.NRGBA{
				G: (uint8)(255),
				A: (uint8)(255),
			})
		}
		for j := maxY - 4; j < maxY; j++ {
			NRGBAImg.Set(i, j, color.NRGBA{
				G: (uint8)(255),
				A: (uint8)(255),
			})
		}
	}
	for j := 5; j < maxY-5; j++ {
		for i := 0; i < 5; i++ {
			NRGBAImg.Set(i, j, color.NRGBA{
				G: (uint8)(255),
				A: (uint8)(255),
			})
		}
		for i := maxX - 4; i < maxX; i++ {
			NRGBAImg.Set(i, j, color.NRGBA{
				G: (uint8)(255),
				A: (uint8)(255),
			})
		}
	}
	for _, v := range data {
		NRGBAImg.Set(5+v.V1, 5+v.V2, color.NRGBA{
			R: (uint8)(255),
			G: (uint8)(255),
			B: (uint8)(255),
			A: (uint8)(255),
		})
	}
	f, ferr := os.Create(filename)
	if ferr != nil {
		panic(ferr)
	}
	png.Encode(f, NRGBAImg)
}

//DataTable :
func DataTable(data []Data) []Data {
	var table []Data
	var check = false
	for _, Vfrom := range data {
		for _, Vto := range table {
			if Vfrom.V1 == Vto.V1 {
				check = true
				Vto.V2++
				break
			}
		}
		if check == false {
			table = append(table, Data{Vfrom.V1, 1})
		}
		check = false
	}
	return table
}

//DataTableInt :
func DataTableInt(data []DataInt) []DataInt {
	var table []DataInt
	var check = false
	for _, Vfrom := range data {
		for _, Vto := range table {
			if Vfrom.V1 == Vto.V1 {
				check = true
				Vto.V2++
				break
			}
		}
		if check == false {
			table = append(table, DataInt{Vfrom.V1, 1})
		}
		check = false
	}
	return table
}

//FlatImg :
func FlatImg(data [][]interface{}) []interface{} {
	var flat []interface{}
	for _, v := range data {
		flat = append(flat, v...)
	}
	return flat
}

//Flat :
func Flat(data [][]interface{}) []interface{} {
	var flat []interface{}
	for _, v := range data {
		flat = append(flat, v...)
	}
	return flat
}

//FlatFloat64 :
func FlatFloat64(data [][]float64) []float64 {
	var flat []float64
	for _, v := range data {
		flat = append(flat, v...)
	}
	return flat
}

//FlatInt :
func FlatInt(data [][]int) []int {
	var flat []int
	for _, v := range data {
		flat = append(flat, v...)
	}
	return flat
}
