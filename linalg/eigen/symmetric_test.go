package eigen

import (
	"math/rand"
	"testing"
	"time"

	"github.com/unixpickle/num-analysis/linalg"
)

var symMat10x10 = &linalg.Matrix{
	Rows: 10,
	Cols: 10,
	Data: []float64{
		-2.17007153473045e+00, -1.78134953679285e+00, -1.84668691578582e+00, -1.92986014687707e+00,
		-1.38510914274212e+00, -2.28127738041142e+00, -2.43236747480013e+00, -1.77802358632556e+00,
		-2.00245801194590e+00, -1.37334257246506e+00, -1.78134953679285e+00, 1.98698378200736e+00,
		1.83604322922824e+00, 1.68049256592970e+00, 1.14772124428194e+00, 1.71396348461814e+00,
		2.47216534076068e+00, 1.50310142655288e+00, 1.67691295950850e+00, 1.84755046819251e+00,
		-1.84668691578582e+00, 1.83604322922824e+00, 2.26381675373077e+00, 2.09474110019186e+00,
		1.24300559159416e+00, 2.36385026924340e+00, 2.72504361975536e+00, 1.72583018622844e+00,
		2.00317297862235e+00, 1.94368488615468e+00, -1.92986014687707e+00, 1.68049256592970e+00,
		2.09474110019186e+00, 5.53062877402673e+00, 2.38623023212995e+00, 3.39568245104898e+00,
		4.00162763626571e+00, 2.62566502250379e+00, 4.16072829966437e+00, 3.02352912696568e+00,
		-1.38510914274212e+00, 1.14772124428194e+00, 1.24300559159416e+00, 2.38623023212995e+00,
		1.51969187836467e+00, 1.85960003638432e+00, 2.13663155324393e+00, 1.84534600968712e+00,
		1.87694630530958e+00, 1.23759212168353e+00, -2.28127738041142e+00, 1.71396348461814e+00,
		2.36385026924340e+00, 3.39568245104898e+00, 1.85960003638432e+00, 3.60001021376507e+00,
		3.03627297433905e+00, 2.55894384615425e+00, 3.34740053477478e+00, 2.15130484636945e+00,
		-2.43236747480013e+00, 2.47216534076068e+00, 2.72504361975536e+00, 4.00162763626571e+00,
		2.13663155324393e+00, 3.03627297433905e+00, 4.61435633264008e+00, 2.74968431909778e+00,
		3.32032430664810e+00, 2.69304419822027e+00, -1.77802358632556e+00, 1.50310142655288e+00,
		1.72583018622844e+00, 2.62566502250379e+00, 1.84534600968712e+00, 2.55894384615425e+00,
		2.74968431909778e+00, 2.99552196313698e+00, 2.34571412022552e+00, 1.51627309920713e+00,
		-2.00245801194590e+00, 1.67691295950850e+00, 2.00317297862235e+00, 4.16072829966437e+00,
		1.87694630530958e+00, 3.34740053477478e+00, 3.32032430664810e+00, 2.34571412022552e+00,
		3.97491391327801e+00, 2.54986679845091e+00, -1.37334257246506e+00, 1.84755046819251e+00,
		1.94368488615468e+00, 3.02352912696568e+00, 1.23759212168353e+00, 2.15130484636945e+00,
		2.69304419822027e+00, 1.51627309920713e+00, 2.54986679845091e+00, 2.97978592140745e+00,
	},
}

func symmetricEigenSolver(m *linalg.Matrix) ([]float64, []linalg.Vector) {
	rand.Seed(time.Now().UnixNano())
	return Symmetric(m)
}

func symmetricPrecEigenSolver(m *linalg.Matrix) ([]float64, []linalg.Vector) {
	rand.Seed(time.Now().UnixNano())
	a, b, _ := SymmetricPrec(m, time.Second, 1e-6)
	return a, b
}

func symmetricTimeEigenSolver(m *linalg.Matrix) ([]float64, []linalg.Vector) {
	rand.Seed(time.Now().UnixNano())
	return SymmetricFixedTime(m, time.Millisecond*30)
}

func TestSymmetricBasic(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 3,
		Cols: 3,
		Data: []float64{
			66, 78, 76,
			78, 93, 92,
			76, 92, 94,
		},
	}
	eigs := []float64{4.81397359013199e-02, 2.99176945337813e+00, 2.49960090810721e+02}
	testEigenSolver(t, symmetricEigenSolver, mat, eigs)
	testEigenSolver(t, symmetricPrecEigenSolver, mat, eigs)
	testEigenSolver(t, symmetricTimeEigenSolver, mat, eigs)
}

func TestSymmetricNullspace(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 3,
		Cols: 3,
		Data: []float64{
			66, 78, 90,
			78, 93, 108,
			90, 108, 126,
		},
	}
	eigs := []float64{0, 1.14141341962985e+00, 2.83858586580370e+02}
	testEigenSolver(t, symmetricEigenSolver, mat, eigs)
	testEigenSolver(t, symmetricPrecEigenSolver, mat, eigs)
	testEigenSolver(t, symmetricTimeEigenSolver, mat, eigs)
}

func TestSymmetric10x10(t *testing.T) {
	eigs := []float64{-3.53320764624989e+00,
		1.94571466978943e-02, 3.94135968791024e-02, 2.79652524908013e-01,
		3.83072877722642e-01, 6.66544542615382e-01, 1.16866047971769e+00,
		1.83799425365499e+00, 2.46391983763316e+00, 2.39701303840477e+01}
	testEigenSolver(t, symmetricEigenSolver, symMat10x10, eigs)
	testEigenSolver(t, symmetricPrecEigenSolver, symMat10x10, eigs)
	testEigenSolver(t, symmetricTimeEigenSolver, symMat10x10, eigs)
}

func TestSymmetricRepeatedEig(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 4,
		Cols: 4,
		Data: []float64{
			1.83316880813395e+00, 7.89456964650591e-01, -8.75952978517168e-01,
			-2.93415346421698e-01,
			7.89456964650591e-01, 1.35250218801815e+00, 1.84867985312246e-01, -7.67609461812816e-01,
			-8.75952978517168e-01, 1.84867985312246e-01, -6.82998775769859e-01,
			1.56466447173890e+00,
			-2.93415346421698e-01, -7.67609461812816e-01, 1.56466447173890e+00,
			4.97327779617766e-01,
		},
	}
	eigs := []float64{1, 1, 3, -2}
	testEigenSolver(t, symmetricEigenSolver, mat, eigs)
	testEigenSolver(t, symmetricPrecEigenSolver, mat, eigs)
	testEigenSolver(t, symmetricTimeEigenSolver, mat, eigs)
}

func TestSymmetricNearEigenvalues(t *testing.T) {
	mat := &linalg.Matrix{
		Rows: 4,
		Cols: 4,
		Data: []float64{
			1.82867576281371e+00, 7.93691175447142e-01, -8.75546706294582e-01,
			-2.95993989236086e-01,
			7.93691175447142e-01, 1.34851190085392e+00, 1.84485117544706e-01, -7.65179368976683e-01,
			-8.75546706294582e-01, 1.84485117544706e-01, -6.83035511904726e-01,
			1.56489763897243e+00,
			-2.95993989236086e-01, -7.65179368976683e-01, 1.56489763897243e+00,
			4.95847848237095e-01,
		},
	}
	eigs := []float64{1, 0.99, 3, -2}
	testEigenSolver(t, symmetricEigenSolver, mat, eigs)
	testEigenSolver(t, symmetricPrecEigenSolver, mat, eigs)
	testEigenSolver(t, symmetricTimeEigenSolver, mat, eigs)
}

func TestSymmetricAsyncCancel(t *testing.T) {
	mat := randomSymMatrix(100)
	res := SymmetricAsync(mat)
	go func() {
		time.Sleep(time.Millisecond * 30)
		close(res.Cancel)
	}()
	timeout := time.After(time.Second)
	for {
		select {
		case _, ok := <-res.Values:
			if !ok {
				return
			}
		case <-timeout:
			t.Error("solver was not cancelled")
			return
		}
	}
}

func BenchmarkSymmetric10x10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		symmetricEigenSolver(symMat10x10)
	}
}

func BenchmarkSymmetric50x50(b *testing.B) {
	mat := randomSymMatrix(50)
	for i := 0; i < b.N; i++ {
		symmetricEigenSolver(mat)
	}
}

func randomSymMatrix(size int) *linalg.Matrix {
	res := linalg.NewMatrix(size, size)
	for i := 0; i < size; i++ {
		for j := 0; j <= i; j++ {
			val := rand.Float64()
			res.Set(i, j, val)
			res.Set(j, i, val)
		}
	}
	return res
}
