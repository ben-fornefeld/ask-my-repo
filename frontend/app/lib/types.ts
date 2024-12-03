export type QueryEventType =
  | "ranking.parsed"
  | "ranking.ranked"
  | "completion.delta"
  | "error";

export interface ParsedChunk {
  FilePath: string;
  Content: string;
}

export interface RankedChunk {
  ParsedChunk: ParsedChunk;
  Score: number;
}

export interface QueryResponseChunk {
  type: QueryEventType;
  parsed_chunk?: ParsedChunk;
  ranked_chunk?: RankedChunk;
  completion?: string;
  error?: string;
}
