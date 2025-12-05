<script lang="ts">
    import { onMount } from 'svelte';
    import { type PaletteTool, PaletteChangeEvent, DEFAULT_TOOL, stampUrl, STAMP_SIZE } from "./Palette";

    const sizes = [ 1, 3, 5, 7, 10, 15 ];
    const stamps = [ 'star' ];

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

            const parsedType = savedTool?.type as PaletteTool['type'] | null;
            const validType = parsedType === 'line' || parsedType === 'stamp' ? parsedType : null;
            const parsedColor = savedTool?.color as PaletteTool['color'] | null;
            const validColor = savedTool?.color === 'black' || savedTool?.color === 'white' ? parsedColor : null;
            const parsedSize = savedTool?.size ? parseInt(savedTool?.size, 10) : NaN;
            const validSize = sizes.includes(parsedSize) || parsedSize === STAMP_SIZE ? parsedSize : null;
            const parsedStamp = savedTool?.stamp as PaletteTool['stamp'] | null;
            const validStamp = parsedType === 'stamp' && stamps.includes(parsedStamp ?? '') ? parsedStamp : null;

            if (validType && validColor && validSize !== null) {
                switch (validType) {
                    case 'line': tool = { type: 'line', color: validColor, size: validSize, stamp: null }; break;
                    case 'stamp': tool = { type: 'stamp', color: validColor, size: 31, stamp: validStamp }; break;
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

        switch (type) {
            case 'line':
                tool = {
                    type: 'line',
                    color: color,
                    size: size,
                    stamp: null
                };
                break;
            case 'stamp':
                tool = {
                    type: 'stamp',
                    color: color,
                    size: STAMP_SIZE,
                    stamp: target.dataset.stamp as PaletteTool['stamp']
                };
                break;
        }

        persist();
        onPaletteChange?.(new PaletteChangeEvent(tool));
    }
</script>

<div class="palette">
    <div class="blacks {tool?.color === 'black' ? 'active' : ''}">
        {#each sizes as size}
            <button title="Line of {size} pixels"
                    class="{tool?.type === 'line' && tool?.size === size ? 'active' : ''}"
                    data-type="line"
                    data-color="black"
                    data-size="{size}"
                    onclick={selectTool}
            ><span></span></button>
        {/each}
        {#each stamps as stamp}
            <button title="Black {stamp}"
                    class="{tool?.stamp === stamp ? 'active' : ''}"
                    data-type="stamp"
                    data-color="black"
                    data-stamp={stamp}
                    data-size={STAMP_SIZE}
                    onclick={selectTool}
            ><img src={stampUrl('black', stamp)} alt={stamp}></button>
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
        {#each stamps as stamp}
            <button title="White {stamp}"
                    class="{tool?.stamp === stamp ? 'active' : ''}"
                    data-type="stamp"
                    data-color="white"
                    data-stamp={stamp}
                    data-size={STAMP_SIZE}
                    onclick={selectTool}
            ><img src={stampUrl('white', stamp)} alt={stamp}></button>
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
