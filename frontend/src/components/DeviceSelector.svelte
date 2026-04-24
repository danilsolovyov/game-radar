<script lang="ts">
    import {
        Select,
        SelectContent,
        SelectTrigger,
        SelectItem,
    } from "$lib/components/ui/select";
    import Badge from "$lib/components/ui/badge/badge.svelte";
    import type { models } from "wailsjs/go/models";

    interface Props {
        devices: models.DeviceSpeakers[];
        selectedDeviceId: string;
        onDeviceSelect: (deviceId: string) => void;
    }

    const { devices, selectedDeviceId, onDeviceSelect }: Props = $props();

    function handleSelect(value: string) {
        onDeviceSelect(value);
    }
</script>

<Select type="single" value={selectedDeviceId} onValueChange={handleSelect}>
    <SelectTrigger class="w-full min-w-0">
        <span class="block truncate">
            {devices.find((d) => d.id === selectedDeviceId)?.name ||
                "Select a device"}
        </span>
    </SelectTrigger>
    <SelectContent class="max-w-[min(80vw,36rem)]">
        {#each devices as device}
            <SelectItem value={device.id}>
                <div class="flex items-center gap-2 min-w-0 w-full">
                    <span class="truncate">{device.name}</span>
                    {#if device.is_default}
                        <Badge class="shrink-0">default</Badge>
                    {/if}
                </div>
            </SelectItem>
        {/each}
    </SelectContent>
</Select>
