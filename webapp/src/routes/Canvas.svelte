<script lang="ts">
    import { onMount } from 'svelte';

    onMount(async () => {
        const canvas = document.getElementById('canvas') as HTMLCanvasElement;
        const context = canvas.getContext('2d');
        if (!context) return;
        context.fillStyle = 'white';
        context.fillRect(0, 0, canvas.width, canvas.height);
    });

    function drawCanvas(event: MouseEvent) {
        let color = 0;
        switch (event.buttons) {
            case 1:
                color = 0;
                break;
            case 2:
                color = 255;
                break;
            case 0:
            default:
                return;
        }
        // console.log('drawCanvas', event);
        const canvas = event.target as HTMLCanvasElement;
        const context = canvas.getContext('2d');
        const x = event.offsetX;
        const y = event.offsetY;
        const imageData = context?.getImageData(x, y, 1, 1);
        const data = imageData?.data;
        if (!data) return;
        data[0] = color; // red
        data[1] = color; // green
        data[2] = color; // blue
        context?.putImageData(imageData, x, y);
    }
</script>

<canvas id="canvas" width="384" height="384" onmousemove={drawCanvas}></canvas>

