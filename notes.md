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