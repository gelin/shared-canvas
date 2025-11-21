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

function encodeImageData(data: ImageDataArray) {
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
        } else if (r < 128 || b < 128 || g < 128) {
            // black pixel
            result += '0';
        } else {
            // white pixel
            result += '1';
        }
    }
    return result;
}
