
class ConstProp {
    public static void main(String[] a) {
	System.out.println(new Doit().doit());
    }
}

class Doit {
    public int doit() {
	int i;
	int j;
    j=2;
	i = 1;
	j = j + i;

    System.out.println(j);
	if (true)
	    System.out.println(111);
	else
	    System.out.println(222);
	while (i < 3) {
	    System.out.println(333);
	    i = i + 1;
	}
	return 0;
    }
}
