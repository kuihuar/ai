package suanfa

import (
	"math/bits"
	"strings"
)

// 51, 52 N-Queens

func SolveNQueens(n int) [][]string {
	// 结果集
	var result [][]string
	// 皇后的攻击位置
	// 保存皇后位置
	board := make([][]string, n)
	attack := make([][]bool, n)
	// 初始化棋盘和攻击范围
	for i := 0; i < n; i++ {
		board[i] = make([]string, n)
		attack[i] = make([]bool, n)
		for j := 0; j < n; j++ {
			board[i][j] = "."
			attack[i][j] = false
		}
	}

	var putQueen = func(board [][]string, row, col int) {
		dx := [8]int{-1, 1, 0, 0, -1, -1, 1, 1}
		dy := [8]int{0, 0, -1, 1, -1, 1, -1, 1}

		attack[row][col] = true // 当前位置放置皇后

		// 标记攻击范围,通过两循环将皇后可能攻击的范围标记
		// 这里i是从1开始的，因为皇后的攻击范围是从皇后位置向1到n-1的距离延伸
		for i := 1; i < len(attack); i++ { // 从皇后位置向1到n-1的距离延伸，多少行
			for j := 0; j < 8; j++ { //遍历8个方向
				newX := row + i*dx[j] //生成新的位置行
				newY := col + i*dy[j] // 生成新的位置列
				// 在横盘范围内，将新的位置标记为true
				if newX >= 0 && newX < len(attack) && newY >= 0 && newY < len(attack) {
					attack[newX][newY] = true
				}
			}
		}

		// 1. 匹配到一次，将当前成功的横盘位置board[row][col]保存到结果里
		board[row][col] = "Q"
	}

	// backtrack 是一个闭包函数，所以它可以访问外部的result变量,所以不应该将result作为参数传递
	var backtrack func(row int)

	// 回溯函数
	// 1. 终止条件：row == n
	// board 表示当前的棋盘状态,存储皇后的位置
	// attack 表示当前的攻击范围
	// row 表示当前处理的行
	// n 表示皇后的数量
	// queen 存储皇后的位置
	// result 存储所有的解法
	backtrack = func(row int) {
		// 终止条件：如果已经处理完所有行，将当前的解加入结果中
		if row == n {
			// 如果已经处理完所有行，将当前的解加入结果中
			solution := make([]string, n)
			for i := range board {
				solution[i] = strings.Join(board[i], "")
			}
			result = append(result, solution)
			return
		}
		// 遍历当前行的所有列，回溯试探皇后可知放置的位置
		for col := 0; col < n; col++ {
			// 如果当前位置可以放置皇后
			// 备份attach数组
			if !attack[row][col] {
				tempAttack := make([][]bool, n)
				for i := range attack {
					tempAttack[i] = append([]bool{}, attack[i]...)
				}
				// 放置皇后修改横盘
				board[row][col] = "Q"
				//更新attack数组
				// 放置皇后修改横盘
				putQueen(board, row, col)
				// 递归处理下一行
				backtrack(row + 1)
				// 回溯，恢复attack数组
				attack = tempAttack
				board[row][col] = "."
			}

		}

	}

	//回溯法递归解决N皇后问题
	backtrack(0)

	return result

}

// 使用三个一维数组来记录皇后的攻击范围
// 1. 行攻击范围，不需要记录，因为每行只能有一个皇后
// 2. 列攻击范围 colAttack[n]bool 表示第i列是否有皇后
// 3. 对角线攻击范围 diagonalAttack[2n-1]bool 表示第i条对角线是否有皇后
// 4. 反对角线攻击范围 antiDiagonalAttack[2n-1]bool 表示第i条反对角线是否有皇后
// 5. 回溯函数 backtrack(row int)
// 4. 反对角线攻击范围

// 对角线类型	索引计算	取值范围	数组长度
// 主对角线	row - col + n - 1	[0, 2n-2]	2n-1
// 副对角线	row + col	[0, 2n-2]	2n-1
func SolveNQueens1(n int) [][]string {
	// 结果集
	var result [][]string
	// 皇后的攻击位置
	// 保存皇后位置
	board := make([][]string, n)
	// 主对角线通过row - col加上n-1来避免负数索引，
	// 而副对角线则直接使用row + col作为索引，
	// 这样都能确保每个对角线有唯一的索引值，不会有冲突
	// 列攻击范围
	colAttack := make([]bool, n) // 列攻击标记
	// 主对角线攻击范围
	// 边界情况	计算示例
	// 最小值	0 - (n-1) = -(n-1)（左上角）
	// 最大值	(n-1) - 0 = n-1（右下角）
	// 因此，row - col 的可能取值为 [-(n-1), n-1]，共 2n-1 个不同的值。
	// 为了将这些值映射到数组索引，我们通过 row - col + n - 1 将范围调整为 [0, 2n-2]，对应数组长度 2n-1
	// 主对角线方向为从左上到右下（↘），同一主对角线上所有位置的 row - col 的值相等。
	// row - col 的可能取值为 [-(n-1), n-1]，共 2n-1 个不同的值
	// 为了将这些值映射到数组索引，
	// 我们通过 row - col + n - 1 将范围调整为 [0, 2n-2]，对应数组长度 2n-1
	diagonalAttack := make([]bool, 2*n-1) // 主对角线攻击标记（row-col）
	// 副对角线方向为从右上到左下（↙），同一副对角线上所有位置的 row + col 的值相等。
	// 对于 n×n 的棋盘，row + col 的可能取值范围为：
	//边界情况	计算示例
	// 最小值	0 + 0 = 0（左上角）
	// 最大值	(n-1) + (n-1) = 2n-2（右下角）
	// 因此，row + col 的可能取值为 [0, 2n-2]，共 2n-1 个不同的值。
	// 直接使用 row + col 作为索引即可，数组长度自然为 2n-1
	antiDiagonalAttack := make([]bool, 2*n-1) // 副对角线攻击标记（row+col）
	// 初始化棋盘
	// 攻击范围不用初始化
	for i := 0; i < n; i++ {
		board[i] = make([]string, n)
		for j := range board[i] {
			// 预填充横盘有
			board[i][j] = "."
		}
	}

	// backtrack 是一个闭包函数，所以它可以访问外部的result变量,所以不应该将result作为参数传递
	var backtrack func(row int)

	// 回溯函数
	// 1. 终止条件：row == n
	// board 表示当前的棋盘状态,存储皇后的位置
	// attack 表示当前的攻击范围
	// row 表示当前处理的行
	// n 表示皇后的数量
	// queen 存储皇后的位置
	// result 存储所有的解法
	backtrack = func(row int) {
		// 终止条件：如果已经处理完所有行，将当前的解加入结果中
		if row == n {
			// 如果已经处理完所有行，将当前的解加入结果中
			solution := make([]string, n)
			for i := range board {
				solution[i] = strings.Join(board[i], "")
			}
			result = append(result, solution)
			return
		}
		// 遍历当前行的所有列，回溯试探皇后可知放置的位置
		for col := 0; col < n; col++ {
			// 计算对角线索引
			diagonalIndex := row - col + n - 1
			antiDiagonalIndex := row + col
			if colAttack[col] || diagonalAttack[diagonalIndex] || antiDiagonalAttack[antiDiagonalIndex] {
				continue
			}

			// 放置皇后修改横盘
			board[row][col] = "Q"

			// 更新攻击范围
			colAttack[col] = true
			diagonalAttack[diagonalIndex] = true
			antiDiagonalAttack[antiDiagonalIndex] = true
			// 如果当前位置可以放置皇后
			// 备份attach数组
			// 递归处理下一行
			backtrack(row + 1)
			// 回溯，恢复attack数组
			board[row][col] = "."
			colAttack[col] = false
			diagonalAttack[diagonalIndex] = false
			antiDiagonalAttack[antiDiagonalIndex] = false

		}

	}

	//回溯法递归解决N皇后问题
	backtrack(0)

	return result

}

func SolveNQueens2(n int) [][]string {
	// 结果集
	result := make([][]string, 0)
	board := make([][]string, n)

	// 初始化棋盘
	for i := range board {
		board[i] = make([]string, n)
		for j := range board[i] {
			board[i][j] = "."
		}
	}

	cols := make([]bool, n)
	diag1 := make([]bool, 2*n-1)
	diag2 := make([]bool, 2*n-1)
	// 如果闭包递归调用自身，必须使用 var 先声明函数，再赋值
	// 否则会出现编译错误：undefined: backtrack
	// 因为 Go 语言的闭包特性，backtrack 函数在声明时还没有被定义
	// 所以在 backtrack 函数内部调用 backtrack 时，会出现编译错误
	// 解决方法是使用 var 先声明函数，再赋值
	var backtrack func(int)
	backtrack = func(row int) {
		if row == n {
			solution := make([]string, n)
			for i := range board {
				solution[i] = strings.Join(board[i], "")
			}
			result = append(result, solution)
			return
		}
		for col := 0; col < n; col++ {
			// 计算对角线索引
			d1 := row - col + n - 1
			d2 := row + col
			// 有可能之前的皇后已经占用了这个位置，所以要跳过
			if cols[col] || diag1[d1] || diag2[d2] {
				continue
			}
			// 放置皇后
			board[row][col] = "Q"
			cols[col] = true
			diag1[d1] = true
			diag2[d2] = true
			// 递归到下一行
			backtrack(row + 1)
			// 回溯
			board[row][col] = "."
			cols[col] = false
			diag1[d1] = false
			diag2[d2] = false
		}
	}

	backtrack(0)
	return result
}

func SolveNQueensLast(n int) [][]string {
	var result [][]string

	board := make([][]string, n)
	colAttack := make([]bool, n)
	diagonalAttack := make([]bool, 2*n-1)
	antiDiagonalAttack := make([]bool, 2*n-1)
	for i := 0; i < n; i++ {
		board[i] = make([]string, n)
		for j := range board[i] {
			board[i][j] = "."
		}
	}

	var backtrack func(row int)
	backtrack = func(row int) {
		if row == n {
			solution := make([]string, n)
			for i := range board {
				solution[i] = strings.Join(board[i], "")
			}
			result = append(result, solution)
			return
		}
		for col := 0; col < n; col++ {
			diagonalIndex := row - col + n - 1
			antiDiagonalIndex := row + col
			if colAttack[col] || diagonalAttack[diagonalIndex] || antiDiagonalAttack[antiDiagonalIndex] {
				continue
			}
			board[row][col] = "Q"
			colAttack[col] = true
			diagonalAttack[diagonalIndex] = true
			antiDiagonalAttack[antiDiagonalIndex] = true

			backtrack(row + 1)
			board[row][col] = "."
			colAttack[col] = false
			diagonalAttack[diagonalIndex] = false
			antiDiagonalAttack[antiDiagonalIndex] = false

		}

	}

	backtrack(0)
	return result
}

// 复杂了
func SolveNQueensBits(n int) [][]string {
	var result [][]string
	board := make([][]string, n)
	for i := 0; i < n; i++ {
		board[i] = make([]string, n)
		for j := range board[i] {
			board[i][j] = "."
		}
	}
	var backtrack func(row, cols, diagonals1, diagonals2 int)
	backtrack = func(row, cols, diagonals1, diagonals2 int) {
		if row == n {
			solution := make([]string, n)
			for i := range board {
				solution[i] = strings.Join(board[i], "")
			}
			result = append(result, solution)
			return
		}
		availablePositions := ((1 << n) - 1) &^ (cols | diagonals1 | diagonals2)
		for availablePositions != 0 {
			// 获取最低位的1
			position := availablePositions & -availablePositions
			// 计算列索引 最低位1的位置
			col := bits.TrailingZeros(uint(position))
			board[row][col] = "Q"
			// 递归到下一行
			backtrack(row+1,
				cols|position,
				(diagonals1|position)<<1, //主对角线线右移
				(diagonals2|position)>>1) //副对角线左移
			//去掉最低位的1
			availablePositions = availablePositions & (availablePositions - 1)
			board[row][col] = "."

		}
	}
	backtrack(0, 0, 0, 0)
	return result
}
func SolveNQueensBits2(n int) int {
	var result int

	// 	手动参数传递更优，不需要手动管理状态
	// 如果定义cols, diagonals1, diagonals2作为局部变量，忘记恢复状态，会导致错误
	// 定义成函数变量，每次调用时都会重新创建，不需要手动管理状态，没有额外开销
	var backtrack func(row, cols, diagonals1, diagonals2 int)
	backtrack = func(row, cols, diagonals1, diagonals2 int) {
		if row == n {
			result++
			return
		}
		availablePositions := ((1 << n) - 1) &^ (cols | diagonals1 | diagonals2)
		for availablePositions != 0 {
			// 获取最低位的1
			position := availablePositions & -availablePositions
			// 计算列索引 最低位1的位置
			// 递归到下一行
			backtrack(row+1,
				cols|position,
				(diagonals1|position)<<1, //主对角线线右移
				(diagonals2|position)>>1) //副对角线左移
			//去掉最低位的1
			availablePositions = availablePositions & (availablePositions - 1)

		}
	}
	backtrack(0, 0, 0, 0)
	return result
}
