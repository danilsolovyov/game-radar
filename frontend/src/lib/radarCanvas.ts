/**
 * Radar canvas rendering module
 * @module lib/radarCanvas
 */

import { colorToRgba, colorToRgbaWithOpacity } from "../utils/color";
import type { models } from "../../wailsjs/go/models";

/**
 * Radar state for canvas rendering
 */
export interface RadarState {
    theme: models.Theme | null;
    sectionBrightness: number[];
    sectionTimers: (ReturnType<typeof setTimeout> | null)[];
    blips: RadarBlip[];
}

/**
 * Blip for radar display
 */
export interface RadarBlip {
    angle: number;
    distance: number;
    intensity: number;
}

/**
 * Creates initial radar state
 */
export function createRadarState(): RadarState {
    return {
        theme: null,
        sectionBrightness: [],
        sectionTimers: [],
        blips: [],
    };
}

/**
 * Resets radar state for a new theme
 */
export function resetRadarState(state: RadarState, theme: models.Theme): void {
    state.theme = theme;
    const sectionCount = theme.section_count || 25;
    state.sectionBrightness = Array(sectionCount).fill(0);
    state.sectionTimers = Array(sectionCount).fill(null);
    state.blips = [];
}

/**
 * Gets section index from angle in degrees
 */
export function getSectionIndex(angleDegrees: number, sectionSizeDegrees: number, sectionCount: number): number {
    const normalizedDegrees = ((angleDegrees % 360) + 360) % 360;
    const offsetDegrees = (normalizedDegrees + sectionSizeDegrees / 2) % 360;
    return Math.floor(offsetDegrees / sectionSizeDegrees) % sectionCount;
}

/**
 * Triggers section brightness
 */
export function triggerSectionBrightness(
    state: RadarState,
    sectionIndex: number,
    intensity: number,
): void {
    const theme = state.theme;
    if (!theme) return;

    if (state.sectionTimers[sectionIndex]) {
        clearTimeout(state.sectionTimers[sectionIndex]);
    }

    state.sectionBrightness[sectionIndex] = Math.min(1, intensity);

    const timeout = theme.section_timeout || 500;
    state.sectionTimers[sectionIndex] = setTimeout(() => {
        state.sectionBrightness[sectionIndex] = 0;
        state.sectionTimers[sectionIndex] = null;
    }, timeout);
}

/**
 * Clears all section timers
 */
export function clearSectionTimers(state: RadarState): void {
    state.sectionTimers.forEach((timer) => {
        if (timer) clearTimeout(timer);
    });
    state.sectionTimers = [];
}

/**
 * Main draw function for radar canvas
 */
export function drawRadar(
    ctx: CanvasRenderingContext2D,
    canvas: HTMLCanvasElement,
    state: RadarState,
): void {
    const theme = state.theme;
    if (!ctx || !theme) return;

    const dpr = window.devicePixelRatio || 1;
    const w = Math.max(1, Math.floor(theme.size * dpr));
    const h = Math.max(1, Math.floor(theme.size * dpr));

    if (canvas.width !== w || canvas.height !== h) {
        canvas.width = w;
        canvas.height = h;
        canvas.style.width = theme.size + "px";
        canvas.style.height = theme.size + "px";
    }

    const c = ctx;
    c.setTransform(1, 0, 0, 1, 0, 0);
    c.clearRect(0, 0, w, h);

    const cx = w / 2;
    const cy = h / 2;
    const r = Math.min(cx, cy);

    // Background
    c.fillStyle = colorToRgba(theme.background_color);
    c.fillRect(0, 0, w, h);

    // Concentric rings
    const radarColorRgba = colorToRgba(theme.radar_color);
    c.strokeStyle = radarColorRgba.replace("1)", `${theme.border_opacity})`);
    c.lineWidth = theme.border_width * dpr;

    for (let i = 0; i <= theme.ring_count; i++) {
        c.beginPath();
        c.arc(cx, cy, (r * i) / theme.ring_count, 0, Math.PI * 2);
        c.stroke();
    }

    // Section size
    const sectionSizeDegrees = theme.section_count ? 360 / theme.section_count : 14.4;

    // Draw sectors
    for (let i = 0; i < theme.section_count; i++) {
        const startAngle = ((i * sectionSizeDegrees - sectionSizeDegrees / 2 - 90) * Math.PI) / 180;
        const endAngle = (((i + 1) * sectionSizeDegrees - sectionSizeDegrees / 2 - 90) * Math.PI) / 180;
        const brightness = state.sectionBrightness[i] || 0;

        const brightOpacity = theme.section_base_opacity + brightness * theme.section_bright_opacity;

        c.fillStyle = radarColorRgba.replace("1)", `${brightOpacity})`);
        c.beginPath();
        c.moveTo(cx, cy);
        c.arc(cx, cy, r, startAngle, endAngle);
        c.closePath();
        c.fill();

        c.strokeStyle = radarColorRgba.replace("1)", `${theme.border_opacity + brightness * 0.5})`);
        c.lineWidth = theme.border_width * dpr;

        c.beginPath();
        c.moveTo(cx, cy);
        c.lineTo(cx + Math.cos(startAngle) * r, cy + Math.sin(startAngle) * r);
        c.stroke();

        c.beginPath();
        c.moveTo(cx, cy);
        c.lineTo(cx + Math.cos(endAngle) * r, cy + Math.sin(endAngle) * r);
        c.stroke();
    }

    // Process blips - always trigger section brightness, but only draw dots if enabled
    const drawBlips = theme.show_blips && theme.blip_size > 0 && theme.blip_opacity > 0;
    
    for (const b of state.blips) {
        const sectionIndex = getSectionIndex(b.angle, sectionSizeDegrees, theme.section_count);
        const intensity = b.intensity || 0.5;
        triggerSectionBrightness(state, sectionIndex, intensity);

        // Only draw blip dots if enabled
        if (drawBlips) {
            const a = ((b.angle - 90) * Math.PI) / 180;
            const dist = Math.min(1, Math.max(0, b.distance));
            const br = theme.blip_size * dpr + (b.intensity || 0) * 2 * dpr;
            const x = cx + Math.cos(a) * r * dist;
            const y = cy + Math.sin(a) * r * dist;

            c.fillStyle = radarColorRgba.replace("1)", `${theme.blip_opacity})`);
            c.beginPath();
            c.arc(x, y, br, 0, Math.PI * 2);
            c.fill();
        }
    }

    // Radar border
    c.strokeStyle = radarColorRgba.replace("1)", `${theme.border_opacity})`);
    c.lineWidth = theme.border_width * dpr;
    c.beginPath();
    c.arc(cx, cy, r, 0, Math.PI * 2);
    c.stroke();
}

/**
 * Adds a blip to radar state
 */
export function addBlip(state: RadarState, data: models.Blip): void {
    if (state.theme) {
        const intensity_multiplier = state.theme.intensity_multiplier || 2;
        const baseDistance = data.distance === 0 ? data.intensity : data.distance;
        const scaledIntensity = data.intensity * intensity_multiplier;
        const blip: RadarBlip = {
            angle: data.angle,
            distance: baseDistance * intensity_multiplier,
            intensity: scaledIntensity,
        };
        state.blips.push(blip);

        const timeout = state.theme.blip_timeout || 500;
        setTimeout(() => {
            state.blips = state.blips.filter(
                (b) => b.angle !== blip.angle || b.distance !== blip.distance || b.intensity !== blip.intensity,
            );
        }, timeout);
    }
}
