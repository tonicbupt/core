package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strings"

	engineapi "github.com/docker/engine-api/client"
	"gitlab.ricebook.net/platform/core/types"
	"golang.org/x/net/context"
)

const (
	letters       = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	shortenLength = 7
	CpuPeriodBase = 100000
)

func RandomString(n int) string {
	r := make([]byte, n)
	for i := 0; i < n; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		// 没那么惨吧
		if err != nil {
			continue
		}
		r[i] = letters[n.Int64()]
	}
	return string(r)
}

func TruncateID(id string) string {
	if len(id) > shortenLength {
		return id[:shortenLength]
	}
	return id
}

func Tail(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func GetGitRepoName(url string) (string, error) {
	if !strings.HasPrefix(url, "git@") || !strings.HasSuffix(url, ".git") {
		return "", fmt.Errorf("Bad git url format %q", url)
	}

	x := strings.SplitN(url, ":", 2)
	if len(x) != 2 {
		return "", fmt.Errorf("Bad git url format %q", url)
	}

	y := strings.SplitN(x[1], "/", 2)
	if len(y) != 2 {
		return "", fmt.Errorf("Bad git url format %q", url)
	}
	return strings.TrimSuffix(y[1], ".git"), nil
}

func GetVersion(image string) string {
	if !strings.Contains(image, ":") {
		return "unknown"
	}

	parts := strings.Split(image, ":")
	if len(parts) != 2 {
		return "unknown"
	}

	return parts[1]
}

// Bind a docker engine client to context
func ToDockerContext(client *engineapi.Client) context.Context {
	return context.WithValue(context.Background(), "engine", client)
}

// Get a docker engine client from a context
func FromDockerContext(ctx context.Context) (*engineapi.Client, bool) {
	client, ok := ctx.Value("engine").(*engineapi.Client)
	return client, ok
}

func SaveFile(content, path string, mode os.FileMode) error {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}

type NodeInfo struct {
	Name    string
	CorePer int
}

type ByCoreNum []NodeInfo

func (a ByCoreNum) Len() int           { return len(a) }
func (a ByCoreNum) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByCoreNum) Less(i, j int) bool { return a[i].CorePer < a[j].CorePer }

func AllocContainerPlan(nodeInfo ByCoreNum, quota int, count int) (map[string]int, error) {
	result := make(map[string]int)
	N := nodeInfo.Len()
	ave := 0
	remain := 0
	flag := -1

	for i := 0; i < N; i++ {
		if nodeInfo[i].CorePer >= quota {
			ave = count / (N - i)
			remain = count % (N - i)
			flag = i
			break
		}
	}

	if flag == -1 {
		return result, fmt.Errorf("Cannot alloc a plan, nodeNum: %d, node1 CorePer: %d, quota: %d ", N, nodeInfo[0].CorePer, quota)
	}

	for i := flag; i < N; i++ {
		result[nodeInfo[i].Name] = ave
	}
	resLen := int64(len(result))
	if remain > 0 {
	r:
		for {
			// 考虑一种情况：不断申请一个 quota 相同的容器
			// 按原来的算法，这个容器会堆积在同一台机上面
			// 加入随机化的选择可以避免这种情况
			step, _ := rand.Int(rand.Reader, big.NewInt(resLen))
			node := nodeInfo[flag+int(step.Int64())].Name
			result[node] += 1
			remain--
			if remain <= 0 {
				break r
			}
		}
	}
	return result, nil
}

func GetNodesInfo(cpumap map[string]types.CPUMap) ByCoreNum {
	result := ByCoreNum{}
	for node, cpu := range cpumap {
		result = append(result, NodeInfo{node, len(cpu) * CpuPeriodBase})
	}
	sort.Sort(result)
	return result
}

// copied from https://gist.github.com/jmervine/d88c75329f98e09f5c87
func safeSplit(s string) []string {
	split := strings.Split(s, " ")

	var result []string
	var inquote string
	var block string
	for _, i := range split {
		if inquote == "" {
			if strings.HasPrefix(i, "'") || strings.HasPrefix(i, "\"") {
				inquote = string(i[0])
				block = strings.TrimPrefix(i, inquote) + " "
			} else {
				result = append(result, i)
			}
		} else {
			if !strings.HasSuffix(i, inquote) {
				block += i + " "
			} else {
				block += strings.TrimSuffix(i, inquote)
				inquote = ""
				result = append(result, block)
				block = ""
			}
		}
	}

	return result
}

func MakeCommandLineArgs(s string) []string {
	r := []string{}
	for _, part := range safeSplit(s) {
		if len(part) == 0 {
			continue
		}
		r = append(r, part)
	}
	return r
}
