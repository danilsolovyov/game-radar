// i18n store for Svelte 5
import { writable, derived, get } from "svelte/store";
import { LOCALES, DEFAULT_LOCALE, type Locale, type LocaleOption } from "./locales";
import { TRANSLATIONS, type Translations } from "./translations";
import { GetLanguage, SetLanguage } from "wailsjs/go/app/App";
import { EventsOn } from "../../../wailsjs/runtime/runtime";

// Storage key for persisted locale
const LOCALE_STORAGE_KEY = "sound-radar-locale";

// Check if we're in a browser environment
const isBrowser = typeof window !== "undefined";

type WailsWindow = Window & {
    go?: {
        app?: {
            App?: unknown;
        };
    };
};

function hasWailsBridge(): boolean {
    if (!isBrowser) return false;
    const w = window as WailsWindow;
    return Boolean(w.go?.app?.App);
}

// Get initial locale from storage or browser settings or default
function getInitialLocale(): Locale {
    if (isBrowser) {
        // Try to get from localStorage first
        const stored = localStorage.getItem(LOCALE_STORAGE_KEY);
        if (stored && LOCALES.some((l) => l.code === stored)) {
            return stored as Locale;
        }

        // Try to detect from browser
        const browserLang = navigator.language.split("-")[0];
        if (LOCALES.some((l) => l.code === browserLang)) {
            return browserLang as Locale;
        }
    }

    return DEFAULT_LOCALE;
}

// Create the locale store
function createLocaleStore() {
    const initialLocale = getInitialLocale();
    const { subscribe, set, update } = writable<Locale>(initialLocale);

    return {
        subscribe,
        set: (locale: Locale) => {
            if (isBrowser) {
                localStorage.setItem(LOCALE_STORAGE_KEY, locale);
            }
            set(locale);
        },
        toggle: () => {
            update((current) => {
                const currentIndex = LOCALES.findIndex((l) => l.code === current);
                const nextIndex = (currentIndex + 1) % LOCALES.length;
                const nextLocale = LOCALES[nextIndex].code;
                if (isBrowser) {
                    localStorage.setItem(LOCALE_STORAGE_KEY, nextLocale);
                }
                return nextLocale;
            });
        },
    };
}

export const locale = createLocaleStore();

// Derived store for translations - named 't' so it can be used as $t in templates
export const t = derived(
    locale,
    ($locale) => TRANSLATIONS[$locale]
);

// Helper function to get current translations
export function getTranslations(): Translations {
    return TRANSLATIONS[get(locale)];
}

// Helper function to get translation by key
export function translate(key: keyof Translations): string {
    const currentLocale = get(locale);
    return TRANSLATIONS[currentLocale][key];
}

// Derived store for available locales
export const availableLocales = derived(
    locale,
    ($locale): LocaleOption[] => LOCALES
);

// Get current locale option
export function getCurrentLocaleOption(): LocaleOption {
    const currentLocale = get(locale);
    return LOCALES.find((l) => l.code === currentLocale) ?? LOCALES[0];
}

// Locale change handler for Wails integration
export async function changeLocale(newLocale: Locale): Promise<void> {
    locale.set(newLocale);

    // If running in browser, notify the backend about locale change
    if (isBrowser) {
        localStorage.setItem(LOCALE_STORAGE_KEY, newLocale);
    }
    if (hasWailsBridge()) {
        await SetLanguage(newLocale);
    }
}

export async function initLocaleSync(): Promise<(() => void) | undefined> {
    if (!hasWailsBridge()) return undefined;

    try {
        const backendLang = (await GetLanguage()).split("-")[0]?.toLowerCase();
        if (backendLang && LOCALES.some((l) => l.code === backendLang)) {
            locale.set(backendLang as Locale);
        }
    } catch {
        // Ignore backend sync errors; fallback locale logic is already applied.
    }

    return EventsOn("language-changed", (lang: string) => {
        const normalized = String(lang ?? "").split("-")[0]?.toLowerCase();
        if (normalized && LOCALES.some((l) => l.code === normalized)) {
            locale.set(normalized as Locale);
        }
    });
}

// Re-export types
export type { Translations };
