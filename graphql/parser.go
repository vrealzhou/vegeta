package graphql

import (
	"strings"

	"github.com/tidwall/gjson"
)

type GQLError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	SourceErr bool   `json:"source_error"`
}

type ResultParser struct {
	paths      []string
	results    []gjson.Result
	errorCode  string
	fromSource string
	urlPath    string
}

func NewParser(urlPath, totalFetchCount, errorCode, fromSource string) *ResultParser {
	return &ResultParser{
		urlPath: urlPath,
		paths: []string{
			totalFetchCount,
			"errors",
		},
		errorCode:  errorCode,
		fromSource: fromSource,
	}
}

func (p *ResultParser) IsGraphQL(url string) bool {
	return strings.Contains(url, p.urlPath)
}

func (p *ResultParser) CheckResult(result []byte) {
	p.results = gjson.GetMany(string(result), p.paths...)
}

func (p *ResultParser) TotalFetchCount() uint64 {
	if p.results == nil {
		return 0
	}
	return uint64(p.results[0].Int())
}

func (p *ResultParser) ParseErrors() []GQLError {
	gqlErrors := make([]GQLError, 0)
	if p.results == nil {
		return gqlErrors
	}
	errs := p.results[1].Array()
	for _, err := range errs {
		errMap := err.Map()
		extensions := errMap["extensions"].Map()
		gqlErr := GQLError{
			Message: errMap["message"].String(),
		}
		if code, ok := extensions[p.errorCode]; ok {
			gqlErr.Code = code.String()
		}
		if fromSource, ok := extensions[p.fromSource]; ok {
			gqlErr.SourceErr = fromSource.Bool()
		}
		gqlErrors = append(gqlErrors, gqlErr)
	}
	return gqlErrors
}
