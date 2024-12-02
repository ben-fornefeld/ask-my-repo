import { useState } from "react";
import { FaGithub } from "react-icons/fa";

interface PromptFormProps {
  onSubmit: (query: string, repoPath: string) => void;
}

export default function PromptForm({ onSubmit }: PromptFormProps) {
  const [query, setQuery] = useState("");
  const [repoPath, setRepoPath] = useState(() => {
    if (typeof window !== "undefined") {
      return localStorage.getItem("repoPath") || "";
    }
    return "";
  });
  const [errors, setErrors] = useState({ query: "", repoPath: "" });

  const validateForm = () => {
    const newErrors = { query: "", repoPath: "" };
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
      onSubmit(query, repoPath);
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
    <div className="flex flex-col gap-2">
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

      <form
        onSubmit={handleSubmit}
        className="flex w-min items-center rounded-base border-2 overflow-x-hidden border-border dark:border-darkBorder font-base shadow-light dark:shadow-dark"
      >
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
      </form>
      {errors && (
        <span className="text-red-500 text-xs ">
          {errors.query || errors.repoPath}
        </span>
      )}
    </div>
  );
}
