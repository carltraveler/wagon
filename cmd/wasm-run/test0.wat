(module
  (type (;0;) (func (result i32)))
  (func (;0;) (type 0) (result i32)
    (local i32 i32)
	block
		i32.const 1
		if (result i32)
			i32.const 1
		else
			i32.const 2
			nop
		end
		drop
	end
	i32.const 9
   )
  (export "loop" (func 0)))
