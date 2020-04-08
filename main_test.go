package main

import (
	"testing"
)

func TestEqual(t *testing.T) {
	testCases := []struct {
		v1     string
		v2     string
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
		v1, err := NewVersion(testCase.v1)
		if err != nil {
			t.Errorf("parse error")
		}
		v2, err := NewVersion(testCase.v2)
		if err != nil {
			t.Errorf("parse error")
		}
		actual := v1.Equal(*v2)
		if testCase.expect != actual {
			t.Errorf("No: %d\nactual:%t\nexpect:%t\nv1: %s\nv2: %s\n", i+1, actual, testCase.expect, testCase.v1, testCase.v2)
		}
	}
}

func TestGreaterThan(t *testing.T) {
	testCases := []struct {
		v1     string
		v2     string
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
		{
			v1:     "1-foo",
			v2:     "1.foo",
			expect: true,
		},
		{
			v1:     "1-1",
			v2:     "1-foo",
			expect: true,
		},
		{
			v1:     "1.1",
			v2:     "1-1",
			expect: true,
		},
	}
	for i, testCase := range testCases {
		v1, err := NewVersion(testCase.v1)
		if err != nil {
			t.Errorf("parse error")
		}
		v2, err := NewVersion(testCase.v2)
		if err != nil {
			t.Errorf("parse error")
		}
		actual := v1.GreaterThan(*v2)
		if testCase.expect != actual {
			t.Errorf("No: %d\nactual:%t\nexpect:%t\nv1: %s\nv2: %s\n", i+1, actual, testCase.expect, testCase.v1, testCase.v2)
		}
	}
}

func TestLessThan(t *testing.T) {
	testCases := []struct {
		v1     string
		v2     string
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
		v1, err := NewVersion(testCase.v1)
		if err != nil {
			t.Errorf("parse error")
		}
		v2, err := NewVersion(testCase.v2)
		if err != nil {
			t.Errorf("parse error")
		}
		actual := v1.LessThan(*v2)
		if testCase.expect != actual {
			t.Errorf("No: %d\nactual:%t\nexpect:%t\nv1: %s\nv2: %s\n", i+1, actual, testCase.expect, testCase.v1, testCase.v2)
		}
	}

}

func TestVersionQualifier(t *testing.T) {
	versionsQualifier := []string{"1-alpha2snapshot", "1-alpha2", "1-alpha-123", "1-beta-2", "1-beta123", "1-m2", "1-m11", "1-rc", "1-cr2",
		"1-rc123", "1-SNAPSHOT", "1", "1-sp", "1-sp2", "1-sp123", "1-abc", "1-def", "1-pom-1", "1-1-snapshot",
		"1-1", "1-2", "1-123"}
	for i := 1; i < len(versionsQualifier); i++ {
		low, err := NewVersion(versionsQualifier[i-1])
		if err != nil {
			t.Errorf("parse error")
		}
		for j := i; j < len(versionsQualifier); j++ {
			high, err := NewVersion(versionsQualifier[j])
			if err != nil {
				t.Errorf("parse error")
			}
			if !low.LessThan(*high) {
				t.Errorf("expected: %s < %s \n", low, high)
			}
			if !high.GreaterThan(*low) {
				t.Errorf("expected: %s > %s \n", high, low)
			}
		}
	}
}

func TestVersionsNumber(t *testing.T) {
	versionsNumber := []string{"2.0", "2-1", "2.0.a", "2.0.0.a", "2.0.2", "2.0.123", "2.1.0", "2.1-a", "2.1b", "2.1-c", "2.1-1", "2.1.0.1",
		"2.2", "2.123", "11.a2", "11.a11", "11.b2", "11.b11", "11.m2", "11.m11", "11", "11.a", "11b", "11c", "11m"}
	for i := 1; i < len(versionsNumber); i++ {
		low, err := NewVersion(versionsNumber[i-1])
		if err != nil {
			t.Errorf("parse error")
		}
		for j := i; j < len(versionsNumber); j++ {
			high, err := NewVersion(versionsNumber[j])
			if err != nil {
				t.Errorf("parse error")
			}
			if !low.LessThan(*high) {
				t.Errorf("expected: %s < %s \n", low.value, high.value)
			}
			if !high.GreaterThan(*low) {
				t.Errorf("expected: %s > %s \n", high, low)
			}
		}
	}
}
