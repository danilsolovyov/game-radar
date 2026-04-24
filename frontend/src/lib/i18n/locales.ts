// Locale definitions
export type Locale = "en" | "ru";

export interface LocaleOption {
    code: Locale;
    name: string;
    flag: string;
}

export const LOCALES: LocaleOption[] = [
    { code: "en", name: "English", flag: "🇺🇸" },
    { code: "ru", name: "Русский", flag: "🇷🇺" },
];

export const DEFAULT_LOCALE: Locale = "en";
