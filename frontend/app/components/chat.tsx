import { useMutation } from "@tanstack/react-query";
import PromptForm from "./prompt-form";
import RenderCompletion from "./render-completion";

export default function Chat() {
  const mutation = useMutation({
    mutationFn: async ({
      query,
      repopath,
    }: {
      query: string;
      repopath: string;
    }) => {
      const body = {
        query,
        repopath,
        ignorepatterns: ["example/"],
        scorethreshold: 0.2,
      };

      const response = await fetch(
        `${import.meta.env.VITE_BACKEND_URL}/query`,
        {
          method: "POST",
          body: JSON.stringify(body),
        }
      );

      return await response.json();
    },
  });

  const handleSubmit = async (query: string, repopath: string) => {
    mutation.mutate({ query, repopath });
  };

  return (
    <div className="flex flex-col gap-4 mx-auto max-w-2xl">
      {mutation.isPending && <div>Loading...</div>}

      {mutation.isError && (
        <div className="text-red-500">
          Error: {(mutation.error as Error).message}
        </div>
      )}

      {mutation.isSuccess && mutation.data && (
        <RenderCompletion completion={mutation.data.results.Completion} />
      )}

      <PromptForm onSubmit={handleSubmit} />
    </div>
  );
}
