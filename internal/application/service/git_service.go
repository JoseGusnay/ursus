package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type GitService struct {
	dbDir string
}

func NewGitService() *GitService {
	home, _ := os.UserHomeDir()
	return &GitService{
		dbDir: filepath.Join(home, ".ursus"),
	}
}

func (s *GitService) run(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = s.dbDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), fmt.Errorf("git %v failed: %w", args, err)
	}
	return string(out), nil
}

func (s *GitService) Init(remote string) (string, error) {
	if _, err := os.Stat(filepath.Join(s.dbDir, ".git")); os.IsNotExist(err) {
		if out, err := s.run("init"); err != nil {
			return out, err
		}
	}
	
	// Create .gitignore if not exists to ensure only the DB is tracked
	gitignorePath := filepath.Join(s.dbDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		os.WriteFile(gitignorePath, []byte("*.db-journal\n*.db-wal\n*.db-shm\n"), 0644)
	}

	return s.run("remote", "add", "origin", remote)
}

func (s *GitService) Push() (string, error) {
	if out, err := s.run("add", "."); err != nil {
		return out, err
	}
	// Note: Commit might fail if there are no changes, but we ignore it for simple sync
	s.run("commit", "-m", "chore: update ursus memories")
	
	return s.run("push", "origin", "main")
}

func (s *GitService) Pull() (string, error) {
	return s.run("pull", "origin", "main")
}

func (s *GitService) IsConfigured() bool {
	_, err := os.Stat(filepath.Join(s.dbDir, ".git"))
	return err == nil
}
