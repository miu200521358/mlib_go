package csv

type CsvModel struct {
	records [][]string
	path    string
}

func NewCsvModel(records [][]string) *CsvModel {
	return &CsvModel{
		records: records,
	}
}

func (d *CsvModel) Records() [][]string {
	return d.records
}

func (d *CsvModel) Cell(row, col int) string {
	if row >= len(d.records) || col >= len(d.records[row]) {
		return ""
	}
	return d.records[row][col]
}

func (d *CsvModel) SetCell(row, col int, value string) {
	if row >= len(d.records) {
		return
	}
	if col >= len(d.records[row]) {
		return
	}
	d.records[row][col] = value
}

func (d *CsvModel) Hash() string {
	return ""
}

func (d *CsvModel) Name() string {
	return ""
}

func (d *CsvModel) Path() string {
	return d.path
}

func (d *CsvModel) SetHash(hash string) {
}

func (d *CsvModel) SetName(name string) {
}

func (d *CsvModel) SetPath(path string) {
	d.path = path
}

func (d *CsvModel) UpdateHash() {
}

func (d *CsvModel) SetRandHash() {
}
