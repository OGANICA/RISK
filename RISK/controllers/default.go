package controllers

import (
	"fmt"
	"math"
	"strconv"

	"github.com/astaxie/beego"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat"
)

type MainController struct {
	beego.Controller
}

type MetricsController struct {
	beego.Controller
}

var chdata string = "REMB22 Student Data AU.csv"

func makematrix(s string) (srcmatrix *mat.Dense, header []string) {
	leng := len(Open(0, chdata)) - 1
	switch s {
	case "v":
	case "r":
		leng = leng - 1
	}
	srcmatrix = mat.NewDense(leng, 8, nil)
	a := make([]string, 0)
	finslice := make([]float64, 0)
	header = make([]string, 0)
	bv := make([]float64, 0)
	ev := make([]float64, 0)
	hpr := 0.0
	for i := 1; i <= 8; i++ {
		a = Open(i, chdata)
		header = append(header, a[0])
		a = a[1:]
		for _, v := range a {
			x, _ := strconv.ParseFloat(v, 64)
			finslice = append(finslice, x)
		}
		switch s {
		case "v":
		case "r":
			bv = finslice[0 : len(finslice)-1]
			ev = finslice[1:]
			finslice = nil
			for h := 0; h <= leng-1; h++ {
				hpr = (ev[h] - bv[h]) / bv[h]
				finslice = append(finslice, hpr)
			}
			bv = nil
			ev = nil
		}
		srcmatrix.SetCol(i-1, finslice)
		finslice = nil
	}
	return srcmatrix, header
}

func covmatrix(s string) (finalmatrix *mat.SymDense, header []string, meanslice []float64, stdevs []float64) {
	dstmatrix := mat.NewSymDense(8, nil)
	a, header := makematrix(s)
	meanslice, stdevs = CovarianceMatrix2(dstmatrix, a, nil)
	return dstmatrix, header, meanslice, stdevs
}

func correlmatrix(s string) (cmatrix *mat.SymDense) {
	dstmatrix := mat.NewSymDense(8, nil)
	a, _ := makematrix(s)
	stat.CorrelationMatrix(dstmatrix, a, nil)
	return dstmatrix
}

func valueatrisk(signif int, mean []float64, stdev []float64) (valueatrisk []float64) {
	vatr := 0.0
	c := 0.0
	switch signif {
	case 5:
		c = 1.65
	case 1:
		c = 2.33
	}
	valueatrisk = make([]float64, 0)
	for h := 0; h <= 7; h++ {
		vatr = (mean[h] - stdev[h]*c) * 100
		valueatrisk = append(valueatrisk, vatr)
	}
	return valueatrisk
}

func roundslice(precision float64, scale float64, source []float64) []float64 {
	precision = math.Pow(10, precision)
	scale = math.Pow(10, scale)
	target := make([]float64, 0)
	for _, v := range source {
		v = v * precision * scale
		v = math.Round(v)
		v = v / precision
		target = append(target, v)
	}
	return target
}

func (c *MetricsController) Get() {
	temp := make(map[interface{}]interface{})
	dstmatrix, names, meanslice, stdevs := covmatrix("r")
	vatr := valueatrisk(5, meanslice, stdevs)
	correlmatrix("r")
	matrix := mat.Formatted(dstmatrix, mat.Prefix(""), mat.Squeeze())
	temp["header"] = names
	temp["matrix"] = matrix
	temp["VaR"] = roundslice(2, 0, vatr)
	temp["meanslice"] = roundslice(2, 2, meanslice)
	temp["stdevs"] = roundslice(2, 2, stdevs)
	temp["varoption"] = 5
	c.Data = temp
	c.TplName = "template2.tpl"
}

func (c *MetricsController) Post() {
	c.TplName = "template2.tpl"
	c.Data["vatr"] = c.GetString("vatr")
	fmt.Println(c.Data)
	x, _ := c.Data["vatr"].(string)
	vatrint64, _ := strconv.ParseInt(x, 0, 64)
	vatrint := int(vatrint64)
	fmt.Println(vatrint)
	temp := make(map[interface{}]interface{})
	dstmatrix, header, meanslice, stdevs := covmatrix("r")
	vatr := valueatrisk(vatrint, meanslice, stdevs)
	correlmatrix("r")
	matrix := mat.Formatted(dstmatrix, mat.Prefix(""), mat.Squeeze())
	temp["header"] = header
	temp["matrix"] = matrix
	temp["VaR"] = roundslice(2, 0, vatr)
	temp["meanslice"] = roundslice(2, 2, meanslice)
	temp["stdevs"] = roundslice(2, 2, stdevs)
	temp["varoption"] = vatrint
	temp["fred"] = fredapisorted("A191RL1Q225SBEA")
	c.Data = temp
}
