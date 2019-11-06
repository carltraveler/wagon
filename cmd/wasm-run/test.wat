(module
  (type (;0;) (func (param i32) (result i32)))
  (type (;1;) (func (result i32)))
  (func (;0;) (type 0) (param i32) (result i32)
    block (result i32)  ;; label = @1
      local.get 0
      i32.const 0
      i32.eq
      if (result i32);; label = @2
		i32.const 1
		i32.const 1
		i32.const 1
		i32.const 1
		i32.const 1
		i32.const 1
		br_if 0
	  else
		i32.const 5
      end
	end
	return)
  (func (;1;) (type 1) (result i32)
    i32.const 0
    call 0)
  (func (;2;) (type 1) (result i32)
    i32.const 1
    call 0)
  (export "invoke" (func 1))
  (export "test2" (func 2)))
