package ast_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/commentcov/commentcov/proto"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	myAst "github.com/commentcov/commentcov-plugin-go/ast"
)

var (
	// blockCmp is go-cmp/cmp custom competitor for *proto.Block.
	blockCmp = cmp.Comparer(func(x, y *proto.Block) bool {
		return cmp.Equal(
			&x,
			&y,
			cmpopts.IgnoreUnexported(proto.Block{}),
		)
	})

	// commentCmp is go-cmp/cmp custom competitor for *proto.Comment.
	commentCmp = cmp.Comparer(func(x, y *proto.Comment) bool {
		return cmp.Equal(
			&x,
			&y,
			cmpopts.IgnoreUnexported(proto.Comment{}),
			blockCmp,
		)
	})

	// coverageItemCmp is go-cmp/cmp custom competitor for *proto.CoverageItem.
	coverageItemCmp = cmp.Comparer(func(x, y *proto.CoverageItem) bool {
		return cmp.Equal(
			&x,
			&y,
			cmpopts.IgnoreUnexported(proto.CoverageItem{}),
			commentCmp,
			blockCmp,
		)
	})
)

// TestProcessFileCoverage is the unittest for ProcessFileCoverage.
//
//nolint:funlen
func TestProcessFileCoverage(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		src      string
		want     []*proto.CoverageItem
	}{

		{
			name:     "file with comment",
			filename: "hoge.go",
			src: `// hoge Header
package hoge // hoge Inline
// Out of hoge

// MyVar Header
var MyVar string = "string" // MyVar Inline

// Out of MyVar
var ( // Out of MyVar
    // MyVar Header
    MyVar string = "string" // MyVar Inline
    // Out of MyVar
) // Out of MyVar

// MyConst Header
const MyConst string = "string" // MyConst Inline

// Out of MyConst
const ( // Out of MyConst
    // MyConst Header
    MyConst string = "string" // MyConst Inline
    // Out of MyConst
) // Out of MyConst

// MyStruct Header
type MyStruct struct { // MyStruct Inline
    // MyStruct Inline a
    a string
    // MyStruct Inline b
    b string
} // MyStruct Inline

// Out of MyStruct
type ( // Out of MyStruct
    // MyStruct Header
    MyStruct struct { // MyStruct Inline
        // MyStruct Inline a
        a string
        // MyStruct Inline b
        b string
    } // MyStruct Inline
    // Out of MyStruct
) // Out of MyStruct

// MyInterface Header
type MyInterface interface { // MyInterface Inline
    // MyInterface Inline a
    a() string
    // MyInterface Inline b
    b() string
} // MyInterface Inline

// Out of MyInterface
type ( // Out of MyInterface
    // MyInterface Header
    MyInterface interface { // MyInterface Inline
        // MyInterface Inline a
        a() string
        // MyInterface Inline b
        b() string
    } // MyInterface Inline
    // Out of MyInterface
) // Out of MyInterface

// MyType Header
type MyType = map[string]int // MyType Inline

// Out of MyType
type ( // Out of MyType
    // MyType Header
    MyType = map[string]int // MyType Inline
    // Out of MyType
) // Out of MyType

// MyFunc Header
func MyFunc() bool { // MyFunc Inline
    return true // MyFunc Inline return
} // MyFunc Inline

// Out of MyFunc
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_FILE,
					TargetBlock: &proto.Block{
						StartLine:   2,
						StartColumn: 1,
						EndLine:     2,
						EndColumn:   1,
					},
					File:       "hoge.go",
					Identifier: "hoge",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   1,
								StartColumn: 1,
								EndLine:     1,
								EndColumn:   15,
							},
							Comment: "hoge Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   2,
								StartColumn: 14,
								EndLine:     2,
								EndColumn:   28,
							},
							Comment: "hoge Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 5,
						EndLine:     6,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   16,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 29,
								EndLine:     6,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   11,
						StartColumn: 5,
						EndLine:     11,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 5,
								EndLine:     10,
								EndColumn:   20,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 29,
								EndLine:     11,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   16,
						StartColumn: 7,
						EndLine:     16,
						EndColumn:   14,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   15,
								StartColumn: 1,
								EndLine:     15,
								EndColumn:   18,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   16,
								StartColumn: 33,
								EndLine:     16,
								EndColumn:   50,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   21,
						StartColumn: 5,
						EndLine:     21,
						EndColumn:   12,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   20,
								StartColumn: 5,
								EndLine:     20,
								EndColumn:   22,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   21,
								StartColumn: 31,
								EndLine:     21,
								EndColumn:   48,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   26,
						StartColumn: 6,
						EndLine:     31,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   25,
								StartColumn: 1,
								EndLine:     25,
								EndColumn:   19,
							},
							Comment: "MyStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   26,
								StartColumn: 24,
								EndLine:     26,
								EndColumn:   42,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   27,
								StartColumn: 5,
								EndLine:     27,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   29,
								StartColumn: 5,
								EndLine:     29,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   31,
								StartColumn: 3,
								EndLine:     31,
								EndColumn:   21,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   36,
						StartColumn: 5,
						EndLine:     41,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   35,
								StartColumn: 5,
								EndLine:     35,
								EndColumn:   23,
							},
							Comment: "MyStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   36,
								StartColumn: 23,
								EndLine:     36,
								EndColumn:   41,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   37,
								StartColumn: 9,
								EndLine:     37,
								EndColumn:   29,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   39,
								StartColumn: 9,
								EndLine:     39,
								EndColumn:   29,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   41,
								StartColumn: 7,
								EndLine:     41,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   46,
						StartColumn: 6,
						EndLine:     51,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   45,
								StartColumn: 1,
								EndLine:     45,
								EndColumn:   22,
							},
							Comment: "MyInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   46,
								StartColumn: 30,
								EndLine:     46,
								EndColumn:   51,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   47,
								StartColumn: 5,
								EndLine:     47,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   49,
								StartColumn: 5,
								EndLine:     49,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   51,
								StartColumn: 3,
								EndLine:     51,
								EndColumn:   24,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   56,
						StartColumn: 5,
						EndLine:     61,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   55,
								StartColumn: 5,
								EndLine:     55,
								EndColumn:   26,
							},
							Comment: "MyInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   56,
								StartColumn: 29,
								EndLine:     56,
								EndColumn:   50,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   57,
								StartColumn: 9,
								EndLine:     57,
								EndColumn:   32,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   59,
								StartColumn: 9,
								EndLine:     59,
								EndColumn:   32,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   61,
								StartColumn: 7,
								EndLine:     61,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   66,
						StartColumn: 6,
						EndLine:     66,
						EndColumn:   29,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   65,
								StartColumn: 1,
								EndLine:     65,
								EndColumn:   17,
							},
							Comment: "MyType Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   66,
								StartColumn: 30,
								EndLine:     66,
								EndColumn:   46,
							},
							Comment: "MyType Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   71,
						StartColumn: 5,
						EndLine:     71,
						EndColumn:   28,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   70,
								StartColumn: 5,
								EndLine:     70,
								EndColumn:   21,
							},
							Comment: "MyType Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   71,
								StartColumn: 29,
								EndLine:     71,
								EndColumn:   45,
							},
							Comment: "MyType Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_FUNCTION,
					TargetBlock: &proto.Block{
						StartLine:   76,
						StartColumn: 1,
						EndLine:     78,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyFunc",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   75,
								StartColumn: 1,
								EndLine:     75,
								EndColumn:   17,
							},
							Comment: "MyFunc Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   76,
								StartColumn: 22,
								EndLine:     76,
								EndColumn:   38,
							},
							Comment: "MyFunc Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   77,
								StartColumn: 17,
								EndLine:     77,
								EndColumn:   40,
							},
							Comment: "MyFunc Inline return\n",
						},
						{
							Block: &proto.Block{
								StartLine:   78,
								StartColumn: 3,
								EndLine:     78,
								EndColumn:   19,
							},
							Comment: "MyFunc Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "file with comment /* */",
			filename: "hoge.go",
			src: `/* hoge Header */
/* hoge Header2 */ package hoge /* hoge Inline */
/* Out of hoge */

/* MyVar Header */
/* MyVar Header2 */ var MyVar string = "string" /* MyVar Inline */

/* Out of MyVar */
/* Out of MyVar */ var ( /* Out of MyVar */
    /* MyVar Header */
    /* MyVar Header2 */ MyVar string = "string" /* MyVar Inline */
    /* Out of MyVar */
) /* Out of MyVar */
/* MyConst Header */

/* MyConst Header */
/* MyConst Header2 */ const MyConst string = "string" /* MyConst Inline */

/* Out of MyConst */
/* Out of MyConst */ const ( /* Out of MyConst */
    /* MyConst Header */
    /* MyConst Header2 */ MyConst string = "string" /* MyConst Inline */
    /* Out of MyConst */
) /* Out of MyConst */

/* MyStruct Header */
/* MyStruct Header2 */ type MyStruct struct { /* MyStruct Inline */
    /* MyStruct Inline a */
    a string
    /* MyStruct Inline b */
    b string
} /* MyStruct Inline */

/* Out of MyStruct */
/* Out of MyStruct2 */ type ( /* Out of MyStruct */
    /* MyStruct Header */
    /* MyStruct Header2 */ MyStruct struct { /* MyStruct Inline */
        /* MyStruct Inline a */
        a string
        /* MyStruct Inline b */
        b string
    } /* MyStruct Inline */
    /* Out of MyStruct */
) /* Out of MyStruct */

/* MyInterface Header */
/* MyInterface Header2 */ type MyInterface interface { /* MyInterface Inline */
    /* MyInterface Inline a */
    a() string
    /* MyInterface Inline b */
    b() string
} /* MyInterface Inline */

/* Out of MyInterface */
/* Out of MyInterface */ type ( /* Out of MyInterface */
    /* MyInterface Header */
    /* MyInterface Header2 */ MyInterface interface { /* MyInterface Inline */
        /* MyInterface Inline a */
        a() string
        /* MyInterface Inline b */
        b() string
    } /* MyInterface Inline */
    /* Out of MyInterface */
) /* Out of MyInterface */

/* MyType Header */
/* MyType Header2 */ type MyType = map[string]int /* MyType Inline */

/* Out of MyType */
/* Out of MyType */ type ( /* Out of MyType */
    /* MyType Header */
    /* MyType Header2 */ MyType = map[string]int /* MyType Inline */
    /* Out of MyType */
) /* Out of MyType */

/* MyFunc Header */
/* MyFunc Header2 */ func MyFunc() bool { /* MyFunc Inline */
    return true /* MyFunc Inline return */
} /* MyFunc Inline */

/* Out of MyFunc */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_FILE,
					TargetBlock: &proto.Block{
						StartLine:   2,
						StartColumn: 20,
						EndLine:     2,
						EndColumn:   20,
					},
					File:       "hoge.go",
					Identifier: "hoge",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   1,
								StartColumn: 1,
								EndLine:     2,
								EndColumn:   19,
							},
							Comment: "hoge Header\n hoge Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   2,
								StartColumn: 33,
								EndLine:     2,
								EndColumn:   50,
							},
							Comment: "hoge Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 25,
						EndLine:     6,
						EndColumn:   30,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 1,
								EndLine:     6,
								EndColumn:   20,
							},
							Comment: "MyVar Header\n MyVar Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 49,
								EndLine:     6,
								EndColumn:   67,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   11,
						StartColumn: 25,
						EndLine:     11,
						EndColumn:   30,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 5,
								EndLine:     11,
								EndColumn:   24,
							},
							Comment: "MyVar Header\n MyVar Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 49,
								EndLine:     11,
								EndColumn:   67,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   17,
						StartColumn: 29,
						EndLine:     17,
						EndColumn:   36,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   16,
								StartColumn: 1,
								EndLine:     17,
								EndColumn:   22,
							},
							Comment: "MyConst Header\n MyConst Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   17,
								StartColumn: 55,
								EndLine:     17,
								EndColumn:   75,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   22,
						StartColumn: 27,
						EndLine:     22,
						EndColumn:   34,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   21,
								StartColumn: 5,
								EndLine:     22,
								EndColumn:   26,
							},
							Comment: "MyConst Header\n MyConst Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   22,
								StartColumn: 53,
								EndLine:     22,
								EndColumn:   73,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   27,
						StartColumn: 29,
						EndLine:     32,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   26,
								StartColumn: 1,
								EndLine:     27,
								EndColumn:   23,
							},
							Comment: "MyStruct Header\n MyStruct Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   27,
								StartColumn: 47,
								EndLine:     27,
								EndColumn:   68,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   28,
								StartColumn: 5,
								EndLine:     28,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   30,
								StartColumn: 5,
								EndLine:     30,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   32,
								StartColumn: 3,
								EndLine:     32,
								EndColumn:   24,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   37,
						StartColumn: 28,
						EndLine:     42,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   36,
								StartColumn: 5,
								EndLine:     37,
								EndColumn:   27,
							},
							Comment: "MyStruct Header\n MyStruct Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   37,
								StartColumn: 46,
								EndLine:     37,
								EndColumn:   67,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   38,
								StartColumn: 9,
								EndLine:     38,
								EndColumn:   32,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   40,
								StartColumn: 9,
								EndLine:     40,
								EndColumn:   32,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   42,
								StartColumn: 7,
								EndLine:     42,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   47,
						StartColumn: 32,
						EndLine:     52,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   46,
								StartColumn: 1,
								EndLine:     47,
								EndColumn:   26,
							},
							Comment: "MyInterface Header\n MyInterface Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   47,
								StartColumn: 56,
								EndLine:     47,
								EndColumn:   80,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   48,
								StartColumn: 5,
								EndLine:     48,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   50,
								StartColumn: 5,
								EndLine:     50,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   52,
								StartColumn: 3,
								EndLine:     52,
								EndColumn:   27,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   57,
						StartColumn: 31,
						EndLine:     62,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   56,
								StartColumn: 5,
								EndLine:     57,
								EndColumn:   30,
							},
							Comment: "MyInterface Header\n MyInterface Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   57,
								StartColumn: 55,
								EndLine:     57,
								EndColumn:   79,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   58,
								StartColumn: 9,
								EndLine:     58,
								EndColumn:   35,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   60,
								StartColumn: 9,
								EndLine:     60,
								EndColumn:   35,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   62,
								StartColumn: 7,
								EndLine:     62,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   67,
						StartColumn: 27,
						EndLine:     67,
						EndColumn:   50,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   66,
								StartColumn: 1,
								EndLine:     67,
								EndColumn:   21,
							},
							Comment: "MyType Header\n MyType Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   67,
								StartColumn: 51,
								EndLine:     67,
								EndColumn:   70,
							},
							Comment: "MyType Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   72,
						StartColumn: 26,
						EndLine:     72,
						EndColumn:   49,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   71,
								StartColumn: 5,
								EndLine:     72,
								EndColumn:   25,
							},
							Comment: "MyType Header\n MyType Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   72,
								StartColumn: 50,
								EndLine:     72,
								EndColumn:   69,
							},
							Comment: "MyType Inline\n",
						},
					},
				},

				{
					Scope: proto.CoverageItem_PUBLIC_FUNCTION,
					TargetBlock: &proto.Block{
						StartLine:   77,
						StartColumn: 22,
						EndLine:     79,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyFunc",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   76,
								StartColumn: 1,
								EndLine:     77,
								EndColumn:   21,
							},
							Comment: "MyFunc Header\n MyFunc Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   77,
								StartColumn: 43,
								EndLine:     77,
								EndColumn:   62,
							},
							Comment: "MyFunc Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   78,
								StartColumn: 17,
								EndLine:     78,
								EndColumn:   43,
							},
							Comment: "MyFunc Inline return\n",
						},
						{
							Block: &proto.Block{
								StartLine:   79,
								StartColumn: 3,
								EndLine:     79,
								EndColumn:   22,
							},
							Comment: "MyFunc Inline\n",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			got := myAst.ProcessFileCoverage(tt.filename, fset, f)
			if diff := cmp.Diff(tt.want, got, coverageItemCmp); diff != "" {
				t.Errorf("*proto.CoverageItem values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

// TestProcessPackageCoverage is the unittest for ProcessPackageCoverage.
//
//nolint:funlen
func TestProcessPackageCoverage(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		src      string
		want     *proto.CoverageItem
	}{

		{
			name:     "package with comment",
			filename: "hoge.go",
			src: `// hoge Header
package hoge // hoge Inline
// Out of hoge

// Out of hoge
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_FILE,
				TargetBlock: &proto.Block{
					StartLine:   2,
					StartColumn: 1,
					EndLine:     2,
					EndColumn:   1,
				},
				File:       "hoge.go",
				Identifier: "hoge",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   1,
							StartColumn: 1,
							EndLine:     1,
							EndColumn:   15,
						},
						Comment: "hoge Header\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   2,
							StartColumn: 14,
							EndLine:     2,
							EndColumn:   28,
						},
						Comment: "hoge Inline\n",
					},
				},
			},
		},

		{
			name:     "package with multi comments",
			filename: "hoge.go",
			src: `// hoge Header1
// hoge Header2
package hoge // hoge Inline
// Out of hoge

// Out of hoge
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_FILE,
				TargetBlock: &proto.Block{
					StartLine:   3,
					StartColumn: 1,
					EndLine:     3,
					EndColumn:   1,
				},
				File:       "hoge.go",
				Identifier: "hoge",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   1,
							StartColumn: 1,
							EndLine:     2,
							EndColumn:   16,
						},
						Comment: "hoge Header1\nhoge Header2\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   3,
							StartColumn: 14,
							EndLine:     3,
							EndColumn:   28,
						},
						Comment: "hoge Inline\n",
					},
				},
			},
		},

		{
			name:     "package with comment /* */",
			filename: "hoge.go",
			src: `/* hoge Header */
/* hoge Header2 */ package hoge /* hoge Inline */
/* Out of hoge */

/* Out of hoge */
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_FILE,
				TargetBlock: &proto.Block{
					StartLine:   2,
					StartColumn: 20,
					EndLine:     2,
					EndColumn:   20,
				},
				File:       "hoge.go",
				Identifier: "hoge",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   1,
							StartColumn: 1,
							EndLine:     2,
							EndColumn:   19,
						},
						Comment: "hoge Header\n hoge Header2\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   2,
							StartColumn: 33,
							EndLine:     2,
							EndColumn:   50,
						},
						Comment: "hoge Inline\n",
					},
				},
			},
		},

		{
			name:     "package with multi comments /* */",
			filename: "hoge.go",
			src: `/*
hoge Header
*/
package hoge /* hoge Inline */
/* Out of hoge */

/* Out of hoge */
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_FILE,
				TargetBlock: &proto.Block{
					StartLine:   4,
					StartColumn: 1,
					EndLine:     4,
					EndColumn:   1,
				},
				File:       "hoge.go",
				Identifier: "hoge",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   1,
							StartColumn: 1,
							EndLine:     3,
							EndColumn:   3,
						},
						Comment: "hoge Header\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   4,
							StartColumn: 14,
							EndLine:     4,
							EndColumn:   31,
						},
						Comment: "hoge Inline\n",
					},
				},
			},
		},

		{
			name:     "package with header comment",
			filename: "hoge.go",
			src: `// hoge Header
package hoge
// Out of hoge

// Out of hoge
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_FILE,
				TargetBlock: &proto.Block{
					StartLine:   2,
					StartColumn: 1,
					EndLine:     2,
					EndColumn:   1,
				},
				File:       "hoge.go",
				Identifier: "hoge",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   1,
							StartColumn: 1,
							EndLine:     1,
							EndColumn:   15,
						},
						Comment: "hoge Header\n",
					},
				},
				InlineComments: []*proto.Comment{},
			},
		},

		{
			name:     "package with inline comment",
			filename: "hoge.go",
			src: `package hoge // hoge Inline
// Out of hoge

// Out of hoge
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_FILE,
				TargetBlock: &proto.Block{
					StartLine:   1,
					StartColumn: 1,
					EndLine:     1,
					EndColumn:   1,
				},
				File:           "hoge.go",
				Identifier:     "hoge",
				Extension:      ".go",
				HeaderComments: []*proto.Comment{},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   1,
							StartColumn: 14,
							EndLine:     1,
							EndColumn:   28,
						},
						Comment: "hoge Inline\n",
					},
				},
			},
		},

		{
			name:     "package without comment",
			filename: "hoge.go",
			src: `// Out of hoge

package hoge

// Out of hoge
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_FILE,
				TargetBlock: &proto.Block{
					StartLine:   3,
					StartColumn: 1,
					EndLine:     3,
					EndColumn:   1,
				},
				File:           "hoge.go",
				Identifier:     "hoge",
				Extension:      ".go",
				HeaderComments: []*proto.Comment{},
				InlineComments: []*proto.Comment{},
			},
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		t.Run(tt.name, func(t *testing.T) {
			got := myAst.ProcessPackageCoverage(tt.filename, fset, f)
			if diff := cmp.Diff(tt.want, got, coverageItemCmp); diff != "" {
				t.Errorf("*proto.CoverageItem values are mismatch (-want +got):%s\n", diff)
			}
		})
	}
}

// TestProcessFunctionCoverage is the unittest for ProcessFunctionCoverage.
//
//nolint:funlen
func TestProcessFunctionCoverage(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		src      string
		want     *proto.CoverageItem
	}{

		{
			name:     "func with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyFunc

// MyFunc Header
func MyFunc() bool { // MyFunc Inline
    return true // MyFunc Inline return
} // MyFunc Inline

// Out of MyFunc
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PUBLIC_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   5,
					StartColumn: 1,
					EndLine:     7,
					EndColumn:   2,
				},
				File:       "hoge.go",
				Identifier: "MyFunc",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   4,
							StartColumn: 1,
							EndLine:     4,
							EndColumn:   17,
						},
						Comment: "MyFunc Header\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   5,
							StartColumn: 22,
							EndLine:     5,
							EndColumn:   38,
						},
						Comment: "MyFunc Inline\n",
					},
					{
						Block: &proto.Block{
							StartLine:   6,
							StartColumn: 17,
							EndLine:     6,
							EndColumn:   40,
						},
						Comment: "MyFunc Inline return\n",
					},
					{
						Block: &proto.Block{
							StartLine:   7,
							StartColumn: 3,
							EndLine:     7,
							EndColumn:   19,
						},
						Comment: "MyFunc Inline\n",
					},
				},
			},
		},

		{
			name:     "func with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyFunc

// MyFunc Header1
// MyFunc Header2
func MyFunc() bool { // MyFunc Inline
    return true // MyFunc Inline return
} // MyFunc Inline

// Out of MyFunc
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PUBLIC_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   6,
					StartColumn: 1,
					EndLine:     8,
					EndColumn:   2,
				},
				File:       "hoge.go",
				Identifier: "MyFunc",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   4,
							StartColumn: 1,
							EndLine:     5,
							EndColumn:   18,
						},
						Comment: "MyFunc Header1\nMyFunc Header2\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   6,
							StartColumn: 22,
							EndLine:     6,
							EndColumn:   38,
						},
						Comment: "MyFunc Inline\n",
					},
					{
						Block: &proto.Block{
							StartLine:   7,
							StartColumn: 17,
							EndLine:     7,
							EndColumn:   40,
						},
						Comment: "MyFunc Inline return\n",
					},
					{
						Block: &proto.Block{
							StartLine:   8,
							StartColumn: 3,
							EndLine:     8,
							EndColumn:   19,
						},
						Comment: "MyFunc Inline\n",
					},
				},
			},
		},

		{
			name:     "func with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyFunc */

/* MyFunc Header */
/* MyFunc Header2 */ func MyFunc() bool { /* MyFunc Inline */
    return true /* MyFunc Inline return */
} /* MyFunc Inline */

/* Out of MyFunc */
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PUBLIC_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   5,
					StartColumn: 22,
					EndLine:     7,
					EndColumn:   2,
				},
				File:       "hoge.go",
				Identifier: "MyFunc",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   4,
							StartColumn: 1,
							EndLine:     5,
							EndColumn:   21,
						},
						Comment: "MyFunc Header\n MyFunc Header2\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   5,
							StartColumn: 43,
							EndLine:     5,
							EndColumn:   62,
						},
						Comment: "MyFunc Inline\n",
					},
					{
						Block: &proto.Block{
							StartLine:   6,
							StartColumn: 17,
							EndLine:     6,
							EndColumn:   43,
						},
						Comment: "MyFunc Inline return\n",
					},
					{
						Block: &proto.Block{
							StartLine:   7,
							StartColumn: 3,
							EndLine:     7,
							EndColumn:   22,
						},
						Comment: "MyFunc Inline\n",
					},
				},
			},
		},

		{
			name:     "func with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyFunc */

/*
MyFunc Header
*/
func MyFunc() bool { /* MyFunc Inline */
    return true /* MyFunc Inline return */
} /* MyFunc Inline */

/* Out of MyFunc */
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PUBLIC_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   7,
					StartColumn: 1,
					EndLine:     9,
					EndColumn:   2,
				},
				File:       "hoge.go",
				Identifier: "MyFunc",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   4,
							StartColumn: 1,
							EndLine:     6,
							EndColumn:   3,
						},
						Comment: "MyFunc Header\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   7,
							StartColumn: 22,
							EndLine:     7,
							EndColumn:   41,
						},
						Comment: "MyFunc Inline\n",
					},
					{
						Block: &proto.Block{
							StartLine:   8,
							StartColumn: 17,
							EndLine:     8,
							EndColumn:   43,
						},
						Comment: "MyFunc Inline return\n",
					},
					{
						Block: &proto.Block{
							StartLine:   9,
							StartColumn: 3,
							EndLine:     9,
							EndColumn:   22,
						},
						Comment: "MyFunc Inline\n",
					},
				},
			},
		},

		{
			name:     "func with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyFunc

// MyFunc Header
func MyFunc() bool {
    return true
}

// Out of MyFunc
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PUBLIC_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   5,
					StartColumn: 1,
					EndLine:     7,
					EndColumn:   2,
				},
				File:       "hoge.go",
				Identifier: "MyFunc",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   4,
							StartColumn: 1,
							EndLine:     4,
							EndColumn:   17,
						},
						Comment: "MyFunc Header\n",
					},
				},
				InlineComments: []*proto.Comment{},
			},
		},

		{
			name:     "func with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyFunc

func MyFunc() bool { /* MyFunc Inline */
    return true /* MyFunc Inline return */
} /* MyFunc Inline */

// Out of MyFunc
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PUBLIC_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   4,
					StartColumn: 1,
					EndLine:     6,
					EndColumn:   2,
				},
				File:           "hoge.go",
				Identifier:     "MyFunc",
				Extension:      ".go",
				HeaderComments: []*proto.Comment{},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   4,
							StartColumn: 22,
							EndLine:     4,
							EndColumn:   41,
						},
						Comment: "MyFunc Inline\n",
					},
					{
						Block: &proto.Block{
							StartLine:   5,
							StartColumn: 17,
							EndLine:     5,
							EndColumn:   43,
						},
						Comment: "MyFunc Inline return\n",
					},
					{
						Block: &proto.Block{
							StartLine:   6,
							StartColumn: 3,
							EndLine:     6,
							EndColumn:   22,
						},
						Comment: "MyFunc Inline\n",
					},
				},
			},
		},

		{
			name:     "func private with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of myFunc

// myFunc Header
func myFunc() bool { // myFunc Inline
    return true // myFunc Inline return
} // myFunc Inline

// Out of myFunc
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PRIVATE_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   5,
					StartColumn: 1,
					EndLine:     7,
					EndColumn:   2,
				},
				File:       "hoge.go",
				Identifier: "myFunc",
				Extension:  ".go",
				HeaderComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   4,
							StartColumn: 1,
							EndLine:     4,
							EndColumn:   17,
						},
						Comment: "myFunc Header\n",
					},
				},
				InlineComments: []*proto.Comment{
					{
						Block: &proto.Block{
							StartLine:   5,
							StartColumn: 22,
							EndLine:     5,
							EndColumn:   38,
						},
						Comment: "myFunc Inline\n",
					},
					{
						Block: &proto.Block{
							StartLine:   6,
							StartColumn: 17,
							EndLine:     6,
							EndColumn:   40,
						},
						Comment: "myFunc Inline return\n",
					},
					{
						Block: &proto.Block{
							StartLine:   7,
							StartColumn: 3,
							EndLine:     7,
							EndColumn:   19,
						},
						Comment: "myFunc Inline\n",
					},
				},
			},
		},

		{
			name:     "func without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyFunc

func MyFunc() bool {
    return true
}

// Out of MyFunc
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PUBLIC_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   4,
					StartColumn: 1,
					EndLine:     6,
					EndColumn:   2,
				},
				File:           "hoge.go",
				Identifier:     "MyFunc",
				Extension:      ".go",
				HeaderComments: []*proto.Comment{},
				InlineComments: []*proto.Comment{},
			},
		},

		{
			name:     "func with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

// nolint:funlen
func MyFunc() bool {
    return true
}
`,
			want: &proto.CoverageItem{
				Scope: proto.CoverageItem_PUBLIC_FUNCTION,
				TargetBlock: &proto.Block{
					StartLine:   4,
					StartColumn: 1,
					EndLine:     6,
					EndColumn:   2,
				},
				File:           "hoge.go",
				Identifier:     "MyFunc",
				Extension:      ".go",
				HeaderComments: []*proto.Comment{},
				InlineComments: []*proto.Comment{},
			},
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, decl := range f.Decls {
			if d, ok := decl.(*ast.FuncDecl); ok {
				t.Run(tt.name, func(t *testing.T) {
					got := myAst.ProcessFunctionCoverage(tt.filename, fset, f, d)
					if diff := cmp.Diff(tt.want, got, coverageItemCmp); diff != "" {
						t.Errorf("*proto.CoverageItem values are mismatch (-want +got):%s\n", diff)
					}
				})
			} else {
				panic("not expected to be called")
			}
		}
	}
}

// TestProcessVariableCoverage_Var is the unittest for ProcessVariableCoverage: var.
//
//nolint:funlen
func TestProcessVariableCoverage_Var(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		src      string
		want     []*proto.CoverageItem
	}{

		{
			name:     "var with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// MyVar Header
var MyVar string = "string" // MyVar Inline

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   16,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 29,
								EndLine:     5,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// MyVar Header1
// MyVar Header2
var MyVar string = "string" // MyVar Inline

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 5,
						EndLine:     6,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   17,
							},
							Comment: "MyVar Header1\nMyVar Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 29,
								EndLine:     6,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var from func with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// MyVar Header
var MyVar = func() string { // MyVar Inline
    return "string" // MyVar Inline return
}() // MyVar Inline
// Out of MyVar

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   16,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 29,
								EndLine:     5,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyVar */

/* MyVar Header */
/* MyVar Header2 */ var MyVar string = "string" /* MyVar Inline */

/* Out of MyVar */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 25,
						EndLine:     5,
						EndColumn:   30,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   20,
							},
							Comment: "MyVar Header\n MyVar Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 49,
								EndLine:     5,
								EndColumn:   67,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyVar */

/*
MyVar Header
*/
var MyVar string = "string" /* MyVar Inline */

/* Out of MyVar */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     7,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     6,
								EndColumn:   3,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 29,
								EndLine:     7,
								EndColumn:   47,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// MyVar Header
var MyVar string = "string"

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   16,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "var with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

var MyVar string = "string" // MyVar Inline

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 5,
						EndLine:     4,
						EndColumn:   10,
					},
					File:           "hoge.go",
					Identifier:     "MyVar",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 29,
								EndLine:     4,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var private with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of myVar

// myVar Header
var myVar string = "string" // myVar Inline

// Out of myVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PRIVATE_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "myVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   16,
							},
							Comment: "myVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 29,
								EndLine:     5,
								EndColumn:   44,
							},
							Comment: "myVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

var MyVar string = "string"

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 5,
						EndLine:     4,
						EndColumn:   10,
					},
					File:           "hoge.go",
					Identifier:     "MyVar",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "var with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

// nolint:hoge
var MyVar string = "string"
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 5,
						EndLine:     4,
						EndColumn:   10,
					},
					File:           "hoge.go",
					Identifier:     "MyVar",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "var () with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// Out of MyVar
var ( // Out of MyVar
    // MyVar Header
    MyVar string = "string" // MyVar Inline
    // Out of MyVar
) // Out of MyVar

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     7,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   20,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 29,
								EndLine:     7,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var () with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// Out of MyVar
var ( // Out of MyVar
    // MyVar Header1
    // MyVar Header2
    MyVar string = "string" // MyVar Inline
    // Out of MyVar
) // Out of MyVar

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   8,
						StartColumn: 5,
						EndLine:     8,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   21,
							},
							Comment: "MyVar Header1\nMyVar Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 29,
								EndLine:     8,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var () from func with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// Out of MyVar
var ( // Out of MyVar
    // MyVar Header
    MyVar = func() string { // MyVar Inline
        return "string" // MyVar Inline return
	}() // MyVar Inline
    // Out of MyVar
) // Out of MyVar

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     7,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   20,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 29,
								EndLine:     7,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var () with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyVar */

/* Out of MyVar */
/* Out of MyVar */ var ( /* Out of MyVar */
    /* MyVar Header */
    /* MyVar Header2 */ MyVar string = "string" /* MyVar Inline */
    /* Out of MyVar */
) /* Out of MyVar */

/* Out of MyVar */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 25,
						EndLine:     7,
						EndColumn:   30,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   24,
							},
							Comment: "MyVar Header\n MyVar Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 49,
								EndLine:     7,
								EndColumn:   67,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var () with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyVar */

/* Out of MyVar */
var ( /* Out of MyVar */
    /*
    MyVar Header
    */
    MyVar string = "string" /* MyVar Inline */
    /* Out of MyVar */
) /* Out of MyVar */

/* Out of MyVar */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   9,
						StartColumn: 5,
						EndLine:     9,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   7,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 29,
								EndLine:     9,
								EndColumn:   47,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var () with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// Out of MyVar
var ( // Out of MyVar
    // MyVar Header
    MyVar string = "string"
) // Out of MyVar

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     7,
						EndColumn:   10,
					},
					File:       "hoge.go",
					Identifier: "MyVar",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   20,
							},
							Comment: "MyVar Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "var () with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

// Out of MyVar
var ( // Out of MyVar
    MyVar string = "string" // MyVar Inline
    // Out of MyVar
) // Out of MyVar

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 5,
						EndLine:     6,
						EndColumn:   10,
					},
					File:           "hoge.go",
					Identifier:     "MyVar",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 29,
								EndLine:     6,
								EndColumn:   44,
							},
							Comment: "MyVar Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "var () without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyVar

var (
    MyVar string = "string"
)

// Out of MyVar
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   10,
					},
					File:           "hoge.go",
					Identifier:     "MyVar",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "var () with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

var (
    // nolint:xxx
    MyVar string = "string"
)
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   10,
					},
					File:           "hoge.go",
					Identifier:     "MyVar",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, decl := range f.Decls {
			if d := decl.(*ast.GenDecl); d.Tok == token.VAR {
				t.Run(tt.name, func(t *testing.T) {
					got := myAst.ProcessVariableCoverage(tt.filename, fset, f, d)
					if diff := cmp.Diff(tt.want, got, coverageItemCmp); diff != "" {
						t.Errorf("[]*proto.CoverageItem values are mismatch (-want +got):%s\n", diff)
					}
				})
			} else {
				panic("not expected to be called")
			}
		}
	}
}

// TestProcessVariableCoverage_Const is the unittest for ProcessVariableCoverage: const.
//
//nolint:funlen
func TestProcessVariableCoverage_Const(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		src      string
		want     []*proto.CoverageItem
	}{

		{
			name:     "const with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

// MyConst Header
const MyConst string = "string" // MyConst Inline

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 7,
						EndLine:     5,
						EndColumn:   14,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   18,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 33,
								EndLine:     5,
								EndColumn:   50,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

// MyConst Header1
// MyConst Header2
const MyConst string = "string" // MyConst Inline

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 7,
						EndLine:     6,
						EndColumn:   14,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   19,
							},
							Comment: "MyConst Header1\nMyConst Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 33,
								EndLine:     6,
								EndColumn:   50,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyConst */

/* MyConst Header */
/* MyConst Header2 */ const MyConst string = "string" /* MyConst Inline */

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 29,
						EndLine:     5,
						EndColumn:   36,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   22,
							},
							Comment: "MyConst Header\n MyConst Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 55,
								EndLine:     5,
								EndColumn:   75,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyConst */

/*
MyConst Header
*/
const MyConst string = "string" /* MyConst Inline */

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 7,
						EndLine:     7,
						EndColumn:   14,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     6,
								EndColumn:   3,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 33,
								EndLine:     7,
								EndColumn:   53,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

// MyConst Header
const MyConst string = "string"

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 7,
						EndLine:     5,
						EndColumn:   14,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   18,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "const with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

const MyConst string = "string" // MyConst Inline

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 7,
						EndLine:     4,
						EndColumn:   14,
					},
					File:           "hoge.go",
					Identifier:     "MyConst",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 33,
								EndLine:     4,
								EndColumn:   50,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const private with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of myConst

// myConst Header
const myConst string = "string" // myConst Inline

// Out of myConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PRIVATE_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 7,
						EndLine:     5,
						EndColumn:   14,
					},
					File:       "hoge.go",
					Identifier: "myConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   18,
							},
							Comment: "myConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 33,
								EndLine:     5,
								EndColumn:   50,
							},
							Comment: "myConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

const MyConst string = "string"

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 7,
						EndLine:     4,
						EndColumn:   14,
					},
					File:           "hoge.go",
					Identifier:     "MyConst",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "const with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

// nolint:xxx
const MyConst string = "string"
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 7,
						EndLine:     4,
						EndColumn:   14,
					},
					File:           "hoge.go",
					Identifier:     "MyConst",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "const () with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

// Out of MyConst
const ( // Out of MyConst
    // MyConst Header
    MyConst string = "string" // MyConst Inline
    // Out of MyConst
) // Out of MyConst

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     7,
						EndColumn:   12,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   22,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 31,
								EndLine:     7,
								EndColumn:   48,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const () with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

// Out of MyConst
const ( // Out of MyConst
    // MyConst Header1
    // MyConst Header2
    MyConst string = "string" // MyConst Inline
    // Out of MyConst
) // Out of MyConst

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   8,
						StartColumn: 5,
						EndLine:     8,
						EndColumn:   12,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   23,
							},
							Comment: "MyConst Header1\nMyConst Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 31,
								EndLine:     8,
								EndColumn:   48,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const () with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyConst */

/* Out of MyConst */
/* Out of MyConst */ const ( /* Out of MyConst */
    /* MyConst Header */
    /* MyConst Header2 */ MyConst string = "string" /* MyConst Inline */
    /* Out of MyConst */
) /* Out of MyConst */

/* Out of MyConst */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 27,
						EndLine:     7,
						EndColumn:   34,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   26,
							},
							Comment: "MyConst Header\n MyConst Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 53,
								EndLine:     7,
								EndColumn:   73,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const () with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyConst */

/* Out of MyConst */
/* Out of MyConst */ const ( /* Out of MyConst */
    /*
    MyConst Header
    */
    MyConst string = "string" /* MyConst Inline */
    /* Out of MyConst */
) /* Out of MyConst */

/* Out of MyConst */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   9,
						StartColumn: 5,
						EndLine:     9,
						EndColumn:   12,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   7,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 31,
								EndLine:     9,
								EndColumn:   51,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const () with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

// Out of MyConst
const ( // Out of MyConst
    // MyConst Header
    MyConst string = "string"
    // Out of MyConst
) // Out of MyConst

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     7,
						EndColumn:   12,
					},
					File:       "hoge.go",
					Identifier: "MyConst",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   22,
							},
							Comment: "MyConst Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "const () with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

const ( // Out of MyConst
    MyConst string = "string" // MyConst Inline
    // Out of MyConst
) // Out of MyConst

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   12,
					},
					File:           "hoge.go",
					Identifier:     "MyConst",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 31,
								EndLine:     5,
								EndColumn:   48,
							},
							Comment: "MyConst Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const () enum with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of Color

const ( // Out of Color
    // Red Header
    Red Color = iota // Red Inline
    // Blue Header
    Blue // Blue Inline
    // Yellow Header
    Yellow // Yellow Inline
    // Out of Color
) // Out of Color

// Out of Color
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 5,
						EndLine:     6,
						EndColumn:   8,
					},
					File:       "hoge.go",
					Identifier: "Red",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 5,
								EndLine:     5,
								EndColumn:   18,
							},
							Comment: "Red Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 22,
								EndLine:     6,
								EndColumn:   35,
							},
							Comment: "Red Inline\n",
						},
					},
				},
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   8,
						StartColumn: 5,
						EndLine:     8,
						EndColumn:   9,
					},
					File:       "hoge.go",
					Identifier: "Blue",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   19,
							},
							Comment: "Blue Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 10,
								EndLine:     8,
								EndColumn:   24,
							},
							Comment: "Blue Inline\n",
						},
					},
				},
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   10,
						StartColumn: 5,
						EndLine:     10,
						EndColumn:   11,
					},
					File:       "hoge.go",
					Identifier: "Yellow",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 5,
								EndLine:     9,
								EndColumn:   21,
							},
							Comment: "Yellow Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 12,
								EndLine:     10,
								EndColumn:   28,
							},
							Comment: "Yellow Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "const () without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyConst

const (
    MyConst string = "string"
)

// Out of MyConst
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   12,
					},
					File:           "hoge.go",
					Identifier:     "MyConst",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "const () with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

const (
    // nolint:xxx
    MyConst string = "string"
)
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_VARIABLE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   12,
					},
					File:           "hoge.go",
					Identifier:     "MyConst",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, decl := range f.Decls {
			if d := decl.(*ast.GenDecl); d.Tok == token.CONST {
				t.Run(tt.name, func(t *testing.T) {
					got := myAst.ProcessVariableCoverage(tt.filename, fset, f, d)
					if diff := cmp.Diff(tt.want, got, coverageItemCmp); diff != "" {
						t.Errorf("[]*proto.CoverageItem values are mismatch (-want +got):%s\n", diff)
					}
				})
			} else {
				panic("not expected to be called")
			}
		}
	}
}

// TestProcessTypeCoverage_Struct is the unittest for ProcessTypeCoverage: Struct.
//
//nolint:funlen
func TestProcessTypeCoverage_Struct(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		src      string
		want     []*proto.CoverageItem
	}{

		{
			name:     "type struct with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

// MyStruct Header
type MyStruct struct { // MyStruct Inline
    // MyStruct Inline a
    a string
    // MyStruct Inline b
    b string
} // MyStruct Inline

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     10,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   19,
							},
							Comment: "MyStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 24,
								EndLine:     5,
								EndColumn:   42,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 3,
								EndLine:     10,
								EndColumn:   21,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type struct with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

// MyStruct Header1
// MyStruct Header2
type MyStruct struct { // MyStruct Inline
    // MyStruct Inline a
    a string
    // MyStruct Inline b
    b string
} // MyStruct Inline

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 6,
						EndLine:     11,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   20,
							},
							Comment: "MyStruct Header1\nMyStruct Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 24,
								EndLine:     6,
								EndColumn:   42,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 5,
								EndLine:     9,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 3,
								EndLine:     11,
								EndColumn:   21,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type struct with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyStruct */

/* MyStruct Header */
/* MyStruct Header2 */ type MyStruct struct { /* MyStruct Inline */
    /* MyStruct Inline a */
    a string
    /* MyStruct Inline b */
    b string
} /* MyStruct Inline */

/* Out of MyStruct */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 29,
						EndLine:     10,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   23,
							},
							Comment: "MyStruct Header\n MyStruct Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 47,
								EndLine:     5,
								EndColumn:   68,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 3,
								EndLine:     10,
								EndColumn:   24,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type struct with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyStruct */

/*
MyStruct Header
*/
type MyStruct struct { /* MyStruct Inline */
    /* MyStruct Inline a */
    a string
    /* MyStruct Inline b */
    b string
} /* MyStruct Inline */

/* Out of MyStruct */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 6,
						EndLine:     12,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     6,
								EndColumn:   3,
							},
							Comment: "MyStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 24,
								EndLine:     7,
								EndColumn:   45,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 5,
								EndLine:     10,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   12,
								StartColumn: 3,
								EndLine:     12,
								EndColumn:   24,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type struct with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

// MyStruct Header
type MyStruct struct {
    a string
    b string
}

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     8,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   19,
							},
							Comment: "MyStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type struct with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

type MyStruct struct { // MyStruct Inline
    // MyStruct Inline a
    a string
    // MyStruct Inline b
    b string
} // MyStruct Inline

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     9,
						EndColumn:   2,
					},
					File:           "hoge.go",
					Identifier:     "MyStruct",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 24,
								EndLine:     4,
								EndColumn:   42,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 5,
								EndLine:     5,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 3,
								EndLine:     9,
								EndColumn:   21,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type private struct with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of myStruct

// myStruct Header
type myStruct struct { // myStruct Inline
    // myStruct Inline a
    a string
    // myStruct Inline b
    b string
} // myStruct Inline

// Out of myStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PRIVATE_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     10,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "myStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   19,
							},
							Comment: "myStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 24,
								EndLine:     5,
								EndColumn:   42,
							},
							Comment: "myStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   25,
							},
							Comment: "myStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   25,
							},
							Comment: "myStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 3,
								EndLine:     10,
								EndColumn:   21,
							},
							Comment: "myStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type struct without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

type MyStruct struct {
    a string
    b string
}

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     7,
						EndColumn:   2,
					},
					File:           "hoge.go",
					Identifier:     "MyStruct",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type struct with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

// nolint:xxx
type MyStruct struct {
    a string
    b string
}
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     7,
						EndColumn:   2,
					},
					File:           "hoge.go",
					Identifier:     "MyStruct",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () struct with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

// Out of MyStruct
type ( // Out of MyStruct
    // MyStruct Header
    MyStruct struct { // MyStruct Inline
        // MyStruct Inline a
        a string
        // MyStruct Inline b
        b string
    } // MyStruct Inline
    // Out of MyStruct
) // Out of MyStruct

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     12,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   23,
							},
							Comment: "MyStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 23,
								EndLine:     7,
								EndColumn:   41,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 9,
								EndLine:     8,
								EndColumn:   29,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 9,
								EndLine:     10,
								EndColumn:   29,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   12,
								StartColumn: 7,
								EndLine:     12,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () struct with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

// Out of MyStruct
type ( // Out of MyStruct
    // MyStruct Header1
    // MyStruct Header2
    MyStruct struct { // MyStruct Inline
        // MyStruct Inline a
        a string
        // MyStruct Inline b
        b string
    } // MyStruct Inline
    // Out of MyStruct
) // Out of MyStruct

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   8,
						StartColumn: 5,
						EndLine:     13,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   24,
							},
							Comment: "MyStruct Header1\nMyStruct Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 23,
								EndLine:     8,
								EndColumn:   41,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 9,
								EndLine:     9,
								EndColumn:   29,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 9,
								EndLine:     11,
								EndColumn:   29,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   13,
								StartColumn: 7,
								EndLine:     13,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () struct with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyStruct */

/* Out of MyStruct */
/* Out of MyStruct */ type ( /* Out of MyStruct */
    /* MyStruct Header */
    /* MyStruct Header2 */ MyStruct struct { /* MyStruct Inline */
        /* MyStruct Inline a */
        a string
        /* MyStruct Inline b */
        b string
    } /* MyStruct Inline */
    /* Out of MyStruct */
) /* Out of MyStruct */

/* Out of MyStruct */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 28,
						EndLine:     12,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   27,
							},
							Comment: "MyStruct Header\n MyStruct Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 46,
								EndLine:     7,
								EndColumn:   67,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 9,
								EndLine:     8,
								EndColumn:   32,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 9,
								EndLine:     10,
								EndColumn:   32,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   12,
								StartColumn: 7,
								EndLine:     12,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () struct with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyStruct */

/* Out of MyStruct */
/* Out of MyStruct */ type ( /* Out of MyStruct */
    /*
    MyStruct Header
    */
    MyStruct struct { /* MyStruct Inline */
        /* MyStruct Inline a */
        a string
        /* MyStruct Inline b */
        b string
    } /* MyStruct Inline */
    /* Out of MyStruct */
) /* Out of MyStruct */

/* Out of MyStruct */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   9,
						StartColumn: 5,
						EndLine:     14,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   7,
							},
							Comment: "MyStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 23,
								EndLine:     9,
								EndColumn:   44,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 9,
								EndLine:     10,
								EndColumn:   32,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   12,
								StartColumn: 9,
								EndLine:     12,
								EndColumn:   32,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   14,
								StartColumn: 7,
								EndLine:     14,
								EndColumn:   28,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () struct with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

// Out of MyStruct
type ( // Out of MyStruct
    // MyStruct Header
    MyStruct struct {
        a string
        b string
    }
    // Out of MyStruct
) // Out of MyStruct

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     10,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyStruct",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   23,
							},
							Comment: "MyStruct Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () struct with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

// Out of MyStruct
type ( // Out of MyStruct
    MyStruct struct { // MyStruct Inline
        // MyStruct Inline a
        a string
        // MyStruct Inline b
        b string
    } // MyStruct Inline
    // Out of MyStruct
) // Out of MyStruct

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 5,
						EndLine:     11,
						EndColumn:   6,
					},
					File:           "hoge.go",
					Identifier:     "MyStruct",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 23,
								EndLine:     6,
								EndColumn:   41,
							},
							Comment: "MyStruct Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 9,
								EndLine:     7,
								EndColumn:   29,
							},
							Comment: "MyStruct Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 9,
								EndLine:     9,
								EndColumn:   29,
							},
							Comment: "MyStruct Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 7,
								EndLine:     11,
								EndColumn:   25,
							},
							Comment: "MyStruct Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () struct without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyStruct

type (
    MyStruct struct {
        a string
        b string
    }
)

// Out of MyStruct
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     8,
						EndColumn:   6,
					},
					File:           "hoge.go",
					Identifier:     "MyStruct",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () struct with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

// nolint:xxx
type (
    MyStruct struct {
        a string
        b string
    }
)
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     8,
						EndColumn:   6,
					},
					File:           "hoge.go",
					Identifier:     "MyStruct",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, decl := range f.Decls {
			if d := decl.(*ast.GenDecl); d.Tok == token.TYPE {
				t.Run(tt.name, func(t *testing.T) {
					got := myAst.ProcessTypeCoverage(tt.filename, fset, f, d)
					if diff := cmp.Diff(tt.want, got, coverageItemCmp); diff != "" {
						t.Errorf("[]*proto.CoverageItem values are mismatch (-want +got):%s\n", diff)
					}
				})
			} else {
				panic("not expected to be called")
			}
		}
	}
}

// TestProcessTypeCoverage_Interface is the unittest for ProcessTypeCoverage: Interface.
//
//nolint:funlen
func TestProcessTypeCoverage_Interface(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		src      string
		want     []*proto.CoverageItem
	}{

		{
			name:     "type interface with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

// MyInterface Header
type MyInterface interface { // MyInterface Inline
    // MyInterface Inline a
    a() string
    // MyInterface Inline b
    b() string
} // MyInterface Inline

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     10,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   22,
							},
							Comment: "MyInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 30,
								EndLine:     5,
								EndColumn:   51,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 3,
								EndLine:     10,
								EndColumn:   24,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type interface with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

// MyInterface Header1
// MyInterface Header2
type MyInterface interface { // MyInterface Inline
    // MyInterface Inline a
    a() string
    // MyInterface Inline b
    b() string
} // MyInterface Inline

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 6,
						EndLine:     11,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   23,
							},
							Comment: "MyInterface Header1\nMyInterface Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 30,
								EndLine:     6,
								EndColumn:   51,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 5,
								EndLine:     9,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 3,
								EndLine:     11,
								EndColumn:   24,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type interface with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyInterface */

/* MyInterface Header */
/* MyInterface Header2 */ type MyInterface interface { /* MyInterface Inline */
    /* MyInterface Inline a */
    a() string
    /* MyInterface Inline b */
    b() string
} /* MyInterface Inline */

/* Out of MyInterface */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 32,
						EndLine:     10,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   26,
							},
							Comment: "MyInterface Header\n MyInterface Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 56,
								EndLine:     5,
								EndColumn:   80,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 3,
								EndLine:     10,
								EndColumn:   27,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type interface with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyInterface */

/*
MyInterface Header
*/
type MyInterface interface { /* MyInterface Inline */
    /* MyInterface Inline a */
    a() string
    /* MyInterface Inline b */
    b() string
} /* MyInterface Inline */

/* Out of MyInterface */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 6,
						EndLine:     12,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     6,
								EndColumn:   3,
							},
							Comment: "MyInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 30,
								EndLine:     7,
								EndColumn:   54,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 5,
								EndLine:     10,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   12,
								StartColumn: 3,
								EndLine:     12,
								EndColumn:   27,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type interface with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

// MyInterface Header
type MyInterface interface {
    a() string
    b() string
}

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     8,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   22,
							},
							Comment: "MyInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type interface with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

type MyInterface interface { // MyInterface Inline
    // MyInterface Inline a
    a() string
    // MyInterface Inline b
    b() string
} // MyInterface Inline

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     9,
						EndColumn:   2,
					},
					File:           "hoge.go",
					Identifier:     "MyInterface",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 30,
								EndLine:     4,
								EndColumn:   51,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 5,
								EndLine:     5,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 3,
								EndLine:     9,
								EndColumn:   24,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type private interface with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of myInterface

// myInterface Header
type myInterface interface { // myInterface Inline
    // myInterface Inline a
    a() string
    // myInterface Inline b
    b() string
} // myInterface Inline

// Out of myInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PRIVATE_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     10,
						EndColumn:   2,
					},
					File:       "hoge.go",
					Identifier: "myInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   22,
							},
							Comment: "myInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 30,
								EndLine:     5,
								EndColumn:   51,
							},
							Comment: "myInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   28,
							},
							Comment: "myInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   28,
							},
							Comment: "myInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 3,
								EndLine:     10,
								EndColumn:   24,
							},
							Comment: "myInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type interface without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

type MyInterface interface {
    a() string
    b() string
}

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     7,
						EndColumn:   2,
					},
					File:           "hoge.go",
					Identifier:     "MyInterface",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type interface with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

// nolint:xxx
type MyInterface interface {
    a() string
    b() string
}
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     7,
						EndColumn:   2,
					},
					File:           "hoge.go",
					Identifier:     "MyInterface",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () interface with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

// Out of MyInterface
type ( // Out of MyInterface
    // MyInterface Header
    MyInterface interface { // MyInterface Inline
        // MyInterface Inline a
        a() string
        // MyInterface Inline b
        b() string
    } // MyInterface Inline
    // Out of MyInterface
) // Out of MyInterface

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     12,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   26,
							},
							Comment: "MyInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 29,
								EndLine:     7,
								EndColumn:   50,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 9,
								EndLine:     8,
								EndColumn:   32,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 9,
								EndLine:     10,
								EndColumn:   32,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   12,
								StartColumn: 7,
								EndLine:     12,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () interface with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

// Out of MyInterface
type ( // Out of MyInterface
    // MyInterface Header1
    // MyInterface Header2
    MyInterface interface { // MyInterface Inline
        // MyInterface Inline a
        a() string
        // MyInterface Inline b
        b() string
    } // MyInterface Inline
    // Out of MyInterface
) // Out of MyInterface

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   8,
						StartColumn: 5,
						EndLine:     13,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   27,
							},
							Comment: "MyInterface Header1\nMyInterface Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 29,
								EndLine:     8,
								EndColumn:   50,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 9,
								EndLine:     9,
								EndColumn:   32,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 9,
								EndLine:     11,
								EndColumn:   32,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   13,
								StartColumn: 7,
								EndLine:     13,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () interface with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyInterface */

/* Out of MyInterface */
/* Out of MyInterface */ type ( /* Out of MyInterface */
    /* MyInterface Header */
    /* MyInterface Header2 */ MyInterface interface { /* MyInterface Inline */
        /* MyInterface Inline a */
        a() string
        /* MyInterface Inline b */
        b() string
    } /* MyInterface Inline */
    /* Out of MyInterface */
) /* Out of MyInterface */

/* Out of MyInterface */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 31,
						EndLine:     12,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   30,
							},
							Comment: "MyInterface Header\n MyInterface Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 55,
								EndLine:     7,
								EndColumn:   79,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 9,
								EndLine:     8,
								EndColumn:   35,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 9,
								EndLine:     10,
								EndColumn:   35,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   12,
								StartColumn: 7,
								EndLine:     12,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () interface with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyInterface */

/* Out of MyInterface */
/* Out of MyInterface */ type ( /* Out of MyInterface */
    /*
    MyInterface Header
    */
    MyInterface interface { /* MyInterface Inline */
        /* MyInterface Inline a */
        a() string
        /* MyInterface Inline b */
        b() string
    } /* MyInterface Inline */
    /* Out of MyInterface */
) /* Out of MyInterface */

/* Out of MyInterface */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   9,
						StartColumn: 5,
						EndLine:     14,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   7,
							},
							Comment: "MyInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 29,
								EndLine:     9,
								EndColumn:   53,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   10,
								StartColumn: 9,
								EndLine:     10,
								EndColumn:   35,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   12,
								StartColumn: 9,
								EndLine:     12,
								EndColumn:   35,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   14,
								StartColumn: 7,
								EndLine:     14,
								EndColumn:   31,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () interface with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

// Out of MyInterface
type ( // Out of MyInterface
    // MyInterface Header
    MyInterface interface {
        a() string
        b() string
    }
    // Out of MyInterface
) // Out of MyInterface

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     10,
						EndColumn:   6,
					},
					File:       "hoge.go",
					Identifier: "MyInterface",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   26,
							},
							Comment: "MyInterface Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () interface with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

// Out of MyInterface
type ( // Out of MyInterface
    MyInterface interface { // MyInterface Inline
        // MyInterface Inline a
        a() string
        // MyInterface Inline b
        b() string
    } // MyInterface Inline
    // Out of MyInterface
) // Out of MyInterface

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 5,
						EndLine:     11,
						EndColumn:   6,
					},
					File:           "hoge.go",
					Identifier:     "MyInterface",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 29,
								EndLine:     6,
								EndColumn:   50,
							},
							Comment: "MyInterface Inline\n",
						},
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 9,
								EndLine:     7,
								EndColumn:   32,
							},
							Comment: "MyInterface Inline a\n",
						},
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 9,
								EndLine:     9,
								EndColumn:   32,
							},
							Comment: "MyInterface Inline b\n",
						},
						{
							Block: &proto.Block{
								StartLine:   11,
								StartColumn: 7,
								EndLine:     11,
								EndColumn:   28,
							},
							Comment: "MyInterface Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () interface without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyInterface

type MyInterface interface {
    a() string
    b() string
}

// Out of MyInterface
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     7,
						EndColumn:   2,
					},
					File:           "hoge.go",
					Identifier:     "MyInterface",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () interface with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

// nolint:xxx
type MyInterface interface {
    a() string
    b() string
}
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_CLASS,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     7,
						EndColumn:   2,
					},
					File:           "hoge.go",
					Identifier:     "MyInterface",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, decl := range f.Decls {
			if d := decl.(*ast.GenDecl); d.Tok == token.TYPE {
				t.Run(tt.name, func(t *testing.T) {
					got := myAst.ProcessTypeCoverage(tt.filename, fset, f, d)
					if diff := cmp.Diff(tt.want, got, coverageItemCmp); diff != "" {
						t.Errorf("[]*proto.CoverageItem values are mismatch (-want +got):%s\n", diff)
					}
				})
			} else {
				panic("not expected to be called")
			}
		}
	}
}

// TestProcessTypeCoverage_TypeAlias is the unittest for ProcessTypeCoverage: Type Alias.
//
//nolint:funlen
func TestProcessTypeCoverage_TypeAlias(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		src      string
		want     []*proto.CoverageItem
	}{

		{
			name:     "type alias with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

// MyType Header
type MyType = map[string]int // MyType Inline

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     5,
						EndColumn:   29,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   17,
							},
							Comment: "MyType Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 30,
								EndLine:     5,
								EndColumn:   46,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type alias with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

// MyType Header1
// MyType Header2
type MyType = map[string]int // MyType Inline

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 6,
						EndLine:     6,
						EndColumn:   29,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   18,
							},
							Comment: "MyType Header1\nMyType Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 30,
								EndLine:     6,
								EndColumn:   46,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type alias with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyType */

/* MyType Header */
/* MyType Header2 */ type MyType = map[string]int /* MyType Inline */

/* Out of MyType */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 27,
						EndLine:     5,
						EndColumn:   50,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     5,
								EndColumn:   21,
							},
							Comment: "MyType Header\n MyType Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 51,
								EndLine:     5,
								EndColumn:   70,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type alias with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyType */

/*
MyType Header
*/
type MyType = map[string]int /* MyType Inline */

/* Out of MyType */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 6,
						EndLine:     7,
						EndColumn:   29,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     6,
								EndColumn:   3,
							},
							Comment: "MyType Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 30,
								EndLine:     7,
								EndColumn:   49,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type alias with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

// MyType Header
type MyType = map[string]int

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     5,
						EndColumn:   29,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   17,
							},
							Comment: "MyType Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type alias with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

type MyType = map[string]int // MyType Inline

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     4,
						EndColumn:   29,
					},
					File:           "hoge.go",
					Identifier:     "MyType",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 30,
								EndLine:     4,
								EndColumn:   46,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type private alias with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of myType

// myType Header
type myType = map[string]int // myType Inline

// Out of myType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PRIVATE_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 6,
						EndLine:     5,
						EndColumn:   29,
					},
					File:       "hoge.go",
					Identifier: "myType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   4,
								StartColumn: 1,
								EndLine:     4,
								EndColumn:   17,
							},
							Comment: "myType Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   5,
								StartColumn: 30,
								EndLine:     5,
								EndColumn:   46,
							},
							Comment: "myType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type alias without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

type MyType = map[string]int

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     4,
						EndColumn:   29,
					},
					File:           "hoge.go",
					Identifier:     "MyType",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type alias with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

// nolint:xxx
type MyType = map[string]int

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   4,
						StartColumn: 6,
						EndLine:     4,
						EndColumn:   29,
					},
					File:           "hoge.go",
					Identifier:     "MyType",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () alias with comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

// Out of MyType
type ( // Out of MyType
    // MyType Header
    MyType = map[string]int // MyType Inline
    // Out of MyType
) // Out of MyType

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     7,
						EndColumn:   28,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   21,
							},
							Comment: "MyType Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 29,
								EndLine:     7,
								EndColumn:   45,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () alias with multi comments",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

// Out of MyType
type ( // Out of MyType
    // MyType Header1
    // MyType Header2
    MyType = map[string]int // MyType Inline
    // Out of MyType
) // Out of MyType

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   8,
						StartColumn: 5,
						EndLine:     8,
						EndColumn:   28,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   22,
							},
							Comment: "MyType Header1\nMyType Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   8,
								StartColumn: 29,
								EndLine:     8,
								EndColumn:   45,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () alias with comment /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyType */

/* Out of MyType */
/* Out of MyType */ type ( /* Out of MyType */
    /* MyType Header */
    /* MyType Header2 */ MyType = map[string]int /* MyType Inline */
    /* Out of MyType */
) /* Out of MyType */

/* Out of MyType */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 26,
						EndLine:     7,
						EndColumn:   49,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     7,
								EndColumn:   25,
							},
							Comment: "MyType Header\n MyType Header2\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   7,
								StartColumn: 50,
								EndLine:     7,
								EndColumn:   69,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () alias with multi comments /* */",
			filename: "hoge.go",
			src: `package hoge
/* Out of MyType */

/* Out of MyType */
/* Out of MyType */ type ( /* Out of MyType */
    /*
    MyType Header
    */
    MyType = map[string]int /* MyType Inline */
    /* Out of MyType */
) /* Out of MyType */

/* Out of MyType */
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   9,
						StartColumn: 5,
						EndLine:     9,
						EndColumn:   28,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     8,
								EndColumn:   7,
							},
							Comment: "MyType Header\n",
						},
					},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   9,
								StartColumn: 29,
								EndLine:     9,
								EndColumn:   48,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () alias with header comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

// Out of MyType
type ( // Out of MyType
    // MyType Header
    MyType = map[string]int
    // Out of MyType
) // Out of MyType

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   7,
						StartColumn: 5,
						EndLine:     7,
						EndColumn:   28,
					},
					File:       "hoge.go",
					Identifier: "MyType",
					Extension:  ".go",
					HeaderComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 5,
								EndLine:     6,
								EndColumn:   21,
							},
							Comment: "MyType Header\n",
						},
					},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () alias with inline comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

// Out of MyType
type ( // Out of MyType
    MyType = map[string]int // MyType Inline
    // Out of MyType
) // Out of MyType

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   6,
						StartColumn: 5,
						EndLine:     6,
						EndColumn:   28,
					},
					File:           "hoge.go",
					Identifier:     "MyType",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{
						{
							Block: &proto.Block{
								StartLine:   6,
								StartColumn: 29,
								EndLine:     6,
								EndColumn:   45,
							},
							Comment: "MyType Inline\n",
						},
					},
				},
			},
		},

		{
			name:     "type () alias without comment",
			filename: "hoge.go",
			src: `package hoge
// Out of MyType

type (
    MyType = map[string]int
)

// Out of MyType
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   28,
					},
					File:           "hoge.go",
					Identifier:     "MyType",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},

		{
			name:     "type () alias with nolint annotation",
			filename: "hoge.go",
			src: `package hoge

type (
    // nolint:xxx
    MyType = map[string]int
)
`,
			want: []*proto.CoverageItem{
				{
					Scope: proto.CoverageItem_PUBLIC_TYPE,
					TargetBlock: &proto.Block{
						StartLine:   5,
						StartColumn: 5,
						EndLine:     5,
						EndColumn:   28,
					},
					File:           "hoge.go",
					Identifier:     "MyType",
					Extension:      ".go",
					HeaderComments: []*proto.Comment{},
					InlineComments: []*proto.Comment{},
				},
			},
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, decl := range f.Decls {
			if d := decl.(*ast.GenDecl); d.Tok == token.TYPE {
				t.Run(tt.name, func(t *testing.T) {
					got := myAst.ProcessTypeCoverage(tt.filename, fset, f, d)
					if diff := cmp.Diff(tt.want, got, coverageItemCmp); diff != "" {
						t.Errorf("[]*proto.CoverageItem values are mismatch (-want +got):%s\n", diff)
					}
				})
			} else {
				panic("not expected to be called")
			}
		}
	}
}

// TestIsHeader is the unittest for IsHeader.
func TestIsHeader(t *testing.T) {
	tests := []struct {
		name  string
		src   string
		block *proto.Block
		want  bool
	}{

		{
			name: "header comment",
			src: `package hoge
// This Is Header
type MyType = map[string]int
`,
			block: &proto.Block{
				StartLine:   3,
				StartColumn: 6,
				EndLine:     3,
				EndColumn:   29,
			},
			want: true,
		},

		{
			name: "inline comment",
			src: `package hoge

type MyType = map[string]int // This Is Inline
`,
			block: &proto.Block{
				StartLine:   3,
				StartColumn: 6,
				EndLine:     3,
				EndColumn:   29,
			},
			want: false,
		},

		{
			name: "out of comment",
			src: `package hoge
// Out of MyType

type MyType = map[string]int

// Out of MyType
`,
			block: &proto.Block{
				StartLine:   4,
				StartColumn: 6,
				EndLine:     4,
				EndColumn:   29,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, cg := range f.Comments {
			t.Run(tt.name, func(t *testing.T) {
				got := myAst.IsHeader(fset, cg, tt.block)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("bool values are mismatch (-want +got):%s\n", diff)
				}
			})
		}
	}
}

// TestIsInline is the unittest for IsInline.
func TestIsInline(t *testing.T) {
	tests := []struct {
		name  string
		src   string
		block *proto.Block
		want  bool
	}{

		{
			name: "inline comment",
			src: `package hoge

type MyType = map[string]int // This Is Inline
`,
			block: &proto.Block{
				StartLine:   3,
				StartColumn: 6,
				EndLine:     3,
				EndColumn:   29,
			},
			want: true,
		},

		{
			name: "header comment",
			src: `package hoge
// This Is Header
type MyType = map[string]int
`,
			block: &proto.Block{
				StartLine:   3,
				StartColumn: 6,
				EndLine:     3,
				EndColumn:   29,
			},
			want: false,
		},

		{
			name: "out of comment",
			src: `package hoge
// Out of MyType

type MyType = map[string]int

// Out of MyType
`,
			block: &proto.Block{
				StartLine:   4,
				StartColumn: 6,
				EndLine:     4,
				EndColumn:   29,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		fset := token.NewFileSet()
		f, err := parser.ParseFile(fset, "", tt.src, parser.ParseComments)
		if err != nil {
			panic(err)
		}

		for _, cg := range f.Comments {
			t.Run(tt.name, func(t *testing.T) {
				got := myAst.IsInline(fset, cg, tt.block)
				if diff := cmp.Diff(tt.want, got); diff != "" {
					t.Errorf("bool values are mismatch (-want +got):%s\n", diff)
				}
			})
		}
	}
}
