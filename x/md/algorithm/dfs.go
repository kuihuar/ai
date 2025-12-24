package algorithm

import "fmt"

// TreeNode 定义二叉树节点结构
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// preorderTraversal 实现前序遍历
func preorderTraversal(root *TreeNode) []int {
	var result []int
	var dfs func(node *TreeNode)
	dfs = func(node *TreeNode) {
		if node == nil {
			return
		}
		// 先访问根节点
		result = append(result, node.Val)
		// 递归访问左子树
		dfs(node.Left)
		// 递归访问右子树
		dfs(node.Right)
	}
	dfs(root)
	return result
}

func main() {
	// 构建一个简单的二叉树
	root := &TreeNode{Val: 1}
	root.Left = &TreeNode{Val: 2}
	root.Right = &TreeNode{Val: 3}
	root.Left.Left = &TreeNode{Val: 4}
	root.Left.Right = &TreeNode{Val: 5}

	// 进行前序遍历
	result := preorderTraversal(root)
	fmt.Println(result)
}
