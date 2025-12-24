package suanfa

import (
	"fmt"
)

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func NewTreeNode(val int) *TreeNode {
	return &TreeNode{Val: val}
}

// 1. InorderTraversal 中序遍历后，升序，可以只保留前一个节点的值，与当前节点的值进行比较
// 2. 中序遍历的时间复杂度为O(n)，空间复杂度为O(n)
// 2. Recursion 递归实现
// 3. 递归实现的时间复杂度为O(n)，空间复杂度为O(n)

// 98. 验证二叉搜索树
// 递归实现
// 1. 递归实现的时间复杂度为O(n)，空间复杂度为O(n),因为用到了递归栈
func ValidateBinarySearchTree(root *TreeNode) bool {
	fmt.Println("ValidateBinarySearchTree")

	return helper(root, nil, nil)
	// return helper1(root, -1<<63, 1<<63-1)
	// return helper1(root, math.MinInt64, math.MaxInt64)
}

func helper(root *TreeNode, lower, upper *TreeNode) bool {
	if root == nil {
		return true
	}
	if lower != nil && root.Val <= lower.Val {
		return false
	}
	if upper != nil && root.Val >= upper.Val {
		return false
	}
	return helper(root.Left, lower, root) && helper(root.Right, root, upper)
}

func helper1(root *TreeNode, lower, upper int) bool {
	if root == nil {
		return true
	}
	if root.Val <= lower || root.Val >= upper {
		return false
	}
	// 左子树： helper1(root.Left, lower, root.Val)   root.Left < root.Val
	// - root.Left：左子树的根节点
	// - lower：左子树的下限，即父节点的值，这个是下限的值，是继承过来的
	// - root.Val：左子树的上限，即父节点的值，这个是上限的值，要改的

	// 右子树： helper1(root.Right, root.Val, upper) root.Right > root.Val
	// - root.Right：右子树的根节点
	// - root.Val：右子树的下限，即父节点的值，这个是下限的值，要改的
	// - upper：右子树的上限，即父节点的值，这个是上限的值，是继承过来的

	// 总结
	// 左子树： 上限更新为父节点值，下限继承父节点的下限
	// 右子树： 下限更新为父节点值，上限继承父节点的上限

	// 左子树：继承父节点的下限 lower，上限设为当前节点值 root.Val（左子树所有节点值必须严格小于 root.Val）。
	// 右子树：下限设为当前节点值 root.Val，继承父节点的上限 upper（右子树所有节点值必须严格大于 root.Val）

	return helper1(root.Left, lower, root.Val) && helper1(root.Right, root.Val, upper)
}

// 2. 迭代实现
// 1. 迭代实现的时间复杂度为O(n)，空间复杂度为O(n)
func ValidateBinarySearchTree1(root *TreeNode) bool {
	if root == nil {
		return true
	}
	stack := []struct {
		node         *TreeNode
		lower, upper int
	}{
		{root, -1 << 63, 1<<63 - 1},
	}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if current.node == nil {
			continue
		}
		if current.node.Val <= current.lower || current.node.Val >= current.upper {
			return false
		}
		stack = append(stack, struct {
			node  *TreeNode
			lower int
			upper int
		}{
			node:  current.node.Left,
			lower: current.lower,
			upper: current.node.Val,
		})
		stack = append(stack, struct {
			node  *TreeNode
			lower int
			upper int
		}{
			node:  current.node.Right,
			lower: current.node.Val,
			upper: current.upper,
		})
	}
	return true
}

// 3. 中序遍历实现
// 左子树，当前节点，右子树， 每个元素都比前一个元素大
// 1. 中序遍历的时间复杂度为O(n)，空间复杂度为O(n)，因为用到了栈
// 先遍历左子树，然后遍历当前节点，最后遍历右子树
// 中序遍历以后得到的是一个升序的序列
func ValidateBinarySearchTree2(root *TreeNode) bool {

	var result []int

	var inorder func(root *TreeNode)

	inorder = func(root *TreeNode) {
		if root == nil {
			return
		}
		inorder(root.Left)
		result = append(result, root.Val)
		inorder(root.Right)
	}
	inorder(root)
	return true
}

// 中序遍历实现
func ValidateBinarySearchTree3(root *TreeNode) bool {
	var inorder func(root *TreeNode) bool
	var prev *TreeNode
	inorder = func(node *TreeNode) bool {
		if node == nil {
			return true
		}
		// 遍历左子树（递归）
		if !inorder(node.Left) {
			return false
		}
		// 检查当前节点是否大于前驱节点
		if prev != nil && node.Val <= prev.Val {
			return false
		}
		// 更新 prev 为当前节点，继续递归遍历右子树
		prev = node
		if !inorder(node.Right) {
			return false
		}
		return true
	}
	return inorder(root)

}

// 最近公共祖先(236)
// 1. 递归实现
// 树的问题，一般都是递归实现
// 1. 定义子问题，子树包信PQ节点

// - 1. 左子树或右子树包含PQ节点
// - 2. 当前节点是P或Q节点，且左子树或右子树包含另一个节点
// 时间复杂度：O(n)，其中 n 是二叉树的节点数。在递归遍历二叉树时，每个节点最多被访问一次。

func LowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
	fmt.Printf("find root: %+v, p: %v, q: %v \n", root, p, q)
	// 1. 递归终止条件
	// 1. 当root为空时，返回nil

	if root == nil {
		return root
	}
	// 2. 递归终止条件
	// 2. 当root等于p或q时，返回root
	// 接下来，如果p等于root或者q等于root，就直接返回root。
	// 这应该意味着，如果当前节点是p或q中的一个，那么这个节点可能就是LCA。
	// 例如，如果p在q的子树中，那么p就是两者的LCA
	if p.Val == root.Val || q.Val == root.Val {
		fmt.Printf("finded root: %+v, p: %v, q: %v \n", root, p, q)
		return root
	}
	left := LowestCommonAncestor(root.Left, p, q)
	right := LowestCommonAncestor(root.Right, p, q)

	//假设p在左子树，q在右子树，那么left递归会找到p，right递归会找到q。
	//这时候，left和right都不为nil，所以返回root作为它们的LCA。这符合LCA的定义，因为root是同时包含p和q的最深节点。
	if left == nil {
		fmt.Printf("find left eq nil, return right: %+v\n", right)
		return right
	}
	if right == nil {
		fmt.Printf("find right eq nil, return left: %+v\n", left)
		return left
	}
	// left != nil && right != nil
	fmt.Printf("finded left and right eq nil return root: %+v\n", root)
	return root
}

// 二叉搜索树最近公共祖先(235)
// 1. 递归实现
func LowestCommonAncestor1(root, p, q *TreeNode) *TreeNode {

	if p.Val < root.Val && q.Val < root.Val {
		return LowestCommonAncestor1(root.Left, p, q)
	}
	if p.Val > root.Val && q.Val > root.Val {
		return LowestCommonAncestor1(root.Right, p, q)
	}
	return root
}

// 二叉搜索树最近公共祖先(235)
// 1. 迭代实现

func LowestCommonAncestor2(root, p, q *TreeNode) *TreeNode {

	for root != nil {
		if p.Val < root.Val && q.Val < root.Val {
			root = root.Left
		}
		if p.Val > root.Val && q.Val > root.Val {
			root = root.Right
		}
		return root
	}
	return nil
}

func PreOrderTraversal(root *TreeNode) {
	var traversePath []int
	if root == nil {
		return
	}
	fmt.Printf("root: %v \n", root.Val)
	traversePath = append(traversePath, root.Val)
	PreOrderTraversal(root.Left)
	PreOrderTraversal(root.Right)
}

func InOrderTraversal(root *TreeNode) {
	var traversePath []int
	if root == nil {
		return
	}
	InOrderTraversal(root.Left)
	fmt.Printf("root: %v \n", root.Val)
	traversePath = append(traversePath, root.Val)
	InOrderTraversal(root.Right)
}
func PostOrderTraversal(root *TreeNode) {
	var traversePath []int
	if root == nil {
		return
	}
	PostOrderTraversal(root.Left)
	PostOrderTraversal(root.Right)
	fmt.Printf("root: %v \n", root.Val)
	traversePath = append(traversePath, root.Val)
}
