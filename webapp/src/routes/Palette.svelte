<script lang="ts">
    import { onMount } from 'svelte';
    import { type PaletteTool, PaletteChangeEvent, DEFAULT_TOOL } from "./Palette";

    const sizes = [ 1, 3, 5, 7, 10, 15 ]

    let { tool = DEFAULT_TOOL, onPaletteChange } = $props();

    const LS_TOOL_KEY = 'palette.tool';

    const persist = () => {
        try {
            localStorage.setItem(LS_TOOL_KEY, JSON.stringify(tool));
        } catch (_) {
            // ignore storage errors (e.g., privacy mode)
        }
    }

    onMount(() => {
        try {
            const savedTool = JSON.parse(localStorage.getItem(LS_TOOL_KEY) || '{}');
            console.log('Restored tool:', savedTool);
            const validColor = savedTool?.color === 'black' || savedTool?.color === 'white' ? savedTool?.color : null;
            const parsedSize = savedTool?.size ? parseInt(savedTool?.size, 10) : NaN;
            const validSize = sizes.includes(parsedSize) ? parsedSize : null;

            if (validColor && validSize !== null) {
                tool = {
                    type: 'line',
                    color: validColor,
                    size: validSize,
                    stamp: null
                }
            }

            // Notify parent about the (possibly restored) selection
            onPaletteChange?.(new PaletteChangeEvent(tool));
        } catch (_) {
            // If localStorage is unavailable, just emit current defaults
            onPaletteChange?.(new PaletteChangeEvent(tool));
        }
    });

    const selectTool = (e: MouseEvent) => {
        const target = e.currentTarget as HTMLElement;
        const type = target.dataset.type as PaletteTool['type'];
        const color = target.dataset.color as PaletteTool['color'];
        const size = parseInt(target.dataset.size as string);
        if (type === 'line') {
            tool = {
                type: 'line',
                color: color,
                size: size,
                stamp: null
            }
        }
        persist();
        onPaletteChange?.(new PaletteChangeEvent(tool));
    }
</script>

<div class="palette">
    <div class="blacks {tool?.color === 'black' ? 'active' : ''}">
        {#each sizes as s}
            <button title="Line of {s} pixels"
                    class="{tool?.type === 'line' && s === tool?.size ? 'active' : ''}"
                    data-type="line"
                    data-color="black"
                    data-size="{s}"
                    onclick={selectTool}
            ><span></span></button>
        {/each}
    </div>
    <div class="whites {tool?.color === 'white' ? 'active' : ''}">
        {#each sizes as s}
            <button title="Line of {s} pixels"
                    class="{tool?.type === 'line' && s === tool?.size ? 'active' : ''}"
                    data-type="line"
                    data-color="white"
                    data-size="{s}"
                    onclick={selectTool}
            ><span></span></button>
        {/each}
    </div>
</div>

<style>
    .palette {
        display: flex;
        flex-direction: row;
        gap: 0.5rem;
    }
    .palette .blacks, .palette .whites {
        display: flex;
        flex-direction: column;
        gap: 1rem;
    }
    .palette button {
        flex-shrink: 0;
        padding: 2px;
        width: 2rem;
        height: 2rem;
        display: flex;
        justify-content: center;
        align-items: center;
    }
    .palette button span {
        display: block;
        border-radius: 50%;
    }
    .palette button[data-size="1"] span {
        width: 1px;
        height: 1px;
    }
    .palette button[data-size="3"] span {
        width: 3px;
        height: 3px;
    }
    .palette button[data-size="5"] span {
        width: 5px;
        height: 5px;
    }
    .palette button[data-size="7"] span {
        width: 7px;
        height: 7px;
    }
    .palette button[data-size="10"] span {
        width: 10px;
        height: 10px;
    }
    .palette button[data-size="15"] span {
        width: 15px;
        height: 15px;
    }
    .palette .blacks button {
        background: color-mix(in oklab, canvas, canvasText 20%);
    }
    .palette .blacks button span {
        background: black;
    }
    .palette .blacks button:hover {
        background: color-mix(in oklab, canvas, canvasText 40%);
    }
    .palette .blacks.active button.active {
        border: 1px solid canvasText;
    }
    .palette .whites button {
        background: color-mix(in oklab, canvas, canvasText 20%);
    }
    .palette .whites button span {
        background: white;
    }
    .palette .whites button:hover {
        background: color-mix(in oklab, canvas, canvasText 40%);
    }
    .palette .whites.active button.active {
        border: 1px solid canvasText;
    }

    @media (max-width: 800px) {
        .palette {
            flex-direction: column;
        }
        .palette .blacks, .palette .whites {
            flex-direction: row;
        }
    }
</style>
