# cna
Repo for small programs during the course

## Go code

Install go from here: <https://golang.org/dl/>

then `go run $file` or build it to save one second `go build .` then `./$myexec -flag 1 flag 2`

### Examples

```bash
$time go run . -init-port 10 -end-port 1600
Going to probe all the IPs within this network 192.168.0.0 /24
From port 10 to 1600
### Bingo:  192.168.0.XXX 22 is listening! ###
Done

real    0m2.352s
user    0m7.020s
sys     0m1.770s
```

Or

```bash
$time go build .

real    0m0.242s
user    0m0.387s
sys     0m0.097s

$./portSweep -init-port 10 -end-port 1600
Going to probe all the IPs within this network 192.168.0.0 /24
From port 10 to 1600
### Bingo:  192.168.0.XXX 22 is listening! ###
Done

real    0m2.286s
user    0m7.249s
sys     0m1.960s
```
