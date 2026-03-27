package entities

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/utils"
)

type Commit struct {
	Hash         string
	TreeHash     string
	ParentHashes []string
	Author       User
	Committer    User
	Message      string
	repoPath     string
}

func NewCommit(treeHash string, parentHashes []string, author User, committer User, message string, repoPath string) *Commit {

	return &Commit{
		TreeHash:     treeHash,
		ParentHashes: parentHashes,
		Author:       author,
		Committer:    committer,
		Message:      message,
		repoPath:     repoPath,
	}
}

func (c *Commit) Commit() (string, error) {

	oldCwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if c.repoPath != "" {
		if err := os.Chdir(c.repoPath); err != nil {
			return "", err
		}
		defer func() {
			os.Chdir(oldCwd)
		}()
	}

	content := c.serialize()

	hash, err := utils.HashObject(content, common.COMMIT, true)
	if err != nil {
		return "", err
	}

	return hash, nil
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

	commit := NewCommit(treeHash, parentHashes, author, committer, message, c.repoPath)
	hash, err := commit.Commit()
	if err != nil {
		return "", err
	}

	return hash, nil
}

func inflateGitObject(data []byte) ([]byte, error) {

	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func(r io.ReadCloser) {
		r.Close()
	}(r)

	result, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func parseCommit(data []byte) (*Commit, error) {

	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		return nil, fmt.Errorf("invalid commit object: no null byte found")
	}

	content := data[nullIndex+1:]

	lines := bytes.Split(content, []byte{utils.NEWLINE})

	commit := &Commit{
		ParentHashes: []string{},
	}

	var messageStartIndex int
	for i, line := range lines {
		if len(line) == 0 {
			messageStartIndex = i + 1
			break
		}

		parts := bytes.SplitN(line, []byte{' '}, 2)
		if len(parts) < 2 {
			continue
		}

		key := string(parts[0])
		value := string(parts[1])

		switch key {
		case utils.TREE:
			commit.TreeHash = value
		case utils.PARENT:
			commit.ParentHashes = append(commit.ParentHashes, value)
		case utils.AUTHOR:
			user, err := parseUser(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse author: %w", err)
			}
			commit.Author = user
		case utils.COMMITTER:
			user, err := parseUser(value)
			if err != nil {
				return nil, fmt.Errorf("failed to parse committer: %w", err)
			}
			commit.Committer = user
		}
	}

	if messageStartIndex < len(lines) {
		commit.Message = string(bytes.Join(lines[messageStartIndex:], []byte{utils.NEWLINE}))
	}

	return commit, nil
}

func parseUser(line string) (User, error) {

	emailStart := bytes.IndexByte([]byte(line), '<')
	emailEnd := bytes.IndexByte([]byte(line), '>')

	if emailStart == -1 || emailEnd == -1 {
		return User{}, fmt.Errorf("invalid user format: %s", line)
	}

	name := string(bytes.TrimSpace([]byte(line[:emailStart])))
	email := line[emailStart+1 : emailEnd]

	remaining := string(bytes.TrimSpace([]byte(line[emailEnd+1:])))

	var timestamp int64
	var timezone string

	_, err := fmt.Sscanf(remaining, "%d %s", &timestamp, &timezone)
	if err != nil {
		return User{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	user := User{
		Name:      name,
		Email:     email,
		Timestamp: time.Unix(timestamp, 0),
	}

	return user, nil
}

func getMainRef(repoPath string) (string, error) {
	headPath := filepath.Join(repoPath, utils.RootShipDir, utils.MainHeadPath)

	headBytes, err := os.ReadFile(headPath)
	if err != nil {
		return "", err
	}

	ref := strings.TrimSpace(string(headBytes))
	return ref, nil
}

func loadCommitByHash(repoBasePath, hash string) (*Commit, error) {
	hash = strings.TrimSpace(hash)
	if hash == "" {
		return nil, nil
	}
	if len(hash) < 3 {
		return nil, fmt.Errorf("invalid commit hash: %s", hash)
	}

	folder := hash[:2]
	file := hash[2:]

	hashPath := filepath.Join(repoBasePath, ".ship", "objects", folder, file)
	data, err := os.ReadFile(hashPath)
	if err != nil {
		return nil, err
	}

	object, err := inflateGitObject(data)
	if err != nil {
		return nil, err
	}

	commit, err := parseCommit(object)
	if err != nil {
		return nil, err
	}

	commit.Hash = hash
	commit.repoPath = repoBasePath
	return commit, nil
}

func LoadCommits(path string) ([]*Commit, error) {

	repoBasePath, err := utils.ShipHasBeenInitRecursive(path)
	if err != nil {
		return nil, err
	}
	if repoBasePath == "" {
		return nil, fmt.Errorf("not a ship repository (or any of the parent directories)")
	}

	head, err := ResolveHead(repoBasePath)
	if err != nil {
		return nil, err
	}

	hash := strings.TrimSpace(head.Hash)
	if hash == "" {
		return []*Commit{}, nil
	}

	commits := make([]*Commit, 0)
	stack := []string{hash}
	visited := make(map[string]bool)

	for len(stack) > 0 {
		currentHash := strings.TrimSpace(stack[len(stack)-1])
		stack = stack[:len(stack)-1]

		if currentHash == "" || visited[currentHash] {
			continue
		}

		commit, err := loadCommitByHash(repoBasePath, currentHash)
		if err != nil {
			return nil, err
		}
		if commit == nil {
			continue
		}

		visited[currentHash] = true
		commits = append(commits, commit)

		for i := len(commit.ParentHashes) - 1; i >= 0; i-- {
			parentHash := strings.TrimSpace(commit.ParentHashes[i])
			if parentHash == "" || visited[parentHash] {
				continue
			}
			stack = append(stack, parentHash)
		}
	}

	return commits, nil
}
