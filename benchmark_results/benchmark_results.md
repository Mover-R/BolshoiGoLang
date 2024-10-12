goos: linux
goarch: amd64
pkg: BolshiGoLang/benchmark_results
cpu: AMD Ryzen 5 5600H with Radeon Graphics         
BenchmarkGet
BenchmarkGet/10
BenchmarkGet/10-12      39464103                29.51 ns/op            0 B/op          0 allocs/op
BenchmarkGet/100
BenchmarkGet/100-12     33210428                33.12 ns/op            0 B/op          0 allocs/op
BenchmarkGet/1000
BenchmarkGet/1000-12    27404025                42.78 ns/op            0 B/op          0 allocs/op
BenchmarkGet/10000
BenchmarkGet/10000-12   22341577                52.76 ns/op            0 B/op          0 allocs/op

BenchmarkSet
BenchmarkSet/10
BenchmarkSet/10-12       2331159               512.2 ns/op           150 B/op          5 allocs/op
BenchmarkSet/100
BenchmarkSet/100-12      2273115               520.8 ns/op           150 B/op          5 allocs/op
BenchmarkSet/1000
BenchmarkSet/1000-12     2251317               543.2 ns/op           150 B/op          5 allocs/op
BenchmarkSet/10000
BenchmarkSet/10000-12    2054314               584.5 ns/op           151 B/op          5 allocs/op

BenchmarkSetGet
BenchmarkSetGet/10
BenchmarkSetGet/10-12    2238235               543.4 ns/op           144 B/op          4 allocs/op
BenchmarkSetGet/100
BenchmarkSetGet/100-12           2129730               568.6 ns/op           153 B/op          4 allocs/op
BenchmarkSetGet/1000
BenchmarkSetGet/1000-12          2075341               583.0 ns/op           153 B/op          4 allocs/op
BenchmarkSetGet/10000
BenchmarkSetGet/10000-12         1891461               636.0 ns/op           154 B/op          4 allocs/op
PASS
ok      BolshiGoLang/benchmark_results  20.718s