import clsx from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: any[]) {
  return twMerge(clsx(...inputs));
}

export function compareVersion(a: string, b: string) {
  const partsA = a.split(".").map(Number);
  const partsB = b.split(".").map(Number);

  const maxLen = Math.max(partsA.length, partsB.length);

  for (let i = 0; i < maxLen; i++) {
    const numA = partsA[i] || 0;
    const numB = partsB[i] || 0;

    if (numA > numB) return 1;
    if (numA < numB) return -1;
  }
  return 0;
}
