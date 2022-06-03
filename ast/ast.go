package ast

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/terakoya76/commentcov/proto"
)

// FileToCoverageItems is the logic of the plugin.
// it converts file to CoverageItems.
func FileToCoverageItems(logger hclog.Logger, file string) ([]*proto.CoverageItem, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		return []*proto.CoverageItem{}, err
	}

	items := ProcessFileCoverage(file, fset, f)
	return items, nil
}

// ProcessFileCoverage measures the comment coverage for the entire given file.
func ProcessFileCoverage(file string, fset *token.FileSet, f *ast.File) []*proto.CoverageItem {
	ci := ProcessPackageCoverage(file, fset, f)
	items := []*proto.CoverageItem{
		ci,
	}

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			ci = ProcessFunctionCoverage(file, fset, f, d)
			items = append(items, ci)

		case *ast.GenDecl:
			switch d.Tok {
			case token.IMPORT:
				// not check the coverage when token.IMPORT

			case token.CONST:
				cis := ProcessVariableCoverage(file, fset, f, d)
				items = append(items, cis...)

			case token.VAR:
				cis := ProcessVariableCoverage(file, fset, f, d)
				items = append(items, cis...)

			case token.TYPE:
				cis := ProcessTypeCoverage(file, fset, f, d)
				items = append(items, cis...)

			case token.ADD, token.ADD_ASSIGN, token.AND, token.AND_ASSIGN, token.AND_NOT, token.AND_NOT_ASSIGN, token.ARROW, token.ASSIGN, token.BREAK, token.CASE, token.CHAN, token.CHAR, token.COLON, token.COMMA, token.COMMENT, token.CONTINUE, token.DEC, token.DEFAULT, token.DEFER, token.DEFINE, token.ELLIPSIS, token.ELSE, token.EOF, token.EQL, token.FALLTHROUGH, token.FLOAT, token.FOR, token.FUNC, token.GEQ, token.GO, token.GOTO, token.GTR, token.IDENT, token.IF, token.ILLEGAL, token.IMAG, token.INC, token.INT, token.INTERFACE, token.LAND, token.LBRACE, token.LBRACK, token.LEQ, token.LOR, token.LPAREN, token.LSS, token.MAP, token.MUL, token.MUL_ASSIGN, token.NEQ, token.NOT, token.OR, token.OR_ASSIGN, token.PACKAGE, token.PERIOD, token.QUO, token.QUO_ASSIGN, token.RANGE, token.RBRACE, token.RBRACK, token.REM, token.REM_ASSIGN, token.RETURN, token.RPAREN, token.SELECT, token.SEMICOLON, token.SHL, token.SHL_ASSIGN, token.SHR, token.SHR_ASSIGN, token.STRING, token.STRUCT, token.SUB, token.SUB_ASSIGN, token.SWITCH, token.TILDE, token.XOR, token.XOR_ASSIGN: //nolint:lll
				// nothing to do
			}
		}
	}

	return items
}

// ProcessPackageCoverage measures the package level comment coverage.
func ProcessPackageCoverage(file string, fset *token.FileSet, f *ast.File) *proto.CoverageItem {
	sp := fset.Position(f.Package)
	block := &proto.Block{
		StartLine:   uint32(sp.Line),
		StartColumn: uint32(sp.Column),
		EndLine:     uint32(sp.Line),
		EndColumn:   uint32(sp.Column),
	}

	hcs := []*proto.Comment{}
	ics := []*proto.Comment{}

	for _, cg := range f.Comments {
		csp := fset.Position(cg.Pos())
		cep := fset.Position(cg.End())

		if IsHeader(fset, cg, block) {
			d := &proto.Comment{
				Comment: strings.TrimLeft(cg.Text(), " "),
				Block: &proto.Block{
					StartLine:   uint32(csp.Line),
					StartColumn: uint32(csp.Column),
					EndLine:     uint32(cep.Line),
					EndColumn:   uint32(cep.Column),
				},
			}
			hcs = append(hcs, d)
		}

		if IsInline(fset, cg, block) {
			d := &proto.Comment{
				Comment: strings.TrimLeft(cg.Text(), " "),
				Block: &proto.Block{
					StartLine:   uint32(csp.Line),
					StartColumn: uint32(csp.Column),
					EndLine:     uint32(cep.Line),
					EndColumn:   uint32(cep.Column),
				},
			}
			ics = append(ics, d)
		}
	}

	return &proto.CoverageItem{
		Scope:          proto.CoverageItem_FILE,
		TargetBlock:    block,
		File:           file,
		Identifier:     f.Name.Name,
		Extension:      filepath.Ext(file),
		HeaderComments: hcs,
		InlineComments: ics,
	}
}

// ProcessFunctionCoverage measures the comment coverage of functions.
func ProcessFunctionCoverage(file string, fset *token.FileSet, f *ast.File, fdecl *ast.FuncDecl) *proto.CoverageItem {
	sp := fset.Position(fdecl.Pos())
	ep := fset.Position(fdecl.End())
	block := &proto.Block{
		StartLine:   uint32(sp.Line),
		StartColumn: uint32(sp.Column),
		EndLine:     uint32(ep.Line),
		EndColumn:   uint32(ep.Column),
	}
	identifier := fdecl.Name.Name

	var scope proto.CoverageItem_Scope
	if ast.IsExported(identifier) {
		scope = proto.CoverageItem_PUBLIC_FUNCTION
	} else {
		scope = proto.CoverageItem_PRIVATE_FUNCTION
	}

	hcs := []*proto.Comment{}
	ics := []*proto.Comment{}
	for _, cg := range f.Comments {
		csp := fset.Position(cg.Pos())
		cep := fset.Position(cg.End())

		if IsHeader(fset, cg, block) {
			d := &proto.Comment{
				Comment: strings.TrimLeft(cg.Text(), " "),
				Block: &proto.Block{
					StartLine:   uint32(csp.Line),
					StartColumn: uint32(csp.Column),
					EndLine:     uint32(cep.Line),
					EndColumn:   uint32(cep.Column),
				},
			}
			hcs = append(hcs, d)
		}

		if IsInline(fset, cg, block) {
			d := &proto.Comment{
				Comment: strings.TrimLeft(cg.Text(), " "),
				Block: &proto.Block{
					StartLine:   uint32(csp.Line),
					StartColumn: uint32(csp.Column),
					EndLine:     uint32(cep.Line),
					EndColumn:   uint32(cep.Column),
				},
			}
			ics = append(ics, d)
		}
	}

	return &proto.CoverageItem{
		Scope:          scope,
		TargetBlock:    block,
		File:           file,
		Identifier:     identifier,
		Extension:      filepath.Ext(file),
		HeaderComments: hcs,
		InlineComments: ics,
	}
}

// ProcessVariableCoverage measures the comment coverage of variables.
func ProcessVariableCoverage(file string, fset *token.FileSet, f *ast.File, gdecl *ast.GenDecl) []*proto.CoverageItem {
	items := make([]*proto.CoverageItem, 0)

	for _, s := range gdecl.Specs {
		vs := s.(*ast.ValueSpec)

		for _, name := range vs.Names {
			identifier := name.Name
			sp := fset.Position(name.Pos())
			ep := fset.Position(name.End())
			block := &proto.Block{
				StartLine:   uint32(sp.Line),
				StartColumn: uint32(sp.Column),
				EndLine:     uint32(ep.Line),
				EndColumn:   uint32(ep.Column),
			}

			var scope proto.CoverageItem_Scope
			if ast.IsExported(identifier) {
				scope = proto.CoverageItem_PUBLIC_VARIABLE
			} else {
				scope = proto.CoverageItem_PRIVATE_VARIABLE
			}

			hcs := []*proto.Comment{}
			ics := []*proto.Comment{}
			for _, cg := range f.Comments {
				csp := fset.Position(cg.Pos())
				cep := fset.Position(cg.End())

				if IsHeader(fset, cg, block) {
					d := &proto.Comment{
						Comment: strings.TrimLeft(cg.Text(), " "),
						Block: &proto.Block{
							StartLine:   uint32(csp.Line),
							StartColumn: uint32(csp.Column),
							EndLine:     uint32(cep.Line),
							EndColumn:   uint32(cep.Column),
						},
					}
					hcs = append(hcs, d)
				}

				if IsInline(fset, cg, block) {
					d := &proto.Comment{
						Comment: strings.TrimLeft(cg.Text(), " "),
						Block: &proto.Block{
							StartLine:   uint32(csp.Line),
							StartColumn: uint32(csp.Column),
							EndLine:     uint32(cep.Line),
							EndColumn:   uint32(cep.Column),
						},
					}
					ics = append(ics, d)
				}
			}

			items = append(items, &proto.CoverageItem{
				Scope:          scope,
				TargetBlock:    block,
				File:           file,
				Identifier:     identifier,
				Extension:      filepath.Ext(file),
				HeaderComments: hcs,
				InlineComments: ics,
			})
		}
	}

	return items
}

// ProcessTypeCoverage measures the comment coverage of type declarations.
func ProcessTypeCoverage(file string, fset *token.FileSet, f *ast.File, gdecl *ast.GenDecl) []*proto.CoverageItem {
	items := make([]*proto.CoverageItem, 0)

	for _, s := range gdecl.Specs {
		ts := s.(*ast.TypeSpec)
		sp := fset.Position(ts.Pos())
		ep := fset.Position(ts.End())
		block := &proto.Block{
			StartLine:   uint32(sp.Line),
			StartColumn: uint32(sp.Column),
			EndLine:     uint32(ep.Line),
			EndColumn:   uint32(ep.Column),
		}
		Identifier := ts.Name.Name

		var scope proto.CoverageItem_Scope
		switch ts.Type.(type) {
		case *ast.ArrayType, *ast.ChanType, *ast.FuncType, *ast.MapType:
			if ast.IsExported(Identifier) {
				scope = proto.CoverageItem_PUBLIC_TYPE
			} else {
				scope = proto.CoverageItem_PRIVATE_TYPE
			}

		case *ast.InterfaceType, *ast.StructType:
			if ast.IsExported(Identifier) {
				scope = proto.CoverageItem_PUBLIC_CLASS
			} else {
				scope = proto.CoverageItem_PRIVATE_CLASS
			}
		}

		hcs := []*proto.Comment{}
		ics := []*proto.Comment{}
		for _, cg := range f.Comments {
			csp := fset.Position(cg.Pos())
			cep := fset.Position(cg.End())

			if IsHeader(fset, cg, block) {
				d := &proto.Comment{
					Comment: strings.TrimLeft(cg.Text(), " "),
					Block: &proto.Block{
						StartLine:   uint32(csp.Line),
						StartColumn: uint32(csp.Column),
						EndLine:     uint32(cep.Line),
						EndColumn:   uint32(cep.Column),
					},
				}
				hcs = append(hcs, d)
			}

			if IsInline(fset, cg, block) {
				d := &proto.Comment{
					Comment: strings.TrimLeft(cg.Text(), " "),
					Block: &proto.Block{
						StartLine:   uint32(csp.Line),
						StartColumn: uint32(csp.Column),
						EndLine:     uint32(cep.Line),
						EndColumn:   uint32(cep.Column),
					},
				}
				ics = append(ics, d)
			}
		}

		items = append(items, &proto.CoverageItem{
			Scope:          scope,
			TargetBlock:    block,
			File:           file,
			Identifier:     Identifier,
			Extension:      filepath.Ext(file),
			HeaderComments: hcs,
			InlineComments: ics,
		})
	}

	return items
}

// IsHeader returns true if the given *ast.CommentGroup is belonged to the given *proto.Block as HeaderComments.
func IsHeader(fset *token.FileSet, cg *ast.CommentGroup, b *proto.Block) bool {
	csp := fset.Position(cg.Pos())
	cep := fset.Position(cg.End())

	return (cep.Line == (int(b.StartLine)-1) && csp.Column <= int(b.StartColumn)) ||
		(cep.Line == int(b.StartLine) && cep.Column < int(b.StartColumn))
}

// IsInline returns true if the given *ast.CommentGroup is belonged to the given *proto.Block as InlineComments.
func IsInline(fset *token.FileSet, cg *ast.CommentGroup, b *proto.Block) bool {
	csp := fset.Position(cg.Pos())
	cep := fset.Position(cg.End())

	return int(b.StartLine) <= csp.Line && cep.Line <= int(b.EndLine)
}
