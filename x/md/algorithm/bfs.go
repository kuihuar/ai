package algorithm

type Object struct {
	refs  []*Object
	color string // "white", "grey", "black"
}

func BFSMark(root *Object) {
	queue := []*Object{root}
	root.color = "grey"

	for len(queue) > 0 {
		//取出头问对象
		obj := queue[0]
		queue = queue[1:]

		for _, child := range obj.refs {
			if child.color == "white" {
				// 将白色子对象标记为灰色，并加入队列
				child.color = "grey"
				queue = append(queue, child)
			}
		}
		// 将当前对象标记为黑色
		obj.color = "black"
	}
}
