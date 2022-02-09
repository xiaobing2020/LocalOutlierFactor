package main

import (
    "fmt"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/plotutil"
    "gonum.org/v1/plot/vg"
    "math"
    "sort"
    "strconv"
)

// Keep2Decimals 浮点数保留2位小数
func Keep2Decimals(num float64) float64 {
    str := fmt.Sprintf("%.2f", num)
    ret, _ := strconv.ParseFloat(str,64)
    return ret
}

// 计算两个n维向量之间的距离
func distance(p1, p2 []float64, order int) float64 {
    // 传入两个向量维数不一样或维数不为n返回0
    if len(p1) != len(p2) || len(p1) != order {
        return -1
    }
    var squareSum float64
    for i := 0; i < order; i++ {
        squareSum += (p2[i] - p1[i]) * (p2[i] - p1[i])
    }
    return math.Sqrt(squareSum)
}

// 计算点p的第k距离 dataSet内无重复点
func d_k_p(p []float64, dataSet [][]float64, k int, order int) float64 {
    size := len(dataSet)
    if size < k + 1 {
        return -1
    }
    var d float64
    var dSet []float64
    for i := 0; i < size; i++ {
        d = distance(p ,dataSet[i], order)
        if d == -1 { // 传入数据维数不一致
            return -1
        } else if d == 0 { // 遍历到p点 p属于dataSet
            continue
        } else {
            dSet = append(dSet, Keep2Decimals(d))
        }
    }
    sort.Float64s(dSet)
    return dSet[k - 1]
}

// 返回点p的第k距离邻域集合
func n_k_p(p []float64, dataSet [][]float64, k int, order int) [][]float64 {
    var nkpSet [][]float64
    var d float64
    size := len(dataSet)
    s := d_k_p(p, dataSet, k, order) //计算第k距离
    for i := 0; i < size; i++ {
        d = distance(p ,dataSet[i], order)
        if d == 0 { // 遍历到p点 p属于dataSet
            continue
        } else if d <= s {
            nkpSet = append(nkpSet, dataSet[i])
        }
    }
    return nkpSet
}
// 返回点A相对点P的可达距离
// dkp 点p在其数据集内的第k距离
func reachDistance(P []float64, A []float64, dkp float64, order int) float64 {
    d := distance(P, A, order)
    if d > dkp {
        return d
    } else {
        return dkp
    }
}

// 计算点p局部可达密度
func lrd_k_P(p []float64, Nkp [][]float64, dkp float64, order int) float64 {
    size := len(Nkp)
    var sumReachDistance float64
    for i := 0; i < size; i++ {
        sumReachDistance += reachDistance(p, Nkp[i], dkp, order)
    }
    return 1.0 / (sumReachDistance / float64(size))
}

// LocalOutlierFactor 局部离群因子计算
func LocalOutlierFactor(p []float64, dataSet [][]float64, k, order int) float64 {
    dkp := d_k_p(p, dataSet, k, order)
    nkp := n_k_p(p, dataSet, k, order)
    lrdkp := lrd_k_P(p, nkp, dkp, order)
    size := len(nkp)
    var sumLrdip float64
    for i := 0; i < size; i++ {
        dki := d_k_p(nkp[i], dataSet, k, order)
        nki := n_k_p(nkp[i], dataSet, k, order)
        lrdki := lrd_k_P(nkp[i], nki, dki, order)
        sumLrdip += lrdki
    }
    return Keep2Decimals(sumLrdip / (float64(size) * lrdkp))
}

// 画图
func drawPoints(data [][]float64) {
    points := make(plotter.XYs, 100)
    for i := 0; i < len(data); i++ {
        points[i].X = data[i][0]
        points[i].Y = data[i][1]
    }
    plt := plot.New()
    if err := plotutil.AddScatters(plt,
        "x3", points,
    ); err != nil {
        panic(err)
    }

    if err := plt.Save(5*vg.Inch, 5*vg.Inch, "source.png"); err != nil {
        panic(err)
    }
}

func main() {
    // 测试用例
    source := [][]float64{{1.1, 1.1}, {1.2, 1.15}, {2.1, 2.1}, {2.1, 5.1}, {1.3,1.2},
        {3.1,4.1},{4.2,3.9}, {4.1, 4.2}, {4.2,4.1},{6.1,3.9}}
    var r float64
    ar := make([]float64, 10)
    for i:=0;i<10;i++ {
        r = LocalOutlierFactor(source[i],source,4,2)
        ar[i] = r
    }
    fmt.Println(ar)
    drawPoints(source)
}
