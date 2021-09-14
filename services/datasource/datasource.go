package datasource

import (
	"github.com/xuri/excelize/v2"
)

func New () ([][]string, error) {
	f, err := excelize.OpenFile("source.xlsx")
	if err != nil {
        return nil, err
    }

	rows, err := f.GetRows("source")
    if err != nil {
        return nil, err
    }

	return rows, nil
}
