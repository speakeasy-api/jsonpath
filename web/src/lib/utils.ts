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
