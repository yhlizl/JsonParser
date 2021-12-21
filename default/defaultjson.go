package defaultjson

import (
	"fmt"
	"io/fs"
	"os"
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

func initExcel(data map[string]map[string]string,keyList map[string]bool){
    f := excelize.NewFile()
    // Create a new sheet.
    index := f.NewSheet("default")
    // Set value of a cell.
    f.SetCellValue("default", "A1", "internalName")
    // Set active sheet of the workbook.
    f.SetActiveSheet(index)
    // get key list
	keyStrList:=[]string{}
	for i:=range keyList{
		keyStrList=append(keyStrList,i)
	}
	//sort list
	sort.Strings(keyStrList)
	// tranverse to create data value
	x:=2
	for device:=range data {
		deviceCells:=XY(1,x)
		f.SetCellValue("default", deviceCells, device)
		for y,v :=range  keyStrList{
			cells:=XY(y+2,1)
			f.SetCellValue("default", cells, v)
			resultCells:=XY(y+2,x)
			if _,ok:=data[device][v];ok{
				
				f.SetCellValue("default", resultCells, data[device][v])
			}else{
				f.SetCellValue("default", resultCells, "---")

			}
			
		}
		x++
	}

	// Save spreadsheet by the given path.

    if err := f.SaveAs("result.xlsx"); err != nil {
        fmt.Println(err)
    }
}
func DefaultJson(location string) {


err:=filepath.WalkDir(location, visit)
if err !=nil{
	fmt.Println("err for walk root json:", err)
}
fmt.Println(pathList)

walkFunction(pathList)

}


func visit(path string, di fs.DirEntry, err error) error {
    fmt.Printf("Visited: %s %s %t\n", path,di.Name(),di.IsDir())
	if !di.IsDir()&&strings.Contains(di.Name(),"defaults.json") {
		pathList=append(pathList,path)
	}
    return nil
}

func walkFunction(pathlist []string)  {

	resultMap:=map[string]map[string]string{}
	list:=map[string]bool{}
for _,p := range pathList{
	file , err := os.Open(p)
	if err == nil {
		} else {
		
		}
stat , _ := file.Stat()

bs := make([]byte, stat.Size())
_ ,err = file.Read(bs)
if err != nil{
    fmt.Println(err)
}
nameList := gjson.Parse(string(bs)).Map()

for device,data := range nameList {
	if _,ok:=resultMap[device];!ok{
		resultMap[device]=map[string]string{}
	}
	params:=data.Map()
	for key,value:=range params{
		resultMap[device][key]=value.String()
		list[key]=true
	}
	
}
	// fmt.Println(resultMap,list)
}
	initExcel(resultMap,list)
}
