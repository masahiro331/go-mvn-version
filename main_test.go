package main

import (
	"testing"
)

func TestEqual(t *testing.T) {
	testCases := []struct {
		v1     Version
		v2     Version
		expect bool
	}{
		{
			v1:     "1--1",
			v2:     "1-0-1",
			expect: true,
		},
		{
			v1:     "1--1",
			v2:     "1-1",
			expect: false,
		},
		{
			v1:     "1..1",
			v2:     "1.0.1",
			expect: true,
		},
		{
			v1:     "1.1",
			v2:     "1.0.1",
			expect: false,
		},
		{
			v1:     "1-ga",
			v2:     "1",
			expect: true,
		},
		{
			v1:     "1-sp",
			v2:     "1",
			expect: false,
		},
		{
			v1:     "1.0.0.RELEASE",
			v2:     "1",
			expect: true,
		},
		{
			v1:     "1.0.0.FINAL",
			v2:     "1",
			expect: true,
		},
		{
			v1:     "1.0.0.FINAL",
			v2:     "1.RELEASE",
			expect: true,
		},
		{
			v1:     "1.2.3-a1b1-m1",
			v2:     "1.2.3-alpha-1-beta-1-milestone-1",
			expect: true,
		},
		{
			v1:     "1.2.3-rc",
			v2:     "1.2.3-cr",
			expect: true,
		},
		{
			v1:     "5.0.0.RELEASE",
			v2:     "4.9.9.RELEASE",
			expect: false,
		},
		{
			v1:     "-1",
			v2:     "1",
			expect: false,
		},
		{
			v1:     "-1",
			v2:     "0-1",
			expect: true,
		},
		{
			v1:     "1-0.3",
			v2:     "1",
			expect: true,
		},
	}
	for i, testCase := range testCases {
		actual := testCase.v1.Equal(testCase.v2)
		if testCase.expect != actual {
			t.Errorf("No: %d\nactual:%t\nexpect:%t\nv1: %s\nv2: %s\n", i+1, actual, testCase.expect, testCase.v1, testCase.v2)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	testCases := []struct {
		v1     Version
		v2     Version
		expect bool
	}{
		{
			v1:     "1--1",
			v2:     "1-0-1",
			expect: false,
		},
		{
			v1:     "1--1",
			v2:     "1-1",
			expect: false,
		},
		{
			v1:     "1..1",
			v2:     "1.0.1",
			expect: false,
		},
		{
			v1:     "1.1",
			v2:     "1.0.1",
			expect: true,
		},
		{
			v1:     "1.0.1",
			v2:     "1.1",
			expect: false,
		},
		{
			v1:     "1-ga",
			v2:     "1",
			expect: false,
		},
		{
			v1:     "1-sp",
			v2:     "1",
			expect: true,
		},
		{
			v1:     "1.0.0.RELEASE",
			v2:     "1",
			expect: false,
		},
		{
			v1:     "1.0.0.FINAL",
			v2:     "1",
			expect: false,
		},
		{
			v1:     "1.0.0.FINAL",
			v2:     "1.RELEASE",
			expect: false,
		},
		{
			v1:     "1.2.3",
			v2:     "1.2.3-a1",
			expect: true,
		},
		{
			v1:     "1.2.3-b1",
			v2:     "1.2.3-a1",
			expect: true,
		},
		{
			v1:     "1.2.3-m1",
			v2:     "1.2.3-b1",
			expect: true,
		},
		{
			v1:     "1.2.3-rc",
			v2:     "1.2.3-m1",
			expect: true,
		},
		{
			v1:     "1.2.3-a2",
			v2:     "1.2.3-a1",
			expect: true,
		},
		{
			v1:     "1.2.3-b1",
			v2:     "1.2.3-a2",
			expect: true,
		},
		{
			v1:     "1.2.3",
			v2:     "1.2.3-cr",
			expect: true,
		},
		{
			v1:     "5.0.0.RELEASE",
			v2:     "4.9.9.RELEASE",
			expect: true,
		},
		{
			v1:     "1",
			v2:     "-1",
			expect: true,
		},
		{
			v1:     "1-0.3",
			v2:     "1",
			expect: true,
		},
	}
	for i, testCase := range testCases {
		actual := testCase.v1.GreaterThan(testCase.v2)
		if testCase.expect != actual {
			t.Errorf("No: %d\nactual:%t\nexpect:%t\nv1: %s\nv2: %s\n", i+1, actual, testCase.expect, testCase.v1, testCase.v2)
		}
	}
}

func TestLessThan(t *testing.T) {
	testCases := []struct {
		v1     Version
		v2     Version
		expect bool
	}{
		{
			v1:     "1--1",
			v2:     "1-0-1",
			expect: false,
		},
		{
			v1:     "1--1",
			v2:     "1-1",
			expect: true,
		},
		{
			v1:     "1..1",
			v2:     "1.0.1",
			expect: false,
		},
		{
			v1:     "1.0.1",
			v2:     "1.0.11",
			expect: true,
		},
		{
			v1:     "1.0.1",
			v2:     "1.1",
			expect: true,
		},
		{
			v1:     "1-0",
			v2:     "1-sp",
			expect: true,
		},
		{
			v1:     "1.2.3-a1",
			v2:     "1.2.3",
			expect: true,
		},
		{
			v1:     "1.2.3-a1",
			v2:     "1.2.3-b1",
			expect: true,
		},
		{
			v1:     "1.2.3-b1",
			v2:     "1.2.3-m1",
			expect: true,
		},
		{
			v1:     "1.2.3-m1",
			v2:     "1.2.3-rc",
			expect: true,
		},
		{
			v1:     "1.2.3-a1",
			v2:     "1.2.3-a2",
			expect: true,
		},
		{
			v1:     "1.2.3-a2",
			v2:     "1.2.3-b1",
			expect: true,
		},
		{
			v1:     "1.2.3-cr",
			v2:     "1.2.3",
			expect: true,
		},
		{
			v1:     "4.9.9.RELEASE",
			v2:     "5.0.0.RELEASE",
			expect: true,
		},
		{
			v1:     "0-1",
			v2:     "1",
			expect: true,
		},
	}
	for i, testCase := range testCases {
		actual := testCase.v1.LessThan(testCase.v2)
		if testCase.expect != actual {
			t.Errorf("No: %d\nactual:%t\nexpect:%t\nv1: %s\nv2: %s\n", i+1, actual, testCase.expect, testCase.v1, testCase.v2)
		}
	}

}
