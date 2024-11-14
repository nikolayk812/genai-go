# Generative AI with Testcontainers for Golang

This project demonstrates how to use [Testcontainers for Golang](https://github.com/testcontainers/testcontainers-go) to create a seamless development environment for building Generative AI applications.

## Project Structure

1. [`01-hello-world`](./01-hello-world): Contains a simple example of using a language model to generate text.
1. [`02-streaming`](./02-streaming): Contains an example of using a language model to generate text in streaming mode.
1. [`03-chat`](./03-chat): Contains an example of using a language model to generate text in a chat application.
1. [`04-vision-model`](./04-vision-model): Contains an example of using a vision model to generate text from images.
1. [`05-augmented-generation`](./05-augmented-generation): Contains an example of augmenting the prompt with additional information to generate more accurate text.
1. [`06-embeddings`](./06-embeddings): Contains an example of generating embeddings from text and calculating similarity between them.
1. [`07-rag`](./07-rag): Contains an example of applying RAG (Retrieval-Augmented Generation) to generate better responses.
1. [`08-testing`](./08-testing): Contains an example with the evolution of testing our Generative AI applications, from an old school approach to a more modern one using Validator Agents.

## Prerequisites

- Go 1.23 or higher
- Docker

## Setup

1. Clone the repository:
    ```sh
    git clone https://github.com/mdelapenya/generative-ai-with-testcontainers.git
    cd generative-ai-with-testcontainers
    ```

## Running the Examples

To run the examples, navigate to the desired directory and run the `go run .` command. For example, to run the `1-hello-world` example:

```sh
cd 1-hello-world
go run .
```
