// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package controllers

import (
	"math"

	"gonum.org/v1/gonum/stat"

	"encoding/csv"
	"os"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

// CovarianceMatrix calculates the covariance matrix (also known as the
// variance-covariance matrix) calculated from a matrix of data, x, using
// a two-pass algorithm. The result is stored in dst.
//
// If weights is not nil the weighted covariance of x is calculated. weights
// must have length equal to the number of rows in input data matrix and
// must not contain negative elements.
// The dst matrix must either be zero-sized or have the same number of
// columns as the input data matrix.
func CovarianceMatrix2(dst *mat.SymDense, x mat.Matrix, weights []float64) (Means, Stddevs []float64) {
	// This is the matrix version of the two-pass algorithm. It doesn't use the
	// additional floating point error correction that the Covariance function uses
	// to reduce the impact of rounding during centering.
	var Meanslice []float64
	var Stddevslice []float64
	r, c := x.Dims()

	if dst.IsZero() {
		dst = (dst.GrowSym(c).(*mat.SymDense))
	} else if n := dst.Symmetric(); n != c {
		panic(mat.ErrShape)
	}

	var xt mat.Dense
	xt.Clone(x.T())
	// Subtract the mean of each of the columns.
	for i := 0; i < c; i++ {
		v := xt.RawRowView(i)
		// This will panic with ErrShape if len(weights) != len(v), so
		// we don't have to check the size later.
		mean := stat.Mean(v, weights)
		stdev := stat.StdDev(v, nil)
		floats.AddConst(-mean, v)
		Meanslice = append(Meanslice, mean)
		Stddevslice = append(Stddevslice, stdev)
	}

	if weights == nil {
		// Calculate the normalization factor
		// scaled by the sample size.
		dst.SymOuterK(1/(float64(r)-1), &xt)
		return Meanslice, Stddevslice
	}

	// Multiply by the sqrt of the weights, so that multiplication is symmetric.
	sqrtwts := make([]float64, r)
	for i, w := range weights {
		if w < 0 {
			panic("stat: negative covariance matrix weights")
		}
		sqrtwts[i] = math.Sqrt(w)
	}
	// Weight the rows.
	for i := 0; i < c; i++ {
		v := xt.RawRowView(i)
		floats.Mul(v, sqrtwts)
	}

	// Calculate the normalization factor
	// scaled by the weighted sample size.
	dst.SymOuterK(1/(floats.Sum(weights)-1), &xt)
	return Meanslice, Stddevslice
}

// CorrelationMatrix returns the correlation matrix calculated from a matrix
// of data, x, using a two-pass algorithm. The result is stored in dst.
//
// If weights is not nil the weighted correlation of x is calculated. weights
// must have length equal to the number of rows in input data matrix and
// must not contain negative elements.
// The dst matrix must either be zero-sized or have the same number of
// columns as the input data matrix.

// Open comment
func Open(column int, filename string) []string {
	//filename := "REMB22 Student Data UK.csv"
	f, _ := os.Open(filename)
	thuy := make([]string, 0)
	defer f.Close()
	Lines, _ := csv.NewReader(f).ReadAll()
	for _, v := range Lines {
		thuy = append(thuy, v[column])
	}
	return thuy
}
