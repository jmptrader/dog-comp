package main

type Sum struct {
}

func (this *Sum) Sum(a int, b int) int {

	sum := a + b
	return sum

}

func main() {
	fmt.Println(new(Sum).Sum(10, 15))
}
