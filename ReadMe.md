## Introduction
Fun is a minimal functional programming language that runs on the ecosystem of Golang. Its syntax is similar to that of Standard ML.

## Motto
**Programming should be fun**

## Goals
- Easy to learn
- Fun to program
- Minimal design, less is more
- Be practical <details><summary>why?</summary>A pure language is good for research, but an industrial language has to be practical to be widely used.</details>
  - Expressive
  - Programming in the large scale

## Features
- [x] Lexical analysis
- [x] Parsing
- [x] Semantics analysis
  - [x] Rename identifiers
  - [x] Semantics checks
- [ ] Type system
  - [x] Parametric polymorphism
  - [ ] Type inference
    - [x] Infix op expression
      - [ ] Type variables in arithmetic expressions
    - [x] Function application
    - [x] If-then-else
    - [x] Fn (lambda) expression
    - [x] Let-in expression
    - [x] Tuples
    - [x] Sequence
    - [ ] Patterns
      - [x] Constant pattern
      - [x] Var pattern
      - [ ] Tuple pattern
    - [ ] Type annotation
  - [ ] Record
  - [ ] Data type
  - [ ] List
  - [ ] Subtyping (structural subtyping)
- [ ] Code generation
  - [ ] Go ast
- [ ] Module
  - [ ] Import statement
  - [ ] Export annotation (or keyword)
- [ ] Package & Distribution
  - [ ] Package declaration
  - [ ] Import map
- [ ] Go packages interoperability
  - [ ] Go types mapping
  - [ ] Call Go functions
  - [ ] Call Go methods
