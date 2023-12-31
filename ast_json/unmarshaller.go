package ast_json

import (
	"go/ast"
	"go/token"
)

var StringToToken = map[string]token.Token{}

func init() {
	for t := token.ILLEGAL; t <= token.TILDE; t++ {
		StringToToken[t.String()] = t
	}
}

type Unmarshaller struct {
	Options
	fset       *token.FileSet
	references map[int]any
}

func NewUnmarshaller(options Options) *Unmarshaller {
	return &Unmarshaller{
		Options:    options,
		fset:       token.NewFileSet(),
		references: make(map[int]any),
	}
}

func wrapUnmarshal[T INode, R any](um *Unmarshaller, node *T, marshal func() *R) *R {
	if node == nil {
		return nil
	}

	if !um.WithReferences {
		return marshal()
	}

	refId := (*node).GetRefId()
	if refId == 0 {
		return marshal()
	}

	if ref, ok := um.references[refId]; ok {
		return ref.(*R)
	}
	result := marshal()
	um.references[refId] = result
	return result
}

func (um *Unmarshaller) FileSet() *token.FileSet {
	return um.fset
}

func (um *Unmarshaller) UnmarshalPositionNode(node *PositionNode) token.Pos {
	if !um.WithPositions {
		return token.NoPos
	}
	if node == nil {
		return token.NoPos
	}
	pos := token.NoPos
	um.fset.Iterate(func(f *token.File) bool {
		if f.Name() == node.Filename {
			pos = f.Pos(node.Offset)
			return false
		}
		return true
	})
	return pos
}

func (um *Unmarshaller) UnmarshalCommentNode(node *CommentNode) *ast.Comment {
	return wrapUnmarshal(um, node, func() *ast.Comment {
		return &ast.Comment{
			Slash: um.UnmarshalPositionNode(node.Slash),
			Text:  node.Text,
		}
	})
}

func (um *Unmarshaller) UnmarshalCommentNodes(nodes []*CommentNode) []*ast.Comment {
	if nodes == nil {
		return nil
	}
	comments := make([]*ast.Comment, len(nodes))
	for index, node := range nodes {
		comments[index] = um.UnmarshalCommentNode(node)
	}
	return comments
}

func (um *Unmarshaller) UnmarshalCommentGroupNode(node *CommentGroupNode) *ast.CommentGroup {
	if !um.WithComments {
		return nil
	}
	return wrapUnmarshal(um, node, func() *ast.CommentGroup {
		return &ast.CommentGroup{
			List: um.UnmarshalCommentNodes(node.List),
		}
	})
}

func (um *Unmarshaller) UnmarshalCommentGroupNodes(comments []*CommentGroupNode) []*ast.CommentGroup {
	if comments == nil {
		return nil
	}
	commentGroups := make([]*ast.CommentGroup, len(comments))
	for i, comment := range comments {
		commentGroups[i] = um.UnmarshalCommentGroupNode(comment)
	}
	return commentGroups
}

func (um *Unmarshaller) UnmarshalFieldNode(node *FieldNode) *ast.Field {
	return wrapUnmarshal(um, node, func() *ast.Field {
		return &ast.Field{
			Doc:     um.UnmarshalCommentGroupNode(node.Doc),
			Names:   um.UnmarshalIdentNodes(node.Names),
			Type:    um.UnmarshalExpr(node.Type),
			Tag:     um.UnmarshalBasicLitNode(node.Tag),
			Comment: um.UnmarshalCommentGroupNode(node.Comment),
		}
	})
}

func (um *Unmarshaller) UnmarshalFieldNodes(nodes []*FieldNode) []*ast.Field {
	if nodes == nil {
		return nil
	}
	fields := make([]*ast.Field, len(nodes))
	for index, node := range nodes {
		fields[index] = um.UnmarshalFieldNode(node)
	}
	return fields
}

func (um *Unmarshaller) UnmarshalFieldListNode(node *FieldListNode) *ast.FieldList {
	return wrapUnmarshal(um, node, func() *ast.FieldList {
		return &ast.FieldList{
			Opening: um.UnmarshalPositionNode(node.Opening),
			List:    um.UnmarshalFieldNodes(node.List),
			Closing: um.UnmarshalPositionNode(node.Closing),
		}
	})
}

func (um *Unmarshaller) UnmarshalBadExprNode(node *BadExprNode) *ast.BadExpr {
	return wrapUnmarshal(um, node, func() *ast.BadExpr {
		return &ast.BadExpr{
			From: um.UnmarshalPositionNode(node.From),
			To:   um.UnmarshalPositionNode(node.To),
		}
	})
}

func (um *Unmarshaller) UnmarshalIdentNode(node *IdentNode) *ast.Ident {
	return wrapUnmarshal(um, node, func() *ast.Ident {
		return &ast.Ident{
			NamePos: um.UnmarshalPositionNode(node.NamePos),
			Name:    node.Name,
		}
	})
}

func (um *Unmarshaller) UnmarshalIdentNodes(nodes []*IdentNode) []*ast.Ident {
	if nodes == nil {
		return nil
	}
	idents := make([]*ast.Ident, len(nodes))
	for index, node := range nodes {
		idents[index] = um.UnmarshalIdentNode(node)
	}
	return idents
}

func (um *Unmarshaller) UnmarshalEllipsisNode(node *EllipsisNode) *ast.Ellipsis {
	return wrapUnmarshal(um, node, func() *ast.Ellipsis {
		return &ast.Ellipsis{
			Ellipsis: um.UnmarshalPositionNode(node.Ellipsis),
			Elt:      um.UnmarshalExpr(node.Elt),
		}
	})
}
func (um *Unmarshaller) UnmarshalBasicLitNode(node *BasicLitNode) *ast.BasicLit {
	return wrapUnmarshal(um, node, func() *ast.BasicLit {
		kind, ok := StringToToken[node.Kind]
		if !ok {
			panic("unsupported token kind " + node.Kind)
		}
		return &ast.BasicLit{
			ValuePos: um.UnmarshalPositionNode(node.ValuePos),
			Kind:     kind,
			Value:    node.Value,
		}
	})
}

func (um *Unmarshaller) UnmarshalFuncLitNode(node *FuncLitNode) *ast.FuncLit {
	return wrapUnmarshal(um, node, func() *ast.FuncLit {
		return &ast.FuncLit{
			Type: um.UnmarshalFuncTypeNode(node.Type),
			Body: um.UnmarshalBlockStmtNode(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalCompositeLitNode(node *CompositeLitNode) *ast.CompositeLit {
	return wrapUnmarshal(um, node, func() *ast.CompositeLit {
		return &ast.CompositeLit{
			Type:       um.UnmarshalExpr(node.Type),
			Lbrace:     um.UnmarshalPositionNode(node.Lbrace),
			Elts:       um.UnmarshalExprNodes(node.Elts),
			Rbrace:     um.UnmarshalPositionNode(node.Rbrace),
			Incomplete: node.Incomplete,
		}
	})
}

func (um *Unmarshaller) UnmarshalParenExprNode(node *ParenExprNode) *ast.ParenExpr {
	return wrapUnmarshal(um, node, func() *ast.ParenExpr {
		return &ast.ParenExpr{
			Lparen: um.UnmarshalPositionNode(node.Lparen),
			X:      um.UnmarshalExpr(node.X),
			Rparen: um.UnmarshalPositionNode(node.Rparen),
		}
	})
}

func (um *Unmarshaller) UnmarshalSelectorExprNode(node *SelectorExprNode) *ast.SelectorExpr {
	return wrapUnmarshal(um, node, func() *ast.SelectorExpr {
		return &ast.SelectorExpr{
			X:   um.UnmarshalExpr(node.X),
			Sel: um.UnmarshalIdentNode(node.Sel),
		}
	})
}

func (um *Unmarshaller) UnmarshalIndexExprNode(node *IndexExprNode) *ast.IndexExpr {
	return wrapUnmarshal(um, node, func() *ast.IndexExpr {
		return &ast.IndexExpr{
			X:      um.UnmarshalExpr(node.X),
			Lbrack: um.UnmarshalPositionNode(node.Lbrack),
			Index:  um.UnmarshalExpr(node.Index),
			Rbrack: um.UnmarshalPositionNode(node.Rbrack),
		}
	})
}

func (um *Unmarshaller) UnmarshalExprNodes(nodes []IExprNode) []ast.Expr {
	if nodes == nil {
		return nil
	}
	exprs := make([]ast.Expr, len(nodes))
	for index, node := range nodes {
		exprs[index] = um.UnmarshalExpr(node)
	}
	return exprs
}

func (um *Unmarshaller) UnmarshalIndexListExprNode(node *IndexListExprNode) *ast.IndexListExpr {
	return wrapUnmarshal(um, node, func() *ast.IndexListExpr {
		return &ast.IndexListExpr{
			X:       um.UnmarshalExpr(node.X),
			Lbrack:  um.UnmarshalPositionNode(node.Lbrack),
			Indices: um.UnmarshalExprNodes(node.Indices),
			Rbrack:  um.UnmarshalPositionNode(node.Rbrack),
		}
	})
}

func (um *Unmarshaller) UnmarshalSliceExprNode(node *SliceExprNode) *ast.SliceExpr {
	return wrapUnmarshal(um, node, func() *ast.SliceExpr {
		return &ast.SliceExpr{
			X:      um.UnmarshalExpr(node.X),
			Lbrack: um.UnmarshalPositionNode(node.Lbrack),
			Low:    um.UnmarshalExpr(node.Low),
			High:   um.UnmarshalExpr(node.High),
			Max:    um.UnmarshalExpr(node.Max),
			Slice3: node.Slice3,
			Rbrack: um.UnmarshalPositionNode(node.Rbrack),
		}
	})
}

func (um *Unmarshaller) UnmarshalTypeAssertExprNode(node *TypeAssertExprNode) *ast.TypeAssertExpr {
	return wrapUnmarshal(um, node, func() *ast.TypeAssertExpr {
		return &ast.TypeAssertExpr{
			X:      um.UnmarshalExpr(node.X),
			Lparen: um.UnmarshalPositionNode(node.Lparen),
			Type:   um.UnmarshalExpr(node.Type),
			Rparen: um.UnmarshalPositionNode(node.Rparen),
		}
	})
}

func (um *Unmarshaller) UnmarshalCallExprNode(node *CallExprNode) *ast.CallExpr {
	return wrapUnmarshal(um, node, func() *ast.CallExpr {
		return &ast.CallExpr{
			Fun:      um.UnmarshalExpr(node.Fun),
			Lparen:   um.UnmarshalPositionNode(node.Lparen),
			Args:     um.UnmarshalExprNodes(node.Args),
			Ellipsis: um.UnmarshalPositionNode(node.Ellipsis),
			Rparen:   um.UnmarshalPositionNode(node.Rparen),
		}
	})
}

func (um *Unmarshaller) UnmarshalStarExprNode(node *StarExprNode) *ast.StarExpr {
	return wrapUnmarshal(um, node, func() *ast.StarExpr {
		return &ast.StarExpr{
			Star: um.UnmarshalPositionNode(node.Star),
			X:    um.UnmarshalExpr(node.X),
		}
	})
}

func (um *Unmarshaller) UnmarshalUnaryExprNode(node *UnaryExprNode) *ast.UnaryExpr {
	return wrapUnmarshal(um, node, func() *ast.UnaryExpr {
		return &ast.UnaryExpr{
			OpPos: um.UnmarshalPositionNode(node.OpPos),
			Op:    StringToToken[node.Op],
			X:     um.UnmarshalExpr(node.X),
		}
	})
}

func (um *Unmarshaller) UnmarshalBinaryExprNode(node *BinaryExprNode) *ast.BinaryExpr {
	return wrapUnmarshal(um, node, func() *ast.BinaryExpr {
		return &ast.BinaryExpr{
			X:     um.UnmarshalExpr(node.X),
			OpPos: um.UnmarshalPositionNode(node.OpPos),
			Op:    StringToToken[node.Op],
			Y:     um.UnmarshalExpr(node.Y),
		}
	})
}

func (um *Unmarshaller) UnmarshalKeyValueExprNode(node *KeyValueExprNode) *ast.KeyValueExpr {
	return wrapUnmarshal(um, node, func() *ast.KeyValueExpr {
		return &ast.KeyValueExpr{
			Key:   um.UnmarshalExpr(node.Key),
			Colon: um.UnmarshalPositionNode(node.Colon),
			Value: um.UnmarshalExpr(node.Value),
		}
	})
}

func (um *Unmarshaller) UnmarshalArrayTypeNode(node *ArrayTypeNode) *ast.ArrayType {
	return wrapUnmarshal(um, node, func() *ast.ArrayType {
		return &ast.ArrayType{
			Lbrack: um.UnmarshalPositionNode(node.Lbrack),
			Len:    um.UnmarshalExpr(node.Len),
			Elt:    um.UnmarshalExpr(node.Elt),
		}
	})
}

func (um *Unmarshaller) UnmarshalStructTypeNode(node *StructTypeNode) *ast.StructType {
	return wrapUnmarshal(um, node, func() *ast.StructType {
		return &ast.StructType{
			Struct:     um.UnmarshalPositionNode(node.Struct),
			Fields:     um.UnmarshalFieldListNode(node.Fields),
			Incomplete: node.Incomplete,
		}
	})
}

func (um *Unmarshaller) UnmarshalFuncTypeNode(node *FuncTypeNode) *ast.FuncType {
	return wrapUnmarshal(um, node, func() *ast.FuncType {
		return &ast.FuncType{
			Func:       um.UnmarshalPositionNode(node.Func),
			TypeParams: um.UnmarshalFieldListNode(node.TypeParams),
			Params:     um.UnmarshalFieldListNode(node.Params),
			Results:    um.UnmarshalFieldListNode(node.Results),
		}
	})
}

func (um *Unmarshaller) UnmarshalInterfaceTypeNode(node *InterfaceTypeNode) *ast.InterfaceType {
	return wrapUnmarshal(um, node, func() *ast.InterfaceType {
		return &ast.InterfaceType{
			Interface:  um.UnmarshalPositionNode(node.Interface),
			Methods:    um.UnmarshalFieldListNode(node.Methods),
			Incomplete: node.Incomplete,
		}
	})
}

func (um *Unmarshaller) UnmarshalMapTypeNode(node *MapTypeNode) *ast.MapType {
	return wrapUnmarshal(um, node, func() *ast.MapType {
		return &ast.MapType{
			Map:   um.UnmarshalPositionNode(node.Map),
			Key:   um.UnmarshalExpr(node.Key),
			Value: um.UnmarshalExpr(node.Value),
		}
	})
}

var StringToChanDir = map[string]ast.ChanDir{
	"SEND": ast.SEND,
	"RECV": ast.RECV,
	"BOTH": ast.SEND | ast.RECV,
}

func (um *Unmarshaller) UnmarshalChanTypeNode(node *ChanTypeNode) *ast.ChanType {
	return wrapUnmarshal(um, node, func() *ast.ChanType {
		return &ast.ChanType{
			Begin: um.UnmarshalPositionNode(node.Begin),
			Arrow: um.UnmarshalPositionNode(node.Arrow),
			Dir:   StringToChanDir[node.Dir],
			Value: um.UnmarshalExpr(node.Value),
		}
	})
}

func (um *Unmarshaller) UnmarshalBadStmtNode(node *BadStmtNode) *ast.BadStmt {
	return wrapUnmarshal(um, node, func() *ast.BadStmt {
		return &ast.BadStmt{
			From: um.UnmarshalPositionNode(node.From),
			To:   um.UnmarshalPositionNode(node.To),
		}
	})
}

func (um *Unmarshaller) UnmarshalDeclStmtNode(node *DeclStmtNode) *ast.DeclStmt {
	return wrapUnmarshal(um, node, func() *ast.DeclStmt {
		return &ast.DeclStmt{
			Decl: um.UnmarshalDecl(node.Decl),
		}
	})
}

func (um *Unmarshaller) UnmarshalEmptyStmtNode(node *EmptyStmtNode) *ast.EmptyStmt {
	return wrapUnmarshal(um, node, func() *ast.EmptyStmt {
		return &ast.EmptyStmt{
			Semicolon: um.UnmarshalPositionNode(node.Semicolon),
			Implicit:  node.Implicit,
		}
	})
}

func (um *Unmarshaller) UnmarshalLabeledStmtNode(node *LabeledStmtNode) *ast.LabeledStmt {
	return wrapUnmarshal(um, node, func() *ast.LabeledStmt {
		return &ast.LabeledStmt{
			Label: um.UnmarshalIdentNode(node.Label),
			Colon: um.UnmarshalPositionNode(node.Colon),
			Stmt:  um.UnmarshalStmt(node.Stmt),
		}
	})
}

func (um *Unmarshaller) UnmarshalExprStmtNode(node *ExprStmtNode) *ast.ExprStmt {
	return wrapUnmarshal(um, node, func() *ast.ExprStmt {
		return &ast.ExprStmt{
			X: um.UnmarshalExpr(node.X),
		}
	})
}

func (um *Unmarshaller) UnmarshalSendStmtNode(node *SendStmtNode) *ast.SendStmt {
	return wrapUnmarshal(um, node, func() *ast.SendStmt {
		return &ast.SendStmt{
			Chan:  um.UnmarshalExpr(node.Chan),
			Arrow: um.UnmarshalPositionNode(node.Arrow),
			Value: um.UnmarshalExpr(node.Value),
		}
	})
}

func (um *Unmarshaller) UnmarshalIncDecStmtNode(node *IncDecStmtNode) *ast.IncDecStmt {
	return wrapUnmarshal(um, node, func() *ast.IncDecStmt {
		return &ast.IncDecStmt{
			X:      um.UnmarshalExpr(node.X),
			TokPos: um.UnmarshalPositionNode(node.TokPos),
			Tok:    StringToToken[node.Tok],
		}
	})
}

func (um *Unmarshaller) UnmarshalAssignStmtNode(node *AssignStmtNode) *ast.AssignStmt {
	return wrapUnmarshal(um, node, func() *ast.AssignStmt {
		return &ast.AssignStmt{
			Lhs:    um.UnmarshalExprNodes(node.Lhs),
			TokPos: um.UnmarshalPositionNode(node.TokPos),
			Tok:    StringToToken[node.Tok],
			Rhs:    um.UnmarshalExprNodes(node.Rhs),
		}
	})
}

func (um *Unmarshaller) UnmarshalGoStmtNode(node *GoStmtNode) *ast.GoStmt {
	return wrapUnmarshal(um, node, func() *ast.GoStmt {
		return &ast.GoStmt{
			Go:   um.UnmarshalPositionNode(node.Go),
			Call: um.UnmarshalCallExprNode(node.Call),
		}
	})
}

func (um *Unmarshaller) UnmarshalDeferStmtNode(node *DeferStmtNode) *ast.DeferStmt {
	return wrapUnmarshal(um, node, func() *ast.DeferStmt {
		return &ast.DeferStmt{
			Defer: um.UnmarshalPositionNode(node.Defer),
			Call:  um.UnmarshalCallExprNode(node.Call),
		}
	})
}

func (um *Unmarshaller) UnmarshalReturnStmtNode(node *ReturnStmtNode) *ast.ReturnStmt {
	return wrapUnmarshal(um, node, func() *ast.ReturnStmt {
		return &ast.ReturnStmt{
			Return:  um.UnmarshalPositionNode(node.Return),
			Results: um.UnmarshalExprNodes(node.Results),
		}
	})
}

func (um *Unmarshaller) UnmarshalBranchStmtNode(node *BranchStmtNode) *ast.BranchStmt {
	return wrapUnmarshal(um, node, func() *ast.BranchStmt {
		return &ast.BranchStmt{
			TokPos: um.UnmarshalPositionNode(node.TokPos),
			Tok:    StringToToken[node.Tok],
			Label:  um.UnmarshalIdentNode(node.Label),
		}
	})
}

func (um *Unmarshaller) UnmarshalStmtNodes(nodes []IStmtNode) []ast.Stmt {
	if nodes == nil {
		return nil
	}
	stmts := make([]ast.Stmt, len(nodes))
	for i, node := range nodes {
		stmts[i] = um.UnmarshalStmt(node)
	}
	return stmts
}

func (um *Unmarshaller) UnmarshalBlockStmtNode(node *BlockStmtNode) *ast.BlockStmt {
	return wrapUnmarshal(um, node, func() *ast.BlockStmt {
		return &ast.BlockStmt{
			Lbrace: um.UnmarshalPositionNode(node.Lbrace),
			List:   um.UnmarshalStmtNodes(node.List),
			Rbrace: um.UnmarshalPositionNode(node.Rbrace),
		}
	})
}

func (um *Unmarshaller) UnmarshalIfStmtNode(node *IfStmtNode) *ast.IfStmt {
	return wrapUnmarshal(um, node, func() *ast.IfStmt {
		return &ast.IfStmt{
			If:   um.UnmarshalPositionNode(node.If),
			Init: um.UnmarshalStmt(node.Init),
			Cond: um.UnmarshalExpr(node.Cond),
			Body: um.UnmarshalBlockStmtNode(node.Body),
			Else: um.UnmarshalStmt(node.Else),
		}
	})
}

func (um *Unmarshaller) UnmarshalCaseClauseNode(node *CaseClauseNode) *ast.CaseClause {
	return wrapUnmarshal(um, node, func() *ast.CaseClause {
		return &ast.CaseClause{
			Case:  um.UnmarshalPositionNode(node.Case),
			List:  um.UnmarshalExprNodes(node.List),
			Colon: um.UnmarshalPositionNode(node.Colon),
			Body:  um.UnmarshalStmtNodes(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalSwitchStmtNode(node *SwitchStmtNode) *ast.SwitchStmt {
	return wrapUnmarshal(um, node, func() *ast.SwitchStmt {
		return &ast.SwitchStmt{
			Switch: um.UnmarshalPositionNode(node.Switch),
			Init:   um.UnmarshalStmt(node.Init),
			Tag:    um.UnmarshalExpr(node.Tag),
			Body:   um.UnmarshalBlockStmtNode(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalTypeSwitchStmtNode(node *TypeSwitchStmtNode) *ast.TypeSwitchStmt {
	return wrapUnmarshal(um, node, func() *ast.TypeSwitchStmt {
		return &ast.TypeSwitchStmt{
			Switch: um.UnmarshalPositionNode(node.Switch),
			Init:   um.UnmarshalStmt(node.Init),
			Assign: um.UnmarshalStmt(node.Assign),
			Body:   um.UnmarshalBlockStmtNode(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalCommClauseNode(node *CommClauseNode) *ast.CommClause {
	return wrapUnmarshal(um, node, func() *ast.CommClause {
		return &ast.CommClause{
			Case:  um.UnmarshalPositionNode(node.Case),
			Comm:  um.UnmarshalStmt(node.Comm),
			Colon: um.UnmarshalPositionNode(node.Colon),
			Body:  um.UnmarshalStmtNodes(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalSelectStmtNode(node *SelectStmtNode) *ast.SelectStmt {
	return wrapUnmarshal(um, node, func() *ast.SelectStmt {
		return &ast.SelectStmt{
			Select: um.UnmarshalPositionNode(node.Select),
			Body:   um.UnmarshalBlockStmtNode(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalForStmtNode(node *ForStmtNode) *ast.ForStmt {
	return wrapUnmarshal(um, node, func() *ast.ForStmt {
		return &ast.ForStmt{
			For:  um.UnmarshalPositionNode(node.For),
			Init: um.UnmarshalStmt(node.Init),
			Cond: um.UnmarshalExpr(node.Cond),
			Post: um.UnmarshalStmt(node.Post),
			Body: um.UnmarshalBlockStmtNode(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalRangeStmtNode(node *RangeStmtNode) *ast.RangeStmt {
	return wrapUnmarshal(um, node, func() *ast.RangeStmt {
		return &ast.RangeStmt{
			For:    um.UnmarshalPositionNode(node.For),
			Key:    um.UnmarshalExpr(node.Key),
			Value:  um.UnmarshalExpr(node.Value),
			TokPos: um.UnmarshalPositionNode(node.TokPos),
			Tok:    StringToToken[node.Tok],
			X:      um.UnmarshalExpr(node.X),
			Body:   um.UnmarshalBlockStmtNode(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalImportSpecNode(node *ImportSpecNode) *ast.ImportSpec {
	return wrapUnmarshal(um, node, func() *ast.ImportSpec {
		return &ast.ImportSpec{
			Doc:     um.UnmarshalCommentGroupNode(node.Doc),
			Name:    um.UnmarshalIdentNode(node.Name),
			Path:    um.UnmarshalBasicLitNode(node.Path),
			Comment: um.UnmarshalCommentGroupNode(node.Comment),
			EndPos:  um.UnmarshalPositionNode(node.EndPos),
		}
	})
}

func (um *Unmarshaller) UnmarshalImportSpecNodes(imports []*ImportSpecNode) []*ast.ImportSpec {
	if imports == nil {
		return nil
	}
	specs := make([]*ast.ImportSpec, len(imports))
	for i, node := range imports {
		specs[i] = um.UnmarshalImportSpecNode(node)
	}
	return specs
}

func (um *Unmarshaller) UnmarshalValueSpecNode(node *ValueSpecNode) *ast.ValueSpec {
	return wrapUnmarshal(um, node, func() *ast.ValueSpec {
		return &ast.ValueSpec{
			Doc:     um.UnmarshalCommentGroupNode(node.Doc),
			Names:   um.UnmarshalIdentNodes(node.Names),
			Type:    um.UnmarshalExpr(node.Type),
			Values:  um.UnmarshalExprNodes(node.Values),
			Comment: um.UnmarshalCommentGroupNode(node.Comment),
		}
	})
}

func (um *Unmarshaller) UnmarshalTypeSpecNode(node *TypeSpecNode) *ast.TypeSpec {
	return wrapUnmarshal(um, node, func() *ast.TypeSpec {
		return &ast.TypeSpec{
			Doc:        um.UnmarshalCommentGroupNode(node.Doc),
			Name:       um.UnmarshalIdentNode(node.Name),
			TypeParams: um.UnmarshalFieldListNode(node.TypeParams),
			Assign:     um.UnmarshalPositionNode(node.Assign),
			Type:       um.UnmarshalExpr(node.Type),
			Comment:    um.UnmarshalCommentGroupNode(node.Comment),
		}
	})
}

func (um *Unmarshaller) UnmarshalSpecNodes(nodes []ISpecNode) []ast.Spec {
	if nodes == nil {
		return nil
	}
	specs := make([]ast.Spec, len(nodes))
	for i, node := range nodes {
		specs[i] = um.UnmarshalSpec(node)
	}
	return specs
}

func (um *Unmarshaller) UnmarshalBadDeclNode(node *BadDeclNode) *ast.BadDecl {
	return wrapUnmarshal(um, node, func() *ast.BadDecl {
		return &ast.BadDecl{
			From: um.UnmarshalPositionNode(node.From),
			To:   um.UnmarshalPositionNode(node.To),
		}
	})
}

func (um *Unmarshaller) UnmarshalGenDeclNode(node *GenDeclNode) *ast.GenDecl {
	return wrapUnmarshal(um, node, func() *ast.GenDecl {
		return &ast.GenDecl{
			Doc:    um.UnmarshalCommentGroupNode(node.Doc),
			TokPos: um.UnmarshalPositionNode(node.TokPos),
			Tok:    StringToToken[node.Tok],
			Lparen: um.UnmarshalPositionNode(node.Lparen),
			Specs:  um.UnmarshalSpecNodes(node.Specs),
			Rparen: um.UnmarshalPositionNode(node.Rparen),
		}
	})
}

func (um *Unmarshaller) UnmarshalFuncDeclNode(node *FuncDeclNode) *ast.FuncDecl {
	return wrapUnmarshal(um, node, func() *ast.FuncDecl {
		return &ast.FuncDecl{
			Doc:  um.UnmarshalCommentGroupNode(node.Doc),
			Recv: um.UnmarshalFieldListNode(node.Recv),
			Name: um.UnmarshalIdentNode(node.Name),
			Type: um.UnmarshalFuncTypeNode(node.Type),
			Body: um.UnmarshalBlockStmtNode(node.Body),
		}
	})
}

func (um *Unmarshaller) UnmarshalDeclNodes(nodes []IDeclNode) []ast.Decl {
	if nodes == nil {
		return nil
	}
	decls := make([]ast.Decl, len(nodes))
	for i, node := range nodes {
		decls[i] = um.UnmarshalDecl(node)
	}
	return decls
}

func (um *Unmarshaller) UnmarshalFileNode(node *FileNode) *ast.File {
	return wrapUnmarshal(um, node, func() *ast.File {
		um.fset = node.FileSet
		var imports []*ast.ImportSpec = nil
		if um.WithImports {
			imports = um.UnmarshalImportSpecNodes(node.Imports)
		}
		return &ast.File{
			Doc:        um.UnmarshalCommentGroupNode(node.Doc),
			Package:    um.UnmarshalPositionNode(node.Package),
			Name:       um.UnmarshalIdentNode(node.Name),
			Decls:      um.UnmarshalDeclNodes(node.Decls),
			Imports:    imports,
			Unresolved: um.UnmarshalIdentNodes(node.Unresolved),
			Comments:   um.UnmarshalCommentGroupNodes(node.Comments),
		}
	})
}

func (um *Unmarshaller) UnmarshalExpr(expr IExprNode) ast.Expr {
	if expr == nil {
		return nil
	}
	return expr.UnmarshalExpr(um)
}

func (um *Unmarshaller) UnmarshalStmt(stmt IStmtNode) ast.Stmt {
	if stmt == nil {
		return nil
	}
	return stmt.UnmarshalStmt(um)
}

func (um *Unmarshaller) UnmarshalSpec(spec ISpecNode) ast.Spec {
	if spec == nil {
		return nil
	}
	return spec.UnmarshalSpec(um)
}

func (um *Unmarshaller) UnmarshalDecl(decl IDeclNode) ast.Decl {
	if decl == nil {
		return nil
	}
	return decl.UnmarshalDecl(um)
}
