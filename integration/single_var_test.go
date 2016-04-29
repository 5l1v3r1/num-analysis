package integration

import (
	"math"
	"testing"
)

func TestIntegrateDegree0(t *testing.T) {
	actual := IntegrateDegree(quadraticTestFunc, Interval{-2, 10}, 1, 0)
	if math.Abs(actual-203) > 1e-10 {
		t.Error("expected integral 203 but got", actual)
	}

	actual = IntegrateDegree(quadraticTestFunc, Interval{10, -2}, 0.5, 0)
	if math.Abs(actual+203.75) > 1e-10 {
		t.Error("expected integral -203.75 but got", actual)
	}
}

func TestIntegrateDegree1(t *testing.T) {
	actual := IntegrateDegree(quadraticTestFunc, Interval{5, 10 + 1e-9}, 0.25, 1)
	expected := 201.71875
	if math.Abs(actual-expected) > 1e-8 {
		t.Error("expected integral", expected, "but got", actual)
	}
}

func TestIntegralDegree2(t *testing.T) {
	actual := IntegrateDegree(quadraticTestFunc, Interval{5, 10 + 1e-9}, 0.25, 2)
	expected := 605.0 / 3.0
	if math.Abs(actual-expected) > 1e-8 {
		t.Error("expected integral", expected, "but got", actual)
	}

	actual = IntegrateDegree(cubicTestFunc, Interval{5, 10 + 1e-9}, 0.25, 2)
	expected = -81955.0 / 12.0
	if math.Abs(actual-expected) > 1e-8 {
		t.Error("expected integral", expected, "but got", actual)
	}

	actual = IntegrateDegree(quarticTestFunc, Interval{5, 10 + 1e-9}, 0.25, 2)
	notExpected := 383045.0 / 12.0
	if math.Abs(actual-notExpected) < 1e-8 {
		t.Error("second degree integral should have failed")
	}
}

func quadraticTestFunc(f float64) float64 {
	return math.Pow(f, 2) - 2*f - 3
}

func cubicTestFunc(f float64) float64 {
	return math.Pow(f, 2) - 2*f - 3 - 3*math.Pow(f, 3)
}

func quarticTestFunc(f float64) float64 {
	return math.Pow(f, 2) - 2*f - 3 - 3*math.Pow(f, 3) + 2*math.Pow(f, 4)
}
