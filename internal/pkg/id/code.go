package id

// 很神奇的唯一性ID生成方式
func NewCode(id uint64, options ...func(*CodeOptions)) string {
	ops := getCodeOptionsOrSetDefault(nil)
	for _, f := range options {
		// 闭包函数保留了赋值逻辑, 如WithCodeChars([]rune(defaultABC))则已保留了字符数组的信息
		// 当将默认CodeOptions作为参数传入闭包函数时, CodeOptions的对应字段将被赋予上层闭包函数定义在参数中的字符数组
		f(ops)
	}
	// 将连续自增的id映射到更分散的数值空间, 一定程度破坏连续的线性关系
	// n1 与 字符数组长度互质, 根据贝祖定理, 保证了[(原id*n1) mod len(chars)]能遍历所有余数
	id = id*uint64(ops.n1) + ops.salt

	var code []rune
	slIdx := make([]byte, ops.l)

	charLen := len(ops.chars)
	charLenUI := uint64(charLen)

	// 扩散阶段
	// 每个slIdx位置的余数都被slIdx[0]影响 (非线性叠加)
	// 改变任意一个id都会影响多个位置的索引值 (雪崩效应)
	// 示例: id = 1时可能索引序列为[3,7,2,...], 而id = 2时可能则为[5,10,...].
	for i := 0; i < ops.l; i++ {
		slIdx[i] = byte(id % charLenUI)                          // get each number
		slIdx[i] = (slIdx[i] + byte(i)*slIdx[0]) % byte(charLen) // let units digit affect other digit
		id /= charLenUI                                          // right shift
	}

	// 混淆阶段(https://en.wikipedia.org/wiki/Permutation_box)
	// 由于n2与长度l互质, 实现code的完全置换
	// 如原始索引为[0,1,2,3,4]会被置换为[4,3,2,1,0]
	for i := 0; i < ops.l; i++ {
		idx := (byte(i) * byte(ops.n2)) % byte(ops.l)
		code = append(code, ops.chars[slIdx[idx]])
	}
	return string(code)
}
