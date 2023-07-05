package validator_test

import (
	"testing"

	"github.com/bwc00/strv-go-newsletter-shakleya-mohammed/util/validator"
)

type testCase struct {
	name     string
	input    interface{}
	expected string
}

var tests = []*testCase{
	{
		name: `required`,
		input: struct {
			Title string `json:"title" form:"required"`
		}{},
		expected: "title is a required field",
	},
	{
		name: `max`,
		input: struct {
			Course string `json:"course" form:"max=7"`
		}{Course: "CS-0001."},
		expected: "course must be a maximum of 7 in length",
	},
	{
		name: `email`,
		input: struct {
			Email string `json:"email" form:"email"`
		}{Email: "shaki2632.com"},
		expected: "email must be a valid Email",
	},
	{
		name: `alpha_zero`,
		input: struct {
			Name string `json:"name" form:"alpha_zero"`
		}{Name: "Some Name 2"},
		expected: "name can only contain alphabetic and space characters",
	},
}

func TestToErrResponse(t *testing.T) {
	vr := validator.New()

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := vr.Struct(tc.input)
			if errResp := validator.ToErrResponse(err); errResp == nil || len(errResp.Errors) != 1 {
				t.Fatalf(`Expected:"{[%v]}", Got:"%v"`, tc.expected, errResp)
			} else if errResp.Errors[0] != tc.expected {
				t.Fatalf(`Expected:"%v", Got:"%v"`, tc.expected, errResp.Errors[0])
			}
		})
	}
}
