export type PaletteTool = {
    type: 'line' | 'stamp';
    color: 'white' | 'black';
    size: number;
    stamp: 'star' | null;
}

export const DEFAULT_TOOL: PaletteTool = { type: 'line', color: 'black', size: 3, stamp: null };

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

