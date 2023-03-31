package dalkeeth

import (
	"context"
	"errors"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var NotImplemented = errors.New("Not implemented")

func countSetValues(values []*Value) int {
	count := 0
	for i := 0; i < len(values); i++ {
		if !values[i].isSet {
			count++
		}
	}
	return count
}

func rawValues(values []*Value) []any {
	v := make([]any, len(values)-countSetValues(values))
	index := 0
	for i := 0; i < len(values); i++ {
		if values[i].isSet {
			v[index] = values[i].value
			index++
		}
	}
	return v
}

func rawWantedValues(values []*Value) ([]any, error) {
	var v []any
	var err error
	for i := 0; i < len(values); i++ {
		if values[i].isWanted {
			values[i].value, err = makeValue(values[i].field.fieldType)
			if err != nil {
				return nil, err
			}
			v = append(v, values[i].value)
		}
	}
	return v, nil
}

func actual(v *Value) {
	switch p := v.value.(type) {
	case *int64:
		v.value = *p
		log.Println("actual int64", p, *p)
	}
}

func makeValue(fieldType FieldType) (any, error) {
	switch fieldType {
	case IntType:
		var v int64
		return &v, nil

	case StringType:
		var v string
		return &v, nil

	case FloatType:
		var v float64
		return &v, nil

	case BoolType:
		var v bool
		return &v, nil
	}

	return nil, errors.New("Unknown field type")

}

func cancelIfError(err error, cancel context.CancelFunc) {
	if err != nil {
		cancel()
	}
}

// From: https://stackoverflow.com/a/32482941
func DirSize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size, err
}

func urlDir(dir, urlRaw string) string {
	u, err := url.Parse(urlRaw)
	if err != nil {
		log.Println("Error parsing url:", urlRaw)
		log.Fatal(err)
	}
	path, err := url.JoinPath(dir, u.Hostname(), u.EscapedPath())
	if err != nil {
		log.Fatal(err)
	}
	path = strings.TrimRight(path, ".git")
	return path
}

func zeroLength[V any](a []V, label string) error {
	if len(a) == 0 {
		return errors.New("Need len>0 for " + label)
	}
	return nil
}
