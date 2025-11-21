<script lang="ts">
    import { onMount } from 'svelte';

    export let width = 384
    export let height = 384

    let mainCanvas: HTMLCanvasElement;
    let mainContext: CanvasRenderingContext2D | null;

    let drawCanvas: HTMLCanvasElement;
    let drawContext: CanvasRenderingContext2D | null;

    let isDrawing = false;
    let prev = { x: 0, y: 0 };
    let lineWidth = 3;

    onMount(() => {
        mainContext = mainCanvas.getContext('2d');
        if (!mainContext) return;
        mainContext.fillStyle = 'white';
        mainContext.fillRect(0, 0, mainCanvas.width, mainCanvas.height);

        drawContext = drawCanvas.getContext('2d');
        if (!drawContext) return;
        drawContext.fillStyle = 'white';
        drawContext.fillRect(0, 0, drawCanvas.width, drawCanvas.height);
        drawContext.strokeStyle = 'black';
        drawContext.lineWidth = lineWidth;
        drawContext.lineCap = 'round';
    })

    const handleMove = (({ offsetX: x1, offsetY: y1, buttons }: MouseEvent) => {
        if (!drawContext) return;
        if (buttons == 1) {
            if (isDrawing) {
                const { x, y } = prev;
                drawContext.beginPath();
                drawContext.moveTo(x, y);
                drawContext.lineTo(x1, y1);
                drawContext.closePath();
                drawContext.stroke();

                const drawX = Math.min(x, x1) - lineWidth / 2;
                const drawY = Math.min(y, y1) - lineWidth / 2;
                const drawWidth = Math.abs(x - x1) + lineWidth;
                const drawHeight = Math.abs(y - y1) + lineWidth;
                const imageData = drawContext.getImageData(drawX, drawY, drawWidth, drawHeight);
                onDraw({ x: drawX, y: drawY, data: imageData });

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

    type DrawEvent = {
        x: number,
        y: number,
        data: ImageData
    }

    const onDraw = async (e: DrawEvent) => {
        if (!mainContext) return;
        const image = await createImageBitmap(e.data);
        mainContext.drawImage(image, e.x, e.y);
        // if (!drawContext) return;
        // drawContext.clearRect(0, 0, drawCanvas.width, drawCanvas.height);
    }
</script>

<canvas id="drawCanvas"
        bind:this={drawCanvas}
        {width}
        {height}
        onmousemove={handleMove}
        onmouseleave={handleEnd}
></canvas>
<canvas
        id="mainCanvas"
        bind:this={mainCanvas}
        {width}
        {height}
></canvas>

<style>
    #drawCanvas {
        border: 1px solid blue;
        /*position: absolute;*/
    }
</style>
