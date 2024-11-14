# 03-chat

Contains a simple example of using a language model to generate text in an interactive chat application.

## Libraries Involved

- `github.com/testcontainers/testcontainers-go`: [Testcontainers for Golang](https://github.com/testcontainers/testcontainers-go) is library for running Docker containers for integration tests.
- `github.com/testcontainers/testcontainers-go/modules/ollama`: A module for running Ollama language models using Testcontainers.
- `github.com/tmc/langchaingo`: A library for interacting with language models.
- `github.com/tmc/langchaingo/llms/ollama`: A specific implementation of the language model interface for Ollama.

## Code Explanation

The code in `main.go` sets up and runs a containerized Ollama language model using Testcontainers, then uses the model to generate text based on an interactive prompt.

### Main Functions

- `main()`: The entry point of the application. It calls the `run()` function and logs any errors.
- `run()`: The main logic of the application. It performs the following steps:
  1. Runs an Ollama container using Testcontainers. The image used is `mdelapenya/llama3.2:0.3.13-1b`, loading the `llama3.2:1b` model.
  2. Retrieves the connection string for the running container.
  3. Creates a new Ollama language model instance.
  4. Defines an infinite loop to interact with the language model in a chat-like manner.
  5. Generates the content and prints it to the console based on the user's input.
  6. Exits the interactive loop if the user types `exit`, `quit`, or hits `Ctrl+C`.

## Running the Example

To run the example, navigate to the `03-chat` directory and run the following command:

```sh
go run .
```

The application will start a containerized Ollama language model and generate text based on the interactive prompt.
The generated text will be displayed in the console, and the user can continue interacting with the model until they choose to exit.

```shell
go run .

You: what is the capital of Japan
The capital of Japan is Tokyo.
You: ^C
Interrupt signal received, ending chat session
```
