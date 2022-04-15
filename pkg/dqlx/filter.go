package dqlx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type FuncType string

var (
	eq FuncType = "eq" // Done
	// ie
	ge FuncType = "ge" // Done
	gt FuncType = "gt" // Done
	le FuncType = "le" // Done
	lt FuncType = "lt" // Done

	has FuncType = "has"  // Done
	typ FuncType = "type" // Done
	// term
	allofterms FuncType = "allofterms" // Done
	anyofterms FuncType = "anyofterms" // Done
	// trigram
	regexp FuncType = "regexp" // Done
	match  FuncType = "match"  // Done
	// fulltext
	alloftext FuncType = "alloftext" // Done
	anyoftext FuncType = "anyoftext" // Done

	between FuncType = "between" // Done
	uid     FuncType = "uid"     // Done
	uidIn   FuncType = "uid_in"  // Done
)

type Filter interface {
	Dqlizer
	Type() FuncType
}

type eqExpr struct {
	key    interface{}
	values []interface{}
}

func (eqExpr *eqExpr) Dql() (string, []interface{}, error) {
	if len(eqExpr.values) == 0 {
		return "", nil, fmt.Errorf("eq should have one value at least")
	}

	var key string
	switch cast := eqExpr.key.(type) {
	case string:
		key = Escape(cast)
	case *ValExpr:
		key = cast.String()
	case *CountExpr:
		key = cast.String()
	default:
		return "", nil, fmt.Errorf("eq accepts only string, val() or count() as value, given %T", cast)
	}

	placeholders := make([]string, len(eqExpr.values))
	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = symbolValuePlaceholder
	}

	return string(eqExpr.Type()) + "(" + key + "," + strings.Join(placeholders, ",") + ")", eqExpr.values, nil
}

func (eqExpr *eqExpr) Type() FuncType {
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
	ieType FuncType
	key    interface{}
	value  interface{}
}

func (ie *ie) Dql() (string, []interface{}, error) {
	var key string
	switch cast := ie.key.(type) {
	case string:
		key = Escape(cast)
	case *ValExpr:
		key = cast.String()
	case *CountExpr:
		key = cast.String()
	default:
		return "", nil, fmt.Errorf("ie accepts only string, val() or count() as value, given %T", cast)
	}

	return string(ie.Type()) + "(" + key + "," + symbolValuePlaceholder + ")", []interface{}{ie.value}, nil
}

func (ie *ie) Type() FuncType {
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

func (hasExpr *hasExpr) Dql() (string, []interface{}, error) {
	return string(hasExpr.Type()) + "(" + Escape(hasExpr.predicate) + ")", nil, nil
}

func (hasExpr *hasExpr) Type() FuncType {
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

func (typeExpr *typeExpr) Dql() (string, []interface{}, error) {
	return string(typeExpr.Type()) + "(" + symbolValuePlaceholder + ")", []interface{}{typeExpr.dgraphType}, nil
}

func (typeExpr *typeExpr) Type() FuncType {
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
	termType  FuncType
	predicate string
	terms     string
}

func (term *term) Dql() (string, []interface{}, error) {
	return string(term.Type()) + "(" + Escape(term.predicate) +
		"," + symbolValuePlaceholder + ")", []interface{}{term.terms}, nil
}

func (term *term) Type() FuncType {
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
	trigramType FuncType
	predicate   string
	pattern     string
}

func (trigram *trigram) Dql() (string, []interface{}, error) {
	return string(trigram.Type()) + "(" + Escape(trigram.predicate) +
		"," + symbolValuePlaceholder + ")", []interface{}{trigram.pattern}, nil
}

func (trigram *trigram) Type() FuncType {
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

func (matchExpr *matchExpr) Dql() (string, []interface{}, error) {
	placeholders := make([]string, 2)
	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = symbolValuePlaceholder
	}

	return string(matchExpr.Type()) + "(" + Escape(matchExpr.predicate) + "," +
		strings.Join(placeholders, ",") + ")", []interface{}{matchExpr.pattern, matchExpr.distance}, nil
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
	fulltextType FuncType
	predicate    string
	text         string
}

func (fulltext *fulltext) Dql() (string, []interface{}, error) {
	return string(fulltext.Type()) + "(" + Escape(fulltext.predicate) + "," +
		symbolValuePlaceholder + ")", []interface{}{fulltext.text}, nil
}

func (fulltext *fulltext) Type() FuncType {
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

func (betweenExpr *betweenExpr) Dql() (string, []interface{}, error) {
	placeholders := make([]string, 2)
	for i := 0; i < len(placeholders); i++ {
		placeholders[i] = symbolValuePlaceholder
	}

	return string(betweenExpr.Type()) + "(" + Escape(betweenExpr.predicate) + "," +
		strings.Join(placeholders, ",") + ")", []interface{}{betweenExpr.lowerbound, betweenExpr.upperbound}, nil
}

func (betweenExpr *betweenExpr) Type() FuncType {
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

func (uidExpr *uidExpr) Dql() (string, []interface{}, error) {
	placeholders := make([]string, len(uidExpr.values))

	for i := 0; i < len(placeholders); i++ {
		// Suppose that variables begin with uppercase letters
		isStartUpper := unicode.IsUpper([]rune(uidExpr.values[i])[0])
		if isStartUpper {
			placeholders[i] = Escape(uidExpr.values[i])
		} else {
			placeholders[i] = uidExpr.values[i]
		}
	}

	return string(uidExpr.Type()) + "(" + strings.Join(placeholders, ",") + ")", nil, nil
}

func (uidExpr *uidExpr) Type() FuncType {
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

func (uidInExpr *uidInExpr) Dql() (string, []interface{}, error) {
	var (
		args []interface{}
	)

	valVal := reflect.ValueOf(uidInExpr.value)
	if valVal.Kind() == reflect.Array || valVal.Kind() == reflect.Slice {
		listValue, err := toInterfaceSlice(uidInExpr.value)
		if err != nil {
			return "", nil, err
		}

		placeholders := make([]string, len(listValue))
		for index, value := range listValue {
			placeholders[index] = symbolValuePlaceholder
			args = append(args, value)
		}

		return string(uidInExpr.Type()) + "(" + Escape(uidInExpr.predicate) + ",[" +
			strings.Join(placeholders, ",") + "])", args, nil
	}

	return string(uidInExpr.Type()) + "(" + Escape(uidInExpr.predicate) + "," +
		symbolValuePlaceholder + ")", []interface{}{uidInExpr.value}, nil
}

func (uidInExpr *uidInExpr) Type() FuncType {
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

func (conjunction *conjunction) Dql() (string, []interface{}, error) {
	var (
		filters []string
		args    []interface{}
	)
	for _, part := range conjunction.filters {
		partQ, partArgs, err := part.Dql()
		if err != nil {
			return "", nil, err
		}
		if partQ != "" {
			filters = append(filters, partQ)
			args = append(args, partArgs...)
		}
	}

	if len(filters) == 1 {
		return filters[0], args, nil
	} else if len(filters) > 1 {
		return "(" + strings.Join(filters, conjunction.connector) + ")", args, nil
	}

	return "", nil, nil
}
func (conjunction *conjunction) Type() FuncType {
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

func (notExpr *notExpr) Dql() (string, []interface{}, error) {
	q, args, err := notExpr.filter.Dql()
	if err != nil {
		return "", nil, err
	}

	return "NOT " + q, args, nil
}

func (notExpr *notExpr) Type() FuncType {
	return "not"
}

func Not(filter Filter) *notExpr {
	return &notExpr{filter: filter}
}
