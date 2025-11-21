<script lang="ts">
    import { onMount } from 'svelte';

    export let width = 384
    export let height = 384

    let canvas: HTMLCanvasElement;
    let context: CanvasRenderingContext2D | null;
    let isDrawing = false;
    let prev = { x: 0, y: 0 };
    let lineWidth = 3;

    onMount(() => {
        context = canvas.getContext('2d');
        if (!context) return;
        context.fillStyle = 'white';
        context.fillRect(0, 0, canvas.width, canvas.height);
        context.strokeStyle = 'black';
        context.lineWidth = lineWidth;
        context.lineCap = 'round';
    })

    const handleMove = (({ offsetX: x1, offsetY: y1, buttons }: MouseEvent) => {
        if (!context) return;
        if (buttons == 1) {
            if (isDrawing) {
                const { x, y } = prev;
                context.beginPath();
                context.moveTo(x, y);
                context.lineTo(x1, y1);
                context.closePath();
                context.stroke();
                prev = { x: x1, y: y1 };
            } else {
                isDrawing = true;
                prev = { x: x1, y: y1 };
            }
        } else {
            isDrawing = false;
        }
    });

    const handleEnd = () => {
        isDrawing = false;
    }
</script>

<canvas
        id="canvas"
        bind:this={canvas}
        {width}
        {height}
        onmousemove={handleMove}
        onmouseleave={handleEnd}
></canvas>

