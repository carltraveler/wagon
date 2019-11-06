(module
  (type (;0;) (func (result i32)))
  (func (;0;) (type 0) (result i32)
    (local i32 i32)
    loop  ;; label = @1
      local.get 1 (;b;)
      local.get 0 (;a;)
      i32.add 
      local.set 1  (; b = b + a;) 
      local.get 0
      i32.const 1
      i32.add
      local.set 0 (; int a = a + 1;)
      local.get 0
      i32.const 5
      i32.lt_s   (; if a < 5 continnue;)
      if  ;; label = @2
        br 1 (;@1;)
      end
    end
	local.get 1)
  (export "loop" (func 0)))
