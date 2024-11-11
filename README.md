# Generative AI with Testcontainers for Golang

This project demonstrates how to use Testcontainers for Golang to create a seamless development environment for building Generative AI applications.

## Project Structure

1. [`01-hello-world`](./01-hello-world): Contains a simple example of using a language model to generate text.
2. [`02-streaming`](./02-streaming): Contains an example of using a language model to generate text in streaming mode.
3. [`03-chat`](./03-chat): Contains an example of using a language model to generate text in a chat application.
4. [`04-vision-model`](./04-vision-model): Contains an example of using a vision model to generate text from images.
5. [`05-augmented-generation`](./05-augmented-generation): Contains an example of augmenting the prompt with additional information to generate more accurate text.
6. [`06-embeddings`](./06-embeddings): Contains an example of generating embeddings from text and calculating similarity between them.

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
