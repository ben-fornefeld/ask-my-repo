import { brutalistMotion } from "../lib/utils";
import { motion } from "framer-motion";
import ReactMarkdown from "react-markdown";

interface RenderCompletionProps {
  completion: string;
}

export default function RenderCompletion({
  completion,
}: RenderCompletionProps) {
  return (
    <motion.div
      variants={brutalistMotion}
      initial="hidden"
      animate="visible"
      className="prose dark:prose-invert max-w-none p-4 bg-white dark:bg-secondaryBlack rounded-base border-2 border-border dark:border-darkBorder shadow-light dark:shadow-dark"
    >
      <ReactMarkdown>{completion}</ReactMarkdown>
    </motion.div>
  );
}
