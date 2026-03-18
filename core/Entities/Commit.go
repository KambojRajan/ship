package entities

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/KambojRajan/ship/core/common"
	"github.com/KambojRajan/ship/core/utils"
)

type Commit struct {
	TreeHash     string
	ParentHashes []string
	Author       User
	Committer    User
	Message      string
	repoPath     string
}

func NewCommit(treeHash string, parentHashes []string, author User, committer User, message string, repoPath string) *Commit {
	log.Printf("[NewCommit] Creating new commit with treeHash: %s", treeHash)
	log.Printf("[NewCommit] Parent hashes count: %d", len(parentHashes))
	log.Printf("[NewCommit] Author: %s <%s>", author.Name, author.Email)
	log.Printf("[NewCommit] Committer: %s <%s>", committer.Name, committer.Email)
	log.Printf("[NewCommit] Message: %s", message)

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
	log.Println("[Commit] Starting commit operation")
	log.Printf("[Commit] Tree hash: %s", c.TreeHash)

	// Change to repo directory to write objects with correct paths
	oldCwd, err := os.Getwd()
	if err != nil {
		log.Printf("[Commit] ERROR: Failed to get current directory: %v", err)
		return "", err
	}

	if c.repoPath != "" {
		if err := os.Chdir(c.repoPath); err != nil {
			log.Printf("[Commit] ERROR: Failed to change directory to %s: %v", c.repoPath, err)
			return "", err
		}
		defer func() {
			if err := os.Chdir(oldCwd); err != nil {
				log.Printf("[Commit] ERROR: Failed to restore directory: %v", err)
			}
		}()
	}

	content := c.serialize()
	log.Printf("[Commit] Serialized content size: %d bytes", len(content))

	hash, err := utils.HashObject(content, common.COMMIT, true)
	if err != nil {
		log.Printf("[Commit] ERROR: Failed to hash object: %v", err)
		return "", err
	}

	log.Printf("[Commit] Successfully created commit with hash: %s", hash)
	return hash, nil
}

func (c *Commit) serialize() []byte {
	log.Println("[serialize] Starting commit serialization")
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("%s ", utils.TREE))
	buffer.WriteString(c.TreeHash)
	buffer.WriteByte(utils.NEWLINE)
	log.Printf("[serialize] Added tree hash: %s", c.TreeHash)

	for i, parentHash := range c.ParentHashes {
		buffer.WriteString(fmt.Sprintf("%s ", utils.PARENT))
		buffer.WriteString(parentHash)
		buffer.WriteByte(utils.NEWLINE)
		log.Printf("[serialize] Added parent hash #%d: %s", i+1, parentHash)
	}

	buffer.WriteString(fmt.Sprintf("%s ", utils.AUTHOR))
	buffer.WriteString(c.Author.String())
	buffer.WriteByte(utils.NEWLINE)
	log.Printf("[serialize] Added author: %s", c.Author.String())

	buffer.WriteString(fmt.Sprintf("%s ", utils.COMMITTER))
	buffer.WriteString(c.Committer.String())
	buffer.WriteByte(utils.NEWLINE)
	log.Printf("[serialize] Added committer: %s", c.Committer.String())

	buffer.WriteByte(utils.NEWLINE)

	buffer.WriteString(c.Message)
	log.Printf("[serialize] Added message: %s", c.Message)
	log.Printf("[serialize] Total serialized size: %d bytes", buffer.Len())

	return buffer.Bytes()
}

func (c *Commit) CommitTree(treeHash string, parentHashes []string, author User, message string) (string, error) {
	log.Println("[CommitTree] Starting commit tree operation")
	log.Printf("[CommitTree] Tree hash: %s", treeHash)
	log.Printf("[CommitTree] Parent hashes count: %d", len(parentHashes))
	log.Printf("[CommitTree] Author: %s <%s>", author.Name, author.Email)
	log.Printf("[CommitTree] Message: %s", message)

	committer := c.Committer
	if committer.Timestamp.IsZero() {
		log.Println("[CommitTree] Committer timestamp is zero, using author timestamp")
		committer.Timestamp = author.Timestamp
	}

	commit := NewCommit(treeHash, parentHashes, author, committer, message, c.repoPath)
	hash, err := commit.Commit()
	if err != nil {
		log.Printf("[CommitTree] ERROR: Failed to commit: %v", err)
		return "", err
	}

	log.Printf("[CommitTree] Successfully created commit tree with hash: %s", hash)
	return hash, nil
}

func inflateGitObject(data []byte) ([]byte, error) {
	log.Printf("[inflateGitObject] Starting inflation of %d bytes", len(data))

	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		log.Printf("[inflateGitObject] ERROR: Failed to create zlib reader: %v", err)
		return nil, err
	}
	defer func(r io.ReadCloser) {
		err := r.Close()
		if err != nil {
			log.Printf("[inflateGitObject] ERROR: Failed to close zlib reader: %v", err)
		}
	}(r)

	result, err := io.ReadAll(r)
	if err != nil {
		log.Printf("[inflateGitObject] ERROR: Failed to read inflated data: %v", err)
		return nil, err
	}

	log.Printf("[inflateGitObject] Successfully inflated to %d bytes", len(result))
	return result, nil
}

func parseCommit(data []byte) (*Commit, error) {
	log.Printf("[parseCommit] Starting to parse commit data (%d bytes)", len(data))

	nullIndex := bytes.IndexByte(data, 0)
	if nullIndex == -1 {
		log.Println("[parseCommit] ERROR: No null byte found in commit object")
		return nil, fmt.Errorf("invalid commit object: no null byte found")
	}
	log.Printf("[parseCommit] Found null byte at index %d", nullIndex)

	content := data[nullIndex+1:]
	log.Printf("[parseCommit] Content size after header: %d bytes", len(content))

	lines := bytes.Split(content, []byte{utils.NEWLINE})
	log.Printf("[parseCommit] Split content into %d lines", len(lines))

	commit := &Commit{
		ParentHashes: []string{},
	}

	var messageStartIndex int
	for i, line := range lines {
		if len(line) == 0 {
			log.Printf("[parseCommit] Found empty line at index %d (message starts next)", i)
			messageStartIndex = i + 1
			break
		}

		parts := bytes.SplitN(line, []byte{' '}, 2)
		if len(parts) < 2 {
			log.Printf("[parseCommit] Skipping line %d (insufficient parts): %s", i, string(line))
			continue
		}

		key := string(parts[0])
		value := string(parts[1])
		log.Printf("[parseCommit] Line %d: key=%s, value=%s", i, key, value)

		switch key {
		case utils.TREE:
			commit.TreeHash = value
			log.Printf("[parseCommit] Set tree hash: %s", value)
		case utils.PARENT:
			commit.ParentHashes = append(commit.ParentHashes, value)
			log.Printf("[parseCommit] Added parent hash: %s (total: %d)", value, len(commit.ParentHashes))
		case utils.AUTHOR:
			user, err := parseUser(value)
			if err != nil {
				log.Printf("[parseCommit] ERROR: Failed to parse author: %v", err)
				return nil, fmt.Errorf("failed to parse author: %w", err)
			}
			commit.Author = user
			log.Printf("[parseCommit] Set author: %s <%s>", user.Name, user.Email)
		case utils.COMMITTER:
			user, err := parseUser(value)
			if err != nil {
				log.Printf("[parseCommit] ERROR: Failed to parse committer: %v", err)
				return nil, fmt.Errorf("failed to parse committer: %w", err)
			}
			commit.Committer = user
			log.Printf("[parseCommit] Set committer: %s <%s>", user.Name, user.Email)
		}
	}

	if messageStartIndex < len(lines) {
		commit.Message = string(bytes.Join(lines[messageStartIndex:], []byte{utils.NEWLINE}))
		log.Printf("[parseCommit] Parsed message (%d chars): %s", len(commit.Message), commit.Message)
	} else {
		log.Println("[parseCommit] No message found in commit")
	}

	log.Println("[parseCommit] Successfully parsed commit")
	return commit, nil
}

func parseUser(line string) (User, error) {
	log.Printf("[parseUser] Parsing user from line: %s", line)

	emailStart := bytes.IndexByte([]byte(line), '<')
	emailEnd := bytes.IndexByte([]byte(line), '>')

	if emailStart == -1 || emailEnd == -1 {
		log.Printf("[parseUser] ERROR: Invalid user format, missing email delimiters")
		return User{}, fmt.Errorf("invalid user format: %s", line)
	}

	name := string(bytes.TrimSpace([]byte(line[:emailStart])))
	email := line[emailStart+1 : emailEnd]
	log.Printf("[parseUser] Parsed name: %s, email: %s", name, email)

	remaining := string(bytes.TrimSpace([]byte(line[emailEnd+1:])))
	log.Printf("[parseUser] Remaining timestamp info: %s", remaining)

	var timestamp int64
	var timezone string

	_, err := fmt.Sscanf(remaining, "%d %s", &timestamp, &timezone)
	if err != nil {
		log.Printf("[parseUser] ERROR: Failed to parse timestamp: %v", err)
		return User{}, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	log.Printf("[parseUser] Parsed timestamp: %d, timezone: %s", timestamp, timezone)

	user := User{
		Name:      name,
		Email:     email,
		Timestamp: time.Unix(timestamp, 0),
	}

	log.Printf("[parseUser] Successfully parsed user: %s <%s> at %v", user.Name, user.Email, user.Timestamp)
	return user, nil
}

func getMainRef(repoPath string) (string, error) {
	log.Printf("[getMainRef] Getting main ref from repo path: %s", repoPath)

	headPath := filepath.Join(repoPath, utils.RootShipDir, utils.MainHeadPath)
	log.Printf("[getMainRef] Reading HEAD from: %s", headPath)

	headBytes, err := os.ReadFile(headPath)
	if err != nil {
		log.Printf("[getMainRef] ERROR: Failed to read HEAD file: %v", err)
		return "", err
	}

	ref := string(headBytes)
	log.Printf("[getMainRef] Successfully retrieved ref: %s", ref)
	return ref, nil
}

func LoadCommits(path string) ([]*Commit, error) {
	log.Printf("[LoadCommits] Starting to load commits from path: %s", path)

	hash, err := getMainRef(path)
	if err != nil {
		log.Printf("[LoadCommits] ERROR: Failed to get main ref: %v", err)
		return nil, err
	}
	log.Printf("[LoadCommits] Got commit hash: %s", hash)

	repoBasePath, err := utils.ShipHasBeenInitRecursive(path)
	if err != nil {
		log.Printf("[LoadCommits] ERROR: Failed to find repo base path: %v", err)
		return nil, err
	}
	log.Printf("[LoadCommits] Repo base path: %s", repoBasePath)

	folder := hash[:2]
	file := hash[2:]
	log.Printf("[LoadCommits] Object folder: %s, file: %s", folder, file)

	hashPath := filepath.Join(repoBasePath, ".ship", "objects", folder, file)
	log.Printf("[LoadCommits] Reading object from: %s", hashPath)

	data, err := os.ReadFile(hashPath)
	if err != nil {
		log.Printf("[LoadCommits] ERROR: Failed to read object file: %v", err)
		return nil, err
	}
	log.Printf("[LoadCommits] Read %d bytes from object file", len(data))

	object, err := inflateGitObject(data)
	if err != nil {
		log.Printf("[LoadCommits] ERROR: Failed to inflate object: %v", err)
		return nil, err
	}
	log.Printf("[LoadCommits] Inflated object to %d bytes", len(object))

	commit, err := parseCommit(object)
	if err != nil {
		log.Printf("[LoadCommits] ERROR: Failed to parse commit: %v", err)
		return nil, err
	}
	log.Printf("[LoadCommits] Successfully parsed commit")

	commits := []*Commit{commit}
	log.Printf("[LoadCommits] Returning %d commit(s)", len(commits))
	return commits, nil
}
