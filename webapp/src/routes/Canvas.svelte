<script lang="ts">
    import { onMount } from 'svelte';

    type Event = {
        offsetX: number;
        offsetY: number;
    }

    let canvas: HTMLCanvasElement;
    let context: CanvasRenderingContext2D | null;
    let isDrawing = false;
    let start = { x: 0, y: 0 };

    onMount(() => {
        context = canvas.getContext('2d');
        if (!context) return;
        context.fillStyle = 'white';
        context.fillRect(0, 0, canvas.width, canvas.height);
        context.strokeStyle = 'black';
        context.lineWidth = 3;  // TODO
        context.lineCap = 'round';
    })

    const handleStart = (({ offsetX: x, offsetY: y }: Event) => {
        isDrawing = true;
        start = { x, y };
    })

    const handleEnd = () => {
        isDrawing = false
    }
    const handleMove = (({ offsetX: x1, offsetY: y1 }: Event) => {
        if (!isDrawing) return;
        const { x, y } = start;
        if (!context) return;
        context.beginPath();
        context.moveTo(x, y);
        context.lineTo(x1, y1);
        context.closePath();
        context.stroke();

        start = { x: x1, y: y1 };
    })
</script>

<canvas
        id="canvas"
        bind:this={canvas}
        width="384"
        height="384"
        onmousedown={handleStart}
        onmouseup={handleEnd}
        onmouseleave={handleEnd}
        onmousemove={handleMove}
></canvas>

