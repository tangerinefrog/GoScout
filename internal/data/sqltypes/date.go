package sqltypes

import "time"

const timeFormat = "2006-01-02 15:04:05-07:00"

func ParseTimeFromSql(input string) (time.Time, error) {
	t, err := time.Parse(timeFormat, input)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func TimeToSqlFormat(t time.Time) string {
	return t.Format(timeFormat)
}
