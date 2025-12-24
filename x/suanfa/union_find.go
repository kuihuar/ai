package suanfa

type UnionFind struct {
	parent []int
	rank   []int
}

func (u *UnionFind) Find(x int) int {
	if u.parent[x] != x {
		u.parent[x] = u.Find(u.parent[x])
	}
	return u.parent[x]
}

func (u *UnionFind) Union(x, y int) {
	rootX := u.Find(x)
	rootY := u.Find(y)

	if rootX == rootY {
		return
	}
	if u.rank[rootX] > u.rank[rootY] {
		u.parent[rootY] = rootX
	} else if u.rank[rootX] < u.rank[rootY] {
		u.parent[rootX] = rootY
	} else {
		u.parent[rootY] = rootX
		u.rank[rootX]++
	}
}
func (u *UnionFind) IsConnected(x, y int) bool {
	return u.Find(x) == u.Find(y)
}

// 200岛屿数量
// 1. 并查集
// 2. 深度优先搜索(染色)
// 3. 广度优先搜索(染色)

// 547 friend circles

// func numIslands(grid [][]byte) int {
// 	m := len(grid)
// 	n := len(grid[0])
// 	uf := NewUnionFind(m * n)
// 	for i := 0; i < m; i++ {
// 		for j := 0; j < n; j++ {
// 			if grid[i][j] == '1' {
// 				if i > 0 && grid[i-1][j] == '1' {
// 					uf.Union(i*n+j, (i-1)*n+j)
// 				}
// 				if j > 0 && grid[i][j-1] == '1' {
// 					uf.Union(i*n+j, i*n+j-1)
// 				}
// 			}
// 		}
// 	}
// 	count := 0
// 	for i := 0; i < m; i++ {
// 		for j := 0; j < n; j++ {
// 			if grid[i][j] == '1' {
// 				count++
// 			}
// 		}
// 	}
// 	for i := 0; i < m; i++ {
// 		for j := 0; j < n; j++ {
// 			if grid[i][j] == '1' {
// 				if uf.IsConnected(i*n+j, 0) {
// 					count--
// 				}
// 			}
// 		}
// 	}
// 	return count
// }
