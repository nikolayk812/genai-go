# 07-rag

Contains a simple example of using a language model to answer questions based on a given prompt using RAG (Retrieval-Augmented Generation).

## Libraries Involved

- `github.com/testcontainers/testcontainers-go`: A library for running Docker containers for integration tests.
- `github.com/testcontainers/testcontainers-go/modules/ollama`: A module for running Ollama language models using Testcontainers.
- `github.com/testcontainers/testcontainers-go/modules/weaviate`: A module for running Weaviate vector search engines using Testcontainers.
- `github.com/tmc/langchaingo`: A library for interacting with language models.
- `github.com/tmc/langchaingo/llms/ollama`: A specific implementation of the language model interface for Ollama.
- `github.com/tmc/langchaingo/vectorstores`: An interface for interacting with vector search engines.
- `github.com/tmc/langchaingo/vectorstores/weaviate`: A specific implementation of the vector store interface for Weaviate.

## Code Explanation

The code in `main.go` sets up and runs two containerized Ollama language models and a Weaviate vector store using Testcontainers, then uses one of the models to generate the embeddings for a set of texts. It then uses the Weaviate vector store to search for similar embeddings and generate text based on the augmented prompt using RAG.

### Main Functions

- `main()`: The entry point of the application. It calls the `run()` function and logs any errors.
- `run()`: The main logic of the application. It performs the following steps:
  1. Runs an Ollama container using Testcontainers for Golang. The image used is `ilopezluna/all-minilm:0.3.13-22m`, loading the `all-minilm:22m` model.
  2. Retrieves the connection string for the running container.
  3. Creates a new Ollama language model instance, which is used as the embedder for the RAG model.
  4. Runs a Weaviate container using Testcontainers for Golang. The image used is `semitechnologies/weaviate:1.27.2`, and it is used to store and retrieve embeddings for the RAG.
  5. Ingests some example data into the Weaviate vector store.
  6. Performs a search in Weaviate to retrieve the most similar embeddings to a query.
  7. If there are no results, the program exits with an error message.
  8. If there are results, the program builds a chat language model using Ollama (image `ilopezluna/llama3.2:0.3.13-1b` and model `llama3.2:1b`).
  9. Using the relevant content from the Weaviate search results, the program generates a streaming response to the user's prompt.

## Running the Example

To run the example, navigate to the `07-rag` directory and run the following command:

```sh
go run .
```

The application will start two containerized Ollama language models and generate text based on the augmented prompt using RAG. The generated text will be displayed in the console.

```shell
What is your favourite sport?

Answer the question considering the following relevant content:
I like football

Based on the information provided, I can infer that you enjoy playing or watching football. However, since you didn't specify a particular aspect of football (e.g., team, competition level), I'll provide some general responses.

If you're interested in discussing football, here are a few questions to get started:

* What position do you play or prefer?
* Do you have a favorite team or player?
* Are you more into the competitive aspect (e.g., Premier League) or the recreational side of the sport?

Feel free to share your thoughts, and I'll be happy to engage in a conversation about football!% 
```
