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
  2. Retrieves the connection string for the running container.
  3. Creates a new Ollama language model instance, which is used as the chat model.
  4. The chat model is asked directly for a response to a fixed question. The model does not have any context about the question.
  1. Runs an Ollama container using Testcontainers for Golang. The image used is `mdelapenya/all-minilm:0.3.13-22m`, loading the `all-minilm:22m` model.
  3. From this Ollama container, it creates a new Ollama language model instance, which is used as the embedder for the RAG model.
  4. Runs a store container using Testcontainers for Golang, and it is used to store and retrieve embeddings for the RAG.
  5. Ingests some markdown documents about Testcontainers Cloud into the vector store, using the embedder.
  6. Performs a search in the store to retrieve the most similar embeddings to the original fixed question.
  7. If there are no results, the program exits with an error message.
  8. If there are results, the program builds a chat language model using Ollama (image `mdelapenya/llama3.2:0.3.13-1b` and model `llama3.2:1b`).
  9. Using the relevant content from the store search results, the program generates a streaming response to the user's prompt.

## Running the Example

To run the example, navigate to the `07-rag` directory and run the following command:

```sh
go run .
```

The application will start two containerized Ollama language models and generate text based on the augmented prompt using RAG. The generated text will be displayed in the console.

```shell
2024/11/12 18:16:20 How I can enable verbose logging in Testcontainers Desktop?
>> Straight answer:
To enable verbose logging in Testcontainers Desktop, you can set the `TESTCONTAINERS_LOG_LEVEL` environment variable to "DEBUG" before running the application. This will increase the log level and provide more detailed output.

Alternatively, you can also use the `-Dtestcontainers.log.level=DEBUG` JVM option when starting the Docker container.

Note that Testcontainers Desktop provides a built-in logging mechanism, so you may need to adjust other logging settings or configurations to achieve the desired verbosity.
2024/11/12 18:17:35 Ingesting document: knowledge/txt/simple-local-development-with-testcontainers-desktop.txt
2024/11/12 18:17:35 Ingesting document: knowledge/txt/tc-guide-introducing-testcontainers.txt
2024/11/12 18:17:35 Ingesting document: knowledge/txt/tcc.txt
2024/11/12 18:17:37 Ingested 3 documents
2024/11/12 18:17:37 run: build ragged chat: similarity search: empty response
exit status 1
```
