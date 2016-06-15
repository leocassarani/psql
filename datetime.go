package psql

import "fmt"

// Now returns an Expression representing a call to the 0-arity date/time function now().
func Now() fnCall {
	return fnCall{"now"}
}

type fnCall struct {
	name string
}

func (f fnCall) ToSQLExpr(*Params) string {
	return fmt.Sprintf("%s()", f.name)
}

func (f fnCall) Relations() []string {
	return nil
}

// DatePart returns an Expression representing a call to the date/time
// function date_part(), which extracts the given field from expr.
func DatePart(field DateField, expr Expression) datePart {
	return datePart{field, expr}
}

type datePart struct {
	field DateField
	expr  Expression
}

func (d datePart) ToSQLExpr(p *Params) string {
	return fmt.Sprintf("date_part('%s', %s)", d.field, d.expr.ToSQLExpr(p))
}

func (d datePart) Relations() []string {
	return d.expr.Relations()
}

type DateField int

const (
	CenturyField DateField = iota
	DayField
	DecadeField
	DayOfWeekField
	DayOfYearField
	EpochField
	HourField
	ISODayOfWeekField
	ISOYearField
	MicrosecondsField
	MillenniumField
	MillisecondsField
	MinuteField
	MonthField
	QuarterField
	SecondField
	TimeZoneField
	TimeZoneHourField
	TimeZoneMinuteField
	WeekField
	YearField
)

func (d DateField) String() string {
	switch d {
	case CenturyField:
		return "century"
	case DayField:
		return "day"
	case DecadeField:
		return "decade"
	case DayOfWeekField:
		return "dow"
	case DayOfYearField:
		return "doy"
	case EpochField:
		return "epoch"
	case HourField:
		return "hour"
	case ISODayOfWeekField:
		return "isodow"
	case ISOYearField:
		return "isoyear"
	case MicrosecondsField:
		return "microseconds"
	case MillenniumField:
		return "millennium"
	case MillisecondsField:
		return "milliseconds"
	case MinuteField:
		return "minute"
	case MonthField:
		return "month"
	case QuarterField:
		return "quarter"
	case SecondField:
		return "second"
	case TimeZoneField:
		return "timezone"
	case TimeZoneHourField:
		return "timezone_hour"
	case TimeZoneMinuteField:
		return "timezone_minute"
	case WeekField:
		return "week"
	case YearField:
		return "year"
	default:
		panic("unknown DateField")
	}
}

// DateTrunc returns an Expression representing a call to the date/time
// function date_trunc(), which truncates expr to the given precision.
func DateTrunc(precision DatePrecision, expr Expression) dateTrunc {
	return dateTrunc{precision, expr}
}

type dateTrunc struct {
	precision DatePrecision
	expr      Expression
}

func (d dateTrunc) ToSQLExpr(p *Params) string {
	return fmt.Sprintf("date_trunc('%s', %s)", d.precision, d.expr.ToSQLExpr(p))
}

func (d dateTrunc) Relations() []string {
	return d.expr.Relations()
}

type DatePrecision int

const (
	MicrosecondsPrecision DatePrecision = iota
	MillisecondsPrecision
	SecondPrecision
	MinutePrecision
	HourPrecision
	DayPrecision
	WeekPrecision
	MonthPrecision
	QuarterPrecision
	YearPrecision
	DecadePrecision
	CenturyPrecision
	MillenniumPrecision
)

func (d DatePrecision) String() string {
	switch d {
	case MicrosecondsPrecision:
		return "microseconds"
	case MillisecondsPrecision:
		return "milliseconds"
	case SecondPrecision:
		return "second"
	case MinutePrecision:
		return "minute"
	case HourPrecision:
		return "hour"
	case DayPrecision:
		return "day"
	case WeekPrecision:
		return "week"
	case MonthPrecision:
		return "month"
	case QuarterPrecision:
		return "quarter"
	case YearPrecision:
		return "year"
	case DecadePrecision:
		return "decade"
	case CenturyPrecision:
		return "century"
	case MillenniumPrecision:
		return "millennium"
	default:
		panic("unknown DatePrecision")
	}
}
