package config

type DataSourceName struct {
	buf []byte
}

func NewDSN(ss ...string) (d *DataSourceName) {
	d = new(DataSourceName)
	if len(ss) > 0 {
		d = &DataSourceName{
			buf: []byte(ss[0]),
		}
	}
	return
}

func (d *DataSourceName) DSN() string {
	return d.String()
}

func (d *DataSourceName) String() string {
	return string(d.buf)
}

func (d *DataSourceName) Write(p []byte) (n int, err error) {
	n = len(p)
	d.buf = append(d.buf, p...)
	return
}

type CaseString = DataSourceName
