import { writable } from "svelte/store";

export type Type = 'line' | 'stamp';

export const COLORS = [ 'black', 'white' ];
export type Color = typeof COLORS[number];

export const SIZES = [ 1, 3, 5, 7, 10, 15 ];

export const STAMP_SIZE = 31;
export const STAMP_HALF_SIZE = 16;
export const STAMP_FONT_SIZE = 35;

export const STAMPS = [ 'star', 'heart' ];
export type Stamp = typeof STAMPS[number];
const STAMP_GLYPHS : Record<Stamp, string> = {
    star: '★',
    heart: '♥',
}

export class PaletteTool {
    readonly #type: Type;
    readonly #color: Color;
    readonly #size: number;
    readonly #stamp: Stamp | null;

    constructor(type: Type, color: Color, size: number, stamp: Stamp | null = null) {
        this.#type = type;
        this.#color = color;
        this.#size = size;
        this.#stamp = stamp;
    }

    get type() {
        return this.#type;
    }
    get color() {
        return this.#color;
    }
    get size() {
        return this.#size;
    }
    get stamp() {
        return this.#stamp;
    }

    get stampUrl() {
        if (this.#type !== 'stamp') return null;
        return stampUrl(this.#color, this.#stamp);
    }

    toJSON() {
        return {
            type: this.#type,
            color: this.#color,
            size: this.#size,
            stamp: this.#stamp,
        };
    }
}

export const stampUrl = (color: Color, stamp: Stamp | null): string | null => {
    if (!stamp) return null;
    const glyph = STAMP_GLYPHS[stamp];
    if (!glyph) return null;
    const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="${STAMP_SIZE}" height="${STAMP_SIZE}" viewBox="0 0 ${STAMP_SIZE} ${STAMP_SIZE}"><text x="50%" y="50%" dominant-baseline="central" text-anchor="middle" font-family="sans-serif" font-size="${STAMP_FONT_SIZE}" fill="${color}">${glyph}</text></svg>`;
    return `data:image/svg+xml;utf8,${svg}`;
};

export const DEFAULT_TOOL: PaletteTool = new PaletteTool('line', 'black', 3);

const LS_TOOL_KEY = 'palette.tool';

const loadTool = () => {
    const savedTool = JSON.parse(localStorage.getItem(LS_TOOL_KEY) || '{}');

    const parsedType = savedTool?.type as PaletteTool['type'] | null;
    const validType = parsedType === 'line' || parsedType === 'stamp' ? parsedType : null;
    const parsedColor = savedTool?.color as PaletteTool['color'] | null;
    const validColor = savedTool?.color === 'black' || savedTool?.color === 'white' ? parsedColor : null;
    const parsedSize = savedTool?.size ? parseInt(savedTool?.size, 10) : NaN;
    const validSize = SIZES.includes(parsedSize) || parsedSize === STAMP_SIZE ? parsedSize : null;
    const parsedStamp = savedTool?.stamp as PaletteTool['stamp'] | null;
    const validStamp = parsedType === 'stamp' && STAMPS.includes(parsedStamp ?? '') ? parsedStamp : null;

    if (validType && validColor && validSize !== null) {
        switch (validType) {
            case 'line': return new PaletteTool('line', validColor, validSize);
            case 'stamp': return new PaletteTool('stamp', validColor, STAMP_SIZE, validStamp);
        }
    }
}

export const tool = writable(loadTool() || DEFAULT_TOOL)

tool.subscribe((t) => {
    try {
        localStorage.setItem(LS_TOOL_KEY, JSON.stringify(t.toJSON()));
    } catch (_) {
        // ignore storage errors (e.g., privacy mode)
    }
});
