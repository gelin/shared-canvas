<script lang="ts">
    import { PaletteTool, tool, STAMPS, STAMP_SIZE, SIZES, stampUrl } from "./Palette";

    const selectTool = (e: MouseEvent) => {
        const target = e.currentTarget as HTMLElement;
        const type = target.dataset.type as PaletteTool['type'];
        const color = target.dataset.color as PaletteTool['color'];
        const size = parseInt(target.dataset.size as string);
        const stamp = target.dataset.stamp as PaletteTool['stamp'];

        switch (type) {
            case 'line':
                $tool = new PaletteTool('line', color, size);
                break;
            case 'stamp':
                $tool = new PaletteTool('stamp', color, STAMP_SIZE, stamp);
                break;
        }
    }
</script>

<div class="palette">
    {#each SIZES as size}
        <div class="pair">
            <button title="Black line {size}px width"
                    class="black {$tool.type === 'line' && $tool.size === size && $tool.color === 'black' ? 'active' : ''}"
                    data-type="line"
                    data-color="black"
                    data-size="{size}"
                    onclick={selectTool}
            ><span></span></button>
            <button title="White line {size}px width"
                    class="white {$tool.type === 'line' && $tool.size === size && $tool.color === 'white' ? 'active' : ''}"
                    data-type="line"
                    data-color="white"
                    data-size="{size}"
                    onclick={selectTool}
            ><span></span></button>
        </div>
    {/each}
    {#each STAMPS as stamp}
        <div class="pair">
            <button title="Black {stamp}"
                    class="black {$tool.stamp === stamp && $tool.color === 'black' ? 'active' : ''}"
                    data-type="stamp"
                    data-color="black"
                    data-stamp={stamp}
                    data-size={STAMP_SIZE}
                    onclick={selectTool}
            ><img src={stampUrl('black',stamp)} alt={stamp}></button>
            <button title="White {stamp}"
                    class="white {$tool.stamp === stamp && $tool.color === 'white' ? 'active' : ''}"
                    data-type="stamp"
                    data-color="white"
                    data-stamp={stamp}
                    data-size={STAMP_SIZE}
                    onclick={selectTool}
            ><img src={stampUrl('white',stamp)} alt={stamp}></button>
        </div>
    {/each}
</div>

<style>
    .palette {
        display: flex;
        flex-direction: column;
        flex-wrap: wrap;
        gap: clamp(0.2rem, 0.8vh, 1rem);
    }

    .palette .pair {
        display: flex;
        flex-direction: row;
        flex-wrap: wrap;
        gap: clamp(0.2rem, 0.5vw, 1rem);
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

    .palette button img {
        width: 25px;
    }

    .palette button.black {
        background: color-mix(in oklab, canvas, canvasText 20%);
    }

    .palette button.black span {
        background: black;
    }

    .palette button.black:hover {
        background: color-mix(in oklab, canvas, canvasText 40%);
    }

    .palette button.black.active {
        border: 2px solid canvasText;
    }

    .palette button.white {
        background: color-mix(in oklab, canvas, canvasText 20%);
    }

    .palette button.white span {
        background: white;
    }

    .palette button.white:hover {
        background: color-mix(in oklab, canvas, canvasText 40%);
    }

    .palette button.white.active {
        border: 2px solid canvasText;
    }

    @media (max-width: 800px) {
        .palette {
            flex-direction: row;
            gap: clamp(0.2rem, 0.8vw, 1rem);
        }

        .palette .pair {
            flex-direction: column;
            gap: clamp(0.1rem, 0.5vh, 1rem);
        }
    }
</style>
