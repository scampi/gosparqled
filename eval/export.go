package eval

import (
    "os"
    "bufio"
    "github.com/golang/glog"
    "strings"
    "strconv"
    "time"
    "path/filepath"
    "regexp"
    "sort"
)

// AggregatedMeasurement represents an aggregation of several Measure
type AggregatedMeasurement struct {
    // Name is the identifier of the set of aggregated Measures
    Name string
    // Jaccard is the average Jaccard similarity of the set of Measures
    // with regards to the set of gold standard recommendations
    Jaccard float64
    // ElapsedTime is the average time taken for retrieving recommendations
    ElapsedTime time.Duration
}

// loadGold reads a list of Recommendations from the given file.
// Each line represents the recommendations for one query.
func loadGold(gold string) [][]Recommendation {
    glog.Infof("Loading Gold [%v]\n", gold)
    fi, err := os.Open(gold)
    if err != nil { glog.Fatal(err) }
    defer fi.Close()
    s := bufio.NewScanner(fi)
    var allrecs [][]Recommendation
    for s.Scan() {
        var recs []Recommendation
        for _,v := range strings.Split(s.Text(), "}") {
            if v == "]" { // end
                break
            }
            lbrace := strings.Index(v, "{")
            space := strings.LastIndex(v, " ")
            if lbrace == -1 || space == -1 {
                glog.Fatalf("Invalid gold: %v\n", v)
            }
            item := v[lbrace+1:space]
            count, _ := strconv.Atoi(v[space+1:])
            recs = append(recs, Recommendation{ Item: item, Count: count })
        }
        allrecs = append(allrecs, recs)
    }
    return allrecs
}

// loadMesure reads a list of Measurements from the given file.
// Each line represents the Measurement for one query.
func loadMeasure(result string) []Measurement {
    fi, err := os.Open(result)
    if err != nil { glog.Fatal(err) }
    defer fi.Close()
    s := bufio.NewScanner(fi)
    var measurements []Measurement
    for s.Scan() {
        m := strings.SplitN(s.Text(), " ", 6)
        min, _ := strconv.Atoi(m[0])
        max, _ := strconv.Atoi(m[1])
        avg, _ := strconv.ParseFloat(m[2], 32)
        length, _ := strconv.Atoi(m[3])
        elapsedTime, _ := time.ParseDuration(m[4])
        var recs []Recommendation
        for _,v := range strings.Split(m[5], "}") {
            if v == "]" { // end
                break
            }
            lbrace := strings.Index(v, "{")
            space := strings.LastIndex(v, " ")
            if lbrace == -1 || space == -1 {
                glog.Fatalf("Invalid measurement: %v\n", v)
            }
            item := v[lbrace+1:space]
            count, _ := strconv.Atoi(v[space+1:])
            recs = append(recs, Recommendation{ Item: item, Count: count })
        }
        measurements = append(measurements, Measurement{ Min: min, Max: max, Avg: float32(avg), Length: length, ElapsedTime: elapsedTime, Recs: recs })
    }
    return measurements
}

// jaccard computes the Jaccard similarity of the two Recommendations sets.
func jaccard(a []Recommendation, b []Recommendation) float64 {
    occ := make(map[string]int)
    for _,r := range a {
        occ[r.Item]++
    }
    for _,r := range b {
        occ[r.Item]++
    }
    inter := float64(0)
    for _,v := range occ {
        if v == 2 {
            inter++
        }
    }
    return inter / float64(len(occ))
}

// compare evaluates the list of Measurements from results
// and returns its AggregatedMeasurement
func compare(gold [][]Recommendation, results string) AggregatedMeasurement {
    glog.Infof("Processing results [%s]\n", results)
    measures := loadMeasure(results)
    if len(gold) != len(measures) {
        glog.Fatalf("Gold=%v Measures=%v", len(gold), len(measures))
    }
    es, j := int64(0), float64(0)
    for i, im := range measures {
        es += int64(im.ElapsedTime)
        j += jaccard(gold[i], im.Recs)
        glog.Infof("Jaccard: %v\n", jaccard(gold[i], im.Recs))
    }
    glog.Infof("Avg: %v\n", j/float64(len(measures)))
    return AggregatedMeasurement{ Name: results, Jaccard: j/float64(len(measures)), ElapsedTime: time.Duration(es/int64(len(measures))) }
}

// bySign sorts AggregatedMeasurements by the query complexity
type bySign []AggregatedMeasurement
func (am bySign) Len() int {
    return len(am)
}
func (am bySign) Swap(i, j int) {
    am[i], am[j] = am[j], am[i]
}
func (am bySign) Less(i, j int) bool {
    signi := signRe.FindString(am[i].Name)
    signj := signRe.FindString(am[j].Name)
    si := strings.Split(signi[strings.Index(signi, "_")+1:], "-")
    sj := strings.Split(signj[strings.Index(signj, "_")+1:], "-")
    if len(si) < len(sj) { return true }
    if len(si) > len(sj) { return false }
    for ind := 0; ind < len(si); ind++ {
        stari, _ := strconv.Atoi(si[ind])
        starj, _ := strconv.Atoi(sj[ind])
        if stari < starj {
            return true
        } else if stari > starj {
            return false
        }
    }
    lastSlashi := strings.LastIndex(am[i].Name, "/")
    lastSlashj := strings.LastIndex(am[j].Name, "/")
    return am[i].Name[:lastSlashi] < am[j].Name[:lastSlashj]
}
var signRe = regexp.MustCompile("_[0-9]*[0-9-]*[0-9]")

// toLatex exports results to output as LATEX table rows
func toLatex(output string, results map[string][]AggregatedMeasurement) {
    out, err := os.Create(output)
    if err != nil { glog.Fatal(err) }
    defer out.Close()
    w := bufio.NewWriter(out)
    defer w.Flush()
    header := ""
    cmr, cmrI := "", 3
    for _,v := range results {
        sort.Sort(bySign(v))
        for i, am := range v {
            if i % 2 == 0 {
                sign := signRe.FindString(am.Name)
                w.WriteString(" & \\phantom{a} & \\multicolumn{2}{c}{" + sign[strings.Index(sign, "_")+1:] + "}")
                header += "c@{}rr"
                cmr += "\\cmidrule{" + strconv.Itoa(cmrI) + "-" + strconv.Itoa(cmrI+1) + "}"
                cmrI += 3
            }
        }
        w.WriteString(" \\\\\n")
        break
    }
    w.WriteString(header + "\n")
    w.WriteString(cmr + "\n")
    for k,v := range results {
        sort.Sort(bySign(v))
        w.WriteString(" \\multirow{2}{*}{" + k + "}")
        for i, am := range v {
            if i % 2 == 0 {
                w.WriteString(" & \\phantom{a} & ")
            } else {
                w.WriteString(" & ")
            }
            w.Write(strconv.AppendFloat([]byte{}, am.Jaccard, 'f', 2, 64))
        }
        w.WriteString(" \\\\\n")
        for i, am := range v {
            if i % 2 == 0 {
                w.WriteString(" & \\phantom{a} & ")
            } else {
                w.WriteString(" & ")
            }
            w.Write(strconv.AppendInt([]byte{}, int64(am.ElapsedTime / time.Millisecond), 10))
        }
        w.WriteString(" \\\\\n")
    }
}

// Export evaluates and export in latex format the results of the experiment
// into the output file.
func Export(output string, gold string, results []string) {
    fi, err := os.Open(gold)
    if err != nil { glog.Fatal(err) }
    defer fi.Close()
    goldFiles, err := filepath.Glob(gold + "/query_*")
    if err != nil { glog.Fatal(err) }

    limit := make(map[string][]AggregatedMeasurement)
    limitRe := regexp.MustCompile("limit[0-9]+$")
    for _, file := range goldFiles {
        gd := loadGold(file)
        _, name := filepath.Split(file)
        for _,res := range results {
            matches, err := filepath.Glob(filepath.Join(res, name)  + "-limit*")
            if err != nil { glog.Fatal(err) }
            for _, match := range matches {
                aggm := compare(gd, match)
                l := limitRe.FindString(match)
                limit[l] = append(limit[l], aggm)
            }
        }
    }
    toLatex(output, limit)
}

