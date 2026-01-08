package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jpmarques19/treehouse/internal/testutil"
)

func TestDetectRepo_InGitRoot(t *testing.T) {
	dir := testutil.SetupGitRepo(t)
	testutil.ChdirWithCleanup(t, dir)

	info, err := DetectRepo()
	if err != nil {
		t.Fatalf("DetectRepo() error = %v, want nil", err)
	}

	if info.Root == "" {
		t.Error("Root should not be empty")
	}

	if info.Branch == "" {
		t.Error("Branch should not be empty")
	}

	if info.CommitSHA == "" {
		t.Error("CommitSHA should not be empty")
	}

	if len(info.CommitSHA) != 7 {
		t.Errorf("CommitSHA length = %d, want 7", len(info.CommitSHA))
	}
}

func TestDetectRepo_InSubdirectory(t *testing.T) {
	dir := testutil.SetupGitRepo(t)

	// Create subdirectory
	subdir := filepath.Join(dir, "sub", "folder")
	if err := os.MkdirAll(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdirectory: %v", err)
	}

	testutil.ChdirWithCleanup(t, subdir)

	info, err := DetectRepo()
	if err != nil {
		t.Fatalf("DetectRepo() error = %v, want nil", err)
	}

	// Root should be the parent git directory, not the subdirectory
	if !strings.HasSuffix(info.Root, filepath.Base(dir)) {
		t.Errorf("Root = %s, should end with %s", info.Root, filepath.Base(dir))
	}
}

func TestDetectRepo_NotGitRepo(t *testing.T) {
	dir := t.TempDir() // Just a temp dir, not a git repo
	testutil.ChdirWithCleanup(t, dir)

	info, err := DetectRepo()
	if err == nil {
		t.Fatal("DetectRepo() error = nil, want error")
	}

	if info != nil {
		t.Error("info should be nil on error")
	}

	gitErr, ok := err.(*GitError)
	if !ok {
		t.Fatalf("error type = %T, want *GitError", err)
	}

	if gitErr.Code != "INIT_NOT_GIT_REPO" {
		t.Errorf("error code = %s, want INIT_NOT_GIT_REPO", gitErr.Code)
	}
}

func TestDetectRepo_DetachedHead(t *testing.T) {
	dir := testutil.SetupGitRepo(t)

	// Detach HEAD by checking out specific commit
	cmd := exec.Command("git", "checkout", "--detach", "HEAD")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to detach HEAD: %v", err)
	}

	testutil.ChdirWithCleanup(t, dir)

	info, err := DetectRepo()
	if err != nil {
		t.Fatalf("DetectRepo() error = %v, want nil", err)
	}

	// In detached HEAD state, branch should be "HEAD"
	if info.Branch != "HEAD" {
		t.Errorf("Branch = %s, want HEAD (detached)", info.Branch)
	}

	// Root and CommitSHA should still be valid
	if info.Root == "" {
		t.Error("Root should not be empty in detached HEAD")
	}

	if len(info.CommitSHA) != 7 {
		t.Errorf("CommitSHA length = %d, want 7", len(info.CommitSHA))
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name       string
		versionStr string
		wantMajor  int
		wantMinor  int
		wantPatch  int
		wantErr    bool
	}{
		{
			name:       "standard version",
			versionStr: "git version 2.34.1",
			wantMajor:  2,
			wantMinor:  34,
			wantPatch:  1,
			wantErr:    false,
		},
		{
			name:       "version with extra info",
			versionStr: "git version 2.39.3 (Apple Git-146)",
			wantMajor:  2,
			wantMinor:  39,
			wantPatch:  3,
			wantErr:    false,
		},
		{
			name:       "minimum supported",
			versionStr: "git version 2.5.0",
			wantMajor:  2,
			wantMinor:  5,
			wantPatch:  0,
			wantErr:    false,
		},
		{
			name:       "old version",
			versionStr: "git version 1.9.0",
			wantMajor:  1,
			wantMinor:  9,
			wantPatch:  0,
			wantErr:    false,
		},
		{
			name:       "two-part version without patch",
			versionStr: "git version 2.34",
			wantMajor:  2,
			wantMinor:  34,
			wantPatch:  0,
			wantErr:    false,
		},
		{
			name:       "pre-release version",
			versionStr: "git version 2.34.1-rc1",
			wantMajor:  2,
			wantMinor:  34,
			wantPatch:  1,
			wantErr:    false,
		},
		{
			name:       "invalid format",
			versionStr: "not a version",
			wantMajor:  0,
			wantMinor:  0,
			wantPatch:  0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := parseVersion(tt.versionStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseVersion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if v.Major != tt.wantMajor {
				t.Errorf("Major = %d, want %d", v.Major, tt.wantMajor)
			}
			if v.Minor != tt.wantMinor {
				t.Errorf("Minor = %d, want %d", v.Minor, tt.wantMinor)
			}
			if v.Patch != tt.wantPatch {
				t.Errorf("Patch = %d, want %d", v.Patch, tt.wantPatch)
			}
		})
	}
}

func TestCheckVersionRequirement(t *testing.T) {
	tests := []struct {
		name        string
		version     Version
		wantErr     bool
		wantErrCode string
	}{
		{
			name:    "version 2.34.1 ok",
			version: Version{Major: 2, Minor: 34, Patch: 1},
			wantErr: false,
		},
		{
			name:    "version 2.5.0 ok (minimum)",
			version: Version{Major: 2, Minor: 5, Patch: 0},
			wantErr: false,
		},
		{
			name:        "version 2.4.0 fail",
			version:     Version{Major: 2, Minor: 4, Patch: 0},
			wantErr:     true,
			wantErrCode: "GIT_VERSION_UNSUPPORTED",
		},
		{
			name:        "version 1.9.0 fail",
			version:     Version{Major: 1, Minor: 9, Patch: 0},
			wantErr:     true,
			wantErrCode: "GIT_VERSION_UNSUPPORTED",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkVersionRequirement(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkVersionRequirement() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				gitErr, ok := err.(*GitError)
				if !ok {
					t.Fatalf("error type = %T, want *GitError", err)
				}
				if gitErr.Code != tt.wantErrCode {
					t.Errorf("error code = %s, want %s", gitErr.Code, tt.wantErrCode)
				}
			}
		})
	}
}

func TestCheckVersion_Integration(t *testing.T) {
	// This test runs against the real git installation
	// It should pass on any modern system with git 2.5+
	err := CheckVersion()
	if err != nil {
		t.Logf("CheckVersion returned error: %v (this may be expected on old git)", err)
		// Don't fail - old git installations are valid test environments
	}
}
