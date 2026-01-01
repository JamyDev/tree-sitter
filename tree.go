package sitter

import (
	"fmt"
)

type Tree struct {
	ts *TreeSitter
	t  uint64
}

func newTree(ts *TreeSitter, t uint64) *Tree {
	return &Tree{ts, t}
}

func (t *Tree) RootNode() (*Node, error) {
	// allocate tsnode 24 bytes
	nodePtr, err := t.ts.call(_malloc, 24)
	if err != nil {
		return nil, fmt.Errorf("allocating node: %w", err)
	}

	// Capture return values to debug WASM calling convention
	retVals, err := t.ts.call(_treeRootNode, nodePtr[0], t.t)
	if err != nil {
		return nil, fmt.Errorf("getting tree root node: %w", err)
	}

	// Debug: Log return values
	fmt.Fprintf(t.ts.out, "DEBUG: ts_tree_root_node returned %d values: %v\n", len(retVals), retVals)

	return newNode(t.ts, nodePtr[0]), nil
}
