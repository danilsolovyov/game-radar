<script lang="ts">
    import { onMount, onDestroy } from "svelte";
    import { EventsOn } from "wailsjs/runtime/runtime";
    import { GetTheme } from "wailsjs/go/app/App";
    import { models } from "wailsjs/go/models";

    import {
        createRadarState,
        resetRadarState,
        drawRadar,
        addBlip,
        clearSectionTimers,
        type RadarState,
    } from "./lib/radarCanvas";
    import { colorToRgbaWithOpacity } from "./utils/color";

    let canvas: HTMLCanvasElement;
    let ctx: CanvasRenderingContext2D | null = null;
    let animationFrame = 0;
    let radarState: RadarState = createRadarState();

    function loop() {
        if (ctx && canvas) {
            drawRadar(ctx, canvas, radarState);
        }
        animationFrame = requestAnimationFrame(loop);
    }

    onMount(async () => {
        ctx = canvas.getContext("2d");
        const theme = await GetTheme();
        radarState.theme = theme;
        resetRadarState(radarState, theme);
        animationFrame = requestAnimationFrame(loop);

        EventsOn("radar-data", (data: models.Blip) => {
            addBlip(radarState, data);
        });
    });

    onDestroy(() => {
        cancelAnimationFrame(animationFrame);
        clearSectionTimers(radarState);
        radarState.blips = [];
        if (ctx) {
            ctx.clearRect(0, 0, canvas.width, canvas.height);
        }
        if (canvas) {
            canvas.width = 0;
            canvas.height = 0;
        }
    });
</script>

<div class="radar-container">
    <canvas
        bind:this={canvas}
        style={`box-shadow: inset 0 0 24px ${colorToRgbaWithOpacity(radarState.theme?.radar_color, 0.2)};`}
    ></canvas>
</div>

<style>
    .radar-container {
        display: inline-block;
    }

    canvas {
        display: block;
        border-radius: 9999px;
    }
</style>
