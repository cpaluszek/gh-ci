package github

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cpaluszek/gh-ci/cache"
)

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

type Steplog struct {
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	Status    string     `json:"status"`
	Duration  string     `json:"duration"`
	Logs      []LogEntry `json:"logs"`
	Collapsed bool       `json:"collapsed"`
}

type GitHubRunInfo struct {
	User  string `json:"user"`
	Repo  string `json:"repo"`
	RunID string `json:"run_id"`
	JobID string `json:"job_id,omitempty"`
}

type GitHubStep struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	Conclusion  string    `json:"conclusion"`
	Number      int       `json:"number"`
	StartedAt   time.Time `json:"started_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type GitHubJob struct {
	Name  string       `json:"name"`
	Steps []GitHubStep `json:"steps"`
}

type GitHubJobsResponse struct {
	Jobs []GitHubJob `json:"jobs"`
}

var (
	fileNameRegex  = regexp.MustCompile(`(\d+)_([\w\s\-\.]+)\.txt$`)
	timestampRegex = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)`)
)

func ParseGitHubURL(rawURL string) (*GitHubRunInfo, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}
	if u.Host != "github.com" {
		return nil, fmt.Errorf("URL must be from github.com")
	}
	pathRegex := regexp.MustCompile(`^/([^/]+)/([^/]+)/actions/runs/(\d+)(?:/job/(\d+))?`)
	matches := pathRegex.FindStringSubmatch(u.Path)
	if len(matches) < 4 {
		return nil, fmt.Errorf("invalid GitHub Actions URL format")
	}
	info := &GitHubRunInfo{
		User:  matches[1],
		Repo:  matches[2],
		RunID: matches[3],
	}
	if len(matches) > 4 && matches[4] != "" {
		info.JobID = matches[4]
	}
	return info, nil
}

func GetLogs(user, repo, runID, attempt string, jobName string) ([]Steplog, error) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("GITHUB_TOKEN not set")
	}

	cacheKey := fmt.Sprintf("logs:%s:%s:%s:%s:%s", user, repo, runID, attempt, jobName)
	c, err := cache.LoadCache()
	if err != nil {
		return nil, fmt.Errorf("failed to load cache: %v", err)
	}

	if cachedData, found := c.Get(cacheKey); found {
		if steplogs, ok := cachedData.([]Steplog); ok {
			return steplogs, nil
		}
	}

	zipCacheKey := fmt.Sprintf("logs:%s:%s:%s:%s", user, repo, runID, attempt)
	zipPath, foundFile := c.GetFileCache(zipCacheKey)
	var zipData []byte

	if !foundFile {
		logsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%s/attempts/%s/logs", user, repo, runID, attempt)
		req, err := http.NewRequest("GET", logsURL, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Accept", "application/vnd.github.v3+json")
		req.Header.Set("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
		}
		zipData, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		zipPath, err = c.SetFileCache(zipCacheKey, zipData, time.Hour)
		if err != nil {
			return nil, fmt.Errorf("failed to cache zip file: %v", err)
		}
		defer os.Remove(zipPath)
	} else {
		zipData, err = os.ReadFile(zipPath)
		if err != nil {
			return nil, err
		}
	}

	// fetch metadata for steps
	stepMeta, err := FetchGitHubJobSteps(user, repo, runID, token, jobName)
	if err != nil {
		return nil, err
	}

	steplogs, err := ParseZipLogs(zipData, stepMeta, jobName)
	if err != nil {
		return nil, err
	}

	if err := c.Set(cacheKey, steplogs, 30*time.Minute); err != nil {
		fmt.Printf("Warning: failed to cache parsed logs: %v\n", err)
	}

	return steplogs, nil
}

func FetchGitHubJobSteps(owner, repo, runID, token, jobname string) (map[int]GitHubStep, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%s/jobs", owner, repo, runID)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error: %s", body)
	}

	var jobsResponse GitHubJobsResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobsResponse); err != nil {
		return nil, err
	}

	stepMap := make(map[int]GitHubStep)

	for _, job := range jobsResponse.Jobs {
		if job.Name == jobname {
			for _, step := range job.Steps {
				stepMap[step.Number] = step
			}
			break
		}
	}

	return stepMap, nil
}

func ParseZipLogs(zipData []byte, stepMeta map[int]GitHubStep, jobName string) ([]Steplog, error) {
	reader, err := zip.NewReader(bytes.NewReader(zipData), int64(len(zipData)))
	if err != nil {
		return nil, err
	}

	var steplogs []Steplog

	for _, file := range reader.File {
		if !strings.Contains(file.Name, jobName) || !strings.HasSuffix(file.Name, ".txt") {
			continue
		}

		match := fileNameRegex.FindStringSubmatch(file.Name)
		if len(match) < 3 {
			continue
		}

		stepNumber, err := strconv.Atoi(match[1])
		if err != nil {
			continue
		}
		stepTitle := strings.ReplaceAll(match[2], "_", " ")

		rc, err := file.Open()
		if err != nil {
			continue
		}
		defer rc.Close()
		// Delete the file after reading
		defer os.Remove(file.Name)

		scanner := bufio.NewScanner(rc)
		var logs []LogEntry

		for scanner.Scan() {
			line := scanner.Text()
			timestamp := ""
			level := "info"

			if tsMatch := timestampRegex.FindStringSubmatch(line); len(tsMatch) > 1 {
				timestamp = tsMatch[1]
				line = strings.TrimPrefix(line, timestamp+" ")
			}

			lower := strings.ToLower(line)
			switch {
			case strings.Contains(lower, "error"):
				level = "error"
			case strings.Contains(lower, "warn"):
				level = "warning"
			}

			logs = append(logs, LogEntry{
				Timestamp: timestamp,
				Level:     level,
				Message:   line,
			})
		}

		meta := stepMeta[stepNumber]
		duration := ""
		if !meta.StartedAt.IsZero() && !meta.CompletedAt.IsZero() {
			d := meta.CompletedAt.Sub(meta.StartedAt)
			duration = d.String()
		}

		steplogs = append(steplogs, Steplog{
			Number:    stepNumber,
			Title:     stepTitle,
			Status:    meta.Conclusion,
			Duration:  duration,
			Logs:      logs,
			Collapsed: false,
		})
	}

	// Add Metadata for steps without logs
	for _, step := range stepMeta {
		exists := false
		for _, st := range steplogs {
			if st.Number == step.Number {
				exists = true
				break
			}
		}
		if !exists {
			duration := ""
			if !step.StartedAt.IsZero() && !step.CompletedAt.IsZero() {
				d := step.CompletedAt.Sub(step.StartedAt)
				duration = d.String()
			}

			steplogs = append(steplogs, Steplog{
				Number:    step.Number,
				Title:     step.Name,
				Status:    step.Conclusion,
				Duration:  duration,
				Logs:      []LogEntry{},
				Collapsed: false,
			})
		}
	}
	ret := []Steplog{}
	for _, steplog := range steplogs {
		if steplog.Status != "" {
			ret = append(ret, steplog)
		}
	}
	steplogs = ret

	// Sort by step number
	sort.Slice(steplogs, func(i, j int) bool {
		if steplogs[i].Number == steplogs[j].Number {
			return steplogs[i].Title < steplogs[j].Title
		}
		return steplogs[i].Number < steplogs[j].Number
	})

	return steplogs, nil
}
