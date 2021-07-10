package excel

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strings"
)

type Reader interface {
	Fields() []string
	SheetName() string
	Each(callback func(dataMap map[string]string) bool) error
}

type excelReaderImpl struct {
	sheetName       string
	dataStartedFrom int
	fieldMap        map[int]string
	fields          []string
	excelFile       *excelize.File
}

func New(fileName string) (Reader, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}

	reader := &excelReaderImpl{
		excelFile:       f,
		sheetName:       f.GetSheetName(0),
		fieldMap:        make(map[int]string),
		dataStartedFrom: 0,
	}

	if r, err := f.Rows(reader.sheetName); err != nil {
		fmt.Println(err)
		return nil, err
	} else {
		for rowCount := 0; r.Next(); rowCount++ {
			if columns, err := r.Columns(); err != nil {
				fmt.Println(err)
				return nil, err
			} else {
				haveData := false
				for _, content := range columns {
					if strings.TrimSpace(content) != `` {
						haveData = true
					}
				}
				if haveData {
					for i, content := range columns {
						data := strings.TrimSpace(content)

						if data != `` {
							reader.fieldMap[i] = data
							reader.fields = append(reader.fields, data)
						}
					}
					reader.dataStartedFrom = rowCount
					//fmt.Println("RowCount:", rowCount)
					break
				}
			}
		}
	}
	return reader, nil
}

func (reader *excelReaderImpl) Fields() []string {
	return reader.fields
}

func (reader *excelReaderImpl) SheetName() string {
	return reader.sheetName
}

func (reader *excelReaderImpl) Each(callback func(dataMap map[string]string) bool) error {
	if r, err := reader.excelFile.Rows(reader.sheetName); err != nil {
		return err
	} else {
		for rowCount := 0; r.Next(); rowCount++ {
			if rowCount <= reader.dataStartedFrom {
				_, _ = r.Columns()
			} else {
				if columns, err := r.Columns(); err == nil {
					dataMap := make(map[string]string)

					for i, content := range columns {
						data := strings.TrimSpace(content)

						if data != `` {
							if key, found := reader.fieldMap[i]; found {
								dataMap[key] = data
							}
						}
					}
					if len(dataMap) > 0 {
						if !callback(dataMap) {
							return nil
						}
					}
				}
			}
		}
		return nil
	}
}
