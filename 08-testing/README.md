# 08-testing

Contains a simple example of using a language model to answer questions based on a given prompt using RAG (Retrieval-Augmented Generation).

## Libraries Involved

- `github.com/testcontainers/testcontainers-go`: A library for running Docker containers for integration tests.
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

### Main Functions

- `main()`: The entry point of the application. It calls the `run()` function and logs any errors.
- `run()`: The main logic of the application. It performs the following steps:
  1. Runs an Ollama container using Testcontainers for Golang. The image used is `mdelapenya/llama3.2:0.3.13-1b`, loading the `llama3.2:1b` model.
  1. Retrieves the connection string for the running container.
  1. Creates a new Ollama language model instance, which is used as the chat model.
  1. The chat model is asked directly for a response to a fixed question. The model does not have any context about the question.
  1. Runs an Ollama container using Testcontainers for Golang. The image used is `mdelapenya/bge-m3:0.3.13-567m`, loading the `bge-m3:567m` model, which is useful for large text generation.
  1. From this Ollama container, it creates a new Ollama language model instance, which is used as the embedder for the RAG model.
  1. Runs a store container using Testcontainers for Golang, and it is used to store and retrieve embeddings for the RAG.
  1. Ingests some markdown documents about Testcontainers Cloud into the vector store, using the embedder.
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
2024/11/14 01:40:28 How I can enable verbose logging in Testcontainers Desktop?
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
2024/11/14 01:40:42 Ingesting document: knowledge/txt/simple-local-development-with-testcontainers-desktop.txt
2024/11/14 01:40:42 Ingesting document: knowledge/txt/tc-guide-introducing-testcontainers.txt
2024/11/14 01:40:42 Ingesting document: knowledge/txt/tcc.txt
2024/11/14 01:40:45 Ingested 3 documents
2024/11/14 01:40:45 Relevant documents for RAG: 3
>> Ragged answer:
To enable verbose logging in Testcontainers Desktop, you can follow these steps:

1. Open the Testcontainers Desktop application.
2. Click on the three dots (`...`) next to the "Testcontainers" label in the top-right corner of the window.
3. Select "Preferences" from the context menu.
4. In the Preferences dialog box, click on the "Logging" tab.
5. Under the "Logging level" dropdown menu, select "DEBUG".
6. You can also adjust other logging levels such as "INFO", "WARNING", and "ERROR" to suit your needs.

Alternatively, you can also enable verbose logging by adding the following property to your `testcontainers.properties` file (usually located in the same directory as your test class):

testcontainers.log.level=DEBUG

This will set the logging level to DEBUG for all Testcontainers components. You can adjust this value to suit your needs.

Additionally, you can also enable verbose logging by adding the following JVM option when running your tests:

-Dtestcontainers.log.level=DEBUG

Note that the logging level is not limited to just "DEBUG". You can use other values such as "INFO", "WARNING", and "ERROR" to control the verbosity of the logs.% 
```

## How to test this (1)

To test this, what we would usually do is to create a test file with two tests, one for the straight answer and another for the ragged answer. We would then run the tests and check if the output matches the expected output. Just take a look at the `main_test.go` file in the `08-testing` directory, and its `Test1` test function.

```shell
go test -timeout 600s -run ^Test1$ github.com/mdelapenya/genai-testcontainers-go/testing -v -count=1
```
