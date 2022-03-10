package dqlx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type funcType string

var (
	eq funcType = "eq" // Done
	// ie
	ge funcType = "ge" // Done
	gt funcType = "gt" // Done
	le funcType = "le" // Done
	lt funcType = "lt" // Done

	has funcType = "has"  // Done
	typ funcType = "type" // Done
	// term
	allofterms funcType = "allofterms" // Done
	anyofterms funcType = "anyofterms" // Done
	// trigram
	regexp funcType = "regexp" // Done
	match  funcType = "match"  // Done
	// fulltext
	alloftext funcType = "alloftext" // Done
	anyoftext funcType = "anyoftext" // Done

	between funcType = "between" // Done
	uid     funcType = "uid"     // Done
	uidIn   funcType = "uid_in"  // Done
)

type Filter interface {
	dqlizer
	Type() funcType
}

type eqExpr struct {
	key    interface{}
	values []interface{}
}

func (eqExpr *eqExpr) toDQL() (query string, args []interface{}, err error) {
	if len(eqExpr.values) == 0 {
		return "", nil, fmt.Errorf("eq should have one value at least")
	}

	var key string
	switch keyCast := eqExpr.key.(type) {
	case string:
		key = escapePredicate(keyCast)
	case *valExpr:
		key = keyCast.String()
	case *countExpr:
		key = keyCast.String()
	default:
		return "", nil, fmt.Errorf("eq accepts only string, val() or count() as value, given %T", keyCast)
	}

	placeholders := make([]string, len(eqExpr.values))
	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = symbolValuePlaceholder
	}

	query = string(eqExpr.Type()) + "(" + key + "," +
		strings.Join(placeholders, ",") + ")"

	return query, eqExpr.values, nil
}

func (eqExpr *eqExpr) Type() funcType {
	return eq
}

// Eq represents the eq expression,
// Expression: eq(key, value, [value1, value2, ..., valueN])
// Expression: eq(val(varName),value), eq(count(predicate), value)
func Eq(key interface{}, values ...interface{}) *eqExpr {
	return &eqExpr{
		key:    key,
		values: values,
	}
}

type ie struct {
	ieType funcType
	key    interface{}
	value  interface{}
}

func (ie *ie) toDQL() (query string, args []interface{}, err error) {
	var key string
	switch keyCast := ie.key.(type) {
	case string:
		key = escapePredicate(keyCast)
	case *valExpr:
		key = keyCast.String()
	case *countExpr:
		key = keyCast.String()
	default:
		return "", nil, fmt.Errorf("ie accepts only string, val() or count() as value, given %T", keyCast)
	}

	query = string(ie.Type()) + "(" + key + "," +
		symbolValuePlaceholder + ")"

	return query, []interface{}{ie.value}, nil
}

func (ie *ie) Type() funcType {
	return ie.ieType
}

// Le represents the le expression,
// Expression: le(key, value)
func Le(key interface{}, value interface{}) *ie {
	return &ie{
		ieType: le,
		key:    key,
		value:  value,
	}
}

// Lt represents the lt expression,
// Expression: lt(name, value)
func Lt(key interface{}, value interface{}) *ie {
	return &ie{
		ieType: lt,
		key:    key,
		value:  value,
	}
}

// Ge represents the ge expression,
// Expression: ge(name, value)
func Ge(key interface{}, value interface{}) *ie {
	return &ie{
		ieType: ge,
		key:    key,
		value:  value,
	}
}

// Gt represents the gt expression,
// Expression: gt(name, value)
func Gt(key interface{}, value interface{}) *ie {
	return &ie{
		ieType: gt,
		key:    key,
		value:  value,
	}
}

type hasExpr struct {
	predicate string
}

func (hasExpr *hasExpr) toDQL() (query string, args []interface{}, err error) {
	query = string(hasExpr.Type()) + "(" + escapePredicate(hasExpr.predicate) + ")"

	return query, args, nil
}

func (hasExpr *hasExpr) Type() funcType {
	return has
}

// Has represents the has expression,
// Expression: has(predicate)
func Has(predicate string) *hasExpr {
	return &hasExpr{
		predicate: predicate,
	}
}

type typeExpr struct {
	dgraphType string
}

func (typeExpr *typeExpr) toDQL() (query string, args []interface{}, err error) {
	query = string(typeExpr.Type()) + "(" + symbolValuePlaceholder + ")"

	return query, []interface{}{typeExpr.dgraphType}, nil
}

func (typeExpr *typeExpr) Type() funcType {
	return typ
}

// Type represents the type expression,
// Expression: type(dgraphType)
func Type(dgraphType string) *typeExpr {
	return &typeExpr{
		dgraphType: dgraphType,
	}
}

type term struct {
	termType  funcType
	predicate string
	terms     string
}

func (term *term) toDQL() (query string, args []interface{}, err error) {
	query = string(term.Type()) + "(" + escapePredicate(term.predicate) +
		"," + symbolValuePlaceholder + ")"

	return query, []interface{}{term.terms}, nil
}

func (term *term) Type() funcType {
	return term.termType
}

// AllOfTerms represents the allofterms expression,
// Expression: allofterms(predicate, terms)
func AllOfTerms(predicate string, terms string) *term {
	return &term{
		termType:  allofterms,
		predicate: predicate,
		terms:     terms,
	}
}

// AnyOfTerms represents the anyofterms expression,
// Expression: anyofterms(predicate, terms)
func AnyOfTerms(predicate string, terms string) *term {
	return &term{
		termType:  anyofterms,
		predicate: predicate,
		terms:     terms,
	}
}

type trigram struct {
	trigramType funcType
	predicate   string
	pattern     string
}

func (trigram *trigram) toDQL() (query string, args []interface{}, err error) {
	query = string(trigram.Type()) + "(" + escapePredicate(trigram.predicate) +
		"," + symbolValuePlaceholder + ")"

	return query, []interface{}{trigram.pattern}, nil
}

func (trigram *trigram) Type() funcType {
	return trigram.trigramType
}

type regexpExpr struct {
	trigram
}

// Regexp represents the regexp expression,
// Expression: regexp(predicate, /pattern/)
func Regexp(predicate string, pattern string) *regexpExpr {
	return &regexpExpr{trigram{
		trigramType: regexp,
		predicate:   predicate,
		pattern:     pattern,
	}}
}

type matchExpr struct {
	trigram
	distance int
}

func (matchExpr *matchExpr) toDQL() (query string, args []interface{}, err error) {
	placeholders := make([]string, 2)
	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = symbolValuePlaceholder
	}
	query = string(matchExpr.Type()) + "(" + escapePredicate(matchExpr.predicate) + "," +
		strings.Join(placeholders, ",") + ")"

	return query, []interface{}{matchExpr.pattern, matchExpr.distance}, nil
}

// Match represents the match expression,
// Expression: match(predicate, pattern, distance)
func Match(predicate string, pattern string, distance int) *matchExpr {
	return &matchExpr{
		trigram: trigram{
			trigramType: match,
			predicate:   predicate,
			pattern:     pattern,
		},
		distance: distance,
	}
}

type fulltext struct {
	fulltextType funcType
	predicate    string
	text         string
}

func (fulltext *fulltext) toDQL() (query string, args []interface{}, err error) {
	query = string(fulltext.Type()) + "(" + escapePredicate(fulltext.predicate) + "," +
		symbolValuePlaceholder + ")"

	return query, []interface{}{fulltext.text}, nil
}

func (fulltext *fulltext) Type() funcType {
	return fulltext.fulltextType
}

// AllOfText represents the match expression,
// Expression: alloftext(predicate, text)
func AllOfText(predicate string, text string) *fulltext {
	return &fulltext{
		fulltextType: alloftext,
		predicate:    predicate,
		text:         text,
	}
}

// AnyOfText represents the match expression,
// Expression: anyoftext(predicate, text)
func AnyOfText(predicate string, text string) *fulltext {
	return &fulltext{
		fulltextType: anyoftext,
		predicate:    predicate,
		text:         text,
	}
}

type betweenExpr struct {
	predicate  string
	lowerbound interface{}
	upperbound interface{}
}

func (betweenExpr *betweenExpr) toDQL() (query string, args []interface{}, err error) {
	placeholders := make([]string, 2)
	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = symbolValuePlaceholder
	}
	query = string(betweenExpr.Type()) + "(" + escapePredicate(betweenExpr.predicate) + "," +
		strings.Join(placeholders, ",") + ")"

	return query, []interface{}{betweenExpr.lowerbound, betweenExpr.upperbound}, nil
}

func (betweenExpr *betweenExpr) Type() funcType {
	return between
}

// Between represents the between expression,
// Expression: between(predicate, lowerbound, upperbound)
func Between(predicate string, lowerbound interface{}, upperbound interface{}) *betweenExpr {
	return &betweenExpr{
		predicate:  predicate,
		lowerbound: lowerbound,
		upperbound: upperbound,
	}
}

type uidExpr struct {
	values []string
}

func (uidExpr *uidExpr) toDQL() (query string, args []interface{}, err error) {
	placeholders := make([]string, len(uidExpr.values))

	for i := 0; i < len(placeholders); i++ {
		// Suppose that variables begin with uppercase letters, and use this strange method for the time being
		isStartUpper := unicode.IsUpper([]rune(uidExpr.values[i])[0])
		if isStartUpper {
			placeholders[i] = escapePredicate(uidExpr.values[i])
		} else {
			placeholders[i] = uidExpr.values[i]
		}
	}
	query = string(uidExpr.Type()) + "(" + strings.Join(placeholders, ",") + ")"

	return query, args, nil
}

func (uidExpr *uidExpr) Type() funcType {
	return uid
}

// Uid returns uid expression
// Expression: uid(value1, [value2, value3...])
func Uid(values ...string) *uidExpr {
	return &uidExpr{
		values: values,
	}
}

type uidInExpr struct {
	predicate string
	value     interface{}
}

func (uidInExpr *uidInExpr) toDQL() (query string, args []interface{}, err error) {
	valVal := reflect.ValueOf(uidInExpr.value)
	isListType := valVal.Kind() == reflect.Array || valVal.Kind() == reflect.Slice

	if isListType {
		var listValue []interface{}

		listValue, err = toInterfaceSlice(uidInExpr.value)

		if err != nil {
			return "", nil, err
		}

		placeholders := make([]string, len(listValue))
		for index, value := range listValue {
			placeholders[index] = symbolValuePlaceholder
			args = append(args, value)
		}
		query = string(uidInExpr.Type()) + "(" + escapePredicate(uidInExpr.predicate) + ",[" +
			strings.Join(placeholders, ",") + "])"

		return
	}

	query = string(uidInExpr.Type()) + "(" + escapePredicate(uidInExpr.predicate) + "," +
		symbolValuePlaceholder + ")"

	return query, []interface{}{uidInExpr.value}, nil
}

func (uidInExpr *uidInExpr) Type() funcType {
	return uidIn
}

// UidIn represents the uid_in expression,
// Expression: uid_in(predicate, value) or uid_in(predicate, [value1, value2...])
func UidIn(predicate string, value interface{}) *uidInExpr {
	return &uidInExpr{
		predicate: predicate,
		value:     value,
	}
}

func toInterfaceSlice(slice interface{}) ([]interface{}, error) {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil, errors.New("toInterfaceSlice given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil, nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret, nil
}

type conjunction struct {
	filters   []Filter
	connector string
}

func (conjunction *conjunction) toDQL() (query string, args []interface{}, err error) {
	if len(conjunction.filters) == 0 {
		return "", args, nil
	}

	var filters []string
	for _, part := range conjunction.filters {
		partDql, partArgs, err := part.toDQL()
		if err != nil {
			return "", nil, err
		}
		if partDql != "" {
			filters = append(filters, partDql)
			args = append(args, partArgs...)
		}
	}

	filterNum := len(filters)
	if filterNum == 1 {
		query = filters[0]
	} else if filterNum > 1 {
		query = "(" + strings.Join(filters, conjunction.connector) + ")"
	}

	return query, args, nil
}
func (conjunction *conjunction) Type() funcType {
	return "conjunction"
}

func And(filters ...Filter) *conjunction {
	return &conjunction{
		filters:   filters,
		connector: " AND ",
	}
}

func Or(filters ...Filter) *conjunction {
	return &conjunction{
		filters:   filters,
		connector: " OR ",
	}
}

type notExpr struct {
	filter Filter
}

func (notExpr *notExpr) toDQL() (query string, args []interface{}, err error) {
	filterDql, filterArgs, err := notExpr.filter.toDQL()
	if err != nil {
		return "", nil, err
	}

	query = "NOT " + filterDql

	return query, filterArgs, nil
}

func (notExpr *notExpr) Type() funcType {
	return "not"
}

func Not(filter Filter) *notExpr {
	return &notExpr{filter: filter}
}
