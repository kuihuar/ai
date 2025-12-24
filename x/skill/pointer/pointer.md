
1. [3]unitptr 是数组
2. unsafe.Pointer 这是go中用于绕过类型安全的指针类型，通用指针类型，绕过类型系统检查，直接操作内存
3. res := *(*[]byte)(unsafe.Pointer(&b))
 - (*[]byte)(unsafe.Pointer(&b)) 是将unsafe.Pointer类型转化为 *[]byte   
 -  *[]byte 它是指向字节切片的指针
 -  (*[]byte) 指针解引用，得到实际的[]byte 值

4. & 获取亦是的内存地址，返回一个指针
5. * 解引用，通过指针访问指向的内存值
6.  unsafe.Pointer go中通用指针类型，绕过类型系统检查，直接操作内存



// 安全使用uintptr
type MyStruct struct{ Field int }

func SafePointerConversion() {
	obj := &MyStruct{Field: 42}
	// 转换为uintptr,必须确保obj存活
	addr := uintptr(unsafe.Pointer(obj))

	// 模拟其它操作和，在此期间存活，防止GC
	runtime.KeepAlive(obj)
	//转换回指针
	ptr := (*MyStruct)(unsafe.Pointer(addr))
	fmt.Println(ptr.Field)
}


// 安全使用uintptr
type MyStruct struct{ Field int }

func SafePointerConversion() {
	obj := &MyStruct{Field: 42}
	// 转换为uintptr,必须确保obj存活
	addr := uintptr(unsafe.Pointer(obj))

	// 模拟其它操作和，在此期间存活，防止GC
	runtime.KeepAlive(obj)
	//转换回指针
	ptr := (*MyStruct)(unsafe.Pointer(addr))
	fmt.Println(ptr.Field)
}

func wrong1() {
	num := 10
	ptr := unsafe.Pointer(&num)

	floatPtr := (*float64)(ptr)

	fmt.Println(*floatPtr)
}

func testSizeof() {
	var num int = 10

	size := unsafe.Sizeof(num)
	fmt.Printf("int 类型亦是 num 所占字节数： %d\n", size)
}

type Person struct {
	name string
	age  int
}

func testOffset() {
	p := Person{}
	offset := unsafe.Offsetof(p.age)
	fmt.Printf("Person 结构体， age 字段的偏移量:%d\n", offset)
}

func testAlignof() {
	var num int = 10
	align := unsafe.Alignof(num)
	fmt.Printf("int类型亦是num的内存对齐系数是 %d\n", align)
	// unsafe.Slice()
}

func testSlice() {
	var arr [5]int = [5]int{1, 2, 3, 4, 5}
	//从数组的第一个元素创建一个新的切片
	slice := unsafe.Slice(&arr[0], 3) // [1,2,3]
	fmt.Println(slice)

	// unsafe.SliceData 用于获取切片底层数组的指针，
	// 它返回一个指向切片第一个元素的打针，通过这个指针可以直接访问切片的底层内存

	slice1 := []int{1, 2, 3}

	//获取底层数组的指针
	ptr := unsafe.SliceData(slice1)
	//通过指针访问切片的第一个元素
	firtElement := *(*int)(ptr) // 1
	fmt.Println(firtElement)

}

func testString() {
	// unsafe.StringData, 函数用于获取字符串底层字节数组的指针。它返回 个指向字符串第一个字节的指针
	// 通过这个指针可以直接访问字符串的底层内存

	str := "hello"
	// 获取字符串底层字节数组的指针
	ptr := unsafe.StringData(str)
	// 通过指针返回第一个字节
	fistByte := *(*byte)(ptr) // h

	fmt.Printf("%c\n", fistByte)

	// unsafe.String用于从一个字节切片和长度创建一个新的字符串。
	// 它允许你直接从字节切片的内存创建字符串，而无需数据复制

	bytes := []byte{'h', 'e', 'l', 'l', 'o'}

	str1 := unsafe.String(unsafe.SliceData(bytes), len(bytes))

	fmt.Println(str1) // hello

}

func testAdd() {
	alice := Person{name: "alice", age: 30}

	ptr := unsafe.Pointer(&alice)

	ageOffset := unsafe.Offsetof(alice.age)

	agePtr := unsafe.Add(ptr, ageOffset)

	age := (*int)(agePtr)

	fmt.Println(*age)
}

uintptr 和 unsafe.Pointer

uintptr 通常用于需要将指针算术，并将结果存储为整数的场景。
unsafe.Pointer 用于在不同类型的指针之间进行转换，
运行