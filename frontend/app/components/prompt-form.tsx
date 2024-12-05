import { useState } from "react";
import { FaGithub } from "react-icons/fa";

interface PromptFormProps {
  onSubmit: (query: string, repoPath: string, ignorePatterns: string[]) => void;
}

export default function PromptForm({ onSubmit }: PromptFormProps) {
  const [query, setQuery] = useState("");
  const [repoPath, setRepoPath] = useState(() => {
    if (typeof window !== "undefined") {
      return localStorage.getItem("repoPath") || "";
    }
    return "";
  });
  const [ignorePatterns, setIgnorePatterns] = useState<string>(() => {
    if (typeof window !== "undefined") {
      return JSON.parse(localStorage.getItem("ignorePatterns") || "[]");
    }
    return [];
  });
  const [errors, setErrors] = useState({
    query: "",
    repoPath: "",
    ignorePatterns: "",
  });

  const validateForm = () => {
    const newErrors = { query: "", repoPath: "", ignorePatterns: "" };
    let isValid = true;

    if (!query.trim()) {
      newErrors.query = "Please enter a question";
      isValid = false;
    }

    if (!repoPath.trim()) {
      newErrors.repoPath = "Please enter a repository URL";
      isValid = false;
    } else if (!repoPath.match(/^https:\/\/github\.com\/[\w-]+\/[\w-]+$/)) {
      newErrors.repoPath = "Please enter a valid GitHub repository URL";
      isValid = false;
    }

    setErrors(newErrors);
    return isValid;
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (validateForm()) {
      const patterns = ignorePatterns
        .split("\n")
        .map((p) => p.trim())
        .filter((p) => p !== "");

      onSubmit(query, repoPath, patterns);
    }
  };

  const handleRepoPathChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setRepoPath(newValue);
    if (typeof window !== "undefined") {
      localStorage.setItem("repoPath", newValue);
    }
  };

  return (
    <form className="flex flex-col gap-2" onSubmit={handleSubmit}>
      <div className="relative flex w-min items-center rounded-base border-2 overflow-x-hidden border-border dark:border-darkBorder font-base shadow-light dark:shadow-dark">
        <span className="absolute left-3 text-gray-500">
          <FaGithub />
        </span>
        <input
          className={`bg-white dark:bg-secondaryBlack w-[250px] p-[5px] pl-10 outline-none text-xs ${
            errors.repoPath ? "border-red-500" : ""
          }`}
          type="text"
          name="repoPath"
          id="repoPath"
          placeholder="https://github.com/owner/repository"
          value={repoPath}
          onChange={handleRepoPathChange}
        />
      </div>

      <div className="flex w-min items-center rounded-base border-2 overflow-x-hidden border-border dark:border-darkBorder font-base shadow-light dark:shadow-dark">
        <input
          className={`bg-white dark:bg-secondaryBlack w-[400px] min-w-[14ch] p-[10px] outline-none ${
            errors.query ? "border-red-500" : ""
          }`}
          type="text"
          name="query"
          id="query"
          placeholder="Start asking..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
        <button
          className="border-l-2 text-text border-border dark:border-darkBorder bg-main p-[10px] sm:px-5 px-3"
          type="submit"
          aria-label="Submit Prompt"
        >
          Go
        </button>
      </div>

      <div className="relative flex w-min items-center rounded-base border-2 overflow-x-hidden border-border dark:border-darkBorder font-base shadow-light dark:shadow-dark">
        <textarea
          className={`bg-white dark:bg-secondaryBlack w-[350px] p-[5px] outline-none text-xs resize-y min-h-[60px]`}
          name="ignorePatterns"
          id="ignorePatterns"
          placeholder="Enter ignore patterns (one per line):&#10;*.test.ts"
          value={ignorePatterns}
          onChange={(e) => {
            setIgnorePatterns(e.target.value);
            if (typeof window !== "undefined") {
              localStorage.setItem(
                "ignorePatterns",
                JSON.stringify(e.target.value)
              );
            }
          }}
          onKeyDown={(e) => {
            if (e.shiftKey && e.key === "Enter") {
              e.stopPropagation();
            }
          }}
        />
      </div>

      {errors && (
        <span className="text-red-500 text-xs ">
          {errors.query || errors.repoPath}
        </span>
      )}
    </form>
  );
}
