1. 数据移动
MOVQ（移动 64 位数据）
作用：在寄存器和内存之间复制数据。
LEAQ（加载有效地址）
作用：计算内存地址并存入寄存器，不实际访问内存。
2. 算术与逻辑指令
ADDQ / SUBQ（加法/减法）
INCQ / DECQ（自增/自减）
3. 控制流指令
CALL（函数调用）
作用：调用函数（如 Go 函数或运行时方法）。
JMP / JEQ / JNE（跳转指令）
作用：条件或无条件跳转。
4. 栈操作指令
PUSHQ / POPQ（压栈/弹栈）
5. 测试与比较指令
TESTQ / CMPQ（测试/比较）

6. Go 汇编特有指令
FUNCDATA / PCDATA
作用：存储调试或垃圾回收信息（由编译器生成）。
NOSPLIT / NOFRAME（函数属性）
作用：标记函数不需要栈分裂或帧指针


7. 接口与类型相关操作
ITAB（接口表）
作用：存储接口类型与具体类型的方法映射。

Go 汇编中的关键符号
符号	含义
SB	静态基址（Static Base），全局符号的基址
SP	栈指针（Stack Pointer）
AX, BX, CX	通用寄存器（用于数据操作）
g	Goroutine 上下文指针



SB 静态地址
SP 指针
MOVEQ 	移动64位数据
CALL 	函数调用
TEXT	定义函数
SRODATA 用于只读数据   **数据段**
LEAQ 把一个内存地址加载到指定的寄存器当中，而非加载该地址所存储的值 加载有效地址 load effective add ress

MOV 用于把内存地址所存储的值加载到寄存器中




MOVEQ AX, addr 	将64位数据从AX寄存器移动到内存地址addr
TESTB AL, (AX)  检查AX指向的内存是否可读（防空指针）
CALL fun（SB）  调用函数
LEAQ sym, reg  	加载符号地址到寄存器
RET 函数返回





# command-line-arguments
main.(*Student).GrouUp STEXT nosplit size=19 args=0x8 locals=0x0 funcid=0x0 align=0x0; 
															指令    函数定义
	0x0000 00000 (/Users/jianfenliu/Workspace/x/main.go:16)	TEXT	main.(*Student).GrouUp(SB), NOSPLIT|NOFRAME|ABIInternal, $0-8 
	0x0000 00000 (/Users/jianfenliu/Workspace/x/main.go:16)	FUNCDATA	$0, gclocals·wgcWObbY2HYnK2SU/U22lA==(SB)
	0x0000 00000 (/Users/jianfenliu/Workspace/x/main.go:16)	FUNCDATA	$1, gclocals·J5F+7Qw7O7ve2QcWC7DpeQ==(SB)
	0x0000 00000 (/Users/jianfenliu/Workspace/x/main.go:16)	FUNCDATA	$5, main.(*Student).GrouUp.arginfo1(SB)
	0x0000 00000 (/Users/jianfenliu/Workspace/x/main.go:16)	MOVQ	AX, main.s+8(SP) ;将接收指针(*Student)存入栈 
	0x0005 00005 (/Users/jianfenliu/Workspace/x/main.go:17)	TESTB	AL, (AX)         ;检查指针是否这空
	0x0007 00007 (/Users/jianfenliu/Workspace/x/main.go:17)	TESTB	AL, (AX)         ; 冗余检查
	0x0009 00009 (/Users/jianfenliu/Workspace/x/main.go:17)	MOVQ	(AX), CX         ; 从Student结构体加载字段到CX （将结构体的字段加载到寄存器CX）
	0x000c 00012 (/Users/jianfenliu/Workspace/x/main.go:17)	INCQ	CX               ; CX中的值加1
	0x000f 00015 (/Users/jianfenliu/Workspace/x/main.go:17)	MOVQ	CX, (AX)         ; 将结果写加结构体字段
	0x0012 00018 (/Users/jianfenliu/Workspace/x/main.go:19)	RET                      ; 返回
	0x0000 48 89 44 24 08 84 00 84 00 48 8b 08 48 ff c1 48  H.D$.....H..H..H
	0x0010 89 08 c3                                         ...
main.main STEXT size=266 args=0x0 locals=0x88 funcid=0x0 align=0x0
	0x0000 00000 (/Users/jianfenliu/Workspace/x/main.go:21)	TEXT	main.main(SB), ABIInternal, $136-0
	0x0000 00000 (/Users/jianfenliu/Workspace/x/main.go:21)	LEAQ	-8(SP), R12
	0x0005 00005 (/Users/jianfenliu/Workspace/x/main.go:21)	CMPQ	R12, 16(R14)
	0x0009 00009 (/Users/jianfenliu/Workspace/x/main.go:21)	PCDATA	$0, $-2
	0x0009 00009 (/Users/jianfenliu/Workspace/x/main.go:21)	JLS	251
	0x000f 00015 (/Users/jianfenliu/Workspace/x/main.go:21)	PCDATA	$0, $-1
	0x000f 00015 (/Users/jianfenliu/Workspace/x/main.go:21)	PUSHQ	BP
	0x0010 00016 (/Users/jianfenliu/Workspace/x/main.go:21)	MOVQ	SP, BP
	0x0013 00019 (/Users/jianfenliu/Workspace/x/main.go:21)	ADDQ	$-128, SP
	0x0017 00023 (/Users/jianfenliu/Workspace/x/main.go:21)	FUNCDATA	$0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
	0x0017 00023 (/Users/jianfenliu/Workspace/x/main.go:21)	FUNCDATA	$1, gclocals·GsKDNj2n4ZlUlJ9f5t57bw==(SB)
	0x0017 00023 (/Users/jianfenliu/Workspace/x/main.go:21)	FUNCDATA	$2, main.main.stkobj(SB)
	0x0017 00023 (/Users/jianfenliu/Workspace/x/main.go:25)	LEAQ	type:main.Student(SB), AX
	0x001e 00030 (/Users/jianfenliu/Workspace/x/main.go:25)	PCDATA	$1, $0
	0x001e 00030 (/Users/jianfenliu/Workspace/x/main.go:25)	NOP
	0x0020 00032 (/Users/jianfenliu/Workspace/x/main.go:25)	CALL	runtime.newobject(SB) ; 在堆上分配 Student实例
	0x0025 00037 (/Users/jianfenliu/Workspace/x/main.go:25)	MOVQ	AX, main..autotmp_3+96(SP) ; 临时变量存储指针
	0x002a 00042 (/Users/jianfenliu/Workspace/x/main.go:25)	MOVQ	$18, (AX) ; 初始化Student.age = 18
	0x0031 00049 (/Users/jianfenliu/Workspace/x/main.go:25)	MOVQ	main..autotmp_3+96(SP), CX
	0x0036 00054 (/Users/jianfenliu/Workspace/x/main.go:25)	MOVQ	CX, main..autotmp_1+120(SP)
																		接口表
	0x003b 00059 (/Users/jianfenliu/Workspace/x/main.go:25)	LEAQ	go:itab.*main.Student,main.Person(SB), DX  ; 接口动态配发， 将*Student 类型的Person接口的接口表itab 加载到DX寄存顺
	0x0042 00066 (/Users/jianfenliu/Workspace/x/main.go:25)	MOVQ	DX, main.qcrao+24(SP)
	0x0047 00071 (/Users/jianfenliu/Workspace/x/main.go:25)	MOVQ	CX, main.qcrao+32(SP)
	0x004c 00076 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVUPS	X15, main..autotmp_2+104(SP)
	0x0052 00082 (/Users/jianfenliu/Workspace/x/main.go:26)	LEAQ	main..autotmp_2+104(SP), CX
	0x0057 00087 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	CX, main..autotmp_5+64(SP)
	0x005c 00092 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	main.qcrao+24(SP), CX
	0x0061 00097 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	main.qcrao+32(SP), DX
	0x0066 00102 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	CX, main..autotmp_6+48(SP)
	0x006b 00107 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	DX, main..autotmp_6+56(SP)
	0x0070 00112 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	CX, main..autotmp_7+40(SP)
	0x0075 00117 (/Users/jianfenliu/Workspace/x/main.go:26)	CMPQ	main..autotmp_7+40(SP), $0
	0x007b 00123 (/Users/jianfenliu/Workspace/x/main.go:26)	JNE	127
	0x007d 00125 (/Users/jianfenliu/Workspace/x/main.go:26)	JMP	143
	0x007f 00127 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	main..autotmp_7+40(SP), DX
	0x0084 00132 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	8(DX), DX
	0x0088 00136 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	DX, main..autotmp_7+40(SP)
	0x008d 00141 (/Users/jianfenliu/Workspace/x/main.go:26)	JMP	145
	0x008f 00143 (/Users/jianfenliu/Workspace/x/main.go:26)	JMP	145
	0x0091 00145 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	main..autotmp_5+64(SP), DX
	0x0096 00150 (/Users/jianfenliu/Workspace/x/main.go:26)	TESTB	AL, (DX)
	0x0098 00152 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	main..autotmp_7+40(SP), SI
	0x009d 00157 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	main..autotmp_6+56(SP), DI
	0x00a2 00162 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	SI, (DX)
	0x00a5 00165 (/Users/jianfenliu/Workspace/x/main.go:26)	CMPL	runtime.writeBarrier(SB), $0
	0x00ac 00172 (/Users/jianfenliu/Workspace/x/main.go:26)	PCDATA	$0, $-2
	0x00ac 00172 (/Users/jianfenliu/Workspace/x/main.go:26)	JEQ	176
	0x00ae 00174 (/Users/jianfenliu/Workspace/x/main.go:26)	JMP	178
	0x00b0 00176 (/Users/jianfenliu/Workspace/x/main.go:26)	JMP	196
	0x00b2 00178 (/Users/jianfenliu/Workspace/x/main.go:26)	CALL	runtime.gcWriteBarrier2(SB)
	0x00b7 00183 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	DI, (R11)
	0x00ba 00186 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	8(DX), SI
	0x00be 00190 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	SI, 8(R11)
	0x00c2 00194 (/Users/jianfenliu/Workspace/x/main.go:26)	JMP	196
	0x00c4 00196 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	DI, 8(DX)
	0x00c8 00200 (/Users/jianfenliu/Workspace/x/main.go:26)	PCDATA	$0, $-1
	0x00c8 00200 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	main..autotmp_5+64(SP), AX
	0x00cd 00205 (/Users/jianfenliu/Workspace/x/main.go:26)	TESTB	AL, (AX)
	0x00cf 00207 (/Users/jianfenliu/Workspace/x/main.go:26)	JMP	209
	0x00d1 00209 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	AX, main..autotmp_4+72(SP)
	0x00d6 00214 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	$1, main..autotmp_4+80(SP)
	0x00df 00223 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	$1, main..autotmp_4+88(SP)
	0x00e8 00232 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVL	$1, BX
	0x00ed 00237 (/Users/jianfenliu/Workspace/x/main.go:26)	MOVQ	BX, CX
	0x00f0 00240 (/Users/jianfenliu/Workspace/x/main.go:26)	CALL	fmt.Println(SB) ; 调用打印
	0x00f5 00245 (/Users/jianfenliu/Workspace/x/main.go:27)	SUBQ	$-128, SP
	0x00f9 00249 (/Users/jianfenliu/Workspace/x/main.go:27)	POPQ	BP
	0x00fa 00250 (/Users/jianfenliu/Workspace/x/main.go:27)	RET
	0x00fb 00251 (/Users/jianfenliu/Workspace/x/main.go:27)	NOP
	0x00fb 00251 (/Users/jianfenliu/Workspace/x/main.go:21)	PCDATA	$1, $-1
	0x00fb 00251 (/Users/jianfenliu/Workspace/x/main.go:21)	PCDATA	$0, $-2
	0x00fb 00251 (/Users/jianfenliu/Workspace/x/main.go:21)	NOP
	0x0100 00256 (/Users/jianfenliu/Workspace/x/main.go:21)	CALL	runtime.morestack_noctxt(SB)
	0x0105 00261 (/Users/jianfenliu/Workspace/x/main.go:21)	PCDATA	$0, $-1
	0x0105 00261 (/Users/jianfenliu/Workspace/x/main.go:21)	JMP	0
	0x0000 4c 8d 64 24 f8 4d 3b 66 10 0f 86 ec 00 00 00 55  L.d$.M;f.......U
	0x0010 48 89 e5 48 83 c4 80 48 8d 05 00 00 00 00 66 90  H..H...H......f.
	0x0020 e8 00 00 00 00 48 89 44 24 60 48 c7 00 12 00 00  .....H.D$`H.....
	0x0030 00 48 8b 4c 24 60 48 89 4c 24 78 48 8d 15 00 00  .H.L$`H.L$xH....
	0x0040 00 00 48 89 54 24 18 48 89 4c 24 20 44 0f 11 7c  ..H.T$.H.L$ D..|
	0x0050 24 68 48 8d 4c 24 68 48 89 4c 24 40 48 8b 4c 24  $hH.L$hH.L$@H.L$
	0x0060 18 48 8b 54 24 20 48 89 4c 24 30 48 89 54 24 38  .H.T$ H.L$0H.T$8
	0x0070 48 89 4c 24 28 48 83 7c 24 28 00 75 02 eb 10 48  H.L$(H.|$(.u...H
	0x0080 8b 54 24 28 48 8b 52 08 48 89 54 24 28 eb 02 eb  .T$(H.R.H.T$(...
	0x0090 00 48 8b 54 24 40 84 02 48 8b 74 24 28 48 8b 7c  .H.T$@..H.t$(H.|
	0x00a0 24 38 48 89 32 83 3d 00 00 00 00 00 74 02 eb 02  $8H.2.=.....t...
	0x00b0 eb 12 e8 00 00 00 00 49 89 3b 48 8b 72 08 49 89  .......I.;H.r.I.
	0x00c0 73 08 eb 00 48 89 7a 08 48 8b 44 24 40 84 00 eb  s...H.z.H.D$@...
	0x00d0 00 48 89 44 24 48 48 c7 44 24 50 01 00 00 00 48  .H.D$HH.D$P....H
	0x00e0 c7 44 24 58 01 00 00 00 bb 01 00 00 00 48 89 d9  .D$X.........H..
	0x00f0 e8 00 00 00 00 48 83 ec 80 5d c3 0f 1f 44 00 00  .....H...]...D..
	0x0100 e8 00 00 00 00 e9 f6 fe ff ff                    ..........
	rel 3+0 t=R_USEIFACE type:*main.Student+0
	rel 26+4 t=R_PCREL type:main.Student+0
	rel 33+4 t=R_CALL runtime.newobject+0
	rel 62+4 t=R_PCREL go:itab.*main.Student,main.Person+0
	rel 167+4 t=R_PCREL runtime.writeBarrier+-1
	rel 179+4 t=R_CALL runtime.gcWriteBarrier2+0
	rel 241+4 t=R_CALL fmt.Println+0 // 调用打印
	rel 257+4 t=R_CALL runtime.morestack_noctxt+0
main.Person.GrouUp STEXT dupok size=98 args=0x10 locals=0x10 funcid=0x16 align=0x0
	0x0000 00000 (<autogenerated>:1)	TEXT	main.Person.GrouUp(SB), DUPOK|WRAPPER|ABIInternal, $16-16
	0x0000 00000 (<autogenerated>:1)	CMPQ	SP, 16(R14)
	0x0004 00004 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0004 00004 (<autogenerated>:1)	JLS	50
	0x0006 00006 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0006 00006 (<autogenerated>:1)	PUSHQ	BP
	0x0007 00007 (<autogenerated>:1)	MOVQ	SP, BP
	0x000a 00010 (<autogenerated>:1)	SUBQ	$8, SP
	0x000e 00014 (<autogenerated>:1)	MOVQ	32(R14), R12
	0x0012 00018 (<autogenerated>:1)	TESTQ	R12, R12
	0x0015 00021 (<autogenerated>:1)	JNE	81
	0x0017 00023 (<autogenerated>:1)	NOP
	0x0017 00023 (<autogenerated>:1)	FUNCDATA	$0, gclocals·IuErl7MOXaHVn7EZYWzfFA==(SB)
	0x0017 00023 (<autogenerated>:1)	FUNCDATA	$1, gclocals·J5F+7Qw7O7ve2QcWC7DpeQ==(SB)
	0x0017 00023 (<autogenerated>:1)	FUNCDATA	$5, main.Person.GrouUp.arginfo1(SB)
	0x0017 00023 (<autogenerated>:1)	MOVQ	AX, main.~p0+24(SP)
	0x001c 00028 (<autogenerated>:1)	MOVQ	BX, main.~p0+32(SP)
	0x0021 00033 (<autogenerated>:1)	TESTB	AL, (AX)
	0x0023 00035 (<autogenerated>:1)	MOVQ	24(AX), CX
	0x0027 00039 (<autogenerated>:1)	MOVQ	BX, AX
	0x002a 00042 (<autogenerated>:1)	PCDATA	$1, $1
	0x002a 00042 (<autogenerated>:1)	CALL	CX
	0x002c 00044 (<autogenerated>:1)	ADDQ	$8, SP
	0x0030 00048 (<autogenerated>:1)	POPQ	BP
	0x0031 00049 (<autogenerated>:1)	RET
	0x0032 00050 (<autogenerated>:1)	NOP
	0x0032 00050 (<autogenerated>:1)	PCDATA	$1, $-1
	0x0032 00050 (<autogenerated>:1)	PCDATA	$0, $-2
	0x0032 00050 (<autogenerated>:1)	MOVQ	AX, 8(SP)
	0x0037 00055 (<autogenerated>:1)	MOVQ	BX, 16(SP)
	0x003c 00060 (<autogenerated>:1)	NOP
	0x0040 00064 (<autogenerated>:1)	CALL	runtime.morestack_noctxt(SB)
	0x0045 00069 (<autogenerated>:1)	PCDATA	$0, $-1
	0x0045 00069 (<autogenerated>:1)	MOVQ	8(SP), AX
	0x004a 00074 (<autogenerated>:1)	MOVQ	16(SP), BX
	0x004f 00079 (<autogenerated>:1)	JMP	0
	0x0051 00081 (<autogenerated>:1)	LEAQ	24(SP), R13
	0x0056 00086 (<autogenerated>:1)	CMPQ	(R12), R13
	0x005a 00090 (<autogenerated>:1)	JNE	23
	0x005c 00092 (<autogenerated>:1)	MOVQ	SP, (R12)
	0x0060 00096 (<autogenerated>:1)	JMP	23
	0x0000 49 3b 66 10 76 2c 55 48 89 e5 48 83 ec 08 4d 8b  I;f.v,UH..H...M.
	0x0010 66 20 4d 85 e4 75 3a 48 89 44 24 18 48 89 5c 24  f M..u:H.D$.H.\$
	0x0020 20 84 00 48 8b 48 18 48 89 d8 ff d1 48 83 c4 08   ..H.H.H....H...
	0x0030 5d c3 48 89 44 24 08 48 89 5c 24 10 0f 1f 40 00  ].H.D$.H.\$...@.
	0x0040 e8 00 00 00 00 48 8b 44 24 08 48 8b 5c 24 10 eb  .....H.D$.H.\$..
	0x0050 af 4c 8d 6c 24 18 4d 39 2c 24 75 bb 49 89 24 24  .L.l$.M9,$u.I.$$
	0x0060 eb b5                                            ..
	rel 2+0 t=R_USEIFACEMETHOD type:main.Person+96
	rel 42+0 t=R_CALLIND +0
	rel 65+4 t=R_CALL runtime.morestack_noctxt+0
go:cuinfo.producer.main SDWARFCUINFO dupok size=0
	0x0000 2d 4e 20 2d 6c 20 2d 73 68 61 72 65 64 20 72 65  -N -l -shared re
	0x0010 67 61 62 69                                      gabi
	                 只读数据段
runtime.interequal·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=R_ADDR runtime.interequal+0
runtime.memequal64·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=R_ADDR runtime.memequal64+0
runtime.gcbits.0100000000000000 SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
type:.namedata.*main.Person. SRODATA dupok size=14
	0x0000 01 0c 2a 6d 61 69 6e 2e 50 65 72 73 6f 6e        ..*main.Person
type:*main.Person SRODATA size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 f4 71 4a 8f 08 08 08 36 00 00 00 00 00 00 00 00  .qJ....6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Person.+0
	rel 48+8 t=R_ADDR type:main.Person+0
runtime.gcbits.0200000000000000 SRODATA dupok size=8
	0x0000 02 00 00 00 00 00 00 00                          ........
type:.namedata.*func()- SRODATA dupok size=9
	0x0000 00 07 2a 66 75 6e 63 28 29                       ..*func()
type:*func() SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 75 ac 29 27 08 08 08 36 00 00 00 00 00 00 00 00  u.)'...6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func()-+0
	rel 48+8 t=R_ADDR type:func()+0
type:func() SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 fe fa b9 80 02 08 08 33 00 00 00 00 00 00 00 00  .......3........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00                                      ....
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func()-+0
	rel 44+4 t=RelocType(-32763) type:*func()+0
type:.importpath.main. SRODATA dupok size=6
	0x0000 00 04 6d 61 69 6e                                ..main
type:.namedata.GrouUp. SRODATA dupok size=8
	0x0000 01 06 47 72 6f 75 55 70                          ..GrouUp
type:main.Person SRODATA size=104
	0x0000 10 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0010 32 2d bf 60 07 08 08 14 00 00 00 00 00 00 00 00  2-.`............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 18 00 00 00 00 00 00 00  ................
	0x0060 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.interequal·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0200000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Person.+0
	rel 44+4 t=R_ADDROFF type:*main.Person+0
	rel 48+8 t=R_ADDR type:.importpath.main.+0
	rel 56+8 t=R_ADDR type:main.Person+96
	rel 80+4 t=R_ADDROFF type:.importpath.main.+0
	rel 96+4 t=R_ADDROFF type:.namedata.GrouUp.+0
	rel 100+4 t=R_ADDROFF type:func()+0
type:.namedata.*main.Student. SRODATA dupok size=15
	0x0000 01 0d 2a 6d 61 69 6e 2e 53 74 75 64 65 6e 74     ..*main.Student
runtime.gcbits. SRODATA dupok size=0
type:.namedata.age- SRODATA dupok size=5
	0x0000 00 03 61 67 65                                   ..age
type:main.Student SRODATA size=120
	0x0000 08 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 eb 29 b9 7a 0f 08 08 19 00 00 00 00 00 00 00 00  .).z............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00 28 00 00 00 00 00 00 00  ........(.......
	0x0060 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0070 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Student.+0
	rel 44+4 t=R_ADDROFF type:*main.Student+0
	rel 48+8 t=R_ADDR type:.importpath.main.+0
	rel 56+8 t=R_ADDR type:main.Student+96
	rel 80+4 t=R_ADDROFF type:.importpath.main.+0
	rel 96+8 t=R_ADDR type:.namedata.age-+0
	rel 104+8 t=R_ADDR type:int+0
type:.namedata.*func(*main.Student)- SRODATA dupok size=22
	0x0000 00 14 2a 66 75 6e 63 28 2a 6d 61 69 6e 2e 53 74  ..*func(*main.St
	0x0010 75 64 65 6e 74 29                                udent)
type:*func(*main.Student) SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 40 78 c5 a4 08 08 08 36 00 00 00 00 00 00 00 00  @x.....6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.Student)-+0
	rel 48+8 t=R_ADDR type:func(*main.Student)+0
type:func(*main.Student) SRODATA dupok size=64
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 fd ac aa d3 02 08 08 33 00 00 00 00 00 00 00 00  .......3........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 01 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*func(*main.Student)-+0
	rel 44+4 t=RelocType(-32763) type:*func(*main.Student)+0
	rel 56+8 t=R_ADDR type:*main.Student+0
type:*main.Student SRODATA size=88
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 2d 16 dc d9 09 08 08 36 00 00 00 00 00 00 00 00  -......6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 01 00 01 00  ................
	0x0040 10 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0050 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*main.Student.+0
	rel 48+8 t=R_ADDR type:main.Student+0
	rel 56+4 t=R_ADDROFF type:.importpath.main.+0
	rel 72+4 t=R_ADDROFF type:.namedata.GrouUp.+0
	rel 76+4 t=R_METHODOFF type:func()+0
	rel 80+4 t=R_METHODOFF main.(*Student).GrouUp+0
	rel 84+4 t=R_METHODOFF main.(*Student).GrouUp+0
go:itab.*main.Student,main.Person SRODATA dupok size=32
	0x0000 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0010 2d 16 dc d9 00 00 00 00 00 00 00 00 00 00 00 00  -...............
	rel 0+8 t=R_ADDR type:main.Person+0
	rel 8+8 t=R_ADDR type:*main.Student+0
	rel 24+8 t=RelocType(-32767) main.(*Student).GrouUp+0
go:cuinfo.packagename.main SDWARFCUINFO dupok size=0
	0x0000 6d 61 69 6e                                      main
main..inittask SNOPTRDATA size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+0 t=R_INITORDER fmt..inittask+0
type:.namedata.*[1]interface {}- SRODATA dupok size=18
	0x0000 00 10 2a 5b 31 5d 69 6e 74 65 72 66 61 63 65 20  ..*[1]interface 
	0x0010 7b 7d                                            {}
runtime.nilinterequal·f SRODATA dupok size=8
	0x0000 00 00 00 00 00 00 00 00                          ........
	rel 0+8 t=R_ADDR runtime.nilinterequal+0
type:[1]interface {} SRODATA dupok size=72
	0x0000 10 00 00 00 00 00 00 00 10 00 00 00 00 00 00 00  ................
	0x0010 6e 20 6a 3d 02 08 08 11 00 00 00 00 00 00 00 00  n j=............
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0040 01 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.nilinterequal·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0200000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*[1]interface {}-+0
	rel 44+4 t=RelocType(-32763) type:*[1]interface {}+0
	rel 48+8 t=R_ADDR type:interface {}+0
	rel 56+8 t=R_ADDR type:[]interface {}+0
type:*[1]interface {} SRODATA dupok size=56
	0x0000 08 00 00 00 00 00 00 00 08 00 00 00 00 00 00 00  ................
	0x0010 a8 0e 57 36 08 08 08 36 00 00 00 00 00 00 00 00  ..W6...6........
	0x0020 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00  ................
	0x0030 00 00 00 00 00 00 00 00                          ........
	rel 24+8 t=R_ADDR runtime.memequal64·f+0
	rel 32+8 t=R_ADDR runtime.gcbits.0100000000000000+0
	rel 40+4 t=R_ADDROFF type:.namedata.*[1]interface {}-+0
	rel 48+8 t=R_ADDR type:[1]interface {}+0
gclocals·wgcWObbY2HYnK2SU/U22lA== SRODATA dupok size=10
	0x0000 02 00 00 00 01 00 00 00 01 00                    ..........
gclocals·J5F+7Qw7O7ve2QcWC7DpeQ== SRODATA dupok size=8
	0x0000 02 00 00 00 00 00 00 00                          ........
main.(*Student).GrouUp.arginfo1 SRODATA static dupok size=3
	0x0000 00 08 ff                                         ...
gclocals·g2BeySu+wFnoycgXfElmcg== SRODATA dupok size=8
	0x0000 01 00 00 00 00 00 00 00                          ........
gclocals·GsKDNj2n4ZlUlJ9f5t57bw== SRODATA dupok size=10
	0x0000 01 00 00 00 0d 00 00 00 00 00                    ..........
main.main.stkobj SRODATA static size=24
	0x0000 01 00 00 00 00 00 00 00 e8 ff ff ff 10 00 00 00  ................
	0x0010 10 00 00 00 00 00 00 00                          ........
	rel 20+4 t=R_ADDROFF runtime.gcbits.0200000000000000+0
gclocals·IuErl7MOXaHVn7EZYWzfFA== SRODATA dupok size=10
	0x0000 02 00 00 00 02 00 00 00 02 00                    ..........
main.Person.GrouUp.arginfo1 SRODATA static dupok size=7
	0x0000 fe 00 08 08 08 fd ff                             .......
