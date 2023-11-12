package ast_json

const MATRIX_SIZE = 93
const (
	NodeID = iota
	PositionNodeID
	CommentNodeID
	CommentGroupNodeID
	FieldNodeID
	FieldNodeAliasID
	FieldListNodeID
	BadExprNodeID
	IdentNodeID
	EllipsisNodeID
	EllipsisNodeAliasID
	BasicLitNodeID
	FuncLitNodeID
	CompositeLitNodeID
	CompositeLitNodeAliasID
	ParenExprNodeID
	ParenExprNodeAliasID
	SelectorExprNodeID
	SelectorExprNodeAliasID
	IndexExprNodeID
	IndexExprNodeAliasID
	IndexListExprNodeID
	IndexListExprNodeAliasID
	SliceExprNodeID
	SliceExprNodeAliasID
	TypeAssertExprNodeID
	TypeAssertExprNodeAliasID
	CallExprNodeID
	CallExprNodeAliasID
	StarExprNodeID
	StarExprNodeAliasID
	UnaryExprNodeID
	UnaryExprNodeAliasID
	BinaryExprNodeID
	BinaryExprNodeAliasID
	KeyValueExprNodeID
	KeyValueExprNodeAliasID
	ArrayTypeNodeID
	ArrayTypeNodeAliasID
	StructTypeNodeID
	FuncTypeNodeID
	InterfaceTypeNodeID
	MapTypeNodeID
	MapTypeNodeAliasID
	ChanTypeNodeID
	ChanTypeNodeAliasID
	BadStmtNodeID
	DeclStmtNodeID
	DeclStmtNodeAliasID
	EmptyStmtNodeID
	LabeledStmtNodeID
	LabeledStmtNodeAliasID
	ExprStmtNodeID
	ExprStmtNodeAliasID
	SendStmtNodeID
	SendStmtNodeAliasID
	IncDecStmtNodeID
	IncDecStmtNodeAliasID
	AssignStmtNodeID
	AssignStmtNodeAliasID
	GoStmtNodeID
	DeferStmtNodeID
	ReturnStmtNodeID
	ReturnStmtNodeAliasID
	BranchStmtNodeID
	BlockStmtNodeID
	BlockStmtNodeAliasID
	IfStmtNodeID
	IfStmtNodeAliasID
	CaseClauseNodeID
	CaseClauseNodeAliasID
	SwitchStmtNodeID
	SwitchStmtNodeAliasID
	TypeSwitchStmtNodeID
	TypeSwitchStmtNodeAliasID
	CommClauseNodeID
	CommClauseNodeAliasID
	SelectStmtNodeID
	ForStmtNodeID
	ForStmtNodeAliasID
	RangeStmtNodeID
	RangeStmtNodeAliasID
	ImportSpecNodeID
	ValueSpecNodeID
	ValueSpecNodeAliasID
	TypeSpecNodeID
	TypeSpecNodeAliasID
	BadDeclNodeID
	GenDeclNodeID
	GenDeclNodeAliasID
	FuncDeclNodeID
	FileNodeID
	FileNodeAliasID
	PackageNodeID
)

var nodeTypesMap map[string]int = map[string]int{
	"Node":                    NodeID,
	"PositionNode":            PositionNodeID,
	"CommentNode":             CommentNodeID,
	"CommentGroupNode":        CommentGroupNodeID,
	"FieldNode":               FieldNodeID,
	"FieldNodeAlias":          FieldNodeAliasID,
	"FieldListNode":           FieldListNodeID,
	"BadExprNode":             BadExprNodeID,
	"IdentNode":               IdentNodeID,
	"EllipsisNode":            EllipsisNodeID,
	"EllipsisNodeAlias":       EllipsisNodeAliasID,
	"BasicLitNode":            BasicLitNodeID,
	"FuncLitNode":             FuncLitNodeID,
	"CompositeLitNode":        CompositeLitNodeID,
	"CompositeLitNodeAlias":   CompositeLitNodeAliasID,
	"ParenExprNode":           ParenExprNodeID,
	"ParenExprNodeAlias":      ParenExprNodeAliasID,
	"SelectorExprNode":        SelectorExprNodeID,
	"SelectorExprNodeAlias":   SelectorExprNodeAliasID,
	"IndexExprNode":           IndexExprNodeID,
	"IndexExprNodeAlias":      IndexExprNodeAliasID,
	"IndexListExprNode":       IndexListExprNodeID,
	"IndexListExprNodeAlias":  IndexListExprNodeAliasID,
	"SliceExprNode":           SliceExprNodeID,
	"SliceExprNodeAlias":      SliceExprNodeAliasID,
	"TypeAssertExprNode":      TypeAssertExprNodeID,
	"TypeAssertExprNodeAlias": TypeAssertExprNodeAliasID,
	"CallExprNode":            CallExprNodeID,
	"CallExprNodeAlias":       CallExprNodeAliasID,
	"StarExprNode":            StarExprNodeID,
	"StarExprNodeAlias":       StarExprNodeAliasID,
	"UnaryExprNode":           UnaryExprNodeID,
	"UnaryExprNodeAlias":      UnaryExprNodeAliasID,
	"BinaryExprNode":          BinaryExprNodeID,
	"BinaryExprNodeAlias":     BinaryExprNodeAliasID,
	"KeyValueExprNode":        KeyValueExprNodeID,
	"KeyValueExprNodeAlias":   KeyValueExprNodeAliasID,
	"ArrayTypeNode":           ArrayTypeNodeID,
	"ArrayTypeNodeAlias":      ArrayTypeNodeAliasID,
	"StructTypeNode":          StructTypeNodeID,
	"FuncTypeNode":            FuncTypeNodeID,
	"InterfaceTypeNode":       InterfaceTypeNodeID,
	"MapTypeNode":             MapTypeNodeID,
	"MapTypeNodeAlias":        MapTypeNodeAliasID,
	"ChanTypeNode":            ChanTypeNodeID,
	"ChanTypeNodeAlias":       ChanTypeNodeAliasID,
	"BadStmtNode":             BadStmtNodeID,
	"DeclStmtNode":            DeclStmtNodeID,
	"DeclStmtNodeAlias":       DeclStmtNodeAliasID,
	"EmptyStmtNode":           EmptyStmtNodeID,
	"LabeledStmtNode":         LabeledStmtNodeID,
	"LabeledStmtNodeAlias":    LabeledStmtNodeAliasID,
	"ExprStmtNode":            ExprStmtNodeID,
	"ExprStmtNodeAlias":       ExprStmtNodeAliasID,
	"SendStmtNode":            SendStmtNodeID,
	"SendStmtNodeAlias":       SendStmtNodeAliasID,
	"IncDecStmtNode":          IncDecStmtNodeID,
	"IncDecStmtNodeAlias":     IncDecStmtNodeAliasID,
	"AssignStmtNode":          AssignStmtNodeID,
	"AssignStmtNodeAlias":     AssignStmtNodeAliasID,
	"GoStmtNode":              GoStmtNodeID,
	"DeferStmtNode":           DeferStmtNodeID,
	"ReturnStmtNode":          ReturnStmtNodeID,
	"ReturnStmtNodeAlias":     ReturnStmtNodeAliasID,
	"BranchStmtNode":          BranchStmtNodeID,
	"BlockStmtNode":           BlockStmtNodeID,
	"BlockStmtNodeAlias":      BlockStmtNodeAliasID,
	"IfStmtNode":              IfStmtNodeID,
	"IfStmtNodeAlias":         IfStmtNodeAliasID,
	"CaseClauseNode":          CaseClauseNodeID,
	"CaseClauseNodeAlias":     CaseClauseNodeAliasID,
	"SwitchStmtNode":          SwitchStmtNodeID,
	"SwitchStmtNodeAlias":     SwitchStmtNodeAliasID,
	"TypeSwitchStmtNode":      TypeSwitchStmtNodeID,
	"TypeSwitchStmtNodeAlias": TypeSwitchStmtNodeAliasID,
	"CommClauseNode":          CommClauseNodeID,
	"CommClauseNodeAlias":     CommClauseNodeAliasID,
	"SelectStmtNode":          SelectStmtNodeID,
	"ForStmtNode":             ForStmtNodeID,
	"ForStmtNodeAlias":        ForStmtNodeAliasID,
	"RangeStmtNode":           RangeStmtNodeID,
	"RangeStmtNodeAlias":      RangeStmtNodeAliasID,
	"ImportSpecNode":          ImportSpecNodeID,
	"ValueSpecNode":           ValueSpecNodeID,
	"ValueSpecNodeAlias":      ValueSpecNodeAliasID,
	"TypeSpecNode":            TypeSpecNodeID,
	"TypeSpecNodeAlias":       TypeSpecNodeAliasID,
	"BadDeclNode":             BadDeclNodeID,
	"GenDeclNode":             GenDeclNodeID,
	"GenDeclNodeAlias":        GenDeclNodeAliasID,
	"FuncDeclNode":            FuncDeclNodeID,
	"FileNode":                FileNodeID,
	"FileNodeAlias":           FileNodeAliasID,
	"PackageNode":             PackageNodeID,
}
