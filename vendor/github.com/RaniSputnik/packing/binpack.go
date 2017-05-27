package packing

type BinPacker struct {
	root *node
}

func (b *BinPacker) Fit(w int, h int, blocks ...Block) error {
	b.root = &node{x: 0, y: 0, w: w, h: h}

	var err error
	for _, block := range blocks {
		bw, bh := block.Size()
		if bw > b.root.w || bh > b.root.h {
			return ErrInputTooLarge
		}

		if n := b.findNode(b.root, bw, bh); n != nil {
			b.splitNode(n, bw, bh)
			block.Place(n.x, n.y)
		} else {
			err = ErrOutOfRoom
		}
	}

	return err
}

func (b *BinPacker) findNode(root *node, w int, h int) *node {
	if root.used {
		if r := b.findNode(root.right, w, h); r != nil {
			return r
		}
		return b.findNode(root.down, w, h)
	} else if (w <= root.w) && (h <= root.h) {
		return root
	} else {
		return nil
	}
}

func (b *BinPacker) splitNode(n *node, w int, h int) {
	n.used = true
	n.right = &node{x: n.x + w, y: n.y, w: n.w - w, h: h}
	n.down = &node{x: n.x, y: n.y + h, w: n.w, h: n.h - h}
}
