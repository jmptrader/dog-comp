package util

type Graph struct {
	nodes       []*Node
	gname       string
	Node_String func(interface{}) string //print Node
}

type Node struct {
	data        interface{}
	edges       []*Edge
	indegree    int
	outdegree   int
	succ        map[*Node]bool //succ set
	prev        map[*Node]bool //prev set
	Node_String func(interface{}) string
}

func (this *Graph) Node_new(data interface{}) *Node {
	o := new(Node)
	o.data = data
	o.edges = make([]*Edge, 0)
	o.indegree = 0
	o.outdegree = 0
	o.succ = make(map[*Node]bool)
	o.prev = make(map[*Node]bool)
	o.Node_String = this.Node_String
	return o
}
func (this *Node) GetSucc() map[*Node]bool {
	return this.succ
}
func (this *Node) String() string {
	if this.Node_String != nil {
		return this.Node_String(this.data)
	} else {
		return ""
	}
}
func (this *Node) GetData() interface{} {
	return this.data
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
	if from.succ[to] == false {
		from.succ[to] = true
	}
	if to.prev[from] == false {
		to.prev[from] = true
	}
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

func (this *Graph) Quasi_reverse() []*Node {
	nodes := make([]*Node, 0)
	visited := make(map[*Node]bool)
	var dfs func(*Node)

	//use DFS
	dfs = func(n *Node) {
		if visited[n] == true {
			return
		}
		visited[n] = true
		for _, e := range n.edges {
			dfs(e.to)
		}
		nodes = append(nodes, n)
	}

	dfs(this.nodes[0])
	//check
	if len(nodes) != len(this.nodes) {
		panic("impossible")
	}
	return nodes
}
