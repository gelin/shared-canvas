<script lang="ts">
    import { onMount, tick } from 'svelte';
    import { fade } from 'svelte/transition';
    import { type DrawEvent, eventToMessage, messageToEvent } from "$lib/ws-canvas";
    import { wsClient, type WSDrawMessage } from "$lib/ws";

    const socket = wsClient;

    let width = $state(400);
    let height = $state(100);
    let ready = $state(false);

    let viewCanvas: HTMLCanvasElement;
    let viewContext: CanvasRenderingContext2D | null;

    let drawCanvas: HTMLCanvasElement;
    let drawContext: CanvasRenderingContext2D | null;

    let isDrawing = false;
    let prev = { x: 0, y: 0 };

    let { lineWidth = 3, color = 'black' } = $props();

    onMount(() => {
        viewContext = viewCanvas.getContext('2d', {
            alpha: false,
        });
        drawContext = drawCanvas.getContext('2d', {
            willReadFrequently: true,
        });
        initCanvas();

        socket.connect();
        socket.subscribe((message: WSDrawMessage) => {
            if (message.method === 'draw') {
                onDraw(messageToEvent(message));
            } else if (message.method === 'init') {
                onInit(messageToEvent(message).data);
            }
        });
    });

    const initCanvas = () => {
        if (!viewContext) return;
        // viewContext.imageSmoothingEnabled = false;
        viewContext.fillStyle = 'white';
        viewContext.fillRect(0, 0, viewCanvas.width, viewCanvas.height);

        if (!drawContext) return;
        // drawContext.imageSmoothingEnabled = false;
        drawContext.clearRect(0, 0, drawCanvas.width, drawCanvas.height);
        drawContext.strokeStyle = color;
        drawContext.lineWidth = lineWidth;
        drawContext.lineCap = 'round';
    };

    const handleClick = ({ offsetX: x1, offsetY: y1, buttons }: MouseEvent) => {
        if (!drawContext) return;
        drawContext.strokeStyle = color;
        drawContext.lineWidth = lineWidth;
        drawContext.beginPath();
        drawContext.moveTo(x1, y1);
        drawContext.lineTo(x1, y1);
        drawContext.stroke();

        sendDraw(x1, y1, x1, y1);
        prev = { x: x1, y: y1 };
    };

    const handleMove = ({ offsetX: x1, offsetY: y1, buttons }: MouseEvent) => {
        if (!drawContext) return;
        if (buttons == 1) {
            if (isDrawing) {
                const { x, y } = prev;
                drawContext.strokeStyle = color;
                drawContext.lineWidth = lineWidth;
                drawContext.beginPath();
                drawContext.moveTo(x, y);
                drawContext.lineTo(x1, y1);
                drawContext.stroke();

                sendDraw(x, y, x1, y1);

                prev = { x: x1, y: y1 };
            } else {
                isDrawing = true;
                prev = { x: x1, y: y1 };
            }
        } else {
            isDrawing = false;
        }
    };

    const handleTouchStart = (e: TouchEvent) => {
        const coords = toCanvasCoords(e);
        if (!coords) return;
        prev = coords;
        isDrawing = true;
    };

    const handleTouch = (e: TouchEvent) => {
        const coords = toCanvasCoords(e);
        if (!coords) return;
        if (!drawContext) return null;
        drawContext.strokeStyle = color;
        drawContext.lineWidth = lineWidth;
        drawContext.beginPath();
        drawContext.moveTo(prev.x, prev.y);
        drawContext.lineTo(coords.x, coords.y);
        drawContext.stroke();

        sendDraw(prev.x, prev.y, coords.x, coords.y);
        prev = coords;
    };

    const toCanvasCoords = ({ touches }: TouchEvent): { x: number; y: number } | null => {
        if (touches.length === 0) return null;
        if (!drawCanvas) return null;
        const rect = drawCanvas.getBoundingClientRect();
        return { x: touches[0].clientX - rect.left, y: touches[0].clientY - rect.top };
    };

    const handleEnd = () => {
        isDrawing = false;
    };

    const sendDraw = (x0: number, y0: number, x1: number, y1: number) => {
        if (!drawContext) return;
        const drawX = Math.floor(Math.min(x0, x1) - lineWidth);
        const drawY = Math.floor(Math.min(y0, y1) - lineWidth);
        const drawWidth = Math.ceil(Math.abs(x0 - x1) + 2 * lineWidth);
        const drawHeight = Math.ceil(Math.abs(y0 - y1) + 2 * lineWidth);
        const imageData = drawContext.getImageData(drawX, drawY, drawWidth, drawHeight);
        const message = eventToMessage({ x: drawX, y: drawY, data: imageData });
        socket.send(message);
    };

    const onDraw = async (e: DrawEvent) => {
        if (!viewContext) return;
        // need to convert to image for transparency and composition to work
        const image = await createImageBitmap(e.data);
        viewContext.drawImage(image, e.x, e.y);
        if (!drawContext) return;
        drawContext.clearRect(e.x, e.y, e.data.width, e.data.height);
    };

    const onInit = async (data: ImageData) => {
        width = data.width;
        height = data.height;
        await tick();   // wait for DOM changes to be applied
        initCanvas();
        if (!viewContext) return;
        viewContext.putImageData(data, 0, 0);
        ready = true;
    };
</script>

{#if !ready}
<p class="loading" transition:fade>Loading...</p>
{/if}
<canvas id="drawCanvas"
        bind:this={drawCanvas}
        {width}
        {height}
        style="width: {width}px; height: {height}px;"
        onclick={handleClick}
        onmousemove={handleMove}
        onmouseleave={handleEnd}
        ontouchstart={handleTouchStart}
        ontouchmove={handleTouch}
        ontouchend={handleEnd}
></canvas>
<canvas
        id="mainCanvas"
        bind:this={viewCanvas}
        {width}
        {height}
        style="width: {width}px; height: {height}px;"
></canvas>

<style>
    #drawCanvas {
        /*border: 1px solid blue;*/
        /*background: white;*/
        position: absolute;
        cursor: crosshair;
    }

    .loading {
        position: absolute;
        width: 400px;
        color: black;
        padding: 0 1rem;
        font-size: 1.8rem;
        text-align: center;
    }

    canvas {
        transition: width 0.3s ease-in-out, height 0.5s ease-in-out;
    }
</style>
