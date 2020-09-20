package agent

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"

	srlndk "github.com/srl-wim/protos"
)

type yangGit struct {
	NetworkInstance struct {
		Value string `json:"value"`
	} `json:"network-instance"`
	Organization struct {
		Value string `json:"value"`
	} `json:"organization"`
	OWner struct {
		Value string `json:"value"`
	} `json:"owner"`
	Repo struct {
		Value string `json:"value"`
	} `json:"repo"`
	File struct {
		Value string `json:"value"`
	} `json:"file"`
	Token struct {
		Value string `json:"value"`
	} `json:"token"`
	Author struct {
		Value string `json:"value"`
	} `json:"author"`
	AuthorEmail struct {
		Value string `json:"value"`
	} `json:"author_email"`
	Branch struct {
		Value string `json:"value"`
	} `json:"branch"`
	Action struct {
		Value string `json:"value"`
	} `json:"action"`
	OperState  string            `json:"oper_state"`
	Statistics yangGitStatistics `json:"statistics"`
}

type yangGitStatistics struct {
	Success struct {
		Value uint64 `json:"value"`
	} `json:"success"`
	Failure struct {
		Value uint64 `json:"value"`
	} `json:"failures"`
}

type cfgTranxEntry struct {
	Op   srlndk.SdkMgrOperation
	Key  *[]string
	Data *string
}

type jsParser struct {
	re1 *regexp.Regexp
	re2 *regexp.Regexp
}

func newJSParser() *jsParser {
	return &jsParser{
		re1: regexp.MustCompile(`\{(.*?)\}`),
		re2: regexp.MustCompile(`\"(.*?)\"`),
	}
}
func (p *jsParser) jsPathToPathKeys(jsPath string, keys []string) (string, []string) {
	if jsPath == "" {
		return "", make([]string, 0)
	}
	submatchall1 := p.re1.FindAllString(jsPath, -1)
	for _, element1 := range submatchall1 {
		submatchall2 := p.re2.FindAllString(element1, -1)
		for _, element2 := range submatchall2 {
			element2 = strings.Replace(element2, "\"", "", -1)
			keys = append(keys, element2)
		}
		jsPath = strings.Replace(jsPath, element1, "", -1)
	}
	return jsPath, keys
}

// HandleConfigEvent function
func (a *Agent) HandleConfigEvent(op srlndk.SdkMgrOperation, key *srlndk.ConfigKey, data *string) {
	// log handle config event
	log.Printf("handleConfigEvent: %v, jspath: %v, keys: %v \n", op, key.GetJsPath(), key.GetKeys())

	p := newJSParser()

	newJsPath, newKey := p.jsPathToPathKeys(key.GetJsPath(), key.GetKeys())

	log.Printf("handleConfigEvent: %v, NewJsPath: %v, NewKeys: %v \n", op, key.GetJsPath(), key.GetKeys())

	// handle end of commit operation
	if key.GetJsPath() != ".commit.end" {
		// Append the information in a tarnsaction map
		a.Config.cfgTranxMap[newJsPath] = append(a.Config.cfgTranxMap[newJsPath], cfgTranxEntry{Op: op, Key: &newKey, Data: data})
		return
	}

	// handle test agent configuration event
	for _, item := range a.Config.cfgTranxMap[".git"] {
		a.HandleGitConfigEvent(item.Op, item.Key, item.Data)
	}

	// Delete all current candidate list. Reinitialize the map
	a.Config.cfgTranxMap = make(map[string][]cfgTranxEntry)
}

// HandleGitConfigEvent function
func (a *Agent) HandleGitConfigEvent(op srlndk.SdkMgrOperation, key *[]string, data *string) {
	log.Printf(".git jsPath %v with operation %v", *key, op)
	if data == nil {
		if op == srlndk.SdkMgrOperation_Delete {
			// Handle delete Configuration
			log.Printf("Handle Delete Config")
			jsPath := ".git"
			a.deleteTelemetry(&jsPath)

		}
		return
	}

	if data != nil {
		// log data received from the yang server
		log.Printf("git data %v", *data)
	}

	var ydata *yangGit
	if err := json.Unmarshal([]byte(*data), &ydata); err != nil {
		log.Fatalf("Can not unmarchal config data: %v error %v", *data, err)
		os.Exit(1)
	}
	log.Printf("YANG Unmarchal: %v", *ydata)

	// handle the create or change yang operation event
	if op != srlndk.SdkMgrOperation_Delete {
		log.Printf("Handle Create or Delete Config")
		// update the global variable yang structures with the information from the configserver
		a.Config.YangConfig = ydata
		log.Printf("YANG config data structure: %v", *a.Config.YangConfig)

		log.Printf("YANG Token: %s \n", a.Config.YangConfig.Token.Value)
		if a.Github.token != nil {
			log.Printf("CONFIG Token: %s \n", *a.Github.token)
		}

		if a.Config.YangConfig.Token.Value != "" && a.Github.token == nil {
			log.Printf("GIT connect")
			a.GitClient()
			a.Github.token = &a.Config.YangConfig.Token.Value
			a.Config.YangConfig.OperState = "OPER_STATE_up"
		}

		log.Printf("Action: %s \n", a.Config.YangConfig.Action.Value)
		if a.Config.YangConfig.Action.Value == "branch" {
			log.Print("Git Branch")

			if err := a.GetRef(&a.Config.YangConfig.Branch.Value); err != nil {
				log.Printf("Error: Unable to get/create the commit reference: %s\n", err)
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}
			if a.Github.Ref == nil {
				log.Printf("Error: No error where returned but the reference is nil")
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}
			a.Config.YangConfig.Statistics.Success.Value++
		}

		if a.Config.YangConfig.Action.Value == "commit" {
			log.Print("Git commit")

			if err := a.GetRef(&a.Config.YangConfig.Branch.Value); err != nil {
				log.Printf("Error Unable to get/create the commit reference: %s\n", err)
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}
			if a.Github.Ref == nil {
				log.Printf("Error: No error where returned but the reference is nil")
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}

			if err := a.GetTree(); err != nil {
				log.Printf("Error Unable to create the tree based on the provided files: %s\n", err)
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}

			if err := a.PushCommit(a.Github.Ref, a.Github.Tree); err != nil {
				log.Printf("Error Unable to create the commit: %s\n", err)
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}
			a.Config.YangConfig.Statistics.Success.Value++
		}

		if a.Config.YangConfig.Action.Value == "pr" {
			log.Print("Git pull request")

			if err := a.GetRef(&a.Config.YangConfig.Branch.Value); err != nil {
				log.Printf("Error Unable to get/create the commit reference: %s\n", err)
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}
			if a.Github.Ref == nil {
				log.Printf("Error: No error where returned but the reference is nil")
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}

			if err := a.GetTree(); err != nil {
				log.Printf("Error: Unable to create the tree based on the provided files: %s\n", err)
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}

			if err := a.CreatePR(&a.Config.YangConfig.Branch.Value); err != nil {
				log.Printf("Error while creating the pull request: %s", err)
				a.Config.YangConfig.Statistics.Failure.Value++
				a.updateConfigTelemetry()
				return
			}
			a.Config.YangConfig.Statistics.Success.Value++
		}
		a.updateConfigTelemetry()
	}
}

func (a *Agent) updateConfigTelemetry() {
	jsPath := ".git"
		jsData, err := json.Marshal(a.Config.YangConfig)
		if err != nil {
			log.Fatalf("Can not marshal config data:%v error %s", *a.Config.YangConfig, err)
		}
		jsonStr := string(jsData)

		log.Printf("Sending telemetry update js_path: %s js_data: %s", jsPath, jsonStr)
		a.updateTelemetry(&jsPath, &jsonStr)
}
