
type Sum struct {
}

func main() {
    fmt.Println(new(Sum).sub(10,15))
}
func (this *Sum) sum(c int, a int, b int) int {

    if c <b {
	sum := a + b * c
    }else {
        sum:=c
    }
	return sum

}

func (this *Sum) sub(i int, j int) bool {
    for i&&j{
        ret := true
    }
    return ret
}
