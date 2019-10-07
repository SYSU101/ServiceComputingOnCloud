package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type flagValidator struct {
	Condition bool
	ErrorMsg  string
}
type flagMatcher struct {
	Pattern string
	Next    func(*[]string) error
	Usage   string
}

var startPageN, endPageN, linesPerPage *int
var formFeed, showHelpMsg *bool
var toDestination, filename *string

var validators []*flagValidator

var matchers []*flagMatcher

func initMatchers() {
	matchers = []*flagMatcher{
		&flagMatcher{
			Pattern: "^-s[0-9]+$",
			Next: func(flagstream *[]string) error {
				if res, err := strconv.Atoi((*flagstream)[0][2:]); err == nil {
					startPageN = &res
					*flagstream = (*flagstream)[1:]
					return nil
				} else {
					return err
				}
			},
			Usage: "-s<number>\tSpecify the number of start page",
		},
		&flagMatcher{
			Pattern: "^-e[0-9]+$",
			Next: func(flagstream *[]string) error {
				if res, err := strconv.Atoi((*flagstream)[0][2:]); err == nil {
					endPageN = &res
					*flagstream = (*flagstream)[1:]
					return nil
				} else {
					return err
				}
			},
			Usage: "-e<number>\tSpecify the number of end page",
		},
		&flagMatcher{
			Pattern: "^-l[0-9]+$",
			Next: func(flagstream *[]string) error {
				if res, err := strconv.Atoi((*flagstream)[0][2:]); err == nil {
					linesPerPage = &res
					*flagstream = (*flagstream)[1:]
					return nil
				} else {
					return err
				}
			},
			Usage: "-l<number>\tSpecify the number of lines in each page(default is 72)",
		},
		&flagMatcher{
			Pattern: "^-f$",
			Next: func(flagstream *[]string) error {
				res := true
				formFeed = &res
				*flagstream = (*flagstream)[1:]
				return nil
			},
			Usage: "-f\t\tUse ascci \\f as page breaking flag",
		},
		&flagMatcher{
			Pattern: "^-d.*$",
			Next: func(flagstream *[]string) error {
				res := (*flagstream)[0][2:]
				toDestination = &res
				*flagstream = (*flagstream)[1:]
				return nil
			},
			Usage: "-d\t\tSpecify the destination printer to send result",
		},
		&flagMatcher{
			Pattern: "^-h$",
			Next: func(flagstream *[]string) error {
				res := true
				showHelpMsg = &res
				*flagstream = (*flagstream)[1:]
				return nil
			},
			Usage: "-h\t\tShow help messages",
		},
		&flagMatcher{
			Pattern: ".*",
			Next: func(flagstream *[]string) error {
				filename = &(*flagstream)[0]
				*flagstream = (*flagstream)[1:]
				return nil
			},
			Usage: "[filename]\tSpecify the input file",
		},
	}
}

func parseFlags(flagstream *[]string) error {
	initMatchers()
	for len(*flagstream) > 0 {
		for _, matcher := range matchers {
			if matched, err := regexp.MatchString(matcher.Pattern, (*flagstream)[0]); err == nil && matched {
				if err = matcher.Next(flagstream); err != nil {
					return err
				} else {
					break
				}
			}
		}
	}
	return checkFlags()
}

func printUsage() {
	fmt.Println("Usage:")
	for _, matcher := range matchers {
		fmt.Println("\t", matcher.Usage)
	}
	fmt.Println()
}

func initValidatiors() {
	validators = []*flagValidator{
		&flagValidator{Condition: startPageN == nil, ErrorMsg: "flag --start-page is mandatory"},
		&flagValidator{Condition: endPageN == nil, ErrorMsg: "flag --end-page is mandatory"},
		&flagValidator{Condition: formFeed != nil && linesPerPage != nil, ErrorMsg: "flag --form-feed and flag --lines-per-page are inconsistent!"},
	}
}

func checkFlags() error {
	if showHelpMsg != nil && *showHelpMsg {
		return nil
	}
	initValidatiors()
	for _, validator := range validators {
		if validator.Condition {
			return errors.New(validator.ErrorMsg)
		}
	}
	if formFeed == nil && linesPerPage == nil {
		temp1 := 72
		linesPerPage = &temp1
		temp2 := false
		formFeed = &temp2
	}
	if showHelpMsg == nil {
		temp := false
		showHelpMsg = &temp
	}
	return nil
}
