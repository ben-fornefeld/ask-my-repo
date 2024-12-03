import { useState, useCallback } from "react";
import PromptForm from "./prompt-form";
import RenderCompletion from "./render-completion";
import type {
  QueryResponseChunk,
  ParsedChunk,
  RankedChunk,
} from "../lib/types";

interface ChatState {
  isLoading: boolean;
  error?: string;
  parsedChunks: ParsedChunk[];
  rankedChunks: RankedChunk[];
  completion?: string;
}

export default function Chat() {
  const [state, setState] = useState<ChatState>({
    isLoading: false,
    parsedChunks: [],
    rankedChunks: [],
  });

  const handleSubmit = useCallback(async (query: string, repopath: string) => {
    // Reset state
    setState({
      isLoading: true,
      parsedChunks: [],
      rankedChunks: [],
    });

    const body = {
      query,
      repopath,
      ignorepatterns: ["example/"],
      scorethreshold: 0.2,
    };

    try {
      const response = await fetch(
        `${import.meta.env.VITE_BACKEND_URL}/query`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(body),
          credentials: "include",
        }
      );

      if (!response.ok || !response.body) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const reader = response.body!.getReader();
      const decoder = new TextDecoder("utf-8");
      let buffer = "";

      while (true as const) {
        const { value, done } = await reader.read();
        if (done) break;

        buffer += decoder.decode(value, { stream: true });

        const messages = buffer.split("\n\n");
        buffer = messages.pop() || "";

        for (const message of messages) {
          if (message.startsWith("data: ")) {
            const jsonStr = message.slice(6);
            const chunk: QueryResponseChunk = JSON.parse(jsonStr);

            setState((prevState) => {
              switch (chunk.type) {
                case "ranking.parsed":
                  return chunk.parsed_chunk
                    ? {
                        ...prevState,
                        parsedChunks: [
                          ...prevState.parsedChunks,
                          chunk.parsed_chunk,
                        ],
                      }
                    : prevState;

                case "ranking.ranked":
                  return chunk.ranked_chunk
                    ? {
                        ...prevState,
                        rankedChunks: [
                          ...prevState.rankedChunks,
                          chunk.ranked_chunk,
                        ],
                      }
                    : prevState;

                case "completion.delta":
                  return chunk.completion
                    ? {
                        ...prevState,
                        completion:
                          (prevState.completion || "") + chunk.completion,
                        isLoading: false,
                      }
                    : prevState;

                case "error":
                  return {
                    ...prevState,
                    error: chunk.error,
                    isLoading: false,
                  };

                default:
                  return prevState;
              }
            });
          }
        }
      }
    } catch (error) {
      setState((prev) => ({
        ...prev,
        isLoading: false,
        error:
          error instanceof Error ? error.message : "Unknown error occurred",
      }));
    }
  }, []);

  return (
    <div className="flex flex-col gap-4 mx-auto max-w-2xl">
      {state.isLoading && <div>Loading...</div>}

      {state.error && <div className="text-red-500">Error: {state.error}</div>}

      {state.parsedChunks.length > 0 && (
        <div className="text-gray-500">
          Processed {state.parsedChunks.length} files...
        </div>
      )}

      {state.rankedChunks.length > 0 && (
        <div className="text-gray-500">
          Found {state.rankedChunks.length} relevant files...
        </div>
      )}

      {state.completion && <RenderCompletion completion={state.completion} />}

      <PromptForm onSubmit={handleSubmit} />
    </div>
  );
}
