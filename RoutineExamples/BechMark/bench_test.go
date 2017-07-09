package BechMark

import "testing"

func BenchmarkForBenchChan(b *testing.B) {
    for n:= 0; n < b.N; n++ {
        ForBenchChan()
    }
}



func BenchmarkForBenchMutext(b *testing.B) {
    for n:= 0; n < b.N; n++ {
        ForBenchMutext()
    }
}