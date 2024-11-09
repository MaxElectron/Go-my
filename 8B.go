//go:build !solution

package treeiter

func DoInOrder[nodeTp interface {
	Left() *nodeTp
	Right() *nodeTp
}](node *nodeTp, action func(*nodeTp)) {
	if node == nil {
		return
	}

	DoInOrder((*node).Left(), action)
	action(node)
	DoInOrder((*node).Right(), action)
}
