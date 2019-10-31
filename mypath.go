package mypath

import (
	"fmt"
	"io"
	"log"
	"os"
	fp "path/filepath"
	"regexp"
	"strconv"
	"strings"
)

//OldAndNew : old filename and extension, new filename, new extention
type OldAndNew struct {
	old       string
	new       string
	newExtnsn string
}

func stringCheck(strslc []string, kwrd string) bool {
	var flag = false
	for _, v := range strslc {
		if v == kwrd {
			flag = true
		}
	}
	return flag
}

//MkDir : make Folder in current location. If same name folder exist, add number at end
func MkDir(fdn string) string {
	wkfd, wErr := os.Getwd()
	if wErr != nil {
		log.Fatal(wErr)
	}
	wkfdFs := wkfd + "\\*"
	wkfdList, wkErr := fp.Glob(wkfdFs)
	if wkErr != nil {
		log.Fatal(wkErr)
	}
	wkfdFsRegE := `\S+\.\S+$`
	wkfdFsRegEP, regErr := regexp.Compile(wkfdFsRegE)
	if regErr != nil {
		log.Fatal(regErr)
	}

	var wkfdSlc []string
	for _, v := range wkfdList {
		if wkfdFsRegEP.MatchString(v) {
		} else {
			wkfdSlc = append(wkfdSlc, v)
		}
	}
	for _, v := range wkfdSlc {
		fmt.Println(v)
	}

	wkfdSlcBase := func(paths []string) []string {
		var bases []string
		for _, v := range paths {
			bases = append(bases, fp.Base(v))
		}
		return bases
	}(wkfdSlc)

	stringCheck := func(strslc []string, kwrd string) bool {
		var flag = false
		for _, v := range strslc {
			if v == kwrd {
				flag = true
			}
		}
		return flag
	}

	var FDN string
	if stringCheck(wkfdSlcBase, fdn) == false {
		FDN = fdn
	} else {
		var fnum = 0
		var fdnFlag = true
		for fdnFlag == true {
			fnum++
			FDN = fdn + "_" + strconv.Itoa(fnum)
			fdnFlag = stringCheck(wkfdSlcBase, FDN)
		}
	}
	fmt.Println("FDN : ", FDN)
	mkDirErr := os.Mkdir(FDN, 0777)
	if mkDirErr != nil {
		return "MkDir Error"
	}
	return FDN
}

//MoveToFolder :
func MoveToFolder(infoSlice []OldAndNew, folderName string) {
	path, _ := os.Getwd()
	destinatioin := path + "\\" + folderName + "\\"
	os.Mkdir(folderName, 0777)
	for _, v := range infoSlice {
		src, _ := os.Open(v.old)
		defer src.Close()
		dstName := destinatioin + v.new + v.newExtnsn
		fmt.Println(dstName)
		dst, _ := os.Create(dstName)
		_, cErr := io.Copy(dst, src)
		if cErr != nil {
			log.Fatal(cErr)
		}
	}
}

//GetFilePathNameSlice :
func GetFilePathNameSlice(ext ...string) []string {
	base, _ := os.Getwd()
	var strslc []string
	for _, v := range ext {
		strslc = append(strslc, base+"\\*."+v)
	}
	var pathsSlice []string
	for _, v := range strslc {
		paths, _ := fp.Glob(v)
		pathsSlice = append(pathsSlice, paths...)
	}
	return pathsSlice
}

//GetFileNameAndExtsn :
func GetFileNameAndExtsn(name string) (string, string) {
	regP, basename := regexp.MustCompile(`\.\S+$`), fp.Base(name)
	extsn := regP.FindString(basename)
	nameWitoutExtsn := strings.Replace(basename, extsn, "", 1)
	return nameWitoutExtsn, extsn
}

//FileNumbering : change filename with number
func FileNumbering(fn []string) {
	var baseName []string
	for _, v := range fn {
		baseName = append(baseName, fp.Base(v))
	}

	getNumber := func(n, i int) string {
		var kn = n
		var cn = 0
		var ki = i
		var ci = 0
		var prefix = ""
		for kn >= 1 {
			kn = kn / 10
			cn++
		}
		for ki >= 1 {
			ki = ki / 10
			ci++
		}
		for d := cn - ci; d > 0; d-- {
			prefix = prefix + "0"
		}
		return prefix + strconv.Itoa(i)
	}

	i, n := 0, len(baseName)
	for _, v := range baseName {
		i++
		rp, _ := regexp.Compile(`\.\S+$`)
		newName := getNumber(n, i) + rp.FindString(v)
		fmt.Println(newName)
		os.Rename(v, newName)
	}

}
