import testSound from "./assets/audio/test_sound.wav";
import { GetSelectedDevice } from "wailsjs/go/app/App";
import { safeLogError } from "./utils/logger";

type AudioPreviewState = {
    paused: boolean;
    ready: boolean;
};

type Listener = (state: AudioPreviewState) => void;

const state: AudioPreviewState = {
    paused: true,
    ready: false
};

const listeners = new Set<Listener>();

let audio: HTMLAudioElement | null = null;
let audioUrl: string | null = null;
let initialized = false;
let loading = false;
let lastSinkId = "";

function emit(): void {
    const snapshot = { ...state };
    for (const listener of listeners) {
        listener(snapshot);
    }
}

function setState(patch: Partial<AudioPreviewState>): void {
    Object.assign(state, patch);
    emit();
}

function normalizeName(value: string): string {
    return value.trim().toLowerCase();
}

function toSinkAudio(element: HTMLAudioElement): HTMLAudioElement & {
    setSinkId?: (sinkId: string) => Promise<void>;
    sinkId?: string;
} {
    return element as HTMLAudioElement & {
        setSinkId?: (sinkId: string) => Promise<void>;
        sinkId?: string;
    };
}

async function listAudioOutputs(): Promise<MediaDeviceInfo[]> {
    return navigator.mediaDevices.enumerateDevices().then((devices) => {
        return devices.filter((d) => d.kind === "audiooutput");
    });
}

function hasUsableOutputDevices(outputs: MediaDeviceInfo[]): boolean {
    return outputs.some((d) => Boolean(d.deviceId) && Boolean(d.label));
}

async function ensureMediaDeviceLabels(): Promise<MediaDeviceInfo[]> {
    let outputs = await listAudioOutputs();
    if (hasUsableOutputDevices(outputs)) {
        return outputs;
    }

    try {
        const stream = await navigator.mediaDevices.getUserMedia({ audio: true, video: false });
        stream.getTracks().forEach((t) => t.stop());
    } catch {}

    outputs = await listAudioOutputs();
    return outputs;
}

async function syncAudioOutputDevice(): Promise<void> {
    if (!audio) return;
    const sinkAudio = toSinkAudio(audio);
    if (typeof sinkAudio.setSinkId !== "function") return;

    try {
        const selected = await GetSelectedDevice();
        const targetName = normalizeName(String(selected?.name ?? ""));
        if (!targetName) return;

        const outputs = await ensureMediaDeviceLabels();
        if (outputs.length === 0 || !hasUsableOutputDevices(outputs)) return;

        const exact = outputs.find((d) => normalizeName(d.label) === targetName);
        const partial = outputs.find((d) => {
            const label = normalizeName(d.label);
            return label.includes(targetName) || targetName.includes(label);
        });
        const target = exact ?? partial;
        if (!target?.deviceId) return;
        if (lastSinkId === target.deviceId || sinkAudio.sinkId === target.deviceId) return;

        await sinkAudio.setSinkId(target.deviceId);
        lastSinkId = target.deviceId;
    } catch (err) {
        safeLogError("Failed to select audio output sink", err);
    }
}

export function subscribeAudioPreview(listener: Listener): () => void {
    listeners.add(listener);
    listener({ ...state });
    return () => {
        listeners.delete(listener);
    };
}

export function initAudioPreview(): void {
    if (initialized || loading) return;
    loading = true;

    audio = new Audio();

    audio.addEventListener("error", (err) => {
        safeLogError("Audio playback error", err);
        setState({ paused: true });
    });

    audio.addEventListener("ended", () => {
        if (!audio) return;
        audio.currentTime = 0;
        setState({ paused: true });
    });

    fetch(testSound)
        .then((response) => response.blob())
        .then((data) => {
            audioUrl = URL.createObjectURL(data);
            if (audio) {
                audio.src = audioUrl;
            }
            initialized = true;
            setState({ ready: true, paused: true });
            void syncAudioOutputDevice();
        })
        .catch((err) => {
            safeLogError("Failed to load audio", err);
        })
        .finally(() => {
            loading = false;
        });
}

export function toggleAudioPreview(): void {
    if (!audio) return;
    const currentAudio = audio;

    if (state.paused) {
        syncAudioOutputDevice()
            .catch((err) => {
                safeLogError("Failed to sync audio output sink", err);
            })
            .finally(() => {
                currentAudio
                    .play()
                    .then(() => {
                        setState({ paused: false });
                    })
                    .catch((err) => {
                        safeLogError("Failed to play audio", err);
                        setState({ paused: true });
                    });
            });
        return;
    }

    currentAudio.pause();
    currentAudio.currentTime = 0;
    setState({ paused: true });
}

export function disposeAudioPreview(): void {
    if (audio) {
        audio.pause();
        audio.src = "";
    }
    if (audioUrl) {
        URL.revokeObjectURL(audioUrl);
    }
    audio = null;
    audioUrl = null;
    initialized = false;
    loading = false;
    lastSinkId = "";
    setState({ paused: true, ready: false });
}
