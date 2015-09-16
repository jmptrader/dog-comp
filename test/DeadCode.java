class DeadCode { 
	public static void main(String[] a) {
        System.out.println(new Doit().doit());
    }
}

class Doit {
    public int doit() {
    	int i;
    	i=1;
        if (!!!!!!true)
          System.out.println(111);
        else 
          System.out.println(222);
        while(i<3)
        {
        	System.out.println(333);
        	i=i+1;
        }
        while(false&&true)
        {
            System.out.println(444);
        }
        return 0;
    }
}
