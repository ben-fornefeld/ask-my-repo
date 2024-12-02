import { cn } from "@/lib/utils";

export default function Card({
  children,
  className,
}: {
  children: React.ReactNode;
  className?: string;
}) {
  return (
    <div
      className={cn(
        "flex flex-col items-center gap-4 rounded-base border-2 border-border dark:border-darkBorder p-4 shadow-light dark:shadow-dark",
        className
      )}
    >
      {children}
    </div>
  );
}
