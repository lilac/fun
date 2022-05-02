let
    val foo = fn x y => x + y
    fun print f a b = println_int (f a b)
in
    print(foo, 10, 42);
    print(fn x y => x - y, 10, 42)
end
