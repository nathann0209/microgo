# High level Overview:
Here's how we can structure it:

`main.go` - The entry point of the application
`git_handler.go` - Handles Git operations
`data_processor.go` - Processes YAML files and caches data
`config.go` - Contains configuration constants
# Code
`main.go`
```
package main

import (
    "fmt"
    "time"
)

func main() {
    InitConfig()
    InitGitHandler()
    InitDataProcessor()

    // Start the periodic update in a goroutine
    go periodicUpdate()

    // Your web server setup and handlers would go here
    // When serving a request, you'd call GetDataForUser()

    // For demonstration purposes, we'll just run indefinitely
    for {
        time.Sleep(time.Hour)
    }
}

 // Set up periodic updates
go func() {
    ticker := time.NewTicker(updateInterval)
    for range ticker.C {
        updateData(repo)
    }
}()

// Initial data update
updateData(repo)
```
`git_handler.go`
```
package main

import (
    "fmt"
    "github.com/go-git/go-git/v5"
    "github.com/go-git/go-git/v5/plumbing"
    "github.com/go-git/go-git/v5/storage/memory"
)

var (
    repo *git.Repository
)

func InitGitHandler() {
    // Repository will be initialized when first needed
}

func getRepo() (*git.Repository, error) {
    if repo == nil {
        var err error
        repo, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
            URL:           repoURL,
            ReferenceName: plumbing.HEAD,
            SingleBranch:  true,
            Depth:         1,
        })
        if err != nil {
            return nil, fmt.Errorf("error cloning repository: %v", err)
        }
    }
    return repo, nil
}

func pullRepo() error {
    r, err := getRepo()
    if err != nil {
        return err
    }

    w, err := r.Worktree()
    if err != nil {
        return fmt.Errorf("error getting worktree: %v", err)
    }

    err = w.Pull(&git.PullOptions{})
    if err != nil && err != git.NoErrAlreadyUpToDate {
        return fmt.Errorf("error pulling updates: %v", err)
    }

    return nil
}
```
`data_processor.go`
```
package main

import (
    "encoding/json"
    "fmt"
    "gopkg.in/yaml.v2"
    "path/filepath"
    "sync"
)

var (
    cachedData          map[string]string
    lastProcessedCommit string
    mu                  sync.RWMutex
)

func InitDataProcessor() {
    cachedData = make(map[string]string)
}

func UpdateDataIfNeeded() error {
    mu.Lock()
    defer mu.Unlock()

    r, err := getRepo()
    if err != nil {
        return err
    }

    err = pullRepo()
    if err != nil {
        return err
    }

    ref, err := r.Head()
    if err != nil {
        return fmt.Errorf("error getting HEAD: %v", err)
    }

    currentCommit := ref.Hash().String()
    if currentCommit == lastProcessedCommit {
        return nil // No changes, no need to update
    }

    // Process and update the data
    newData, err := processConfigFiles(r)
    if err != nil {
        return fmt.Errorf("error processing config files: %v", err)
    }

    // Update the cache and last processed commit
    cachedData = newData
    lastProcessedCommit = currentCommit

    return nil
}

func processConfigFiles(r *git.Repository) (map[string]string, error) {
    newData := make(map[string]string)

    tree, err := r.Worktree()
    if err != nil {
        return nil, fmt.Errorf("error getting worktree: %v", err)
    }

    files, err := tree.Filesystem.ReadDir(configPath)
    if err != nil {
        return nil, fmt.Errorf("error reading config directory: %v", err)
    }

    for _, file := range files {
        if !file.IsDir() && filepath.Ext(file.Name()) == ".yaml" {
            content, err := tree.Filesystem.Open(filepath.Join(configPath, file.Name()))
            if err != nil {
                return nil, fmt.Errorf("error opening file %s: %v", file.Name(), err)
            }
            defer content.Close()

            var yamlData map[string]interface{}
            err = yaml.NewDecoder(content).Decode(&yamlData)
            if err != nil {
                return nil, fmt.Errorf("error parsing YAML from %s: %v", file.Name(), err)
            }

            jsonData, err := json.Marshal(yamlData)
            if err != nil {
                return nil, fmt.Errorf("error converting to JSON for %s: %v", file.Name(), err)
            }

            newData[file.Name()] = string(jsonData)
        }
    }

    return newData, nil
}

func GetDataForUser() map[string]string {
    mu.RLock()
    defer mu.RUnlock()
    return cachedData
}
```
`config.go`
```
package main

const (
    repoURL    = "ssh://git@your-repo-url.git"
    configPath = "path/to/your/config/files"
)

func InitConfig() {
    // This function can be used to load configuration from environment variables or config files
}
```
* * *
* * *
# What each file does 
Go through each Go file and describe the functions within them.

1. `main.go`:

- `main()`:
  - Initializes the configuration, Git handler, and data processor.
  - Starts the periodic update process in a separate goroutine.
  - This is where you would set up your web server and request handlers.
  - Currently, it just runs indefinitely as a placeholder.

- `periodicUpdate()`:
  - Runs in an infinite loop.
  - Calls `UpdateDataIfNeeded()` every 5 minutes.
  - If there's an error updating the data, it prints the error message.

2. `git_handler.go`:

- `InitGitHandler()`:
  - Currently empty, but could be used for any initialization needed for the Git handler.

- `getRepo()`:
  - Checks if the repository has been cloned.
  - If not, it clones the repository into memory.
  - Returns the repository object or an error.

- `pullRepo()`:
  - Gets the repository object using `getRepo()`.
  - Retrieves the worktree from the repository.
  - Attempts to pull the latest changes.
  - Returns an error if there's a problem pulling, unless it's just that the repo is already up to date.

3. `data_processor.go`:

- `InitDataProcessor()`:
  - Initializes the `cachedData` map.

- `UpdateDataIfNeeded()`:
  - Locks the mutex to ensure thread-safety.
  - Gets the repository and pulls the latest changes.
  - Checks if the current commit is different from the last processed commit.
  - If there are changes, it processes the config files and updates the cached data.
  - Updates the `lastProcessedCommit` to the current commit.

- `processConfigFiles()`:
  - Takes a Git repository object as input.
  - Reads all YAML files in the specified config directory.
  - For each YAML file:
    - Reads the content
    - Parses the YAML into a map
    - Converts the map to JSON
    - Stores the JSON string in the `newData` map
  - Returns the `newData` map or an error.

- `GetDataForUser()`:
  - Uses a read lock to ensure thread-safety.
  - Returns the current `cachedData` map.

4. `config.go`:

- `InitConfig()`:
  - Currently empty, but could be used to load configuration from environment variables or config files.

The `config.go` file also defines two constants:
- `repoURL`: The URL of the Git repository to clone and pull from.
- `configPath`: The path within the repository where the config files are located.

This structure allows for a clear separation of concerns:
- Git operations are handled in `git_handler.go`
- Data processing and caching are managed in `data_processor.go`
- The main application flow and periodic updates are controlled in `main.go`
- Configuration is centralized in `config.go`

This design makes the code more maintainable and easier to extend or modify in the future.

Extra: 

```
func BuildKafkaConfig(rawTopics *RawTopicsConfig, rawServices map[string]RawServiceConfig) *KafkaConfig {
    kafkaConfig := &KafkaConfig{
        Topics:   make(map[string]TopicInfo),
        Services: make(map[string]ServiceInfo),
    }

    // Process topics
    for topicName, rawTopic := range rawTopics.Topics {
        kafkaConfig.Topics[topicName] = TopicInfo{
            Partitions: rawTopic.Prod.Partitions,
            Producers:  []string{},
            Consumers:  make(map[string][]string),
        }
    }

    // Process services
    for serviceName, rawService := range rawServices {
        serviceInfo := ServiceInfo{
            ProducedTopics: make(map[string]ProducerInfo),
            ConsumedTopics: make(map[string]ConsumerInfo),
        }

        // Process producers
        for eventName, eventInfo := range rawService.Bus.Producers {
            for _, topic := range eventInfo.Topics {
                if _, exists := serviceInfo.ProducedTopics[topic]; !exists {
                    serviceInfo.ProducedTopics[topic] = ProducerInfo{
                        PartitionAlgorithm: "unknown", // You might need to get this from somewhere else
                        EventTypes:         []string{},
                    }
                }
                serviceInfo.ProducedTopics[topic].EventTypes = append(serviceInfo.ProducedTopics[topic].EventTypes, eventName)
                
                // Update Topics
                if !contains(kafkaConfig.Topics[topic].Producers, serviceName) {
                    kafkaConfig.Topics[topic].Producers = append(kafkaConfig.Topics[topic].Producers, serviceName)
                }
            }
        }

        // Process consumers
        for eventName, eventInfo := range rawService.Bus.Consumers {
            for _, topic := range eventInfo.Topics {
                if _, exists := serviceInfo.ConsumedTopics[topic]; !exists {
                    serviceInfo.ConsumedTopics[topic] = ConsumerInfo{
                        ConsumerGroup:      eventInfo.GroupID,
                        ConsumedEventTypes: []string{},
                    }
                }
                serviceInfo.ConsumedTopics[topic].ConsumedEventTypes = append(serviceInfo.ConsumedTopics[topic].ConsumedEventTypes, eventName)
                
                // Update Topics
                if _, exists := kafkaConfig.Topics[topic].Consumers[eventInfo.GroupID]; !exists {
                    kafkaConfig.Topics[topic].Consumers[eventInfo.GroupID] = []string{}
                }
                if !contains(kafkaConfig.Topics[topic].Consumers[eventInfo.GroupID], serviceName) {
                    kafkaConfig.Topics[topic].Consumers[eventInfo.GroupID] = append(kafkaConfig.Topics[topic].Consumers[eventInfo.GroupID], serviceName)
                }
            }
        }

        kafkaConfig.Services[serviceName] = serviceInfo
    }

    return kafkaConfig
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```
