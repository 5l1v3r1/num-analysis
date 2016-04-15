package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/unixpickle/num-analysis/mvroots"
)

var polyTerm = regexp.MustCompile("^(-?[0-9\\.]*)(x(\\^([0-9*]))?)?$")

func ParsePolynomial(s string) (mvroots.Polynomial, error) {
	coeffs := map[int]float64{}

	isOperator := false
	negateTerm := false
	terms := strings.Fields(s)

	for _, term := range terms {
		if isOperator {
			if term == "+" {
				negateTerm = false
			} else if term == "-" {
				negateTerm = true
			} else {
				return nil, fmt.Errorf("invalid term: %s", term)
			}
			isOperator = false
			continue
		}

		match := polyTerm.FindStringSubmatch(term)
		if match == nil {
			return nil, fmt.Errorf("invalid term format: %s", term)
		}

		coefficient := 1.0
		if len(match[1]) > 0 {
			coeff, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				return nil, fmt.Errorf("invalid coefficient in term: %s", term)
			}
			coefficient = coeff
		}

		degree := 0
		switch match[2] {
		case "":
		case "x":
			degree = 1
		default:
			deg, err := strconv.Atoi(match[4])
			if err != nil {
				return nil, fmt.Errorf("invalid exponent in term: %s", term)
			}
			degree = deg
		}

		if negateTerm {
			coefficient *= -1
		}
		coeffs[degree] += coefficient

		isOperator = true
	}

	if !isOperator {
		return nil, errors.New("polynomial ended with an operator")
	}
	return polynomialFromMap(coeffs), nil
}

func polynomialFromMap(m map[int]float64) mvroots.Polynomial {
	degree := 0
	for d := range m {
		if d > degree {
			degree = d
		}
	}
	res := make(mvroots.Polynomial, degree+1)
	for d, c := range m {
		res[d] = complex(c, 0)
	}
	return res
}
