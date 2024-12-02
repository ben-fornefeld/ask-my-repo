import Chat from "../components/chat";
import type { MetaFunction } from "@remix-run/node";

export const meta: MetaFunction = () => {
  return [
    { title: "Ask my Repo" },
    {
      name: "description",
      content: "Interface for a go ranking context retrival service",
    },
  ];
};

export default function Index() {
  return (
    <div className="flex min-h-screen items-center justify-center py-24">
      <main>
        <Chat />
      </main>
    </div>
  );
}
