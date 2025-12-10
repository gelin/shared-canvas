import { writable } from "svelte/store";

export const COLORS = [ 'black', 'white' ];

export const SIZES = [ 1, 3, 5, 7, 10, 15 ];

export const STAMP_SIZE = 31;
export const STAMP_HALF_SIZE = 16;
export const STAMPS = [ 'star' ];

export class PaletteTool {
    readonly #type: 'line' | 'stamp';
    readonly #color: 'black' | 'white';
    readonly #size: number;
    readonly #stamp: 'star' | null;

    constructor(type: 'line' | 'stamp', color: 'black' | 'white', size: number, stamp: 'star' | null = null) {
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
        return `/stamps/${this.#color}-${this.#stamp}.svg`;
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

export const stampUrl = (color: string, stamp: string): string => {
    return `/stamps/${color}-${stamp}.svg`;
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
        // TODO fix this
        localStorage.setItem(LS_TOOL_KEY, JSON.stringify(t.toJSON()));
    } catch (_) {
        // ignore storage errors (e.g., privacy mode)
    }
});
