# 06-embeddings

Contains a simple example of using a language model to calculate the embeddings for a given set of texts, and calculate the similarity between them.

## Libraries Involved

- `github.com/testcontainers/testcontainers-go`: A library for running Docker containers for integration tests.
- `github.com/testcontainers/testcontainers-go/modules/ollama`: A module for running Ollama language models using Testcontainers.
- `github.com/tmc/langchaingo`: A library for interacting with language models.
- `github.com/tmc/langchaingo/llms/ollama`: A specific implementation of the language model interface for Ollama.

## Code Explanation

The code in `main.go` sets up and runs a containerized Ollama language model using Testcontainers, then uses the model to obtain an embedder and calculate the embeddings for a set of texts. It then calculates the similarity between the embeddings of the texts.

### Main Functions

- `main()`: The entry point of the application. It calls the `run()` function and logs any errors.
- `run()`: The main logic of the application. It performs the following steps:
  1. Runs an Ollama container using Testcontainers for Golang. The image used is `ilopezluna/all-minilm:0.3.13-22m`, loading the `all-minilm:22m` model.
  2. Retrieves the connection string for the running container.
  3. Creates a new Ollama language model instance.
  4. Defines a set of texts for which we want to calculate the embeddings.
  5. Calculates the embeddings for the texts.
  6. Calculates the similarity between the embeddings of the texts, displaying the results in the console.

## Running the Example

To run the example, navigate to the `06-embeddings` directory and run the following command:

```sh
go run .
```

The application will start a containerized Ollama language model and generate the embeddings for the provided texts.
It will then calculate the similarity between the embeddings and display the results in the console.

```shell
Similarities:
A cat is a small domesticated carnivorous mammal ~ A cat is a small domesticated carnivorous mammal = 1.00
A cat is a small domesticated carnivorous mammal ~ A tiger is a large carnivorous feline mammal = 0.68
A cat is a small domesticated carnivorous mammal ~ Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container = 0.01
A cat is a small domesticated carnivorous mammal ~ Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code. = 0.07
A tiger is a large carnivorous feline mammal ~ A cat is a small domesticated carnivorous mammal = 0.68
A tiger is a large carnivorous feline mammal ~ A tiger is a large carnivorous feline mammal = 1.00
A tiger is a large carnivorous feline mammal ~ Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container = -0.04
A tiger is a large carnivorous feline mammal ~ Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code. = -0.05
Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container ~ A cat is a small domesticated carnivorous mammal = 0.01
Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container ~ A tiger is a large carnivorous feline mammal = -0.04
Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container ~ Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container = 1.00
Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container ~ Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code. = 0.54
Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code. ~ A cat is a small domesticated carnivorous mammal = 0.07
Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code. ~ A tiger is a large carnivorous feline mammal = -0.05
Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code. ~ Testcontainers is a Go package that supports JUnit tests, providing lightweight, throwaway instances of common databases, web browsers, or anything else that can run in a Docker container = 0.54
Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code. ~ Docker is a platform designed to help developers build, share, and run container applications. We handle the tedious setup, so you can focus on the code. = 1.00
```
