package pisco

import (
	"fmt"
	"fuzzy/pisco/ansicolor-master"
	"math"
	"os"
	"reflect"
)

type TableFrame struct {
	table  [][]float64
	size   int
	indexs []string
	labels []string
}

func New(table [][]float64, opts ...map[string]interface{}) TableFrame {
	if len(opts) > 1 {
		panic("panic")
	}

	index := ls(opts[0], "indexs", len(table))
	labels := ls(opts[0], "label", len(table[0]))

	return TableFrame{table: table, size: len(table[0]), indexs: index, labels: labels}
}

func (f TableFrame) Print(precision ...int) {
	var i int
	var d int
	var x int

	w := ansicolor.NewAnsiColorWriter(os.Stdout)

	if len(precision) == 1 {
		d = precision[0]
	} else {
		d = 2
	}

	for _, label := range f.labels {
		if len(label) > i {
			i = len(label) + 1
		}
	}

	for _, t := range f.table {
		for _, v := range t {
			s := len(fmt.Sprintf(fmt.Sprintf("%%.%df", d), v))
			if s > i {
				i = s + 1
			}
		}
	}

	for _, idx := range f.indexs {
		s := len(idx)
		if s > x {
			x = s + 1
		}
	}

	text := fmt.Sprintf("%%%ds", x)
	text = fmt.Sprintf(text, " ")
	fmt.Fprintf(w, fmt.Sprintf("%%s %s %%s", text), "\x1b[47;30m", "\x1b[0m")

	for _, txt := range f.labels {
		text = fmt.Sprintf("%%%ds", i)

		if len(txt) > i {
			txt = txt[:i-1] + "."
		}

		text = fmt.Sprintf(text, txt)
		fmt.Fprintf(w, fmt.Sprintf("%%s %s %%s", text), "\x1b[47;30m", "\x1b[0m")
	}

	println()

	for indx, line := range f.table {
		text = fmt.Sprintf("%%%ds", x)
		text = fmt.Sprintf(text, f.indexs[indx])
		fmt.Fprintf(w, fmt.Sprintf("%%s %s %%s", text), "\x1b[47;30m", "\x1b[0m")

		for _, column := range line {
			text = fmt.Sprintf("%%%d.%df ", i+1, d)
			fmt.Printf(text, column)
		}

		println()
	}

	println()

}

func (f *TableFrame) MaxColumn(inPlace bool) TableFrame {
	if !inPlace {
		t := f.copyTable()
		columMax(&t)

		return t
	}

	columMax(f)

	return TableFrame{}
}

func (f *TableFrame) MinColumn(inPlace bool) TableFrame {
	if !inPlace {
		t := f.copyTable()
		columMin(&t)

		return t
	}

	columMin(f)

	return TableFrame{}
}

func (f *TableFrame) MeanColumn(inPlace bool) TableFrame {
	if !inPlace {
		t := f.copyTable()
		columMean(&t)

		return t
	}

	columMean(f)

	return TableFrame{}
}

func (f *TableFrame) MaxLine(inPlace bool) TableFrame {
	if !inPlace {
		t := f.copyTable()
		lineMax(&t)

		return t
	}

	lineMax(f)

	return TableFrame{}
}

func (f *TableFrame) MinLine(inPlace bool) TableFrame {
	if !inPlace {
		t := f.copyTable()
		lineMin(&t)

		return t
	}

	lineMin(f)

	return TableFrame{}
}

func (f *TableFrame) MeanLine(inPlace bool) TableFrame {
	if !inPlace {
		t := f.copyTable()
		lineMean(&t)

		return t
	}

	lineMean(f)

	return TableFrame{}
}

func (f *TableFrame) copyTable() TableFrame {
	table := make([][]float64, 0, len(f.table))
	for _, t := range f.table {
		add := make([]float64, len(t))
		copy(add, t)
		table = append(table, add)
	}

	return TableFrame{table, f.size, f.indexs, f.labels}
}

func columMax(table *TableFrame) {
	var max float64
	column := make([]float64, 0, len(table.table[0]))

	for i := 0; i < len(table.table[0]); i++ {
		max = 0
		for j := 0; j < len(table.table); j++ {
			if table.table[j][i] > max {
				max = table.table[j][i]
			}
		}
		column = append(column, max)
	}

	table.table = append(table.table, column)
	table.indexs = append(table.indexs, "MAX")
}

func columMin(table *TableFrame) {
	var min float64
	column := make([]float64, 0, len(table.table[0]))

	for i := 0; i < len(table.table[0]); i++ {
		min = math.MaxFloat32
		for j := 0; j < len(table.table); j++ {
			if table.table[j][i] < min {
				min = table.table[j][i]
			}
		}
		column = append(column, min)
	}

	table.table = append(table.table, column)
	table.indexs = append(table.indexs, "MIN")
}

func columMean(table *TableFrame) {
	var mean float64
	column := make([]float64, 0, len(table.table[0]))

	for i := 0; i < len(table.table[0]); i++ {
		mean = 0
		for j := 0; j < len(table.table); j++ {
			mean += table.table[j][i]
		}
		column = append(column, mean/float64(len(table.table)))
	}

	table.table = append(table.table, column)
	table.indexs = append(table.indexs, "MEAN")
}

func lineMax(table *TableFrame) {
	var max float64
	line := make([]float64, 0, len(table.table))

	for i := 0; i < len(table.table); i++ {
		max = 0
		for j := 0; j < len(table.table[0]); j++ {
			if table.table[i][j] > max {
				max = table.table[i][j]
			}
		}
		line = append(line, max)
	}

	for i := 0; i < len(table.table); i++ {
		table.table[i] = append(table.table[i], line[i])
	}

	table.labels = append(table.labels, "MAX")
	table.size += 1
}

func lineMin(table *TableFrame) {
	var min float64
	line := make([]float64, 0, len(table.table))

	for i := 0; i < len(table.table); i++ {
		min = math.MaxFloat32
		for j := 0; j < len(table.table[0]); j++ {
			if table.table[i][j] < min {
				min = table.table[i][j]
			}
		}
		line = append(line, min)
	}

	for i := 0; i < len(table.table); i++ {
		table.table[i] = append(table.table[i], line[i])
	}

	table.labels = append(table.labels, "MIN")
	table.size += 1
}

func lineMean(table *TableFrame) {
	var mean float64
	line := make([]float64, 0, len(table.table))

	for i := 0; i < len(table.table); i++ {
		mean = 0
		for j := 0; j < len(table.table[0]); j++ {
			mean += table.table[i][j]
		}
		line = append(line, mean/float64(len(table.table[0])))
	}

	for i := 0; i < len(table.table); i++ {
		table.table[i] = append(table.table[i], line[i])
	}

	table.labels = append(table.labels, "MIN")
	table.size += 1
}

func ls(opt map[string]interface{}, str string, size int) []string {
	lts := make([]string, 0, size)

	if val, ok := opt[str]; ok {
		switch val.(type) {
		case []string:
			s := reflect.ValueOf(val)

			if s.Len() != size {
				panic("Panic")
			}

			for i := 0; i < s.Len(); i++ {
				lts = append(lts, s.Index(i).String())
			}
		}
	} else {
		for i := 0; i < size; i++ {
			lts = append(lts, fmt.Sprintf("%d", i))
		}
	}

	return lts
}
