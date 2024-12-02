import ReactMarkdown from "react-markdown";

interface RenderCompletionProps {
  completion: string;
}

export default function RenderCompletion({
  completion,
}: RenderCompletionProps) {
  return (
    <div className="prose dark:prose-invert max-w-none">
      <ReactMarkdown>{completion}</ReactMarkdown>
    </div>
  );
}
