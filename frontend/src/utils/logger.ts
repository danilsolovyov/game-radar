/**
 * Logging utilities wrapping Wails runtime logging
 * @module utils/logger
 */

import {
    LogError,
    LogInfo,
    LogWarning,
} from "../../wailsjs/runtime/runtime";

/**
 * Safely logs an error to Wails runtime and console
 * @param message - Error message
 * @param error - Optional error object or additional data
 */
export function safeLogError(message: string, error?: unknown): void {
    const eventLike =
        typeof error === "object" &&
        error !== null &&
        "type" in error &&
        "target" in error;

    // Log to Wails runtime
    try {
        LogError(message);
    } catch (err) {
        console.error("[Logger] Failed to log via Wails:", err);
    }

    // Log detailed error if provided
    if (error instanceof Error) {
        console.error(message, error);
        try {
            LogError(`${message}: ${error.message}`);
        } catch (err) {
            console.error("[Logger] Failed to log details:", err);
        }
    } else if (error !== undefined) {
        console.error(message, error);
        // For Event-like objects, do not duplicate the error in Wails logs
        // as "[object Event]": that is noise without useful context.
        if (!eventLike) {
            try {
                LogError(`${message}: ${String(error)}`);
            } catch (err) {
                console.error("[Logger] Failed to log details:", err);
            }
        }
    } else {
        console.error(message);
    }
}

/**
 * Logs a warning message
 * @param message - Warning message
 */
export function safeLogWarn(message: string): void {
    console.warn(message);
    try {
        LogWarning(`WARN: ${message}`);
    } catch (err) {
        console.error("[Logger] Failed to log warning via Wails:", err);
    }
}

/**
 * Logs an informational message
 * @param message - Info message
 */
export function safeLogInfo(message: string): void {
    console.log(message);
    try {
        LogInfo(`INFO: ${message}`);
    } catch (err) {
        console.error("[Logger] Failed to log info via Wails:", err);
    }
}
