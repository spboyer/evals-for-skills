interface BadgeProps {
  variant: "success" | "danger" | "warning" | "neutral";
  children: string;
}

const variants: Record<BadgeProps["variant"], string> = {
  success:
    "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300",
  danger: "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-300",
  warning:
    "bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-300",
  neutral:
    "bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300",
};

export function Badge({ variant, children }: BadgeProps) {
  return (
    <span
      className={`inline-flex items-center rounded-full px-2 py-0.5 text-xs font-medium ${variants[variant]}`}
    >
      {children}
    </span>
  );
}
