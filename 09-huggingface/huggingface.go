package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	tcexec "github.com/testcontainers/testcontainers-go/exec"
)

// ollamaLifecycleHook is a testcontainers.ContainerLifecycleHooks that installs the Huggingface tooling
// (python3-pip and huggingface-hub) and downloads the model from Huggingface.
// It also creates the Ollama model file from the downloaded model file,
// creating the Ollama model from that file.
var ollamaLifecycleHook = func(model string, modelFile string) testcontainers.ContainerLifecycleHooks {
	execs := [][]string{
		{"apt-get", "update"},
		{"apt-get", "upgrade", "-y"},
		{"apt-get", "install", "-y", "python3-pip"},
		{"pip", "install", "huggingface-hub"},
		{"huggingface-cli", "download", model, modelFile, "--local-dir", "."},
		{"sh", "-c", fmt.Sprintf("echo '%s' > Modelfile", "FROM "+modelFile)},
		{"ollama", "create", modelFile, "-f", "Modelfile"},
		{"rm", modelFile},
	}

	return testcontainers.ContainerLifecycleHooks{
		PostStarts: []testcontainers.ContainerHook{
			func(ctx context.Context, ctr testcontainers.Container) error {
				var errs []error
				for _, exec := range execs {
					code, _, err := ctr.Exec(ctx, exec, tcexec.Multiplexed())
					if err != nil {
						errs = append(errs, err)
						continue
					}

					if code != 0 {
						errs = append(errs, fmt.Errorf("exec %v returned %d", exec, code))
					}
				}

				return errors.Join(errs...)
			},
		},
	}
}

// WithHuggingfaceModel is a testcontainers.CustomizeRequestOption that adds a lifecycle hook
// to the container to download a model from Huggingface.
// Please make sure the model file exists on Huggingface with the given name, case sensitive.
// The model name is used to download the model from Huggingface, and the model file is used
// to create the Ollama model.
//
// Example:
//
//	modelName = "DavidAU/DistiLabelOrca-TinyLLama-1.1B-Q8_0-GGUF"
//	modelFile = "distilabelorca-tinyllama-1.1b.Q8_0.gguf"
func WithHuggingfaceModel(modelName string, modelFile string) testcontainers.CustomizeRequestOption {
	return func(req *testcontainers.GenericContainerRequest) error {
		req.LifecycleHooks = append(req.LifecycleHooks, ollamaLifecycleHook(modelName, modelFile))
		return nil
	}
}
