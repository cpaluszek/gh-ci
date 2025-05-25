package github

import (
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
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

type GitHubJobsResponse struct {
	Jobs []Job `json:"jobs"`
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

func (c *Client) GetLogs(user, repo, runID, attempt string, jobName string) ([]Steplog, error) {
	cacheKey := fmt.Sprintf("logs:%s:%s:%s:%s:%s", user, repo, runID, attempt, jobName)
	cache, err := cache.LoadCache()
	if err != nil {
		return nil, fmt.Errorf("failed to load cache: %v", err)
	}

	if cachedData, found := cache.Get(cacheKey); found {
		if steplogs, ok := cachedData.([]Steplog); ok {
			return steplogs, nil
		}
	}

	zipCacheKey := fmt.Sprintf("logs:%s:%s:%s:%s", user, repo, runID, attempt)
	zipPath, foundFile := cache.GetFileCache(zipCacheKey)
	var zipData []byte

	if !foundFile {
		logsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%s/attempts/%s/logs", user, repo, runID, attempt)

		resp, err := c.Client.Request(http.MethodGet, logsURL, nil)
		if err != nil {
			log.Printf("failed to fetch logs: %v", err)
			return nil, err
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, 1024)) // Read up to 1KB for error logging
			log.Printf("failed to download logs: status code %d from %s. Response: %s", resp.StatusCode, resp.Request.URL, string(bodyBytes))
			return nil, fmt.Errorf("failed to download logs: status code %d from %s", resp.StatusCode, resp.Request.URL)
		}

		zipData, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("failed to read log response body: %v", err)
			return nil, fmt.Errorf("failed to read log response body: %w", err)
		}
		log.Printf("Fetched logs successfully, size: %d bytes\n", len(zipData))

		// TODO: use .config file log TTL
		zipPath, err = cache.SetFileCache(zipCacheKey, zipData, time.Hour)
		if err != nil {
			return nil, fmt.Errorf("failed to cache zip file: %v", err)
		}
		defer func() {
			removeErr := os.Remove(zipPath)
			if removeErr == nil {
				err = removeErr
			}
		}()
	} else {
		log.Printf("Using cached zip file: %s\n", zipPath)
		zipData, err = os.ReadFile(zipPath)
		if err != nil {
			return nil, err
		}
	}

	// fetch metadata for steps
	stepMeta, err := c.FetchGitHubJobSteps(user, repo, runID, jobName)
	if err != nil {
		return nil, err
	}

	steplogs, err := ParseZipLogs(zipData, stepMeta, jobName)
	if err != nil {
		return nil, err
	}

	if err := cache.Set(cacheKey, steplogs, 30*time.Minute); err != nil {
		fmt.Printf("Warning: failed to cache parsed logs: %v\n", err)
	}

	return steplogs, nil
}

func (c *Client) FetchGitHubJobSteps(owner, repo, runID, jobname string) (map[int]Step, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/actions/runs/%s/jobs", owner, repo, runID)

	var jobsResponse GitHubJobsResponse

	err := c.Client.Get(apiURL, &jobsResponse)
	if err != nil {
		return nil, err
	}

	stepMap := make(map[int]Step)

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

func ParseZipLogs(zipData []byte, stepMeta map[int]Step, jobName string) ([]Steplog, error) {
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

		err = rc.Close()
		if err != nil {
			continue
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
