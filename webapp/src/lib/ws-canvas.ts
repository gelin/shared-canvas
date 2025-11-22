import type { WSDrawMessage } from "$lib/ws";

export type DrawEvent = {
    x: number,
    y: number,
    data: ImageData
}

export const eventToMessage = (e: DrawEvent): WSDrawMessage => {
    return {
        method: 'draw',
        params: {
            x: e.x, y: e.y,
            w: e.data.width, h: e.data.height,
            p: encodeImageData(e.data.data),
        }
    }
}

const encodeImageData = (data: ImageDataArray) => {
    // every 4 numbers represent a pixel: rbga
    let result = '';
    for (let i = 0; i < data.length; i += 4) {
        const r = data[i];
        const g = data[i + 1];
        const b = data[i + 2];
        const a = data[i + 3];
        if (a < 128) {
            // transparent pixel
            result += '_';  // ¯\_(ツ)_/¯
        } else if ((r + g + b) / 3 < 128) {
            // black pixel
            result += '0';
        } else {
            // white pixel
            result += '1';
        }
    }
    return result;
}

export const messageToEvent = (m: WSDrawMessage): DrawEvent => {
    const dataArray = decodeImageData(m.params.p);
    const imageData = new ImageData(dataArray, m.params.w, m.params.h);
    return {
        x: m.params.x,
        y: m.params.y,
        data: imageData,
    }
}

const decodeImageData = (pixels: string): ImageDataArray => {
    // every pixed is represented by 4 numbers: rbga
    let result = new Uint8ClampedArray(pixels.length * 4);
    for (let i = 0; i < pixels.length; i += 1) {
        const c = pixels[i];
        switch (c) {
            case '_':
                // transparent pixel
                result[i * 4] = 0;
                result[i * 4 + 1] = 0;
                result[i * 4 + 2] = 0;
                result[i * 4 + 3] = 0;
                break;
            case '0':
                // black pixel
                result[i * 4] = 0;
                result[i * 4 + 1] = 0;
                result[i * 4 + 2] = 0;
                result[i * 4 + 3] = 255;
                break;
            case '1':
                // white pixel
                result[i * 4] = 255;
                result[i * 4 + 1] = 255;
                result[i * 4 + 2] = 255;
                result[i * 4 + 3] = 255;
                break;
        }
    }
    return result;
}
