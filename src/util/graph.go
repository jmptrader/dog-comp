package util

type Graph struct {
	nodes       []*Node
	gname       string
	Node_String func(interface{}) string
}

type Node struct {
	data        interface{}
	edges       []*Edge
	indegree    int
	outdegree   int
	Node_String func(interface{}) string
}

func (this *Graph) Node_new(data interface{}) *Node {
	o := new(Node)
	o.data = data
	o.edges = make([]*Edge, 0)
	o.indegree = 0
	o.outdegree = 0
	o.Node_String = this.Node_String
	return o
}
func (this *Node) String() string {
	if this.Node_String != nil {
		return this.Node_String(this.data)
	} else {
		return ""
	}
}

type Edge struct {
	from *Node
	to   *Node
}

func (this *Edge) String() string {
	return ""
}

func Graph_new(name string) *Graph {
	o := new(Graph)
	o.gname = name
	o.nodes = make([]*Node, 0)
	return o
}

func (this *Graph) LookupNode(data interface{}) *Node {
	for _, n := range this.nodes {
		if n.data == data {
			return n
		}
	}
	return nil
}

func (this *Graph) AddNode(data interface{}) {
	if this.LookupNode(data) != nil {
		panic("dup graph data")
	}
	this.nodes = append(this.nodes, this.Node_new(data))
}

func (this *Graph) addEdge(from *Node, to *Node) {
	from.outdegree++
	to.indegree++
	from.edges = append(from.edges, &Edge{from, to})
}

func (this *Graph) AddEdge(from interface{}, to interface{}) {
	f := this.LookupNode(from)
	t := this.LookupNode(to)
	if f == nil || t == nil {
		panic("from || to is nil")
	}
	this.addEdge(f, t)
}

func (this *Graph) Visualize() {
	dot := Dot_new()
	fname := this.gname
	for _, n := range this.nodes {
		if n.indegree == 0 && n.outdegree == 0 {
			dot.InsertOne(n.String())
		}
		for _, e := range n.edges {
			//dot.Insert(e.from.String(), e.to.String())
			dot.Insert(e.from.String(), e.to.String())
		}
	}
	dot.Visualize(fname)
}
