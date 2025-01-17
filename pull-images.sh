#! /bin/bash

docker pull mdelapenya/llama3.2:0.5.4-1b
docker pull mdelapenya/llama3.2:0.5.4-3b
docker pull mdelapenya/qwen2:0.5.4-0.5b
docker pull mdelapenya/moondream:0.5.4-1.8b
docker pull mdelapenya/all-minilm:0.5.4-22m

docker pull semitechnologies/weaviate:1.27.2
docker pull pgvector/pgvector:pg16
