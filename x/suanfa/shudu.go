package suanfa

//36 37 数独

// 回溯法 + 剪枝
func SolveSudoku(board [][]byte) {
	// 初始化一个 9x9 的布尔数组，用于记录每个数字是否已经在当前行、列或 3x3 子网格中出现
	// 标记当前行、列或 3x3 子网格中已经使用的数字
	rowUsed := make([][]bool, 9)
	colUsed := make([][]bool, 9)
	boxUsed := make([][]bool, 9)

	for i := 0; i < 9; i++ {
		rowUsed[i] = make([]bool, 9)
		colUsed[i] = make([]bool, 9)
		boxUsed[i] = make([]bool, 9)
	}
	// 初始化 rowUsed、colUsed 和 boxUsed 数组
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] != '.' {
				// 将 '1' 到 '9' 转换为 0 到 8 的索引
				num := board[i][j] - '1'
				// 第i行已经存在num+1
				rowUsed[i][num] = true
				// 第i行已经存在num+1
				colUsed[j][num] = true
				// 格子已经存在num+1
				boxUsed[i/3*3+j/3][num] = true
			}
		}
	}
	// 回溯函数
	var backtrack func(row, col int) bool
	backtrack = func(row, col int) bool {
		// 找到下一个未填充的位置，跳转到下一个 '.'的位置
		for row < 9 && board[row][col] != '.' {
			row, col = row+(col+1)/9, (col+1)%9
		}

		// 如果已经填充完所有位置，返回 true， 终止条件
		if row == 9 {
			return true
		}
		// 尝试填充当前位置
		for num := 0; num < 9; num++ {
			// 如果当前数字在当前行、列或 3x3 子网格中已经出现，跳过，这段是剪枝
			if !rowUsed[row][num] && !colUsed[col][num] && !boxUsed[row/3*3+col/3][num] {
				// 将0～8 转换字符 '1'～'9'
				// 尝试填充数字
				board[row][col] = byte(num) + '1'
				rowUsed[row][num] = true
				colUsed[col][num] = true
				boxUsed[row/3*3+col/3][num] = true
				// 递归尝试填充下一个位置（空格）
				if backtrack(row, col) {
					return true
				}
				// 回溯撤消
				board[row][col] = '.'
				rowUsed[row][num] = false
				colUsed[col][num] = false
				boxUsed[row/3*3+col/3][num] = false
			}
		}
		return false
	}
	backtrack(0, 0)
}

func SoloveSudoku2(board [][]byte) {
	rowUsed := make([][]bool, 9)
	colUsed := make([][]bool, 9)
	boxUsed := make([][]bool, 9)

	for i := 0; i < 9; i++ {
		rowUsed[i] = make([]bool, 9)
		colUsed[i] = make([]bool, 9)
		boxUsed[i] = make([]bool, 9)
	}

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] != '.' {
				num := board[i][j] - '1'
				rowUsed[i][num] = true
				colUsed[j][num] = true
				boxUsed[i/3*3+j/3][num] = true
			}
		}
	}
	var backtrack func() bool

	backtrack = func() bool {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if board[i][j] == '.' {
					for num := 0; num < 9; num++ {

						boxIndex := i/3*3 + j/3

						if !rowUsed[i][num] && !colUsed[j][num] && !boxUsed[boxIndex][num] {
							board[i][j] = byte(num) + '1'
							rowUsed[i][num] = true
							colUsed[j][num] = true
							boxUsed[boxIndex][num] = true
							if backtrack() {
								return true
							}
							// 回溯，撤销选择
							board[i][j] = '.'
							rowUsed[i][num] = false
							colUsed[j][num] = false
							boxUsed[boxIndex][num] = false
						}
					}
					// 尝试了所有可能的数字，仍然无法找到解决方案，返回 false
					return false
				}
			}
		}
		return true
	}

	backtrack()
}

func SolveSudoku3(board [][]byte) {
	rowUsed := make([][]bool, 9)
	colUsed := make([][]bool, 9)
	boxUsed := make([][]bool, 9)
	for i := 0; i < 9; i++ {
		rowUsed[i] = make([]bool, 9)
		colUsed[i] = make([]bool, 9)
		boxUsed[i] = make([]bool, 9)
	}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if board[i][j] != '.' {
				num := board[i][j] - '1'
				rowUsed[i][num] = true
				colUsed[j][num] = true
				boxUsed[i/3*3+j/3][num] = true
			}
		}
	}

	var backtrack func(row, col int) bool
	backtrack = func(row int, col int) bool {

		// 两种可能
		// 1. col < 8,未到行尾，
		// row, col = row + 0, col+1  列名移一格，行不变
		// 2. col == 8,到行尾，
		// row, col = row + 1, 0  // 跳到下一行开头
		for row < 9 && board[row][col] != '.' {
			row, col = row+(col+1)/9, (col+1)%9
		}
		// 等效逻辑

		for row < 9 && board[row][col] != '.' {
			if col < 8 {
				col++
			} else {
				row++
				col = 0
			}
		}
		for row < 9 && board[row][col] != '.' {
			col++
			if col >= 9 {
				col = 0
				row++
			}
		}

		for row < 9 && board[row][col] != '.' {
			row, col = row+(col+1)/9, (col+1)%9
		}
		if row == 9 {
			return true
		}

		for num := 0; num < 9; num++ {
			boxIndex := (row/3)*3 + col/3
			if !rowUsed[row][num] && !colUsed[col][num] && !boxUsed[boxIndex][num] {
				board[row][col] = byte(num) + '1'
				rowUsed[row][num] = true
				colUsed[col][num] = true
				boxUsed[boxIndex][num] = true
				if backtrack(row, col) {
					return true
				}
				board[row][col] = '.'
				rowUsed[row][num] = false
				colUsed[col][num] = false
				boxUsed[boxIndex][num] = false
			}
		}

		return false
	}

	backtrack(0, 0)

}

func isValidSudoku(board [][]byte) bool {

	var rows [9][9]int        // 记录每行数字1-9的出现次数
	var columns [9][9]int     // 记录每列数字1-9的出现次数
	var subboxes [3][3][9]int // 记录每个3x3小宫格数字1-9的出现次数

	for i, row := range board {
		for j, c := range row {
			if c == '.' {
				continue
			}
			// 数字转换
			// 例如：'5' → '5'-'1' = 4
			numIndex := c - '1'
			rows[i][numIndex]++
			columns[j][numIndex]++
			subboxes[i/3][j/3][numIndex]++

			if rows[i][numIndex] > 1 || columns[j][numIndex] > 1 || subboxes[i/3][j/3][numIndex] > 1 {
				return false
			}

		}
	}

	return true
}
