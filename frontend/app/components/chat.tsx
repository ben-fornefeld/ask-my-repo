import { useState, useCallback, useRef, useEffect } from "react";
import PromptForm from "./prompt-form";
import RenderCompletion from "./render-completion";
import type {
  QueryResponseChunk,
  ParsedChunk,
  RankedChunk,
} from "../lib/types";
import { motion } from "framer-motion";
import { brutalistSlideMotion } from "../lib/utils";
import Badge from "./ui/badge";

interface ChatState {
  isLoading: boolean;
  error?: string;
  parsedChunks: ParsedChunk[];
  rankedChunks: RankedChunk[];
  completion?: string;
  completionBuffer: string;
}

export default function Chat() {
  const [state, setState] = useState<ChatState>({
    isLoading: false,
    parsedChunks: [],
    rankedChunks: [],
    completionBuffer: "",
  });

  const batchTimeoutRef = useRef<NodeJS.Timeout | null>(null);
  const completionQueueRef = useRef<string>("");
  const isProcessingRef = useRef<boolean>(false);

  const insertCompletionInBatches = useCallback(async (content: string) => {
    const batchSize = 2;
    const intervalMs = 6;
    let index = 0;

    return new Promise<void>((resolve) => {
      const intervalId = setInterval(() => {
        if (index < content.length) {
          const batch = content.slice(index, index + batchSize);
          setState((prevState) => ({
            ...prevState,
            completion: (prevState.completion || "") + batch,
          }));
          index += batchSize;
        } else {
          clearInterval(intervalId);
          resolve();
        }
      }, intervalMs);
    });
  }, []);

  const processQueue = useCallback(async () => {
    if (isProcessingRef.current || !completionQueueRef.current) {
      return;
    }

    isProcessingRef.current = true;
    const content = completionQueueRef.current;
    completionQueueRef.current = "";

    await insertCompletionInBatches(content);
    isProcessingRef.current = false;

    if (completionQueueRef.current) {
      processQueue();
    }
  }, [insertCompletionInBatches]);

  const handleSubmit = useCallback(
    async (query: string, repopath: string) => {
      setState({
        isLoading: true,
        parsedChunks: [],
        rankedChunks: [],
        completionBuffer: "",
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

              switch (chunk.type) {
                case "ranking.parsed": {
                  if (chunk.parsed_chunk) {
                    setState((prevState) => ({
                      ...prevState,
                      parsedChunks: [
                        ...prevState.parsedChunks,
                        chunk.parsed_chunk!,
                      ],
                    }));
                  }
                  break;
                }

                case "ranking.ranked": {
                  if (chunk.ranked_chunk) {
                    setState((prevState) => ({
                      ...prevState,
                      rankedChunks: [
                        ...prevState.rankedChunks,
                        chunk.ranked_chunk!,
                      ],
                    }));
                  }
                  break;
                }

                case "completion.delta": {
                  if (chunk.completion) {
                    completionQueueRef.current += chunk.completion;
                    await processQueue();
                  }
                  break;
                }

                case "error": {
                  setState((prevState) => ({
                    ...prevState,
                    error: chunk.error,
                    isLoading: false,
                  }));
                  break;
                }
              }
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
      } finally {
        isProcessingRef.current = false;
        setState((prev) => ({
          ...prev,
          isLoading: false,
        }));
      }
    },
    [processQueue]
  );

  useEffect(() => {
    return () => {
      if (batchTimeoutRef.current) {
        clearTimeout(batchTimeoutRef.current);
      }
    };
  }, []);

  return (
    <div className="flex flex-col gap-4 mx-auto max-w-2xl">
      {state.isLoading && <div>Loading...</div>}

      {state.error && <div className="text-red-500">Error: {state.error}</div>}

      {state.parsedChunks.length > 0 && (
        <motion.div
          variants={brutalistSlideMotion}
          initial="hidden"
          animate="visible"
          exit="exit"
        >
          Processed {state.parsedChunks.length} files...
        </motion.div>
      )}

      {state.rankedChunks.length > 0 && (
        <motion.div
          variants={brutalistSlideMotion}
          initial="hidden"
          animate="visible"
          exit="exit"
          className="bg-white w-full flex gap-3 flex-wrap dark:bg-secondaryBlack p-2 rounded-base border-2 border-border dark:border-darkBorder shadow-light dark:shadow-dark"
        >
          <p className="w-full">
            Found {state.rankedChunks.length} relevant files...
          </p>
          {state.rankedChunks.map((chunk) => (
            <Badge
              key={chunk.ParsedChunk.FilePath}
              text={chunk.ParsedChunk.FilePath}
              className="text-xs"
            />
          ))}
        </motion.div>
      )}

      {state.completion && (
        <RenderCompletion completion={state.completion || ""} />
      )}

      <PromptForm onSubmit={handleSubmit} />
    </div>
  );
}
