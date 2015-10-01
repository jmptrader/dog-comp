
class CopyProp {
    public static void main(String[] a) {
	System.out.println(new Doit().doit());
    }
}

class Doit {
    public int doit() {
	int i;
	int j;
	int k;
	int q;
	int p;
	int x;
	int y;
	int z;
	
	j=2;
	i = j;
	k = i;
	q = k;
	p = q;
	x = p;
	y = x;
	z = y;
	j = z + z;
	if (true)
	    System.out.println(111);
	else
	    System.out.println(222);
	while (i < 3) {
	    System.out.println(333);
	    i = p + 1;
	    j = i +1 ;
        q = j;
	    k = q;//this should be opt -> k=j;
	    p = k+j;
        i = k+j;//this should be opt by subexp
	}
	return 0;
    }
}
