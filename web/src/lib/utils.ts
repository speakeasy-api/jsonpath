import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function guessDocumentLanguage(content: string): DocumentLanguage {
  try {
    JSON.parse(content);
    return "json";
  } catch {
    return "yaml";
  }
}

export function formatDocument(doc: string, indentWidth: number = 2): string {
  const docLang = guessDocumentLanguage(doc);

  if (docLang === "json")
    return JSON.stringify(JSON.parse(doc), null, indentWidth);

  return doc;
}

export function arraysEqual<T>(a: T[], b: T[]): boolean {
  // Check if the arrays have the same length
  if (a.length !== b.length) {
    return false;
  }

  // Compare each element in the arrays
  for (let i = 0; i < a.length; i++) {
    if (a[i] !== b[i]) {
      return false;
    }
  }

  return true;
}
