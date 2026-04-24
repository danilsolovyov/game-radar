/**
 * Color conversion utilities
 * @module utils/color
 */

import { color, models } from "../../wailsjs/go/models";

/**
 * Converts a color object to CSS rgba string
 * @param colorValue - Color object with R, G, B, A properties
 * @returns CSS rgba string
 */
export function colorToRgba(
    colorValue: color.RGBA | string | null | undefined,
): string {
    if (!colorValue) return "rgba(0, 0, 0, 0)";

    if (typeof colorValue === "string") return colorValue;

    const r = colorValue.R ?? 0;
    const g = colorValue.G ?? 0;
    const b = colorValue.B ?? 0;
    const a = colorValue.A !== undefined ? colorValue.A / 255 : 1;

    return `rgba(${r}, ${g}, ${b}, ${a})`;
}

/**
 * Converts a color object to CSS rgba string with custom opacity
 * @param colorValue - Color object with R, G, B, A properties
 * @param opacity - Opacity value (0-1)
 * @returns CSS rgba string with custom opacity
 */
export function colorToRgbaWithOpacity(
    colorValue: color.RGBA | string | null | undefined,
    opacity: number,
): string {
    if (!colorValue) return `rgba(0, 255, 180, ${opacity})`;
    if (typeof colorValue === "string") return colorValue;

    const r = colorValue.R ?? 0;
    const g = colorValue.G ?? 0;
    const b = colorValue.B ?? 0;
    return `rgba(${r}, ${g}, ${b}, ${opacity})`;
}

/**
 * Safely parses a numeric value
 * @param value - Value to parse
 * @param fallback - Fallback value if parsing fails
 * @returns Parsed number or fallback
 */
export function numberValue(value: unknown, fallback: number): number {
    if (typeof value === "number" && Number.isFinite(value)) return value;
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : fallback;
}

/**
 * Normalizes a color object with fallback values
 * @param input - Input color object
 * @param fallback - Fallback RGBA values
 * @returns Normalized RGBA object
 */
export function normalizeColor(
    input: unknown,
    fallback: color.RGBA,
): color.RGBA {
    const source = input as color.RGBA;
    return new color.RGBA({
        R: numberValue(source?.R, fallback.R),
        G: numberValue(source?.G, fallback.G),
        B: numberValue(source?.B, fallback.B),
        A: numberValue(source?.A, fallback.A),
    });
}

/**
 * Clamps a byte value to 0-255 range
 * @param value - Value to clamp
 * @returns Clamped byte value
 */
export function clampByte(value: number): number {
    return Math.max(0, Math.min(255, Math.round(numberValue(value, 0))));
}

/**
 * Converts a byte to hex string with padding
 * @param value - Byte value (0-255)
 * @returns Hex string (e.g., "ff")
 */
export function toHexByte(value: number): string {
    return clampByte(value).toString(16).padStart(2, "0");
}

/**
 * Converts RGBA to hex color string (without alpha)
 * @param colorValue - RGBA color object
 * @returns Hex color string (e.g., "#00ffb4")
 */
export function rgbaToHex(colorValue: color.RGBA): string {
    return `#${toHexByte(colorValue.R)}${toHexByte(colorValue.G)}${toHexByte(colorValue.B)}`;
}

/**
 * Parses a hex color string to RGB object
 * @param hex - Hex color string (with or without #)
 * @returns RGB object or null if invalid
 */
export function hexToRgb(
    hex: string,
): { r: number; g: number; b: number } | null {
    const normalized = hex.trim().replace("#", "");
    if (!/^[0-9a-fA-F]{6}$/.test(normalized)) {
        return null;
    }

    return {
        r: parseInt(normalized.slice(0, 2), 16),
        g: parseInt(normalized.slice(2, 4), 16),
        b: parseInt(normalized.slice(4, 6), 16),
    };
}

/**
 * Applies hex color to RGBA target object
 * @param target - Target RGBA object to modify
 * @param hex - Hex color string to apply
 */
export function applyHexToColor(target: color.RGBA, hex: string): void {
    const rgb = hexToRgb(hex);
    if (!rgb) return;
    target.R = rgb.r;
    target.G = rgb.g;
    target.B = rgb.b;
}

/**
 * Default background color for radar
 */
export const DEFAULT_BACKGROUND_COLOR = new color.RGBA({
    R: 0,
    G: 20,
    B: 10,
    A: 50,
});

/**
 * Default radar color (neon green)
 */
export const DEFAULT_RADAR_COLOR = new color.RGBA({
    R: 0,
    G: 255,
    B: 180,
    A: 255,
});

/**
 * Creates a deep clone of a Theme object
 * @param theme - Theme to clone
 * @returns Cloned theme
 */
export function cloneTheme(theme: models.Theme): models.Theme {
    return new models.Theme({
        name: String(theme?.name ?? ""),
        background_color: normalizeColor(theme?.background_color, DEFAULT_BACKGROUND_COLOR),
        radar_color: normalizeColor(theme?.radar_color, DEFAULT_RADAR_COLOR),
        border_opacity: numberValue(theme?.border_opacity, 0.1),
        border_width: numberValue(theme?.border_width, 2),
        section_base_opacity: numberValue(theme?.section_base_opacity, 0),
        section_bright_opacity: numberValue(theme?.section_bright_opacity, 0.9),
        section_timeout: numberValue(theme?.section_timeout, 500),
        section_count: numberValue(theme?.section_count, 25),
        ring_count: numberValue(theme?.ring_count, 3),
        show_blips: theme?.show_blips ?? true,
        blip_opacity: numberValue(theme?.blip_opacity, 0.5),
        blip_timeout: numberValue(theme?.blip_timeout, 500),
        blip_size: numberValue(theme?.blip_size, 3),
        size: numberValue(theme?.size, 320),
        pos_x: numberValue(theme?.pos_x, 30),
        pos_y: numberValue(theme?.pos_y, 30),
        intensity_multiplier: numberValue(theme?.intensity_multiplier, 1),
    });
}
