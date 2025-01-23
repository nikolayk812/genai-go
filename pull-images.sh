#! /bin/bash

readonly OLLAMA_VERSION="0.5.4"

docker pull ollama/ollama:${OLLAMA_VERSION}

docker pull testcontainers/ryuk:0.11.0 &
docker pull mdelapenya/llama3.2:${OLLAMA_VERSION}-1b &
docker pull mdelapenya/llama3.2:${OLLAMA_VERSION}-3b &
docker pull mdelapenya/qwen2:${OLLAMA_VERSION}-0.5b &
docker pull mdelapenya/moondream:${OLLAMA_VERSION}-1.8b &
docker pull mdelapenya/all-minilm:${OLLAMA_VERSION}-22m &
docker pull semitechnologies/weaviate:1.27.2 &
docker pull pgvector/pgvector:pg16 &

wait
