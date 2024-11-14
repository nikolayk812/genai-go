# 08-testing

Contains a simple example of using a language model to validate the answers of other language models, using a Validator Agent.

## Libraries Involved

- `github.com/testcontainers/testcontainers-go`: [Testcontainers for Golang](https://github.com/testcontainers/testcontainers-go) is library for running Docker containers for integration tests.
- `github.com/testcontainers/testcontainers-go/modules/ollama`: A module for running Ollama language models using Testcontainers.
- `github.com/testcontainers/testcontainers-go/modules/postgres`: A module for running PgVector vector search engines using Testcontainers.
- `github.com/testcontainers/testcontainers-go/modules/weaviate`: A module for running Weaviate vector search engines using Testcontainers.
- `github.com/tmc/langchaingo`: A library for interacting with language models.
- `github.com/tmc/langchaingo/llms/ollama`: A specific implementation of the language model interface for Ollama.
- `github.com/tmc/langchaingo/vectorstores`: An interface for interacting with vector search engines.
- `github.com/tmc/langchaingo/vectorstores/pgvector`: A specific implementation of the vector store interface for PgVector.
- `github.com/tmc/langchaingo/vectorstores/weaviate`: A specific implementation of the vector store interface for Weaviate.

## Code Explanation

The code in `main.go` prints out two different responses for the same task: one for talking to a model in a straight manner, and the second using RAG. For that, it sets up and runs two containerized Ollama language models and a vector store using Testcontainers, then uses one of the models to generate the embeddings for a set of texts. It then uses the selected vector store to search for similar embeddings and generate text based on the augmented prompt using RAG.

The vector store to use is `weaviate` by default, but it can be changed to `pgvector` by setting the `VECTOR_STORE` environment variable to `pgvector`. 

- The image used for Weaviate is `semitechnologies/weaviate:1.27.2`.
- The image used for PgVector is `pgvector/pgvector:pg16`.

We are adding tests to demonstrate how to validate the answers of the language models. We will use a Validator Agent to do so.

### Main Functions

- `main()`: The entry point of the application. It calls the `run()` function and logs any errors.
- `run()`: The main logic of the application. It performs the following steps:
  1. Runs an Ollama container using Testcontainers. The image used is `mdelapenya/llama3.2:0.3.13-1b`, loading the `llama3.2:1b` model.
  1. Retrieves the connection string for the running container.
  1. Creates a new Ollama language model instance, which is used as the chat model.
  1. The chat model is asked directly for a response to a fixed question. The model does not have any context about the question.
  1. Runs an Ollama container using Testcontainers. The image used is `mdelapenya/all-minilm:0.3.13-2m`, loading the `all-minilm:22m` model, which is useful for large text generation.
  1. From this Ollama container, it creates a new Ollama language model instance, which is used as the embedder for the RAG model.
  1. Runs a store container using Testcontainers, and it is used to store and retrieve embeddings for the RAG.
  1. Ingests some markdown documents about Testcontainers Cloud into the vector store, using the embedder. The files are ingested using chunks of 1024 characters.
  1. Performs a search in the store to retrieve the most similar embeddings to the original fixed question.
  1. If there are no results, the program exits with an error message.
  1. If there are results, the program builds a chat language model using Ollama (image `mdelapenya/llama3.2:0.3.13-1b` and model `llama3.2:1b`).
  1. Using the relevant content from the store search results, the program generates a streaming response to the user's prompt.

## Running the Example

To run the example, navigate to the `08-testing` directory and run the following command:

```sh
go run .
```

The application will start two containerized Ollama language models and generate text based on the augmented prompt using RAG. The generated text will be displayed in the console.

```shell
2024/11/14 12:54:19 How I can enable verbose logging in Testcontainers Desktop?
>> Straight answer:
 To enable verbose logging in Testcontainers Desktop, you can set the `TESTCONTAINERS_LOG_LEVEL` environment variable to "DEBUG" before running the application. This will increase the log level and provide more detailed output.

Alternatively, you can also specify a logging configuration file that includes the desired log level. For example, you can create a file named `logging.conf` with the following content:

[loggers]
keys=root

[handlers]
keys=console

[formatters]
keys=simple

[logger_root]
level=DEBUG
qname=
level=DEBUG
type=console
args=-v

[handler_console]
class=StreamHandler
args=-v
level=
formatter=simple
args=

[formatter_simple]
format=%(asctime)s - %(name)-15j - %(levelname)-8s - %(message)s

Then, you can run Testcontainers Desktop with the `--config` option to specify the logging configuration file.
2024/11/14 12:54:34 Ingesting document: knowledge/txt/simple-local-development-with-testcontainers-desktop.txt
2024/11/14 12:54:34 Ingesting document: knowledge/txt/tc-guide-introducing-testcontainers.txt
2024/11/14 12:54:34 Ingesting document: knowledge/txt/tcc.txt
2024/11/14 12:54:44 Ingested 82 documents
2024/11/14 12:54:44 Relevant documents for RAG: 3
>> Ragged answer:
 To enable verbose logging in Testcontainers Desktop, you can add a property to your per-user configuration file (~/.testcontainers.properties) with the key "cloud.logs.verbose" and set its value to true. Alternatively, you can use the --verbose flag when running the client or set the environmental variable TC_CLOUD_LOGS_VERBOSE=true.
```

## How to test this (1): String comparison

To test this, what we would usually do is to create a test file with two tests, one for the straight answer and another for the ragged answer. We would then run the tests and check if the output matches the expected output. Just take a look at the `main_test.go` file in the `08-testing` directory, and its `Test1_oldSchool` test function, and then run the tests:

```shell
go test -timeout 600s -run ^Test1_oldSchool$ github.com/mdelapenya/genai-testcontainers-go/testing -v -count=1
```

## How to test this (2): Embeddings

In a second iteration, we remembered that we now know how to create emebeddings and calculate the cosine similarity. So we create two tests more in the test file, one for the straight answer and another for the ragged answer. We would then run the tests and check if the cosine similarity is higher thatn 0.8 (our threshold). Just take a look at the `main_test.go` file in the `08-testing` directory, and its `Test2_embeddings` test function, and then run the tests:

```shell
go test -timeout 600s -run ^Test2_embeddings$ github.com/mdelapenya/genai-testcontainers-go/testing -v -count=1
```

## How to test this (3): Validator Agents

Finally, in a third iteration, we realised that we have a lot of power with LLMs, and it would be cool to use one to validate the answers. We could be as strict as needed defining the System and User prompts, in order for the validator agent to be very specific about the answer. We can even provide an output format for the answer, so the validator agent can check if the answer is correct. Just take a look at the `main_test.go` file in the `08-testing` directory, and its `Test3_validatorAgent` test function, and then run the tests:

```shell
go test -timeout 600s -run ^Test3_validatorAgent$ github.com/mdelapenya/genai-testcontainers-go/testing -v -count=1
```
