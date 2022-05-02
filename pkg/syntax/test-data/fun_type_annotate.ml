fun a (x:int) y = x + y
fun b (x:int) (y:int):t = x + y
fun c x y:t = x + y
val d = fn (x:int) y => x + y
val e = fn (x:int) (y:int):t => x + y
val f = fn x y:t => x + y
val g = fn _: (int ->t ->t) => f
