package codegen

import (
	fast "github.com/lilac/fun-lang/pkg/ast"
	"go/ast"
	"go/token"
)

func GenFunction(fun *fast.Fn) *ast.FuncLit {
	/*
		func (n int) int {
			if n > 2 {
				return fib(n-1) + fib(n-2)
			} else {
				return 1
			}
		}
	*/
	return &ast.FuncLit{
		Type: &ast.FuncType{
			Func: token.NoPos,
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{
							{
								Name: "n",
							},
						},
						Type: &ast.Ident{
							Name: "int",
						},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.Ident{
							Name: "int",
						},
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X: &ast.Ident{
							Name: "n",
						},
						Op: token.GTR,
						Y: &ast.BasicLit{
							Kind:  token.INT,
							Value: "2",
						},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.BinaryExpr{
										X: &ast.CallExpr{
											Fun: &ast.Ident{
												Name: "fib",
											},
											Args: []ast.Expr{
												&ast.BinaryExpr{
													X: &ast.Ident{
														Name: "n",
													},
													Op: token.SUB,
													Y: &ast.BasicLit{
														Kind:  token.INT,
														Value: "1",
													},
												},
											},
											Ellipsis: 0,
										},
										Op: token.ADD,
										Y: &ast.CallExpr{
											Fun: &ast.Ident{
												Name: "fib",
											},
											Args: []ast.Expr{
												&ast.BinaryExpr{
													X: &ast.Ident{
														Name: "n",
													},
													Op: token.SUB,
													Y: &ast.BasicLit{
														Kind:  token.INT,
														Value: "2",
													},
												},
											},
											Ellipsis: 0,
										},
									},
								},
							},
						},
					},
					Else: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ReturnStmt{
								Results: []ast.Expr{
									&ast.BasicLit{
										Kind:  token.INT,
										Value: "1",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
