<script lang="ts">
    import {
        GetAudioDevices,
        GetSelectedDevice,
        SetDevice,
    } from "wailsjs/go/app/App";
    import { models } from "wailsjs/go/models";
    import { onMount } from "svelte";
    import { SvelteMap } from "svelte/reactivity";
    import { safeLogError, safeLogInfo } from "./utils/logger";
    import { t } from "$lib/i18n";

    import Card from "$lib/components/ui/card/card.svelte";
    import CardContent from "$lib/components/ui/card/card-content.svelte";
    import CardHeader from "$lib/components/ui/card/card-header.svelte";
    import CardTitle from "$lib/components/ui/card/card-title.svelte";
    import Table from "$lib/components/ui/table/table.svelte";
    import TableBody from "$lib/components/ui/table/table-body.svelte";
    import TableRow from "$lib/components/ui/table/table-row.svelte";
    import TableCell from "$lib/components/ui/table/table-cell.svelte";
    import Button from "$lib/components/ui/button/button.svelte";
    import Badge from "$lib/components/ui/badge/badge.svelte";
    import ChevronUp from "@lucide/svelte/icons/chevron-up";
    import ChevronDown from "@lucide/svelte/icons/chevron-down";
    import RefreshCcw from "@lucide/svelte/icons/refresh-ccw";

    import DeviceSelector from "./components/DeviceSelector.svelte";
    import AudioPlayer from "./AudioPlayer.svelte";

    let devices: models.DeviceSpeakers[] = $state([]);
    let deviceMap = new SvelteMap<string, models.DeviceSpeakers>();
    let error = $state("");
    let loading = $state(false);
    let collapsed = $state(true);
    let selectedDeviceId = $state("");

    function getDeviceById(id: string): models.DeviceSpeakers | undefined {
        return deviceMap.get(id);
    }

    const selectedDevice = $derived.by(() => {
        devices;
        return getDeviceById(selectedDeviceId);
    });

    const infoRows = $derived.by(() => {
        const device = selectedDevice;
        return [
            {
                label: $t.pcm,
                value: device?.format_pcm ? `${device.format_pcm} bit` : "—",
            },
            { label: $t.rate, value: device?.rate ? `${device.rate} Hz` : "—" },
            { label: $t.channels, value: device?.channels ?? "—" },
            {
                label: $t.defaultPeriod,
                value: device?.default_period
                    ? `${device.default_period} ns`
                    : "—",
            },
            {
                label: $t.minimumPeriod,
                value: device?.minimum_period
                    ? `${device.minimum_period} ns`
                    : "—",
            },
            { label: $t.latency, value: device?.latency ?? "—" },
        ];
    });

    async function handleDeviceSelect(deviceId: string) {
        selectedDeviceId = deviceId;
        const device = getDeviceById(deviceId);
        if (device) {
            await SetDevice(device);
            safeLogInfo(`Selected device: ${deviceId}`);
        } else {
            safeLogError(`Device with id ${deviceId} not found`);
        }
    }

    async function loadDevices() {
        if (loading) return;
        loading = true;
        error = "";

        try {
            let backendSelectedId = "";
            try {
                const backendSelected = await GetSelectedDevice();
                backendSelectedId = backendSelected?.id || "";
            } catch (getSelectedErr) {
                safeLogError(
                    "Failed to get backend selected device",
                    getSelectedErr,
                );
            }

            const result = await GetAudioDevices();
            devices = Array.isArray(result)
                ? result.map((item) => new models.DeviceSpeakers(item))
                : [];

            deviceMap.clear();
            devices.forEach((device) => {
                deviceMap.set(device.id, device);
            });

            const currentSelectedId =
                selectedDeviceId && deviceMap.has(selectedDeviceId)
                    ? selectedDeviceId
                    : "";
            const backendSelectedValidId =
                backendSelectedId && deviceMap.has(backendSelectedId)
                    ? backendSelectedId
                    : "";
            const nextSelectedId =
                currentSelectedId ||
                backendSelectedValidId ||
                devices.find((device) => device.is_default)?.id ||
                devices[0]?.id ||
                "";

            selectedDeviceId = nextSelectedId;

            if (nextSelectedId && nextSelectedId !== backendSelectedId) {
                try {
                    const device = getDeviceById(nextSelectedId);
                    if (device) {
                        await SetDevice(device);
                    }
                } catch (setErr) {
                    safeLogError(
                        "Failed to set active device after load",
                        setErr,
                    );
                }
            }
        } catch (err) {
            error = err instanceof Error ? err.message : String(err);
            safeLogError("Failed to load audio devices", err);
        } finally {
            loading = false;
        }
    }

    function toggleCollapseInfo() {
        collapsed = !collapsed;
    }

    onMount(() => {
        loadDevices();
    });
</script>

<Card class="w-full min-w-0 min-h-0 mx-auto overflow-hidden border-2">
    <CardHeader class="space-y-2 pb-3">
        <div class="flex items-center justify-between">
            <div>
                <CardTitle>
                    <h2 class="text-lg font-semibold">
                        {$t.audioDevicesTitle}
                    </h2>
                </CardTitle>
                <p class="text-sm text-muted-foreground mt-0.5">
                    {$t.playbackDevice}
                </p>
            </div>
            {#if selectedDevice?.is_default}
                <Badge variant="secondary" class="font-medium"
                    >{$t.defaultDevice}</Badge
                >
            {/if}
        </div>
    </CardHeader>

    <CardContent class="space-y-4 min-h-0 overflow-auto">
        <div class="flex items-center gap-2 p-3 rounded-lg bg-muted/30 border">
            <div class="flex-1 min-w-0">
                <DeviceSelector
                    {devices}
                    {selectedDeviceId}
                    onDeviceSelect={handleDeviceSelect}
                />
            </div>
            <Button
                onclick={loadDevices}
                variant="outline"
                size="icon"
                class="shrink-0"
                title={$t.refreshDevices}
            >
                <RefreshCcw class="h-4 w-4" />
            </Button>
            <div class="shrink-0">
                <AudioPlayer />
            </div>
        </div>

        <div class="rounded-lg border bg-muted/20 p-3 space-y-2">
            <div class="flex items-center justify-between">
                <span
                    class="text-xs uppercase tracking-wide text-muted-foreground font-medium"
                >
                    {$t.deviceInfo}
                </span>
            </div>
            <p class="text-sm font-semibold truncate">
                {selectedDevice?.name || $t.noDevicesFound}
            </p>
            <p class="text-xs text-muted-foreground">
                {selectedDevice?.channels
                    ? `${selectedDevice.channels} ch • ${selectedDevice.rate} Hz`
                    : $t.noDevicesFound}
            </p>
        </div>

        <Button
            variant="outline"
            onclick={toggleCollapseInfo}
            class="w-full justify-center"
        >
            {#if collapsed}
                <ChevronDown class="h-4 w-4 mr-1" />{$t.expandInfo}
            {:else}
                <ChevronUp class="h-4 w-4 mr-1" />{$t.collapseInfo}
            {/if}
        </Button>

        {#if !collapsed}
            <div class="rounded-lg border bg-muted/10 p-3">
                <Table>
                    <TableBody>
                        {#each infoRows as row (row.label)}
                            <TableRow class="hover:bg-muted/30">
                                <TableCell class="text-muted-foreground text-xs"
                                    >{row.label}</TableCell
                                >
                                <TableCell class="font-medium text-sm"
                                    >{row.value}</TableCell
                                >
                            </TableRow>
                        {/each}
                    </TableBody>
                </Table>
            </div>
        {/if}
    </CardContent>
</Card>
