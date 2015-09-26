package cfg

import (
	"../util"
)

func Quasi_reverse(g *util.Graph) []Block {
	f_blocks := make([]Block, 0)
	nodes := g.Quasi_reverse()

	for _, n := range nodes {
		bb := n.GetData()
		if b, ok := bb.(Block); ok {
			f_blocks = append(f_blocks, b)
		} else {
			panic("impossible")
		}
	}

	return f_blocks
}

func GenGraph(mm Method) *util.Graph {
	var f_blocks map[*util.Label]Block
	var graph *util.Graph
	var f_b Block

	createMap := func(bb Block) {
		switch b := bb.(type) {
		case *BlockSingle:
			f_blocks[b.Label_id] = b
		default:
			panic("impossible")
		}
	}

	graph_Edge := func(bb Block) {
		do_Transfer := func(tt Transfer) {
			switch t := tt.(type) {
			case *Goto:
				to := f_blocks[t.Label_id]
				graph.AddEdge(f_b, to)
			case *If:
				to1 := f_blocks[t.Truee]
				to2 := f_blocks[t.Falsee]
				graph.AddEdge(f_b, to1)
				graph.AddEdge(f_b, to2)
			case *Return:
				//no need
			default:
				panic("impossible")
			}
		}

		switch b := bb.(type) {
		case *BlockSingle:
			f_b = b
			do_Transfer(b.Trans)
		default:
			panic("impossible")
		}
	}

	graph_Node := func(bb Block) {
		graph.AddNode(bb)
	}

	switch m := mm.(type) {
	case *MethodSingle:
		f_blocks = make(map[*util.Label]Block)
		graph = util.Graph_new(m.ClassId + "_" + m.Name)
		graph.Node_String = PP
		for _, b := range m.Blocks {
			createMap(b)
		}
		for _, b := range m.Blocks {
			graph_Node(b)
		}
		for _, b := range m.Blocks {
			graph_Edge(b)
		}
	default:
		panic("impossible")
	}

	return graph
}
