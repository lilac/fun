# Notes for contributors

Following is some notes that are intended for potential project contributors, to understand the architecture and design motives of the language and the compiler.

## Why go?

**Fun** is a _functional programming language_, but why is it compiled to Golang, and embracing the Golang ecosystem?

If **Fun** is a toy language, it could target at any runtimes, and there are lots of options. But **Fun** is ambitious to be a real industrial language, then multiple factors need to be considered. Popular options are listed below.

<table>
  <tbody><tr>
    <td> <i>Option</i> </td>
    <td>Pros</td>
    <td> <i>Cons</i> </td>
    <td> Example </td>
  </tr>
  <tr>
    <td>Mainstream virtual machine like JVM or .Net</td>
    <td>
    <ul>
      <li>Mature ecosystem</li>
      <li>Easy to implement</li>
    </ul>
    </td>
    <td> <i>Heavy dependency on the VM</i> </td>
    <td> Scala, F# </td>
  </tr>
  <tr>
    <td>Create the language's own runtime</td>
    <td>Tailored environment is flexible</td>
    <td><ul>
    <li><i>High bar on implementation cost</i></li>
    <li><i>Have to recreate the whole ecosystem</i></li>
    </ul></td>
    <td> Haskell, Ocaml </td>
  </tr>
  </tbody>
</table>

**Fun** wants to have the benefits of VM without writing a VM, so we are experimenting a new way: compiling to Golang. Note that transpiling to another language is not new, the GHC compiler also compiles to C.

Compared to C, Golang has all the tooling and features of a modern programming language, eg.

- Build system and module manager
- Garbage collector
- Green thread (go routine)

That's why we chose Golang as the target platform.

## Components

- [Abstract syntax tree](./pkg/ast/notes.md)
- [Type system](pkg/types/notes.md)
- [Lexing](./pkg/syntax/lexer.go)
- [Parsing](./pkg/syntax/parser.go)
- [Alpha transformation](./pkg/alpha/notes.md)
- [Type inference](./pkg/typing/notes.md)
- Code generation

## Roadmap

- Golang compatibility
  - Investigate [type reconstruction with structural subtyping](http://cristal.inria.fr/~simonet/publis/simonet-aplas03.pdf)
  - Retreat to local type inference
- Extend it with modern features
  - Default parameter of function
  - For comprehension
  - Extension method
  - Elegant exception/error handling

## Syntax

- Expression as Statement (do statement)
  - ~~Add a semicolon `;` to the end of an expression, to make it a statement~~.
  - Example
    ```ocaml
    do println "ab"
    (* equivalent to this statement. *)
    val _ = println "ab"
    ```
- Let-end expression
  - Make the `in` clause of `let-in-end` expression optional. When the `in` clause is skipped, the result of the expression is `()`.
  - Example
    ```ocaml
    let
      val a = 1
      do println a
      val b = "done"
      do println b
    end
    ```
- For expression
  - The `for` expression is like the [for comprehension of Elixir](https://elixir-lang.org/getting-started/comprehensions.html).
  - It supports 2 forms.
    - `for <pattern> in <exp> yield <exp>`
    - `for <pattern> in <exp> <block>`
  - Examples
    ```ocaml
    val lines = for x in [1, 2, 3] yield string x
    (* A complex comprehension *)
    val nums = for i in [1, 2, 3] do
      if i mod 2 = 0 then yield i
      else yield i * 2
    end
    ```
- Block expression
  - In procedural programming style, block (of statements) is the building blocks of program. Though many functional programming  languages do not support it, we think it's essential for practical programming.
  - A block contains a sequence of expressions, that are executed one by one.
  - Depending on the context, a block may contain control functions like `break`, `continue`, or `yield`.
  - Example
    ```ocaml
    fun log x y = do
      x := 1;
      println "x = $x";
      println "y = $y"
    end
    ```
- While expression
  - A procedural style loop.
  - Example
    ```ocaml
    let
      val a = 3
      in
      while i < 3 do
        println i
      end;
      a
    end
    ```
- For loop
  - The procedural version of for expression
  - Example
    ```ocaml
    for i in [1, 2, 3] do
      if i mod 2 = 0 then println i
    end
    ```
- Try expression
- Defer expression
- Raise/throw expression

## Thoughts

- Should we implement break/continue/return statements via escape continuation?
  - If so, maybe we can add a variant of `let in end` expression, and borrow the `let/ec` syntax from Racket.
  - In addition, if _escape continuation_ is supported, then all these statements plus exception handling can be built
    on it. Refer to [this doc](https://matt.might.net/articles/implementing-exceptions/).