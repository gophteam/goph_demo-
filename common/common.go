package common

import (
	"io/ioutil"
	"math"
	"reflect"
	"regexp"
	"strings"
	"time"

	"encoding/json"
	"fmt"
	"net/http"
	// gophCommon "github.com/gophteam/goph/common"
)

var EXPOSED_VAR string = "exposed"

// EmailValidity !
type EmailValidity struct {
	Email   string `json:"email"`
	Valid   bool   `json:"valid"`
	Exist   bool   `json:"exist"`
	Remarks string `json:"remarks"`
}

// ValidateEmail !
func ValidateEmail(email string) (EmailValidity, error) {
	// https://files.gopilipinas.org/mailvalidate/goph_validate_email.php?token=j0intherev0luti0n&email=jerickgonito@gmail.com

	email = strings.TrimSpace(email)
	emailValidity := EmailValidity{}
	emailValidity.Email = email
	switch true {
	case len(email) <= 5: // x@x.x
		fallthrough
	case !strings.Contains(email, "@"):
		fallthrough
	case !strings.Contains(email, "."):
		emailValidity.Valid = false
		emailValidity.Exist = false
		emailValidity.Remarks = "Invalid email format"
		return emailValidity, fmt.Errorf(emailValidity.Remarks)
	}

	token := "j0intherev0luti0n"
	endPoint := fmt.Sprintf("https://files.gopilipinas.org/mailvalidate/goph_validate_email.php?token=%s&email=%s", token, email)

	res, err := http.Post(endPoint, "application/json", nil)
	if err != nil {
		return emailValidity, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(body, &emailValidity); err != nil {
		return emailValidity, err
	}

	return emailValidity, nil
}

// IsValidEmail !
func IsValidEmail(email string) bool {
	emailValidity, err := ValidateEmail(email)
	if err != nil {
		// handle error
		return false
	}
	return emailValidity.Valid && emailValidity.Exist
}

// GetInclusiveYM !
func GetInclusiveYM(y1 int, y2 int, m1 int, m2 int) map[int][]int {
	yms := map[int][]int{}

	y := y1
	m := m1

	for {
		if _, ok := yms[y]; !ok {
			yms[y] = []int{}
		}
		yms[y] = append(yms[y], m)

		m++
		if m > 12 {
			m = 1
			y++
		}

		if y > y2 || (y == y2 && m > m2) {
			break
		}
	}

	return yms
}

// GetInclusiveDates !
func GetInclusiveDates(dateFrom time.Time, dateTo time.Time, groupBy string) []string {
	inclusives := []string{}
	// isSameYear := dateFrom.Format("2006") == dateTo.Format("2006")
	dt := dateFrom

	switch groupBy {
	case "d":
		// format := "Jan 02"
		// if !isSameYear {
		// 	format += ", 2006"
		// }

		for {
			inclusives = append(inclusives, dt.Format("2006-01-02"))
			// inclusives[dt.Format("2006-01-02")] = dt.Format(format)
			dt = dt.AddDate(0, 0, 1)
			if dt.Unix() > dateTo.Unix() {
				break
			}
		}
		break
	case "m":
		// format := "Jan"
		// if !isSameYear {
		// 	format += ", 2006"
		// }

		dt, _ = time.Parse("2006-01-02", dt.Format("2006-01")+"-01")
		for {
			inclusives = append(inclusives, dt.Format("2006-01"))
			// inclusives[dt.Format("2006-01")] = dt.Format(format)
			dt = dt.AddDate(0, 1, 0)
			if dt.Unix() > dateTo.Unix() {
				break
			}
		}
		break
	case "y":

		break
	}

	return inclusives
}

// StructToMap !
func StructToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = StructToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}

// SanitizeDate !
func SanitizeDate(iDate interface{}) *time.Time {
	date, ok := iDate.(time.Time)
	if !ok {
		if dt, ok := iDate.(*time.Time); ok && dt != nil {
			date = *dt
		} else {
			return nil
		}
	}

	newDate := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), time.Local)
	return &newDate
}

// GetAge !
func GetAge(birthdate *time.Time) int {
	if birthdate == nil {
		return 0
	}

	now := time.Now()
	return int(math.Floor(now.Sub(*birthdate).Hours() / 24 / 365))
}

// ValidatePasswordStrength !
func ValidatePasswordStrength(password string) bool {
	// check for lower case charter
	if matched, _ := regexp.MatchString(`[a-z]`, password); !matched {
		return false
	}

	// check for upper case character
	if matched, _ := regexp.MatchString(`[A-Z]`, password); !matched {
		return false
	}

	// check for digits character
	if matched, _ := regexp.MatchString(`[\d]`, password); !matched {
		return false
	}

	// check for non-letter and non-digit character
	if matched, _ := regexp.MatchString(`[\W]`, password); !matched {
		return false
	}

	if len(password) < 8 {
		return false
	}

	return true
}

// ToSnakeCase !
func ToSnakeCase(str string) string {
	// return gophCommon.ToSnakeCase(str)

	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	v := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	v = matchAllCap.ReplaceAllString(v, "${1}_${2}")
	return strings.ToLower(v)
}

// GeneratePages !
func GeneratePages(totalRecord int, recordPerPage int) []map[string]int {
	pages := []map[string]int{}
	totalPage := math.Ceil(float64(totalRecord) / float64(recordPerPage))
	for i := 0; i < int(totalPage); i++ {
		pages = append(pages, map[string]int{
			"page":   i + 1,
			"offset": i * recordPerPage,
		})
	}
	return pages
}

// EnumerateDates : Returns an array of dates base on passed paremeters, dateTo can be nil and days can be negative.
func EnumerateDates(dateFrom time.Time, dateTo *time.Time, days int) []time.Time {
	dateFrom = time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, dateFrom.Location())
	if dateTo != nil {
		dtValue := *dateTo
		tmpDT := time.Date(dtValue.Year(), dtValue.Month(), dtValue.Day(), 0, 0, 0, 0, dtValue.Location())
		dateTo = &tmpDT
	} else {
		if days == 0 {
			now := time.Now()
			tmpDT := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			dateTo = &tmpDT
		} else {
			if days > 0 {
				tmpDT := dateFrom.AddDate(0, 0, days)
				dateTo = &tmpDT
			} else {
				tmpDF := dateFrom.AddDate(0, 0, days)
				dateTo = &dateFrom
				dateFrom = tmpDF
			}
		}
	}

	dates := []time.Time{}

	tmpDate := dateFrom
	for !tmpDate.After(*dateTo) {
		dates = append(dates, tmpDate)
		tmpDate = tmpDate.AddDate(0, 0, 1)
	}
	return dates
}

// EnumerateWeeks : Returns an array of 2 dates base on pass paremeters, dateTo can be nil and week can be negative.
func EnumerateWeeks(dateFrom time.Time, dateTo *time.Time, weeks int) [][]time.Time {
	days := 0
	if weeks != 0 {
		days = 7 * weeks
	}

	dates := EnumerateDates(dateFrom, dateTo, days)
	resWeeks := [][]time.Time{}
	weekRange := []time.Time{}
	for i, date := range dates {
		if len(weekRange) == 0 {
			weekRange = append(weekRange, date)
		}

		// Saturday : last day of the week
		if date.Weekday() == time.Saturday || i == (len(dates)-1) {
			weekRange = append(weekRange, time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location()).AddDate(0, 0, 1).Add(-time.Second))
			resWeeks = append(resWeeks, weekRange)
			weekRange = []time.Time{}
		}
	}

	return resWeeks
}

// EnumerateMonths : Returns an array of 2 dates base on pass paremeters, dateTo can be nil and months can be negative.
func EnumerateMonths(dateFrom time.Time, dateTo *time.Time, months int) [][]time.Time {
	dateFrom = time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, dateFrom.Location())
	if dateTo != nil {
		dtValue := *dateTo
		tmpDT := time.Date(dtValue.Year(), dtValue.Month(), dtValue.Day(), 0, 0, 0, 0, dtValue.Location())
		dateTo = &tmpDT
	} else {
		if months == 0 {
			now := time.Now()
			tmpDT := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
			dateTo = &tmpDT
		} else {
			if months > 0 {
				tmpDT := dateFrom.AddDate(0, months, 0)
				dateTo = &tmpDT
			} else {
				tmpDF := dateFrom.AddDate(0, months, 0)
				dateTo = &dateFrom
				dateFrom = tmpDF
			}
		}
	}

	tmpDF := dateFrom
	resMonths := [][]time.Time{}
	monthRange := []time.Time{}
	monthStart := dateFrom

	for !tmpDF.After(*dateTo) {
		if len(monthRange) == 0 {
			monthRange = append(monthRange, monthStart)
		}

		if tmpDF.Month() == dateTo.Month() && tmpDF.Year() == dateTo.Year() {
			tmpDate := *dateTo
			monthRange = append(monthRange, time.Date(tmpDate.Year(), tmpDate.Month(), tmpDate.Day(), 0, 0, 0, 0, tmpDate.Location()).AddDate(0, 0, 1).Add(-time.Second))
		} else {
			monthRange = append(monthRange, time.Date(monthStart.Year(), (monthStart.Month()+1), 1, 0, 0, 0, 0, monthStart.Location()).Add(-time.Second))
		}

		resMonths = append(resMonths, monthRange)
		monthStart = time.Date(monthStart.Year(), (monthStart.Month() + 1), 1, 0, 0, 0, 0, monthStart.Location())
		monthRange = []time.Time{}
		tmpDF = tmpDF.AddDate(0, 1, 0)
	}

	return resMonths
}

// EnumerateYears : Returns an array of 2 dates base on pass paremeters, yearTo can be nil and years can be negative.
func EnumerateYears(yearFrom, yearTo, years int) [][]time.Time {
	resYears := [][]time.Time{}
	yearDif := yearTo - yearFrom
	if yearDif == 0 {
		tmpDF := time.Date(yearFrom, time.Month(1), 1, 0, 0, 0, 0, time.Now().Location())
		tmpDT := time.Date((yearFrom + 1), time.Month(1), 1, 0, 0, 0, 0, time.Now().Location()).Add(-time.Second)
		resYears = append(resYears, []time.Time{tmpDF, tmpDT})
	} else {
		for i := 0; i <= yearDif; i++ {
			tmpYear := (yearFrom + i)
			tmpDF := time.Date(tmpYear, time.Month(1), 1, 0, 0, 0, 0, time.Now().Location())
			tmpDT := time.Date(tmpYear+1, time.Month(1), 1, 0, 0, 0, 0, time.Now().Location()).Add(-time.Second)
			resYears = append(resYears, []time.Time{tmpDF, tmpDT})
		}
	}

	return resYears
}

// ReturnJSON returns json to the client
func ReturnJSON(w http.ResponseWriter, r *http.Request, v interface{}, minified bool) {
	// b, err := json.MarshalIndent(v, "", "  ")
	var b []byte
	var err error
	if minified {
		b, _ = json.Marshal(v)
	} else {
		b, _ = json.MarshalIndent(v, "", "  ")
	}
	if err != nil {
		response := map[string]interface{}{
			"status":    "error",
			"error_msg": fmt.Sprintf("unable to encode JSON. %s", err),
		}
		if minified {
			b, _ = json.Marshal(response)
		} else {
			b, _ = json.MarshalIndent(response, "", "  ")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
		return
	}
	w.Write(b)
}

// GroupRecord !
type GroupRecord struct {
	Group  string      `json:"group"`
	Values interface{} `json:"values"`
}

// GroupFieldKeys !
type GroupFieldKeys struct {
	Field string
	Keys  []string
}

// GroupRecords !
func GroupRecords(records interface{}, groupFieldKeys []GroupFieldKeys, outputForm interface{}) ([]GroupRecord, error) {
	fields := []string{}
	for _, groupFieldKey := range groupFieldKeys {
		if field := groupFieldKey.Field; field != "" {
			fields = append(fields, field)
		}
	}

	groupRecords := []GroupRecord{}
	if len(fields) == 0 {
		return groupRecords, fmt.Errorf("no fields found")
	}

	rRecords := reflect.ValueOf(records)
	if rRecords.Kind() != reflect.Slice {
		return groupRecords, fmt.Errorf("'records' must be slice")
	}

	rValueType := reflect.TypeOf(outputForm)
	rValuesType := reflect.New(reflect.SliceOf(rValueType)).Elem().Type()

	mainGroup := GroupRecord{Values: []GroupRecord{}}
	keyIndexes := map[string]int{}

	for i := 0; i < rRecords.Len(); i++ {
		inProcessGroup := &mainGroup

		rRecord := rRecords.Index(i)
		keyIndex := ""
		for _, field := range fields {
			iKey := rRecord.FieldByName(field).Interface()
			key := fmt.Sprint(iKey)

			if keyIndex != "" {
				keyIndex += "__"
			}
			keyIndex += key

			i, ok := keyIndexes[keyIndex]
			if !ok {
				values := inProcessGroup.Values.([]GroupRecord)
				values = append(values, GroupRecord{Group: key, Values: []GroupRecord{}})
				inProcessGroup.Values = values
				i = len(values) - 1
				keyIndexes[keyIndex] = i
			}

			inProcessGroup = &inProcessGroup.Values.([]GroupRecord)[i]
		}

		if reflect.TypeOf(inProcessGroup.Values) != rValuesType {
			inProcessGroup.Values = reflect.New(rValuesType).Elem().Interface()
		}
		rValues := reflect.ValueOf(inProcessGroup.Values)

		rValue := reflect.New(rValueType).Elem()
		if rValue.Type() == rRecord.Type() {
			rValue = rRecord
		} else {
			for ii := 0; ii < rValue.NumField(); ii++ {
				fName := rValue.Type().Field(ii).Name
				if rFieldValue := rRecord.FieldByName(fName); rFieldValue.IsValid() {
					rValue.Field(ii).Set(rFieldValue)
				}
			}
		}

		rValues = reflect.Append(rValues, rValue)
		inProcessGroup.Values = rValues.Interface()
	}

	return mainGroup.Values.([]GroupRecord), nil
}

// EncryptKey :
func EncryptKey(key string) string {
	return ""
}

// DecryptKey :
func DecryptKey(key string) string {
	return ""
}
