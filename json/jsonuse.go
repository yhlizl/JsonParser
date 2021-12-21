package jsonuse

import (
	"fmt"
	"io/fs"
	"os"
	pathuse "path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	excelize "github.com/xuri/excelize/v2"
)
var pathList = []string{}

func XY(y, x int) string {
	var r []byte
	if x < 26 {
		r = make([]byte, 1, 5)
		r[0] = 'A' + byte(x)
	} else if x < 27*26 {
		r = make([]byte, 2, 5)
		r[0] = 'A' - 1 + byte(x/26)
		r[1] = 'A' + byte(x%26)
	} else if x < 16384 {
		r = make([]byte, 3, 5)
		r[2] = 'A' + byte(x%26)
		x /= 26
		r[0] = 'A' - 1 + byte(x/26)
		r[1] = 'A' - 1 + byte(x%26)
	} else {
		panic(fmt.Errorf("more than 16384 columns: %d", x))
	}
	return string(strconv.AppendUint(r, uint64(y+1), 10))
}

func initExcel(data map[string]map[string]string,deviceList []string){
    f := excelize.NewFile()
    // Create a new sheet.
    index := f.NewSheet("Compare")
    // Set value of a cell.
    _=f.SetCellValue("Compare", "A1", "internalName")
    _=f.SetCellValue("Compare", "B1", "ReviewName")
    // Set active sheet of the workbook.
    f.SetActiveSheet(index)
    // tranverse to create data value
	//sort list
	sort.Strings(deviceList)
	y:=2
	for internalName:=range data {
		RVcells:=XY(y,1)
		f.SetCellValue("Compare", RVcells, internalName)
		for dint :=range  deviceList{
			cells:=XY(1,dint+2)
			f.SetCellValue("Compare", cells, deviceList[dint])
			resultCells:=XY(y,dint+2)
			if _,ok:=data[internalName][deviceList[dint]];ok{
				
				f.SetCellValue("Compare", resultCells, data[internalName][deviceList[dint]])
			}else{
				f.SetCellValue("Compare", resultCells, "---")

			}
		}
		y++
	}

	// Save spreadsheet by the given path.

    if err := f.SaveAs("result.xlsx"); err != nil {
        fmt.Println(err)
    }
}
func Jsonuse(location string) {


err:=filepath.WalkDir(location, visit)
if err !=nil{
	fmt.Println("err for walk root json:", err)
}
fmt.Println(pathList)

walkFunction(pathList)

}


func visit(path string, di fs.DirEntry, err error) error {
    fmt.Printf("Visited: %s %s %t\n", path,di.Name(),di.IsDir())
	if !di.IsDir()&&strings.Contains(di.Name(),"json") {
		pathList=append(pathList,path)
	}
    return nil
}

func walkFunction(pathlist []string)  {

	resultMap:=map[string]map[string]string{}
	list:=[]string{}
for _,p := range pathList{
	file , err := os.Open(p)
	if err == nil {
		} else {
		
		}
stat , _ := file.Stat()
name:=pathuse.Base(file.Name())[:len(pathuse.Base(file.Name()))-5];
list=append(list,name)
bs := make([]byte, stat.Size())
_ ,err = file.Read(bs)
if err != nil{
    fmt.Println(err)
}
result := gjson.Get(string(bs), "#.internalName")
reviewName := gjson.Get(string(bs), "#.reviewName")
//raw := gjson.Get(string(bs), "#")
//fmt.Println(len(result.Array()),len(reviewName.Array()),len(choices.Array()))
// fmt.Println(choices)
for i := range result.Array() {
	if _,ok:=resultMap[result.Array()[i].String()+";"+reviewName.Array()[i].String()];!ok{
		resultMap[result.Array()[i].String()+";"+reviewName.Array()[i].String()]=map[string]string{}
	}
	num:=strconv.Itoa(i)
	choices := gjson.Get(string(bs), num+".choices")
	interSlice:=[]string{}
	for j:=0;j< len(choices.Array());j++ {
		inter := gjson.Get(choices.Array()[j].Raw, "internalValue").String()
		rv := gjson.Get(choices.Array()[j].Raw, "reviewValue").String()
		str:=rv+"<<"+inter+">>"
		interSlice=append(interSlice,str)
	}
	//fmt.Println(name,result.Array()[i],reviewName.Array()[i],interSlice)
	interStr:=strings.Join(interSlice,";")
	resultMap[result.Array()[i].String()+"("+reviewName.Array()[i].String()+")"][name]=interStr

}

}
initExcel(resultMap,list)
	}
