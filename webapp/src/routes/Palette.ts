export type PaletteTool = {
    type: 'line' | 'stamp';
    color: 'black' | 'white';
    size: number;
    stamp: 'star' | null;
}

export const DEFAULT_TOOL: PaletteTool = { type: 'line', color: 'black', size: 3, stamp: null };

export const STAMP_SIZE = 31;

export const stampUrl = (color: string, stamp: string) => {
    return `/stamps/${color}-${stamp}.svg`;
};

export class PaletteChangeEvent extends Event {
    readonly #tool: PaletteTool;

    constructor(tool: PaletteTool) {
        super("palettechange");
        this.#tool = tool;
    }

    get tool() {
        return this.#tool;
    }
}

