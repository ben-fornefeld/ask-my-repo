# Ask My Repo

A Go server and modern frontend for efficient Large Language Model (LLM) code context retrieval, using model-based ranking. It enables users to query public GitHub repositories, returning the most relevant code sections for their questions.

## Features

- **Backend (Go)**
  - Clones public GitHub repositories, parses and chunks text/code files.
  - Uses custom ignore patterns for file selection.
  - Ranks code chunks' relevance to a user query using LLMs (Anthropic, Replicate).
  - REST API endpoint for queries with CORS support.
  - Modular design with parser, ranking, and completion engines.

- **Frontend (Remix, React, TypeScript)**
  - User interface to input a GitHub repo, question, and ignore patterns.
  - Sends queries to the backend and displays ranked answers.
  - Built with Remix, Tailwind CSS, and React Query for modern developer experience.

## Getting Started

### Prerequisites

- Go (for backend)
- Node.js, npm or bun (for frontend)
- Doppler (for environment management)
- API keys for Anthropic and Replicate (set as environment variables)

### Development

In the project root, use the provided Makefile:

```sh
make dev
