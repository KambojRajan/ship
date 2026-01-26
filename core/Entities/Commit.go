package entities

import (
	"bytes"
	"fmt"

	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/utils"
)

type Commit struct {
	TreeHash     string
	ParentHashes []string
	Author       User
	Committer    User
	Message      string
}

func NewCommit(treeHash string, parentHashes []string, author User, committer User, message string) *Commit {
	return &Commit{
		TreeHash:     treeHash,
		ParentHashes: parentHashes,
		Author:       author,
		Committer:    committer,
		Message:      message,
	}
}

func (c *Commit) Commit() (string, error) {
	content := c.serialize()
	return utils.HashObject(content, common.COMMIT, true)
}

func (c *Commit) serialize() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s ", utils.TREE))
	buffer.WriteString(c.TreeHash)
	buffer.WriteByte(utils.NEWLINE)

	for _, parentHash := range c.ParentHashes {
		buffer.WriteString(fmt.Sprintf("%s ", utils.PARENT))
		buffer.WriteString(parentHash)
		buffer.WriteByte(utils.NEWLINE)
	}

	buffer.WriteString(fmt.Sprintf("%s ", utils.AUTHOR))
	buffer.WriteString(c.Author.String())
	buffer.WriteByte(utils.NEWLINE)

	buffer.WriteString(fmt.Sprintf("%s ", utils.COMMITTER))
	buffer.WriteString(c.Committer.String())
	buffer.WriteByte(utils.NEWLINE)

	buffer.WriteByte(utils.NEWLINE)

	buffer.WriteString(c.Message)

	return buffer.Bytes()
}

func (c *Commit) CommitTree(treeHash string, parentHashes []string, author User, message string) (string, error) {
	committer := c.Committer
	if committer.Timestamp.IsZero() {
		committer.Timestamp = author.Timestamp
	}

	commit := NewCommit(treeHash, parentHashes, author, committer, message)
	return commit.Commit()
}
