package main

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/fatih/structs"
)

// GetTableFields GetTableFields
func GetTableFields(v interface{}, exclude, duplicateKey []string) (fields, updateOnDuplicate []string) {
	for _, field := range structs.Fields(v) {
		if field.Tag("db") == "-" {
			continue
		}

		if IsStrInSlice(field.Tag("db"), exclude) {
			continue
		}
		fields = append(fields, field.Tag("db"))

		if IsStrInSlice(field.Tag("db"), duplicateKey) {
			continue
		}
		updateOnDuplicate = append(updateOnDuplicate, field.Tag("db")+`=VALUES(`+field.Tag("db")+`)`)
	}

	return
}

// IsStrInSlice IsStrInSlice
func IsStrInSlice(target string, slice []string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

// ToLowerCamel ToLowerCamel
func ToLowerCamel(name string) string {
	name = strings.ReplaceAll(name, "ID", "Id")
	for i, v := range name {
		return string(unicode.ToLower(v)) + name[i+1:]
	}
	return name
}

// AnySliceToInterfaceSlice AnySliceToInterfaceSlice
func AnySliceToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	ret := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
