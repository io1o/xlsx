package xlsx

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"github.com/unidoc/unioffice/spreadsheet"
)

func createOption(optionFns []OptionFn) *Option {
	option := &Option{}

	for _, fn := range optionFns {
		fn(option)
	}

	return option
}

// Option defines the option for the xlsx processing.
type Option struct {
	TemplateWorkbook *spreadsheet.Workbook
	Workbook         *spreadsheet.Workbook
	httpUpload       *upload
	Validations      map[string][]string
	Placeholder      bool
}

// OptionFn defines the func to change the option.
type OptionFn func(*Option)

// WithTemplate defines the template excel file for writing template.
func WithTemplate(f interface{}) OptionFn {
	wb, err := parseExcel(f)
	if err != nil {
		logrus.Warnf("failed to open template excel %v", err)
		return func(*Option) {}
	}

	return func(o *Option) { o.TemplateWorkbook = wb }
}

// AsPlaceholder defines the template excel file for writing template in placeholder mode.
func AsPlaceholder() OptionFn {
	return func(o *Option) { o.Placeholder = true }
}

// WithExcel defines the input excel file for reading.
func WithExcel(f interface{}) OptionFn {
	wb, err := parseExcel(f)
	if err != nil {
		logrus.Warnf("failed to open excel %v", err)
		return func(*Option) {}
	}

	return func(o *Option) { o.Workbook = wb }
}

func parseExcel(f interface{}) (wb *spreadsheet.Workbook, err error) {
	var bs []byte

	switch ft := f.(type) {
	case string:
		return spreadsheet.Open(ft)
	case []byte:
		return spreadsheet.Read(bytes.NewReader(ft), int64(len(ft)))
	case io.Reader:
		bs, err = ioutil.ReadAll(ft)
		if err != nil {
			return nil, err
		}

		return spreadsheet.Read(bytes.NewReader(bs), int64(len(bs)))
	default:
		return nil, fmt.Errorf("unknown excel file format")
	}
}

// WithValidations defines the validations for the cells.
func WithValidations(v map[string][]string) OptionFn {
	return func(o *Option) { o.Validations = v }
}
