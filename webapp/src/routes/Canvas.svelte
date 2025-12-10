<script lang="ts">
    import { onMount, tick } from 'svelte';
    import { fade } from 'svelte/transition';
    import { type DrawEvent, eventToMessage, messageToEvent } from "$lib/ws-canvas";
    import { wsClient, type WSDrawMessage } from "$lib/ws";
    import { tool, STAMP_HALF_SIZE, STAMP_SIZE } from "./Palette";

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

    let stampImage: Image | null;

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
                onDrawMessage(messageToEvent(message));
            } else if (message.method === 'init') {
                onInitMessage(messageToEvent(message));
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
        initContext();
    };

    const initContext = () => {
        if (!drawContext) return;
        drawContext.strokeStyle = $tool.color;
        drawContext.lineWidth = $tool.size;
        drawContext.lineCap = 'round';
    };

    tool.subscribe((t) => {
        if (t.type === 'stamp') {
            const img = new Image();
            img.src = t.stampUrl;
            stampImage = img;
        } else {
            stampImage = null;
        }
    });

    const handleClick = (e: MouseEvent) => {
        if (!drawContext) return;
        const coords = getMouseCanvasCoords(e);
        if (!coords) return;

        initContext();

        switch ($tool.type) {
            case 'line':
                drawContext.beginPath();
                drawContext.moveTo(coords.x, coords.y);
                drawContext.lineTo(coords.x, coords.y);
                drawContext.stroke();
                sendDrawLine(coords.x, coords.y, coords.x, coords.y);
            break;
            case 'stamp':
                drawContext.drawImage(stampImage, coords.x - STAMP_HALF_SIZE, coords.y - STAMP_HALF_SIZE);
                sendDrawStamp(coords.x, coords.y);
                break;
        }

        prev = coords;
    };

    const handleMove = (e: MouseEvent & { buttons: number }) => {
        if (!drawContext) return;
        if ($tool.type !== 'line') return;
        const coords = getMouseCanvasCoords(e);
        if (!coords) return;

        if (e.buttons == 1) {
            if (isDrawing) {
                initContext();
                drawContext.beginPath();
                drawContext.moveTo(prev.x, prev.y);
                drawContext.lineTo(coords.x, coords.y);
                drawContext.stroke();

                sendDrawLine(prev.x, prev.y, coords.x, coords.y);
            } else {
                isDrawing = true;
            }
            prev = coords;
        } else {
            isDrawing = false;
        }
    };

    const getMouseCanvasCoords = (e: MouseEvent): { x: number; y: number } | null => {
        if (!drawCanvas) return null;
        const rect = drawCanvas.getBoundingClientRect();
        const scaleX = drawCanvas.width / rect.width;
        const scaleY = drawCanvas.height / rect.height;
        return {
            x: (e.clientX - rect.left) * scaleX,
            y: (e.clientY - rect.top) * scaleY,
        };
    };

    const handleTouchStart = (e: TouchEvent) => {
        const coords = getTouchCanvasCoords(e);
        if (!coords) return;

        prev = coords;
        isDrawing = true;
        if ($tool.type === 'stamp') {
            if (!drawContext) return;
            initContext();
            drawContext.drawImage(stampImage, coords.x - STAMP_HALF_SIZE, coords.y - STAMP_HALF_SIZE);
            sendDrawStamp(coords.x, coords.y);
        }
    };

    const handleTouch = (e: TouchEvent) => {
        const coords = getTouchCanvasCoords(e);
        if (!coords) return;

        if ($tool.type === 'line') {
            if (!drawContext) return null;
            initContext();
            drawContext.beginPath();
            drawContext.moveTo(prev.x, prev.y);
            drawContext.lineTo(coords.x, coords.y);
            drawContext.stroke();
            sendDrawLine(prev.x, prev.y, coords.x, coords.y);
        }

        prev = coords;
    };

    const getTouchCanvasCoords = ({ touches }: TouchEvent): { x: number; y: number } | null => {
        if (touches.length === 0) return null;
        if (!drawCanvas) return null;
        const rect = drawCanvas.getBoundingClientRect();
        const scaleX = drawCanvas.width / rect.width;
        const scaleY = drawCanvas.height / rect.height;
        return {
            x: (touches[0].clientX - rect.left) * scaleX,
            y: (touches[0].clientY - rect.top) * scaleY,
        };
    };

    const handleEnd = () => {
        isDrawing = false;
    };

    const sendDrawLine = (x0: number, y0: number, x1: number, y1: number) => {
        if (!drawContext) return;
        const lineWidth = drawContext.lineWidth;
        const drawX = Math.floor(Math.min(x0, x1) - lineWidth);
        const drawY = Math.floor(Math.min(y0, y1) - lineWidth);
        const drawWidth = Math.ceil(Math.abs(x0 - x1) + 2 * lineWidth);
        const drawHeight = Math.ceil(Math.abs(y0 - y1) + 2 * lineWidth);
        const imageData = drawContext.getImageData(drawX, drawY, drawWidth, drawHeight);
        const message = eventToMessage({ x: drawX, y: drawY, data: imageData });
        socket.send(message);
    };

    const sendDrawStamp = (x: number, y: number) => {
        if (!drawContext) return;
        const drawX = Math.floor(x - STAMP_HALF_SIZE);
        const drawY = Math.floor(y - STAMP_HALF_SIZE);
        const drawWidth = STAMP_SIZE;
        const drawHeight = STAMP_SIZE;
        const imageData = drawContext.getImageData(drawX, drawY, drawWidth, drawHeight);
        const message = eventToMessage({ x: drawX, y: drawY, data: imageData });
        socket.send(message);
    };

    const onDrawMessage = async (e: DrawEvent) => {
        if (!viewContext) return;
        // need to convert to image for transparency and composition to work
        const image = await createImageBitmap(e.data);
        viewContext.drawImage(image, e.x, e.y);
        if (!drawContext) return;
        drawContext.clearRect(e.x, e.y, e.data.width, e.data.height);
    };

    const onInitMessage = async (e: DrawEvent) => {
        width = e.data.width;
        height = e.data.height;
        await tick();   // wait for DOM changes to be applied
        initCanvas();
        if (!viewContext) return;
        viewContext.putImageData(e.data, 0, 0);
        ready = true;
    };
</script>

{#if !ready}
<p class="loading" transition:fade>Loading...</p>
{/if}
<div class="canvas-wrap" style="max-width: {width}px;">
    <canvas id="drawCanvas"
            bind:this={drawCanvas}
            {width}
            {height}
            style="cursor: {$tool.type === 'stamp' ? `url(${$tool.stampUrl}) ${STAMP_HALF_SIZE} ${STAMP_HALF_SIZE},` : ''} crosshair;"
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
    ></canvas>
</div>

<style>
    .canvas-wrap {
        position: relative;
        width: 100%;
        /* prevent accidental horizontal scroll around the canvas */
        overflow: hidden;
    }

    #drawCanvas {
        /*border: 1px solid blue;*/
        /*background: white;*/
        position: absolute;
        top: 0;
        left: 0;
        cursor: crosshair;
    }

    .loading {
        position: absolute;
        width: 100%;
        color: black;
        padding: 0 1rem;
        font-size: 1.8rem;
        text-align: center;
    }

    canvas {
        display: block;
        width: 100%;
        height: auto;
        transition: width 0.3s ease-in-out, height 0.5s ease-in-out;
    }
</style>
