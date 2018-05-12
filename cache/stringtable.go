package cache

import "github.com/galaco/source-tools-common/texdatastringtable"

var texDataStringTable texdatastringtable.TexDataStringTable

func CreateTexDataStringTable(stringData string, stringTable []int32) {
	texDataStringTable = *texdatastringtable.NewTable(stringData, stringTable)
}

func GetTexDataStringTable() *texdatastringtable.TexDataStringTable {
	return &texDataStringTable
}