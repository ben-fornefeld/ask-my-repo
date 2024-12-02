import Card from "../components/ui/card";

export default function Navbar() {
  return (
    <nav className="fixed bottom-4 left-4 z-50">
      <Card className="bg-white mx-auto px-8 py-3">
        <h1 className="text-xl whitespace-nowrap">Ask my Repo</h1>
      </Card>
    </nav>
  );
}
