package psql

import "testing"

func TestDatePart(t *testing.T) {
	cases := []struct {
		field DateField
		sql   string
	}{
		{
			CenturyField,
			`date_part('century', now())`,
		},
		{
			DayField,
			`date_part('day', now())`,
		},
		{
			DecadeField,
			`date_part('decade', now())`,
		},
		{
			DayOfWeekField,
			`date_part('dow', now())`,
		},
		{
			DayOfYearField,
			`date_part('doy', now())`,
		},
		{
			EpochField,
			`date_part('epoch', now())`,
		},
		{
			HourField,
			`date_part('hour', now())`,
		},
		{
			ISODayOfWeekField,
			`date_part('isodow', now())`,
		},
		{
			ISOYearField,
			`date_part('isoyear', now())`,
		},
		{
			MicrosecondsField,
			`date_part('microseconds', now())`,
		},
		{
			MillenniumField,
			`date_part('millennium', now())`,
		},
		{
			MillisecondsField,
			`date_part('milliseconds', now())`,
		},
		{
			MinuteField,
			`date_part('minute', now())`,
		},
		{
			MonthField,
			`date_part('month', now())`,
		},
		{
			QuarterField,
			`date_part('quarter', now())`,
		},
		{
			SecondField,
			`date_part('second', now())`,
		},
		{
			TimeZoneField,
			`date_part('timezone', now())`,
		},
		{
			TimeZoneHourField,
			`date_part('timezone_hour', now())`,
		},
		{
			TimeZoneMinuteField,
			`date_part('timezone_minute', now())`,
		},
		{
			WeekField,
			`date_part('week', now())`,
		},
		{
			YearField,
			`date_part('year', now())`,
		},
	}

	for i, tc := range cases {
		sql := DatePart(tc.field, Now()).ToSQLExpr(nil)
		if sql != tc.sql {
			t.Errorf("text case %d: expected %q, got %q", i+1, tc.sql, sql)
		}
	}
}

func TestDateTrunc(t *testing.T) {
	cases := []struct {
		precision DatePrecision
		sql       string
	}{
		{
			MicrosecondsPrecision,
			`date_trunc('microseconds', now())`,
		},
		{
			MillisecondsPrecision,
			`date_trunc('milliseconds', now())`,
		},
		{
			SecondPrecision,
			`date_trunc('second', now())`,
		},
		{
			MinutePrecision,
			`date_trunc('minute', now())`,
		},
		{
			HourPrecision,
			`date_trunc('hour', now())`,
		},
		{
			DayPrecision,
			`date_trunc('day', now())`,
		},
		{
			WeekPrecision,
			`date_trunc('week', now())`,
		},
		{
			MonthPrecision,
			`date_trunc('month', now())`,
		},
		{
			QuarterPrecision,
			`date_trunc('quarter', now())`,
		},
		{
			YearPrecision,
			`date_trunc('year', now())`,
		},
		{
			DecadePrecision,
			`date_trunc('decade', now())`,
		},
		{
			CenturyPrecision,
			`date_trunc('century', now())`,
		},
		{
			MillenniumPrecision,
			`date_trunc('millennium', now())`,
		},
	}

	for i, tc := range cases {
		sql := DateTrunc(tc.precision, Now()).ToSQLExpr(nil)
		if sql != tc.sql {
			t.Errorf("text case %d: expected %q, got %q", i+1, tc.sql, sql)
		}
	}
}
